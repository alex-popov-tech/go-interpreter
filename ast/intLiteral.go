package ast

import (
	"monkey/token"
)

type IntLiteral struct {
	Token token.Token
	Value int64
}

func (this *IntLiteral) expressionNode()      {}
func (this *IntLiteral) TokenLiteral() string { return this.Token.Literal }
func (this *IntLiteral) String() string       { return this.Token.Literal }
