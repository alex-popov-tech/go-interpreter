package cmd

import (
	"fmt"

	"monkey/evaluator"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
)

func Run(content string) {
	l := lexer.New(content)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		printParserErrors(p.Errors())
		return
	}

	scope := object.NewGlobalScope()
	output := evaluator.Eval(scope, program)
	fmt.Println(output.Inspect())
}
