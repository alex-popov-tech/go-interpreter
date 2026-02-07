# Go Interpreter

A Monkey programming language interpreter written in Go, following the book **["Writing An Interpreter In Go"](https://interpreterbook.com/)** by Thorsten Ball.

## About

This project implements a tree-walking interpreter for the Monkey programming language. The interpreter follows a classic pipeline architecture:

```
Source Code -> Lexer -> Tokens -> Parser -> AST -> Evaluator -> Result
```

## Currently Supported Features

### Data Types
- **Integers** with underscore separators: `10_000`
- **Booleans**: `true`, `false`
- **Strings** with concatenation (`+`) and repetition (`*`)
- **Null**: `null`

### Operators
- **Arithmetic**: `+`, `-`, `*`, `/`
- **Comparison**: `==`, `!=`, `<`, `>`
- **Logical**: `&&`, `||`
- **Prefix**: `-`, `!`
- **Assignment**: `=`

### Control Flow
- **If/else expressions**: `if (condition) { ... } else { ... }`

### Functions
- **First-class functions**: `let add = fn(x, y) { x + y };`
- **Closures** with lexical scoping
- **Higher-order functions**

### Statements
- **Let statements**: `let x = 5;`
- **Return statements**: `return x + y;`
- **Expression statements**

## TODO (Chapter 4: Extending the Interpreter)

- [ ] Arrays: `[1, 2, 3]`
- [ ] Array indexing: `arr[0]`
- [ ] Hash maps: `{"key": "value"}`
- [ ] Hash indexing: `hash["key"]`
- [ ] Built-in functions:
  - [ ] `len(s)` - get length of string or array
  - [ ] `first(arr)` - get first element of array
  - [ ] `last(arr)` - get last element of array
  - [ ] `rest(arr)` - get array without first element
  - [ ] `push(arr, val)` - append value to array
  - [ ] `puts(val)` - print value to stdout

## Usage

```bash
# Run the REPL
go run main.go

# Run a Monkey script file
go run main.go run script.monkey

# Run all tests
go test ./...
```

## Example

```monkey
let fibonacci = fn(n) {
    if (n < 2) {
        return n;
    }
    return fibonacci(n - 1) + fibonacci(n - 2);
};

fibonacci(10);
```

## Resources

- [Writing An Interpreter In Go](https://interpreterbook.com/) by Thorsten Ball
- [Writing A Compiler In Go](https://compilerbook.com/) (sequel)
