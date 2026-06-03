//go:build darwin || freebsd || linux || netbsd || openbsd
// +build darwin freebsd linux netbsd openbsd

package runcmd

import (
	"os"
	"os/exec"
	"syscall"
	"testing"
)

var testCommands = []testCommand{
	{
		command: &Command{
			CommandLine: `echo "BAR=${BAR}!"`,
			ForceShell:  "/bin/bash",
			Env:         Env{"BAR": "foo"},
			Logfile:     "out.log",
		},
		successExpected: true,
	},
	{
		command: &Command{
			CommandLine: `command-not-found`,
			Logfile:     "out.log",
		},
		successExpected: false,
	},
	{
		command: &Command{
			Exe:        "/bin/echo",
			Args:       []string{"hello"},
			WorkingDir: "/tmp",
			Logfile:    "out.log",
		},
		successExpected: true,
	},
}

// TestStartWithPresetSysProcAttr exercises the else-branch in start() where
// SysProcAttr is already set on the Cmd before start() is called.
func TestStartWithPresetSysProcAttr(t *testing.T) {
	cmd := exec.Command("/bin/echo", "hello")
	cmd.SysProcAttr = &syscall.SysProcAttr{}
	process, err := start(cmd)
	if err != nil {
		t.Fatalf("start with preset SysProcAttr: %v", err)
	}
	if _, err := process.Wait(); err != nil {
		t.Fatalf("process.Wait: %v", err)
	}
}

// TestStartBuildCmdError covers the early-return path in Start() when
// buildCmd fails (no Exe and no CommandLine configured).
func TestStartBuildCmdError(t *testing.T) {
	if err := (&Command{}).Start(); err == nil {
		t.Fatal("Start() with empty Command should return an error")
	}
}

// TestStartNonExecutable covers the start()-failure path in Start(), including
// the lf.Close() cleanup when the logfile was opened but the process could not
// be launched.
func TestStartNonExecutable(t *testing.T) {
	f, err := os.CreateTemp("", "runcmd-nonexec-*")
	if err != nil {
		t.Fatal(err)
	}
	name := f.Name()
	f.Close()
	defer os.Remove(name)
	// File has no execute bit: os/exec will fail with EACCES.
	// If we are somehow running as root and the exec succeeds, skip rather
	// than fail — the goal is to exercise the error path, not to mandate
	// that root follows permission checks.
	cmd := &Command{
		Exe:     name,
		Logfile: "out.log",
	}
	if err := cmd.Start(); err == nil {
		if cmd.Process != nil {
			cmd.Process.Wait()
		}
		t.Skip("process started despite missing execute bit (running as privileged user?)")
	}
}
