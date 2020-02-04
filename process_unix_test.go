// +build darwin freebsd linux netbsd openbsd

package runcmd

var testCommands = []testCommand{
	{
		command: &Command{
			CommandLine: `echo "BAR=${BAR}!"`,
			ForceShell:  "/bin/bash",
			Env:         Env{"BAR": "foo"},
			Logfile:     "out.log",
		},
		success: true,
	},
}
