package repl

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"monkey/lexer"
	"monkey/parser"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	fmt.Printf("Hello bro! Thin is the Monkey programming language!\n")
	fmt.Printf("Feel free to type in commands:\n")
	for {
		fmt.Fprintf(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		if line == "q" || line == "quit" {
			os.Exit(0)
		}
		l := lexer.New(line)
		p := parser.New(l)
		program := p.ParseProgram()

		if len(p.Errors()) > 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		fmt.Fprintf(out, "%s", program.String())
	}
}

func printParserErrors(out io.Writer, errors []error) {
	for _, msg := range errors {
		fmt.Fprintf(out, "\t%s\n", msg.Error())
	}
}
