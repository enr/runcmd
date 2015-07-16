// +build windows

package runcmd

var testdata2 = []cnl{
	{command: &Command{},
		name:    "",
		logfile: "",
	},
	{command: &Command{
		CommandLine: `echo "home=%HOME%"`,
		UseEnv:      true,
	},
		name:    "cmd-c-echo-home-home",
		logfile: "cmd-c-echo-home-home",
	},
}
