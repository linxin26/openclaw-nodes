package darwin

import (
	"context"

	"github.com/openclaw/openclaw-node/internal/device/capabilities/location"
	"github.com/openclaw/openclaw-node/internal/device/core"
)

type LocationProvider struct{}

func (p *LocationProvider) Descriptor() core.CapabilityDescriptor {
	return core.CapabilityDescriptor{Name: "location", DisplayName: "Location", Description: "macOS location provider skeleton.", Commands: []string{"location.get"}, Tier: 2}
}

func (p *LocationProvider) Permission() core.PermissionState { return core.PermissionRestricted }
func (p *LocationProvider) Availability() core.Availability {
	return core.Availability{Available: false, Reason: "macOS provider not yet implemented"}
}
func (p *LocationProvider) Get(context.Context, bool) (*location.Result, error) {
	return nil, core.ErrRestricted
}
