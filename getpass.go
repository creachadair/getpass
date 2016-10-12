// Package getpass supports reading text from a terminal with echo disabled,
// suitable for prompting a user for a passphrase.
package getpass

import (
	"fmt"
	"os"

	"bufio"

	"bitbucket.org/creachadair/getpass/echo"
)

// Readline reads a single line of text from f with echo disabled, returning
// the line without its trailing newline.
func Readline(f *os.File) (string, error) {
	fd := f.Fd()
	if err := echo.Disable(fd); err != nil {
		return "", err
	}
	defer echo.Enable(fd)

	rd := bufio.NewScanner(f)
	if !rd.Scan() {
		return "", rd.Err()
	}
	return rd.Text(), nil
}

// Prompt prints the given prompt string to os.Stderr, then calls Readline to
// read a line of text from f with echo disabled.
func Prompt(prompt string, f *os.File) (string, error) {
	fmt.Fprint(os.Stderr, prompt)
	return Readline(f)
}
