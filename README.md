# getpass

This repository provides a Go package to read passwords from the terminal with
echo disabled. This implementation uses cgo, and relies on the POSIX terminal
API functions defined in `<termios.h>` on systems that support them.

View documentation on [GoDoc](http://godoc.org/bitbucket.org/creachadair/getpass).
