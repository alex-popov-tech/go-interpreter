<img alt="go-interpreter heading image" src="https://github.com/user-attachments/assets/e69ed722-8166-48de-9b6f-9ab264963328" />

# Go Interpreter

A Monkey programming language interpreter written in Go, following the book **["Writing An Interpreter In Go"](https://interpreterbook.com/)** by Thorsten Ball.

## About

This project implements a tree-walking interpreter for the Monkey programming language. The interpreter follows a classic pipeline architecture:

```
Source Code -> Lexer -> Tokens -> Parser -> AST -> Evaluator -> Result
```

## Features

### Data Types
- **Integers** with underscore separators: `10_000`
- **Booleans**: `true`, `false`
- **Strings** with single, double, and backtick quotes + escape characters
- **Arrays**: `[1, 2, 3]`
- **Hashes**: `#{"name": "Monkey", "version": 1}`
- **Null**: `null`

### Operators
- **Arithmetic**: `+`, `-`, `*`, `/`
- **Comparison**: `==`, `!=`, `<`, `>`
- **Logical**: `&&`, `||`
- **Prefix**: `-`, `!`, `+`
- **Assignment**: `=` (right-associative, supports chaining: `x = y = 5`)

### Control Flow
- **If/else expressions**: `if (x > 5) { "big" } else { "small" }`
- **Else-if chains**: `if (x > 10) { "big" } else if (x > 5) { "medium" } else { "small" }`

### Functions
- **First-class functions**: `let add = fn(x, y) { x + y };`
- **Closures** with lexical scoping
- **Higher-order functions**: functions that accept and return functions
- **Immediate invocation**: `fn(x) { x * 2 }(5)`

### Collections
- **Array indexing**: `[1, 2, 3][0]` -> `1`
- **String indexing**: `"hello"[1]` -> `"e"` (Unicode-aware)
- **Hash access**: `#{"key": "value"}["key"]` -> `"value"` (keys: strings, ints, bools)

### String Operations
- **Concatenation**: `"hello" + " " + "world"`
- **Repetition**: `"ab" * 3` -> `"ababab"`
- **Coercion**: `"count: " + 5` -> `"count: 5"`

### Built-in Functions
| Function | Description |
|----------|-------------|
| `len(s)` | Length of a string or array |
| `first(arr)` | First element of an array |
| `last(arr)` | Last element of an array |
| `rest(arr)` | New array without the first element |
| `push(arr, val)` | New array with value appended |
| `puts(val, ...)` | Print values to stdout |
| `readFile(path)` | Read file contents as a string |
| `writeFile(path, content)` | Write a string to a file |

### Statements
- **Let statements**: `let x = 5;`
- **Return statements**: `return x + y;`
- **Expression statements**

## Book Progress

All 4 chapters of **"Writing An Interpreter In Go"** are complete:

- [x] **Chapter 1** -- Lexing
- [x] **Chapter 2** -- Parsing (Pratt parser with operator precedence)
- [x] **Chapter 3** -- Evaluation (tree-walking evaluator with scope chain)
- [x] **Chapter 4** -- Extending the Interpreter (strings, arrays, hashes, built-in functions)

## Usage

```bash
# Run the REPL
go run main.go

# Run a Monkey script file
go run main.go run script.monkey

# Run all tests
go test ./...
```

## Examples

```monkey
// Recursive fibonacci
let fibonacci = fn(n) {
    if (n < 2) {
        return n;
    }
    return fibonacci(n - 1) + fibonacci(n - 2);
};
fibonacci(10);
```

```monkey
// Closures and higher-order functions
let newCounter = fn() {
    let count = 0;
    return fn() {
        count = count + 1;
        return count;
    };
};
let counter = newCounter();
counter(); // 1
counter(); // 2
counter(); // 3
```

```monkey
// Functional patterns with arrays
let map = fn(arr, f) {
    if (len(arr) == 0) { return []; }
    return push(map(rest(arr), f), f(first(arr)));
};

let double = fn(x) { x * 2 };
map([1, 2, 3], double); // [2, 4, 6]
```

```monkey
// Hash maps
let person = #{"name": "Monkey", "age": 1, "alive": true};
puts(person["name"]); // Monkey
```

## Resources

- [Writing An Interpreter In Go](https://interpreterbook.com/) by Thorsten Ball
- [Writing A Compiler In Go](https://compilerbook.com/) (sequel)
