![Build Status](https://github.com/tischda/watchenv/actions/workflows/build.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/tischda/watchenv)](https://goreportcard.com/report/github.com/tischda/watchenv)

# watchenv

Windows utility that prints a message when environment variables are changed.

It creates a hidden window that listens for `WM_SETTINGCHANGED` messages.

### Install

~~~
go install github.com/tischda/watchenv@latest
~~~

### Usage

~~~
Usage: watchenv [ version | --version | --help ]

Listens for WM_SETTINGCHANGE messages and prints a message when environment variables change.
Press CTRL+C to exit.

OPTIONS:

  -h, -help
        display this help message
  -v, -version
        print version and exit

EXAMPLES:

  $ watchenv
    2025/08/29 18:54:57 Environment changed
    2025/08/29 18:55:00 Environment changed
    2025/08/29 18:55:09 Environment changed
~~~
