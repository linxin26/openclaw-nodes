package wails

import (
	"reflect"
	"testing"

	"github.com/openclaw/openclaw-node/internal/config"
	appcrypto "github.com/openclaw/openclaw-node/internal/crypto"
	"github.com/openclaw/openclaw-node/internal/device"
	deviceruntime "github.com/openclaw/openclaw-node/internal/device/runtime"
	"github.com/openclaw/openclaw-node/internal/protocol"
)

func newTestApp(t *testing.T) *App {
	t.Helper()
	identity, err := appcrypto.GenerateIdentity()
	if err != nil {
		t.Fatalf("GenerateIdentity() error = %v", err)
	}
	cfg := config.Default()
	cfg.Gateway = "gateway.example:443"
	cfg.TLS = true
	device.Bootstrap(cfg)
	device.RegisterProtocolHandlers(protocol.RegisterHandler)
	client := protocol.NewClient(cfg.Gateway, &protocol.Identity{DeviceID: identity.DeviceID, Role: "node"}, identity, protocol.ConnectOptions{})
	return NewApp(t.TempDir(), identity, cfg, client)
}

func TestGetStatusReturnsOfflineByDefault(t *testing.T) {
	appInstance = nil
	app := newTestApp(t)
	status := app.GetStatus()
	if status.Status != StatusOffline {
		t.Fatalf("status.Status = %q, want %q", status.Status, StatusOffline)
	}
	if status.Gateway != "gateway.example:443" {
		t.Fatalf("status.Gateway = %q", status.Gateway)
	}
	if !status.TLS {
		t.Fatal("status.TLS = false, want true")
	}
	if _, ok := status.Capabilities["camera"]; !ok {
		t.Fatal("camera capability missing from status")
	}
}

func TestGetConfigReturnsEditableFields(t *testing.T) {
	appInstance = nil
	app := newTestApp(t)
	cfg := app.GetConfig()
	if cfg.Gateway != "gateway.example:443" {
		t.Fatalf("cfg.Gateway = %q", cfg.Gateway)
	}
	if cfg.Port != 18789 {
		t.Fatalf("cfg.Port = %d, want 18789", cfg.Port)
	}
	if cfg.Token != "" {
		t.Fatalf("cfg.Token = %q, want empty string", cfg.Token)
	}
	if len(cfg.Capabilities) == 0 {
		t.Fatal("cfg.Capabilities is empty")
	}
}

func TestGetDeviceInfoUsesIdentity(t *testing.T) {
	appInstance = nil
	app := newTestApp(t)
	info := app.GetDeviceInfo()
	if info.DeviceID == "" {
		t.Fatal("info.DeviceID is empty")
	}
	if info.Mode != "node" {
		t.Fatalf("info.Mode = %q, want node", info.Mode)
	}
	if info.Version != deviceruntime.Default().Metadata().Version {
		t.Fatalf("info.Version = %q", info.Version)
	}
}

func TestNormalizeGatewayUsesConfiguredPort(t *testing.T) {
	got := normalizeGateway("localhost", 38789)
	if got != "localhost:38789" {
		t.Fatalf("normalizeGateway() = %q, want localhost:38789", got)
	}
}

func TestNormalizeGatewayPreservesExplicitPort(t *testing.T) {
	got := normalizeGateway("localhost:40123", 38789)
	if got != "localhost:40123" {
		t.Fatalf("normalizeGateway() = %q, want localhost:40123", got)
	}
}

func TestSaveConfigUpdatesClientToken(t *testing.T) {
	appInstance = nil
	app := newTestApp(t)
	err := app.SaveConfig(&Config{Gateway: "localhost", Port: 38789, Token: "gateway-secret", TLS: false, Discovery: "manual", Capabilities: map[string]bool{"camera": true}, CapabilityOptions: map[string]CapabilityOption{"photos": {Provider: "windows", Path: "C:/photos"}}})
	if err != nil {
		t.Fatalf("SaveConfig() error = %v", err)
	}
	cfg := app.GetConfig()
	if cfg.Token != "gateway-secret" {
		t.Fatalf("cfg.Token = %q, want gateway-secret", cfg.Token)
	}
	if cfg.CapabilityOptions["photos"].Path != "C:/photos" {
		t.Fatalf("photos path = %q", cfg.CapabilityOptions["photos"].Path)
	}
	clientToken := reflect.ValueOf(app.client).Elem().FieldByName("opts").FieldByName("Token").String()
	if clientToken != "gateway-secret" {
		t.Fatalf("client token = %q, want gateway-secret", clientToken)
	}
}

func TestGetCapabilitiesReflectsRuntimeState(t *testing.T) {
	appInstance = nil
	app := newTestApp(t)
	items := app.GetCapabilities()
	if len(items) == 0 {
		t.Fatal("GetCapabilities() returned empty list")
	}
	var motion *CapabilityInfo
	for _, item := range items {
		if item.Key == "motion" {
			motion = item
			break
		}
	}
	if motion == nil {
		t.Fatal("motion capability missing")
	}
	if motion.Available {
		t.Fatalf("motion.Available = true, want false")
	}
	if motion.Permission == "" {
		t.Fatal("motion permission empty")
	}
}
