package parser

import (
	"fmt"

	"monkey/ast"
	"monkey/token"
)

func (p *Parser) parseExpression(precedence int) (ast.Expression, error) {
	defer untrace(trace(fmt.Sprintf("parseExpression = token '%s'", p.currentToken.Literal)))

	prefixParseFn, ok := p.prefixParseFns[p.currentToken.Type]
	if !ok || prefixParseFn == nil {
		return nil, fmt.Errorf(
			"no prefix parse function for '%s' found",
			p.currentToken.Type,
		)
	}

	leftExp, err := prefixParseFn()
	if err != nil {
		return nil, err
	}

	for p.peekToken.Type != token.SEMICOLON && precedence < p.peekPrecedence() {
		trace(fmt.Sprintf("parseExpression - parse infix, token %s", p.currentToken.Literal))
		infixParseFn, ok := p.infixParseFns[p.peekToken.Type]
		if !ok || infixParseFn == nil {
			return nil, fmt.Errorf(
				"no infix parse function for '%s' found",
				p.currentToken.Type,
			)
		}

		// go to operator token
		p.nextToken()

		leftExp, err = infixParseFn(leftExp)
		if err != nil {
			return nil, err
		}
		untrace(fmt.Sprintf("parseExpression - parse infix, token %s", p.currentToken.Literal))
	}

	return leftExp, nil
}
