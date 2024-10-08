//go:build darwin || freebsd || linux || netbsd || openbsd
// +build darwin freebsd linux netbsd openbsd

package runcmd

import (
	"os"
	"os/exec"
	"syscall"
)

// ^C on Unix sends a signal to every process in the process group.
// To keep alive the child process you have to put it in a different process group.
// You do that by setting Setpgid to true in the Cmd.SysProcAttr field.
// From Golang docs:
//
//	Setpgid: Set process group ID to new pid (SYSV setpgrp)
func start(cmd *exec.Cmd) (*os.Process, error) {
	keepAliveChild := true
	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{
			Setpgid: keepAliveChild,
		}
	} else {
		cmd.SysProcAttr.Setpgid = keepAliveChild
	}
	err := cmd.Start()
	if err != nil {
		return nil, err
	}
	process := cmd.Process
	return process, nil
}
