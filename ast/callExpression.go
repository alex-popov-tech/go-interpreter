package ast

import (
	"fmt"
	"strings"

	"monkey/token"
)

type CallExpression struct {
	Token     token.Token // '(' token
	FnIdentifier Expression // identifier like 'add' or fn expression
	Arguments []*Identifier
}

func (this *CallExpression) expressionNode() {}

func (this *CallExpression) TokenLiteral() string { return this.Token.Literal }

func (this *CallExpression) String() string {
	args := []string{}
	for _, arg := range this.Arguments {
		args = append(args, arg.String())
	}
	return fmt.Sprintf("(%s)", strings.Join(args, ", "))
}
