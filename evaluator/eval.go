package evaluator

import (
	"fmt"
	"strings"

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
	case *ast.CallExpression:
		// resolve argumentValues
		argumentValues := []object.Object{}
		for _, a := range node.Arguments {
			argVal := Eval(scope, a)
			argVal = resolveIdentIfNeeded(scope, argVal)
			if isType(object.ERROR, argVal) {
				return argVal
			}
			argumentValues = append(argumentValues, argVal)
		}

		if builtinFn, isThere := builtins[node.FnIdentifier.String()]; isThere {
			return builtinFn.Function(argumentValues...)
		}

		calleeObj := Eval(scope, node.FnIdentifier)
		calleeObj = resolveIdentIfNeeded(scope, calleeObj)
		if isType(object.ERROR, calleeObj) {
			return calleeObj
		}

		fnObject, isFnObject := calleeObj.(*object.FnObject)
		if !isFnObject {
			return &object.ErrorObject{
				Message: &object.StringObject{
					Value: fmt.Sprintf(
						"'%s' is not a function, got %s",
						node.FnIdentifier.String(),
						calleeObj.Type(),
					),
				},
			}
		}

		// create new scope and populate it with arguments
		inner := fnObject.LexicalScope.Spawn()
		for i := 0; i < len(fnObject.Arguments); i++ {
			inner.Add(fnObject.Arguments[i].Value, argumentValues[i])
		}

		var result object.Object = object.NULL_OBJECT
		for i := range fnObject.Body.Statements {
			result = Eval(inner, fnObject.Body.Statements[i])
			if isOneOfTypes(result, object.ERROR) {
				return result
			}
			if isOneOfTypes(result, object.RETURN) {
				return result.(*object.ReturnObject).Value
			}
		}
		return result
	case *ast.ArrayExpression:
		objects := []object.Object{}
		for _, a := range node.Elements {
			obj := Eval(scope, a)
			obj = resolveIdentIfNeeded(scope, obj)
			if isType(object.ERROR, obj) {
				return obj
			}
			objects = append(objects, obj)
		}
		return &object.ArrayObject{Items: objects}
	case *ast.HashExpression:
		m := map[any]object.Object{}

		for key, val := range node.Map {
			keyObject := Eval(scope, key)
			keyObject = resolveIdentIfNeeded(scope, keyObject)
			if isType(object.ERROR, keyObject) {
				return keyObject
			}
			if !isOneOfTypes(keyObject, object.STRING, object.INT, object.BOOL) {
				return &object.ErrorObject{
					Message: &object.StringObject{
						Value: fmt.Sprintf(
							"hash keys must be of type [%s], but was %s",
							strings.Join(
								[]string{
									string(object.STRING),
									string(object.INT),
									string(object.BOOL),
								},
								", ",
							),
							keyObject.Type(),
						),
					},
				}
			}

			valObject := Eval(scope, val)
			valObject = resolveIdentIfNeeded(scope, valObject)
			if isType(object.ERROR, valObject) {
				return valObject
			}

			switch keyObject := keyObject.(type) {
			case *object.StringObject:
				m[keyObject.Value] = valObject
			case *object.IntObject:
				m[keyObject.Value] = valObject
			case *object.BoolObject:
				m[keyObject.Value] = valObject
			default:
				return &object.ErrorObject{
					Message: &object.StringObject{
						Value: fmt.Sprintf("cannot use key for hash of type %s", keyObject.Type()),
					},
				}
			}
		}
		return &object.HashObject{Map: m}

	case *ast.IndexExpression:
		return evalIndexExpression(scope, node)

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
	case *ast.Identifier:
		return &object.IdentifierObject{Value: node.Value}
	case *ast.FnExpression:
		return &object.FnObject{Arguments: node.Arguments, Body: node.Body, LexicalScope: scope}

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
