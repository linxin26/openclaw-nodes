package motion

import "testing"

func TestMotionPlaceholderState(t *testing.T) {
	state := NewService().State(true)
	if state.Availability.Available {
		t.Fatalf("Availability = %#v", state.Availability)
	}
	if state.Permission != "not_applicable" {
		t.Fatalf("Permission = %q", state.Permission)
	}
}
