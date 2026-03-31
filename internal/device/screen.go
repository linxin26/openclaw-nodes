package device

import (
	"context"
	"encoding/json"

	capscreen "github.com/openclaw/openclaw-node/internal/device/capabilities/screen"
	"github.com/openclaw/openclaw-node/internal/protocol"
)

func handleScreenSnapshot(params json.RawMessage) (*protocol.InvokeResult, error) {
	var args capscreen.SnapshotRequest
	if err := decodeParams(params, &args); err != nil {
		return protocol.NewInvokeResultError("INVALID_PARAMS", "Invalid JSON: "+err.Error())
	}
	ctx, cancel := context.WithTimeout(context.Background(), defaultExecTimeout)
	defer cancel()
	payload, err := Runtime().Screen().Snapshot(ctx, args)
	if err != nil {
		return invokeError(err, "failed to capture screen")
	}
	return protocol.NewInvokeResultOK(payload)
}
