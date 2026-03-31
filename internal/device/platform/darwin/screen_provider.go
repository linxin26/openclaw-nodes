package darwin

import (
	"context"

	"github.com/openclaw/openclaw-node/internal/device/capabilities/screen"
	"github.com/openclaw/openclaw-node/internal/device/core"
)

type ScreenProvider struct{}

func (p *ScreenProvider) Descriptor() core.CapabilityDescriptor {
	return core.CapabilityDescriptor{Name: "screen", DisplayName: "Screen", Description: "macOS screen provider skeleton.", Commands: []string{"screen.snapshot"}, Tier: 1}
}

func (p *ScreenProvider) Permission() core.PermissionState { return core.PermissionRestricted }
func (p *ScreenProvider) Availability() core.Availability {
	return core.Availability{Available: false, Reason: "macOS provider not yet implemented"}
}
func (p *ScreenProvider) Snapshot(context.Context, screen.SnapshotRequest) (core.ImagePayload, error) {
	return core.ImagePayload{}, core.ErrRestricted
}
