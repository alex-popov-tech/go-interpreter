package object

type IdentifierObject struct {
	Value string
}

func (this IdentifierObject) Inspect() string {
	return this.Value
}

func (this IdentifierObject) Type() ObjectType {
	return IDENT
}
