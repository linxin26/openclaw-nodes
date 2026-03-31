package protocol

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/websocket"
)

func TestSendInvokeResultWritesNodeInvokeResultRequest(t *testing.T) {
	upgrader := websocket.Upgrader{}
	frameCh := make(chan Frame, 1)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Errorf("upgrade failed: %v", err)
			return
		}
		defer conn.Close()
		var frame Frame
		if err := conn.ReadJSON(&frame); err == nil {
			frameCh <- frame
		}
	}))
	defer server.Close()

	wsURL := "ws" + server.URL[len("http"):]
	clientConn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("dial failed: %v", err)
	}
	defer clientConn.Close()

	client := &Client{conn: clientConn}
	req := &InvokeRequest{
		ID:      "invoke-1",
		NodeID:  "node-1",
		Command: "camera.snap",
	}
	result, err := NewInvokeResultOK(map[string]interface{}{"handled": true})
	if err != nil {
		t.Fatalf("NewInvokeResultOK failed: %v", err)
	}

	client.sendInvokeResult(req, result)

	frame := <-frameCh
	if frame.Type != "req" {
		t.Fatalf("frame.Type = %q, want req", frame.Type)
	}
	if frame.Method != "node.invoke.result" {
		t.Fatalf("frame.Method = %q, want node.invoke.result", frame.Method)
	}

	var params map[string]interface{}
	if err := json.Unmarshal(frame.Params, &params); err != nil {
		t.Fatalf("unmarshal params failed: %v", err)
	}
	if params["id"] != "invoke-1" {
		t.Fatalf("params.id = %#v, want invoke-1", params["id"])
	}
	if params["nodeId"] != "node-1" {
		t.Fatalf("params.nodeId = %#v, want node-1", params["nodeId"])
	}
	if params["ok"] != true {
		t.Fatalf("params.ok = %#v, want true", params["ok"])
	}
	payload, ok := params["payload"].(map[string]interface{})
	if !ok || payload["handled"] != true {
		t.Fatalf("params.payload = %#v, want handled=true", params["payload"])
	}
}
