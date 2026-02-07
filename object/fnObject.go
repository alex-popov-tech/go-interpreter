package object

import (
	"fmt"
	"strings"

	"monkey/ast"
)

type FnObject struct {
	LexicalScope *Scope
	Arguments    []*ast.Identifier
	Body         *ast.BlockExpression
}

func (this FnObject) Inspect() string {
	args := []string{}
	for _, arg := range this.Arguments {
		args = append(args, arg.String())
	}
	return fmt.Sprintf("fn (%s) %s", strings.Join(args, ", "), this.Body.String())
}

func (this FnObject) Type() ObjectType {
	return FN
}
