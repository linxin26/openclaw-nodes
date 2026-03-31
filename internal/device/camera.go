package device

import (
	"context"
	"encoding/json"
	"time"

	capcamera "github.com/openclaw/openclaw-node/internal/device/capabilities/camera"
	"github.com/openclaw/openclaw-node/internal/device/core"
	"github.com/openclaw/openclaw-node/internal/protocol"
)

const defaultExecTimeout = 30 * time.Second

type Camera = capcamera.Device

func handleCameraList(params json.RawMessage) (*protocol.InvokeResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultExecTimeout)
	defer cancel()
	items, err := Runtime().Camera().List(ctx)
	if err != nil {
		return invokeError(err, "failed to list cameras")
	}
	return protocol.NewInvokeResultOK(map[string]interface{}{"cameras": items})
}

func handleCameraSnap(params json.RawMessage) (*protocol.InvokeResult, error) {
	var args capcamera.SnapRequest
	if err := decodeParams(params, &args); err != nil {
		return protocol.NewInvokeResultError("INVALID_PARAMS", "Invalid JSON: "+err.Error())
	}
	ctx, cancel := context.WithTimeout(context.Background(), defaultExecTimeout)
	defer cancel()
	payload, err := Runtime().Camera().Snap(ctx, args)
	if err != nil {
		return invokeError(err, "failed to capture image")
	}
	return protocol.NewInvokeResultOK(payload)
}

func handleCameraClip(params json.RawMessage) (*protocol.InvokeResult, error) {
	var args capcamera.ClipRequest
	if err := decodeParams(params, &args); err != nil {
		return protocol.NewInvokeResultError("INVALID_PARAMS", "Invalid JSON: "+err.Error())
	}
	ctx, cancel := context.WithTimeout(context.Background(), defaultExecTimeout)
	defer cancel()
	payload, err := Runtime().Camera().Clip(ctx, args)
	if err != nil {
		return invokeError(err, "failed to capture video")
	}
	return protocol.NewInvokeResultOK(payload)
}

var _ core.ImagePayload
