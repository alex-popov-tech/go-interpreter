package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

func IdentifierToType(word string) TokenType {
	t, ok := map[string]TokenType{
		"fn":     FUNCTION,
		"let":    LET,
		"if":     IF,
		"else":   ELSE,
		"return": RETURN,
		"true":   TRUE,
		"false":  FALSE,
	}[word]
	if ok {
		return t
	}
	return IDENTIFIER
}

func New(tokenType TokenType, literal string) Token {
	return Token{Type: tokenType, Literal: literal}
}

func Empty() Token {
	return Token{Type: EOF, Literal: ""}
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// identifiers and literals
	IDENTIFIER = "IDENT"
	INT        = "INT"
	STRING     = "STRING"

	// operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	SLASH    = "/"
	ASTERISK = "*"
	EQ       = "=="
	NOT_EQ   = "!="
	AND      = "&&"
	OR       = "||"

	// delimiters
	COMMA     = ","
	SEMICOLON = ";"

	LPAREN   = "("
	RPAREN   = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	LBRKT    = "["
	RBRKT    = "]"
	LT       = "<"
	LT_OR_EQ = "<="
	GT       = ">"
	GT_OR_EQ = ">="

	HASH  = "#"
	COLON = ":"

	// keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
)
