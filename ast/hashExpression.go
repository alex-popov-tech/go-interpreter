package ast

import (
	"fmt"
	"strings"

	"monkey/token"
)

type HashExpression struct {
	Token token.Token // '(' token
	Map   map[Expression]Expression
}

func (this *HashExpression) expressionNode() {}

func (this HashExpression) TokenLiteral() string { return this.Token.Literal }

func (this HashExpression) String() string {
	strs := []string{}
	for k, v := range this.Map {
		strs = append(strs, fmt.Sprintf("%s:%s", k, v))
	}
	return fmt.Sprintf("#{%s}", strings.Join(strs, ", "))
}
