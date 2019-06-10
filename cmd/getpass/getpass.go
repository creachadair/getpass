package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/creachadair/getpass"
)

var prompt = flag.String("prompt", "Password: ", "Prompt string")

func main() {
	flag.Parse()
	pw, err := getpass.Prompt(*prompt)
	if err != nil {
		log.Fatalf("getpass: %v", err)
	}
	fmt.Println(pw)
}
