package windows

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/openclaw/openclaw-node/internal/device/capabilities/calendar"
	"github.com/openclaw/openclaw-node/internal/device/core"
)

type CalendarProvider struct {
	Path string
}

func (p *CalendarProvider) Descriptor() core.CapabilityDescriptor {
	return core.CapabilityDescriptor{Name: "calendar", DisplayName: "Calendar", Description: "Browse and create local calendar events.", Commands: []string{"calendar.events", "calendar.add"}, Tier: 2}
}

func (p *CalendarProvider) Permission() core.PermissionState { return core.PermissionGranted }
func (p *CalendarProvider) Availability() core.Availability {
	return core.Availability{Available: true}
}

func (p *CalendarProvider) Events(ctx context.Context, req calendar.EventsRequest) ([]calendar.Event, error) {
	_ = ctx
	data, err := os.ReadFile(p.calendarPath())
	if err != nil {
		return nil, core.ErrCapabilityUnavailable
	}
	lines := strings.Split(string(data), "\n")
	items := make([]calendar.Event, 0, req.Limit)
	for _, line := range lines {
		if strings.HasPrefix(line, "SUMMARY:") {
			items = append(items, calendar.Event{
				ID:    fmt.Sprintf("%d", len(items)+1),
				Title: strings.TrimPrefix(line, "SUMMARY:"),
				Start: time.Now().UnixMilli(),
				End:   time.Now().Add(time.Hour).UnixMilli(),
			})
		}
		if len(items) >= req.Limit {
			break
		}
	}
	return items, nil
}

func (p *CalendarProvider) Add(ctx context.Context, req calendar.AddRequest) (*calendar.AddResult, error) {
	_ = ctx
	if req.Title == "" || req.Start == 0 {
		return nil, core.ErrCapabilityUnavailable
	}
	path := p.calendarPath()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, core.ErrCapabilityUnavailable
	}
	id := uuid.New().String()
	event := fmt.Sprintf("BEGIN:VEVENT\nUID:%s\nDTSTART:%s\nDTEND:%s\nSUMMARY:%s\nDESCRIPTION:%s\nLOCATION:%s\nEND:VEVENT\n", id, time.UnixMilli(req.Start).Format("20060102T150405Z"), time.UnixMilli(req.End).Format("20060102T150405Z"), req.Title, req.Description, req.Location)
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, core.ErrCapabilityUnavailable
	}
	defer f.Close()
	if _, err := f.WriteString(event); err != nil {
		return nil, core.ErrCapabilityUnavailable
	}
	return &calendar.AddResult{Success: true, EventID: id}, nil
}

func (p *CalendarProvider) calendarPath() string {
	if p.Path != "" {
		return p.Path
	}
	return filepath.Join(os.Getenv("USERPROFILE"), "Documents", "OpenClaw", "calendar.ics")
}
