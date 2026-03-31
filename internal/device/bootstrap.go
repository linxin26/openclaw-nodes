package device

import (
	"encoding/json"

	"github.com/openclaw/openclaw-node/internal/config"
	"github.com/openclaw/openclaw-node/internal/device/core"
	deviceruntime "github.com/openclaw/openclaw-node/internal/device/runtime"
	"github.com/openclaw/openclaw-node/internal/protocol"
)

func Bootstrap(cfg *config.Config) *deviceruntime.Runtime {
	return deviceruntime.MustBootstrap(cfg)
}

func Runtime() *deviceruntime.Runtime {
	return deviceruntime.Default()
}

func RegisterProtocolHandlers(register func(string, protocol.InvokeHandler)) {
	register("camera.list", handleCameraList)
	register("camera.snap", handleCameraSnap)
	register("camera.clip", handleCameraClip)
	register("photos.latest", handlePhotosLatest)
	register("screen.snapshot", handleScreenSnapshot)
	register("location.get", handleLocationGet)
	register("notifications.list", handleNotificationsList)
	register("notifications.actions", handleNotificationsActions)
	register("calendar.events", handleCalendarEvents)
	register("calendar.add", handleCalendarAdd)
	register("motion.activity", handleMotionActivity)
	register("motion.pedometer", handleMotionPedometer)
	register("sms.send", handleSmsSend)
	register("sms.search", handleSmsSearch)
}

func invokeError(err error, unavailableMsg string) (*protocol.InvokeResult, error) {
	switch err {
	case nil:
		return nil, nil
	case core.ErrRestricted:
		return protocol.NewInvokeResultError("PERMISSION_RESTRICTED", unavailableMsg)
	case core.ErrNotSupported:
		return protocol.NewInvokeResultError("NOT_SUPPORTED", unavailableMsg)
	default:
		return protocol.NewInvokeResultError("CAPABILITY_UNAVAILABLE", unavailableMsg)
	}
}

func decodeParams(raw json.RawMessage, target interface{}) error {
	if raw == nil {
		return nil
	}
	return json.Unmarshal(raw, target)
}
