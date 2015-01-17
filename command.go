package runcmd

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

type Command struct {
	// the executable
	Exe string
	// args passed to executable
	Args []string
	// the command is run in a shell, that is prepend `/bin/sh -c` or `cmd /C` to the command line
	CommandLine string
	WorkingDir  string
	// only for *nix: if not set, runcmd uses env[$SHELL] or defaults to /bin/sh
	ForceShell string
	// custom environment variables. these are overwritten from .env file if UseEnv is true
	Env Env
	// only for *nix: if true .profile file in the working dir is sourced
	UseProfile bool // dovrebbe essere Profile string : path
	// only for *nix: if true .env file in the working dir is used to initialize env vars
	UseEnv bool // dovrebbe essere EnvFile string: path
}

func (c *Command) String() string {
	return fmt.Sprintf("%s# %s", c.WorkingDir, c.FullCommand())
}

func (c *Command) FullCommand() string {
	if c.Exe == "" && c.CommandLine == "" {
		return ""
	}
	if c.Exe == "" {
		c.useShell()
	}
	return strings.TrimSpace(c.Exe + " " + strings.Join(c.Args, " "))
}

// Run starts the specified command and waits for it to complete.
func (c *Command) Run() *ExecResult {

	// shell, args := shellAndArgs()
	var bout, berr bytes.Buffer
	streams := &Streams{
		out: &bout,
		err: &berr,
	}
	result := &ExecResult{
		fullCommand: c.FullCommand(),
		streams:     streams,
	}
	cmd, err := c.buildCmd()
	if err != nil {
		result.err = err
		return result
	}
	cmd.Stdout = &bout
	cmd.Stderr = &berr
	if c.WorkingDir != "" {
		cmd.Dir = c.WorkingDir
	}

	if c.UseEnv {
		flagEnv := filepath.Join(cmd.Dir, ".env")
		env, _ := ReadEnv(flagEnv)
		cmd.Env = env.asArray()
	} else if len(c.Env) > 0 {
		cmd.Env = c.Env.asArray()
	}
	// if runtime.GOOS == "windows" {
	// 	r.Environment = mergeEnvironment(r.Environment)
	// }

	if err := cmd.Run(); err != nil {
		result.err = err
		return result
	}
	return result
}

func (c *Command) buildCmd() (*exec.Cmd, error) {
	if c.Exe == "" && c.CommandLine == "" {
		return nil, fmt.Errorf("error creating command: no Exe nor CommandLine")
	}
	if c.Exe != "" {
		return exec.Command(c.Exe, c.Args...), nil
	}
	// if commandline, use shell
	c.useShell()
	cmd := exec.Command(c.Exe, c.Args...)
	return cmd, nil
}
