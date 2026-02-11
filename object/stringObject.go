package object

type StringObject struct {
	Value string
}

func (this StringObject) Inspect() string {
	return this.Value
}

func (this StringObject) Type() ObjectType {
	return STRING
}
