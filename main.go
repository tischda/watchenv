package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

// https://goreleaser.com/cookbooks/using-main.version/
var (
	name    string
	version string
	date    string
	commit  string
)

// flags
type Config struct {
	help    bool
	version bool
}

func initFlags() *Config {
	cfg := &Config{}
	flag.BoolVar(&cfg.help, "?", false, "")
	flag.BoolVar(&cfg.help, "help", false, "displays this help message")
	flag.BoolVar(&cfg.version, "v", false, "")
	flag.BoolVar(&cfg.version, "version", false, "print version and exit")
	return cfg
}

func main() {
	cfg := initFlags()
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: "+name+` [OPTIONS]

Listens for WM_SETTINGCHANGE messages and prints a message when environment variables change.
Press CTRL+C to exit.

OPTIONS:

  -h, --help
        display this help message
  -v, --version
        print version and exit

EXAMPLES:`)

		fmt.Fprintln(os.Stderr, "\n  $ "+name+`
	2025/08/29 18:54:50	Listening for environment changes. Press CTRL+C to exit.
    2025/08/29 18:54:57 Environment changed
    2025/08/29 18:55:00 Environment changed
    2025/08/29 18:55:09 Environment changed`)
	}
	flag.Parse()

	if flag.Arg(0) == "version" || cfg.version {
		fmt.Printf("%s %s, built on %s (commit: %s)\n", name, version, date, commit)
		return
	}

	if cfg.help {
		flag.Usage()
		return
	}

	if len(os.Args) > 1 {
		flag.Usage()
		os.Exit(1)
	}

	log.Println("Listening for environment changes. Press CTRL+C to exit.")
	watch()
}
