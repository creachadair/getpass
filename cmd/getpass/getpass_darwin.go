//go:build darwin

package main

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

func guiPrompt(prompt string) (string, error) {
	cmd := exec.Command("osascript", "-s", "ho")
	cmd.Stdin = strings.NewReader(fmt.Sprintf(`display dialog %q default answer "" hidden answer true`, prompt))
	raw, err := cmd.Output()
	out := strings.TrimRight(string(raw), "\n")
	if err != nil {
		if strings.Contains(out, "User canceled") {
			return "", errors.New("user cancelled request")
		}
		return "", err
	}
	const needle = "text returned:"
	if _, after, ok := strings.Cut(out, needle); ok {
		return after, nil
	}
	return "", errors.New("missing user response")
}
