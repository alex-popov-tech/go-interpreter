package parser

import (
	"fmt"

	"monkey/ast"
	"monkey/token"
)

func (p *Parser) parseStatement() (ast.Statement, error) {
	switch p.currentToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseExpressionStatement() (*ast.ExpressionStatement, error) {
	defer untrace(trace("parseExpressionStatement"))
	statement := &ast.ExpressionStatement{Token: p.currentToken}

	expr, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, fmt.Errorf("could not parse expression statement: %s", err)
	}
	statement.Expression = expr

	p.finishStatement()
	return statement, nil
}

func (p *Parser) parseLetStatement() (*ast.LetStatement, error) {
	defer untrace(trace("parseLetStatement"))

	if token.LET != p.currentToken.Type {
		return nil, fmt.Errorf("expected %s, got %s", token.LET, p.currentToken.Type)
	}
	statement := &ast.LetStatement{Token: p.currentToken}

	// move to identifier
	p.nextToken()
	if token.IDENTIFIER != p.currentToken.Type {
		return nil, fmt.Errorf("expected %s, got %s", token.IDENTIFIER, p.currentToken.Type)
	}
	statement.Identifier = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}

	// move to '='
	p.nextToken()
	if token.ASSIGN != p.currentToken.Type {
		return nil, fmt.Errorf("expected %s, got %s", token.ASSIGN, p.currentToken.Type)
	}

	// move to expression
	p.nextToken()
	expr, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, fmt.Errorf("could not parse let statement: %s", err)
	}
	statement.Value = expr

	p.finishStatement()
	return statement, nil
}

func (p *Parser) parseReturnStatement() (*ast.ReturnStatement, error) {
	defer untrace(trace(fmt.Sprintf("parseReturnStatement '%s'", p.currentToken.Literal)))

	if token.RETURN != p.currentToken.Type {
		return nil, fmt.Errorf("expected %s, got %s", token.RETURN, p.currentToken.Type)
	}
	res := &ast.ReturnStatement{Token: p.currentToken}

	if token.SEMICOLON != p.peekToken.Type {
		// move to expression
		p.nextToken()
		expr, err := p.parseExpression(LOWEST)
		if err != nil {
			return nil, fmt.Errorf("could not parse return statement: %s", err)
		}
		res.Value = expr
	}
	p.finishStatement()
	return res, nil
}
