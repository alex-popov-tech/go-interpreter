package object

type IfObject struct {
	Value Object
}

func (this IfObject) Inspect() string {
	return this.Value.Inspect()
}

func (this IfObject) Type() ObjectType {
	return IF
}
