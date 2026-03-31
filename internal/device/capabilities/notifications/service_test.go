package notifications

import (
	"context"
	"testing"

	"github.com/openclaw/openclaw-node/internal/device/core"
)

type fakeProvider struct {
	items []Notification
	err   error
}

func (p fakeProvider) Descriptor() core.CapabilityDescriptor {
	return core.CapabilityDescriptor{Name: "notifications", Commands: []string{"notifications.list", "notifications.actions"}, Tier: 2}
}
func (p fakeProvider) Permission() core.PermissionState { return core.PermissionGranted }
func (p fakeProvider) Availability() core.Availability  { return core.Availability{Available: true} }
func (p fakeProvider) List(context.Context, ListRequest) ([]Notification, error) {
	return p.items, p.err
}
func (p fakeProvider) Action(context.Context, ActionRequest) error { return p.err }

func TestNotificationsServiceNormalizesResultShape(t *testing.T) {
	svc := NewService(fakeProvider{items: []Notification{{ID: "1", Title: "hello"}}})
	items, err := svc.List(context.Background(), ListRequest{})
	if err != nil || len(items) != 1 {
		t.Fatalf("items=%#v err=%v", items, err)
	}
}

func TestNotificationsServiceMapsProviderErrors(t *testing.T) {
	svc := NewService(fakeProvider{err: core.ErrNotSupported})
	if err := svc.Action(context.Background(), ActionRequest{}); err != core.ErrNotSupported {
		t.Fatalf("err = %v", err)
	}
}
