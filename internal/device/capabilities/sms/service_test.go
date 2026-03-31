package sms

import "testing"

func TestSMSPlaceholderState(t *testing.T) {
	state := NewService().State(true)
	if state.Availability.Available {
		t.Fatalf("Availability = %#v", state.Availability)
	}
	if state.Permission != "not_supported" {
		t.Fatalf("Permission = %q", state.Permission)
	}
}
