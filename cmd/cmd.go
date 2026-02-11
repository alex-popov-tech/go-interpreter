package cmd

import (
	"fmt"
	"os"
)

func Cmd() {
	if len(os.Args) < 2 {
		fmt.Println("Missing command")
	}

	switch os.Args[1] {
	case "version":
		PrintVersion()
	case "repl":
		Repl()
	case "run":
		if len(os.Args) < 3 {
			fmt.Println("Missing file path")
			os.Exit(1)
		}
		content, err := toFilePath(os.Args[2])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		Run(content)
	default:
		fmt.Println("Unknkown command", fmt.Sprintf("%v", os.Args[1:]))
		os.Exit(1)
	}
}
