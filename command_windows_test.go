// +build windows

package runcmd

var testdata2 = []cnl{
	{command: &Command{
		Logfile: "out.log",
	},
		name:    "",
		logfile: "out.log",
	},
	{command: &Command{
		CommandLine: `echo "home=%HOME%"`,
		UseEnv:      true,
	},
		name:    "cmd-c-echo-home-home",
		logfile: "runcmd-cmd-c-echo-home-home.log",
	},
}
