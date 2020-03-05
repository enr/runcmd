package runcmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

var (
	unixTemplate = `#!/bin/bash
cat %s

cat %s >&2

exit %d
`

	winTemplate = `@echo off
type %s

type %s >&2

exit /B %d
`
)

func mockCmdOutput(cmdName, stdout string, stderr string, exitCode int) {
	_ = os.Mkdir("build", os.ModePerm)
	tmpDir, err := ioutil.TempDir("build", cmdName+"-")
	if err != nil {
		log.Fatal(err)
	}
	pathToFakeBin := filepath.Join(tmpDir, cmdName)

	var template string
	switch runtime.GOOS {
	case "windows":
		pathToFakeBin += ".bat"
		template = winTemplate
	default:
		template = unixTemplate
	}

	stdoutFile := filepath.Join(tmpDir, "stdout.txt")
	_ = ioutil.WriteFile(stdoutFile, []byte(stdout), os.ModePerm)
	// c.Assert(err, check.IsNil)

	stderrFile := filepath.Join(tmpDir, "stderr.txt")
	_ = ioutil.WriteFile(stderrFile, []byte(stderr), os.ModePerm)

	fakeBin, err := os.OpenFile(
		pathToFakeBin,
		os.O_RDWR|os.O_CREATE|os.O_TRUNC,
		0777,
	)
	// c.Assert(err, check.IsNil)

	_, err = fmt.Fprintf(fakeBin, template, stdoutFile, stderrFile, exitCode)
	// c.Assert(err, check.IsNil)

	err = fakeBin.Close()
	// c.Assert(err, check.IsNil)

	newPath := fmt.Sprintf("%s%c%s", tmpDir, os.PathListSeparator, os.Getenv("PATH"))
	err = os.Setenv("PATH", newPath)
	// c.Assert(err, check.IsNil)
}

func Test01(t *testing.T) {
	// t.Skip()
	success := true
	exitCode := 0
	stdout := "stdout"
	stderr := "stderr"
	executable := "ping"

	mockCmdOutput(executable, stdout, stderr, exitCode)

	args := []string{"-al"}
	command := &Command{
		Exe:  executable,
		Args: args,
	}
	res := command.Run()

	if res.Success() != success {
		t.Fatalf("%s: expected success %t but got %t", command, success, res.Success())
	}
	expectedCode := exitCode
	actualCode := res.ExitStatus()
	if actualCode != expectedCode {
		t.Fatalf("%s: expected exit code %d but got %d", command, expectedCode, actualCode)
	}
	assertStringContains(t, res.Stdout().String(), stdout)
	assertStringContains(t, res.Stderr().String(), stderr)

}
