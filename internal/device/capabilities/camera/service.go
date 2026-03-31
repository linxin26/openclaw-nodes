package camera

import (
	"context"
	"strings"

	"github.com/openclaw/openclaw-node/internal/device/core"
)

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

func (s *Service) List(ctx context.Context) ([]Device, error) {
	devices, err := s.provider.List(ctx)
	if err != nil {
		return nil, normalizeProviderError(err)
	}
	for i := range devices {
		devices[i].Position = normalizePosition(devices[i].Position)
	}
	return devices, nil
}

func (s *Service) Snap(ctx context.Context, req SnapRequest) (core.ImagePayload, error) {
	if req.MaxWidth <= 0 {
		req.MaxWidth = 1920
	}
	payload, err := s.provider.Snap(ctx, req)
	if err != nil {
		return core.ImagePayload{}, normalizeProviderError(err)
	}
	if payload.Format == "" {
		payload.Format = "jpeg"
	}
	if payload.Width == 0 {
		payload.Width = req.MaxWidth
	}
	return payload, nil
}

func (s *Service) Clip(ctx context.Context, req ClipRequest) (core.VideoPayload, error) {
	if req.MaxWidth <= 0 {
		req.MaxWidth = 1920
	}
	if req.DurationMs <= 0 {
		req.DurationMs = 5000
	}
	payload, err := s.provider.Clip(ctx, req)
	if err != nil {
		return core.VideoPayload{}, normalizeProviderError(err)
	}
	if payload.Format == "" {
		payload.Format = "mp4"
	}
	if payload.Width == 0 {
		payload.Width = req.MaxWidth
	}
	if payload.DurationMs == 0 {
		payload.DurationMs = req.DurationMs
	}
	return payload, nil
}

func normalizeProviderError(err error) error {
	switch err {
	case nil:
		return nil
	case core.ErrCapabilityUnavailable, core.ErrRestricted, core.ErrNotSupported:
		return err
	default:
		return core.ErrCapabilityUnavailable
	}
}

func normalizePosition(position string) string {
	switch strings.ToLower(strings.TrimSpace(position)) {
	case "front", "rear", "external":
		return strings.ToLower(position)
	case "back":
		return "rear"
	default:
		return "external"
	}
}
