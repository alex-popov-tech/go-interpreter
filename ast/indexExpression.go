package ast

import (
	"fmt"

	"monkey/token"
)

type IndexExpression struct {
	Token           token.Token // '[' token
	Identifier      Expression  // identifier like 'arr'/'myMap' or arr/hash literal
	IndexExpression Expression
}

func (this *IndexExpression) expressionNode() {}

func (this IndexExpression) TokenLiteral() string { return this.Token.Literal }

func (this IndexExpression) String() string {
	return fmt.Sprintf("%s[%s]", this.Identifier.String(), this.IndexExpression.String())
}
