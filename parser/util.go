package parser

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"monkey/ast"
	"monkey/token"
)

var ident = 0

func trace(s string) string {
	if len(os.Getenv("TRACE")) == 0 {
		return ""
	}
	ident += 2
	fmt.Printf("%s> %s\n", strings.Repeat(" ", ident), s)
	return s
}

func untrace(s string) {
	if len(os.Getenv("TRACE")) == 0 {
		return
	}
	fmt.Printf("%s< %s\n", strings.Repeat(" ", ident), s)
	ident -= 2
}

func parseint(intAsString string) (int64, error) {
	cleanNumber := strings.ReplaceAll(intAsString, "_", "")
	i, err := strconv.ParseInt(cleanNumber, 10, 64)
	if err != nil {
		return 0, err
	}
	return i, nil
}

func parsebool(boolAsString string) bool {
	// This should never fail because lexer guarantees "true" or "false"
	i, err := strconv.ParseBool(boolAsString)
	if err != nil {
		panic(err) // Internal error - lexer bug
	}
	return i
}

// isBlockLikeExpression returns true if the expression ends with '}' and
// parseBlockExpression has already advanced past it. This affects cursor
// position when finishing statement parsing.
func isBlockLikeExpression(expr ast.Expression) bool {
	switch expr.(type) {
	case *ast.IfExpression, *ast.BlockExpression, *ast.FnExpression:
		return true
	default:
		return false
	}
}

// finishStatement handles the common logic for ending a statement:
// - Advances cursor if expression didn't end with a block
// - Checks for semicolon (required or optional based on semicolonRequired)
// - Skips any trailing semicolons
func (p *Parser) finishStatement(expr ast.Expression, semicolonRequired bool) error {
	isBlockLike := isBlockLikeExpression(expr)

	if !isBlockLike {
		p.nextToken()
	}

	if semicolonRequired || !isBlockLike {
		if p.currentToken.Type != token.SEMICOLON {
			return fmt.Errorf("expected ;, got %s", p.currentToken.Type)
		}
	}
	p.skipSemicolons()
	return nil
}
