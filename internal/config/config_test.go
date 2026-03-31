package config

import (
	"path/filepath"
	"testing"
)

func TestDefault(t *testing.T) {
	cfg := Default()
	if cfg.Port != 18789 {
		t.Errorf("Expected Port = 18789, got %d", cfg.Port)
	}
	if cfg.Discovery != "auto" {
		t.Errorf("Expected Discovery = 'auto', got '%s'", cfg.Discovery)
	}
	if cfg.Capabilities == nil {
		t.Fatal("Expected Capabilities to not be nil")
	}
	if !cfg.Capabilities["camera"] {
		t.Error("Expected camera capability to be true")
	}
	if cfg.CapabilityOptions == nil {
		t.Fatal("Expected CapabilityOptions to not be nil")
	}
}

func TestLoadNonExistent(t *testing.T) {
	_, err := Load("/nonexistent/path/config.yaml")
	if err == nil {
		t.Error("Expected error for non-existent file")
	}
}

func TestLoadAndSave(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "config.yaml")
	cfg := Default()
	cfg.Gateway = "localhost:9999"
	cfg.CapabilityOptions["photos"] = CapabilityOption{Provider: "windows", Path: filepath.Join(tmpDir, "photos")}
	if err := cfg.Save(path); err != nil {
		t.Fatalf("Save failed: %v", err)
	}
	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if loaded.Gateway != cfg.Gateway {
		t.Errorf("Gateway mismatch: got %s, want %s", loaded.Gateway, cfg.Gateway)
	}
	if loaded.Port != cfg.Port {
		t.Errorf("Port mismatch: got %d, want %d", loaded.Port, cfg.Port)
	}
	if !loaded.Capabilities["camera"] {
		t.Fatal("camera capability should remain enabled")
	}
	if loaded.CapabilityOptions["photos"].Path == "" {
		t.Fatal("photos path option not persisted")
	}
}
