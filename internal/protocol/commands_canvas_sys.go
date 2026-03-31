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

// Canvas handlers - return UNSUPPORTED since canvas requires full UI environment
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

func handleSystemNotify(params json.RawMessage) (*InvokeResult, error) {
	var args struct {
		Title string `json:"title"`
		Body  string `json:"body"`
	}
	if params != nil {
		json.Unmarshal(params, &args) //nolint:errcheck
	}
	fmt.Printf("System notify: %s - %s\n", args.Title, args.Body)
	result, err := NewInvokeResultOK(map[string]interface{}{"sent": true})
	return result, err
}

func handleDebugLogs(params json.RawMessage) (*InvokeResult, error) {
	result, err := NewInvokeResultOK(map[string]interface{}{"logs": []string{}})
	return result, err
}

func handleDebugEd25519(params json.RawMessage) (*InvokeResult, error) {
	if GlobalProtocol == nil || GlobalProtocol.Identity == nil {
		return NewInvokeResultError("NOT_INITIALIZED", "Identity not set")
	}
	result, err := NewInvokeResultOK(map[string]interface{}{
		"deviceId":   GlobalProtocol.Identity.DeviceID,
		"clientId":   GlobalProtocol.Identity.ClientID,
		"signedAtMs": GlobalProtocol.Identity.SignedAtMs,
	})
	return result, err
}
