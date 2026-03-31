package photos

import (
	"context"

	"github.com/openclaw/openclaw-node/internal/device/core"
)

type Entry struct {
	ID        string `json:"id"`
	Path      string `json:"path"`
	Thumbnail string `json:"thumbnail,omitempty"`
	CreatedAt int64  `json:"createdAt"`
	Size      int    `json:"size"`
	Width     int    `json:"width,omitempty"`
	Height    int    `json:"height,omitempty"`
	Format    string `json:"format"`
}

type LatestRequest struct {
	Limit int   `json:"limit"`
	After int64 `json:"after"`
}

type Provider interface {
	Descriptor() core.CapabilityDescriptor
	Permission() core.PermissionState
	Availability() core.Availability
	DefaultRoot() string
	List(context.Context, string) ([]Entry, error)
}
