package evaluator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
)

// =============================================================================
// Test Helper
// =============================================================================

func evaluate(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	scope := object.NewGlobalScope()
	return Eval(scope, program)
}

// =============================================================================
// Literal Tests
// =============================================================================

func TestIntLiteralEvaluation(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5;", 5},
		{"10;", 10},
		{"10_000;", 10000},
		{"999;", 999},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := evaluate(tt.input)
			require.IsType(t, &object.IntObject{}, result)
			assert.Equal(t, tt.expected, result.(*object.IntObject).Value)
		})
	}
}

func TestBoolLiteralEvaluation(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true;", true},
		{"false;", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := evaluate(tt.input)
			require.IsType(t, &object.BoolObject{}, result)
			assert.Equal(t, tt.expected, result.(*object.BoolObject).Value)
		})
	}
}

func TestStringLiteralEvaluation(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"'hello';", "hello"},
		{"'hello world';", "hello world"},
		{"'';", ""},
		{"'hello \\'mom\\'';", "hello 'mom'"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := evaluate(tt.input)
			require.IsType(t, &object.StringObject{}, result)
			assert.Equal(t, tt.expected, result.(*object.StringObject).Value)
		})
	}
}

// =============================================================================
// Prefix Expression Tests
// =============================================================================

func TestPrefixExpressionEvaluation(t *testing.T) {
	intTests := []struct {
		input    string
		expected int64
	}{
		{"-5;", -5},
		{"-(-5);", 5},
		{"+5;", 5},
		{"+(-5);", -5},
		{"+(5);", 5},
	}

	for _, tt := range intTests {
		t.Run(tt.input, func(t *testing.T) {
			result := evaluate(tt.input)
			require.IsType(t, &object.IntObject{}, result)
			assert.Equal(t, tt.expected, result.(*object.IntObject).Value)
		})
	}

	boolTests := []struct {
		input    string
		expected bool
	}{
		{"!true;", false},
		{"!false;", true},
		{"!!true;", true},
		{"!!false;", false},
		{"!5;", false},   // 5 is truthy
		{"!100;", false}, // 100 is truthy
	}

	for _, tt := range boolTests {
		t.Run(tt.input, func(t *testing.T) {
			result := evaluate(tt.input)
			require.IsType(t, &object.BoolObject{}, result)
			assert.Equal(t, tt.expected, result.(*object.BoolObject).Value)
		})
	}
}

// =============================================================================
// Infix Arithmetic Tests
// =============================================================================

func TestInfixArithmeticEvaluation(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		// Basic operations
		{"5 + 5;", 10},
		{"5 - 5;", 0},
		{"5 * 5;", 25},
		{"10 / 2;", 5},

		// Operator precedence
		{"5 + 5 * 2;", 15},
		{"(5 + 5) * 2;", 20},
		{"5 * 2 + 5;", 15},
		{"5 / 2 + 3;", 5}, // integer division: 5/2 = 2, 2+3 = 5

		// Multiple operations
		{"2 + 3 + 4;", 9},
		{"10 - 3 - 2;", 5},
		{"2 * 3 * 4;", 24},
		{"100 / 5 / 2;", 10},

		// Grouping
		{"(2 + 3) * 4;", 20},
		{"2 * (3 + 4);", 14},
		{"(10 - 5) * (3 + 2);", 25},

		// Negative numbers
		{"-5 + 10;", 5},
		{"10 + -5;", 5},
		{"-5 * -2;", 10},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := evaluate(tt.input)
			require.IsType(t, &object.IntObject{}, result)
			assert.Equal(t, tt.expected, result.(*object.IntObject).Value)
		})
	}
}

// =============================================================================
// Infix Comparison Tests
// =============================================================================

func TestInfixComparisonEvaluation(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		// Integer equality
		{"5 == 5;", true},
		{"5 != 5;", false},
		{"5 == 10;", false},
		{"5 != 10;", true},

		// Boolean equality
		{"true == true;", true},
		{"false == false;", true},
		{"true == false;", false},
		{"true != false;", true},
		{"false != true;", true},

		// Expression equality
		{"(1 + 2) == 3;", true},
		{"(2 * 3) == 6;", true},
		{"(2 * 3) != 5;", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := evaluate(tt.input)
			require.IsType(t, &object.BoolObject{}, result)
			assert.Equal(t, tt.expected, result.(*object.BoolObject).Value)
		})
	}
}

// =============================================================================
// String Operations Tests
// =============================================================================

func TestStringConcatenation(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"'hello' + ' ' + 'world';", "hello world"},
		{"'a' + 'b';", "ab"},
		{"'' + 'a';", "a"},
		{"'a' + '';", "a"},
		// String + Int concatenation
		{"'count: ' + 5;", "count: 5"},
		{"5 + ' items';", "5 items"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := evaluate(tt.input)
			require.IsType(t, &object.StringObject{}, result)
			assert.Equal(t, tt.expected, result.(*object.StringObject).Value)
		})
	}
}

func TestStringMultiplication(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"'ab' * 3;", "ababab"},
		{"3 * 'ab';", "ababab"},
		{"'x' * 5;", "xxxxx"},
		{"'hello' * 1;", "hello"},
		{"'hello' * 0;", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := evaluate(tt.input)
			require.IsType(t, &object.StringObject{}, result)
			assert.Equal(t, tt.expected, result.(*object.StringObject).Value)
		})
	}
}

func TestStringComparison(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"'hello' == 'hello';", true},
		{"'hello' == 'world';", false},
		{"'hello' != 'world';", true},
		{"'hello' != 'hello';", false},
		{"'' == '';", true},
		{"'a' != 'b';", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := evaluate(tt.input)
			require.IsType(t, &object.BoolObject{}, result)
			assert.Equal(t, tt.expected, result.(*object.BoolObject).Value)
		})
	}
}

func TestStringOperationErrors(t *testing.T) {
	t.Run("string * string is invalid", func(t *testing.T) {
		result := evaluate("'a' * 'b';")
		assertError(t, result, "cannot perform operation 'STRING * STRING'")
	})

	t.Run("string - string is invalid", func(t *testing.T) {
		result := evaluate("'a' - 'b';")
		assertError(t, result, "cannot perform operation 'STRING - STRING'")
	})

	t.Run("string / string is invalid", func(t *testing.T) {
		result := evaluate("'a' / 'b';")
		assertError(t, result, "cannot perform operation 'STRING / STRING'")
	})
}

// =============================================================================
// Logical Operator Tests
// =============================================================================

func TestLogicalOperatorEvaluation(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		// Basic AND
		{"true && true;", true},
		{"true && false;", false},
		{"false && true;", false},
		{"false && false;", false},

		// Basic OR
		{"true || true;", true},
		{"true || false;", true},
		{"false || true;", true},
		{"false || false;", false},

		// Truthiness with integers
		{"5 && 3;", true},  // both truthy
		{"5 && 0;", false}, // 0 is falsy
		{"5 || 0;", true},  // 5 is truthy

		// Combined with NOT
		{"!false && true;", true},
		{"!true || false;", false},
		{"!(true && false);", true},
		{"!(false || false);", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := evaluate(tt.input)
			require.IsType(t, &object.BoolObject{}, result)
			assert.Equal(t, tt.expected, result.(*object.BoolObject).Value)
		})
	}
}

// =============================================================================
// If Expression Tests
// =============================================================================

func TestIfExpressionEvaluation(t *testing.T) {
	t.Run("simple if true", func(t *testing.T) {
		result := evaluate("if (true) { 10; };")
		require.IsType(t, &object.IntObject{}, result)
		assert.Equal(t, int64(10), result.(*object.IntObject).Value)
	})

	t.Run("simple if false returns null", func(t *testing.T) {
		result := evaluate("if (false) { 10; };")
		require.IsType(t, object.NullObject{}, result)
	})

	t.Run("if-else true branch", func(t *testing.T) {
		result := evaluate("if (true) { 10; } else { 20; };")
		require.IsType(t, &object.IntObject{}, result)
		assert.Equal(t, int64(10), result.(*object.IntObject).Value)
	})

	t.Run("if-else false branch", func(t *testing.T) {
		result := evaluate("if (false) { 10; } else { 20; };")
		require.IsType(t, &object.IntObject{}, result)
		assert.Equal(t, int64(20), result.(*object.IntObject).Value)
	})

	t.Run("if with equality condition", func(t *testing.T) {
		result := evaluate("if (5 == 5) { 10; };")
		require.IsType(t, &object.IntObject{}, result)
		assert.Equal(t, int64(10), result.(*object.IntObject).Value)
	})

	t.Run("if with inequality condition", func(t *testing.T) {
		result := evaluate("if (5 != 10) { 42; };")
		require.IsType(t, &object.IntObject{}, result)
		assert.Equal(t, int64(42), result.(*object.IntObject).Value)
	})

	t.Run("nested if expression", func(t *testing.T) {
		result := evaluate("if (true) { if (true) { 10; }; };")
		require.IsType(t, &object.IntObject{}, result)
		assert.Equal(t, int64(10), result.(*object.IntObject).Value)
	})

	t.Run("else-if chain first branch", func(t *testing.T) {
		result := evaluate("if (true) { 1; } else if (true) { 2; } else { 3; };")
		require.IsType(t, &object.IntObject{}, result)
		assert.Equal(t, int64(1), result.(*object.IntObject).Value)
	})

	t.Run("else-if chain second branch", func(t *testing.T) {
		result := evaluate("if (false) { 1; } else if (true) { 2; } else { 3; };")
		require.IsType(t, &object.IntObject{}, result)
		assert.Equal(t, int64(2), result.(*object.IntObject).Value)
	})

	t.Run("else-if chain else branch", func(t *testing.T) {
		result := evaluate("if (false) { 1; } else if (false) { 2; } else { 3; };")
		require.IsType(t, &object.IntObject{}, result)
		assert.Equal(t, int64(3), result.(*object.IntObject).Value)
	})

	t.Run("else-if with multiple conditions", func(t *testing.T) {
		result := evaluate(`
			if (false) { 1; }
			else if (false) { 2; }
			else if (true) { 3; }
			else { 4; }
		`)
		require.IsType(t, &object.IntObject{}, result)
		assert.Equal(t, int64(3), result.(*object.IntObject).Value)
	})
}

// =============================================================================
// Block Expression Tests
// =============================================================================

func TestBlockExpressionEvaluation(t *testing.T) {
	t.Run("block returns last value", func(t *testing.T) {
		result := evaluate("{ 1; 2; 3; };")
		require.IsType(t, &object.IntObject{}, result)
		assert.Equal(t, int64(3), result.(*object.IntObject).Value)
	})

	t.Run("block with single expression", func(t *testing.T) {
		result := evaluate("{ 42; };")
		require.IsType(t, &object.IntObject{}, result)
		assert.Equal(t, int64(42), result.(*object.IntObject).Value)
	})

	t.Run("block with computation", func(t *testing.T) {
		result := evaluate("{ 5 + 5; 10 * 2; };")
		require.IsType(t, &object.IntObject{}, result)
		assert.Equal(t, int64(20), result.(*object.IntObject).Value)
	})

	t.Run("empty block returns null", func(t *testing.T) {
		result := evaluate("{};")
		require.IsType(t, object.NullObject{}, result) // 0 is falsy, returns null
	})
}

// =============================================================================
// Truthiness Tests
// =============================================================================

func TestTruthinessEvaluation(t *testing.T) {
	truthyTests := []struct {
		input    string
		expected int64
	}{
		{"if (1) { 10; };", 10},   // 1 is truthy
		{"if (100) { 10; };", 10}, // 100 is truthy
		{"if (-1) { 10; };", 10},  // -1 is truthy (non-zero)
	}

	for _, tt := range truthyTests {
		t.Run(tt.input, func(t *testing.T) {
			result := evaluate(tt.input)
			require.IsType(t, &object.IntObject{}, result)
			assert.Equal(t, tt.expected, result.(*object.IntObject).Value)
		})
	}

	t.Run("zero is falsy", func(t *testing.T) {
		result := evaluate("if (0) { 10; };")
		require.IsType(t, object.NullObject{}, result) // 0 is falsy, returns null
	})
}

// =============================================================================
// Error Cases
// =============================================================================

// Helper to assert error result
func assertError(t *testing.T, result object.Object, expectedMessage string) {
	t.Helper()
	require.IsType(t, &object.ErrorObject{}, result, "expected ErrorObject")
	errObj := result.(*object.ErrorObject)
	assert.Contains(t, errObj.Message.Inspect(), expectedMessage)
}

func TestPrefixExpressionErrors(t *testing.T) {
	t.Run("minus on boolean", func(t *testing.T) {
		result := evaluate("-true;")
		assertError(t, result, "Unsupported operator: - for Bools")
	})

	t.Run("plus on boolean", func(t *testing.T) {
		result := evaluate("+true;")
		assertError(t, result, "Unsupported operator: + for Bools")
	})
}

func TestInfixArithmeticErrors(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "bool + int",
			input:    "true + 5;",
			expected: "cannot perform operation 'BOOL + INT'",
		},
		{
			name:     "int + bool",
			input:    "5 + true;",
			expected: "cannot perform operation 'INT + BOOL'",
		},
		{
			name:     "bool - int",
			input:    "false - 10;",
			expected: "cannot perform operation 'BOOL - INT'",
		},
		{
			name:     "bool * bool",
			input:    "true * false;",
			expected: "cannot perform operation 'BOOL * BOOL'",
		},
		{
			name:     "bool / int",
			input:    "true / 2;",
			expected: "cannot perform operation 'BOOL / INT'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evaluate(tt.input)
			assertError(t, result, tt.expected)
		})
	}
}

func TestIntegerComparisonEdgeCases(t *testing.T) {
	t.Run("int equality uses direct comparison", func(t *testing.T) {
		result := evaluate("5 == 5;")
		require.IsType(t, &object.BoolObject{}, result)
		assert.True(t, result.(*object.BoolObject).Value)
	})

	t.Run("int inequality uses direct comparison", func(t *testing.T) {
		result := evaluate("5 != 3;")
		require.IsType(t, &object.BoolObject{}, result)
		assert.True(t, result.(*object.BoolObject).Value)
	})

	t.Run("zero equals zero", func(t *testing.T) {
		result := evaluate("0 == 0;")
		require.IsType(t, &object.BoolObject{}, result)
		assert.True(t, result.(*object.BoolObject).Value)
	})

	t.Run("negative numbers comparison", func(t *testing.T) {
		result := evaluate("-5 == -5;")
		require.IsType(t, &object.BoolObject{}, result)
		assert.True(t, result.(*object.BoolObject).Value)
	})
}

func TestErrorPropagationInBlocks(t *testing.T) {
	t.Run("error stops block evaluation", func(t *testing.T) {
		result := evaluate("{ 1; true + 5; 3; };")
		assertError(t, result, "cannot perform operation 'BOOL + INT'")
	})

	t.Run("error in nested block propagates", func(t *testing.T) {
		result := evaluate("{ { true + 5; }; 10; };")
		assertError(t, result, "cannot perform operation 'BOOL + INT'")
	})

	t.Run("error in if condition block", func(t *testing.T) {
		result := evaluate("if (true) { true + 5; };")
		assertError(t, result, "cannot perform operation 'BOOL + INT'")
	})

	t.Run("error in else block", func(t *testing.T) {
		result := evaluate("if (false) { 10; } else { true - 1; };")
		assertError(t, result, "cannot perform operation 'BOOL - INT'")
	})
}

func TestErrorPropagationInProgram(t *testing.T) {
	t.Run("error stops program evaluation", func(t *testing.T) {
		result := evaluate("1; true + 5; 3;")
		assertError(t, result, "cannot perform operation 'BOOL + INT'")
	})

	t.Run("statements after error are not evaluated", func(t *testing.T) {
		// If error didn't stop, we'd get 3
		result := evaluate("true * false; 3;")
		assertError(t, result, "cannot perform operation 'BOOL * BOOL'")
	})
}

// =============================================================================
// Return Statement Tests
// =============================================================================

func TestReturnStatementEvaluation(t *testing.T) {
	t.Run("simple return", func(t *testing.T) {
		result := evaluate("return 10;")
		require.IsType(t, &object.ReturnObject{}, result)
		returnObj := result.(*object.ReturnObject)
		require.IsType(t, &object.IntObject{}, returnObj.Value)
		assert.Equal(t, int64(10), returnObj.Value.(*object.IntObject).Value)
	})

	t.Run("return expression", func(t *testing.T) {
		result := evaluate("return 5 + 5;")
		require.IsType(t, &object.ReturnObject{}, result)
		returnObj := result.(*object.ReturnObject)
		require.IsType(t, &object.IntObject{}, returnObj.Value)
		assert.Equal(t, int64(10), returnObj.Value.(*object.IntObject).Value)
	})

	t.Run("return boolean", func(t *testing.T) {
		result := evaluate("return true;")
		require.IsType(t, &object.ReturnObject{}, result)
		returnObj := result.(*object.ReturnObject)
		require.IsType(t, &object.BoolObject{}, returnObj.Value)
		assert.True(t, returnObj.Value.(*object.BoolObject).Value)
	})

	t.Run("return stops program evaluation", func(t *testing.T) {
		result := evaluate("return 10; 20;")
		require.IsType(t, &object.ReturnObject{}, result)
		returnObj := result.(*object.ReturnObject)
		require.IsType(t, &object.IntObject{}, returnObj.Value)
		assert.Equal(t, int64(10), returnObj.Value.(*object.IntObject).Value)
	})

	t.Run("return stops block evaluation", func(t *testing.T) {
		result := evaluate("{ return 10; 20; };")
		require.IsType(t, &object.ReturnObject{}, result)
		returnObj := result.(*object.ReturnObject)
		require.IsType(t, &object.IntObject{}, returnObj.Value)
		assert.Equal(t, int64(10), returnObj.Value.(*object.IntObject).Value)
	})

	t.Run("return in nested block propagates", func(t *testing.T) {
		result := evaluate("{ { return 10; }; 20; };")
		require.IsType(t, &object.ReturnObject{}, result)
		returnObj := result.(*object.ReturnObject)
		require.IsType(t, &object.IntObject{}, returnObj.Value)
		assert.Equal(t, int64(10), returnObj.Value.(*object.IntObject).Value)
	})

	t.Run("return in if block", func(t *testing.T) {
		result := evaluate("if (true) { return 10; }; 20;")
		require.IsType(t, &object.ReturnObject{}, result)
		returnObj := result.(*object.ReturnObject)
		require.IsType(t, &object.IntObject{}, returnObj.Value)
		assert.Equal(t, int64(10), returnObj.Value.(*object.IntObject).Value)
	})
}

// =============================================================================
// Error in Subexpression Tests
// =============================================================================

func TestErrorInSubexpressions(t *testing.T) {
	t.Run("error in left infix operand propagates", func(t *testing.T) {
		result := evaluate("(true + 5) + 10;")
		assertError(t, result, "cannot perform operation 'BOOL + INT'")
	})

	t.Run("error in right infix operand propagates", func(t *testing.T) {
		result := evaluate("10 + (true + 5);")
		assertError(t, result, "cannot perform operation 'BOOL + INT'")
	})

	t.Run("error in return value propagates", func(t *testing.T) {
		result := evaluate("return true + 5;")
		assertError(t, result, "cannot perform operation 'BOOL + INT'")
	})

	t.Run("error in if condition propagates", func(t *testing.T) {
		result := evaluate("if (true + 5) { 10; };")
		assertError(t, result, "cannot perform operation 'BOOL + INT'")
	})

	t.Run("error in else-if condition propagates", func(t *testing.T) {
		result := evaluate("if (false) { 1; } else if (true + 5) { 2; };")
		assertError(t, result, "cannot perform operation 'BOOL + INT'")
	})

	t.Run("error in prefix operand propagates", func(t *testing.T) {
		result := evaluate("-(true + 5);")
		assertError(t, result, "cannot perform operation 'BOOL + INT'")
	})

	t.Run("error in bang prefix operand propagates", func(t *testing.T) {
		result := evaluate("!(true + 5);")
		assertError(t, result, "cannot perform operation 'BOOL + INT'")
	})
}

// =============================================================================
// Complex Expression Tests
// =============================================================================

func TestComplexExpressions(t *testing.T) {
	t.Run("arithmetic in condition", func(t *testing.T) {
		result := evaluate("if ((5 + 5) == 10) { 100; };")
		require.IsType(t, &object.IntObject{}, result)
		assert.Equal(t, int64(100), result.(*object.IntObject).Value)
	})

	t.Run("logical in condition", func(t *testing.T) {
		result := evaluate("if (true && true) { 42; };")
		require.IsType(t, &object.IntObject{}, result)
		assert.Equal(t, int64(42), result.(*object.IntObject).Value)
	})

	t.Run("nested arithmetic", func(t *testing.T) {
		result := evaluate("((2 + 3) * (4 - 1)) / 3;")
		require.IsType(t, &object.IntObject{}, result)
		assert.Equal(t, int64(5), result.(*object.IntObject).Value) // (5 * 3) / 3 = 5
	})
}

// =============================================================================
// Let Statement and Identifier Tests
// =============================================================================

func TestLetStatementEvaluation(t *testing.T) {
	t.Run("let binds value", func(t *testing.T) {
		result := evaluate("let x = 5; x;")
		require.IsType(t, &object.IntObject{}, result)
		assert.Equal(t, int64(5), result.(*object.IntObject).Value)
	})

	t.Run("let with expression", func(t *testing.T) {
		result := evaluate("let x = 5 + 5; x;")
		require.IsType(t, &object.IntObject{}, result)
		assert.Equal(t, int64(10), result.(*object.IntObject).Value)
	})

	t.Run("let with string", func(t *testing.T) {
		result := evaluate("let name = 'hello'; name;")
		require.IsType(t, &object.StringObject{}, result)
		assert.Equal(t, "hello", result.(*object.StringObject).Value)
	})

	t.Run("let with boolean", func(t *testing.T) {
		result := evaluate("let flag = true; flag;")
		require.IsType(t, &object.BoolObject{}, result)
		assert.True(t, result.(*object.BoolObject).Value)
	})

	t.Run("multiple let statements", func(t *testing.T) {
		result := evaluate("let x = 5; let y = 10; x;")
		require.IsType(t, &object.IntObject{}, result)
		assert.Equal(t, int64(5), result.(*object.IntObject).Value)
	})

	t.Run("use identifier in expression", func(t *testing.T) {
		result := evaluate("let x = 5; let y = 10; x + y;")
		require.IsType(t, &object.IntObject{}, result)
		assert.Equal(t, int64(15), result.(*object.IntObject).Value)
	})

	t.Run("identifier referencing identifier", func(t *testing.T) {
		result := evaluate("let x = 5; let y = x; y;")
		require.IsType(t, &object.IntObject{}, result)
		assert.Equal(t, int64(5), result.(*object.IntObject).Value)
	})

	t.Run("complex expression with identifiers", func(t *testing.T) {
		result := evaluate("let a = 2; let b = 3; let c = 4; a * b + c;")
		require.IsType(t, &object.IntObject{}, result)
		assert.Equal(t, int64(10), result.(*object.IntObject).Value)
	})

	t.Run("identifier in condition", func(t *testing.T) {
		result := evaluate("let x = true; if (x) { 42; };")
		require.IsType(t, &object.IntObject{}, result)
		assert.Equal(t, int64(42), result.(*object.IntObject).Value)
	})

	t.Run("identifier comparison", func(t *testing.T) {
		result := evaluate("let x = 5; let y = 5; x == y;")
		require.IsType(t, &object.BoolObject{}, result)
		assert.True(t, result.(*object.BoolObject).Value)
	})
}

func TestIdentifierErrors(t *testing.T) {
	t.Run("undefined identifier", func(t *testing.T) {
		result := evaluate("foobar;")
		assertError(t, result, "identifier foobar not found")
	})

	t.Run("undefined identifier in expression", func(t *testing.T) {
		result := evaluate("5 + unknown;")
		assertError(t, result, "identifier unknown not found")
	})
}

// =============================================================================
// Assignment Operator Tests
// =============================================================================

func TestAssignmentOperatorEvaluation(t *testing.T) {
	t.Run("basic reassignment", func(t *testing.T) {
		result := evaluate("let x = 5; x = 10; x;")
		require.IsType(t, &object.IntObject{}, result)
		assert.Equal(t, int64(10), result.(*object.IntObject).Value)
	})

	t.Run("assignment returns value", func(t *testing.T) {
		result := evaluate("let x = 5; x = 10;")
		require.IsType(t, &object.IntObject{}, result)
		assert.Equal(t, int64(10), result.(*object.IntObject).Value)
	})

	t.Run("self-referencing assignment", func(t *testing.T) {
		result := evaluate("let x = 5; x = x + 1; x;")
		require.IsType(t, &object.IntObject{}, result)
		assert.Equal(t, int64(6), result.(*object.IntObject).Value)
	})

	t.Run("chained assignment", func(t *testing.T) {
		result := evaluate("let x = 0; let y = 0; x = y = 5; x;")
		require.IsType(t, &object.IntObject{}, result)
		assert.Equal(t, int64(5), result.(*object.IntObject).Value)
	})

	t.Run("chained assignment y value", func(t *testing.T) {
		result := evaluate("let x = 0; let y = 0; x = y = 5; y;")
		require.IsType(t, &object.IntObject{}, result)
		assert.Equal(t, int64(5), result.(*object.IntObject).Value)
	})

	t.Run("assignment with expression", func(t *testing.T) {
		result := evaluate("let x = 0; x = 3 + 4 * 2; x;")
		require.IsType(t, &object.IntObject{}, result)
		assert.Equal(t, int64(11), result.(*object.IntObject).Value)
	})

	t.Run("assignment to string", func(t *testing.T) {
		result := evaluate("let s = 'hello'; s = 'world'; s;")
		require.IsType(t, &object.StringObject{}, result)
		assert.Equal(t, "world", result.(*object.StringObject).Value)
	})
}

func TestAssignmentOperatorErrors(t *testing.T) {
	t.Run("assignment to undefined variable", func(t *testing.T) {
		result := evaluate("x = 5;")
		assertError(t, result, "x is not defined")
	})

	t.Run("assignment to literal is invalid", func(t *testing.T) {
		result := evaluate("5 = 10;")
		assertError(t, result, "cannot perform operation")
	})
}

// =============================================================================
// Comparison Operator Tests (< and >)
// =============================================================================

func TestComparisonOperatorEvaluation(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		// Less than
		{"5 < 10;", true},
		{"10 < 5;", false},
		{"5 < 5;", false},

		// Greater than
		{"10 > 5;", true},
		{"5 > 10;", false},
		{"5 > 5;", false},

		// With expressions
		{"(1 + 2) < 5;", true},
		{"(5 - 1) > 3;", true},
		{"2 * 3 < 10;", true},
		{"10 / 2 > 4;", true},

		// Edge cases
		{"0 < 1;", true},
		{"0 > -1;", true},
		{"-5 < 0;", true},
		{"-5 > -10;", true},

		// With identifiers
		{"let x = 5; let y = 10; x < y;", true},
		{"let x = 5; let y = 10; x > y;", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := evaluate(tt.input)
			require.IsType(t, &object.BoolObject{}, result)
			assert.Equal(t, tt.expected, result.(*object.BoolObject).Value)
		})
	}
}

func TestComparisonOperatorErrors(t *testing.T) {
	t.Run("string < string is invalid", func(t *testing.T) {
		result := evaluate("'a' < 'b';")
		assertError(t, result, "cannot perform operation")
	})

	t.Run("bool > bool is invalid", func(t *testing.T) {
		result := evaluate("true > false;")
		assertError(t, result, "cannot perform operation")
	})
}

// =============================================================================
// Closure and Lexical Scoping Tests
// =============================================================================

func TestClosureEvaluation(t *testing.T) {
	t.Run("closure captures outer variable", func(t *testing.T) {
		result := evaluate(`
			let x = 10;
			let addX = fn(y) { x + y; };
			addX(5);
		`)
		require.IsType(t, &object.IntObject{}, result)
		assert.Equal(t, int64(15), result.(*object.IntObject).Value)
	})

	t.Run("closure retains captured value after outer change", func(t *testing.T) {
		result := evaluate(`
			let x = 10;
			let getX = fn() { x; };
			x = 20;
			getX();
		`)
		// Closure captures the scope, so it sees the updated value
		require.IsType(t, &object.IntObject{}, result)
		assert.Equal(t, int64(20), result.(*object.IntObject).Value)
	})

	t.Run("higher-order function returning closure", func(t *testing.T) {
		result := evaluate(`
			let makeAdder = fn(x) { fn(y) { x + y; }; };
			let add5 = makeAdder(5);
			add5(10);
		`)
		require.IsType(t, &object.IntObject{}, result)
		assert.Equal(t, int64(15), result.(*object.IntObject).Value)
	})

	t.Run("multiple closures from same factory", func(t *testing.T) {
		result := evaluate(`
			let makeAdder = fn(x) { fn(y) { x + y; }; };
			let add5 = makeAdder(5);
			let add10 = makeAdder(10);
			add5(1) + add10(1);
		`)
		require.IsType(t, &object.IntObject{}, result)
		assert.Equal(t, int64(17), result.(*object.IntObject).Value) // (5+1) + (10+1) = 17
	})

	t.Run("nested closures", func(t *testing.T) {
		result := evaluate(`
			let a = 1;
			let f = fn() {
				let b = 2;
				fn() {
					let c = 3;
					fn() { a + b + c; };
				};
			};
			f()()();
		`)
		require.IsType(t, &object.IntObject{}, result)
		assert.Equal(t, int64(6), result.(*object.IntObject).Value)
	})

	t.Run("closure with parameter shadowing", func(t *testing.T) {
		result := evaluate(`
			let x = 10;
			let f = fn(x) { x * 2; };
			f(5);
		`)
		require.IsType(t, &object.IntObject{}, result)
		assert.Equal(t, int64(10), result.(*object.IntObject).Value) // 5 * 2, not 10 * 2
	})

	t.Run("counter closure", func(t *testing.T) {
		result := evaluate(`
			let makeCounter = fn() {
				let count = 0;
				fn() {
					count = count + 1;
					count;
				};
			};
			let counter = makeCounter();
			counter();
			counter();
			counter();
		`)
		require.IsType(t, &object.IntObject{}, result)
		assert.Equal(t, int64(3), result.(*object.IntObject).Value)
	})
}
