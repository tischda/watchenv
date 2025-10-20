[![Build Status](https://github.com/tischda/watchenv/actions/workflows/build.yml/badge.svg)](https://github.com/tischda/watchenv/actions/workflows/build.yml)
[![Linter Status](https://github.com/tischda/watchenv/actions/workflows/linter.yml/badge.svg)](https://github.com/tischda/watchenv/actions/workflows/linter.yml)
[![License](https://img.shields.io/github/license/tischda/watchenv)](/LICENSE)
[![Release](https://img.shields.io/github/release/tischda/watchenv.svg)](https://github.com/tischda/watchenv/releases/latest)


# watchenv

Prints a message when Windows environment variables have changed.

It creates a hidden window that listens for `WM_SETTINGCHANGED` messages.

## Install

~~~
go install github.com/tischda/watchenv@latest
~~~

## Usage

~~~
Usage: watchenv [OPTIONS]

Listens for WM_SETTINGCHANGE messages and prints a message when environment variables change.
Press CTRL+C to exit.

OPTIONS:

  -h, --help
        display this help message
  -v, --version
        print version and exit

EXAMPLES:

  $ watchenv
    2025/08/29 18:54:50	Listening for environment changes. Press CTRL+C to exit.
    2025/08/29 18:54:57 Environment changed
    2025/08/29 18:55:00 Environment changed
    2025/08/29 18:55:09 Environment changed
~~~
