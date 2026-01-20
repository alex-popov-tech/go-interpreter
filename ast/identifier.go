package ast

import "monkey/token"

type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}

func (this *Identifier) expressionNode()      {}
func (this *Identifier) TokenLiteral() string { return this.Token.Literal }
func (this *Identifier) String() string {
	return this.Value
}
