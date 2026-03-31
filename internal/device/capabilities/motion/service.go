package motion

import (
	"time"

	"github.com/openclaw/openclaw-node/internal/device/core"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) State(enabled bool) core.CapabilityState {
	return core.CapabilityState{
		Descriptor: core.CapabilityDescriptor{
			Name:        "motion",
			DisplayName: "Motion",
			Description: "Desktop motion sensors placeholder.",
			Commands:    []string{"motion.activity", "motion.pedometer"},
			Tier:        3,
		},
		Permission: core.PermissionNotApplicable,
		Availability: core.Availability{
			Enabled:   enabled,
			Available: false,
			Reason:    "desktop motion sensors are not implemented",
		},
	}
}

func (s *Service) Activity() map[string]interface{} {
	return map[string]interface{}{
		"activity":   "not_supported",
		"confidence": 0,
		"timestamp":  time.Now().UnixMilli(),
	}
}

func (s *Service) Pedometer() map[string]interface{} {
	return map[string]interface{}{
		"steps":       0,
		"distance":    0,
		"startTime":   time.Now().UnixMilli(),
		"timestamp":   time.Now().UnixMilli(),
		"available":   false,
		"placeholder": true,
	}
}
