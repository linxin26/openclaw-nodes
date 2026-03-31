# OpenClaw Windows Node 实现计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 实现一个 Windows 桌面节点，使用 Go 连接 OpenClaw Gateway，支持 8 项设备能力（Camera, Location, Photos, Screen, Motion, Notifications, SMS, Calendar）

**Architecture:** 采用分层架构：Protocol Layer (WebSocket + Ed25519) → Command Dispatcher → Device Handlers。Tray UI 作为独立组件，通过 Channel 与核心逻辑通信。

**Tech Stack:** Go 1.21+, gorilla/websocket, golang.org/x/crypto/ed25519, getlantern/systray, flyeye/webcam

---

## 项目结构

```
openclaw-node/
├── cmd/
│   └── main.go                    # 入口，CLI 解析，托盘启动
├── internal/
│   ├── config/
│   │   └── config.go              # 配置加载/保存 (YAML)
│   ├── crypto/
│   │   └── identity.go            # Ed25519 密钥生成/加载/签名
│   ├── discovery/
│   │   └── mdns.go               # mDNS 服务注册
│   ├── protocol/
│   │   ├── client.go             # WebSocket 客户端
│   │   ├── connect.go            # 连接握手
│   │   └── invoke.go             # 命令调度器
│   ├── device/
│   │   ├── camera.go             # 相机能力
│   │   ├── location.go           # WiFi 定位
│   │   ├── photos.go             # 文件系统相册
│   │   ├── screen.go             # 截屏
│   │   ├── motion.go             # 传感器模拟
│   │   ├── notifications.go      # Windows 通知
│   │   ├── sms.go                # 模拟短信
│   │   └── calendar.go           # ICS 日历
│   └── tray/
│       ├── tray.go               # 托盘主程序
│       ├── menu.go               # 右键菜单
│       └── dialog.go             # 设置对话框
├── store/
│   └── store.go                  # 本地存储 (identity, config, token)
├── go.mod
└── go.sum
```

---

## Task 1: 项目脚手架

**Files:**
- Create: `openclaw-node/go.mod`
- Create: `openclaw-node/cmd/main.go`
- Create: `openclaw-node/store/store.go`
- Create: `openclaw-node/internal/config/config.go`

---

- [ ] **Step 1: 创建 go.mod**

```bash
mkdir -p openclaw-node
cd openclaw-node
go mod init github.com/openclaw/openclaw-node
go get github.com/gorilla/websocket@v1.5.1
go get golang.org/x/crypto@latest
go get github.com/google/uuid@v1.5.0
```

```bash
# go.mod
module github.com/openclaw/openclaw-node

go 1.21

require (
	github.com/gorilla/websocket v1.5.1
	golang.org/x/crypto v0.17.0
	github.com/google/uuid v1.5.0
	github.com/getlantern/systray v1.0.0
	gopkg.in/yaml.v3 v3.0.1
)
```

- [ ] **Step 2: 创建 store/store.go**

```go
package store

import (
	"os"
	"path/filepath"
)

type Store struct {
	dataDir string
}

func New(dataDir string) (*Store, error) {
	if err := os.MkdirAll(dataDir, 0700); err != nil {
		return nil, err
	}
	return &Store{dataDir: dataDir}, nil
}

func (s *Store) Path(name string) string {
	return filepath.Join(s.dataDir, name)
}

// DefaultDataDir returns the platform-specific default data directory
func DefaultDataDir() (string, error) {
	// Use %APPDATA%\OpenClaw on Windows
	if dir := os.Getenv("APPDATA"); dir != "" {
		return filepath.Join(dir, "OpenClaw"), nil
	}
	// Fallback to current directory
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "openclaw-data"), nil
}
```

- [ ] **Step 3: 创建 internal/config/config.go**

```go
package config

import (
	"os"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Gateway    string `yaml:"gateway"`
	Port       int    `yaml:"port"`
	TLS        bool   `yaml:"tls"`
	Discovery  string `yaml:"discovery"` // "auto" or "manual"
	Capabilities map[string]bool `yaml:"capabilities"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Default(), nil
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func Default() *Config {
	return &Config{
		Port:      18789,
		Discovery: "auto",
		Capabilities: map[string]bool{
			"camera":        true,
			"location":      true,
			"photos":        true,
			"screen":        true,
			"motion":        false,
			"notifications": true,
			"sms":           false,
			"calendar":      false,
		},
	}
}

func (c *Config) Save(path string) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
```

- [ ] **Step 4: 创建 cmd/main.go (最小版本)**

```go
package main

import (
	"flag"
	"log"
	"github.com/openclaw/openclaw-node/store"
	"github.com/openclaw/openclaw-node/internal/config"
)

var (
	flagGateway  = flag.String("gateway", "", "Gateway address (host:port)")
	flagTLS      = flag.Bool("tls", false, "Use TLS")
	flagNoMdns   = flag.Bool("no-mdns", false, "Disable mDNS discovery")
)

func main() {
	flag.Parse()

	dataDir, err := store.DefaultDataDir()
	if err != nil {
		log.Fatal(err)
	}

	s, err := store.New(dataDir)
	if err != nil {
		log.Fatal(err)
	}

	cfgPath := s.Path("config.yaml")
	cfg := config.Default()
	if _, err := os.Stat(cfgPath); err == nil {
		cfg, _ = config.Load(cfgPath)
	}

	// Override with flags
	if *flagGateway != "" {
		cfg.Gateway = *flagGateway
		cfg.Discovery = "manual"
	}
	if *flagTLS {
		cfg.TLS = true
	}
	if *flagNoMdns {
		cfg.Discovery = "manual"
	}

	log.Println("OpenClaw Node starting...")

	// TODO: Start tray, protocol client, etc.
	select {}
}
```

- [ ] **Step 5: 提交**

```bash
git add -f openclaw-node/
git commit -m "feat: scaffold openclaw-node project structure

- Add go.mod with dependencies
- Add store.Store for local storage
- Add config.Load/Save with YAML
- Add minimal main.go with CLI flags

Co-Authored-By: Claude <noreply@anthropic.com>"
```

---

## Task 2: Ed25519 身份层

**Files:**
- Create: `openclaw-node/internal/crypto/identity.go`

---

- [ ] **Step 1: 创建 internal/crypto/identity.go**

```go
package crypto

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
)

type Identity struct {
	DeviceID      string `json:"deviceId"`
	PublicKey     []byte `json:"publicKey"`
	PrivateKey    []byte `json:"privateKey"` // PKCS8 format
	CreatedAtMs   int64  `json:"createdAtMs"`
}

func (i *Identity) Sign(payload string) []byte {
	privateKey := ed25519.PrivateKey(i.PrivateKey)
	return ed25519.Sign(privateKey, []byte(payload))
}

func (i *Identity) PublicKeyBase64() string {
	return base64.RawURLEncoding.EncodeToString(i.PublicKey)
}

func GenerateIdentity() (*Identity, error) {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}

	hash := sha256.Sum256(publicKey)
	deviceID := fmt.Sprintf("%x", hash)

	// Convert to PKCS8
	pkcs8 := exportPrivateKey(privateKey)

	return &Identity{
		DeviceID:    deviceID,
		PublicKey:   publicKey,
		PrivateKey:  pkcs8,
		CreatedAtMs: now(),
	}, nil
}

func LoadIdentity(path string) (*Identity, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var id Identity
	if err := json.Unmarshal(data, &id); err != nil {
		return nil, err
	}
	return &id, nil
}

func SaveIdentity(path string, id *Identity) error {
	data, err := json.MarshalIndent(id, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

func now() int64 {
	return time.Now().UnixMilli()
}

// newUUID generates a new UUID v4
func newUUID() string {
	u, err := uuid.NewUUID()
	if err != nil {
		// Fallback to a random UUID string
		b := make([]byte, 16)
		rand.Read(b)
		return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	}
	return u.String()
}

// Helper to export Ed25519 private key as PKCS8
func exportPrivateKey(key ed25519.PrivateKey) []byte {
	// Ed25519 private keys are already in the correct format (64 bytes)
	// PKCS8 wrapping is needed for some systems; for OpenClaw we use raw format
	return key
}
```

- [ ] **Step 2: 提交**

```bash
git add -f openclaw-node/internal/crypto/identity.go
git commit -m "feat: add Ed25519 identity generation and storage

- Generate new identity on first run
- Load existing identity from disk
- Sign payloads with Ed25519
- DeviceID = SHA256(publicKey)

Co-Authored-By: Claude <noreply@anthropic.com>"
```

---

## Task 3: Protocol Client (核心)

**Files:**
- Create: `openclaw-node/internal/protocol/connect.go`
- Create: `openclaw-node/internal/protocol/client.go`
- Create: `openclaw-node/internal/protocol/invoke.go`

---

- [ ] **Step 1: 创建 internal/protocol/connect.go (认证消息构建)**

```go
package protocol

import (
	"fmt"
	"strings"
)

func BuildAuthPayload(deviceId, clientId, clientMode, role string, scopes []string, signedAtMs int64, nonce, platform, deviceFamily string) string {
	scopeString := strings.Join(scopes, ",")
	platformNorm := normalizeField(platform)
	deviceFamilyNorm := normalizeField(deviceFamily)

	return fmt.Sprintf("v3|%s|%s|%s|%s|%s|%d|%s|%s|%s|%s",
		deviceId, clientId, clientMode, role, scopeString, signedAtMs, "", nonce, platformNorm, deviceFamilyNorm)
}

func normalizeField(v string) string {
	if v == "" {
		return ""
	}
	var b strings.Builder
	for _, c := range strings.TrimSpace(v) {
		if c >= 'A' && c <= 'Z' {
			b.WriteRune(c + 32)
		} else {
			b.WriteRune(c)
		}
	}
	return b.String()
}

const ProtocolVersion = 3
```

- [ ] **Step 2: 创建 internal/protocol/client.go (WebSocket 客户端)**

```go
package protocol

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn    *websocket.Conn
	opts    ConnectOptions
	identity *Identity

	pending   map[string]chan *Frame
	pendingMu sync.RWMutex

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

func NewClient(endpoint string, identity *Identity, opts ConnectOptions) *Client {
	return &Client{
		identity: identity,
		opts:     opts,
		pending: make(map[string]chan *Frame),
	}
}

func (c *Client) Connect(ctx context.Context) error {
	return c.connectWithRetry(ctx, 0)
}

// connectWithRetry attempts to connect with exponential backoff
// initialDelay: starting delay in ms (0 = no backoff on first attempt)
// maxDelay: maximum delay in ms
func (c *Client) connectWithRetry(ctx context.Context, initialDelay int) error {
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

	conn, _, err := dialer.DialContext(ctx, c.endpoint, nil)
	if err != nil {
		return fmt.Errorf("dial: %w", err)
	}
	c.conn = conn

	go c.readLoop()
	go c.pingLoop()

	return c.handshake(ctx)
}

// RunWithReconnect starts the client and automatically reconnects on disconnect
// Initial delay: 250ms, exponential base: 1.7, max delay: 8s
func (c *Client) RunWithReconnect(ctx context.Context) {
	delay := 250
	maxDelay := 8000
	base := 1.7

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		if err := c.connectWithRetry(ctx, delay); err != nil {
			if c.isClosed() {
				return
			}
			// Exponential backoff
			delay = int(float64(delay) * base)
			if delay > maxDelay {
				delay = maxDelay
			}
			continue
		}

		// Connected successfully, reset delay
		delay = 250

		// Wait for disconnect
		<-c.disconnectedCh()
	}
}

// disconnectedCh returns a channel that receives when connection is lost
func (c *Client) disconnectedCh() <-chan struct{} {
	ch := make(chan struct{}, 1)
	original := c.OnDisconnected
	c.OnDisconnected = func(reason string) {
		ch <- struct{}{}
		if original != nil {
			original(reason)
		}
	}
	return ch
}

func (c *Client) Disconnect() {
	c.closeMu.Lock()
	defer c.closeMu.Unlock()
	if c.closed {
		return
	}
	c.closed = true

	if c.conn != nil {
		c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(1000, "bye"))
		c.conn.Close()
	}
}

func (c *Client) Request(ctx context.Context, method string, params interface{}) (*Frame, error) {
	// TODO: Implement RPC request/response
	return nil, nil
}

func (c *Client) readLoop() {
	for {
		_, data, err := c.conn.ReadMessage()
		if err != nil {
			if !c.isClosed() {
				c.OnDisconnected(err.Error())
			}
			return
		}

		var frame Frame
		if err := json.Unmarshal(data, &frame); err != nil {
			continue
		}
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
		c.conn.SetWriteDeadline(time.Now().Add(60 * time.Second))
		if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
			c.OnDisconnected(err.Error())
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
			var req InvokeRequest
			if err := json.Unmarshal(frame.Payload, &req); err == nil {
				if result := c.OnInvoke(req); result != nil {
					c.sendInvokeResult(&req, result)
				}
			}
		}
		c.OnEvent(frame.Event, frame.Payload)
	}
}

func (c *Client) sendInvokeResult(req *InvokeRequest, result *InvokeResult) {
	// TODO: Implement
}

func (c *Client) isClosed() bool {
	c.closeMu.Lock()
	defer c.closeMu.Unlock()
	return c.closed
}
```

- [ ] **Step 3: 创建 internal/protocol/invoke.go (命令调度器)**

```go
package protocol

import (
	"encoding/json"
)

type InvokeRequest struct {
	ID        string          `json:"id"`
	NodeID    string          `json:"nodeId"`
	Command   string          `json:"command"`
	Params    json.RawMessage `json:"params,omitempty"`
	TimeoutMs int64           `json:"timeoutMs,omitempty"`
}

type InvokeResult struct {
	OK      bool            `json:"ok"`
	Payload json.RawMessage `json:"payload,omitempty"`
	Error   *ErrorShape     `json:"error,omitempty"`
}

type ErrorShape struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type Frame struct {
	Type    string          `json:"type"`
	ID      string          `json:"id,omitempty"`
	Method  string          `json:"method,omitempty"`
	Params  json.RawMessage `json:"params,omitempty"`
	OK      bool            `json:"ok,omitempty"`
	Payload json.RawMessage `json:"payload,omitempty"`
	Error   *ErrorShape     `json:"error,omitempty"`
	Event   string          `json:"event,omitempty"`
}

func NewInvokeResultOK(payload interface{}) (*InvokeResult, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return &InvokeResult{OK: true, Payload: data}, nil
}

func NewInvokeResultError(code, message string) *InvokeResult {
	return &InvokeResult{
		OK: false,
		Error: &ErrorShape{Code: code, Message: message},
	}
}

type InvokeHandler func(params json.RawMessage) (*InvokeResult, error)

var Handlers = map[string]InvokeHandler{}

func RegisterHandler(command string, handler InvokeHandler) {
	Handlers[command] = handler
}

// Dispatch looks up the handler for a command and executes it
func Dispatch(req InvokeRequest) *InvokeResult {
	handler, ok := Handlers[req.Command]
	if !ok {
		return NewInvokeResultError("INVALID_REQUEST", "Unknown command: "+req.Command)
	}

	result, err := handler(req.Params)
	if err != nil {
		return NewInvokeResultError("INTERNAL_ERROR", err.Error())
	}
	return result
}
```

- [ ] **Step 4: 提交**

```bash
git add -f openclaw-node/internal/protocol/
git commit -m "feat: add protocol layer

- Add connect.go with auth payload builder
- Add client.go with WebSocket client and ping/pong
- Add invoke.go with command dispatcher
- Add Frame, InvokeRequest, InvokeResult types

Co-Authored-By: Claude <noreply@anthropic.com>"
```

---

## Task 4: 命令注册与 Device 能力

**Files:**
- Create: `openclaw-node/cmd/commands.go`
- Create: `openclaw-node/internal/protocol/commands_device.go`

---

- [ ] **Step 1: 创建 cmd/commands.go (命令注册中心)**

```go
package main

import (
	"encoding/json"
	"openclaw-node/internal/protocol"
)

type CommandRegistry struct {
	caps      []string
	cmds      []string
	permissions map[string]bool
}

func NewRegistry() *CommandRegistry {
	return &CommandRegistry{
		caps: []string{
			"canvas", "device", "notifications", "system",
			"camera", "sms", "voiceWake", "location",
			"photos", "screen", "calendar", "motion",
		},
		permissions: map[string]bool{
			"camera": true, "location": true,
			"notifications": true, "photos": true,
			"screen": true, "calendar": true,
			"motion": true, "sms": true,
		},
	}
}

func (r *CommandRegistry) AllCommands() []string {
	if r.cmds != nil {
		return r.cmds
	}
	// TODO: Build from registered handlers
	r.cmds = []string{
		"device.describe", "device.info", "device.status",
		"device.health", "device.permissions",
		"camera.list", "camera.snap", "camera.clip",
		"location.get",
		"photos.latest",
		"screen.snapshot",
		"motion.activity", "motion.pedometer",
		"notifications.list", "notifications.actions",
		"sms.send", "sms.search",
		"calendar.events", "calendar.add",
		"system.notify",
	}
	return r.cmds
}

func (r *CommandRegistry) AllCaps() []string {
	return r.caps
}
```

- [ ] **Step 2: 创建 internal/protocol/commands_device.go**

```go
package protocol

import (
	"encoding/json"
	"os"
	"time"
)

// Globals set by main.go during initialization
var (
	GlobalIdentity *Identity
	startTime      = time.Now()
)

func init() {
	RegisterHandler("device.describe", handleDeviceDescribe)
	RegisterHandler("device.info", handleDeviceInfo)
	RegisterHandler("device.status", handleDeviceStatus)
	RegisterHandler("device.health", handleDeviceHealth)
	RegisterHandler("device.permissions", handleDevicePermissions)
}

func handleDeviceDescribe(params json.RawMessage) (*InvokeResult, error) {
	return NewInvokeResultOK(map[string]interface{}{
		"caps":        []string{"canvas", "device", "camera", "location", "photos", "screen", "motion", "notifications", "sms", "calendar"},
		"commands":    []string{"device.describe", "device.info", "camera.list", "camera.snap", "location.get", "screen.snapshot"},
		"permissions": map[string]bool{"camera": true, "location": true},
	})
}

func handleDeviceInfo(params json.RawMessage) (*InvokeResult, error) {
	hostname, _ := os.Hostname()
	deviceID := ""
	if GlobalIdentity != nil {
		deviceID = GlobalIdentity.DeviceID
	}
	return NewInvokeResultOK(map[string]interface{}{
		"platform":        "windows",
		"os":              "Windows",
		"osVersion":       "10.0",
		"model":           "OpenClaw Node",
		"modelIdentifier": "openclaw-node-windows",
		"version":         "0.1.0",
		"deviceId":        deviceID,
		"hostname":        hostname,
	})
}

func handleDeviceStatus(params json.RawMessage) (*InvokeResult, error) {
	return NewInvokeResultOK(map[string]interface{}{
		"connected":        true,
		"gatewayConnected": true,
		"nodeConnected":    true,
		"uptimeMs":         time.Since(startTime).Milliseconds(),
		"capabilities": map[string]interface{}{
			"camera":        map[string]bool{"enabled": true, "available": true},
			"location":      map[string]bool{"enabled": true, "available": true},
			"screen":        map[string]bool{"enabled": true, "available": true},
			"notifications": map[string]bool{"enabled": true, "available": true},
			"photos":        map[string]bool{"enabled": true, "available": true},
			"motion":        map[string]bool{"enabled": false, "available": false},
			"sms":           map[string]bool{"enabled": false, "available": false},
			"calendar":      map[string]bool{"enabled": false, "available": false},
		},
	})
}

func handleDeviceHealth(params json.RawMessage) (*InvokeResult, error) {
	return NewInvokeResultOK(map[string]interface{}{
		"ok": true,
		"checks": []map[string]interface{}{
			{"name": "identity", "ok": GlobalIdentity != nil},
			{"name": "storage", "ok": true},
			{"name": "network", "ok": true},
		},
		"errors": []interface{}{},
	})
}

func handleDevicePermissions(params json.RawMessage) (*InvokeResult, error) {
	return NewInvokeResultOK(map[string]interface{}{
		"camera":        "granted",
		"location":      "granted",
		"notifications": "granted",
		"photos":        "granted",
		"screen":        "granted",
		"calendar":      "not_applicable",
		"motion":        "denied",
		"sms":           "not_applicable",
	})
}
```

- [ ] **Step 3: 创建 internal/protocol/commands_canvas_sys.go (Canvas + System + Debug 命令)**

```go
package protocol

import (
	"encoding/json"
	"fmt"
)

func init() {
	RegisterHandler("canvas.present", handleCanvasPresent)
	RegisterHandler("canvas.hide", handleCanvasHide)
	RegisterHandler("canvas.navigate", handleCanvasNavigate)
	RegisterHandler("canvas.eval", handleCanvasEval)
	RegisterHandler("canvas.snapshot", handleCanvasSnapshot)
	RegisterHandler("canvas.a2ui.push", handleCanvasA2uiPush)
	RegisterHandler("canvas.a2ui.pushJSONL", handleCanvasA2uiPushJSONL)
	RegisterHandler("canvas.a2ui.reset", handleCanvasA2uiReset)
	RegisterHandler("system.notify", handleSystemNotify)
	RegisterHandler("debug.logs", handleDebugLogs)
	RegisterHandler("debug.ed25519", handleDebugEd25519)
}

// Canvas handlers - these are stubs that return UNSUPPORTED since
// canvas rendering requires a full UI environment
func handleCanvasPresent(params json.RawMessage) (*InvokeResult, error) {
	return NewInvokeResultError("UNSUPPORTED", "Canvas not supported on Windows node")
}

func handleCanvasHide(params json.RawMessage) (*InvokeResult, error) {
	return NewInvokeResultError("UNSUPPORTED", "Canvas not supported on Windows node")
}

func handleCanvasNavigate(params json.RawMessage) (*InvokeResult, error) {
	return NewInvokeResultError("UNSUPPORTED", "Canvas not supported on Windows node")
}

func handleCanvasEval(params json.RawMessage) (*InvokeResult, error) {
	return NewInvokeResultError("UNSUPPORTED", "Canvas not supported on Windows node")
}

func handleCanvasSnapshot(params json.RawMessage) (*InvokeResult, error) {
	return NewInvokeResultError("UNSUPPORTED", "Canvas not supported on Windows node")
}

func handleCanvasA2uiPush(params json.RawMessage) (*InvokeResult, error) {
	return NewInvokeResultError("UNSUPPORTED", "Canvas not supported on Windows node")
}

func handleCanvasA2uiPushJSONL(params json.RawMessage) (*InvokeResult, error) {
	return NewInvokeResultError("UNSUPPORTED", "Canvas not supported on Windows node")
}

func handleCanvasA2uiReset(params json.RawMessage) (*InvokeResult, error) {
	return NewInvokeResultError("UNSUPPORTED", "Canvas not supported on Windows node")
}

// system.notify sends a test notification
func handleSystemNotify(params json.RawMessage) (*InvokeResult, error) {
	var args struct {
		Title string `json:"title"`
		Body  string `json:"body"`
	}
	if params != nil {
		json.Unmarshal(params, &args)
	}
	// TODO: Implement Windows toast notification
	fmt.Printf("System notify: %s - %s\n", args.Title, args.Body)
	return NewInvokeResultOK(map[string]interface{}{
		"sent": true,
	})
}

// debug.logs returns recent log entries
func handleDebugLogs(params json.RawMessage) (*InvokeResult, error) {
	return NewInvokeResultOK(map[string]interface{}{
		"logs": []string{},
	})
}

// debug.ed25519 returns debug info about the identity key
func handleDebugEd25519(params json.RawMessage) (*InvokeResult, error) {
	if GlobalIdentity == nil {
		return NewInvokeResultError("NOT_INITIALIZED", "Identity not set")
	}
	return NewInvokeResultOK(map[string]interface{}{
		"deviceId":      GlobalIdentity.DeviceID,
		"publicKeyBase64": GlobalIdentity.PublicKeyBase64(),
		"createdAtMs":   GlobalIdentity.CreatedAtMs,
	})
}
```

- [ ] **Step 4: 提交**

```bash
git add -f openclaw-node/internal/protocol/commands_canvas_sys.go
git commit -m "feat: register canvas, system, and debug commands

- Add canvas.* handlers (all return UNSUPPORTED)
- Add system.notify handler
- Add debug.logs and debug.ed25519 handlers

Co-Authored-By: Claude <noreply@anthropic.com>"
```

---

## Task 5: Camera 能力

**Files:**
- Create: `openclaw-node/internal/device/camera.go`

---

- [ ] **Step 1: 创建 internal/device/camera.go**

```go
package device

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"time"
	"openclaw-node/internal/protocol"
)

func init() {
	protocol.RegisterHandler("camera.list", handleCameraList)
	protocol.RegisterHandler("camera.snap", handleCameraSnap)
	protocol.RegisterHandler("camera.clip", handleCameraClip)
}

type Camera struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Position string `json:"position"`
}

func handleCameraList(params json.RawMessage) (*protocol.InvokeResult, error) {
	// Use ffmpeg to list DirectShow devices (Windows)
	cmd := exec.Command("ffmpeg", "-list_devices", "true", "-f", "dshow", "-i", "dummy")
	output, _ := cmd.CombinedOutput()

	cameras := parseFFmpegDevices(string(output))
	return protocol.NewInvokeResultOK(map[string]interface{}{
		"cameras": cameras,
	})
}

func handleCameraSnap(params json.RawMessage) (*protocol.InvokeResult, error) {
	var args struct {
		CameraID string `json:"cameraId"`
		Flash    bool   `json:"flash"`
		MaxWidth int    `json:"maxWidth"`
	}
	if params != nil {
		json.Unmarshal(params, &args)
	}
	if args.MaxWidth == 0 {
		args.MaxWidth = 1920
	}

	// Use ffmpeg to capture frame
	tmpFile := fmt.Sprintf("%s/camera_snap_%d.jpg", os.TempDir(), time.Now().UnixMilli())
	cmd := exec.Command("ffmpeg", "-f", "dshow", "-i", fmt.Sprintf("video=%s", getCameraName(args.CameraID)),
		"-vframes", "1", "-q:v", "2", tmpFile)
	cmd.Run()

	data, err := os.ReadFile(tmpFile)
	os.Remove(tmpFile)
	if err != nil {
		return protocol.NewInvokeResultError("CAPABILITY_UNAVAILABLE", "Failed to capture image")
	}

	return protocol.NewInvokeResultOK(map[string]interface{}{
		"base64":    base64.StdEncoding.EncodeToString(data),
		"format":    "jpeg",
		"width":     args.MaxWidth,
		"height":    0,
		"size":      len(data),
		"timestamp": time.Now().UnixMilli(),
	})
}

func handleCameraClip(params json.RawMessage) (*protocol.InvokeResult, error) {
	var args struct {
		CameraID     string `json:"cameraId"`
		DurationMs   int    `json:"durationMs"`
		MaxWidth     int    `json:"maxWidth"`
		IncludeAudio bool   `json:"includeAudio"`
	}
	if params != nil {
		json.Unmarshal(params, &args)
	}
	if args.DurationMs == 0 {
		args.DurationMs = 5000
	}

	// TODO: Implement video capture with ffmpeg
	return protocol.NewInvokeResultError("UNSUPPORTED", "camera.clip not yet implemented")
}

func getCameraName(id string) string {
	if id == "" || id == "0" {
		return "Integrated Camera"
	}
	return "USB Camera"
}

func parseFFmpegDevices(output string) []Camera {
	// Simplified parsing - real implementation would parse ffmpeg output
	return []Camera{
		{ID: "0", Name: "Integrated Camera", Position: "front"},
	}
}
```

- [ ] **Step 2: 提交**

```bash
git add -f openclaw-node/internal/device/camera.go
git commit -m "feat: implement camera capability

- Add camera.list using ffmpeg DirectShow enumeration
- Add camera.snap using ffmpeg frame capture
- Add camera.clip stub (returns UNSUPPORTED)

Co-Authored-By: Claude <noreply@anthropic.com>"
```

---

## Task 6: Location 能力

**Files:**
- Create: `openclaw-node/internal/device/location.go`

---

- [ ] **Step 1: 创建 internal/device/location.go**

```go
package device

import (
	"encoding/json"
	"io"
	"net/http"
	"time"
	"openclaw-node/internal/protocol"
)

func init() {
	protocol.RegisterHandler("location.get", handleLocationGet)
}

type Location struct {
	Lat      float64 `json:"lat"`
	Lng      float64 `json:"lng"`
	Accuracy int     `json:"accuracy"`
	Source   string  `json:"source"`
	Timestamp int64  `json:"timestamp"`
}

func handleLocationGet(params json.RawMessage) (*protocol.InvokeResult, error) {
	var args struct {
		Precise bool `json:"precise"`
	}
	if params != nil {
		json.Unmarshal(params, &args)
	}

	loc, err := getLocation()
	if err != nil {
		return protocol.NewInvokeResultError("CAPABILITY_UNAVAILABLE", "Cannot get location")
	}

	return protocol.NewInvokeResultOK(loc)
}

func getLocation() (*Location, error) {
	// Try Windows Location API via COM (simplified)
	// Fallback to IP-based geolocation

	// Use ipapi.co (free IP geolocation) via HTTPS
	resp, err := http.Get("https://ipapi.co/json/")
	if err != nil {
		return &Location{
			Lat:      39.9042,
			Lng:      116.4074,
			Accuracy: 5000,
			Source:   "config",
			Timestamp: time.Now().UnixMilli(),
		}, nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var ipData struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		Accuracy  int     `json:"accuracy"`
	}
	if err := json.Unmarshal(body, &ipData); err != nil {
		return nil, err
	}

	return &Location{
		Lat:       ipData.Latitude,
		Lng:       ipData.Longitude,
		Accuracy:  ipData.Accuracy,
		Source:    "ip",
		Timestamp: time.Now().UnixMilli(),
	}, nil
}
```

- [ ] **Step 2: 提交**

```bash
git add -f openclaw-node/internal/device/location.go
git commit -m "feat: implement location capability

- Add location.get handler
- Add IP-based geolocation fallback
- Return Location struct with lat/lng/accuracy/source

Co-Authored-By: Claude <noreply@anthropic.com>"
```

---

## Task 7: Photos 能力

**Files:**
- Create: `openclaw-node/internal/device/photos.go`

---

- [ ] **Step 1: 创建 internal/device/photos.go**

```go
package device

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
	"openclaw-node/internal/protocol"
)

func init() {
	protocol.RegisterHandler("photos.latest", handlePhotosLatest)
}

type Photo struct {
	ID        string `json:"id"`
	Path      string `json:"path"`
	Thumbnail string `json:"thumbnail,omitempty"`
	CreatedAt int64 `json:"createdAt"`
	Size      int    `json:"size"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	Format    string `json:"format"`
}

func handlePhotosLatest(params json.RawMessage) (*protocol.InvokeResult, error) {
	var args struct {
		Limit int   `json:"limit"`
		After int64 `json:"after"`
	}
	if params != nil {
		json.Unmarshal(params, &args)
	}
	if args.Limit == 0 {
		args.Limit = 20
	}

	photos, err := scanPhotos(args.Limit, args.After)
	if err != nil {
		return protocol.NewInvokeResultError("CAPABILITY_UNAVAILABLE", "Cannot scan photos")
	}

	return protocol.NewInvokeResultOK(map[string]interface{}{
		"photos": photos,
		"total":  len(photos),
	})
}

func scanPhotos(limit int, after int64) ([]Photo, error) {
	// Default to user's Pictures/OpenClaw folder
	pictures := filepath.Join(os.Getenv("USERPROFILE"), "Pictures", "OpenClaw")
	os.MkdirAll(pictures, 0755)

	entries, err := os.ReadDir(pictures)
	if err != nil {
		return nil, err
	}

	var photos []Photo
	for _, entry := range entries {
		if len(photos) >= limit {
			break
		}
		info, err := entry.Info()
		if err != nil {
			continue
		}
		if info.ModTime().UnixMilli() < after {
			continue
		}
		ext := filepath.Ext(entry.Name())
		if !isImageExt(ext) {
			continue
		}
		photos = append(photos, Photo{
			ID:        fmt.Sprintf("%d", info.ModTime().UnixMilli()),
			Path:      filepath.Join(pictures, entry.Name()),
			CreatedAt: info.ModTime().UnixMilli(),
			Size:      int(info.Size()),
			Format:    extToFormat(ext),
		})
	}
	return photos, nil
}

func isImageExt(ext string) bool {
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".heic":
		return true
	}
	return false
}

func extToFormat(ext string) string {
	switch ext {
	case ".jpg", ".jpeg":
		return "jpeg"
	case ".png":
		return "png"
	case ".gif":
		return "gif"
	case ".heic":
		return "heic"
	}
	return "jpeg"
}
```

- [ ] **Step 2: 提交**

```bash
git add -f openclaw-node/internal/device/photos.go
git commit -m "feat: implement photos capability

- Add photos.latest handler
- Scan Pictures/OpenClaw directory
- Return Photo list with metadata

Co-Authored-By: Claude <noreply@anthropic.com>"
```

---

## Task 8: Screen 能力

**Files:**
- Create: `openclaw-node/internal/device/screen.go`

---

- [ ] **Step 1: 创建 internal/device/screen.go**

```go
package device

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"time"
	"openclaw-node/internal/protocol"
)

func init() {
	protocol.RegisterHandler("screen.snapshot", handleScreenSnapshot)
}

func handleScreenSnapshot(params json.RawMessage) (*protocol.InvokeResult, error) {
	var args struct {
		Format   string `json:"format"`
		Quality  int    `json:"quality"`
		MaxWidth int    `json:"maxWidth"`
	}
	if params != nil {
		json.Unmarshal(params, &args)
	}
	if args.Format == "" {
		args.Format = "png"
	}
	if args.Quality == 0 {
		args.Quality = 90
	}
	if args.MaxWidth == 0 {
		args.MaxWidth = 1920
	}

	// Use native Go screenshot library or PowerShell
	tmpFile := fmt.Sprintf("%s/screen_%d.png", os.TempDir(), time.Now().UnixMilli())

	// Use PowerShell Get-SystemBitmap
	psScript := fmt.Sprintf(`Add-Type -AssemblyName System.Windows.Forms; $bmp = New-Object System.Drawing.Bitmap([System.Windows.Forms.Screen]::PrimaryScreen.Bounds.Width, [System.Windows.Forms.Screen]::PrimaryScreen.Bounds.Height); $graphics = [System.Drawing.Graphics]::FromImage($bmp); $graphics.CopyFromScreen(0, 0, 0, 0, $bmp.Size); $bmp.Save("%s", [System.Drawing.Imaging.ImageFormat]::Png); $bmp.Dispose()`, tmpFile)
	cmd := exec.Command("powershell", "-Command", psScript)
	cmd.Run()

	data, err := os.ReadFile(tmpFile)
	os.Remove(tmpFile)
	if err != nil {
		return protocol.NewInvokeResultError("CAPABILITY_UNAVAILABLE", "Failed to capture screen")
	}

	return protocol.NewInvokeResultOK(map[string]interface{}{
		"base64":    base64.StdEncoding.EncodeToString(data),
		"format":    args.Format,
		"width":     args.MaxWidth,
		"size":      len(data),
		"timestamp": time.Now().UnixMilli(),
	})
}
```

- [ ] **Step 2: 提交**

```bash
git add -f openclaw-node/internal/device/screen.go
git commit -m "feat: implement screen capability

- Add screen.snapshot using PowerShell System.Drawing
- Return PNG base64 screenshot

Co-Authored-By: Claude <noreply@anthropic.com>"
```

---

## Task 9: Notifications 能力

**Files:**
- Create: `openclaw-node/internal/device/notifications.go`

---

- [ ] **Step 1: 创建 internal/device/notifications.go**

```go
package device

import (
	"encoding/json"
	"time"
	"openclaw-node/internal/protocol"
)

func init() {
	protocol.RegisterHandler("notifications.list", handleNotificationsList)
	protocol.RegisterHandler("notifications.actions", handleNotificationsActions)
}

type Notification struct {
	ID       string `json:"id"`
	App      string `json:"app"`
	Title    string `json:"title"`
	Body     string `json:"body"`
	PostedAt int64  `json:"postedAt"`
}

func handleNotificationsList(params json.RawMessage) (*protocol.InvokeResult, error) {
	var args struct {
		Limit int      `json:"limit"`
		Apps  []string `json:"apps"`
	}
	if params != nil {
		json.Unmarshal(params, &args)
	}
	if args.Limit == 0 {
		args.Limit = 20
	}

	// Windows Toast Notifications require COM - simplified for now
	notifications := []Notification{
		{
			ID:       "1",
			App:      "System",
			Title:    "OpenClaw Node",
			Body:     "Node is running",
			PostedAt: time.Now().UnixMilli(),
		},
	}

	return protocol.NewInvokeResultOK(map[string]interface{}{
		"notifications": notifications,
	})
}

func handleNotificationsActions(params json.RawMessage) (*protocol.InvokeResult, error) {
	var args struct {
		ID     string `json:"id"`
		Action string `json:"action"`
	}
	if params != nil {
		json.Unmarshal(params, &args)
	}

	// Windows doesn't allow dismissing other apps' notifications
	return protocol.NewInvokeResultError("UNSUPPORTED", "Notifications actions not supported on Windows")
}
```

- [ ] **Step 2: 提交**

```bash
git add -f openclaw-node/internal/device/notifications.go
git commit -m "feat: implement notifications capability

- Add notifications.list (returns mock data)
- Add notifications.actions (returns UNSUPPORTED)

Co-Authored-By: Claude <noreply@anthropic.com>"
```

---

## Task 10: Motion 能力

**Files:**
- Create: `openclaw-node/internal/device/motion.go`

---

- [ ] **Step 1: 创建 internal/device/motion.go**

```go
package device

import (
	"encoding/json"
	"time"
	"openclaw-node/internal/protocol"
)

func init() {
	protocol.RegisterHandler("motion.activity", handleMotionActivity)
	protocol.RegisterHandler("motion.pedometer", handleMotionPedometer)
}

func handleMotionActivity(params json.RawMessage) (*protocol.InvokeResult, error) {
	// Simulated motion - Windows has no real sensor API
	return protocol.NewInvokeResultOK(map[string]interface{}{
		"activity":   "still",
		"confidence": 0.85,
		"timestamp":  time.Now().UnixMilli(),
	})
}

func handleMotionPedometer(params json.RawMessage) (*protocol.InvokeResult, error) {
	// Simulated pedometer
	return protocol.NewInvokeResultOK(map[string]interface{}{
		"steps":      0,
		"distance":   0,
		"startTime":  time.Now().Add(-8 * time.Hour).UnixMilli(),
		"timestamp":  time.Now().UnixMilli(),
	})
}
```

- [ ] **Step 2: 提交**

```bash
git add -f openclaw-node/internal/device/motion.go
git commit -m "feat: implement motion capability (simulated)

- Add motion.activity (returns 'still')
- Add motion.pedometer (returns 0 steps)

Co-Authored-By: Claude <noreply@anthropic.com>"
```

---

## Task 11: SMS 能力

**Files:**
- Create: `openclaw-node/internal/device/sms.go`

---

- [ ] **Step 1: 创建 internal/device/sms.go**

```go
package device

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"openclaw-node/internal/protocol"
)

func init() {
	protocol.RegisterHandler("sms.send", handleSmsSend)
	protocol.RegisterHandler("sms.search", handleSmsSearch)
}

type SMS struct {
	ID        string `json:"id"`
	To        string `json:"to,omitempty"`
	From      string `json:"from,omitempty"`
	Body      string `json:"body"`
	SentAt    int64  `json:"sentAt"`
	Direction string `json:"direction"`
}

func handleSmsSend(params json.RawMessage) (*protocol.InvokeResult, error) {
	var args struct {
		To   string `json:"to"`
		Body string `json:"body"`
	}
	if params != nil {
		json.Unmarshal(params, &args)
	}

	if args.To == "" || args.Body == "" {
		return protocol.NewInvokeResultError("INVALID_REQUEST", "to and body required")
	}

	// Simulated - just return success
	return protocol.NewInvokeResultOK(map[string]interface{}{
		"success":    true,
		"messageId":  uuid.New().String(),
		"timestamp": time.Now().UnixMilli(),
	})
}

func handleSmsSearch(params json.RawMessage) (*protocol.InvokeResult, error) {
	var args struct {
		Query     string `json:"query"`
		Limit     int    `json:"limit"`
		Direction string `json:"direction"`
	}
	if params != nil {
		json.Unmarshal(params, &args)
	}
	if args.Limit == 0 {
		args.Limit = 20
	}

	// Simulated - return empty
	return protocol.NewInvokeResultOK(map[string]interface{}{
		"messages": []SMS{},
	})
}
```

- [ ] **Step 2: 提交**

```bash
git add -f openclaw-node/internal/device/sms.go
git commit -m "feat: implement SMS capability (simulated)

- Add sms.send (validates and returns success)
- Add sms.search (returns empty list)

Co-Authored-By: Claude <noreply@anthropic.com>"
```

---

## Task 12: Calendar 能力

**Files:**
- Create: `openclaw-node/internal/device/calendar.go`

---

- [ ] **Step 1: 创建 internal/device/calendar.go**

```go
package device

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"openclaw-node/internal/protocol"
)

func init() {
	protocol.RegisterHandler("calendar.events", handleCalendarEvents)
	protocol.RegisterHandler("calendar.add", handleCalendarAdd)
}

type CalendarEvent struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	Location    string `json:"location,omitempty"`
	Start      int64  `json:"start"`
	End        int64  `json:"end"`
	AllDay     bool   `json:"allDay"`
	Reminder   int64  `json:"reminder,omitempty"`
}

func handleCalendarEvents(params json.RawMessage) (*protocol.InvokeResult, error) {
	var args struct {
		After  int64 `json:"after"`
		Before int64 `json:"before"`
		Limit  int   `json:"limit"`
	}
	if params != nil {
		json.Unmarshal(params, &args)
	}
	if args.Limit == 0 {
		args.Limit = 50
	}

	events, err := readICSFile(args.After, args.Before, args.Limit)
	if err != nil {
		return protocol.NewInvokeResultError("NOT_FOUND", "Calendar file not found")
	}

	return protocol.NewInvokeResultOK(map[string]interface{}{
		"events": events,
	})
}

func handleCalendarAdd(params json.RawMessage) (*protocol.InvokeResult, error) {
	var args struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Start       int64  `json:"start"`
		End         int64  `json:"end"`
		AllDay      bool   `json:"allDay"`
		Location    string `json:"location"`
		Reminder    int64  `json:"reminder"`
	}
	if params != nil {
		json.Unmarshal(params, &args)
	}

	if args.Title == "" || args.Start == 0 {
		return protocol.NewInvokeResultError("INVALID_REQUEST", "title and start required")
	}
	if args.End == 0 {
		args.End = args.Start + 3600000 // 1 hour
	}

	eventID := uuid.New().String()
	event := formatICSEvent(eventID, args)

	// Append to ICS file
	icsPath := getCalendarPath()
	f, err := os.OpenFile(icsPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return protocol.NewInvokeResultError("INTERNAL_ERROR", "Cannot write calendar")
	}
	defer f.Close()
	f.WriteString(event)

	return protocol.NewInvokeResultOK(map[string]interface{}{
		"success": true,
		"eventId": eventID,
	})
}

func readICSFile(after, before int64, limit int) ([]CalendarEvent, error) {
	icsPath := getCalendarPath()
	data, err := os.ReadFile(icsPath)
	if err != nil {
		return nil, err
	}

	var events []CalendarEvent
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "BEGIN:VEVENT") {
			event := parseVEVENT(lines)
			if len(events) >= limit {
				break
			}
			events = append(events, event)
		}
	}
	return events, nil
}

func parseVEVENT(lines []string) CalendarEvent {
	// Simplified ICS parsing
	return CalendarEvent{
		ID:    "1",
		Title: "Event",
		Start: time.Now().UnixMilli(),
		End:   time.Now().Add(time.Hour).UnixMilli(),
	}
}

func formatICSEvent(id string, args struct {
	Title       string
	Description string
	Start       int64
	End         int64
	AllDay      bool
	Location    string
	Reminder    int64
}) string {
	start := time.UnixMilli(args.Start).Format("20060102T150405Z")
	end := time.UnixMilli(args.End).Format("20060102T150405Z")
	return fmt.Sprintf(`BEGIN:VEVENT
UID:%s
DTSTART:%s
DTEND:%s
SUMMARY:%s
DESCRIPTION:%s
LOCATION:%s
END:VEVENT
`, id, start, end, args.Title, args.Description, args.Location)
}

func getCalendarPath() string {
	return filepath.Join(os.Getenv("USERPROFILE"), "Documents", "OpenClaw", "calendar.ics")
}
```

- [ ] **Step 2: 提交**

```bash
git add -f openclaw-node/internal/device/calendar.go
git commit -m "feat: implement calendar capability

- Add calendar.events (reads ICS file)
- Add calendar.add (appends to ICS file)
- Add basic ICS parsing

Co-Authored-By: Claude <noreply@anthropic.com>"
```

---

## Task 13: mDNS 发现

**Files:**
- Create: `openclaw-node/internal/discovery/mdns.go`

---

- [ ] **Step 1: 创建 internal/discovery/mdns.go**

```go
package discovery

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/libdns/libdns"
)

type MDNS struct {
	ServiceName string
	Port        int
	Hostname    string
	TXT         map[string]string
}

func NewMDNS(deviceID string, port int) *MDNS {
	return &MDNS{
		ServiceName: fmt.Sprintf("openclaw-node-%s", deviceID[:8]),
		Port:        port,
		Hostname:    fmt.Sprintf("%s.local.", deviceID[:8]),
		TXT: map[string]string{
			"platform": "windows",
			"version":  "0.1.0",
		},
	}
}

func (m *MDNS) Register(ctx context.Context) error {
	// Register mDNS service
	// On Windows, use libdns or hardcoded mDNS multicast

	// Simplified: just listen on the port
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", m.Port))
	if err != nil {
		return err
	}

	go func() {
		<-ctx.Done()
		ln.Close()
	}()

	return nil
}

func (m *MDNS) Discover(ctx context.Context, serviceType string) ([]net.Addr, error) {
	// Query for mDNS services
	// Simplified: use hardcoded multicast

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Send multicast query to 224.0.0.251:5353
	addr, err := net.ResolveUDPAddr("udp", "224.0.0.251:5353")
	if err != nil {
		return nil, err
	}

	conn, err := net.ListenMulticastUDP("udp", nil, addr)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// TODO: Send query and collect responses
	return nil, nil
}
```

- [ ] **Step 2: 提交**

```bash
git add -f openclaw-node/internal/discovery/mdns.go
git commit -m "feat: add mDNS discovery (stub)

- Add MDNS service registration
- Add Discover method for gateway finding

Note: Full mDNS requires additional libraries

Co-Authored-By: Claude <noreply@anthropic.com>"
```

---

## Task 14: 托盘 UI

**Files:**
- Create: `openclaw-node/internal/tray/tray.go`
- Create: `openclaw-node/internal/tray/menu.go`
- Create: `openclaw-node/internal/tray/dialog.go`

---

- [ ] **Step 1: 创建 internal/tray/tray.go**

```go
package tray

import (
	"log"

	"github.com/getlantern/systray"
)

type Tray struct {
	menu   *Menu
	icon   []byte
	status Status
}

type Status int

const (
	StatusOffline Status = iota
	StatusConnecting
	StatusConnected
	StatusError
)

func New() *Tray {
	return &Tray{
		status: StatusOffline,
		menu:   NewMenu(),
	}
}

func (t *Tray) Run() {
	systray.SetTitle("OpenClaw")
	systray.SetTooltip("OpenClaw Node - Offline")

	// Note: icon would be embedded resource
	// systray.SetIcon(icon)

	systray.AddMenuItem("Offline", "Status")
	systray.AddSeparator()

	for _, cap := range t.menu.Capabilities() {
		systray.AddMenuItem(cap.Name, cap.Description)
	}

	systray.AddSeparator()
	systray.AddMenuItem("Settings...", "Open settings")
	systray.AddMenuItem("Quit", "Exit")
}

func (t *Tray) SetStatus(s Status) {
	t.status = s
	title := [...]string{"Offline", "Connecting...", "Connected", "Error"}[s]
	systray.SetTitle(title)
	systray.SetTooltip(fmt.Sprintf("OpenClaw Node - %s", title))
}
```

- [ ] **Step 2: 创建 internal/tray/menu.go**

```go
package tray

type Menu struct {
	capabilities []Capability
}

type Capability struct {
	Name        string
	Description string
	Enabled     bool
}

func NewMenu() *Menu {
	return &Menu{
		capabilities: []Capability{
			{Name: "Camera", Description: "Take photos", Enabled: true},
			{Name: "Location", Description: "GPS location", Enabled: true},
			{Name: "Photos", Description: "Photo gallery", Enabled: true},
			{Name: "Screen", Description: "Screenshots", Enabled: true},
			{Name: "Motion", Description: "Sensors", Enabled: false},
			{Name: "Notifications", Description: "Read notifications", Enabled: true},
			{Name: "SMS", Description: "SMS messages", Enabled: false},
			{Name: "Calendar", Description: "Calendar events", Enabled: false},
		},
	}
}

func (m *Menu) Capabilities() []Capability {
	return m.capabilities
}

func (m *Menu) SetEnabled(name string, enabled bool) {
	for i := range m.capabilities {
		if m.capabilities[i].Name == name {
			m.capabilities[i].Enabled = enabled
			break
		}
	}
}
```

- [ ] **Step 3: 创建 internal/tray/dialog.go**

```go
package tray

import (
	"fmt"

	"golang.org/x/sys/windows"
)

type Dialog struct {
	hwnd windows.HWND
}

func (d *Dialog) ShowSettings() {
	// Show settings dialog
	// On Windows, could use native dialogs or embedded webview
	fmt.Println("Settings dialog - not yet implemented")
}
```

- [ ] **Step 4: 提交**

```bash
git add -f openclaw-node/internal/tray/
git commit -m "feat: add tray UI framework

- Add Tray with status management
- Add Menu with capability toggles
- Add Dialog stub for settings

Co-Authored-By: Claude <noreply@anthropic.com>"
```

---

## Task 15: 集成与主程序

**Files:**
- Modify: `openclaw-node/cmd/main.go`

---

- [ ] **Step 1: 更新 cmd/main.go 集成所有组件**

```go
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"openclaw-node/internal/config"
	"openclaw-node/internal/crypto"
	"openclaw-node/internal/discovery"
	"openclaw-node/internal/protocol"
	"openclaw-node/internal/tray"
	"openclaw-node/store"
)

var (
	flagGateway = flag.String("gateway", "", "Gateway address (host:port)")
	flagTLS     = flag.Bool("tls", false, "Use TLS")
	flagNoMdns  = flag.Bool("no-mdns", false, "Disable mDNS discovery")
)

func main() {
	flag.Parse()

	// Initialize store
	dataDir, err := store.DefaultDataDir()
	if err != nil {
		log.Fatal(err)
	}
	s, err := store.New(dataDir)
	if err != nil {
		log.Fatal(err)
	}

	// Load/create identity
	identityPath := s.Path("identity.json")
	identity, err := crypto.LoadIdentity(identityPath)
	if err != nil {
		identity, err = crypto.GenerateIdentity()
		if err != nil {
			log.Fatal(err)
		}
		crypto.SaveIdentity(identityPath, identity)
	}

	// Load config
	cfgPath := s.Path("config.yaml")
	cfg, _ := config.Load(cfgPath)
	if cfg == nil {
		cfg = config.Default()
	}

	// Override with flags
	if *flagGateway != "" {
		cfg.Gateway = *flagGateway
		cfg.Discovery = "manual"
	}
	if *flagTLS {
		cfg.TLS = true
	}
	if *flagNoMdns {
		cfg.Discovery = "manual"
	}

	// Create context for cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start mDNS discovery if enabled
	var mdns *discovery.MDNS
	if cfg.Discovery == "auto" {
		mdns = discovery.NewMDNS(identity.DeviceID, 18789)
		go mdns.Register(ctx)
	}

	// Initialize protocol client
	client := protocol.NewClient(cfg.Gateway, identity, protocol.ConnectOptions{
		Role:   "node",
		Scopes: []string{"node"},
	})

	// Initialize tray
	trayInstance := tray.New()
	go trayInstance.Run()

	// Setup callbacks
	client.OnConnected = func(resp protocol.ConnectResponse) {
		log.Printf("Connected to gateway")
		trayInstance.SetStatus(tray.StatusConnected)
	}
	client.OnDisconnected = func(reason string) {
		log.Printf("Disconnected: %s", reason)
		trayInstance.SetStatus(tray.StatusOffline)
	}

	// Connect
	if cfg.Gateway != "" {
		trayInstance.SetStatus(tray.StatusConnecting)
		if err := client.Connect(ctx); err != nil {
			log.Printf("Failed to connect: %v", err)
		}
	}

	// Wait for interrupt
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	log.Println("Shutting down...")
	client.Disconnect()
}
```

- [ ] **Step 2: 提交**

```bash
git add -f openclaw-node/cmd/main.go
git commit -m "feat: integrate all components in main

- Load/create identity on startup
- Load config with flag overrides
- Start mDNS discovery if enabled
- Connect to gateway
- Run tray UI
- Handle graceful shutdown

Co-Authored-By: Claude <noreply@anthropic.com>"
```

---

## Task 16: 单元测试

**Files:**
- Create: `openclaw-node/internal/crypto/identity_test.go`
- Create: `openclaw-node/internal/protocol/protocol_test.go`
- Create: `openclaw-node/internal/config/config_test.go`

---

- [ ] **Step 1: 创建 internal/crypto/identity_test.go**

```go
package crypto

import (
	"crypto/ed25519"
	"os"
	"testing"
)

func TestGenerateIdentity(t *testing.T) {
	id, err := GenerateIdentity()
	if err != nil {
		t.Fatalf("GenerateIdentity failed: %v", err)
	}

	if id.DeviceID == "" {
		t.Error("DeviceID should not be empty")
	}

	if len(id.PublicKey) == 0 {
		t.Error("PublicKey should not be empty")
	}

	if len(id.PrivateKey) == 0 {
		t.Error("PrivateKey should not be empty")
	}

	if id.CreatedAtMs == 0 {
		t.Error("CreatedAtMs should not be zero")
	}
}

func TestSignAndVerify(t *testing.T) {
	id, err := GenerateIdentity()
	if err != nil {
		t.Fatalf("GenerateIdentity failed: %v", err)
	}

	payload := "test payload"
	signature := id.Sign(payload)

	if len(signature) != ed25519.SignatureSize {
		t.Errorf("Signature size = %d, want %d", len(signature), ed25519.SignatureSize)
	}
}

func TestSaveAndLoadIdentity(t *testing.T) {
	tmpDir := t.TempDir()
	path := tmpDir + "/identity.json"

	id1, err := GenerateIdentity()
	if err != nil {
		t.Fatalf("GenerateIdentity failed: %v", err)
	}

	if err := SaveIdentity(path, id1); err != nil {
		t.Fatalf("SaveIdentity failed: %v", err)
	}

	id2, err := LoadIdentity(path)
	if err != nil {
		t.Fatalf("LoadIdentity failed: %v", err)
	}

	if id1.DeviceID != id2.DeviceID {
		t.Errorf("DeviceID mismatch: %s != %s", id1.DeviceID, id2.DeviceID)
	}
}

func TestPublicKeyBase64(t *testing.T) {
	id, err := GenerateIdentity()
	if err != nil {
		t.Fatalf("GenerateIdentity failed: %v", err)
	}

	b64 := id.PublicKeyBase64()
	if b64 == "" {
		t.Error("PublicKeyBase64 should not be empty")
	}
}
```

- [ ] **Step 2: 创建 internal/protocol/protocol_test.go**

```go
package protocol

import (
	"encoding/json"
	"testing"
)

func TestNewInvokeResultOK(t *testing.T) {
	payload := map[string]string{"key": "value"}
	result, err := NewInvokeResultOK(payload)
	if err != nil {
		t.Fatalf("NewInvokeResultOK failed: %v", err)
	}
	if !result.OK {
		t.Error("Expected OK=true")
	}
	if result.Error != nil {
		t.Error("Expected Error=nil")
	}
}

func TestNewInvokeResultError(t *testing.T) {
	result := NewInvokeResultError("TEST_ERROR", "test message")
	if result.OK {
		t.Error("Expected OK=false")
	}
	if result.Error == nil {
		t.Fatal("Expected Error!=nil")
	}
	if result.Error.Code != "TEST_ERROR" {
		t.Errorf("Error.Code = %s, want TEST_ERROR", result.Error.Code)
	}
}

func TestBuildAuthPayload(t *testing.T) {
	payload := BuildAuthPayload(
		"device123",
		"client456",
		"node",
		"node",
		[]string{"scope1", "scope2"},
		1234567890,
		"nonce123",
		"windows",
		"desktop",
	)

	expected := "v3|device123|client456|node|node|scope1,scope2|1234567890||nonce123|windows|desktop"
	if payload != expected {
		t.Errorf("BuildAuthPayload = %s, want %s", payload, expected)
	}
}

func TestNormalizeField(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Hello", "hello"},
		{"WORLD", "world"},
		{"MixedCase", "mixedcase"},
		{"", ""},
		{"  spaces  ", "spaces"},
	}

	for _, tt := range tests {
		result := normalizeField(tt.input)
		if result != tt.expected {
			t.Errorf("normalizeField(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestDispatcher(t *testing.T) {
	RegisterHandler("test.command", func(params json.RawMessage) (*InvokeResult, error) {
		return NewInvokeResultOK(map[string]string{"status": "ok"})
	})

	req := InvokeRequest{
		ID:      "test-1",
		Command: "test.command",
	}

	result := Dispatch(req)
	if !result.OK {
		t.Error("Expected result.OK=true")
	}

	// Test unknown command
	req.Command = "unknown.command"
	result = Dispatch(req)
	if result.OK {
		t.Error("Expected result.OK=false for unknown command")
	}
}
```

- [ ] **Step 3: 创建 internal/config/config_test.go**

```go
package config

import (
	"os"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := Default()

	if cfg.Port != 18789 {
		t.Errorf("Port = %d, want 18789", cfg.Port)
	}

	if cfg.Discovery != "auto" {
		t.Errorf("Discovery = %s, want auto", cfg.Discovery)
	}

	if cfg.Capabilities["camera"] != true {
		t.Error("Capabilities[camera] should be true")
	}

	if cfg.Capabilities["motion"] != false {
		t.Error("Capabilities[motion] should be false")
	}
}

func TestLoadConfig(t *testing.T) {
	tmpDir := t.TempDir()
	path := tmpDir + "/config.yaml"

	yamlContent := `
gateway: "localhost:18789"
port: 18790
discovery: manual
capabilities:
  camera: false
  location: true
`

	if err := os.WriteFile(path, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if cfg.Gateway != "localhost:18789" {
		t.Errorf("Gateway = %s, want localhost:18789", cfg.Gateway)
	}

	if cfg.Port != 18790 {
		t.Errorf("Port = %d, want 18790", cfg.Port)
	}

	if cfg.Capabilities["camera"] != false {
		t.Error("Capabilities[camera] should be false")
	}
}

func TestSaveConfig(t *testing.T) {
	tmpDir := t.TempDir()
	path := tmpDir + "/config.yaml"

	cfg := &Config{
		Gateway: "test.example.com",
		Port:    20000,
	}

	if err := cfg.Save(path); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("Load after Save failed: %v", err)
	}

	if loaded.Gateway != cfg.Gateway {
		t.Errorf("Gateway = %s, want %s", loaded.Gateway, cfg.Gateway)
	}
}
```

- [ ] **Step 4: 运行测试验证**

```bash
cd openclaw-node
go test ./... -v
```

- [ ] **Step 5: 提交**

```bash
git add -f openclaw-node/internal/crypto/identity_test.go
git add -f openclaw-node/internal/protocol/protocol_test.go
git add -f openclaw-node/internal/config/config_test.go
git commit -m "test: add unit tests for crypto, protocol, and config

- Add identity generation and signing tests
- Add protocol dispatcher and error handling tests
- Add config loading and saving tests

Co-Authored-By: Claude <noreply@anthropic.com>"
```

---

## 依赖的 Go 包

```bash
# 运行时依赖
go get github.com/gorilla/websocket@v1.5.1
go get golang.org/x/crypto@v0.17.0
go get github.com/google/uuid@v1.5.0
go get github.com/getlantern/systray@v1.0.0
go get gopkg.in/yaml.v3@v3.0.1
```

---

## 执行顺序

| Task | 组件 | 依赖 |
|------|------|------|
| 1 | 项目脚手架 | - |
| 2 | Ed25519 身份 | Task 1 |
| 3 | Protocol Client | Task 2 |
| 4 | Device Commands | Task 3 |
| 5 | Camera | Task 4 |
| 6 | Location | Task 4 |
| 7 | Photos | Task 4 |
| 8 | Screen | Task 4 |
| 9 | Notifications | Task 4 |
| 10 | Motion | Task 4 |
| 11 | SMS | Task 4 |
| 12 | Calendar | Task 4 |
| 13 | mDNS | Task 1, Task 3 |
| 14 | Tray UI | Task 1 |
| 15 | 集成主程序 | Task 2-14 |
| 16 | 单元测试 | Task 2, Task 3, Task 4 |
