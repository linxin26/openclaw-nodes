package wails

import (
	"context"

	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

func emit(ctx context.Context, event string, payload interface{}) {
	if ctx == nil {
		return
	}
	wailsruntime.EventsEmit(ctx, event, payload)
}

func EmitStatusChange(ctx context.Context, status *ConnectionStatus) {
	emit(ctx, "status:change", status)
}

func EmitCapabilityChange(ctx context.Context, capability *CapabilityInfo) {
	emit(ctx, "capability:change", capability)
}

func EmitLog(ctx context.Context, entry *LogEntry) {
	emit(ctx, "log", entry)
}

func EmitActivity(ctx context.Context, entry *ActivityEntry) {
	emit(ctx, "activity", entry)
}

func EmitInvokeComplete(ctx context.Context, method string, success bool, durationMs int64) {
	emit(ctx, "invoke:complete", map[string]interface{}{
		"method":     method,
		"success":    success,
		"durationMs": durationMs,
	})
}

func EmitConfigChange(ctx context.Context, cfg *Config) {
	emit(ctx, "config:change", cfg)
}
