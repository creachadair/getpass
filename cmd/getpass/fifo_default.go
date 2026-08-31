//go:build !windows

package main

import (
	"errors"
	"io/fs"
	"os"

	"golang.org/x/sys/unix"
)

// maybeCreatePipe creates a named pipe at the specified path, if one does not
// already exist there. The flag reports whether a pipe was created.
func maybeCreatePipe(path string) (created bool, err error) {
	if err := unix.Mkfifo(path, 0600); err == nil {
		created = true
		defer func() {
			if err != nil {
				os.Remove(path)
			}
		}()
	} else if !errors.Is(err, fs.ErrExist) {
		return false, &fs.PathError{Op: "mkfifo", Path: path, Err: err}
	}
	fi, err := os.Lstat(path)
	if err != nil {
		return created, err
	} else if fi.Mode().Type() != fs.ModeNamedPipe {
		return created, &fs.PathError{Op: "stat", Path: path, Err: errors.New("not a named pipe")}
	}
	return created, nil
}
