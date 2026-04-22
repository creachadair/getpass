//go:build !darwin

package main

func guiPrompt(prompt string) (string, error) { return "", errNoGUI }
