// +build windows

package runcmd

var testCommands = []testCommand{
	{
		command: &Command{
			CommandLine: `echo "BAR=%BAR%!"`,
			Env:         Env{"BAR": "foo"},
			Logfile:     "out.log",
		},
		successExpected: true,
	},
}
