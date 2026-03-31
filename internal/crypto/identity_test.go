package crypto

import (
	"crypto/ed25519"
	"testing"
)

func TestGenerateIdentity(t *testing.T) {
	id, err := GenerateIdentity()
	if err != nil {
		t.Fatalf("GenerateIdentity failed: %v", err)
	}

	if id.DeviceID == "" {
		t.Error("DeviceID should not be empty")
	}

	if len(id.PublicKey) == 0 {
		t.Error("PublicKey should not be empty")
	}

	if len(id.PrivateKey) == 0 {
		t.Error("PrivateKey should not be empty")
	}

	if id.CreatedAtMs == 0 {
		t.Error("CreatedAtMs should not be zero")
	}
}

func TestSignAndVerify(t *testing.T) {
	id, err := GenerateIdentity()
	if err != nil {
		t.Fatalf("GenerateIdentity failed: %v", err)
	}

	payload := "test payload"
	signature := id.Sign(payload)

	if len(signature) != ed25519.SignatureSize {
		t.Errorf("Signature size = %d, want %d", len(signature), ed25519.SignatureSize)
	}
}

func TestSaveAndLoadIdentity(t *testing.T) {
	tmpDir := t.TempDir()
	path := tmpDir + "/identity.json"

	id1, err := GenerateIdentity()
	if err != nil {
		t.Fatalf("GenerateIdentity failed: %v", err)
	}

	if err := SaveIdentity(path, id1); err != nil {
		t.Fatalf("SaveIdentity failed: %v", err)
	}

	id2, err := LoadIdentity(path)
	if err != nil {
		t.Fatalf("LoadIdentity failed: %v", err)
	}

	if id1.DeviceID != id2.DeviceID {
		t.Errorf("DeviceID mismatch: %s != %s", id1.DeviceID, id2.DeviceID)
	}
}

func TestPublicKeyBase64(t *testing.T) {
	id, err := GenerateIdentity()
	if err != nil {
		t.Fatalf("GenerateIdentity failed: %v", err)
	}

	b64 := id.PublicKeyBase64()
	if b64 == "" {
		t.Error("PublicKeyBase64 should not be empty")
	}
}
