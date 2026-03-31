package screen

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

func (s *Service) Snapshot(ctx context.Context, req SnapshotRequest) (core.ImagePayload, error) {
	if req.Format == "" {
		req.Format = "png"
	}
	req.Format = strings.ToLower(req.Format)
	if req.Quality <= 0 {
		req.Quality = 90
	}
	if req.MaxWidth <= 0 {
		req.MaxWidth = 1920
	}
	payload, err := s.provider.Snapshot(ctx, req)
	if err != nil {
		return core.ImagePayload{}, normalizeScreenError(err)
	}
	if payload.Format == "" {
		payload.Format = req.Format
	}
	if payload.Width == 0 {
		payload.Width = req.MaxWidth
	}
	return payload, nil
}

func normalizeScreenError(err error) error {
	if err == nil {
		return nil
	}
	if err == core.ErrCapabilityUnavailable || err == core.ErrRestricted || err == core.ErrNotSupported {
		return err
	}
	return core.ErrCapabilityUnavailable
}
