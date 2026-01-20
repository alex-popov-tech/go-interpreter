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
	arguments := []*ast.Identifier{}

	// go from '(' to first argument or ')'
	p.nextToken()
	for token.RPAREN != p.currentToken.Type && token.EOF != p.currentToken.Type {
		if token.IDENTIFIER != p.currentToken.Type {
			return nil, fmt.Errorf("expected %s, got %s", token.IDENTIFIER, p.currentToken.Type)
		}
		arguments = append(
			arguments,
			&ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal},
		)

		// if next token is ',' - go over it
		if token.COMMA == p.peekToken.Type {
			p.nextToken()
		}
		p.nextToken()
	}
	res.Arguments = arguments

	return res, nil
}
