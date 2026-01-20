package parser

import (
	"monkey/token"
)

const (
	_ int = iota
	LOWEST
	EQUALS      // =
	LESSGREATER // < or >
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myfunc(X)
)

var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.LPAREN:   CALL,
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) currPrecedence() int {
	if p, ok := precedences[p.currentToken.Type]; ok {
		return p
	}
	return LOWEST
}

// skipCurrentStatement advances the parser to the next statement boundary after an error.
// This enables error recovery so we can report multiple errors in one parse.
func (p *Parser) skipCurrentStatement() {
	// skip until end of file or semicolon
	for p.currentToken.Type != token.EOF || p.currentToken.Type == token.SEMICOLON {
		p.nextToken()
	}
	// if semicolon is found, advance
	if p.currentToken.Type == token.SEMICOLON {
		p.nextToken()
	}
}

func (p *Parser) skipSemicolons() {
	for p.currentToken.Type == token.SEMICOLON {
		p.nextToken()
	}
}
