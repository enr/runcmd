// +build darwin freebsd linux netbsd openbsd

package runcmd

var testdata2 = []cnl{
	{command: &Command{
		Logfile: "out.log",
	},
		name:    "",
		logfile: "out.log",
	},
	{command: &Command{
		CommandLine: `echo "home=$HOME"`,
		UseEnv:      true,
	},
		name:    "bin-bash-c-echo-home-home",
		logfile: "runcmd-bin-bash-c-echo-home-home.log",
	},
}
