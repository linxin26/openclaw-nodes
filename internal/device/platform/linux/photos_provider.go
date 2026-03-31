package linux

import (
	"context"

	"github.com/openclaw/openclaw-node/internal/device/capabilities/photos"
	"github.com/openclaw/openclaw-node/internal/device/core"
)

type PhotosProvider struct{}

func (p *PhotosProvider) Descriptor() core.CapabilityDescriptor {
	return core.CapabilityDescriptor{Name: "photos", DisplayName: "Photos", Description: "Linux photos provider skeleton.", Commands: []string{"photos.latest"}, Tier: 1}
}

func (p *PhotosProvider) Permission() core.PermissionState { return core.PermissionRestricted }
func (p *PhotosProvider) Availability() core.Availability {
	return core.Availability{Available: false, Reason: "linux provider not yet implemented"}
}
func (p *PhotosProvider) DefaultRoot() string { return "" }
func (p *PhotosProvider) List(context.Context, string) ([]photos.Entry, error) {
	return nil, core.ErrRestricted
}
