package main

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/signal"
	"syscall"

	"github.com/creachadair/command"
	"github.com/creachadair/flax"
	"github.com/creachadair/getpass"
	"github.com/creachadair/getpass/gui"
)

var flags struct {
	Prompt  string `flag:"prompt,Prompt string"`
	Confirm bool   `flag:"confirm,Require confirmation (repeat response)"`
	GUI     bool   `flag:"gui,Prompt via a GUI (if available)"`
	Offer   string `flag:"offer,Offer passphrase to named pipe at this path"`
}

var errNoGUI = errors.New("no GUI support is available")

func main() {
	root := &command.C{
		Name: command.ProgramName(),
		Help: `Prompt the user for a passphrase.

By default, prompt the user at the calling terminal and print the resulting
passphrase to stdout. If --prompt is not set, a generic prompt is generated.
With --gui, a graphical prompt is issued where possible.  With --confirm, the
user is prompted twice and an error is reported if the results do not match.

If the user cancels the request, or an interrupt is received, the prompt is
aborted and getpass exits with code 1.

With --offer, create (if necessary) a named pipe at the given path, and open
it. The prompt is deferred until a reader arrives, and the passphrase is
written to the pipe instead of stdout.
`,
		SetFlags: command.Flags(flax.MustBind, &flags),
		Run: command.Adapt(func(env *command.Env) error {
			if flags.Offer != "" {
				fmt.Fprintln(env, "...")
				return offerToPipe(env.Context(), flags.Offer, promptAndWrite)
			}
			return promptAndWrite(os.Stdout)
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

func promptAndWrite(w io.Writer) error {
	label := cmp.Or(flags.Prompt, "Passphrase: ")
	pw, err := call(label)
	if err != nil {
		return fmt.Errorf("getpass: %v", err)
	}
	if flags.Confirm {
		cf, err := call("(confirm) " + label)
		if err != nil {
			return fmt.Errorf("get confirmation: %w", err)
		} else if cf != pw {
			return errors.New("values do not match")
		}
	}
	fmt.Fprint(w, pw)
	return nil
}

func offerToPipe(ctx context.Context, fifoPath string, getpass func(io.Writer) error) error {
	// Pipe fitting. Opening a named pipe or writing blocks until a reader is
	// available, but we want a way to interrupt and clean up. So we'll create
	// the pipe itself first, then service the write in a goroutine. To mitigate
	// TOCTTOU, we will check again after opening that we actually have a pipe.
	didCreate, err := maybeCreatePipe(fifoPath)
	if err != nil {
		return err
	} else if didCreate {
		defer os.Remove(fifoPath) // best-effort
	}

	// Handle signals so we can clean up. If the context ends, we will open the
	// pipe for reading ourselves
	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	stop := context.AfterFunc(ctx, func() {
		r, err := os.Open(fifoPath)
		if err == nil {
			r.Close()
		}
	})

	// N.B.: Not O_TRUNC, since this might not actually be ours.
	f, err := os.OpenFile(fifoPath, os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	// Re-check to avert TOCTTOU. This doesn't prevent A-B-A substitutions, but
	// at least we know we still have the right type.
	if fi, err := f.Stat(); err != nil {
		return err
	} else if fi.Mode().Type() != fs.ModeNamedPipe {
		return fmt.Errorf("path %q is not a named pipe", fifoPath)
	}
	stop()
	if err != nil {
		return err
	} else if ctx.Err() != nil {
		return ctx.Err()
	}
	cancel() // release the signal handler

	// Reaching here, we have a pipe to write to.
	werr := getpass(f)
	return errors.Join(werr, f.Close())
}
