package getpass

import (
	"os"
	"syscall"
)

// ttyIn opens the controlling terminal of the current process for reading.
func ttyIn() (*os.File, error) {
	fd, err := syscall.Open("CONIN$", syscall.O_RDWR, 0)
	if err != nil {
		return nil, err
	}
	return os.NewFile(uintptr(fd), "console-in"), nil
}

// ttyOut opens the controlling terminal of the current process for writing.
func ttyOut() (*os.File, error) {
	fd, err := syscall.Open("CONOUT$", syscall.O_RDWR, 0)
	if err != nil {
		return nil, err
	}
	return os.NewFile(uintptr(fd), "console-out"), nil
}
