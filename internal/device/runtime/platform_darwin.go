package runtime

import (
	"github.com/openclaw/openclaw-node/internal/config"
	"github.com/openclaw/openclaw-node/internal/device/platform/darwin"
)

func darwinProviderSet(cfg *config.Config) ProviderSet {
	_ = cfg
	return ProviderSet{
		Camera:        &darwin.CameraProvider{},
		Photos:        &darwin.PhotosProvider{},
		Screen:        &darwin.ScreenProvider{},
		Location:      &darwin.LocationProvider{},
		Notifications: &darwin.NotificationsProvider{},
		Calendar:      &darwin.CalendarProvider{},
	}
}
