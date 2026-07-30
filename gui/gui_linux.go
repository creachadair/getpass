//go:build linux

package gui

func guiPrompt(prompt string) (string, error) { return promptViaPinentry(prompt) }
