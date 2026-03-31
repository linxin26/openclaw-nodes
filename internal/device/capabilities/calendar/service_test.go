package calendar

import (
	"context"
	"testing"

	"github.com/openclaw/openclaw-node/internal/device/core"
)

type fakeProvider struct {
	events []Event
	add    *AddResult
	err    error
}

func (p fakeProvider) Descriptor() core.CapabilityDescriptor {
	return core.CapabilityDescriptor{Name: "calendar", Commands: []string{"calendar.events", "calendar.add"}, Tier: 2}
}
func (p fakeProvider) Permission() core.PermissionState                       { return core.PermissionGranted }
func (p fakeProvider) Availability() core.Availability                        { return core.Availability{Available: true} }
func (p fakeProvider) Events(context.Context, EventsRequest) ([]Event, error) { return p.events, p.err }
func (p fakeProvider) Add(context.Context, AddRequest) (*AddResult, error)    { return p.add, p.err }

func TestCalendarServiceNormalizesResultShape(t *testing.T) {
	svc := NewService(fakeProvider{events: []Event{{ID: "1", Title: "event"}}})
	items, err := svc.Events(context.Background(), EventsRequest{})
	if err != nil || len(items) != 1 {
		t.Fatalf("items=%#v err=%v", items, err)
	}
}

func TestCalendarServiceMapsProviderErrors(t *testing.T) {
	svc := NewService(fakeProvider{err: core.ErrCapabilityUnavailable})
	_, err := svc.Add(context.Background(), AddRequest{Title: "x", Start: 1})
	if err != core.ErrCapabilityUnavailable {
		t.Fatalf("err = %v", err)
	}
}
