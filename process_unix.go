// +build darwin freebsd linux netbsd openbsd

package runcmd

import (
	"os/exec"
	"syscall"
)

func start(cmd *exec.Cmd) error {
	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{
			Setpgid: true,
		}
	} else {
		cmd.SysProcAttr.Setpgid = true
	}
	err := cmd.Start()
	return err
}
