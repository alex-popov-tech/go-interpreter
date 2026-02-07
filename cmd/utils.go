package cmd

import (
	"fmt"
	"os"
)

func printParserErrors(errors []error) {
	for _, msg := range errors {
		fmt.Fprintf(os.Stderr, "%s\n", msg.Error())
	}
}

func toFilePath(path string) (content string, err error) {
	bytes, err := os.ReadFile(path)
	return string(bytes), err
}
