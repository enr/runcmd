package runcmd

import (
	"path/filepath"
	"strings"
	"testing"
)

type cnl struct {
	command *Command
	name    string
	logfile string
}

func TestGetName(t *testing.T) {
	for _, d := range testdata2 {
		cmd := d.command
		actual := cmd.GetName()
		expected := d.name
		if actual != expected {
			t.Fatalf("%s: name expected %s but got %s", cmd, expected, actual)
		}
	}
}

func TestGetLogfile(t *testing.T) {
	for _, d := range testdata2 {
		cmd := d.command
		actual := filepath.Base(cmd.GetLogfile())
		expected := d.logfile
		if actual != expected {
			t.Fatalf("%s: logfile base expected %s but got %s", cmd, expected, actual)
		}
	}
}

func TestHugeLogName(t *testing.T) {
	command :=
		&Command{
			Exe:    `/usr/local/bin/myapp`,
			Args:   []string{"-a", `"a very long command line"`, strings.Repeat("addanotherarg", 20)},
			UseEnv: true,
		}
	lf := command.GetLogfile()
	if len(lf) > 200 {
		t.Fatalf("%s: logfile name length too huge: %d \n%s", command, len(lf), lf)
	}

}
