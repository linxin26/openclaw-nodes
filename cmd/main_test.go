package main

import (
	"testing"

	"github.com/openclaw/openclaw-node/internal/config"
	occrypto "github.com/openclaw/openclaw-node/internal/crypto"
	"github.com/openclaw/openclaw-node/internal/device"
	"github.com/openclaw/openclaw-node/internal/protocol"
)

func TestBuildConnectOptionsAdvertisesRegistryCommands(t *testing.T) {
	identity, err := occrypto.GenerateIdentity()
	if err != nil {
		t.Fatalf("GenerateIdentity failed: %v", err)
	}
	device.Bootstrap(config.Default())
	device.RegisterProtocolHandlers(protocol.RegisterHandler)
	registry := NewRegistry()
	opts := buildConnectOptions(identity, registry, "test-token")
	if opts.Role != "node" {
		t.Fatalf("opts.Role = %q, want node", opts.Role)
	}
	if len(opts.Commands) == 0 {
		t.Fatal("opts.Commands is empty")
	}
	if len(opts.Caps) == 0 {
		t.Fatal("opts.Caps is empty")
	}
	if len(opts.Permissions) == 0 {
		t.Fatal("opts.Permissions is empty")
	}
	if opts.Token != "test-token" {
		t.Fatalf("opts.Token = %q, want test-token", opts.Token)
	}
}

func TestDeviceHandlersAreRegistered(t *testing.T) {
	device.Bootstrap(config.Default())
	device.RegisterProtocolHandlers(protocol.RegisterHandler)
	for _, command := range []string{"camera.snap", "camera.list", "photos.latest", "screen.snapshot"} {
		if _, ok := protocol.GetHandler(command); !ok {
			t.Fatalf("handler %q is not registered", command)
		}
	}
}
