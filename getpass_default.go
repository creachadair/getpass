//go:build !windows

package getpass

import (
	"os"
)

// ttyIn opens the controlling terminal of the current process for reading.
func ttyIn() (*os.File, error) {
	return os.OpenFile("/dev/tty", os.O_RDWR, 0644)
}

// ttyOut opens the controlling terminal of the current process for writing.
func ttyOut() (*os.File, error) {
	return os.OpenFile("/dev/tty", os.O_RDWR, 0644)
}
