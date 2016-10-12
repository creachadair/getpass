// Package echo allows enabling and disabling terminal echo.
//
// This implementation is based on the POSIX tcgetattr and tcsetattr functions,
// which are accessed via cgo, so it will only work on systems that support
// both cgo and supply those functions.
package echo

/*
#include <unistd.h>
#include <termios.h>
*/
import "C"
import (
	"errors"
	"syscall"

	"golang.org/x/sys/unix"
)

// isatty reports whether desriptor fd is a TTY.  In case of error, the answer
// is assumed to be false.
func isatty(fd uintptr) bool { return C.isatty(C.int(fd)) > 0 }

func tcgetattr(fd uintptr) (*C.struct_termios, error) {
	var info C.struct_termios
	if ec := C.tcgetattr(C.int(fd), &info); ec < 0 {
		return nil, errors.New("tcgetattr failed")
	}
	return &info, nil
}

func tcsetattr(fd uintptr, info *C.struct_termios) error {
	if ec := C.tcsetattr(C.int(fd), C.int(unix.TCSAFLUSH), info); ec < 0 {
		return errors.New("tcsetattr failed")
	}
	return nil
}

const termBits = syscall.ECHO | syscall.ECHOE | syscall.ECHOK

// Enable enables echo for the terminal connected to descriptor fd.
func Enable(fd uintptr) error {
	if !isatty(fd) {
		return syscall.ENOTTY
	}
	info, err := tcgetattr(fd)
	if err != nil {
		return err
	}
	info.c_lflag |= termBits
	return tcsetattr(fd, info)
}

// Disable disables echo for the terminal connected to descriptor fd.
func Disable(fd uintptr) error {
	if !isatty(fd) {
		return syscall.ENOTTY
	}
	info, err := tcgetattr(fd)
	if err != nil {
		return err
	}
	info.c_lflag &^= termBits
	return tcsetattr(fd, info)
}
