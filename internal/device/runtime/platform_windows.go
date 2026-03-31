package runtime

import (
	"net/http"
	"time"

	"github.com/openclaw/openclaw-node/internal/config"
	win "github.com/openclaw/openclaw-node/internal/device/platform/windows"
)

func windowsProviderSet(cfg *config.Config) ProviderSet {
	root := ""
	if cfg != nil {
		if option, ok := cfg.CapabilityOptions["photos"]; ok {
			root = option.Path
		}
	}
	return ProviderSet{
		Camera:        &win.CameraProvider{},
		Photos:        &win.PhotosProvider{Root: root},
		Screen:        &win.ScreenProvider{},
		Location:      &win.LocationProvider{Client: &http.Client{Timeout: 10 * time.Second}},
		Notifications: &win.NotificationsProvider{},
		Calendar:      &win.CalendarProvider{},
	}
}
