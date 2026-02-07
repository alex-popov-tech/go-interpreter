package object

import (
	"fmt"
)

type StringObject struct {
	Value string
}

func (this StringObject) Inspect() string {
	return fmt.Sprintf("%s", this.Value)
}

func (this StringObject) Type() ObjectType {
	return STRING
}
