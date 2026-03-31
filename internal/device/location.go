package device

import (
	"context"
	"encoding/json"
	"time"

	caplocation "github.com/openclaw/openclaw-node/internal/device/capabilities/location"
	"github.com/openclaw/openclaw-node/internal/protocol"
)

const defaultHTTPTimeout = 10 * time.Second

type Location = caplocation.Result

func handleLocationGet(params json.RawMessage) (*protocol.InvokeResult, error) {
	var args struct {
		Precise bool `json:"precise"`
	}
	if err := decodeParams(params, &args); err != nil {
		return protocol.NewInvokeResultError("INVALID_PARAMS", "Invalid JSON: "+err.Error())
	}
	ctx, cancel := context.WithTimeout(context.Background(), defaultHTTPTimeout)
	defer cancel()
	payload, err := Runtime().Location().Get(ctx, args.Precise)
	if err != nil {
		return invokeError(err, "cannot get location")
	}
	return protocol.NewInvokeResultOK(payload)
}
