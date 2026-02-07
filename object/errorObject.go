package object

type ErrorObject struct {
	Message Object
}

func (this ErrorObject) Inspect() string {
	return this.Message.Inspect()
}

func (this ErrorObject) Type() ObjectType {
	return ERROR
}
