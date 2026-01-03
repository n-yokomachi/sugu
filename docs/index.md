# Sugu Language

## Overview

Sugu は学習目的で開発されたシンプルなインタプリタ言語です。
JavaScript に似た構文を持ち、Go で実装されています。

## Quick Start

```bash
# ビルド
go build -o sugu .

# REPL を起動
./sugu

# ファイルを実行
./sugu examples/hello.sugu
```

## Hello World

```javascript
outln("Hello, Sugu!");
```

## Features

### Variables

```javascript
mut x = 10;       // 再代入可能
const PI = 3.14;  // 再代入不可

x = 20;           // OK
PI = 3;           // Error!
```

### Data Types

| Type | Example |
|------|---------|
| number | `42`, `3.14` |
| string | `"hello"` |
| boolean | `true`, `false` |
| null | `null` |
| array | `[1, 2, 3]` |
| map | `{"key": "value"}` |

### Functions

```javascript
func add(a, b) => {
    return a + b;
}

outln(add(1, 2));  // 3
```

### Control Flow

```javascript
// if-else
if (x > 0) {
    outln("positive");
} else {
    outln("non-positive");
}

// while
mut i = 0;
while (i < 5) {
    outln(i);
    i = i + 1;
}

// for
for (mut i = 0; i < 5; i = i + 1) {
    outln(i);
}

// switch
switch (x) {
    case 1: { outln("one"); }
    case 2: { outln("two"); }
    default: { outln("other"); }
}
```

### Operators

| Category | Operators |
|----------|-----------|
| Arithmetic | `+`, `-`, `*`, `/`, `%` |
| Comparison | `==`, `!=`, `<`, `>`, `<=`, `>=` |
| Logical | `&&`, `\|\|`, `!` |

### Builtin Functions

| Function | Description |
|----------|-------------|
| `out(x)` | Print without newline |
| `outln(x)` | Print with newline |
| `in()` | Read user input |
| `len(x)` | Get length of string/array/map |
| `type(x)` | Get type of value |
| `push(arr, x)` | Append to array (returns new array) |
| `pop(arr)` | Remove last element (returns new array) |
| `first(arr)` | Get first element |
| `last(arr)` | Get last element |
| `rest(arr)` | Get all but first element |
| `keys(map)` | Get map keys as array |
| `values(map)` | Get map values as array |

### Comments

```javascript
// Single line comment

//--
Multi-line
comment
--//
```

## AWS Lambda

Sugu は AWS Lambda 上で実行できます。
詳細は [AWS Lambda での Sugu 実行](lambda/usage.md) を参照。

---

## Documentation

- [Language Specification](specification.md) - 言語仕様
- [Backlog](backlog.md) - 将来の改善予定
- [AWS Lambda Usage](lambda/usage.md) - AWS Lambda での実行方法

## Architecture

```
sugu/
├── token/      # Token definitions
├── lexer/      # Lexical analysis
├── ast/        # Abstract Syntax Tree
├── parser/     # Parsing
├── object/     # Object system
├── evaluator/  # Evaluation
├── repl/       # REPL
└── lambda/     # AWS Lambda runtime
```

## License

MIT
