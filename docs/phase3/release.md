# v0.0.3: Error Handling, Type Conversion & File I/O

Phase 3 の実装が完了しました。エラーハンドリング、型変換、ファイル I/O の機能が追加され、より実用的なプログラムが書けるようになりました。

## New Features

### Error Handling (try/catch/throw)
```javascript
try {
    const result = riskyOperation();
} catch (e) {
    outln("Error: " + e);
}

throw "Something went wrong";
```

### Type Conversion Functions
| Function | Description | Example |
|----------|-------------|---------|
| `int(x)` | Convert to integer | `int(3.7)` → `3` |
| `float(x)` | Convert to float | `float("3.14")` → `3.14` |
| `string(x)` | Convert to string | `string(42)` → `"42"` |
| `bool(x)` | Convert to boolean | `bool(0)` → `false` |

### File I/O Functions
| Function | Description |
|----------|-------------|
| `readFile(path)` | Read file contents |
| `writeFile(path, content)` | Write to file |
| `appendFile(path, content)` | Append to file |
| `fileExists(path)` | Check if file exists |

### Index Assignment
```javascript
mut arr = [1, 2, 3];
arr[0] = 10;  // [10, 2, 3]

mut map = {"a": 1};
map["b"] = 2;  // {"a": 1, "b": 2}
```

## Improvements

- `const` variables now properly prevent element assignment (`arr[0] = x` on const array throws error)
- `fileExists()` returns `false` for directories
- Clarified `int()` truncation behavior for negative numbers (trunc, not floor)

## Downloads

| Platform | File |
|----------|------|
| macOS (Apple Silicon) | `sugu-darwin-arm64.zip` |
| Windows (64-bit) | `sugu-windows-amd64.zip` |

## Documentation

- [Language Specification](https://n-yokomachi.github.io/sugu/)
