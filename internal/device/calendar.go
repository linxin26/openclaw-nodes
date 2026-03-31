package device

import (
	"context"
	"encoding/json"

	capcalendar "github.com/openclaw/openclaw-node/internal/device/capabilities/calendar"
	"github.com/openclaw/openclaw-node/internal/protocol"
)

type CalendarEvent = capcalendar.Event

func handleCalendarEvents(params json.RawMessage) (*protocol.InvokeResult, error) {
	var args capcalendar.EventsRequest
	if err := decodeParams(params, &args); err != nil {
		return protocol.NewInvokeResultError("INVALID_PARAMS", "Invalid JSON: "+err.Error())
	}
	items, err := Runtime().Calendar().Events(context.Background(), args)
	if err != nil {
		return invokeError(err, "calendar unavailable")
	}
	return protocol.NewInvokeResultOK(map[string]interface{}{"events": items})
}

func handleCalendarAdd(params json.RawMessage) (*protocol.InvokeResult, error) {
	var args capcalendar.AddRequest
	if err := decodeParams(params, &args); err != nil {
		return protocol.NewInvokeResultError("INVALID_PARAMS", "Invalid JSON: "+err.Error())
	}
	result, err := Runtime().Calendar().Add(context.Background(), args)
	if err != nil {
		return invokeError(err, "cannot write calendar")
	}
	return protocol.NewInvokeResultOK(result)
}
