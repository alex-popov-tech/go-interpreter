package parser

import (
	"fmt"

	"monkey/ast"
	"monkey/token"
)

func (p *Parser) parseInfixExpression(left ast.Expression) (ast.Expression, error) {
	defer untrace(trace(fmt.Sprintf("parseInfixExpression, left is %s", left.String())))
	res := &ast.InfixExpression{Token: p.currentToken, Left: left, Operator: p.currentToken.Literal}

	precedence := p.currPrecedence()

	// Assignment is right-associative: x = y = 5 should parse as x = (y = 5)
	if p.currentToken.Type == token.ASSIGN {
		precedence = precedence - 1
	}

	p.nextToken()

	exrp, err := p.parseExpression(precedence)
	if err != nil {
		return nil, fmt.Errorf("could not parse infix expression: %s", err)
	}
	res.Right = exrp

	return res, nil
}

func (p *Parser) parseCallExpression(left ast.Expression) (ast.Expression, error) {
	defer untrace(
		trace(fmt.Sprintf("parseCallExpression, identifier or literal is %s", left.String())),
	)
	res := &ast.CallExpression{Token: p.currentToken, FnIdentifier: left}
	arguments := []ast.Expression{}

	// go from '(' to first argument or ')'
	p.nextToken()
	for token.RPAREN != p.currentToken.Type && token.EOF != p.currentToken.Type {
		expr, err := p.parseExpression(LOWEST)
		if err != nil {
			return nil, fmt.Errorf("could not parse call argument expression: %s", err)
		}
		arguments = append(arguments, expr)

		// if next token is ',' - go over it
		if token.COMMA == p.peekToken.Type {
			p.nextToken()
		}
		p.nextToken()
	}
	res.Arguments = arguments

	return res, nil
}
