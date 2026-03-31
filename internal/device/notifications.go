package device

import (
	"context"
	"encoding/json"

	capnotifications "github.com/openclaw/openclaw-node/internal/device/capabilities/notifications"
	"github.com/openclaw/openclaw-node/internal/protocol"
)

type Notification = capnotifications.Notification

func handleNotificationsList(params json.RawMessage) (*protocol.InvokeResult, error) {
	var args capnotifications.ListRequest
	if err := decodeParams(params, &args); err != nil {
		return protocol.NewInvokeResultError("INVALID_PARAMS", "Invalid JSON: "+err.Error())
	}
	items, err := Runtime().Notifications().List(context.Background(), args)
	if err != nil {
		return invokeError(err, "notifications unavailable")
	}
	return protocol.NewInvokeResultOK(map[string]interface{}{"notifications": items})
}

func handleNotificationsActions(params json.RawMessage) (*protocol.InvokeResult, error) {
	var args capnotifications.ActionRequest
	if err := decodeParams(params, &args); err != nil {
		return protocol.NewInvokeResultError("INVALID_PARAMS", "Invalid JSON: "+err.Error())
	}
	if err := Runtime().Notifications().Action(context.Background(), args); err != nil {
		return invokeError(err, "notifications actions are not supported")
	}
	return protocol.NewInvokeResultOK(map[string]interface{}{"success": true})
}
