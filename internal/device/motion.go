package device

import (
	"encoding/json"

	"github.com/openclaw/openclaw-node/internal/protocol"
)

func handleMotionActivity(params json.RawMessage) (*protocol.InvokeResult, error) {
	return protocol.NewInvokeResultError("NOT_SUPPORTED", "desktop motion sensors are not supported")
}

func handleMotionPedometer(params json.RawMessage) (*protocol.InvokeResult, error) {
	return protocol.NewInvokeResultError("NOT_SUPPORTED", "desktop motion sensors are not supported")
}
