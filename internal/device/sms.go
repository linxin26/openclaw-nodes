package device

import (
	"encoding/json"

	"github.com/openclaw/openclaw-node/internal/protocol"
)

func handleSmsSend(params json.RawMessage) (*protocol.InvokeResult, error) {
	return protocol.NewInvokeResultError("NOT_SUPPORTED", "desktop sms is not supported")
}

func handleSmsSearch(params json.RawMessage) (*protocol.InvokeResult, error) {
	return protocol.NewInvokeResultError("NOT_SUPPORTED", "desktop sms is not supported")
}
