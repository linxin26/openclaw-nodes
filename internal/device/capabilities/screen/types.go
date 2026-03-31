package screen

import (
	"context"

	"github.com/openclaw/openclaw-node/internal/device/core"
)

type SnapshotRequest struct {
	Format   string `json:"format"`
	Quality  int    `json:"quality"`
	MaxWidth int    `json:"maxWidth"`
}

type Provider interface {
	Descriptor() core.CapabilityDescriptor
	Permission() core.PermissionState
	Availability() core.Availability
	Snapshot(context.Context, SnapshotRequest) (core.ImagePayload, error)
}
