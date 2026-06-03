package runcmd

import (
	"os"
)

// Start start a process without waiting
func (c *Command) Start() error {
	cmd, err := c.buildCmd()
	if err != nil {
		return err
	}
	// Open the logfile once so stdout and stderr always land in the same file.
	// The handle is stored on Command so the GC can finalise it after the
	// process exits and the Command is no longer referenced.
	lf := logfile(c.GetLogfile())
	cmd.Stdout = lf
	cmd.Stderr = lf
	if c.WorkingDir != "" {
		cmd.Dir = c.WorkingDir
	}

	cmd.Env = c.prepareEnv(cmd.Dir)

	process, err := start(cmd)
	if err != nil {
		if lf != nil {
			lf.Close()
		}
		return err
	}
	c.Process = process
	c.logFile = lf
	return nil
}

func logfile(path string) *os.File {
	if path == "" {
		return nil
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0660)
	if err != nil {
		return nil
	}
	return file
}
