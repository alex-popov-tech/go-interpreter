package ast

import (
	"fmt"

	"monkey/token"
)

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Value    Expression
}

func (this *PrefixExpression) expressionNode() {}

func (this *PrefixExpression) TokenLiteral() string { return this.Token.Literal }

func (this *PrefixExpression) String() string {
	return fmt.Sprintf("(%s%s)", this.Operator, this.Value.String())
}
