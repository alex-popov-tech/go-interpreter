# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Monkey programming language interpreter written in Go, following the book "Writing An Interpreter In Go" by Thorsten Ball. All 4 chapters of the book are complete — the interpreter has a working lexer, parser, evaluator, and extended features (strings, arrays, hashes, built-in functions).

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
Source Code → Lexer → Tokens → Parser → AST → Evaluator → Result
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

- **evaluator/**: Tree-walking evaluator. Key files:
  - `eval.go`: Main `Eval()` function, handles all AST node types
  - `infixExpression.go`: Binary operator evaluation (arithmetic, comparison, logical, assignment)
  - `prefixExpression.go`: Unary operator evaluation (`!`, `-`, `+`)
  - `ifExpression.go`: If/else-if/else evaluation with truthiness
  - `indexExpression.go`: Array, string, and hash indexing
  - `builtin.go`: Built-in functions (`len`, `first`, `last`, `rest`, `push`, `puts`, `readFile`, `writeFile`)

- **object/**: Object system for evaluated values. Types: `IntObject`, `BoolObject`, `NullObject`, `StringObject`, `ArrayObject`, `HashObject`, `FnObject`, `BuiltinFnObject`, `ReturnObject`, `ErrorObject`. Includes `Scope` with parent chain for lexical scoping and closures.

- **cmd/**: CLI and REPL. Supports `monkey repl`, `monkey run <file>`, and `monkey version`.

### Parser Design Pattern

The parser uses the Pratt parsing technique:
- `prefixParseFns`: Map of token types to prefix parse functions (identifiers, literals, unary operators)
- `infixParseFns`: Map of token types to infix parse functions (binary operators)
- Precedence levels: `LOWEST < ASSIGN < EQUALS < LESSGREATER < SUM < PRODUCT < BOOL < PREFIX < CALL < INDEX`
- Expression parsing: `parseExpression(precedence int)` recursively builds AST respecting operator precedence

To add a new operator:
1. Add token constant in `token/token.go`
2. Handle in lexer's `NextToken()` switch
3. Register parse function in `parser.New()` via `registerPrefixFn()` or `registerInfixFn()`
4. Add to `precedences` map in `parser/internal.go`

### Monkey Language Features (Currently Supported)

- Let statements: `let x = 5;`
- Return statements: `return 5;`
- Integer literals with underscore separators: `10_000`
- String literals with single, double, and backtick quotes
- Booleans: `true`, `false`
- Arrays: `[1, 2, 3]` with indexing `arr[0]`
- Hashes: `#{"key": "value"}` with indexing `hash["key"]`
- First-class functions with closures: `let add = fn(x, y) { x + y };`
- If/else-if/else expressions
- Prefix expressions: `-x`, `!x`, `+x`
- Infix expressions: `+`, `-`, `*`, `/`, `==`, `!=`, `<`, `>`, `&&`, `||`
- Assignment: `x = 5` (right-associative, variable must exist)
- String operations: concatenation (`+`), repetition (`*`), coercion with ints
- Built-in functions: `len`, `first`, `last`, `rest`, `push`, `puts`, `readFile`, `writeFile`

Note: `<=` and `>=` are lexed but not yet wired into the parser's precedence/infix maps.
