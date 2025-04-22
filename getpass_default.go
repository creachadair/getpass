//go:build !windows

package getpass

import "os"

// TTY opens the controlling terminal of the current process if possible.
func TTY() (*os.File, error) { return os.OpenFile("/dev/tty", os.O_RDWR, 0644) }
