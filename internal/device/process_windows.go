//go:build windows

package device

import (
	"os/exec"
	"syscall"
)

func configureBackgroundCommand(cmd *exec.Cmd) {
	if cmd == nil {
		return
	}
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
}
