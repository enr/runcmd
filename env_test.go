package runcmd

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestReadEnvMissing(t *testing.T) {
	env, err := readEnv("/tmp/runcmd-definitely-does-not-exist-99999.env")
	if err != nil {
		t.Fatalf("readEnv on missing file should not error, got: %v", err)
	}
	if len(env) != 0 {
		t.Fatalf("readEnv on missing file should return empty env, got: %v", env)
	}
}

func TestReadEnvFile(t *testing.T) {
	f, err := os.CreateTemp("", "runcmd-test-*.env")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	if _, err := f.WriteString("RUNCMD_FOO=bar\nRUNCMD_BAZ=qux\n"); err != nil {
		t.Fatal(err)
	}
	f.Close()

	env, err := readEnv(f.Name())
	if err != nil {
		t.Fatalf("readEnv: unexpected error: %v", err)
	}
	if got := env["RUNCMD_FOO"]; got != "bar" {
		t.Fatalf("expected RUNCMD_FOO=bar, got %q", got)
	}
	if got := env["RUNCMD_BAZ"]; got != "qux" {
		t.Fatalf("expected RUNCMD_BAZ=qux, got %q", got)
	}
}

func TestMergeEnvironmentNil(t *testing.T) {
	if mergeEnvironment(nil) != nil {
		t.Fatal("mergeEnvironment(nil) should return nil")
	}
}

func TestMergeEnvironment(t *testing.T) {
	k := fmt.Sprintf(`TEST_K_%d`, time.Now().UnixNano())
	v := fmt.Sprintf(`TEST_V_%d`, time.Now().UnixNano())
	osEnvLine := fmt.Sprintf(`%s=%s`, k, v)
	cmdEnvLine := `COMMAND_TEST_VAR=foo`
	os.Setenv(k, v)
	cmdEnv := []string{
		cmdEnvLine,
	}
	env := mergeEnvironment(cmdEnv)
	fmt.Printf("%q \n", env)
	if len(env) < 2 {
		t.Errorf(`Env does not contains all variables`)
	}
	if !contains(env, osEnvLine) {
		t.Errorf(`Missing value from OS in env: %s`, osEnvLine)
	}
	if !contains(env, cmdEnvLine) {
		t.Errorf(`Missing value from cmd in env: %s`, cmdEnvLine)
	}
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
