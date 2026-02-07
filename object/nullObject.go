package object

type NullObject struct{}

func (this NullObject) Inspect() string {
	return "null"
}

func (this NullObject) Type() ObjectType {
	return NULL
}
