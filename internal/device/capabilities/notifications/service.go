package notifications

import (
	"context"

	"github.com/openclaw/openclaw-node/internal/device/core"
)

type Notification struct {
	ID       string `json:"id"`
	App      string `json:"app"`
	Title    string `json:"title"`
	Body     string `json:"body"`
	PostedAt int64  `json:"postedAt"`
}

type ListRequest struct {
	Limit int      `json:"limit"`
	Apps  []string `json:"apps"`
}

type ActionRequest struct {
	ID     string `json:"id"`
	Action string `json:"action"`
}

type Provider interface {
	Descriptor() core.CapabilityDescriptor
	Permission() core.PermissionState
	Availability() core.Availability
	List(context.Context, ListRequest) ([]Notification, error)
	Action(context.Context, ActionRequest) error
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

func (s *Service) List(ctx context.Context, req ListRequest) ([]Notification, error) {
	if req.Limit <= 0 {
		req.Limit = 20
	}
	items, err := s.provider.List(ctx, req)
	if err != nil {
		if err == core.ErrRestricted || err == core.ErrNotSupported || err == core.ErrCapabilityUnavailable {
			return nil, err
		}
		return nil, core.ErrCapabilityUnavailable
	}
	return items, nil
}

func (s *Service) Action(ctx context.Context, req ActionRequest) error {
	err := s.provider.Action(ctx, req)
	if err == nil || err == core.ErrRestricted || err == core.ErrNotSupported || err == core.ErrCapabilityUnavailable {
		return err
	}
	return core.ErrCapabilityUnavailable
}
