// Package getpass supports reading text from a terminal with echo disabled,
// suitable for prompting a user for a passphrase.
package getpass

import (
	"fmt"
	"os"

	"bufio"

	"bitbucket.org/creachadair/getpass/echo"
)

// TTY opens the controlling terminal of the current process if possible.
func TTY() (*os.File, error) { return os.OpenFile("/dev/tty", os.O_RDWR, 0644) }

// FReadline reads a single line of text from f with echo disabled, returning
// the line without its trailing newline.
func FReadline(f *os.File) (string, error) {
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

// Readline is a shorthand for FReadline using the TTY.
func Readline() (string, error) {
	f, err := TTY()
	if err != nil {
		return "", err
	}
	defer f.Close()
	return FReadline(f)
}

// Prompt prints the prompt string to TTY then calls FReadline.
func Prompt(prompt string) (string, error) {
	f, err := TTY()
	if err != nil {
		return "", err
	}
	defer f.Close()
	fmt.Fprint(f, prompt)
	return FReadline(f)
}
