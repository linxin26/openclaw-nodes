package sms

import "github.com/openclaw/openclaw-node/internal/device/core"

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) State(enabled bool) core.CapabilityState {
	return core.CapabilityState{
		Descriptor: core.CapabilityDescriptor{
			Name:        "sms",
			DisplayName: "SMS",
			Description: "Desktop SMS placeholder.",
			Commands:    []string{"sms.send", "sms.search"},
			Tier:        3,
		},
		Permission: core.PermissionNotSupported,
		Availability: core.Availability{
			Enabled:   enabled,
			Available: false,
			Reason:    "desktop sms is not supported",
		},
	}
}
