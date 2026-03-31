package camera

import (
	"context"
	"testing"
	"time"

	"github.com/openclaw/openclaw-node/internal/device/core"
)

type fakeProvider struct {
	devices []Device
	image   core.ImagePayload
	video   core.VideoPayload
	err     error
}

func (p fakeProvider) Descriptor() core.CapabilityDescriptor {
	return core.CapabilityDescriptor{Name: "camera", Commands: []string{"camera.list", "camera.snap", "camera.clip"}, Tier: 1}
}
func (p fakeProvider) Permission() core.PermissionState       { return core.PermissionGranted }
func (p fakeProvider) Availability() core.Availability        { return core.Availability{Available: true} }
func (p fakeProvider) List(context.Context) ([]Device, error) { return p.devices, p.err }
func (p fakeProvider) Snap(context.Context, SnapRequest) (core.ImagePayload, error) {
	return p.image, p.err
}
func (p fakeProvider) Clip(context.Context, ClipRequest) (core.VideoPayload, error) {
	return p.video, p.err
}

func TestCameraServiceListReturnsNormalizedDevices(t *testing.T) {
	svc := NewService(fakeProvider{devices: []Device{{ID: "1", Name: "Front Cam", Position: "back"}}})
	items, err := svc.List(context.Background())
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if items[0].Position != "rear" {
		t.Fatalf("Position = %q", items[0].Position)
	}
}

func TestCameraServiceSnapReturnsUnifiedImagePayload(t *testing.T) {
	svc := NewService(fakeProvider{image: core.ImagePayload{Base64: "abc", Timestamp: time.Now().UnixMilli()}})
	payload, err := svc.Snap(context.Background(), SnapRequest{})
	if err != nil {
		t.Fatalf("Snap() error = %v", err)
	}
	if payload.Format != "jpeg" || payload.Width != 1920 {
		t.Fatalf("payload = %#v", payload)
	}
}

func TestCameraServiceClipReturnsUnifiedVideoPayload(t *testing.T) {
	svc := NewService(fakeProvider{video: core.VideoPayload{Base64: "abc", Timestamp: time.Now().UnixMilli()}})
	payload, err := svc.Clip(context.Background(), ClipRequest{})
	if err != nil {
		t.Fatalf("Clip() error = %v", err)
	}
	if payload.Format != "mp4" || payload.DurationMs != 5000 {
		t.Fatalf("payload = %#v", payload)
	}
}

func TestCameraServiceMapsUnavailableProviderErrors(t *testing.T) {
	svc := NewService(fakeProvider{err: core.ErrCapabilityUnavailable})
	_, err := svc.List(context.Background())
	if err != core.ErrCapabilityUnavailable {
		t.Fatalf("err = %v", err)
	}
}
