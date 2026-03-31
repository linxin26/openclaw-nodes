package location

import (
	"context"
	"testing"

	"github.com/openclaw/openclaw-node/internal/device/core"
)

type fakeProvider struct {
	result *Result
	err    error
}

func (p fakeProvider) Descriptor() core.CapabilityDescriptor {
	return core.CapabilityDescriptor{Name: "location", Commands: []string{"location.get"}, Tier: 2}
}
func (p fakeProvider) Permission() core.PermissionState           { return core.PermissionGranted }
func (p fakeProvider) Availability() core.Availability            { return core.Availability{Available: true} }
func (p fakeProvider) Get(context.Context, bool) (*Result, error) { return p.result, p.err }

func TestLocationServiceMapsPermissionState(t *testing.T) {
	svc := NewService(fakeProvider{result: &Result{Lat: 1, Lng: 2}})
	result, err := svc.Get(context.Background(), false)
	if err != nil || result.Lat != 1 {
		t.Fatalf("result=%#v err=%v", result, err)
	}
}

func TestLocationServiceMapsUnavailableError(t *testing.T) {
	svc := NewService(fakeProvider{err: core.ErrCapabilityUnavailable})
	_, err := svc.Get(context.Background(), false)
	if err != core.ErrCapabilityUnavailable {
		t.Fatalf("err = %v", err)
	}
}
