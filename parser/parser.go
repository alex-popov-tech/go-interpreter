package parser

import (
	"fmt"

	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
)

type (
	prefixParseFn func() (ast.Expression, error)
	infixParseFn  func(ast.Expression) (ast.Expression, error)
	Parser        struct {
		l *lexer.Lexer

		currentToken token.Token
		peekToken    token.Token

		prefixParseFns map[token.TokenType]prefixParseFn
		infixParseFns  map[token.TokenType]infixParseFn

		errors []error
	}
)

func New(l *lexer.Lexer) *Parser {
	parser := &Parser{
		l:              l,
		errors:         []error{},
		infixParseFns:  make(map[token.TokenType]infixParseFn),
		prefixParseFns: make(map[token.TokenType]prefixParseFn),
	}

	parser.prefixParseFns[token.IDENTIFIER] = parser.parseIdentifierExpression
	parser.prefixParseFns[token.INT] = parser.parseIntLiteralExpression
	parser.prefixParseFns[token.TRUE] = parser.parseBoolLiteralExpression
	parser.prefixParseFns[token.FALSE] = parser.parseBoolLiteralExpression
	parser.prefixParseFns[token.STRING] = parser.parseStringLiteralExpression
	parser.prefixParseFns[token.LPAREN] = parser.parseGroupedExpression
	parser.prefixParseFns[token.BANG] = parser.parsePrefixExpression
	parser.prefixParseFns[token.MINUS] = parser.parsePrefixExpression
	parser.prefixParseFns[token.IF] = parser.parseIfExpression
	parser.prefixParseFns[token.LBRACE] = parser.parseBlockExpression
	parser.prefixParseFns[token.FUNCTION] = parser.parseFnExpression

	parser.infixParseFns[token.PLUS] = parser.parseInfixExpression
	parser.infixParseFns[token.MINUS] = parser.parseInfixExpression
	parser.infixParseFns[token.SLASH] = parser.parseInfixExpression
	parser.infixParseFns[token.ASTERISK] = parser.parseInfixExpression
	parser.infixParseFns[token.EQ] = parser.parseInfixExpression
	parser.infixParseFns[token.NOT_EQ] = parser.parseInfixExpression
	parser.infixParseFns[token.LT] = parser.parseInfixExpression
	parser.infixParseFns[token.GT] = parser.parseInfixExpression

	parser.infixParseFns[token.LPAREN] = parser.parseCallExpression

	// establish a pointers
	parser.currentToken = l.NextToken()
	parser.peekToken = l.NextToken()

	return parser
}

func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) Errors() []error {
	return p.errors
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{Statements: []ast.Statement{}}

	for p.currentToken.Type != token.EOF {
		trace("parseStatement")
		statement, err := p.parseStatement()
		if err != nil {
			p.errors = append(p.errors, err)
			untrace("parseStatement => nil (error)")
			// Error recovery: synchronize and continue parsing
			p.skipCurrentStatement()
			continue
		}
		program.Statements = append(program.Statements, statement)
		untrace(fmt.Sprintf("parseStatement => '%s'", statement.String()))
	}

	return program
}
