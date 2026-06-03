//go:build darwin || freebsd || linux || netbsd || openbsd
// +build darwin freebsd linux netbsd openbsd

package runcmd

import (
	"os"
	"path/filepath"
	"testing"
)

var testdata2 = []cnl{
	{
		command: &Command{
			Logfile: "out.log",
		},
		name:    "",
		logfile: "out.log",
	},
	{
		command: &Command{
			CommandLine: `echo "home=$HOME"`,
			UseEnv:      true,
		},
		name:    "echo-home-home",
		logfile: "runcmd-echo-home-home.log",
	},
	{
		command: &Command{
			Exe:    `/usr/local/bin/myapp`,
			Args:   []string{"-a", `"say hello!"`},
			UseEnv: true,
		},
		name:    "usr-local-bin-myapp-a-say-hello",
		logfile: "runcmd-usr-local-bin-myapp-a-say-hello.log",
	},
}

func TestPrepareEnvWithUseEnv(t *testing.T) {
	dir, err := os.MkdirTemp("", "runcmd-env-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	if err := os.WriteFile(filepath.Join(dir, ".env"), []byte("RUNCMD_TEST_UNIQUE=hello\n"), 0644); err != nil {
		t.Fatal(err)
	}

	env := (&Command{UseEnv: true}).prepareEnv(dir)
	for _, e := range env {
		if e == "RUNCMD_TEST_UNIQUE=hello" {
			return
		}
	}
	t.Fatal("prepareEnv with UseEnv did not include variables from .env file")
}

func TestUseShellNoShellEnv(t *testing.T) {
	orig, had := os.LookupEnv("SHELL")
	os.Unsetenv("SHELL")
	defer func() {
		if had {
			os.Setenv("SHELL", orig)
		}
	}()

	cmd := &Command{CommandLine: "echo hello"}
	cmd.useShell()
	if cmd.Exe != "/bin/sh" {
		t.Fatalf("expected /bin/sh when SHELL is unset, got %q", cmd.Exe)
	}
}
