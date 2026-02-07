package evaluator

import (
	"fmt"

	"monkey/object"
)

func convertToBoolish(it object.Object) bool {
	switch it := it.(type) {
	case *object.IntObject:
		switch it.Value {
		case 0:
			return false
		default:
			return true
		}
	case *object.BoolObject:
		return it.Value
	case *object.StringObject:
		if it.Value == "" {
			return false
		}
		return true
	default:
		panic(
			fmt.Sprintf("you tried to convert to bool unknown type %T", it),
		)
	}
}

func isType(t object.ObjectType, actual ...object.Object) bool {
	for i := range actual {
		if actual[i] == nil || actual[i].Type() != t {
			return false
		}
	}
	return true
}

func isOneOfTypes(val object.Object, types ...object.ObjectType) bool {
	for i := range types {
		if types[i] == val.Type() {
			return true
		}
	}
	return false
}

func isArithmeticOperator(operator string) bool {
	switch operator {
	case "+", "-", "*", "/":
		return true
	default:
		return false
	}
}

func isBoolOperator(operator string) bool {
	switch operator {
	case "&&", "||", "==", "!=":
		return true
	default:
		return false
	}
}

func makeIncorrectOperationError(operator string, left, right object.Object) *object.ErrorObject {
	return &object.ErrorObject{
		Message: &object.StringObject{
			Value: fmt.Sprintf(
				"cannot perform operation '%s %s %s'",
				left.Type(),
				operator,
				right.Type(),
			),
		},
	}
}

func makeBoolObject(val bool) *object.BoolObject {
	if val {
		return &object.TRUE_OBJECT
	} else {
		return &object.FALSE_OBJECT
	}
}

func resolveIdentIfNeeded(scope *object.Scope, o object.Object) object.Object {
	ident, isIdent := o.(*object.IdentifierObject)
	if isIdent {
		res, found := scope.Get(ident.Value)
		if !found {
			return &object.ErrorObject{
				Message: &object.StringObject{
					Value: fmt.Sprintf("identifier %s not found", ident.Value),
				},
			}
		}
		return res
	}
	return o
}
