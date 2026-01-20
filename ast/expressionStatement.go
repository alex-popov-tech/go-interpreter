package ast

import (
	"fmt"

	"monkey/token"
)

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (this *ExpressionStatement) statementNode()       {}
func (this *ExpressionStatement) TokenLiteral() string { return this.Token.Literal }

func (this *ExpressionStatement) String() string {
	return fmt.Sprintf("%s;", this.Expression.String())
}
