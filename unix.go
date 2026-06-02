//go:build darwin || freebsd || linux || netbsd || openbsd
// +build darwin freebsd linux netbsd openbsd

package runcmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func (c *Command) useShell() {
	command := c.CommandLine
	var shell string
	if c.ForceShell != "" {
		shell = c.ForceShell
	} else {
		defaultShell := os.Getenv("SHELL")
		if defaultShell == "" {
			shell = "/bin/sh"
		} else {
			shell = defaultShell
		}
	}
	shellArgument := "-c"
	var shellCommand string
	if c.UseProfile {
		profile := filepath.Join(c.WorkingDir, ".profile")
		// Single-quote the path so shell metacharacters in WorkingDir can't inject commands.
		// Any literal single quotes in the path are escaped as '\''.
		escapedProfile := "'" + strings.ReplaceAll(profile, "'", `'\''`) + "'"
		// "." is portable, "source" is bash only
		shellCommand = fmt.Sprintf(". %s 2>/dev/null; %s", escapedProfile, command)
	} else {
		shellCommand = command
	}
	c.Exe = shell
	c.Args = []string{shellArgument, shellCommand}
}
