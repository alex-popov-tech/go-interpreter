package evaluator

import (
	"monkey/ast"
	"monkey/object"
)

func evalIfExpression(scope *object.Scope, ifExpression *ast.IfExpression) object.Object {
	ifConditionResult := Eval(scope, ifExpression.Condition)
	ifConditionResult = resolveIdentIfNeeded(scope, ifConditionResult)
	if isType(object.ERROR, ifConditionResult) {
		return ifConditionResult
	}

	if convertToBoolish(ifConditionResult) {
		inner := scope.Spawn()
		res := Eval(inner, ifExpression.IfBlock)
		return res
	}

	for i := range ifExpression.ElseIfBlocks {
		elseIfConditionResult := Eval(scope, ifExpression.ElseIfBlocks[i].Condition)
		elseIfConditionResult = resolveIdentIfNeeded(scope, elseIfConditionResult)
		if isType(object.ERROR, elseIfConditionResult) {
			return elseIfConditionResult
		}
		if convertToBoolish(elseIfConditionResult) {
			inner := scope.Spawn()
			res := Eval(inner, ifExpression.ElseIfBlocks[i].Block)
			return res
		}
	}

	if ifExpression.ElseBlock != nil {
		inner := scope.Spawn()
		res := Eval(inner, ifExpression.ElseBlock)
		return res
	}
	return object.NULL_OBJECT
}
