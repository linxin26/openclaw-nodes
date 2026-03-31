package windows

import "testing"

func TestGuessCameraPosition(t *testing.T) {
	if got := guessCameraPosition("Front Webcam"); got != "front" {
		t.Fatalf("guessCameraPosition() = %q", got)
	}
}
