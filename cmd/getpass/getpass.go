package main

import (
	"cmp"
	"errors"
	"flag"
	"fmt"
	"log"

	"github.com/creachadair/getpass"
	"github.com/creachadair/getpass/gui"
)

var (
	prompt    = flag.String("prompt", "", "Prompt string")
	doConfirm = flag.Bool("confirm", false, "Require confirmation (repeat response)")
	doGUI     = flag.Bool("gui", false, "Prompt via a GUI (if available)")
)

var errNoGUI = errors.New("no GUI support is available")

func main() {
	flag.Parse()
	label := cmp.Or(*prompt, "Passphrase: ")
	pw, err := call(label)
	if err != nil {
		log.Fatalf("getpass: %v", err)
	}
	if *doConfirm {
		cf, err := call("(confirm) " + label)
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
		pw, err := gui.Prompt(prompt)
		if err == nil {
			return pw, nil
		} else if !errors.Is(err, errNoGUI) {
			return "", err
		}
	}
	return getpass.Prompt(prompt)
}
