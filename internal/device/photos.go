package device

import (
	"context"
	"encoding/json"

	capphotos "github.com/openclaw/openclaw-node/internal/device/capabilities/photos"
	"github.com/openclaw/openclaw-node/internal/protocol"
)

type Photo = capphotos.Entry

func handlePhotosLatest(params json.RawMessage) (*protocol.InvokeResult, error) {
	var args capphotos.LatestRequest
	if err := decodeParams(params, &args); err != nil {
		return protocol.NewInvokeResultError("INVALID_PARAMS", "Invalid JSON: "+err.Error())
	}
	items, err := Runtime().Photos().Latest(context.Background(), args)
	if err != nil {
		return invokeError(err, "cannot scan photos")
	}
	return protocol.NewInvokeResultOK(map[string]interface{}{"photos": items, "total": len(items)})
}
