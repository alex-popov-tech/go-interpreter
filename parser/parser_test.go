package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"monkey/ast"
	"monkey/lexer"
)

// =============================================================================
// Statement Tests
// =============================================================================
func parseStatements(input string) (statements []ast.Statement, errors []string) {
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	stringErrors := make([]string, len(p.Errors()))
	for i, err := range p.Errors() {
		stringErrors[i] = err.Error()
	}
	return program.Statements, stringErrors
}

func TestLetStatement(t *testing.T) {
	statements, errors := parseStatements("let x = 5;")
	require.Empty(t, errors)

	require.Len(t, statements, 1)
	s := statements[0].(*ast.LetStatement)
	require.IsType(t, &ast.LetStatement{}, s)
	assert.Equal(t, "x", s.Identifier.Value)
	assert.Equal(t, int64(5), s.Value.(*ast.IntLiteral).Value)
}

func TestLetStatementValidation(t *testing.T) {
	// missing identifier
	_, errors := parseStatements("let = 5;")
	require.Equal(t, []string{"expected IDENT, got ="}, errors)

	// missing '='
	_, errors = parseStatements("let nice 5;")
	require.Equal(t, []string{"expected =, got INT"}, errors)
}

func TestReturnStatement(t *testing.T) {
	statements, errors := parseStatements("return hello;")
	require.Empty(t, errors)

	require.Len(t, statements, 1)
	s := statements[0].(*ast.ReturnStatement)
	require.IsType(t, &ast.ReturnStatement{}, s)
	assert.Equal(t, "hello", s.Value.(*ast.Identifier).Value)
}

func TestBlockExpression(t *testing.T) {
	statements, errors := parseStatements(`{
		let x = 5;
		let y = 10;
		return x + y;
	};`)
	require.Empty(t, errors)

	require.Len(t, statements, 1)
	es := statements[0].(*ast.ExpressionStatement)
	require.IsType(t, &ast.ExpressionStatement{}, es)
	bs := es.Expression.(*ast.BlockExpression)
	require.IsType(t, &ast.BlockExpression{}, bs)

	require.IsType(t, &ast.LetStatement{}, bs.Statements[0])
	require.Equal(t, "x", bs.Statements[0].(*ast.LetStatement).Identifier.Value)
	require.IsType(t, &ast.IntLiteral{}, bs.Statements[0].(*ast.LetStatement).Value)
	require.IsType(t, int64(5), bs.Statements[0].(*ast.LetStatement).Value.(*ast.IntLiteral).Value)

	require.IsType(t, &ast.LetStatement{}, bs.Statements[1])
	require.Equal(t, "y", bs.Statements[1].(*ast.LetStatement).Identifier.Value)
	require.IsType(t, &ast.IntLiteral{}, bs.Statements[1].(*ast.LetStatement).Value)
	require.IsType(t, int64(10), bs.Statements[1].(*ast.LetStatement).Value.(*ast.IntLiteral).Value)

	require.IsType(t, &ast.ReturnStatement{}, bs.Statements[2])
	returnStmt := bs.Statements[2].(*ast.ReturnStatement)
	require.IsType(t, &ast.InfixExpression{}, returnStmt.Value)
	infixExpr := returnStmt.Value.(*ast.InfixExpression)
	require.IsType(t, &ast.Identifier{}, infixExpr.Left)
	require.Equal(t, "x", infixExpr.Left.(*ast.Identifier).Value)
	require.Equal(t, "+", infixExpr.Operator)
	require.IsType(t, &ast.Identifier{}, infixExpr.Right)
	require.Equal(t, "y", infixExpr.Right.(*ast.Identifier).Value)
}

func TestEmptyBlockExpression(t *testing.T) {
	statements, errors := parseStatements(`{};`)
	require.Empty(t, errors)

	require.Len(t, statements, 1)
	es := statements[0].(*ast.ExpressionStatement)
	require.IsType(t, &ast.ExpressionStatement{}, es)
	bs := es.Expression.(*ast.BlockExpression)
	require.IsType(t, &ast.BlockExpression{}, bs)
	require.Empty(t, bs.Statements)
}

func TestIfExpression(t *testing.T) {
	statements, errors := parseStatements(`
		if (x < 0) {
			return negative;
		} else if (x == 0) {
			return zero;
		} else if (x > 0) {
			return positive;
		} else {
			return unknown;
		}
	`)
	require.Empty(t, errors)

	require.Len(t, statements, 1)
	expressionStatement := statements[0].(*ast.ExpressionStatement)
	require.IsType(t, &ast.ExpressionStatement{}, expressionStatement)
	ifExpr := expressionStatement.Expression.(*ast.IfExpression)
	require.IsType(t, &ast.IfExpression{}, ifExpr)

	// Verify if condition: x < 0
	require.IsType(t, &ast.InfixExpression{}, ifExpr.Condition)
	ifCond := ifExpr.Condition.(*ast.InfixExpression)
	require.Equal(t, "x", ifCond.Left.(*ast.Identifier).Value)
	require.Equal(t, "<", ifCond.Operator)
	require.Equal(t, int64(0), ifCond.Right.(*ast.IntLiteral).Value)

	// Verify if-block: return negative;
	require.Len(t, ifExpr.IfBlock.Statements, 1)
	ifReturn := ifExpr.IfBlock.Statements[0].(*ast.ReturnStatement)
	require.Equal(t, "negative", ifReturn.Value.(*ast.Identifier).Value)

	// Verify else-if blocks
	require.Len(t, ifExpr.ElseIfBlocks, 2)

	// First else-if: else if (x == 0) { return zero; }
	elseIf1 := ifExpr.ElseIfBlocks[0]
	require.IsType(t, &ast.InfixExpression{}, elseIf1.Condition)
	elseIf1Cond := elseIf1.Condition.(*ast.InfixExpression)
	require.Equal(t, "x", elseIf1Cond.Left.(*ast.Identifier).Value)
	require.Equal(t, "==", elseIf1Cond.Operator)
	require.Equal(t, int64(0), elseIf1Cond.Right.(*ast.IntLiteral).Value)
	require.Len(t, elseIf1.Block.Statements, 1)
	elseIf1Return := elseIf1.Block.Statements[0].(*ast.ReturnStatement)
	require.Equal(t, "zero", elseIf1Return.Value.(*ast.Identifier).Value)

	// Second else-if: else if (x > 0) { return positive; }
	elseIf2 := ifExpr.ElseIfBlocks[1]
	require.IsType(t, &ast.InfixExpression{}, elseIf2.Condition)
	elseIf2Cond := elseIf2.Condition.(*ast.InfixExpression)
	require.Equal(t, "x", elseIf2Cond.Left.(*ast.Identifier).Value)
	require.Equal(t, ">", elseIf2Cond.Operator)
	require.Equal(t, int64(0), elseIf2Cond.Right.(*ast.IntLiteral).Value)
	require.Len(t, elseIf2.Block.Statements, 1)
	elseIf2Return := elseIf2.Block.Statements[0].(*ast.ReturnStatement)
	require.Equal(t, "positive", elseIf2Return.Value.(*ast.Identifier).Value)

	// Verify else-block: return unknown;
	require.NotNil(t, ifExpr.ElseBlock)
	require.Len(t, ifExpr.ElseBlock.Statements, 1)
	elseReturn := ifExpr.ElseBlock.Statements[0].(*ast.ReturnStatement)
	require.Equal(t, "unknown", elseReturn.Value.(*ast.Identifier).Value)
}

func TestFnExpression(t *testing.T) {
	statements, errors := parseStatements(`fn(a, b, c) { return a + b + c; };`)
	require.Empty(t, errors)

	require.Len(t, statements, 1)
	expressionStatement := statements[0].(*ast.ExpressionStatement)
	require.IsType(t, &ast.ExpressionStatement{}, expressionStatement)
	fnExpr := expressionStatement.Expression.(*ast.FnExpression)
	require.IsType(t, &ast.FnExpression{}, fnExpr)

	// Verify arguments: a, b, c
	require.Len(t, fnExpr.Arguments, 3)
	require.Equal(t, "a", fnExpr.Arguments[0].Value)
	require.Equal(t, "b", fnExpr.Arguments[1].Value)
	require.Equal(t, "c", fnExpr.Arguments[2].Value)

	// Verify body: return a + b + c;
	require.Len(t, fnExpr.Body.Statements, 1)
	returnStmt := fnExpr.Body.Statements[0].(*ast.ReturnStatement)
	require.IsType(t, &ast.ReturnStatement{}, returnStmt)

	// Verify return value is infix expression chain: (a + b) + c
	outerInfix := returnStmt.Value.(*ast.InfixExpression)
	require.Equal(t, "+", outerInfix.Operator)
	require.Equal(t, "c", outerInfix.Right.(*ast.Identifier).Value)

	innerInfix := outerInfix.Left.(*ast.InfixExpression)
	require.Equal(t, "+", innerInfix.Operator)
	require.Equal(t, "a", innerInfix.Left.(*ast.Identifier).Value)
	require.Equal(t, "b", innerInfix.Right.(*ast.Identifier).Value)
}

func TestEmptyFnExpression(t *testing.T) {
	statements, errors := parseStatements(`fn() { return 42; };`)
	require.Empty(t, errors)

	require.Len(t, statements, 1)
	expressionStatement := statements[0].(*ast.ExpressionStatement)
	fnExpr := expressionStatement.Expression.(*ast.FnExpression)
	require.IsType(t, &ast.FnExpression{}, fnExpr)

	// Verify empty arguments
	require.Empty(t, fnExpr.Arguments)

	// Verify body: return 42;
	require.Len(t, fnExpr.Body.Statements, 1)
	returnStmt := fnExpr.Body.Statements[0].(*ast.ReturnStatement)
	require.Equal(t, int64(42), returnStmt.Value.(*ast.IntLiteral).Value)
}

func TestFnExpressionEmptyBody(t *testing.T) {
	statements, errors := parseStatements(`fn(x) {};`)
	require.Empty(t, errors)

	require.Len(t, statements, 1)
	expressionStatement := statements[0].(*ast.ExpressionStatement)
	fnExpr := expressionStatement.Expression.(*ast.FnExpression)
	require.IsType(t, &ast.FnExpression{}, fnExpr)

	// Verify single argument
	require.Len(t, fnExpr.Arguments, 1)
	require.Equal(t, "x", fnExpr.Arguments[0].Value)

	// Verify empty body
	require.Empty(t, fnExpr.Body.Statements)
}

func TestCallExpression(t *testing.T) {
	statements, errors := parseStatements(`add(a, b, c);`)
	require.Empty(t, errors)

	require.Len(t, statements, 1)
	expressionStatement := statements[0].(*ast.ExpressionStatement)
	require.IsType(t, &ast.ExpressionStatement{}, expressionStatement)
	callExpr := expressionStatement.Expression.(*ast.CallExpression)
	require.IsType(t, &ast.CallExpression{}, callExpr)

	// Verify function identifier
	fnIdent := callExpr.FnIdentifier.(*ast.Identifier)
	require.Equal(t, "add", fnIdent.Value)

	// Verify arguments: a, b, c
	require.Len(t, callExpr.Arguments, 3)
	require.Equal(t, "a", callExpr.Arguments[0].(*ast.Identifier).Value)
	require.Equal(t, "b", callExpr.Arguments[1].(*ast.Identifier).Value)
	require.Equal(t, "c", callExpr.Arguments[2].(*ast.Identifier).Value)
}

func TestIndexExpression_Identifier(t *testing.T) {
	statements, errors := parseStatements("foo[1]")
	require.Empty(t, errors)

	require.Len(t, statements, 1)
	expressionStatement := statements[0].(*ast.ExpressionStatement)
	require.IsType(t, &ast.ExpressionStatement{}, expressionStatement)
	indexExpr := expressionStatement.Expression.(*ast.IndexExpression)
	require.IsType(t, &ast.IndexExpression{}, indexExpr)

	// Verify function identifier
	ident := indexExpr.Identifier.(*ast.Identifier)
	require.Equal(t, "foo", ident.Value)

	val := indexExpr.IndexExpression.(*ast.IntLiteral)
	require.Equal(t, int64(1), val.Value)
}

func TestIndexExpression_ArrayLiteral(t *testing.T) {
	statements, errors := parseStatements("[1, 2, 3][0]")
	require.Empty(t, errors)

	require.Len(t, statements, 1)
	expressionStatement := statements[0].(*ast.ExpressionStatement)
	indexExpr := expressionStatement.Expression.(*ast.IndexExpression)
	require.IsType(t, &ast.IndexExpression{}, indexExpr)

	// Verify left side is an array literal
	arrExpr := indexExpr.Identifier.(*ast.ArrayExpression)
	require.Len(t, arrExpr.Elements, 3)
	require.Equal(t, int64(1), arrExpr.Elements[0].(*ast.IntLiteral).Value)
	require.Equal(t, int64(2), arrExpr.Elements[1].(*ast.IntLiteral).Value)
	require.Equal(t, int64(3), arrExpr.Elements[2].(*ast.IntLiteral).Value)

	// Verify index
	idx := indexExpr.IndexExpression.(*ast.IntLiteral)
	require.Equal(t, int64(0), idx.Value)
}

func TestIndexExpression_CallExpressionIndex(t *testing.T) {
	statements, errors := parseStatements("foo[bar()]")
	require.Empty(t, errors)

	require.Len(t, statements, 1)
	expressionStatement := statements[0].(*ast.ExpressionStatement)
	indexExpr := expressionStatement.Expression.(*ast.IndexExpression)
	require.IsType(t, &ast.IndexExpression{}, indexExpr)

	// Verify array identifier
	ident := indexExpr.Identifier.(*ast.Identifier)
	require.Equal(t, "foo", ident.Value)

	// Verify index is a call expression
	callExpr := indexExpr.IndexExpression.(*ast.CallExpression)
	require.IsType(t, &ast.CallExpression{}, callExpr)
	fnIdent := callExpr.FnIdentifier.(*ast.Identifier)
	require.Equal(t, "bar", fnIdent.Value)
	require.Empty(t, callExpr.Arguments)
}

func TestIndexExpression_InfixIndex(t *testing.T) {
	statements, errors := parseStatements("foo[a + 1]")
	require.Empty(t, errors)

	require.Len(t, statements, 1)
	expressionStatement := statements[0].(*ast.ExpressionStatement)
	indexExpr := expressionStatement.Expression.(*ast.IndexExpression)
	require.IsType(t, &ast.IndexExpression{}, indexExpr)

	// Verify array identifier
	ident := indexExpr.Identifier.(*ast.Identifier)
	require.Equal(t, "foo", ident.Value)

	// Verify index is an infix expression: a + 1
	infixExpr := indexExpr.IndexExpression.(*ast.InfixExpression)
	require.IsType(t, &ast.InfixExpression{}, infixExpr)
	require.Equal(t, "a", infixExpr.Left.(*ast.Identifier).Value)
	require.Equal(t, "+", infixExpr.Operator)
	require.Equal(t, int64(1), infixExpr.Right.(*ast.IntLiteral).Value)
}

func TestIndexExpression_Nested(t *testing.T) {
	statements, errors := parseStatements("foo[0][1]")
	require.Empty(t, errors)

	require.Len(t, statements, 1)
	expressionStatement := statements[0].(*ast.ExpressionStatement)

	// Outer: (...)[1]
	outerIndex := expressionStatement.Expression.(*ast.IndexExpression)
	require.IsType(t, &ast.IndexExpression{}, outerIndex)
	require.Equal(t, int64(1), outerIndex.IndexExpression.(*ast.IntLiteral).Value)

	// Inner: foo[0]
	innerIndex := outerIndex.Identifier.(*ast.IndexExpression)
	require.IsType(t, &ast.IndexExpression{}, innerIndex)
	require.Equal(t, "foo", innerIndex.Identifier.(*ast.Identifier).Value)
	require.Equal(t, int64(0), innerIndex.IndexExpression.(*ast.IntLiteral).Value)
}

func TestIndexExpression_CallResultIndexed(t *testing.T) {
	statements, errors := parseStatements("bar()[0]")
	require.Empty(t, errors)

	require.Len(t, statements, 1)
	expressionStatement := statements[0].(*ast.ExpressionStatement)
	indexExpr := expressionStatement.Expression.(*ast.IndexExpression)
	require.IsType(t, &ast.IndexExpression{}, indexExpr)

	// Verify left side is a call expression
	callExpr := indexExpr.Identifier.(*ast.CallExpression)
	require.IsType(t, &ast.CallExpression{}, callExpr)
	fnIdent := callExpr.FnIdentifier.(*ast.Identifier)
	require.Equal(t, "bar", fnIdent.Value)
	require.Empty(t, callExpr.Arguments)

	// Verify index
	idx := indexExpr.IndexExpression.(*ast.IntLiteral)
	require.Equal(t, int64(0), idx.Value)
}

func TestEmptyCallExpression(t *testing.T) {
	statements, errors := parseStatements(`foo();`)
	require.Empty(t, errors)

	require.Len(t, statements, 1)
	expressionStatement := statements[0].(*ast.ExpressionStatement)
	callExpr := expressionStatement.Expression.(*ast.CallExpression)
	require.IsType(t, &ast.CallExpression{}, callExpr)

	// Verify function identifier
	fnIdent := callExpr.FnIdentifier.(*ast.Identifier)
	require.Equal(t, "foo", fnIdent.Value)

	// Verify empty arguments
	require.Empty(t, callExpr.Arguments)
}

func TestCallExpressionSingleArg(t *testing.T) {
	statements, errors := parseStatements(`print(x);`)
	require.Empty(t, errors)

	require.Len(t, statements, 1)
	expressionStatement := statements[0].(*ast.ExpressionStatement)
	callExpr := expressionStatement.Expression.(*ast.CallExpression)
	require.IsType(t, &ast.CallExpression{}, callExpr)

	// Verify function identifier
	fnIdent := callExpr.FnIdentifier.(*ast.Identifier)
	require.Equal(t, "print", fnIdent.Value)

	// Verify single argument
	require.Len(t, callExpr.Arguments, 1)
	require.Equal(t, "x", callExpr.Arguments[0].(*ast.Identifier).Value)
}

func TestCallExpressionWithFnLiteral(t *testing.T) {
	statements, errors := parseStatements(`fn(a, b, c) { return a + b + c; }(x, y, z);`)
	require.Empty(t, errors)

	require.Len(t, statements, 1)
	expressionStatement := statements[0].(*ast.ExpressionStatement)
	callExpr := expressionStatement.Expression.(*ast.CallExpression)
	require.IsType(t, &ast.CallExpression{}, callExpr)

	// Verify function is an FnExpression, not an Identifier
	fnExpr := callExpr.FnIdentifier.(*ast.FnExpression)
	require.IsType(t, &ast.FnExpression{}, fnExpr)

	// Verify fn parameters: a, b, c
	require.Len(t, fnExpr.Arguments, 3)
	require.Equal(t, "a", fnExpr.Arguments[0].Value)
	require.Equal(t, "b", fnExpr.Arguments[1].Value)
	require.Equal(t, "c", fnExpr.Arguments[2].Value)

	// Verify call arguments: x, y, z
	require.Len(t, callExpr.Arguments, 3)
	require.Equal(t, "x", callExpr.Arguments[0].(*ast.Identifier).Value)
	require.Equal(t, "y", callExpr.Arguments[1].(*ast.Identifier).Value)
	require.Equal(t, "z", callExpr.Arguments[2].(*ast.Identifier).Value)
}

// =============================================================================
// Expression Tests
// =============================================================================

func TestIdentifierExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"foo;", "foo"},
		{"hello_mom;", "hello_mom"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			statements, errors := parseStatements(tt.input)
			require.Empty(t, errors)

			require.Len(t, statements, 1)
			s := statements[0].(*ast.ExpressionStatement)
			require.IsType(t, ast.ExpressionStatement{}, *s)

			require.IsType(t, &ast.Identifier{}, s.Expression)
			require.Equal(t, tt.expected, s.Expression.(*ast.Identifier).Value)
		})
	}
}

func TestIntLiteralExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5;", 5},
		{"10_000;", 10000},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			statements, errors := parseStatements(tt.input)
			require.Empty(t, errors)

			require.Len(t, statements, 1)
			s := statements[0].(*ast.ExpressionStatement)
			require.IsType(t, ast.ExpressionStatement{}, *s)

			require.IsType(t, &ast.IntLiteral{}, s.Expression)
			require.Equal(t, tt.expected, s.Expression.(*ast.IntLiteral).Value)
		})
	}
}

func TestStringLiteralExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"'hello mom';", "hello mom"},
		{"'hello \\'mom\\'';", "hello 'mom'"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			statements, errors := parseStatements(tt.input)
			require.Empty(t, errors)

			require.Len(t, statements, 1)
			s := statements[0].(*ast.ExpressionStatement)
			require.IsType(t, ast.ExpressionStatement{}, *s)

			require.IsType(t, &ast.StringLiteral{}, s.Expression)
			require.Equal(t, tt.expected, s.Expression.(*ast.StringLiteral).Value)
		})
	}
}

func TestBoolLiteralExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true;", true},
		{"false;", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			statements, errors := parseStatements(tt.input)
			require.Empty(t, errors)

			require.Len(t, statements, 1)
			s := statements[0].(*ast.ExpressionStatement)
			require.IsType(t, ast.ExpressionStatement{}, *s)

			require.IsType(t, &ast.BoolLiteral{}, s.Expression)
			require.Equal(t, tt.expected, s.Expression.(*ast.BoolLiteral).Value)
		})
	}
}

func TestEmptyArrayExpression(t *testing.T) {
	statements, errors := parseStatements(`[];`)
	require.Empty(t, errors)

	require.Len(t, statements, 1)
	expressionStatement := statements[0].(*ast.ExpressionStatement)
	arrExpr := expressionStatement.Expression.(*ast.ArrayExpression)
	require.IsType(t, &ast.ArrayExpression{}, arrExpr)

	require.Empty(t, arrExpr.Elements)
}

func TestArrayExpressionSingleArg(t *testing.T) {
	statements, errors := parseStatements(`[ a ];`)
	require.Empty(t, errors)

	require.Len(t, statements, 1)
	expressionStatement := statements[0].(*ast.ExpressionStatement)
	arrExpr := expressionStatement.Expression.(*ast.ArrayExpression)
	require.IsType(t, &ast.ArrayExpression{}, arrExpr)

	require.Len(t, arrExpr.Elements, 1)
	require.IsType(t, &ast.Identifier{}, arrExpr.Elements[0])
	require.IsType(t, &ast.Identifier{}, arrExpr.Elements[0])
	require.Equal(t, "a", arrExpr.Elements[0].(*ast.Identifier).Value)
}

func TestArrayExpressionMultipleItems(t *testing.T) {
	statements, errors := parseStatements(`[a, 1, true];`)
	require.Empty(t, errors)

	require.Len(t, statements, 1)
	expressionStatement := statements[0].(*ast.ExpressionStatement)
	arrExpr := expressionStatement.Expression.(*ast.ArrayExpression)
	require.IsType(t, &ast.ArrayExpression{}, arrExpr)

	require.Len(t, arrExpr.Elements, 3)

	require.IsType(t, &ast.Identifier{}, arrExpr.Elements[0])
	require.Equal(t, "a", arrExpr.Elements[0].(*ast.Identifier).Value)

	require.IsType(t, &ast.IntLiteral{}, arrExpr.Elements[1])
	require.Equal(t, int64(1), arrExpr.Elements[1].(*ast.IntLiteral).Value)

	require.IsType(t, &ast.BoolLiteral{}, arrExpr.Elements[2])
	require.Equal(t, true, arrExpr.Elements[2].(*ast.BoolLiteral).Value)
}

func TestEmptyHashExpression(t *testing.T) {
	statements, errors := parseStatements("#{}")
	require.Empty(t, errors)

	require.Len(t, statements, 1)
	expressionStatement := statements[0].(*ast.ExpressionStatement)
	hashExpr := expressionStatement.Expression.(*ast.HashExpression)
	require.IsType(t, &ast.HashExpression{}, hashExpr)

	require.Empty(t, hashExpr.Map)
}

func TestSingleKeyHashExpression(t *testing.T) {
	statements, errors := parseStatements("#{foo:bar}")
	require.Empty(t, errors)

	require.Len(t, statements, 1)
	expressionStatement := statements[0].(*ast.ExpressionStatement)
	hashExpr := expressionStatement.Expression.(*ast.HashExpression)
	require.IsType(t, &ast.HashExpression{}, hashExpr)
	require.Len(t, hashExpr.Map, 1)
	require.Equal(t, "#{foo:bar}", hashExpr.String())
}

func TestPrefixExpressions(t *testing.T) {
	tests := []struct {
		input    string
		operator string
		value    int64
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
		{"+15;", "+", 15},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			statements, errors := parseStatements(tt.input)
			require.Empty(t, errors)

			require.Len(t, statements, 1)
			s := statements[0].(*ast.ExpressionStatement)
			require.IsType(t, &ast.ExpressionStatement{}, s)
			e := s.Expression.(*ast.PrefixExpression)
			require.IsType(t, &ast.PrefixExpression{}, e)

			require.Equal(t, tt.operator, e.Operator)
			require.IsType(t, &ast.IntLiteral{}, e.Value)
			require.Equal(t, tt.value, e.Value.(*ast.IntLiteral).Value)
		})
	}
}

func TestInfixExpressions(t *testing.T) {
	tests := []struct {
		input      string
		leftValue  int64
		operator   string
		rightValue int64
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"5 && 5;", 5, "&&", 5},
		{"5 || 5;", 5, "||", 5},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			statements, errors := parseStatements(tt.input)
			require.Empty(t, errors)

			require.Len(t, statements, 1)
			s := statements[0].(*ast.ExpressionStatement)
			require.IsType(t, &ast.ExpressionStatement{}, s)
			e := s.Expression.(*ast.InfixExpression)
			require.IsType(t, &ast.InfixExpression{}, e)

			require.IsType(t, &ast.IntLiteral{}, e.Left)
			require.Equal(t, tt.leftValue, e.Left.(*ast.IntLiteral).Value)

			require.Equal(t, tt.operator, e.Operator)

			require.IsType(t, &ast.IntLiteral{}, e.Right)
			require.Equal(t, tt.rightValue, e.Right.(*ast.IntLiteral).Value)
		})
	}
}

// =============================================================================
// Operator Precedence Tests
// =============================================================================

func TestOperatorPrecedence(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		// Basic arithmetic precedence
		{"5 + 5;", "(5 + 5);"},
		{"5 * 5 + 5;", "((5 * 5) + 5);"},
		{"5 + 5 * 5;", "(5 + (5 * 5));"},
		{"5 - 5 * 5;", "(5 - (5 * 5));"},
		{"5 / 5 + 5;", "((5 / 5) + 5);"},

		// Grouping
		{"(5 + 5) * 5;", "((5 + 5) * 5);"},
		{"-(5 + 5);", "(-(5 + 5));"},

		// Left-to-right associativity
		{"5 - 3 - 1;", "((5 - 3) - 1);"},
		{"5 + 3 + 1;", "((5 + 3) + 1);"},
		{"5 * 3 * 2;", "((5 * 3) * 2);"},
		{"5 / 3 / 2;", "((5 / 3) / 2);"},

		// Prefix operators
		{"-5 + 5;", "((-5) + 5);"},
		{"!true;", "(!true);"},
		{"!false;", "(!false);"},
		{"-a * b;", "((-a) * b);"},

		// Comparison operators
		{"5 > 3 == true;", "((5 > 3) == true);"},
		{"5 < 3 == false;", "((5 < 3) == false);"},
		{"3 + 4 * 5 == 3 * 1 + 4 * 5;", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)));"},

		// Mixed operators
		{"1 + 2 * 3 + 4 / 5 - 6;", "(((1 + (2 * 3)) + (4 / 5)) - 6);"},
		{"5 > 4 == 3 < 4;", "((5 > 4) == (3 < 4));"},
		{"5 > 4 != 3 < 4;", "((5 > 4) != (3 < 4));"},

		// Boolean expressions
		{"true == true;", "(true == true);"},
		{"true != false;", "(true != false);"},
		{"!true == false;", "((!true) == false);"},

		// Logical operators (Note: && and || have higher precedence than arithmetic in this implementation)
		{"true && false;", "(true && false);"},
		{"true || false;", "(true || false);"},
		{"true && false || true;", "((true && false) || true);"},
		{"true || false && true;", "((true || false) && true);"},
		{"!true && false;", "((!true) && false);"},
		// With current precedence: && and || bind tighter than + and *
		{"a + b && c + d;", "((a + (b && c)) + d);"},
		{"a * b || c * d;", "((a * (b || c)) * d);"},

		// Unary plus
		{"+5;", "(+5);"},
		{"+5 + 3;", "((+5) + 3);"},

		// Assignment operator (lowest precedence, right-associative)
		{"x = 5;", "(x = 5);"},
		{"x = 5 + 3;", "(x = (5 + 3));"},
		{"x = y = 5;", "(x = (y = 5));"},
		{"x = 5 * 3 + 2;", "(x = ((5 * 3) + 2));"},

		// Index expressions (highest precedence)
		{"a + b[0];", "(a + b[0]);"},
		{"a * b[0];", "(a * b[0]);"},
		{"bar()[0];", "bar()[0];"},
		{"foo[0][1];", "foo[0][1];"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			sattements, errors := parseStatements(tt.input)
			require.Empty(t, errors)
			require.Len(t, sattements, 1)

			require.Equal(t, tt.expected, sattements[0].String())
		})
	}
}

// =============================================================================
// Edge Cases
// =============================================================================

func TestWhitespaceOnlyInput(t *testing.T) {
	statements, errors := parseStatements("   \n\t  ")
	require.Empty(t, errors)
	require.Empty(t, statements)
}

func TestIndexExpression_MissingClosingBracket(t *testing.T) {
	_, errors := parseStatements("foo[1")
	assert.NotEmpty(t, errors, "expected parser error for missing closing ']'")
}

func TestIndexExpression_EmptyIndex(t *testing.T) {
	_, errors := parseStatements("foo[]")
	assert.NotEmpty(t, errors, "expected parser error for empty index expression")
}

func TestIndexExpression_ArrayLiteralMissingClosingBracket(t *testing.T) {
	_, errors := parseStatements("[1, 2, 3][0")
	assert.NotEmpty(t, errors, "expected parser error for missing closing ']' on index")
}

func TestIndexExpression_MissingClosingBracketWithTrailing(t *testing.T) {
	_, errors := parseStatements("foo[1 + 2")
	assert.NotEmpty(t, errors, "expected parser error for missing closing ']' with infix index")
}
