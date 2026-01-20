package main

import (
	"os"

	"monkey/repl"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "repl" {
		os.Setenv("TRACE", "true")
		repl.Start(os.Stdin, os.Stdout)
	}
}
