package object

type ObjectType string

type Object interface {
	Inspect() string
	Type() ObjectType
}

var (
	INT    = ObjectType("INT")
	BOOL   = ObjectType("BOOL")
	NULL   = ObjectType("NULL")
	IF     = ObjectType("IF")
	RETURN = ObjectType("RETURN")
	ERROR  = ObjectType("ERROR")
	STRING = ObjectType("STRING")
	IDENT  = ObjectType("IDENT")
	FN     = ObjectType("FN")
)

var (
	NULL_OBJECT  = NullObject{}
	TRUE_OBJECT  = BoolObject{Value:true}
	FALSE_OBJECT = BoolObject{Value:false}
)
