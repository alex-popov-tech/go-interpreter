package ast

import (
	"fmt"

	"monkey/token"
)

type ReturnStatement struct {
	Token token.Token
	Value Expression
}

func (this *ReturnStatement) statementNode()       {}
func (this *ReturnStatement) TokenLiteral() string { return this.Token.Literal }
func (this *ReturnStatement) String() string {
	return fmt.Sprintf("return %s;", this.Value.String())
}
