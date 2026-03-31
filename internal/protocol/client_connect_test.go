package protocol

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	occrypto "github.com/openclaw/openclaw-node/internal/crypto"
)

func TestBuildConnectParamsIncludesSignedClientMetadata(t *testing.T) {
	identity, err := occrypto.GenerateIdentity()
	if err != nil {
		t.Fatalf("GenerateIdentity failed: %v", err)
	}

	client := NewClient("ws://example.test", &Identity{DeviceID: identity.DeviceID, Role: "node"}, identity, ConnectOptions{
		Role:   "node",
		Scopes: []string{"node"},
		Token:  "test-token",
		Client: ClientInfo{
			ID:              "node-host",
			DisplayName:     "OpenClaw Node",
			Version:         "0.1.0",
			Platform:        "windows",
			Mode:            "node",
			InstanceID:      "test-instance",
			DeviceFamily:    "desktop",
			ModelIdentifier: "openclaw-node-windows",
		},
	})

	params := client.buildConnectParams(1_700_000_000_000, "test-nonce")
	clientMap, ok := params["client"].(map[string]interface{})
	if !ok {
		t.Fatalf("client params missing or wrong type: %#v", params["client"])
	}

	if got := clientMap["deviceFamily"]; got != "desktop" {
		t.Fatalf("client.deviceFamily = %#v, want %q", got, "desktop")
	}
	if got := clientMap["displayName"]; got != "OpenClaw Node" {
		t.Fatalf("client.displayName = %#v, want %q", got, "OpenClaw Node")
	}
	if got := clientMap["modelIdentifier"]; got != "openclaw-node-windows" {
		t.Fatalf("client.modelIdentifier = %#v, want %q", got, "openclaw-node-windows")
	}

	deviceMap, ok := params["device"].(map[string]interface{})
	if !ok {
		t.Fatalf("device params missing or wrong type: %#v", params["device"])
	}
	if got := deviceMap["nonce"]; got != "test-nonce" {
		t.Fatalf("device.nonce = %#v, want %q", got, "test-nonce")
	}
}

func TestClientCanReconnectAfterDisconnect(t *testing.T) {
	identity, err := occrypto.GenerateIdentity()
	if err != nil {
		t.Fatalf("GenerateIdentity failed: %v", err)
	}

	upgrader := websocket.Upgrader{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Errorf("upgrade failed: %v", err)
			return
		}
		defer conn.Close()

		challengePayload, _ := json.Marshal(map[string]string{"nonce": "test-nonce"})
		if err := conn.WriteJSON(Frame{
			Type:    "event",
			Event:   "connect.challenge",
			Payload: challengePayload,
		}); err != nil {
			t.Errorf("write challenge failed: %v", err)
			return
		}

		var req Frame
		if err := conn.ReadJSON(&req); err != nil {
			t.Errorf("read connect request failed: %v", err)
			return
		}

		connectPayload, _ := json.Marshal(ConnectResponse{OK: true})
		if err := conn.WriteJSON(Frame{
			Type:    "res",
			ID:      req.ID,
			OK:      true,
			Payload: connectPayload,
		}); err != nil {
			t.Errorf("write connect response failed: %v", err)
			return
		}

		time.Sleep(100 * time.Millisecond)
	}))
	defer server.Close()

	wsURL := "ws" + server.URL[len("http"):]
	client := NewClient(wsURL, &Identity{DeviceID: identity.DeviceID, Role: "node"}, identity, ConnectOptions{
		Role:   "node",
		Scopes: []string{"node"},
		Client: ClientInfo{
			ID:           "node-host",
			Version:      "0.1.0",
			Platform:     "windows",
			Mode:         "node",
			InstanceID:   "test-instance",
			DeviceFamily: "desktop",
		},
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		t.Fatalf("first Connect() error = %v", err)
	}
	if !client.IsConnected() {
		t.Fatal("IsConnected() = false after first connect, want true")
	}

	client.Disconnect()
	if client.IsConnected() {
		t.Fatal("IsConnected() = true after disconnect, want false")
	}

	ctx2, cancel2 := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel2()

	if err := client.Connect(ctx2); err != nil {
		t.Fatalf("second Connect() error = %v", err)
	}
	if !client.IsConnected() {
		t.Fatal("IsConnected() = false after reconnect, want true")
	}
}
