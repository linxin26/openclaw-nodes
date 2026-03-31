package main

import (
	"context"
	"embed"
	"log"
	"sort"

	"github.com/openclaw/openclaw-node/internal/config"
	"github.com/openclaw/openclaw-node/internal/crypto"
	"github.com/openclaw/openclaw-node/internal/device"
	"github.com/openclaw/openclaw-node/internal/device/core"
	deviceruntime "github.com/openclaw/openclaw-node/internal/device/runtime"
	"github.com/openclaw/openclaw-node/internal/protocol"
	appwails "github.com/openclaw/openclaw-node/internal/wails"
	"github.com/openclaw/openclaw-node/store"
	wails "github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	windowsoptions "github.com/wailsapp/wails/v2/pkg/options/windows"
	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed all:frontend/dist
var assets embed.FS

type commandRegistry struct {
	caps        []string
	cmds        []string
	permissions map[string]bool
}

func newRegistry() *commandRegistry {
	rt := deviceruntime.Default()
	permissions := map[string]bool{}
	for name, state := range rt.Registry().Permissions() {
		permissions[name] = state == core.PermissionGranted
	}
	cmds := append([]string{"device.describe", "device.info", "device.status", "device.health", "device.permissions"}, rt.Registry().Commands()...)
	sort.Strings(cmds)
	return &commandRegistry{caps: rt.Registry().CapabilityNames(), cmds: cmds, permissions: permissions}
}

func buildGUIConnectOptions(identity *crypto.Identity, registry *commandRegistry) protocol.ConnectOptions {
	meta := deviceruntime.Default().Metadata()
	return protocol.ConnectOptions{Role: "node", Scopes: []string{}, Caps: registry.caps, Commands: registry.cmds, Permissions: registry.permissions, Token: "", Client: protocol.ClientInfo{ID: "node-host", DisplayName: "OpenClaw Node", Version: meta.Version, Platform: meta.Platform, Mode: "node", InstanceID: identity.DeviceID[:8], DeviceFamily: "desktop", ModelIdentifier: meta.ModelIdentifier}}
}

func main() {
	dataDir, err := store.DefaultDataDir()
	if err != nil {
		log.Fatal(err)
	}
	s, err := store.New(dataDir)
	if err != nil {
		log.Fatal(err)
	}
	identityPath := s.Path("identity.json")
	identity, err := crypto.LoadIdentity(identityPath)
	if err != nil {
		identity, err = crypto.GenerateIdentity()
		if err != nil {
			log.Fatal(err)
		}
		if err := crypto.SaveIdentity(identityPath, identity); err != nil {
			log.Fatal(err)
		}
	}
	cfgPath := s.Path("config.yaml")
	cfg, _ := config.Load(cfgPath)
	if cfg == nil {
		cfg = config.Default()
	}
	device.Bootstrap(cfg)
	device.RegisterProtocolHandlers(protocol.RegisterHandler)
	protoIdentity := &protocol.Identity{DeviceID: identity.DeviceID, Role: "node"}
	registry := newRegistry()
	protocol.GlobalProtocol.Identity = &protocol.Identity{DeviceID: identity.DeviceID, ClientID: "node-host", ClientMode: "node", Role: "node", SignedAtMs: identity.CreatedAtMs}
	connectOptions := buildGUIConnectOptions(identity, registry)
	connectOptions.Token = cfg.Token
	client := protocol.NewClient(appwails.NormalizeGatewayAddress(cfg.Gateway, cfg.Port), protoIdentity, identity, connectOptions)
	client.OnInvoke = protocol.Dispatch
	app := appwails.NewApp(dataDir, identity, cfg, client)
	err = wails.Run(&options.App{Title: "OpenClaw Node", Width: 1024, Height: 680, MinWidth: 800, MinHeight: 500, AssetServer: &assetserver.Options{Assets: assets}, BackgroundColour: &options.RGBA{R: 244, G: 250, B: 252, A: 255}, OnStartup: func(ctx context.Context) { app.WailsInit(ctx) }, OnBeforeClose: func(ctx context.Context) bool { wailsruntime.WindowHide(ctx); return false }, Bind: []interface{}{app}, Windows: &windowsoptions.Options{WebviewIsTransparent: false}})
	if err != nil {
		log.Fatal(err)
	}
}
