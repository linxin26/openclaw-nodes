package location

import (
	"context"

	"github.com/openclaw/openclaw-node/internal/device/core"
)

type Result struct {
	Lat       float64 `json:"lat"`
	Lng       float64 `json:"lng"`
	Accuracy  int     `json:"accuracy"`
	Source    string  `json:"source"`
	Timestamp int64   `json:"timestamp"`
}

type Provider interface {
	Descriptor() core.CapabilityDescriptor
	Permission() core.PermissionState
	Availability() core.Availability
	Get(context.Context, bool) (*Result, error)
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

func (s *Service) Get(ctx context.Context, precise bool) (*Result, error) {
	result, err := s.provider.Get(ctx, precise)
	if err != nil {
		if err == core.ErrRestricted || err == core.ErrNotSupported || err == core.ErrCapabilityUnavailable {
			return nil, err
		}
		return nil, core.ErrCapabilityUnavailable
	}
	return result, nil
}
