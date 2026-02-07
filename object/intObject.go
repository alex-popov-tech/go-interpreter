package object

import (
	"fmt"
)

type IntObject struct {
	Value int64
}

func (this IntObject) Inspect() string {
	return fmt.Sprintf("%d", this.Value)
}

func (this IntObject) Type() ObjectType {
	return INT
}

func (this IntObject) String() string {
	return this.Inspect()
}
