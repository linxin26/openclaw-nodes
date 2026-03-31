package runtime

import (
	"github.com/openclaw/openclaw-node/internal/config"
	"github.com/openclaw/openclaw-node/internal/device/platform/linux"
)

func linuxProviderSet(cfg *config.Config) ProviderSet {
	_ = cfg
	return ProviderSet{
		Camera:        &linux.CameraProvider{},
		Photos:        &linux.PhotosProvider{},
		Screen:        &linux.ScreenProvider{},
		Location:      &linux.LocationProvider{},
		Notifications: &linux.NotificationsProvider{},
		Calendar:      &linux.CalendarProvider{},
	}
}
