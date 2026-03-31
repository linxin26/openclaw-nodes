package protocol

import (
	"encoding/json"
	"sync"
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
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
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

type ConnectResponse struct {
	OK    bool   `json:"ok"`
	Token string `json:"token,omitempty"`
}

func NewInvokeResultOK(payload interface{}) (*InvokeResult, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return &InvokeResult{OK: true, Payload: data}, nil
}

func NewInvokeResultError(code, message string) (*InvokeResult, error) {
	return &InvokeResult{
		OK:    false,
		Error: &ErrorShape{Code: code, Message: message},
	}, nil
}

type InvokeHandler func(params json.RawMessage) (*InvokeResult, error)

// Protocol manages command handlers with thread-safe access
type Protocol struct {
	mu       sync.RWMutex
	handlers map[string]InvokeHandler
	Identity *Identity // Device identity for authentication
}

// GlobalProtocol is the global protocol instance for handler registration
// Initialized lazily to ensure thread-safe initialization order
var GlobalProtocol *Protocol

var globalProtocolOnce sync.Once

func getGlobalProtocol() *Protocol {
	globalProtocolOnce.Do(func() {
		GlobalProtocol = &Protocol{
			handlers: make(map[string]InvokeHandler),
		}
	})
	return GlobalProtocol
}

// RegisterHandler registers a command handler (thread-safe)
func RegisterHandler(command string, handler InvokeHandler) {
	getGlobalProtocol().RegisterHandler(command, handler)
}

// RegisterHandler registers a command handler on this protocol
func (p *Protocol) RegisterHandler(command string, handler InvokeHandler) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.handlers[command] = handler
}

// Dispatch looks up the handler for a command and executes it (thread-safe)
func Dispatch(req InvokeRequest) *InvokeResult {
	return getGlobalProtocol().Dispatch(req)
}

// GetHandler returns the handler for a command (thread-safe, exported)
func GetHandler(command string) (InvokeHandler, bool) {
	return getGlobalProtocol().GetHandler(command)
}

// Dispatch looks up the handler for a command and executes it
func (p *Protocol) Dispatch(req InvokeRequest) *InvokeResult {
	p.mu.RLock()
	handler, ok := p.handlers[req.Command]
	p.mu.RUnlock()
	if !ok {
		res, _ := NewInvokeResultError("INVALID_REQUEST", "Unknown command: "+req.Command)
		return res
	}

	result, err := handler(req.Params)
	if err != nil {
		res, _ := NewInvokeResultError("INTERNAL_ERROR", err.Error())
		return res
	}
	return result
}

// GetHandler returns the handler for a command (thread-safe)
func (p *Protocol) GetHandler(command string) (InvokeHandler, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	handler, ok := p.handlers[command]
	return handler, ok
}
