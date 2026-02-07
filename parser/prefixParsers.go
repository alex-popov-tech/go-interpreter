package parser

import (
	"fmt"

	"monkey/ast"
	"monkey/token"
)

func (p *Parser) parsePrefixExpression() (ast.Expression, error) {
	defer untrace(trace("parsePrefixExpression"))
	res := &ast.PrefixExpression{Token: p.currentToken, Operator: p.currentToken.Literal}

	// go to expression itself
	p.nextToken()

	right, err := p.parseExpression(PREFIX)
	if err != nil {
		return nil, fmt.Errorf("could not parse prefix expression: %s", err)
	}
	res.Value = right

	return res, nil
}

func (p *Parser) parseGroupedExpression() (ast.Expression, error) {
	defer untrace(trace("parseGroupedExpression"))
	// go to expression itself
	p.nextToken()

	expression, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, fmt.Errorf("could not parse grouped expression: %s", err)
	}
	if p.peekToken.Literal != token.RPAREN {
		return nil, fmt.Errorf("cannot find ')' after parsing '($expression'")
	}

	// go over ')'
	p.nextToken()

	return expression, nil
}

func (p *Parser) parseIdentifierExpression() (ast.Expression, error) {
	defer untrace(trace("parseIdentifier"))
	return &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}, nil
}

func (p *Parser) parseIntLiteralExpression() (ast.Expression, error) {
	defer untrace(trace(fmt.Sprintf("parseIntLiteral = '%s'", p.currentToken.Literal)))
	// it is number for sure, only possible errors are like 99999999999999999999999999999 or something like that
	value, err := parseint(p.currentToken.Literal)
	if err != nil {
		return nil, fmt.Errorf("could not parse %q as int64: %s", p.currentToken.Literal, err)
	}
	return &ast.IntLiteral{Token: p.currentToken, Value: value}, nil
}

func (p *Parser) parseBoolLiteralExpression() (ast.Expression, error) {
	defer untrace(trace(fmt.Sprintf("parseBoolLiteral = '%s'", p.currentToken.Literal)))
	return &ast.BoolLiteral{Token: p.currentToken, Value: parsebool(p.currentToken.Literal)}, nil
}

func (p *Parser) parseStringLiteralExpression() (ast.Expression, error) {
	return &ast.StringLiteral{Token: p.currentToken, Value: (p.currentToken.Literal)}, nil
}

func (p *Parser) parseIfExpression() (ast.Expression, error) {
	defer untrace(trace("parseIfExpression"))
	res := ast.IfExpression{Token: p.currentToken}

	if token.IF != p.currentToken.Type {
		return nil, fmt.Errorf("expected %s, got %s", token.IF, p.currentToken.Type)
	}

	// proceed to '('
	p.nextToken()
	if token.LPAREN != p.currentToken.Type {
		return nil, fmt.Errorf("expected %s, got %s", token.LPAREN, p.currentToken.Type)
	}

	// proceed to condition expression
	expr, err := p.parseGroupedExpression()
	if err != nil {
		return nil, fmt.Errorf("could not parse if statement condition: %s", err)
	}
	res.Condition = expr

	// proceed to '{'
	p.nextToken()
	if token.LBRACE != p.currentToken.Type {
		return nil, fmt.Errorf("expected %s, got %s", token.LBRACE, p.currentToken.Type)
	}

	// parse 'if' block body
	ifBlock, err := p.parseBlockExpression()
	if err != nil {
		return nil, fmt.Errorf("could not parse if statement body block: %s", err)
	}
	res.IfBlock = ifBlock.(*ast.BlockExpression)
	// advance over '}'
	p.nextToken()

	// parsing else-if blocks
	for p.currentToken.Type == token.ELSE && p.peekToken.Type == token.IF {
		elseIfBlock := &ast.ElseIfBlock{Token: p.currentToken}

		if token.ELSE != p.currentToken.Type {
			return nil, fmt.Errorf("expected %s, got %s", token.ELSE, p.currentToken.Type)
		}
		p.nextToken() // proceed from 'else' to 'if'
		if token.IF != p.currentToken.Type {
			return nil, fmt.Errorf("expected %s, got %s", token.IF, p.currentToken.Type)
		}
		p.nextToken() // proceed from 'if' to '('

		if token.LPAREN != p.currentToken.Type {
			return nil, fmt.Errorf("expected %s, got %s", token.LPAREN, p.currentToken.Type)
		}
		expr, err := p.parseGroupedExpression()
		if err != nil {
			return nil, fmt.Errorf("could not parse if statement else-if condition: %s", err)
		}
		elseIfBlock.Condition = expr

		// proceed to '{'
		p.nextToken()
		if token.LBRACE != p.currentToken.Type {
			return nil, fmt.Errorf("expected %s, got %s", token.LBRACE, p.currentToken.Type)
		}

		// parse 'else-if' block body
		block, err := p.parseBlockExpression()
		if err != nil {
			return nil, fmt.Errorf("could not parse if statement else-if body block: %s", err)
		}
		// advance over '}'
		p.nextToken()
		elseIfBlock.Block = block.(*ast.BlockExpression)
		res.ElseIfBlocks = append(res.ElseIfBlocks, elseIfBlock)
	}

	// if there is no 'else' block - return
	if token.ELSE != p.currentToken.Type {
		return &res, nil
	}

	// proceed to '{' of 'else' block
	p.nextToken()
	if token.LBRACE != p.currentToken.Type {
		return nil, fmt.Errorf("expected %s, got %s", token.LBRACE, p.currentToken.Type)
	}

	elseBlock, err := p.parseBlockExpression()
	if err != nil {
		return nil, fmt.Errorf("could not parse if statement else body block: %s", err)
	}
	// advance over '}'
	p.nextToken()
	res.ElseBlock = elseBlock.(*ast.BlockExpression)

	return &res, nil
}

func (p *Parser) parseBlockExpression() (ast.Expression, error) {
	defer untrace(trace("parseBlockExpression"))
	statements := make([]ast.Statement, 0)
	res := ast.BlockExpression{Token: p.currentToken}

	// go over '{' to first statement
	if token.LBRACE != p.currentToken.Type {
		return nil, fmt.Errorf("expected %s, got %s", token.LBRACE, p.currentToken.Type)
	}
	p.nextToken()

	for token.RBRACE != p.currentToken.Type && token.EOF != p.currentToken.Type {
		statement, err := p.parseStatement()
		if err != nil {
			p.errors = append(p.errors, err)
			// Error recovery: synchronize and continue parsing
			p.skipCurrentStatement()
			continue
		}
		statements = append(statements, statement)
	}
	res.Statements = statements

	if token.RBRACE != p.currentToken.Type {
		return nil, fmt.Errorf("block expression is missing closing '}'")
	}

	return &res, nil
}

func (p *Parser) parseFnExpression() (ast.Expression, error) {
	defer untrace(trace("parseFnExpression"))
	res := &ast.FnExpression{Token: p.currentToken}
	arguments := []*ast.Identifier{}

	// go over 'fn' to '('
	p.nextToken()
	if token.LPAREN != p.currentToken.Type {
		return nil, fmt.Errorf("expected %s, got %s", token.LPAREN, p.currentToken.Type)
	}

	for token.RPAREN != p.currentToken.Type && token.EOF != p.currentToken.Type {
		p.nextToken()
		if token.RPAREN == p.currentToken.Type {
			break
		}
		if token.IDENTIFIER != p.currentToken.Type {
			return nil, fmt.Errorf("expected %s, got %s", token.IDENTIFIER, p.currentToken.Type)
		}
		arguments = append(
			arguments,
			&ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal},
		)

		p.nextToken()
	}
	res.Arguments = arguments

	// it can be EOF or ')'
	if token.EOF == p.currentToken.Type {
		return nil, fmt.Errorf("fn expression is missing closing ')'")
	}

	// its ')', go over it to '{'
	p.nextToken()
	if token.LBRACE != p.currentToken.Type {
		return nil, fmt.Errorf("expected %s, got %s", token.LBRACE, p.currentToken.Type)
	}
	body, err := p.parseBlockExpression()
	if err != nil {
		return nil, fmt.Errorf("could not parse fn statement body block: %s", err)
	}
	res.Body, _ = (body).(*ast.BlockExpression)

	return res, nil
}
