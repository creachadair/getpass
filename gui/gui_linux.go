//go:build linux

package gui

import (
	"bytes"
	"errors"
	"fmt"
	"net/url"
	"os/exec"
	"strings"
)

func guiPrompt(prompt string) (string, error) {
	// Execute the "pinentry" binary from the GPG package.
	// If we can't find a pinentry binary, assume it is not possible to prompt.
	// Otherwise, trust that it knows what to do (i.e., GUI vs. curses).
	ppath, err := exec.LookPath("pinentry")
	if err != nil || ppath == "" {
		return "", ErrNoGUI
	}

	// The pinentry tool speaks the Assuan text protocol on stdio.
	// For our purposes, we just want to set a prompt and then fetch a PIN.
	// That is, we expect this interaction between us (C) and pinentry (S):
	//
	//    C: SETPROMPT <prompt>
	//    S: OK
	//    C: GETPIN
	//    S: D <pin>
	//    S: OK
	//    C: BYE
	//    S: OK closing connection
	//
	// If the user cancels, the server will reply "ERR" instead.
	var req bytes.Buffer
	fmt.Fprintf(&req, "SETPROMPT %s\n", esc.Replace(prompt))
	req.WriteString("GETPIN\nBYE\n")

	cmd := exec.Command(ppath, "--ttyname=/dev/tty")
	cmd.Stdin = &req
	rsp, err := cmd.Output()
	if err != nil {
		return "", err
	}
	var result string
	for line := range strings.SplitSeq(string(rsp), "\n") {
		if tail, ok := strings.CutPrefix(line, "D "); ok {
			result += unescape(tail) // it is possible, though unlikely, to get multiple D lines
		} else if _, ok := strings.CutPrefix(line, "ERR "); ok {
			return "", errors.New("user cancelled request") // probably
		}
	}
	return result, nil
}

var esc = strings.NewReplacer(
	"\n", "%0A",
	"\r", "%0D",
	"%", "%25",
)

func unescape(s string) string {
	q, err := url.QueryUnescape(s)
	if err != nil {
		return s
	}
	return q
}
