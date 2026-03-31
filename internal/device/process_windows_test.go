//go:build windows

package device

import (
	"os/exec"
	"testing"
)

func TestConfigureBackgroundCommandHidesWindow(t *testing.T) {
	cmd := exec.Command("cmd", "/c", "exit", "0")
	configureBackgroundCommand(cmd)

	if cmd.SysProcAttr == nil {
		t.Fatal("SysProcAttr is nil")
	}
	if !cmd.SysProcAttr.HideWindow {
		t.Fatal("HideWindow = false, want true")
	}
}
