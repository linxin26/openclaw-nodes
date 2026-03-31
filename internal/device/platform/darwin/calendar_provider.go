package darwin

import (
	"context"

	"github.com/openclaw/openclaw-node/internal/device/capabilities/calendar"
	"github.com/openclaw/openclaw-node/internal/device/core"
)

type CalendarProvider struct{}

func (p *CalendarProvider) Descriptor() core.CapabilityDescriptor {
	return core.CapabilityDescriptor{Name: "calendar", DisplayName: "Calendar", Description: "macOS calendar provider skeleton.", Commands: []string{"calendar.events", "calendar.add"}, Tier: 2}
}

func (p *CalendarProvider) Permission() core.PermissionState { return core.PermissionRestricted }
func (p *CalendarProvider) Availability() core.Availability {
	return core.Availability{Available: false, Reason: "macOS provider not yet implemented"}
}
func (p *CalendarProvider) Events(context.Context, calendar.EventsRequest) ([]calendar.Event, error) {
	return nil, core.ErrRestricted
}
func (p *CalendarProvider) Add(context.Context, calendar.AddRequest) (*calendar.AddResult, error) {
	return nil, core.ErrRestricted
}
