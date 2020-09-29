package runcmd

import (
	"fmt"
	"os"
	"testing"
	"time"
)

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
