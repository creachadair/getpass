package getpass

import "os"

// TTY opens the controlling terminal of the current process if possible.
func TTY() (*os.File, error) {
	return os.NewFile(os.Stderr.Fd(), "console"), nil
}
