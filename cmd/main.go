//go:build !wails
// +build !wails

package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/openclaw/openclaw-node/internal/config"
	"github.com/openclaw/openclaw-node/internal/crypto"
	"github.com/openclaw/openclaw-node/internal/device"
	deviceruntime "github.com/openclaw/openclaw-node/internal/device/runtime"
	"github.com/openclaw/openclaw-node/internal/discovery"
	"github.com/openclaw/openclaw-node/internal/protocol"
	"github.com/openclaw/openclaw-node/internal/tray"
	appwails "github.com/openclaw/openclaw-node/internal/wails"
	"github.com/openclaw/openclaw-node/store"
)

var (
	flagGateway = flag.String("gateway", "", "Gateway address (host:port)")
	flagToken   = flag.String("token", "", "Gateway auth token")
	flagTLS     = flag.Bool("tls", false, "Use TLS")
	flagNoMdns  = flag.Bool("no-mdns", false, "Disable mDNS discovery")
)

func main() {
	flag.Parse()
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
		crypto.SaveIdentity(identityPath, identity)
	}
	cfgPath := s.Path("config.yaml")
	cfg, _ := config.Load(cfgPath)
	if cfg == nil {
		cfg = config.Default()
	}
	if *flagGateway != "" {
		cfg.Gateway = *flagGateway
		cfg.Discovery = "manual"
	}
	if *flagTLS {
		cfg.TLS = true
	}
	if *flagNoMdns {
		cfg.Discovery = "manual"
	}
	device.Bootstrap(cfg)
	device.RegisterProtocolHandlers(protocol.RegisterHandler)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var mdns *discovery.MDNS
	if cfg.Discovery == "auto" {
		mdns = discovery.NewMDNS(identity.DeviceID, 18789)
		go mdns.Register(ctx)
	}
	protoIdentity := &protocol.Identity{DeviceID: identity.DeviceID, Role: "node"}
	registry := NewRegistry()
	protocol.GlobalProtocol.Identity = &protocol.Identity{DeviceID: identity.DeviceID, ClientID: "node-host", ClientMode: "node", Role: "node", SignedAtMs: identity.CreatedAtMs}
	client := protocol.NewClient(appwails.NormalizeGatewayAddress(cfg.Gateway, cfg.Port), protoIdentity, identity, buildConnectOptions(identity, registry, *flagToken))
	trayInstance := tray.NewWithRuntime(deviceruntime.Default())
	go trayInstance.Run()
	client.OnConnected = func(resp protocol.ConnectResponse) {
		log.Printf("Connected to gateway")
		trayInstance.SetStatus(tray.StatusConnected)
	}
	client.OnDisconnected = func(reason string) {
		log.Printf("Disconnected: %s", reason)
		trayInstance.SetStatus(tray.StatusOffline)
	}
	client.OnInvoke = protocol.Dispatch
	if cfg.Gateway != "" {
		trayInstance.SetStatus(tray.StatusConnecting)
		if err := client.Connect(ctx); err != nil {
			log.Printf("Failed to connect: %v", err)
		}
	}
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	log.Println("Shutting down...")
	client.Disconnect()
}

func buildConnectOptions(identity *crypto.Identity, registry *CommandRegistry, token string) protocol.ConnectOptions {
	meta := deviceruntime.Default().Metadata()
	return protocol.ConnectOptions{Role: "node", Scopes: []string{}, Caps: registry.AllCaps(), Commands: registry.AllCommands(), Permissions: registry.permissions, Token: token, Client: protocol.ClientInfo{ID: "node-host", DisplayName: "OpenClaw Node", Version: meta.Version, Platform: meta.Platform, Mode: "node", InstanceID: identity.DeviceID[:8], DeviceFamily: "desktop", ModelIdentifier: meta.ModelIdentifier}}
}
