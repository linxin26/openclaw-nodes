package runtime

import (
	"github.com/openclaw/openclaw-node/internal/config"
	darwinplatform "github.com/openclaw/openclaw-node/internal/device/platform/darwin"
	linuxplatform "github.com/openclaw/openclaw-node/internal/device/platform/linux"
)

func darwinProviderSet(cfg *config.Config) ProviderSet {
	_ = cfg
	return ProviderSet{
		Camera:        &darwinplatform.CameraProvider{},
		Photos:        &darwinplatform.PhotosProvider{},
		Screen:        &darwinplatform.ScreenProvider{},
		Location:      &darwinplatform.LocationProvider{},
		Notifications: &darwinplatform.NotificationsProvider{},
		Calendar:      &darwinplatform.CalendarProvider{},
	}
}

func linuxProviderSet(cfg *config.Config) ProviderSet {
	_ = cfg
	return ProviderSet{
		Camera:        &linuxplatform.CameraProvider{},
		Photos:        &linuxplatform.PhotosProvider{},
		Screen:        &linuxplatform.ScreenProvider{},
		Location:      &linuxplatform.LocationProvider{},
		Notifications: &linuxplatform.NotificationsProvider{},
		Calendar:      &linuxplatform.CalendarProvider{},
	}
}
