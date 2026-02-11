package lexer

import (
	"testing"

	"monkey/token"
)

type expectedToken struct {
	typ token.TokenType
	lit string
}

func verifyTokens(t *testing.T, input string, expected []expectedToken) {
	t.Helper()
	l := New(input)
	for i, tt := range expected {
		tok := l.NextToken()
		if tok.Type != tt.typ {
			t.Errorf("token[%d]: wrong type. expected=%q, got=%q",
				i, tt.typ, tok.Type)
		}
		if tok.Literal != tt.lit {
			t.Errorf("token[%d]: wrong literal. expected=%q, got=%q",
				i, tt.lit, tok.Literal)
		}
	}
}

func TestNextToken_LetStatements(t *testing.T) {
	input := `let five = 5;
	let ten = 10;`

	expected := []expectedToken{
		{token.LET, "let"},
		{token.IDENTIFIER, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENTIFIER, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	verifyTokens(t, input, expected)
}

func TestNextToken_FunctionDefinition(t *testing.T) {
	input := `let add = fn(x, y) {
		x + y;
	};`

	expected := []expectedToken{
		{token.LET, "let"},
		{token.IDENTIFIER, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},
		{token.LPAREN, "("},
		{token.IDENTIFIER, "x"},
		{token.COMMA, ","},
		{token.IDENTIFIER, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENTIFIER, "x"},
		{token.PLUS, "+"},
		{token.IDENTIFIER, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	verifyTokens(t, input, expected)
}

func TestNextToken_FunctionCall(t *testing.T) {
	input := `let result = add(five, ten);`

	expected := []expectedToken{
		{token.LET, "let"},
		{token.IDENTIFIER, "result"},
		{token.ASSIGN, "="},
		{token.IDENTIFIER, "add"},
		{token.LPAREN, "("},
		{token.IDENTIFIER, "five"},
		{token.COMMA, ","},
		{token.IDENTIFIER, "ten"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	verifyTokens(t, input, expected)
}

func TestNextToken_IfElseStatement(t *testing.T) {
	input := `if (9 < 10) {
	  return true;
	} else {
		return false;
	}`

	expected := []expectedToken{
		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.INT, "9"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.TRUE, "true"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.FALSE, "false"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.EOF, ""},
	}

	verifyTokens(t, input, expected)
}

func TestNextToken_ComparisonOperators(t *testing.T) {
	input := `10 == 10
	9 != 11

	8 >= 12
	6 <= 1`

	expected := []expectedToken{
		{token.INT, "10"},
		{token.EQ, "=="},
		{token.INT, "10"},
		{token.INT, "9"},
		{token.NOT_EQ, "!="},
		{token.INT, "11"},
		{token.INT, "8"},
		{token.GT_OR_EQ, ">="},
		{token.INT, "12"},
		{token.INT, "6"},
		{token.LT_OR_EQ, "<="},
		{token.INT, "1"},
		{token.EOF, ""},
	}

	verifyTokens(t, input, expected)
}

func TestNextToken_String(t *testing.T) {
	input := `
	'hello mom'
	'hello "mom"'
	'hello \'mom\''
`
	expected := []expectedToken{
		{token.STRING, "hello mom"},
		{token.STRING, "hello \"mom\""},
		{token.STRING, "hello 'mom'"},
	}

	verifyTokens(t, input, expected)
}

func TestNextToken_UnterminatedString(t *testing.T) {
	input := "'hello mom"

	expected := []expectedToken{
		{token.STRING, "hello mom"},
		{token.EOF, ""},
	}

	verifyTokens(t, input, expected)
}

func TestNextToken_LogicalOperators(t *testing.T) {
	input := `true && false
	true || false
	a && b || c`

	expected := []expectedToken{
		{token.TRUE, "true"},
		{token.AND, "&&"},
		{token.FALSE, "false"},
		{token.TRUE, "true"},
		{token.OR, "||"},
		{token.FALSE, "false"},
		{token.IDENTIFIER, "a"},
		{token.AND, "&&"},
		{token.IDENTIFIER, "b"},
		{token.OR, "||"},
		{token.IDENTIFIER, "c"},
		{token.EOF, ""},
	}

	verifyTokens(t, input, expected)
}

func TestNextToken_Brackets(t *testing.T) {
	input := `[0]
	["hello"]
	`

	expected := []expectedToken{
		{token.LBRKT, "["},
		{token.INT, "0"},
		{token.RBRKT, "]"},
		{token.LBRKT, "["},
		{token.STRING, "hello"},
		{token.RBRKT, "]"},
		{token.EOF, ""},
	}

	verifyTokens(t, input, expected)
}

func TestNextToken_Hash(t *testing.T) {
	input := "#{ foo: bar }"

	expected := []expectedToken{
		{token.HASH, "#"},
		{token.LBRACE, "{"},
		{token.IDENTIFIER, "foo"},
		{token.COLON, ":"},
		{token.IDENTIFIER, "bar"},
		{token.RBRACE, "}"},
		{token.EOF, ""},
	}

	verifyTokens(t, input, expected)
}
