package runtime

import (
	"testing"

	"github.com/openclaw/openclaw-node/internal/config"
	"github.com/openclaw/openclaw-node/internal/device/core"
)

func TestRegistryBuildsDeviceDescribeFromProviders(t *testing.T) {
	rt := New(config.Default(), "windows")
	if got := len(rt.Registry().CapabilityNames()); got < 8 {
		t.Fatalf("CapabilityNames() len = %d, want >= 8", got)
	}
	if got := len(rt.Registry().Commands()); got < 10 {
		t.Fatalf("Commands() len = %d, want >= 10", got)
	}
}

func TestRegistryBuildsDevicePermissionsFromProviders(t *testing.T) {
	rt := New(config.Default(), "windows")
	permissions := rt.Registry().Permissions()
	if permissions["camera"] != core.PermissionGranted {
		t.Fatalf("permissions[camera] = %q", permissions["camera"])
	}
	if permissions["sms"] != core.PermissionNotSupported {
		t.Fatalf("permissions[sms] = %q", permissions["sms"])
	}
}

func TestRegistryBuildsDeviceStatusFromProviders(t *testing.T) {
	cfg := config.Default()
	cfg.Capabilities["calendar"] = true
	rt := New(cfg, "windows")
	status := rt.Registry().Availability()
	if !status["camera"].Enabled || !status["camera"].Available {
		t.Fatalf("camera availability = %#v", status["camera"])
	}
	if status["motion"].Available {
		t.Fatalf("motion availability = %#v, want unavailable", status["motion"])
	}
}

func TestRuntimeRegistersCommandsFromEnabledProviders(t *testing.T) {
	rt := New(config.Default(), "windows")
	if got := len(rt.Registry().Commands()); got < 10 {
		t.Fatalf("Commands() len = %d, want >= 10", got)
	}
}

func TestRuntimeSkipsDisabledCapabilities(t *testing.T) {
	cfg := config.Default()
	cfg.Capabilities["camera"] = false
	rt := New(cfg, "windows")
	if rt.Registry().Availability()["camera"].Enabled {
		t.Fatal("camera enabled = true, want false")
	}
}

func TestRuntimeSelectsCurrentPlatformProviderSet(t *testing.T) {
	if New(config.Default(), "windows").Platform() != "windows" {
		t.Fatal("windows runtime platform mismatch")
	}
	if New(config.Default(), "darwin").Platform() != "darwin" {
		t.Fatal("darwin runtime platform mismatch")
	}
	if New(config.Default(), "linux").Platform() != "linux" {
		t.Fatal("linux runtime platform mismatch")
	}
}

func TestBuildProviderSetProvidesTier1Providers(t *testing.T) {
	cfg := config.Default()
	if set := BuildProviderSet("darwin", cfg); set.Camera == nil || set.Photos == nil || set.Screen == nil {
		t.Fatal("darwin provider set missing tier1 providers")
	}
	if set := BuildProviderSet("linux", cfg); set.Camera == nil || set.Photos == nil || set.Screen == nil {
		t.Fatal("linux provider set missing tier1 providers")
	}
}
