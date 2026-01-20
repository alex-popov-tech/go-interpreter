package ast

import (
	"fmt"
	"monkey/token"
)

type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (this *InfixExpression) expressionNode()      {}
func (this *InfixExpression) TokenLiteral() string { return this.Token.Literal }
func (this *InfixExpression) String() string {
	return fmt.Sprintf("(%s %s %s)", this.Left.String(), this.Operator, this.Right.String())
}
