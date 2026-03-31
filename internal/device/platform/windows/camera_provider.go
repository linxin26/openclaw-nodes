package windows

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/openclaw/openclaw-node/internal/device/capabilities/camera"
	"github.com/openclaw/openclaw-node/internal/device/core"
)

type CameraProvider struct{}

func (p *CameraProvider) Descriptor() core.CapabilityDescriptor {
	return core.CapabilityDescriptor{Name: "camera", DisplayName: "Camera", Description: "List cameras and capture media.", Commands: []string{"camera.list", "camera.snap", "camera.clip"}, Tier: 1}
}

func (p *CameraProvider) Permission() core.PermissionState { return core.PermissionGranted }
func (p *CameraProvider) Availability() core.Availability  { return core.Availability{Available: true} }

func (p *CameraProvider) List(ctx context.Context) ([]camera.Device, error) {
	cmd := exec.CommandContext(ctx, "ffmpeg", "-list_devices", "true", "-f", "dshow", "-i", "dummy")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, core.ErrCapabilityUnavailable
	}
	return parseFFmpegDevices(string(output)), nil
}

func (p *CameraProvider) Snap(ctx context.Context, req camera.SnapRequest) (core.ImagePayload, error) {
	tmpFile := fmt.Sprintf("%s/camera_snap_%d.jpg", os.TempDir(), time.Now().UnixMilli())
	cmd := exec.CommandContext(ctx, "ffmpeg", "-f", "dshow", "-i", fmt.Sprintf("video=%s", p.getCameraName(req.CameraID)), "-vframes", "1", "-q:v", "2", tmpFile)
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
		Format:    "jpeg",
		Width:     req.MaxWidth,
		Size:      len(data),
		Timestamp: time.Now().UnixMilli(),
	}, nil
}

func (p *CameraProvider) Clip(ctx context.Context, req camera.ClipRequest) (core.VideoPayload, error) {
	tmpFile := fmt.Sprintf("%s/camera_clip_%d.mp4", os.TempDir(), time.Now().UnixMilli())
	cmd := exec.CommandContext(ctx, "ffmpeg", "-y", "-f", "dshow", "-t", fmt.Sprintf("%0.2f", float64(req.DurationMs)/1000), "-i", fmt.Sprintf("video=%s", p.getCameraName(req.CameraID)), tmpFile)
	if err := cmd.Run(); err != nil {
		return core.VideoPayload{}, core.ErrCapabilityUnavailable
	}
	defer os.Remove(tmpFile)
	data, err := os.ReadFile(tmpFile)
	if err != nil {
		return core.VideoPayload{}, core.ErrCapabilityUnavailable
	}
	return core.VideoPayload{
		Base64:       base64.StdEncoding.EncodeToString(data),
		Format:       "mp4",
		Width:        req.MaxWidth,
		Size:         len(data),
		DurationMs:   req.DurationMs,
		Timestamp:    time.Now().UnixMilli(),
		IncludeAudio: req.IncludeAudio,
	}, nil
}

func (p *CameraProvider) getCameraName(id string) string {
	if id == "" || id == "0" {
		return "Integrated Camera"
	}
	return "USB Camera"
}

func parseFFmpegDevices(output string) []camera.Device {
	var cameras []camera.Device
	lines := strings.Split(output, "\n")
	inVideoSection := false
	for _, line := range lines {
		if strings.Contains(line, "DirectShow video devices") {
			inVideoSection = true
			continue
		}
		if strings.Contains(line, "DirectShow audio devices") {
			break
		}
		if !inVideoSection {
			continue
		}
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "\"") && strings.HasSuffix(line, "\"") {
			name := strings.Trim(line, "\"")
			if name != "" {
				cameras = append(cameras, camera.Device{
					ID:       fmt.Sprintf("%d", len(cameras)),
					Name:     name,
					Position: guessCameraPosition(name),
				})
			}
		}
	}
	return cameras
}

func guessCameraPosition(name string) string {
	lower := strings.ToLower(name)
	switch {
	case strings.Contains(lower, "front"), strings.Contains(lower, "webcam"):
		return "front"
	case strings.Contains(lower, "back"), strings.Contains(lower, "rear"):
		return "rear"
	default:
		return "external"
	}
}
