package protocol

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/openclaw/openclaw-node/internal/crypto"
)

type Client struct {
	conn           *websocket.Conn
	serverURL      string
	opts           ConnectOptions
	identity       *Identity
	cryptoIdentity *crypto.Identity // For signing

	pending   map[string]chan *Frame
	pendingMu sync.RWMutex

	pendingChallengeNonce string // Captured challenge nonce for handshake

	OnConnected    func(ConnectResponse)
	OnDisconnected func(string)
	OnEvent        func(string, json.RawMessage)
	OnInvoke       func(InvokeRequest) *InvokeResult

	writeMu sync.Mutex
	closeMu sync.Mutex
	closed  bool
}

type ConnectOptions struct {
	Role        string
	Scopes      []string
	Caps        []string
	Commands    []string
	Permissions map[string]bool
	Token       string
	Client      ClientInfo
}

type ClientInfo struct {
	ID              string
	DisplayName     string
	Version         string
	Platform        string
	Mode            string
	InstanceID      string
	DeviceFamily    string
	ModelIdentifier string
}

type Identity struct {
	DeviceID   string
	ClientID   string
	ClientMode string
	Role       string
	SignedAtMs int64
	Nonce      string
}

func NewClient(endpoint string, identity *Identity, cryptoIdentity *crypto.Identity, opts ConnectOptions) *Client {
	return &Client{
		serverURL:      endpoint,
		identity:       identity,
		cryptoIdentity: cryptoIdentity,
		opts:           opts,
		pending:        make(map[string]chan *Frame),
	}
}

func (c *Client) Connect(ctx context.Context) error {
	return c.connectWithRetry(ctx, 0)
}

// connectWithRetry attempts to connect with exponential backoff
// initialDelay: starting delay in ms (0 = no backoff on first attempt)
func (c *Client) connectWithRetry(ctx context.Context, initialDelay int) error {
	c.closeMu.Lock()
	c.closed = false
	c.closeMu.Unlock()

	if initialDelay > 0 {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Duration(initialDelay) * time.Millisecond):
		}
	}

	dialer := websocket.Dialer{
		HandshakeTimeout: 12 * time.Second,
	}

	conn, _, err := dialer.DialContext(ctx, c.endpoint(), nil)
	if err != nil {
		return fmt.Errorf("dial: %w", err)
	}
	c.conn = conn

	// Set up a preliminary handler to capture challenge BEFORE starting readLoop
	// The gateway sends connect.challenge immediately when WS opens, so we need to be ready
	challengeCh := make(chan string, 1)
	preliminaryHandler := c.OnEvent
	c.OnEvent = func(event string, payload json.RawMessage) {
		if event == "connect.challenge" {
			var challenge struct {
				Nonce string `json:"nonce"`
			}
			json.Unmarshal(payload, &challenge) //nolint:errcheck
			if challenge.Nonce != "" {
				select {
				case challengeCh <- challenge.Nonce:
				default:
				}
			}
		}
		if preliminaryHandler != nil {
			preliminaryHandler(event, payload)
		}
	}

	go c.readLoop()
	go c.pingLoop()

	// Wait for challenge before proceeding
	select {
	case nonce := <-challengeCh:
		c.pendingChallengeNonce = nonce
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(10 * time.Second):
		return fmt.Errorf("challenge timeout")
	}

	return c.handshake(ctx)
}

func (c *Client) endpoint() string {
	// c.serverURL is like "http://localhost:38789" or "localhost:38789"
	// Convert to WebSocket URL
	serverURL := c.serverURL
	if serverURL == "" {
		return "ws://localhost:18789"
	}

	// Add scheme if missing
	if len(serverURL) >= 4 && serverURL[:4] != "http" && serverURL[:4] != "ws:/" {
		if c.opts.Client.Platform == "https" || c.opts.Client.Platform == "wss" {
			serverURL = "wss://" + serverURL
		} else {
			serverURL = "ws://" + serverURL
		}
	} else if len(serverURL) >= 4 && serverURL[:4] == "http" {
		// Replace http with ws, https with wss
		if len(serverURL) >= 5 && serverURL[:5] == "https" {
			serverURL = "wss" + serverURL[5:]
		} else {
			serverURL = "ws" + serverURL[4:]
		}
	}

	return serverURL
}

func (c *Client) IsConnected() bool {
	return c.conn != nil && !c.isClosed()
}

func (c *Client) SetServerURL(serverURL string) {
	c.serverURL = serverURL
}

func (c *Client) SetToken(token string) {
	c.opts.Token = token
}

func (c *Client) Disconnect() {
	c.closeMu.Lock()
	defer c.closeMu.Unlock()
	if c.closed {
		return
	}
	c.closed = true

	if c.conn != nil {
		c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(1000, "bye")) //nolint:errcheck
		c.conn.Close()
		c.conn = nil
	}
}

func (c *Client) readLoop() {
	for {
		_, data, err := c.conn.ReadMessage()
		if err != nil {
			if !c.isClosed() && c.OnDisconnected != nil {
				c.OnDisconnected(err.Error())
			}
			return
		}

		var frame Frame
		if err := json.Unmarshal(data, &frame); err != nil {
			log.Printf("Received non-frame websocket message: %s", string(data))
			continue
		}
		logFrame("recv", &frame)
		c.handleFrame(&frame)
	}
}

func (c *Client) pingLoop() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if c.isClosed() {
			return
		}
		c.conn.SetWriteDeadline(time.Now().Add(60 * time.Second)) //nolint:errcheck
		if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
			if c.OnDisconnected != nil {
				c.OnDisconnected(err.Error())
			}
			return
		}
	}
}

func (c *Client) handleFrame(frame *Frame) {
	switch frame.Type {
	case "res":
		c.pendingMu.RLock()
		ch, ok := c.pending[frame.ID]
		c.pendingMu.RUnlock()
		if ok {
			select {
			case ch <- frame:
			default:
			}
		}
	case "event":
		if frame.Event == "node.invoke.request" {
			log.Printf("Handling node.invoke.request payload=%s", string(frame.Payload))
			var req InvokeRequest
			if err := json.Unmarshal(frame.Payload, &req); err == nil {
				log.Printf("Received node.invoke.request: id=%s nodeId=%s command=%s", req.ID, req.NodeID, req.Command)
				if c.OnInvoke != nil {
					if result := c.OnInvoke(req); result != nil {
						c.sendInvokeResult(&req, result)
					}
				}
			}
		}
		if c.OnEvent != nil {
			c.OnEvent(frame.Event, frame.Payload)
		}
	}
}

func (c *Client) sendInvokeResult(req *InvokeRequest, result *InvokeResult) {
	if c.conn == nil || req == nil || result == nil {
		return
	}

	params := map[string]interface{}{
		"id":     req.ID,
		"nodeId": req.NodeID,
		"ok":     result.OK,
	}
	if len(result.Payload) > 0 {
		var payload interface{}
		if err := json.Unmarshal(result.Payload, &payload); err == nil {
			params["payload"] = payload
		} else {
			params["payloadJSON"] = string(result.Payload)
		}
	}
	if result.Error != nil {
		params["error"] = map[string]interface{}{
			"code":    result.Error.Code,
			"message": result.Error.Message,
		}
	}

	reqFrame := Frame{
		Type:   "req",
		ID:     uuid.New().String(),
		Method: "node.invoke.result",
	}
	paramsJSON, err := json.Marshal(params)
	if err != nil {
		return
	}
	reqFrame.Params = paramsJSON

	c.writeMu.Lock()
	defer c.writeMu.Unlock()
	if err := c.conn.WriteJSON(reqFrame); err != nil {
		log.Printf("Failed to send node.invoke.result: id=%s nodeId=%s command=%s err=%v", req.ID, req.NodeID, req.Command, err)
		return
	}
	logFrame("send", &reqFrame)
	if result.OK {
		log.Printf("Sent node.invoke.result: id=%s nodeId=%s command=%s ok=true", req.ID, req.NodeID, req.Command)
	} else if result.Error != nil {
		log.Printf("Sent node.invoke.result: id=%s nodeId=%s command=%s ok=false error=%s:%s", req.ID, req.NodeID, req.Command, result.Error.Code, result.Error.Message)
	} else {
		log.Printf("Sent node.invoke.result: id=%s nodeId=%s command=%s ok=false", req.ID, req.NodeID, req.Command)
	}
}

func (c *Client) isClosed() bool {
	c.closeMu.Lock()
	defer c.closeMu.Unlock()
	return c.closed
}

// handshake performs the protocol handshake with challenge-response
func (c *Client) handshake(ctx context.Context) error {
	// Use the challenge nonce that was captured in connectWithRetry
	if c.pendingChallengeNonce == "" {
		return fmt.Errorf("no challenge nonce captured")
	}

	// Send connect with signed device info using the challenge nonce
	return c.sendSignedConnect(ctx, c.pendingChallengeNonce)
}

// sendSignedConnect sends the connect request with signed device info
func (c *Client) sendSignedConnect(ctx context.Context, challengeNonce string) error {
	// cryptoIdentity is required for device authentication
	if c.cryptoIdentity == nil {
		return fmt.Errorf("cryptoIdentity is required for device authentication")
	}

	signedAtMs := time.Now().UnixMilli()

	params := c.buildConnectParams(signedAtMs, challengeNonce)

	reqID := uuid.New().String()
	req := Frame{
		Type:   "req",
		ID:     reqID,
		Method: "connect",
	}

	paramsJSON, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("marshal connect params: %w", err)
	}
	req.Params = paramsJSON

	// Register pending response before writing so a fast server response is not lost.
	respCh := make(chan *Frame, 1)
	c.pendingMu.Lock()
	c.pending[reqID] = respCh
	c.pendingMu.Unlock()

	c.writeMu.Lock()
	err = c.conn.WriteJSON(req)
	c.writeMu.Unlock()
	if err != nil {
		c.pendingMu.Lock()
		delete(c.pending, reqID)
		c.pendingMu.Unlock()
		return fmt.Errorf("send connect: %w", err)
	}
	logFrame("send", &req)

	defer func() {
		c.pendingMu.Lock()
		delete(c.pending, reqID)
		c.pendingMu.Unlock()
	}()

	select {
	case resp := <-respCh:
		if resp.OK {
			if c.OnConnected != nil {
				var connResp ConnectResponse
				if resp.Payload != nil {
					json.Unmarshal(resp.Payload, &connResp) //nolint:errcheck
				}
				c.OnConnected(connResp)
			}
			return nil
		} else {
			if resp.Error != nil {
				if requestID, ok := resp.Error.Details["requestId"].(string); ok && requestID != "" {
					return fmt.Errorf("connect rejected: %s - %s (requestId: %s)", resp.Error.Code, resp.Error.Message, requestID)
				}
				return fmt.Errorf("connect rejected: %s - %s", resp.Error.Code, resp.Error.Message)
			}
			return fmt.Errorf("connect failed")
		}
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(10 * time.Second):
		return fmt.Errorf("connect timeout")
	}
}

func (c *Client) buildConnectParams(signedAtMs int64, challengeNonce string) map[string]interface{} {
	client := map[string]interface{}{
		"id":         c.opts.Client.ID,
		"version":    c.opts.Client.Version,
		"platform":   c.opts.Client.Platform,
		"mode":       c.opts.Client.Mode,
		"instanceId": c.opts.Client.InstanceID,
	}
	if c.opts.Client.DisplayName != "" {
		client["displayName"] = c.opts.Client.DisplayName
	}
	if c.opts.Client.DeviceFamily != "" {
		client["deviceFamily"] = c.opts.Client.DeviceFamily
	}
	if c.opts.Client.ModelIdentifier != "" {
		client["modelIdentifier"] = c.opts.Client.ModelIdentifier
	}

	params := map[string]interface{}{
		"minProtocol": ProtocolVersion,
		"maxProtocol": ProtocolVersion,
		"client":      client,
		"role":        c.opts.Role,
	}
	if len(c.opts.Scopes) > 0 {
		params["scopes"] = c.opts.Scopes
	}
	if len(c.opts.Caps) > 0 {
		params["caps"] = c.opts.Caps
	}
	if len(c.opts.Commands) > 0 {
		params["commands"] = c.opts.Commands
	}
	if len(c.opts.Permissions) > 0 {
		params["permissions"] = c.opts.Permissions
	}

	if c.opts.Token != "" {
		params["auth"] = map[string]interface{}{
			"token": c.opts.Token,
		}
	}

	if c.cryptoIdentity == nil {
		// Device info will not be sent without cryptoIdentity
		return params
	}

	authPayload := BuildAuthPayload(
		c.cryptoIdentity.DeviceID,
		c.opts.Client.ID,
		c.opts.Client.Mode,
		c.opts.Role,
		c.opts.Scopes,
		signedAtMs,
		c.opts.Token,
		challengeNonce,
		c.opts.Client.Platform,
		c.opts.Client.DeviceFamily,
	)

	signature := c.cryptoIdentity.Sign(authPayload)
	log.Printf("[INFO] Device authentication: deviceId=%s, signatureLen=%d", c.cryptoIdentity.DeviceID, len(signature))

	params["device"] = map[string]interface{}{
		"id":        c.cryptoIdentity.DeviceID,
		"publicKey": c.cryptoIdentity.PublicKeyBase64(),
		"signature": base64.RawURLEncoding.EncodeToString(signature),
		"signedAt":  signedAtMs,
		"nonce":     challengeNonce,
	}

	return params
}

func logFrame(direction string, frame *Frame) {
	if frame == nil {
		return
	}
	switch frame.Type {
	case "event":
		log.Printf("WS %s frame: type=event event=%s id=%s", direction, frame.Event, frame.ID)
	case "req":
		log.Printf("WS %s frame: type=req method=%s id=%s", direction, frame.Method, frame.ID)
	case "res":
		log.Printf("WS %s frame: type=res id=%s ok=%t", direction, frame.ID, frame.OK)
	default:
		log.Printf("WS %s frame: type=%s id=%s", direction, frame.Type, frame.ID)
	}
}
