package main

import (
	"cmp"
	"errors"
	"flag"
	"fmt"
	"log"

	"github.com/creachadair/getpass"
)

var (
	prompt    = flag.String("prompt", "Password: ", "Prompt string")
	doConfirm = flag.Bool("confirm", false, "Prompt for confirmation")
	doGUI     = flag.Bool("gui", false, "Prompt via the GUI (if available)")
)

var errNoGUI = errors.New("no GUI support is available")

func main() {
	flag.Parse()
	pw, err := call(*prompt)
	if err != nil {
		log.Fatalf("getpass: %v", err)
	}
	if *doConfirm {
		cf, err := call("(confirm) " + *prompt)
		if err != nil {
			log.Fatalf("get confirmation: %v", err)
		} else if cf != pw {
			log.Fatal("Values do not match")
		}
	}
	fmt.Println(pw)
}

func call(prompt string) (string, error) {
	if *doGUI {
		pw, err := guiPrompt(cmp.Or(prompt, "Passphrase:"))
		if err == nil {
			return pw, nil
		} else if !errors.Is(err, errNoGUI) {
			return "", err
		}
	}
	return getpass.Prompt(prompt)
}
