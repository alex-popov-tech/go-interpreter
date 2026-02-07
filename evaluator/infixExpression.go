package evaluator

import (
	"fmt"
	"strings"

	"monkey/object"
)

func evalInfixExpression(
	scope *object.Scope,
	operator string,
	left, right object.Object,
) object.Object {
	resolvedLeft := resolveIdentIfNeeded(scope, left)
	resolvedRight := resolveIdentIfNeeded(scope, right)

	// Check for resolution errors (except for assignment where left stays as ident)
	if operator != "=" {
		if isType(object.ERROR, resolvedLeft) {
			return resolvedLeft
		}
		if isType(object.ERROR, resolvedRight) {
			return resolvedRight
		}
	}

	switch operator {
	case "=":
		if !isType(object.IDENT, left) {
			return makeIncorrectOperationError(operator, left, right)
		}
		ident, _ := left.(*object.IdentifierObject)
		_, isDefined := scope.Get(ident.Value)
		if !isDefined {
			return &object.ErrorObject{
				Message: &object.StringObject{Value: fmt.Sprintf("%s is not defined", ident.Value)},
			}
		}
		scope.Set(ident.Value, right)
		return right
	case "+":
		if !isOneOfTypes(resolvedLeft, object.STRING, object.INT) ||
			!isOneOfTypes(resolvedRight, object.STRING, object.INT) {
			return makeIncorrectOperationError(operator, resolvedLeft, resolvedRight)
		}
		if isType(object.INT, resolvedLeft) && isType(object.INT, resolvedRight) {
			leftInt, _ := resolvedLeft.(*object.IntObject)
			rightInt, _ := resolvedRight.(*object.IntObject)
			return &object.IntObject{Value: leftInt.Value + rightInt.Value}
		}
		return &object.StringObject{Value: resolvedLeft.Inspect() + resolvedRight.Inspect()}
	case "-":
		if !isOneOfTypes(resolvedLeft, object.INT) || !isOneOfTypes(resolvedRight, object.INT) {
			return makeIncorrectOperationError(operator, resolvedLeft, resolvedRight)
		}
		leftInt, _ := resolvedLeft.(*object.IntObject)
		rightInt, _ := resolvedRight.(*object.IntObject)
		return &object.IntObject{Value: leftInt.Value - rightInt.Value}
	case "*":
		if !isOneOfTypes(resolvedLeft, object.STRING, object.INT) ||
			!isOneOfTypes(resolvedRight, object.STRING, object.INT) {
			return makeIncorrectOperationError(operator, resolvedLeft, resolvedRight)
		}
		if isType(object.STRING, resolvedLeft) && isType(object.STRING, resolvedRight) {
			return makeIncorrectOperationError(operator, resolvedLeft, resolvedRight)
		}
		if isType(object.INT, resolvedLeft) && isType(object.INT, resolvedRight) {
			leftInt, _ := resolvedLeft.(*object.IntObject)
			rightInt, _ := resolvedRight.(*object.IntObject)
			return &object.IntObject{Value: leftInt.Value * rightInt.Value}
		}
		leftInt, isLeftInt := resolvedLeft.(*object.IntObject)
		rightInt, _ := resolvedRight.(*object.IntObject)
		leftString, _ := resolvedLeft.(*object.StringObject)
		rightString, _ := resolvedRight.(*object.StringObject)
		if isLeftInt {
			return &object.StringObject{
				Value: strings.Repeat(rightString.Value, int(leftInt.Value)),
			}
		} else {
			return &object.StringObject{
				Value: strings.Repeat(leftString.Value, int(rightInt.Value)),
			}
		}
	case "/":
		if !isOneOfTypes(resolvedLeft, object.INT) || !isOneOfTypes(resolvedRight, object.INT) {
			return makeIncorrectOperationError(operator, resolvedLeft, resolvedRight)
		}
		leftInt, _ := resolvedLeft.(*object.IntObject)
		rightInt, _ := resolvedRight.(*object.IntObject)
		return &object.IntObject{Value: leftInt.Value / rightInt.Value}
	case "<":
		if !isOneOfTypes(resolvedLeft, object.INT) || !isOneOfTypes(resolvedRight, object.INT) {
			return makeIncorrectOperationError(operator, resolvedLeft, resolvedRight)
		}
		leftInt, _ := resolvedLeft.(*object.IntObject)
		rightInt, _ := resolvedRight.(*object.IntObject)
		return makeBoolObject(leftInt.Value < rightInt.Value)

	case ">":
		if !isOneOfTypes(resolvedLeft, object.INT) || !isOneOfTypes(resolvedRight, object.INT) {
			return makeIncorrectOperationError(operator, resolvedLeft, resolvedRight)
		}
		leftInt, _ := resolvedLeft.(*object.IntObject)
		rightInt, _ := resolvedRight.(*object.IntObject)
		return makeBoolObject(leftInt.Value > rightInt.Value)
	case "&&":
		return &object.BoolObject{
			Value: convertToBoolish(resolvedLeft) && convertToBoolish(resolvedRight),
		}
	case "||":
		return &object.BoolObject{
			Value: convertToBoolish(resolvedLeft) || convertToBoolish(resolvedRight),
		}
	case "==":
		leftBool, isLeftBool := resolvedLeft.(*object.BoolObject)
		rightBool, isRightBool := resolvedRight.(*object.BoolObject)
		leftInt, isLeftInt := resolvedLeft.(*object.IntObject)
		rightInt, isRightInt := resolvedRight.(*object.IntObject)
		leftString, isLeftString := resolvedLeft.(*object.StringObject)
		rightString, isRightString := resolvedRight.(*object.StringObject)
		switch {
		case (isLeftInt && isRightInt):
			return &object.BoolObject{Value: leftInt.Value == rightInt.Value}
		case (isLeftBool && isRightBool):
			return &object.BoolObject{Value: leftBool.Value == rightBool.Value}
		case (isLeftString && isRightString):
			return &object.BoolObject{Value: leftString.Value == rightString.Value}
		default:
			return makeIncorrectOperationError(operator, resolvedLeft, resolvedRight)
		}
	case "!=":
		leftBool, isLeftBool := resolvedLeft.(*object.BoolObject)
		rightBool, isRightBool := resolvedRight.(*object.BoolObject)
		leftInt, isLeftInt := resolvedLeft.(*object.IntObject)
		rightInt, isRightInt := resolvedRight.(*object.IntObject)
		leftString, isLeftString := resolvedLeft.(*object.StringObject)
		rightString, isRightString := resolvedRight.(*object.StringObject)
		switch {
		case (isLeftInt && isRightInt):
			return &object.BoolObject{Value: leftInt.Value != rightInt.Value}
		case (isLeftBool && isRightBool):
			return &object.BoolObject{Value: leftBool.Value != rightBool.Value}
		case (isLeftString && isRightString):
			return &object.BoolObject{Value: leftString.Value != rightString.Value}
		default:
			return makeIncorrectOperationError(operator, resolvedLeft, resolvedRight)
		}
	default:
		return makeIncorrectOperationError(operator, resolvedLeft, resolvedRight)
	}
}
