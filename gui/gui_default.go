//go:build !darwin

package gui

func guiPrompt(prompt string) (string, error) { return "", ErrNoGUI }
