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

// for (with ++ and +=)
for (mut i = 0; i < 5; i++) {
    outln(i);
}

// for-in (array)
for (item in [1, 2, 3]) {
    outln(item);
}

// for-in with index
for (i, item in ["a", "b", "c"]) {
    outln(string(i) + ": " + item);
}

// for-in (map)
for (key, value in {"x": 1, "y": 2}) {
    outln(key + " = " + string(value));
}

// switch
switch (x) {
    case 1: { outln("one"); }
    case 2: { outln("two"); }
    default: { outln("other"); }
}

// try-catch
try {
    throw "error";
} catch (e) {
    outln(e);
}
```

### Operators

| Category | Operators |
|----------|-----------|
| Arithmetic | `+`, `-`, `*`, `/`, `%` |
| Postfix | `++`, `--` |
| Compound Assignment | `+=`, `-=`, `*=`, `/=`, `%=` |
| Comparison | `==`, `!=`, `<`, `>`, `<=`, `>=` |
| Logical | `&&`, `\|\|`, `!` |

### Builtin Functions

**I/O**

| Function | Description |
|----------|-------------|
| `out(x)` | Print without newline |
| `outln(x)` | Print with newline |
| `in()` | Read user input |

**Type & Length**

| Function | Description |
|----------|-------------|
| `type(x)` | Get type of value |
| `len(x)` | Get length of string/array/map |

**Type Conversion**

| Function | Description |
|----------|-------------|
| `int(x)` | Convert to integer (truncate toward zero) |
| `float(x)` | Convert to float |
| `string(x)` | Convert to string |
| `bool(x)` | Convert to boolean |

**Array**

| Function | Description |
|----------|-------------|
| `push(arr, x)` | Append to array (returns new array) |
| `pop(arr)` | Remove last element (returns new array) |
| `first(arr)` | Get first element |
| `last(arr)` | Get last element |
| `rest(arr)` | Get all but first element |
| `contains(arr, x)` | Check if element exists in array/string |
| `concat(a, b, ...)` | Concatenate multiple arrays |

**Map**

| Function | Description |
|----------|-------------|
| `keys(map)` | Get map keys as array |
| `values(map)` | Get map values as array |
| `delete(map, key)` | Delete key from map |

**String**

| Function | Description |
|----------|-------------|
| `split(str, sep)` | Split string by separator |
| `join(arr, sep)` | Join array elements with separator |
| `trim(str)` | Remove leading/trailing whitespace |
| `replace(str, old, new)` | Replace all occurrences |
| `substring(str, start, end)` | Get substring (rune-based) |
| `indexOf(str, substr)` | Find substring position (rune-based) |
| `toUpper(str)` | Convert to uppercase |
| `toLower(str)` | Convert to lowercase |

**Math**

| Function | Description |
|----------|-------------|
| `abs(x)` | Absolute value |
| `floor(x)` | Round toward negative infinity |
| `ceil(x)` | Round toward positive infinity |
| `round(x)` | Round to nearest integer |
| `sqrt(x)` | Square root |
| `pow(x, y)` | Power (x^y) |
| `min(a, b, ...)` | Minimum value |
| `max(a, b, ...)` | Maximum value |
| `random()` | Random number [0, 1) |

**File I/O**

| Function | Description |
|----------|-------------|
| `readFile(path)` | Read file contents |
| `writeFile(path, content)` | Write to file |
| `appendFile(path, content)` | Append to file |
| `fileExists(path)` | Check if file exists |

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
