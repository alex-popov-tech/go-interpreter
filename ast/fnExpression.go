package ast

import (
	"fmt"
	"strings"

	"monkey/token"
)

type FnExpression struct {
	Token     token.Token // 'fn' token
	Arguments []*Identifier
	Body      *BlockExpression
}

func (it FnExpression) expressionNode() {}

func (it *FnExpression) TokenLiteral() string { return it.Token.Literal }

func (it *FnExpression) String() string {
	args := []string{}
	for _, arg := range it.Arguments {
		args = append(args, arg.String())
	}
	return fmt.Sprintf("%s(%s)%s", it.TokenLiteral(), strings.Join(args, ", "), it.Body.String())
}
