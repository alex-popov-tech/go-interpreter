package evaluator

import (
	"monkey/object"
)

func evalPrefixExpression(operator string, value object.Object) object.Object {
	switch it := value.(type) {
	case *object.IntObject:
		switch operator {
		case "+":
			return it
		case "-":
			return &object.IntObject{Value: -it.Value}
		case "!":
			return &object.BoolObject{Value: !convertToBoolish(it)}
		default:
			return &object.ErrorObject{
				Message: &object.StringObject{
					Value: "Unsupported operator: " + operator + " for Ints",
				},
			}
		}
	case *object.BoolObject:
		switch operator {
		case "!":
			return &object.BoolObject{Value: !it.Value}
		default:
			return &object.ErrorObject{
				Message: &object.StringObject{
					Value: "Unsupported operator: " + operator + " for Bools",
				},
			}
		}
	default:
		return &object.ErrorObject{
			Message: &object.StringObject{
				Value: "Unsupported operator: " + operator + " for Bools",
			},
		}
	case *object.StringObject:
		switch operator {
		case "!":
			return &object.BoolObject{Value: !convertToBoolish(it)}
		default:
			return &object.ErrorObject{
				Message: &object.StringObject{
					Value: "Unsupported operator: " + operator + " for Strings",
				},
			}
		}
	}
}
