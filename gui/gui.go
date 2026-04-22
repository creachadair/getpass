// Package gui supports prompting a human user for input via a GUI.
package gui

import (
	"errors"
)

// ErrNoGUI is a sentinel error reported by [Prompt] when it is not possible to
// prompt via a GUI in the current environment.
var ErrNoGUI = errors.New("GUI support is not available")

// Prompt creates a GUI element displaying the given prompt string, and
// requesting the user to provide a text input.
//
// If it is not possible to create a GUI element in the current environment,
// Prompt reports [ErrNoGUI] with an empty response.
func Prompt(prompt string) (string, error) {
	if prompt == "" {
		return "", errors.New("empty prompt string")
	}
	return guiPrompt(prompt)
}
