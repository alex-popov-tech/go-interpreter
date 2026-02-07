package ast

import (
	"monkey/token"
)

type BoolLiteral struct {
	Token token.Token
	Value bool
}

func (this *BoolLiteral) expressionNode()      {}
func (this BoolLiteral) TokenLiteral() string { return this.Token.Literal }
func (this BoolLiteral) String() string       { return this.Token.Literal }
