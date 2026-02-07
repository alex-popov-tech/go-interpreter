package evaluator

import (
	"fmt"

	"monkey/ast"
	"monkey/object"
)

func Eval(scope *object.Scope, node ast.Node) object.Object {
	switch node := node.(type) {

	case *ast.Program:
		var result object.Object
		for i := range node.Statements {
			result = Eval(scope, node.Statements[i])

			if isOneOfTypes(result, object.ERROR, object.RETURN) {
				return result
			}
		}
		return result

	// Expressions
	case *ast.ExpressionStatement:
		result := Eval(scope, node.Expression)
		return resolveIdentIfNeeded(scope, result)
	case *ast.PrefixExpression:
		operator := node.Operator
		obj := Eval(scope, node.Value)
		obj = resolveIdentIfNeeded(scope, obj)

		if isType(object.ERROR, obj) {
			return obj
		}
		return evalPrefixExpression(operator, obj)
	case *ast.InfixExpression:
		leftObj := Eval(scope, node.Left)

		if isType(object.ERROR, leftObj) {
			return leftObj
		}
		operator := node.Operator
		rightObj := Eval(scope, node.Right)
		if isType(object.ERROR, rightObj) {
			return rightObj
		}
		return evalInfixExpression(scope, operator, leftObj, rightObj)
	case *ast.IfExpression:
		return evalIfExpression(scope, node)
	case *ast.BlockExpression:
		inner := scope.Spawn()
		var result object.Object = object.NULL_OBJECT
		for i := range node.Statements {
			result = Eval(inner, node.Statements[i])
			if isOneOfTypes(result, object.ERROR, object.RETURN) {
				return result
			}
		}
		return result
	case *ast.Identifier:
		return &object.IdentifierObject{Value: node.Value}
	case *ast.FnExpression:
		return &object.FnObject{Arguments: node.Arguments, Body: node.Body, LexicalScope: scope}
	case *ast.CallExpression:
		// Evaluate the callee expression to get the function object
		calleeObj := Eval(scope, node.FnIdentifier)
		calleeObj = resolveIdentIfNeeded(scope, calleeObj)
		if isType(object.ERROR, calleeObj) {
			return calleeObj
		}

		fnObject, isFnObject := calleeObj.(*object.FnObject)
		if !isFnObject {
			return &object.ErrorObject{
				Message: &object.StringObject{
					Value: fmt.Sprintf("'%s' is not a function, got %s", node.FnIdentifier.String(), calleeObj.Type()),
				},
			}
		}

		// resolve argumentValues
		var argumentValues []object.Object = []object.Object{}
		for _, a := range node.Arguments {
			argVal := Eval(scope, a)
			argVal = resolveIdentIfNeeded(scope, argVal)
			if isType(object.ERROR, argVal) {
				return argVal
			}
			argumentValues = append(argumentValues, argVal)
		}

		// create new scope and populate it with arguments
		inner := fnObject.LexicalScope.Spawn()
		for i := 0; i < len(fnObject.Arguments); i++ {
			inner.Add(fnObject.Arguments[i].Value, argumentValues[i])
		}

		var result object.Object = object.NULL_OBJECT
		for i := range fnObject.Body.Statements {
			result = Eval(inner, fnObject.Body.Statements[i])
			if isOneOfTypes(result, object.ERROR, object.RETURN) {
				return result
			}
		}
		return result
	case *ast.IntLiteral:
		return &object.IntObject{Value: node.Value}
	case *ast.BoolLiteral:
		if node.Value {
			return &object.TRUE_OBJECT
		} else {
			return &object.FALSE_OBJECT
		}
	case *ast.StringLiteral:
		return &object.StringObject{Value: node.Value}

	// Statements
	case *ast.ReturnStatement:
		result := Eval(scope, node.Value)
		result = resolveIdentIfNeeded(scope, result)
		if isType(object.ERROR, result) {
			return result
		}
		return &object.ReturnObject{Value: result}

	case *ast.LetStatement:
		ident := node.Identifier.Value
		val := Eval(scope, node.Value)
		val = resolveIdentIfNeeded(scope, val)
		if isType(object.ERROR, val) {
			return val
		}
		scope.Add(ident, val)
		return val

	default:
		return &object.ErrorObject{
			Message: &object.StringObject{Value: fmt.Sprintf("unknown node type %T", node)},
		}
	}
}
