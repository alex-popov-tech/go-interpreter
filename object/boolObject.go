package object

import (
	"fmt"
)

type BoolObject struct {
	Value bool
}

func (this BoolObject) Inspect() string {
	return fmt.Sprintf("%t", this.Value)
}

func (this BoolObject) Type() ObjectType {
	return BOOL
}
