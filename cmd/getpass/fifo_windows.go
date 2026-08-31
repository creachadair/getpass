//go:build windows

package main

import "errors"

func maybeCreatePipe(path string) (bool, error) {
	return false, errors.New("named pipes are not supported here")
}
