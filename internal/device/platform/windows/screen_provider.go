package windows

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/openclaw/openclaw-node/internal/device/capabilities/screen"
	"github.com/openclaw/openclaw-node/internal/device/core"
)

type ScreenProvider struct{}

func (p *ScreenProvider) Descriptor() core.CapabilityDescriptor {
	return core.CapabilityDescriptor{Name: "screen", DisplayName: "Screen", Description: "Capture desktop screenshots.", Commands: []string{"screen.snapshot"}, Tier: 1}
}

func (p *ScreenProvider) Permission() core.PermissionState { return core.PermissionGranted }
func (p *ScreenProvider) Availability() core.Availability  { return core.Availability{Available: true} }

func (p *ScreenProvider) Snapshot(ctx context.Context, req screen.SnapshotRequest) (core.ImagePayload, error) {
	tmpFile := fmt.Sprintf("%s/screen_%d.png", os.TempDir(), time.Now().UnixMilli())
	script := fmt.Sprintf(`Add-Type -AssemblyName System.Windows.Forms; $bmp = New-Object System.Drawing.Bitmap([System.Windows.Forms.Screen]::PrimaryScreen.Bounds.Width, [System.Windows.Forms.Screen]::PrimaryScreen.Bounds.Height); $graphics = [System.Drawing.Graphics]::FromImage($bmp); $graphics.CopyFromScreen(0, 0, 0, 0, $bmp.Size); $bmp.Save("%s", [System.Drawing.Imaging.ImageFormat]::Png); $bmp.Dispose()`, tmpFile)
	cmd := exec.CommandContext(ctx, "powershell", "-Command", script)
	if err := cmd.Run(); err != nil {
		return core.ImagePayload{}, core.ErrCapabilityUnavailable
	}
	defer os.Remove(tmpFile)
	data, err := os.ReadFile(tmpFile)
	if err != nil {
		return core.ImagePayload{}, core.ErrCapabilityUnavailable
	}
	return core.ImagePayload{
		Base64:    base64.StdEncoding.EncodeToString(data),
		Format:    req.Format,
		Width:     req.MaxWidth,
		Size:      len(data),
		Timestamp: time.Now().UnixMilli(),
	}, nil
}
