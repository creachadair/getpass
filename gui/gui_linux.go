//go:build linux

package gui

import (
	"bytes"
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"strings"
)

func guiPrompt(prompt string) (string, error) {
	desk, ok := os.LookupEnv("XDG_CURRENT_DESKTOP")
	if !ok || !strings.HasSuffix(strings.ToLower(desk), "gnome") {
		return "", ErrNoGUI

		// TODO: Handle KDE too, probably. It is possible the below just works,
		// but I don't have a KDE setup to test it.
	}

	// Execute the "pinentry" binary from the GPG package.
	// This speaks the Assuan text protocol on stdio.
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

	cmd := exec.Command("pinentry")
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
