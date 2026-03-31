package runtime

import "github.com/openclaw/openclaw-node/internal/config"

func BuildProviderSet(platform string, cfg *config.Config) ProviderSet {
	switch platform {
	case "darwin":
		return darwinProviderSet(cfg)
	case "linux":
		return linuxProviderSet(cfg)
	default:
		return windowsProviderSet(cfg)
	}
}
