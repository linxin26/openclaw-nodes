package windows

import (
	"context"
	"time"

	"github.com/openclaw/openclaw-node/internal/device/capabilities/notifications"
	"github.com/openclaw/openclaw-node/internal/device/core"
)

type NotificationsProvider struct{}

func (p *NotificationsProvider) Descriptor() core.CapabilityDescriptor {
	return core.CapabilityDescriptor{Name: "notifications", DisplayName: "Notifications", Description: "Read desktop notifications.", Commands: []string{"notifications.list", "notifications.actions"}, Tier: 2}
}

func (p *NotificationsProvider) Permission() core.PermissionState { return core.PermissionGranted }
func (p *NotificationsProvider) Availability() core.Availability {
	return core.Availability{Available: true}
}

func (p *NotificationsProvider) List(ctx context.Context, req notifications.ListRequest) ([]notifications.Notification, error) {
	_ = ctx
	_ = req
	return []notifications.Notification{{
		ID:       "1",
		App:      "System",
		Title:    "OpenClaw Node",
		Body:     "Node is running",
		PostedAt: time.Now().UnixMilli(),
	}}, nil
}

func (p *NotificationsProvider) Action(ctx context.Context, req notifications.ActionRequest) error {
	_ = ctx
	_ = req
	return core.ErrNotSupported
}
