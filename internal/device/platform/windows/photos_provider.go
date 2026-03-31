package windows

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/openclaw/openclaw-node/internal/device/capabilities/photos"
	"github.com/openclaw/openclaw-node/internal/device/core"
)

type PhotosProvider struct {
	Root string
}

func (p *PhotosProvider) Descriptor() core.CapabilityDescriptor {
	return core.CapabilityDescriptor{Name: "photos", DisplayName: "Photos", Description: "Browse recent local images.", Commands: []string{"photos.latest"}, Tier: 1}
}

func (p *PhotosProvider) Permission() core.PermissionState { return core.PermissionGranted }
func (p *PhotosProvider) Availability() core.Availability  { return core.Availability{Available: true} }

func (p *PhotosProvider) DefaultRoot() string {
	if p.Root != "" {
		return p.Root
	}
	return filepath.Join(os.Getenv("USERPROFILE"), "Pictures", "OpenClaw")
}

func (p *PhotosProvider) List(ctx context.Context, root string) ([]photos.Entry, error) {
	_ = ctx
	if err := os.MkdirAll(root, 0o755); err != nil {
		return nil, core.ErrCapabilityUnavailable
	}
	items, err := os.ReadDir(root)
	if err != nil {
		return nil, core.ErrCapabilityUnavailable
	}
	out := make([]photos.Entry, 0, len(items))
	for _, item := range items {
		info, err := item.Info()
		if err != nil || item.IsDir() {
			continue
		}
		ext := strings.ToLower(filepath.Ext(item.Name()))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" && ext != ".heic" {
			continue
		}
		out = append(out, photos.Entry{
			ID:        fmt.Sprintf("%d", info.ModTime().UnixMilli()),
			Path:      filepath.Join(root, item.Name()),
			CreatedAt: info.ModTime().UnixMilli(),
			Size:      int(info.Size()),
			Format:    ext,
		})
	}
	return out, nil
}
