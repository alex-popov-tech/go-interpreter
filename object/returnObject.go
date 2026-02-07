package object

type ReturnObject struct {
	Value Object
}

func (this ReturnObject) Inspect() string {
	return this.Value.Inspect()
}

func (this ReturnObject) Type() ObjectType {
	return RETURN
}
