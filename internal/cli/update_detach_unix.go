//go:build !windows

package cli

import (
	"os/exec"
	"syscall"
)

func configureDetachedUpdateRefresh(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
}
