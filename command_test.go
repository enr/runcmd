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

func TestCommandString(t *testing.T) {
	cmd := &Command{
		WorkingDir: "/tmp",
		Exe:        "/bin/echo",
		Args:       []string{"hello"},
	}
	s := cmd.String()
	if !strings.Contains(s, "/tmp") {
		t.Fatalf("String() should contain WorkingDir, got: %s", s)
	}
	if !strings.Contains(s, "/bin/echo") {
		t.Fatalf("String() should contain exe, got: %s", s)
	}
}

func TestGetNameExplicit(t *testing.T) {
	cmd := &Command{Name: "my-custom-name"}
	if got := cmd.GetName(); got != "my-custom-name" {
		t.Fatalf("GetName: expected %q, got %q", "my-custom-name", got)
	}
}
