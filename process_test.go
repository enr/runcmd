package runcmd

import (
	"testing"
)

func TestLogfileEmptyPath(t *testing.T) {
	if f := logfile(""); f != nil {
		f.Close()
		t.Fatal("logfile with empty path should return nil")
	}
}

func TestLogfileInvalidPath(t *testing.T) {
	// parent directory does not exist → OpenFile fails → logfile returns nil
	if f := logfile("/runcmd-nonexistent-parent-99999/out.log"); f != nil {
		f.Close()
		t.Fatal("logfile with unwritable path should return nil")
	}
}

type testCommand struct {
	command         *Command
	successExpected bool
}

func TestProcess(t *testing.T) {
	for _, d := range testCommands {
		cmd := d.command
		err := cmd.Start()
		if d.successExpected && err != nil {
			t.Fatalf("%s: success expected but got error %v", cmd, err)
		}
		runningProcess := cmd.Process
		ps, err := runningProcess.Wait()
		if d.successExpected != ps.Success() {
			t.Fatalf("%s: expected success=%t but got %t", cmd, d.successExpected, ps.Success())
		}
		if d.successExpected && err != nil {
			t.Fatalf("%s: success expected but got error %v", cmd, err)
		}
	}
}
