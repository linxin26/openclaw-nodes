//go:build !windows

package device

import "os/exec"

func configureBackgroundCommand(cmd *exec.Cmd) {
}
