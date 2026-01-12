# Phase 3 実装計画

## 概要

Phase 3 では、Sugu を実用的なプログラミング言語にするための必須機能を実装する。

## スコープ

| Step | 機能 | 説明 |
|------|------|------|
| 1 | 代入式 | 配列・マップ要素への代入 |
| 2 | エラーハンドリング | try/catch によるエラー処理 |
| 3 | 型変換関数 | `int()`, `float()`, `string()` |
| 4 | ファイル I/O | ファイルの読み書き |

---

## Step 1: 代入式 ✅

### 目的

配列要素やマップ要素への代入を可能にする。

### 構文

```javascript
mut arr = [1, 2, 3];
arr[0] = 10;           // arr は [10, 2, 3]

mut map = {"a": 1};
map["b"] = 2;          // map は {"a": 1, "b": 2}
map["a"] = 100;        // map は {"a": 100, "b": 2}
```

### 実装内容

- [x] **AST**: `IndexAssignExpression` ノードを追加
- [x] **Parser**: インデックス式の左辺値としての解析
- [x] **Evaluator**: 代入の評価ロジック
- [x] **テスト**: 配列・マップ代入のテスト

### 制約

- `const` で宣言された変数への代入はエラー
- 文字列はイミュータブル（代入不可）

---

## Step 2: エラーハンドリング ✅

### 目的

実行時エラーをキャッチして処理できるようにする。

### 構文

```javascript
try {
    const result = riskyOperation();
} catch (e) {
    outln("Error: " + e);
}
```

### 実装内容

- [x] **Token**: `TRY`, `CATCH`, `THROW` トークンを追加
- [x] **AST**: `TryStatement`, `ThrowStatement` ノードを追加
- [x] **Parser**: try/catch/throw の解析
- [x] **Evaluator**: エラーの伝播とキャッチ
- [x] **テスト**: try/catch/throw のテスト

### throw 文

```javascript
throw "Something went wrong";
```

### 特徴

- `throw` された値は `catch` ブロックでキャッチできる
- 組み込み関数からのエラー（例: `10 / 0`）もキャッチ可能
- キャッチされなかった `throw` はプログラムを終了させる

---

## Step 3: 型変換関数 ✅

### 目的

異なる型間の変換を可能にする。

### 組み込み関数

| 関数 | 説明 | 例 |
|------|------|-----|
| `int(x)` | 整数に変換（小数点以下切り捨て） | `int(3.7)` → `3` |
| `float(x)` | 浮動小数点数に変換 | `float("3.14")` → `3.14` |
| `string(x)` | 文字列に変換 | `string(42)` → `"42"` |
| `bool(x)` | 真偽値に変換 | `bool(0)` → `false` |

### 実装内容

- [x] **Evaluator/builtins.go**: `int()`, `float()`, `string()`, `bool()` 関数を実装
- [x] **テスト**: 各変換関数のテスト

### 変換ルール

```javascript
// int()
int(3.7)      // 3
int("42")     // 42
int(true)     // 1
int(false)    // 0

// float()
float("3.14") // 3.14
float(42)     // 42.0

// string()
string(42)    // "42"
string(true)  // "true"
string(null)  // "null"

// bool()
bool(0)       // false
bool("")      // false
bool(null)    // false
bool([])      // false (空配列)
bool({})      // false (空マップ)
// それ以外は true
```

---

## Step 4: ファイル I/O ✅

### 目的

ファイルの読み書きを可能にする。

### 組み込み関数

| 関数 | 説明 | 例 |
|------|------|-----|
| `readFile(path)` | ファイル内容を文字列で返す | `readFile("data.txt")` |
| `writeFile(path, content)` | ファイルに書き込む | `writeFile("out.txt", "hello")` |
| `appendFile(path, content)` | ファイルに追記 | `appendFile("log.txt", "entry")` |
| `fileExists(path)` | ファイルの存在確認 | `fileExists("config.json")` |

### 実装内容

- [x] **Evaluator/builtins.go**: `readFile()`, `writeFile()`, `appendFile()`, `fileExists()` 関数を実装
- [x] **テスト**: ファイル操作のテスト

### エラーハンドリング

```javascript
try {
    const content = readFile("missing.txt");
} catch (e) {
    outln("File not found: " + e);
}
```

---

## テスト方針

各 Step で以下のテストを作成：

1. 正常系テスト
2. 異常系テスト（エラーケース）
3. エッジケーステスト

---

## 仕様更新

各 Step 完了時に以下を更新：

- `docs/specification.md` - 言語仕様に追記
- `docs/index.md` - 機能一覧に追加
