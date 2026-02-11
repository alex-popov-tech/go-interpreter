package ast

import (
	"fmt"
	"strings"

	"monkey/token"
)

type ArrayExpression struct {
	Token        token.Token // '(' token
	Elements    []Expression
}

func (this *ArrayExpression) expressionNode() {}

func (this ArrayExpression) TokenLiteral() string { return this.Token.Literal }

func (this ArrayExpression) String() string {
	args := []string{}
	for _, arg := range this.Elements {
		args = append(args, arg.String())
	}
	return fmt.Sprintf("[%s]", strings.Join(args, ", "))
}
