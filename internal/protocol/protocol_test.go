package protocol

import (
	"encoding/json"
	"testing"
)

func TestNewInvokeResultOK(t *testing.T) {
	payload := map[string]interface{}{"key": "value"}
	result, err := NewInvokeResultOK(payload)

	if err != nil {
		t.Fatalf("NewInvokeResultOK failed: %v", err)
	}

	if !result.OK {
		t.Error("Expected OK to be true")
	}

	if result.Error != nil {
		t.Error("Expected Error to be nil")
	}
}

func TestNewInvokeResultError(t *testing.T) {
	result, err := NewInvokeResultError("TEST_ERROR", "test message")

	if err != nil {
		t.Fatalf("NewInvokeResultError failed: %v", err)
	}

	if result.OK {
		t.Error("Expected OK to be false")
	}

	if result.Error == nil {
		t.Fatal("Expected Error to not be nil")
	}

	if result.Error.Code != "TEST_ERROR" {
		t.Errorf("Expected Error.Code = 'TEST_ERROR', got '%s'", result.Error.Code)
	}

	if result.Error.Message != "test message" {
		t.Errorf("Expected Error.Message = 'test message', got '%s'", result.Error.Message)
	}
}

func TestRegisterHandler(t *testing.T) {
	handler := func(params json.RawMessage) (*InvokeResult, error) {
		return NewInvokeResultOK(nil)
	}

	RegisterHandler("test.command", handler)

	if _, ok := getGlobalProtocol().GetHandler("test.command"); !ok {
		t.Error("Handler was not registered")
	}
}

func TestDispatch(t *testing.T) {
	RegisterHandler("test.dispatch", func(params json.RawMessage) (*InvokeResult, error) {
		return NewInvokeResultOK(map[string]interface{}{"dispatched": true})
	})

	req := InvokeRequest{
		ID:      "test-1",
		Command: "test.dispatch",
	}

	result := Dispatch(req)

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if !result.OK {
		t.Error("Expected result.OK to be true")
	}
}

func TestDispatchUnknownCommand(t *testing.T) {
	req := InvokeRequest{
		ID:      "test-2",
		Command: "unknown.command",
	}

	result := Dispatch(req)

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result.OK {
		t.Error("Expected result.OK to be false")
	}

	if result.Error == nil {
		t.Fatal("Expected Error to not be nil")
	}

	if result.Error.Code != "INVALID_REQUEST" {
		t.Errorf("Expected Error.Code = 'INVALID_REQUEST', got '%s'", result.Error.Code)
	}
}
