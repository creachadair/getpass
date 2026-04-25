//go:build !darwin && !linux

package gui

func guiPrompt(prompt string) (string, error) { return "", ErrNoGUI }
