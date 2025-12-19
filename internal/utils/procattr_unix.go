//go:build linux || darwin
// +build linux darwin

package utils

import "os/exec"

func ApplySysProcAttr(cmd *exec.Cmd) {
	// no-op on Unix systems
}
