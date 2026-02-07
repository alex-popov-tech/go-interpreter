package lexer

import (
	"monkey/token"
)

type Lexer struct {
	input           string
	currentPosition int
	currentChar     byte
	peekPosition    int
}

func New(input string) *Lexer {
	if len(input) <= 1 {
		panic("Lexer must have input")
	}
	lexer := &Lexer{input: input}
	lexer.currentPosition = 0
	lexer.currentChar = input[0]
	lexer.peekPosition = 1
	return lexer
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9' || ch == '_'
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func (l *Lexer) NextToken() token.Token {
	var t token.Token

	l.skipWhitespaces()

	switch l.currentChar {
	case 0:
		t = token.Empty()

	case '=':
		if l.peekChar() == '=' {
			first := string(l.currentChar)
			// as this token is two-character, skip first one here
			second := string(l.nextChar())
			literal := first + second
			t = token.New(token.EQ, literal)
		} else {
			t = token.New(token.ASSIGN, string(l.currentChar))
		}
	case '+':
		t = token.New(token.PLUS, string(l.currentChar))
	case '-':
		t = token.New(token.MINUS, string(l.currentChar))
	case '*':
		t = token.New(token.ASTERISK, string(l.currentChar))
	case '/':
		t = token.New(token.SLASH, string(l.currentChar))
	case '!':
		if l.peekChar() == '=' {
			first := string(l.currentChar)
			// as this token is two-character, skip first one here
			second := string(l.nextChar())
			literal := first + second
			t = token.New(token.NOT_EQ, literal)
		} else {
			t = token.New(token.BANG, string(l.currentChar))
		}
	case ',':
		t = token.New(token.COMMA, string(l.currentChar))
	case ';':
		t = token.New(token.SEMICOLON, string(l.currentChar))

	case '(':
		t = token.New(token.LPAREN, string(l.currentChar))
	case ')':
		t = token.New(token.RPAREN, string(l.currentChar))
	case '{':
		t = token.New(token.LBRACE, string(l.currentChar))
	case '}':
		t = token.New(token.RBRACE, string(l.currentChar))
	case '<':
		if l.peekChar() == '=' {
			first := string(l.currentChar)
			// as this token is two-character, skip first one here
			second := string(l.nextChar())
			literal := first + second
			t = token.New(token.LT_OR_EQ, literal)
		} else {
			t = token.New(token.LT, string(l.currentChar))
		}
	case '>':
		if l.peekChar() == '=' {
			first := string(l.currentChar)
			// as this token is two-character, skip first one here
			second := string(l.nextChar())
			literal := first + second
			t = token.New(token.GT_OR_EQ, literal)
		} else {
			t = token.New(token.GT, string(l.currentChar))
		}
	case '\'', '"', '`':
		return token.New(token.STRING, l.readString())

	case '&':
		if l.peekChar() == '&' {
			first := string(l.currentChar)
			// as this token is two-character, skip first one here
			second := string(l.nextChar())
			literal := first + second
			t = token.New(token.AND, literal)
		}
	case '|':
		if l.peekChar() == '|' {
			first := string(l.currentChar)
			// as this token is two-character, skip first one here
			second := string(l.nextChar())
			literal := first + second
			t = token.New(token.OR, literal)
		}

	default:
		if isDigit(l.currentChar) {
			number := l.readNumber()
			return token.New(token.INT, number)
		}

		if isLetter(l.currentChar) {
			identifier := l.readIdentifier()
			t := token.IdentifierToType(identifier)
			return token.New(t, identifier)
		}

		t = token.New(token.ILLEGAL, string(l.currentChar))
	}

	l.nextChar()
	return t
}

func (this *Lexer) nextChar() byte {
	if this.peekPosition >= len(this.input) {
		this.currentChar = 0
	} else {
		this.currentChar = this.input[this.peekPosition]
	}
	this.currentPosition = this.peekPosition
	this.peekPosition += 1
	return this.currentChar
}

func (this *Lexer) peekChar() byte {
	if this.peekPosition >= len(this.input) {
		return 0
	} else {
		return this.input[this.peekPosition]
	}
}

func (this *Lexer) readIdentifier() string {
	position := this.currentPosition
	for isLetter(this.currentChar) || isDigit(this.currentChar) {
		this.nextChar()
	}
	return this.input[position:this.currentPosition]
}

func (this *Lexer) readNumber() string {
	position := this.currentPosition
	for isDigit(this.currentChar) {
		this.nextChar()
	}
	return this.input[position:this.currentPosition]
}

func (this *Lexer) readString() string {
	openingQuote := this.currentChar
	acc := ""

	// go over to first string byte
	this.nextChar()

	for this.currentChar != openingQuote && this.currentChar != 0 {
		if !this.isEscapeChar() {
			acc += string(this.currentChar)
			this.nextChar()
		} else {
			this.nextChar()
			acc += string(this.currentChar)
			this.nextChar()
		}
	}
	// go over last quote
	this.nextChar()

	return acc
}

func (this *Lexer) isEscapeChar() bool {
	return this.currentChar == byte('\\')
}

func (this *Lexer) skipWhitespaces() {
	for this.currentChar == ' ' || this.currentChar == '\t' || this.currentChar == '\n' || this.currentChar == '\r' {
		this.nextChar()
	}
}
