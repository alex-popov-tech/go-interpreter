package ast

import (
	"monkey/token"
)

type StringLiteral struct {
	Token token.Token
	Value string
}

func (this *StringLiteral) expressionNode()      {}
func (this *StringLiteral) TokenLiteral() string { return this.Token.Literal }
func (this *StringLiteral) String() string       { return this.Token.Literal }
