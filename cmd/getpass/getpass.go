package main

import (
	"cmp"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/creachadair/command"
	"github.com/creachadair/flax"
	"github.com/creachadair/getpass"
	"github.com/creachadair/getpass/gui"
)

var flags struct {
	Prompt  string `flag:"prompt,Prompt string"`
	Confirm bool   `flag:"confirm,Require confirmation (repeat response)"`
	GUI     bool   `flag:"gui,Prompt via a GUI (if available)"`
}

var errNoGUI = errors.New("no GUI support is available")

func main() {
	root := &command.C{
		Name:     command.ProgramName(),
		Help:     `Prompt the user for a passphrase.`,
		SetFlags: command.Flags(flax.MustBind, &flags),
		Run: command.Adapt(func(env *command.Env) error {
			label := cmp.Or(flags.Prompt, "Passphrase: ")
			pw, err := call(label)
			if err != nil {
				log.Fatalf("getpass: %v", err)
			}
			if flags.Confirm {
				cf, err := call("(confirm) " + label)
				if err != nil {
					return fmt.Errorf("get confirmation: %w", err)
				} else if cf != pw {
					return errors.New("values do not match")
				}
			}
			fmt.Println(pw)
			return nil
		}),
		Commands: []*command.C{
			command.HelpCommand(nil),
			command.VersionCommand(),
		},
	}
	command.RunOrFail(root.NewEnv(nil), os.Args[1:])
}

func call(prompt string) (string, error) {
	if flags.GUI {
		pw, err := gui.Prompt(prompt)
		if err == nil {
			return pw, nil
		} else if !errors.Is(err, errNoGUI) {
			return "", err
		}
	}
	return getpass.Prompt(prompt)
}
