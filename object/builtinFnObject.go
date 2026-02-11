package object

import (
	"fmt"
)

type BuiltinFnObject struct {
	Name     string
	Function func(args ...Object) Object
}

func (this BuiltinFnObject) Inspect() string {
	return fmt.Sprintf("fn %s(...) { ...builtin... }", this.Name)
}

func (this BuiltinFnObject) Type() ObjectType {
	return BUILTIN_FN
}
