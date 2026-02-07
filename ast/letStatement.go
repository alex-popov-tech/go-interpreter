package ast

import (
	"fmt"

	"monkey/token"
)

type LetStatement struct {
	Token token.Token // 'let' token
	Identifier  *Identifier
	Value Expression
}

func (this *LetStatement) statementNode() {}

func (this LetStatement) TokenLiteral() string { return this.Token.Literal }

func (this LetStatement) String() string {
	return fmt.Sprintf("let %s = %s;", this.Identifier.String(), this.Value.String())
}
