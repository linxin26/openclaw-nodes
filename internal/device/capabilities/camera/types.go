package camera

import (
	"context"

	"github.com/openclaw/openclaw-node/internal/device/core"
)

type Device struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Position string `json:"position"`
}

type SnapRequest struct {
	CameraID string `json:"cameraId"`
	Flash    bool   `json:"flash"`
	MaxWidth int    `json:"maxWidth"`
}

type ClipRequest struct {
	CameraID     string `json:"cameraId"`
	DurationMs   int    `json:"durationMs"`
	MaxWidth     int    `json:"maxWidth"`
	IncludeAudio bool   `json:"includeAudio"`
}

type Provider interface {
	Descriptor() core.CapabilityDescriptor
	Permission() core.PermissionState
	Availability() core.Availability
	List(context.Context) ([]Device, error)
	Snap(context.Context, SnapRequest) (core.ImagePayload, error)
	Clip(context.Context, ClipRequest) (core.VideoPayload, error)
}
