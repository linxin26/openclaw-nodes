package darwin

import (
	"context"

	"github.com/openclaw/openclaw-node/internal/device/capabilities/camera"
	"github.com/openclaw/openclaw-node/internal/device/core"
)

type CameraProvider struct{}

func (p *CameraProvider) Descriptor() core.CapabilityDescriptor {
	return core.CapabilityDescriptor{Name: "camera", DisplayName: "Camera", Description: "macOS camera provider skeleton.", Commands: []string{"camera.list", "camera.snap", "camera.clip"}, Tier: 1}
}

func (p *CameraProvider) Permission() core.PermissionState { return core.PermissionRestricted }
func (p *CameraProvider) Availability() core.Availability {
	return core.Availability{Available: false, Reason: "macOS provider not yet implemented"}
}
func (p *CameraProvider) List(context.Context) ([]camera.Device, error) {
	return nil, core.ErrRestricted
}
func (p *CameraProvider) Snap(context.Context, camera.SnapRequest) (core.ImagePayload, error) {
	return core.ImagePayload{}, core.ErrRestricted
}
func (p *CameraProvider) Clip(context.Context, camera.ClipRequest) (core.VideoPayload, error) {
	return core.VideoPayload{}, core.ErrRestricted
}
