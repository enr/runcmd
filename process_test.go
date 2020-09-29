package runcmd

import (
	"testing"
)

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
	}
}
