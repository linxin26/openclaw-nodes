package photos

import (
	"context"
	"testing"

	"github.com/openclaw/openclaw-node/internal/device/core"
)

type fakeProvider struct {
	entries []Entry
	err     error
}

func (p fakeProvider) Descriptor() core.CapabilityDescriptor {
	return core.CapabilityDescriptor{Name: "photos", Commands: []string{"photos.latest"}, Tier: 1}
}
func (p fakeProvider) Permission() core.PermissionState              { return core.PermissionGranted }
func (p fakeProvider) Availability() core.Availability               { return core.Availability{Available: true} }
func (p fakeProvider) DefaultRoot() string                           { return "photos" }
func (p fakeProvider) List(context.Context, string) ([]Entry, error) { return p.entries, p.err }

func TestPhotosServiceFiltersAndNormalizes(t *testing.T) {
	svc := NewService(fakeProvider{entries: []Entry{{ID: "1", CreatedAt: 10, Format: ".jpg"}, {ID: "2", CreatedAt: 20, Format: "png"}}})
	items, err := svc.Latest(context.Background(), LatestRequest{Limit: 1, After: 15})
	if err != nil {
		t.Fatalf("Latest() error = %v", err)
	}
	if len(items) != 1 || items[0].Format != "png" {
		t.Fatalf("items = %#v", items)
	}
}

func TestPhotosServiceMapsUnavailableProviderErrors(t *testing.T) {
	svc := NewService(fakeProvider{err: core.ErrCapabilityUnavailable})
	_, err := svc.Latest(context.Background(), LatestRequest{})
	if err != core.ErrCapabilityUnavailable {
		t.Fatalf("err = %v", err)
	}
}
