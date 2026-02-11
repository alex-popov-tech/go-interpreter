package evaluator

import (
	"fmt"
	"strings"

	"monkey/ast"
	"monkey/object"
)

func evalIndexExpression(scope *object.Scope, node *ast.IndexExpression) object.Object {
	source := Eval(scope, node.Identifier)
	source = resolveIdentIfNeeded(scope, source)
	if !isOneOfTypes(source, object.ARRAY, object.STRING, object.HASH) {
		return &object.ErrorObject{
			Message: &object.StringObject{
				Value: fmt.Sprintf("can index only [%s], but was %s",
					strings.Join(
						[]string{
							string(object.STRING),
							string(object.ARRAY),
							string(object.HASH),
						},
						", ",
					),
					source.Type()),
			},
		}
	}

	index := Eval(scope, node.IndexExpression)
	index = resolveIdentIfNeeded(scope, index)

	switch source := source.(type) {
	case *object.ArrayObject, *object.StringObject:
		if !isType(object.INT, index) {
			return &object.ErrorObject{
				Message: &object.StringObject{
					Value: fmt.Sprintf("index must be INT, but was %s", index.Type()),
				},
			}
		}
		intObject := index.(*object.IntObject)

		if sourceArr, isArray := source.(*object.ArrayObject); isArray {
			if intObject.Value < 0 || intObject.Value >= int64(len(sourceArr.Items)) {
				return &object.ErrorObject{
					Message: &object.StringObject{
						Value: fmt.Sprintf("index %d out of bounds for array of length %d", intObject.Value, len(sourceArr.Items)),
					},
				}
			}
			return sourceArr.Items[intObject.Value]
		} else if sourceStr, isString := source.(*object.StringObject); isString {
			runes := []rune(sourceStr.Value)
			if intObject.Value < 0 || intObject.Value >= int64(len(runes)) {
				return &object.ErrorObject{
					Message: &object.StringObject{
						Value: fmt.Sprintf("index %d out of bounds for string of length %d", intObject.Value, len(runes)),
					},
				}
			}
			return &object.StringObject{Value: string(runes[intObject.Value])}
		}
		panic("Indexed type is not array or string")
	case *object.HashObject:
		if !isOneOfTypes(index, object.STRING, object.INT, object.BOOL) {
			return &object.ErrorObject{
				Message: &object.StringObject{
					Value: fmt.Sprintf(
						"index must be [%s], but was %s",
						strings.Join(
							[]string{
								string(object.STRING),
								string(object.INT),
								string(object.BOOL),
							},
							", ",
						),
						index.Type(),
					),
				},
			}
		}
		switch index := index.(type) {
		case *object.StringObject:
			res, found := source.Map[index.Value]
			if !found {
				return object.NULL_OBJECT
			}
			return res
		case *object.IntObject:
			res, found := source.Map[index.Value]
			if !found {
				return object.NULL_OBJECT
			}
			return res
		case *object.BoolObject:
			res, found := source.Map[index.Value]
			if !found {
				return object.NULL_OBJECT
			}
			return res
		default:
			return &object.ErrorObject{
				Message: &object.StringObject{
					Value: fmt.Sprintf("cannot index %s", index.Type()),
				},
			}
		}
	default:
		return &object.ErrorObject{
			Message: &object.StringObject{
				Value: fmt.Sprintf("cannot index %s", source.Type()),
			},
		}
	}
}
