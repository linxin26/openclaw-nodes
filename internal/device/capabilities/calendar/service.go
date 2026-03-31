package calendar

import (
	"context"

	"github.com/openclaw/openclaw-node/internal/device/core"
)

type Event struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	Location    string `json:"location,omitempty"`
	Start       int64  `json:"start"`
	End         int64  `json:"end"`
	AllDay      bool   `json:"allDay"`
	Reminder    int64  `json:"reminder,omitempty"`
}

type EventsRequest struct {
	After  int64 `json:"after"`
	Before int64 `json:"before"`
	Limit  int   `json:"limit"`
}

type AddRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Start       int64  `json:"start"`
	End         int64  `json:"end"`
	AllDay      bool   `json:"allDay"`
	Location    string `json:"location"`
	Reminder    int64  `json:"reminder"`
}

type AddResult struct {
	Success bool   `json:"success"`
	EventID string `json:"eventId"`
}

type Provider interface {
	Descriptor() core.CapabilityDescriptor
	Permission() core.PermissionState
	Availability() core.Availability
	Events(context.Context, EventsRequest) ([]Event, error)
	Add(context.Context, AddRequest) (*AddResult, error)
}

type Service struct {
	provider Provider
}

func NewService(provider Provider) *Service {
	return &Service{provider: provider}
}

func (s *Service) State(enabled bool) core.CapabilityState {
	availability := s.provider.Availability()
	availability.Enabled = enabled
	return core.CapabilityState{
		Descriptor:   s.provider.Descriptor(),
		Permission:   s.provider.Permission(),
		Availability: availability,
	}
}

func (s *Service) Events(ctx context.Context, req EventsRequest) ([]Event, error) {
	if req.Limit <= 0 {
		req.Limit = 50
	}
	items, err := s.provider.Events(ctx, req)
	if err != nil {
		if err == core.ErrRestricted || err == core.ErrNotSupported || err == core.ErrCapabilityUnavailable {
			return nil, err
		}
		return nil, core.ErrCapabilityUnavailable
	}
	return items, nil
}

func (s *Service) Add(ctx context.Context, req AddRequest) (*AddResult, error) {
	if req.End == 0 {
		req.End = req.Start + 3600000
	}
	result, err := s.provider.Add(ctx, req)
	if err != nil {
		if err == core.ErrRestricted || err == core.ErrNotSupported || err == core.ErrCapabilityUnavailable {
			return nil, err
		}
		return nil, core.ErrCapabilityUnavailable
	}
	return result, nil
}
