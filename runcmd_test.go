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

func failIfErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func mockCmdOutput(cmdName, stdout string, stderr string, exitCode int) {
	err := os.MkdirAll("build", os.ModePerm)
	failIfErr(err)

	tmpDir, err := ioutil.TempDir("build", cmdName+"-")
	failIfErr(err)

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
	err = ioutil.WriteFile(stdoutFile, []byte(stdout), os.ModePerm)
	failIfErr(err)

	stderrFile := filepath.Join(tmpDir, "stderr.txt")
	err = ioutil.WriteFile(stderrFile, []byte(stderr), os.ModePerm)
	failIfErr(err)

	fakeBin, err := os.OpenFile(
		pathToFakeBin,
		os.O_RDWR|os.O_CREATE|os.O_TRUNC,
		0777,
	)
	failIfErr(err)

	_, err = fmt.Fprintf(fakeBin, template, stdoutFile, stderrFile, exitCode)
	failIfErr(err)

	err = fakeBin.Close()
	failIfErr(err)

	newPath := fmt.Sprintf("%s%c%s", tmpDir, os.PathListSeparator, os.Getenv("PATH"))
	err = os.Setenv("PATH", newPath)
	failIfErr(err)
}

type commandRunTestCase struct {
	success    bool
	exitCode   int
	stdout     string
	stderr     string
	executable string
}

var testCommandRun = []commandRunTestCase{
	{
		success:    true,
		exitCode:   0,
		stdout:     "stdout",
		stderr:     "stderr",
		executable: "test1",
	},
	{
		success:    false,
		exitCode:   3,
		stdout:     "stdout",
		stderr:     "stderr",
		executable: "test2",
	},
}

func TestCommandRun(t *testing.T) {

	for _, c := range testCommandRun {
		mockCmdOutput(c.executable, c.stdout, c.stderr, c.exitCode)

		command := &Command{
			Exe: c.executable,
		}
		res := command.Run()

		if res.Success() != c.success {
			t.Fatalf("%s: expected success %t but got %t", command, c.success, res.Success())
		}
		expectedCode := c.exitCode
		actualCode := res.ExitStatus()
		if actualCode != expectedCode {
			t.Fatalf("%s: expected exit code %d but got %d", command, expectedCode, actualCode)
		}
		assertStringContains(t, res.Stdout().String(), c.stdout)
		assertStringContains(t, res.Stderr().String(), c.stderr)
	}
}

func TestCommandStart(t *testing.T) {

	for _, c := range testCommandRun {
		mockCmdOutput(c.executable, c.stdout, c.stderr, c.exitCode)

		command := &Command{
			Exe: c.executable,
		}
		err := command.Start()

		if c.success && err != nil {
			t.Fatalf("%s: success expected but got error %v", command, err)
		}
		runningProcess := command.Process
		ps, err := runningProcess.Wait()
		if c.success != ps.Success() {
			t.Fatalf("%s: expected success=%t but got %t", command, c.success, ps.Success())
		}
		if c.success && err != nil {
			t.Fatalf("%s: success expected but got error %v", command, err)
		}
	}
}
