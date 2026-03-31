package crypto

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Identity struct {
	DeviceID    string `json:"deviceId"`
	PublicKey   []byte `json:"publicKey"`
	PrivateKey  []byte `json:"privateKey"` // PKCS8 format
	CreatedAtMs int64  `json:"createdAtMs"`
}

func (i *Identity) Sign(payload string) []byte {
	privateKey := ed25519.PrivateKey(i.PrivateKey)
	return ed25519.Sign(privateKey, []byte(payload))
}

func (i *Identity) PublicKeyBase64() string {
	return base64.RawURLEncoding.EncodeToString(i.PublicKey)
}

func GenerateIdentity() (*Identity, error) {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}

	hash := sha256.Sum256(publicKey)
	deviceID := fmt.Sprintf("%x", hash)

	// Convert to PKCS8
	pkcs8 := exportPrivateKey(privateKey)

	return &Identity{
		DeviceID:    deviceID,
		PublicKey:   publicKey,
		PrivateKey:  pkcs8,
		CreatedAtMs: now(),
	}, nil
}

func LoadIdentity(path string) (*Identity, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var id Identity
	if err := json.Unmarshal(data, &id); err != nil {
		return nil, err
	}
	return &id, nil
}

func SaveIdentity(path string, id *Identity) error {
	data, err := json.MarshalIndent(id, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

func now() int64 {
	return time.Now().UnixMilli()
}

// exportPrivateKey returns Ed25519 private key bytes.
// Ed25519 keys are used directly in raw format for OpenClaw.
func exportPrivateKey(key ed25519.PrivateKey) []byte {
	return key
}
