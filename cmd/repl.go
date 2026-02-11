package cmd

import (
	"bufio"
	"fmt"
	"os"

	"monkey/evaluator"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
)

const PROMPT = ">> "

func Repl() {
	in := os.Stdin
	out := os.Stdout
	scanner := bufio.NewScanner(in)

	scope := object.NewGlobalScope()
	fmt.Println("Hello bro! This is the Monkey programming language!")
	fmt.Println("Feel free to type in commands:")
	for {
		fmt.Fprint(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		if line == "q" || line == "quit" {
			fmt.Printf("Bye bye!")
			os.Exit(0)
		}
		l := lexer.New(line)
		p := parser.New(l)
		program := p.ParseProgram()

		if len(p.Errors()) > 0 {
			printParserErrors(p.Errors())
			continue
		}

		output := evaluator.Eval(scope, program)
		fmt.Println(output.Inspect())
	}
}
