// Package getpass supports reading text from a terminal with echo disabled,
// suitable for prompting a user for a passphrase.
package getpass

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

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

// Readline reads a single line of text from the TTY of the current process,
// with echo disabled. The line is returned without its trailling newline and
// echo is re-enabled before returning.
func Readline() (string, error) {
	pw, err := ttyIn()
	if err != nil {
		return "", err
	}
	defer pw.Close()
	return freadline(pw)
}

// Prompt prints the prompt string to the TTY of the current process and then
// calls Readline to read a line of text with echo disabled.
func Prompt(prompt string) (string, error) {
	pr, err := ttyOut()
	if err != nil {
		return "", err
	}
	pw, err := ttyIn()
	if err != nil {
		pr.Close()
		return "", err
	}
	defer func() {
		pr.WriteString("\n")
		pr.Close()
		pw.Close()
	}()
	fmt.Fprint(pr, prompt)
	return freadline(pw)
}
