package ast

import (
	"strings"

	"monkey/token"
)

type IfExpression struct {
	Token        token.Token // 'if' token
	Condition    Expression
	IfBlock      *BlockExpression
	ElseIfBlocks []*ElseIfBlock    // optional
	ElseBlock    *BlockExpression // optional
}

func (this *IfExpression) expressionNode() {}

func (this *IfExpression) TokenLiteral() string { return this.Token.Literal }

func (this *IfExpression) String() string {
	sb := strings.Builder{}
	sb.WriteString("if (")
	sb.WriteString(this.Condition.String())
	sb.WriteString(")")
	sb.WriteString(this.IfBlock.String())

	if this.ElseIfBlocks != nil {
		for _, ieb := range this.ElseIfBlocks {
			sb.WriteString(ieb.String())
		}
	}

	if this.ElseBlock != nil {
		sb.WriteString("else")
		sb.WriteString(this.ElseBlock.String())
	}

	return sb.String()
}

type ElseIfBlock struct {
	Token     token.Token // 'else' token
	Condition Expression
	Block     *BlockExpression
}

func (this *ElseIfBlock) String() string {
	sb := strings.Builder{}
	sb.WriteString("else if (")
	sb.WriteString(this.Condition.String())
	sb.WriteString(")")
	sb.WriteString(this.Block.String())
	return sb.String()
}
