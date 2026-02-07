# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Monkey programming language interpreter written in Go, following the book "Writing An Interpreter In Go" by Thorsten Ball. The interpreter is currently in early development with lexer and parser components implemented.

## Common Commands

```bash
# Run all tests
go test ./...

# Run tests for a specific package
go test ./lexer
go test ./parser

# Run a single test by name
go test ./parser -run TestOperatorPrecedenceParsing

# Run the REPL (type 'q' or 'quit' to exit)
go run main.go
```

## Architecture

The interpreter follows a classic pipeline architecture:

```
Source Code → Lexer → Tokens → Parser → AST → (Evaluator - not yet implemented)
```

### Package Structure

- **token/**: Token types and definitions. `TokenType` is a string-based enum, `Token` holds type and literal. `IdentifierToType()` maps keywords to their token types.

- **lexer/**: Single-pass character-by-character scanner. Uses peek (`peekChar()`) for two-character tokens (==, !=, <=, >=). Supports underscore in numbers (`10_000`).

- **ast/**: AST node definitions. Core interfaces: `Node` (base), `Statement`, `Expression`. Each node type in its own file (e.g., `letStatement.go`, `infixExpression.go`).

- **parser/**: Pratt parser with operator precedence. Key files:
  - `parser.go`: Main struct, initialization, `ParseProgram()` entry point
  - `expressionsParsers.go`: `parseExpression()` and precedence handling
  - `prefixParsers.go` / `infixParsers.go`: Parse function registrations
  - `statementsParser.go`: Statement parsing (let, return, expression statements)
  - `util.go`: Tracing helpers for debugging (`trace()`/`untrace()`)

- **repl/**: Read-Eval-Print Loop. Currently parses input but doesn't print AST.

### Parser Design Pattern

The parser uses the Pratt parsing technique:
- `prefixParseFns`: Map of token types to prefix parse functions (identifiers, literals, unary operators)
- `infixParseFns`: Map of token types to infix parse functions (binary operators)
- Precedence levels: `LOWEST < EQUALS < LESSGREATER < SUM < PRODUCT < PREFIX < CALL`
- Expression parsing: `parseExpression(precedence int)` recursively builds AST respecting operator precedence

To add a new operator:
1. Add token constant in `token/token.go`
2. Handle in lexer's `NextToken()` switch
3. Register parse function in `parser.New()` via `registerPrefixFn()` or `registerInfixFn()`
4. Add to `precedences` map in `expressionsParsers.go`

### Monkey Language Features (Currently Supported)

- Let statements: `let x = 5;`
- Return statements: `return 5;`
- Integer literals with underscore separators: `10_000`
- Identifiers
- Prefix expressions: `-x`, `!x`
- Infix expressions: `+`, `-`, `*`, `/`, `==`, `!=`, `<`, `>`
- Keywords: `fn`, `let`, `if`, `else`, `return`, `true`, `false`

Note: `<=` and `>=` are lexed but not yet wired into the parser's precedence/infix maps.
