package runcmd

import (
	"testing"
)

type testCommand struct {
	command *Command
	success bool
}

func TestProcess(t *testing.T) {
	for _, d := range testCommands {
		cmd := d.command
		err := cmd.Start()
		runningProcess := cmd.Process
		if d.success && err != nil {
			t.Fatalf("%s: success expected but got error %v", cmd, err)
		}
		if !d.success && err == nil {
			t.Fatalf("%s: fail expected but got nil error", cmd)
		}
		if d.success && runningProcess == nil {
			t.Fatalf("%s: success expected but got nil process", cmd)
		}
	}
}
