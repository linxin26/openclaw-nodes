package photos

import (
	"context"
	"sort"
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

func (s *Service) Latest(ctx context.Context, req LatestRequest) ([]Entry, error) {
	if req.Limit <= 0 {
		req.Limit = 20
	}
	entries, err := s.provider.List(ctx, s.provider.DefaultRoot())
	if err != nil {
		return nil, normalizePhotosError(err)
	}
	filtered := make([]Entry, 0, len(entries))
	for _, entry := range entries {
		if req.After > 0 && entry.CreatedAt < req.After {
			continue
		}
		entry.Format = normalizeFormat(entry.Format)
		filtered = append(filtered, entry)
	}
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].CreatedAt > filtered[j].CreatedAt
	})
	if len(filtered) > req.Limit {
		filtered = filtered[:req.Limit]
	}
	return filtered, nil
}

func normalizePhotosError(err error) error {
	if err == nil {
		return nil
	}
	if err == core.ErrCapabilityUnavailable || err == core.ErrRestricted || err == core.ErrNotSupported {
		return err
	}
	return core.ErrCapabilityUnavailable
}

func normalizeFormat(format string) string {
	format = strings.TrimPrefix(strings.ToLower(strings.TrimSpace(format)), ".")
	switch format {
	case "jpg", "":
		return "jpeg"
	default:
		return format
	}
}
