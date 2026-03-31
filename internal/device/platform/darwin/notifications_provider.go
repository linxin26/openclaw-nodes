package darwin

import (
	"context"

	"github.com/openclaw/openclaw-node/internal/device/capabilities/notifications"
	"github.com/openclaw/openclaw-node/internal/device/core"
)

type NotificationsProvider struct{}

func (p *NotificationsProvider) Descriptor() core.CapabilityDescriptor {
	return core.CapabilityDescriptor{Name: "notifications", DisplayName: "Notifications", Description: "macOS notifications provider skeleton.", Commands: []string{"notifications.list", "notifications.actions"}, Tier: 2}
}

func (p *NotificationsProvider) Permission() core.PermissionState { return core.PermissionRestricted }
func (p *NotificationsProvider) Availability() core.Availability {
	return core.Availability{Available: false, Reason: "macOS provider not yet implemented"}
}
func (p *NotificationsProvider) List(context.Context, notifications.ListRequest) ([]notifications.Notification, error) {
	return nil, core.ErrRestricted
}
func (p *NotificationsProvider) Action(context.Context, notifications.ActionRequest) error {
	return core.ErrRestricted
}
