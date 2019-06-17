// Package echo allows enabling and disabling terminal echo.
//
// This implementation is based on the POSIX tcgetattr and tcsetattr
// operations, so it will only work on systems that support them via the
// ioctl(2) interface.
package echo

import (
	"syscall"

	"golang.org/x/sys/unix"
)

const termBits = syscall.ECHO | syscall.ECHOE | syscall.ECHOK

// Enable enables echo for the terminal connected to descriptor fd.
func Enable(fd uintptr) error {
	t, err := unix.IoctlGetTermios(int(fd), tcGet)
	if err != nil {
		return err
	}
	t.Lflag |= termBits | syscall.ECHONL
	return unix.IoctlSetTermios(int(fd), tcSet, t)
}

// Disable disables echo for the terminal connected to descriptor fd.
func Disable(fd uintptr) error {
	t, err := unix.IoctlGetTermios(int(fd), tcGet)
	if err != nil {
		return err
	}
	t.Lflag &^= termBits
	t.Lflag |= syscall.ECHONL
	return unix.IoctlSetTermios(int(fd), tcSet, t)
}
