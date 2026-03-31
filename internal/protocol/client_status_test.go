package protocol_test

import (
	"encoding/json"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/openclaw/openclaw-node/internal/config"
	"github.com/openclaw/openclaw-node/internal/device"
	"github.com/openclaw/openclaw-node/internal/protocol"
)

func TestClientIsConnected(t *testing.T) {
	client := &protocol.Client{}
	if client.IsConnected() {
		t.Fatal("IsConnected() = true, want false for zero client")
	}
	clientPtr := client
	_ = clientPtr
}

func TestClientSetServerURL(t *testing.T) {
	client := &protocol.Client{}
	client.SetServerURL("localhost:38789")
}

func TestClientSetToken(t *testing.T) {
	client := &protocol.Client{}
	client.SetToken("new-token")
}

func TestClientStatusUsesRuntimeRegistry(t *testing.T) {
	device.Bootstrap(config.Default())
	device.RegisterProtocolHandlers(protocol.RegisterHandler)
	result := protocol.Dispatch(protocol.InvokeRequest{ID: "status-1", Command: "device.status"})
	if !result.OK {
		t.Fatalf("device.status error = %#v", result.Error)
	}
	var payload struct {
		Capabilities map[string]struct {
			Enabled   bool   `json:"enabled"`
			Available bool   `json:"available"`
			Reason    string `json:"reason"`
		} `json:"capabilities"`
	}
	if err := json.Unmarshal(result.Payload, &payload); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	if _, ok := payload.Capabilities["camera"]; !ok {
		t.Fatal("camera capability missing")
	}
	if payload.Capabilities["motion"].Available {
		t.Fatal("motion should be unavailable")
	}
}

func TestClientStatusReportsDynamicPermissions(t *testing.T) {
	device.Bootstrap(config.Default())
	result := protocol.Dispatch(protocol.InvokeRequest{ID: "perm-1", Command: "device.permissions"})
	if !result.OK {
		t.Fatalf("device.permissions error = %#v", result.Error)
	}
	var permissions map[string]string
	if err := json.Unmarshal(result.Payload, &permissions); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	if permissions["camera"] != "granted" {
		t.Fatalf("camera permission = %q", permissions["camera"])
	}
	if permissions["sms"] != "not_supported" {
		t.Fatalf("sms permission = %q", permissions["sms"])
	}
}

var _ = websocket.Conn{}
