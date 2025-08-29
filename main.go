package main

import (
	"flag"
	"fmt"
	"os"
)

// Set at build time via -ldflags "-X main.AppName=... -X main.AppVersion=... -X main.BuildDate=... -X main.GitCommit=..."
var AppName string
var AppVersion string
var BuildDate string
var GitCommit string

var flagHelp = flag.Bool("help", false, "displays this help message")
var flagVersion = flag.Bool("version", false, "print version and exit")

func init() {
	flag.BoolVar(flagHelp, "h", false, "")
	flag.BoolVar(flagVersion, "v", false, "")
}

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: "+AppName+` [ version | --version | --help ]

Listens for WM_SETTINGCHANGE messages and prints a message when environment variables change.
Press CTRL+C to exit.

OPTIONS:

  -h, -help
        display this help message
  -v, -version
        print version and exit

EXAMPLES:`)

		fmt.Fprintln(os.Stderr, "\n  $ "+AppName+`
    2025/08/29 18:54:57 Environment changed
    2025/08/29 18:55:00 Environment changed
    2025/08/29 18:55:09 Environment changed`)
	}
	flag.Parse()

	if flag.Arg(0) == "version" || *flagVersion {
		fmt.Printf("%s %s, built on %s (commit: %s)\n", AppName, AppVersion, BuildDate, GitCommit)
		return
	}

	if *flagHelp {
		flag.Usage()
		return
	}

	if flag.NArg() > 0 {
		flag.Usage()
		os.Exit(1)
	}

	// Start watching for environment changes.
	watch()
}
