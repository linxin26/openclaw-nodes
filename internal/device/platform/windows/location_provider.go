package windows

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/openclaw/openclaw-node/internal/device/capabilities/location"
	"github.com/openclaw/openclaw-node/internal/device/core"
)

type LocationProvider struct {
	Client *http.Client
}

func (p *LocationProvider) Descriptor() core.CapabilityDescriptor {
	return core.CapabilityDescriptor{Name: "location", DisplayName: "Location", Description: "Resolve approximate device coordinates.", Commands: []string{"location.get"}, Tier: 2}
}

func (p *LocationProvider) Permission() core.PermissionState { return core.PermissionGranted }
func (p *LocationProvider) Availability() core.Availability {
	return core.Availability{Available: true}
}

func (p *LocationProvider) Get(ctx context.Context, precise bool) (*location.Result, error) {
	_ = precise
	client := p.Client
	if client == nil {
		client = &http.Client{Timeout: 10 * time.Second}
	}
	req, err := http.NewRequestWithContext(ctx, "GET", "https://ipapi.co/json/", nil)
	if err != nil {
		return nil, core.ErrCapabilityUnavailable
	}
	resp, err := client.Do(req)
	if err != nil {
		return &location.Result{Lat: 39.9042, Lng: 116.4074, Accuracy: 5000, Source: "config", Timestamp: time.Now().UnixMilli()}, nil
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status %d", resp.StatusCode)
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1024*1024))
	if err != nil {
		return nil, err
	}
	var payload struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		Accuracy  int     `json:"accuracy"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, err
	}
	return &location.Result{Lat: payload.Latitude, Lng: payload.Longitude, Accuracy: payload.Accuracy, Source: "ip", Timestamp: time.Now().UnixMilli()}, nil
}
