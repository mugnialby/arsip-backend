//go:build windows
// +build windows

package utils

import (
	"os/exec"
	"syscall"
)

func ApplySysProcAttr(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
	}
}
