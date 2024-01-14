package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/creachadair/getpass"
)

var (
	prompt    = flag.String("prompt", "Password: ", "Prompt string")
	doConfirm = flag.Bool("confirm", false, "Prompt for confirmation")
)

func main() {
	flag.Parse()
	pw, err := getpass.Prompt(*prompt)
	if err != nil {
		log.Fatalf("getpass: %v", err)
	}
	if *doConfirm {
		cf, err := getpass.Prompt("(confirm) " + *prompt)
		if err != nil {
			log.Fatalf("get confirmation: %v", err)
		} else if cf != pw {
			log.Fatal("Values do not match")
		}
	}
	fmt.Println(pw)
}
