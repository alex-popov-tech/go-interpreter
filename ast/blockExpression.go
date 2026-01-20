package ast

import (
	"strings"

	"monkey/token"
)

type BlockExpression struct {
	Token      token.Token // '{' token
	Statements []Statement
}

func (be BlockExpression) expressionNode() {}

func (be *BlockExpression) TokenLiteral() string { return be.Token.Literal }

func (be *BlockExpression) String() string {
	sb := strings.Builder{}
	sb.WriteString("{")
	for _, s := range be.Statements {
		sb.WriteString(s.String())
	}
	sb.WriteString("}")
	return sb.String()
}
