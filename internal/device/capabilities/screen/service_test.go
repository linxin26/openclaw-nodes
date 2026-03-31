package screen

import (
	"context"
	"testing"
	"time"

	"github.com/openclaw/openclaw-node/internal/device/core"
)

type fakeProvider struct {
	payload core.ImagePayload
	err     error
}

func (p fakeProvider) Descriptor() core.CapabilityDescriptor {
	return core.CapabilityDescriptor{Name: "screen", Commands: []string{"screen.snapshot"}, Tier: 1}
}
func (p fakeProvider) Permission() core.PermissionState { return core.PermissionGranted }
func (p fakeProvider) Availability() core.Availability  { return core.Availability{Available: true} }
func (p fakeProvider) Snapshot(context.Context, SnapshotRequest) (core.ImagePayload, error) {
	return p.payload, p.err
}

func TestScreenServiceNormalizesSnapshotRequest(t *testing.T) {
	svc := NewService(fakeProvider{payload: core.ImagePayload{Base64: "abc", Timestamp: time.Now().UnixMilli()}})
	payload, err := svc.Snapshot(context.Background(), SnapshotRequest{})
	if err != nil {
		t.Fatalf("Snapshot() error = %v", err)
	}
	if payload.Format != "png" || payload.Width != 1920 {
		t.Fatalf("payload = %#v", payload)
	}
}

func TestScreenServiceMapsProviderErrors(t *testing.T) {
	svc := NewService(fakeProvider{err: core.ErrCapabilityUnavailable})
	_, err := svc.Snapshot(context.Background(), SnapshotRequest{})
	if err != core.ErrCapabilityUnavailable {
		t.Fatalf("err = %v", err)
	}
}
