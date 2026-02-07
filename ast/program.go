package ast

import (
	"bytes"
	"fmt"
)

type Program struct {
	Statements []Statement
}

func (this Program) TokenLiteral() string {
	return ""
}

func (this Program) String() string {
	var buf bytes.Buffer

	for _, s := range this.Statements {
		buf.WriteString(fmt.Sprintf("%s\n", s.String()))
	}

	return buf.String()
}
