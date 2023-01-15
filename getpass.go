// Package getpass supports reading text from a terminal with echo disabled,
// suitable for prompting a user for a passphrase.
package getpass

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

// TTY opens the controlling terminal of the current process if possible.
func TTY() (*os.File, error) { return os.OpenFile("/dev/tty", os.O_RDWR, 0644) }

// freadline reads a single line of text from f, a file associated with an open
// terminal, with echo disabled. The line is returned line without its trailing
// newline and echo is (re)enabled before returning.
func freadline(f *os.File) (string, error) {
	fd := f.Fd()
	pw, err := term.ReadPassword(int(fd))
	if err != nil {
		return "", err
	}
	return string(pw), nil
}

// Readline is a shorthand for FReadline using the TTY.
func Readline() (string, error) {
	f, err := TTY()
	if err != nil {
		return "", err
	}
	defer f.Close()
	return freadline(f)
}

// Prompt prints the prompt string to TTY then calls FReadline on it.
func Prompt(prompt string) (string, error) {
	f, err := TTY()
	if err != nil {
		return "", err
	}
	defer f.Close()
	fmt.Fprint(f, prompt)
	return freadline(f)
}
