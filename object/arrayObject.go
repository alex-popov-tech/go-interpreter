package object

import (
	"fmt"
	"strings"
)

type ArrayObject struct {
	Items []Object
}

func (ao ArrayObject) Inspect() string {
	args := []string{}
	for _, arg := range ao.Items {
		args = append(args, arg.Inspect())
	}
	return fmt.Sprintf("[%s]", strings.Join(args, ", "))
}

func (ao ArrayObject) Type() ObjectType {
	return ARRAY
}
