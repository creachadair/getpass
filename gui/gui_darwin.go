//go:build darwin

package gui

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func guiPrompt(prompt string) (string, error) {
	if canProbablySeeUI() {
		return promptViaApplescript(prompt)
	}
	return promptViaPinentry(prompt)
}

func canProbablySeeUI() bool {
	// If there is no TERM_PROGRAM set at all, it often means the program was
	// started by launchd or some other program running in the UI.
	// Otherwise, guess based on the name. TODO(creachadair): Maybe include other
	// popular GUI terminal substitutions here.
	tp, ok := os.LookupEnv("TERM_PROGRAM")
	return !ok || strings.HasPrefix(strings.ToLower(tp), "apple_terminal")
}

func promptViaApplescript(prompt string) (string, error) {
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
