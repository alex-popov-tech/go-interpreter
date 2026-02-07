package parser

import (
	"fmt"
	"os"
	"strconv"
	"strings"
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

func (p *Parser) finishStatement() {
	p.nextToken()
	p.skipSemicolons()
}
