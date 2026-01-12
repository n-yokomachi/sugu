# Phase 4 実装計画

## 概要

Phase 4 では、Sugu を汎用プログラミング言語として使えるようにするための機能を実装する。
文字列操作、数学関数、イテレーション機能の強化を中心に、実用的なプログラムを書くための基盤を整える。

## スコープ

| Step | 機能 | 説明 |
|------|------|------|
| 1 | 剰余演算子の浮動小数点対応 | math.Mod を使用 |
| 2 | delete 関数 | マップからキーを削除 |
| 3 | 文字列操作関数 | split, join, trim, replace など |
| 4 | 数学関数 | abs, floor, ceil, round, random など |
| 5 | for-in ループ | 配列・マップのイテレーション構文 |

---

## Step 1: 剰余演算子の浮動小数点対応

### 目的

剰余演算子 `%` を浮動小数点数に対応させる。

### 現状の問題

```javascript
// 現在の実装（int64 に変換して計算）
5.5 % 2.0  // 1 (誤り)

// 正しい結果
5.5 % 2.0  // 1.5
```

### 実装内容

- [ ] **Evaluator/evaluator.go**: `%` 演算子を `math.Mod` で実装
- [ ] **テスト**: 浮動小数点剰余のテスト

### 実装詳細

```go
// 現在
case "%":
    return &object.Number{Value: float64(int64(leftVal) % int64(rightVal))}

// 改善後
case "%":
    return &object.Number{Value: math.Mod(leftVal, rightVal)}
```

### 仕様更新

`docs/specification.md` の演算子セクションを更新：

```markdown
| `%` | 剰余（浮動小数点対応） | `10 % 3` → `1`, `5.5 % 2.0` → `1.5` |
```

---

## Step 2: delete 関数

### 目的

マップからキーを削除する機能を追加する。

### 組み込み関数

| 関数 | 説明 | 例 |
|------|------|-----|
| `delete(map, key)` | マップからキーを削除 | `delete(map, "key")` |

### 実装内容

- [ ] **Evaluator/builtins.go**: `delete()` 関数を実装
- [ ] **テスト**: delete のテスト

### 動作仕様

```javascript
mut map = {"a": 1, "b": 2};
delete(map, "a");  // map は {"b": 2}
delete(map, "c");  // 存在しないキーは無視（エラーにならない）
```

### 戻り値

- 削除成功時: `true`
- キーが存在しない場合: `false`
- const 変数の場合: エラー

---

## Step 3: 文字列操作関数

### 目的

文字列を操作するための組み込み関数を追加する。

### 組み込み関数

| 関数 | 説明 | 例 |
|------|------|-----|
| `split(str, sep)` | 文字列を区切り文字で分割 | `split("a,b,c", ",")` → `["a", "b", "c"]` |
| `join(arr, sep)` | 配列を区切り文字で結合 | `join(["a", "b"], "-")` → `"a-b"` |
| `trim(str)` | 前後の空白を除去 | `trim("  hello  ")` → `"hello"` |
| `replace(str, old, new)` | 文字列を置換（全て） | `replace("aaa", "a", "b")` → `"bbb"` |
| `substring(str, start, end)` | 部分文字列を取得 | `substring("hello", 1, 4)` → `"ell"` |
| `indexOf(str, substr)` | 部分文字列の位置を返す | `indexOf("hello", "ll")` → `2` |
| `toUpper(str)` | 大文字に変換 | `toUpper("hello")` → `"HELLO"` |
| `toLower(str)` | 小文字に変換 | `toLower("HELLO")` → `"hello"` |

### 実装内容

- [ ] **Evaluator/builtins.go**: 文字列操作関数を実装
- [ ] **テスト**: 各関数のテスト

### 設計詳細

```go
// split - strings.Split を使用
"split": {
    Fn: func(args ...object.Object) object.Object {
        // 引数: (string, separator)
        // 戻り値: Array of strings
    },
},

// join - strings.Join を使用
"join": {
    Fn: func(args ...object.Object) object.Object {
        // 引数: (array, separator)
        // 戻り値: String
    },
},
```

### マルチバイト対応

- `substring` は rune 単位で動作（バイト単位ではない）
- `indexOf` も rune 単位でインデックスを返す

---

## Step 4: 数学関数

### 目的

数値計算のための組み込み関数を追加する。

### 組み込み関数

| 関数 | 説明 | 例 |
|------|------|-----|
| `abs(x)` | 絶対値 | `abs(-5)` → `5` |
| `floor(x)` | 切り捨て | `floor(3.7)` → `3` |
| `ceil(x)` | 切り上げ | `ceil(3.2)` → `4` |
| `round(x)` | 四捨五入 | `round(3.5)` → `4` |
| `min(a, b, ...)` | 最小値 | `min(3, 1, 2)` → `1` |
| `max(a, b, ...)` | 最大値 | `max(3, 1, 2)` → `3` |
| `random()` | 0以上1未満の乱数 | `random()` → `0.5234...` |
| `sqrt(x)` | 平方根 | `sqrt(16)` → `4` |
| `pow(x, y)` | べき乗 | `pow(2, 3)` → `8` |

### 実装内容

- [ ] **Evaluator/builtins.go**: 数学関数を実装
- [ ] **テスト**: 各関数のテスト

### 設計詳細

```go
// math パッケージを使用
import "math"
import "math/rand"

"abs": {
    Fn: func(args ...object.Object) object.Object {
        // math.Abs を使用
    },
},

"random": {
    Fn: func(args ...object.Object) object.Object {
        // rand.Float64() を使用
    },
},
```

### 注意点

- `floor` と `int()` の違い: `floor(-3.7)` → `-4`, `int(-3.7)` → `-3`
- `random()` は毎回異なる値を返す（初期シードは time.Now() を使用）

---

## Step 5: for-in ループ

### 目的

配列やマップを簡潔にイテレートできる構文を追加する。

### 構文

```javascript
// 配列のイテレーション
const arr = [1, 2, 3];
for (item in arr) {
    outln(item);  // 1, 2, 3
}

// インデックス付きイテレーション
for (i, item in arr) {
    outln(i + ": " + string(item));  // "0: 1", "1: 2", "2: 3"
}

// マップのイテレーション
const map = {"a": 1, "b": 2};
for (key in map) {
    outln(key);  // "a", "b" (順序不定)
}

// キーと値のイテレーション
for (key, value in map) {
    outln(key + ": " + string(value));
}
```

### 実装内容

- [ ] **Token**: `IN` トークンを追加
- [ ] **AST**: `ForInStatement` ノードを追加
- [ ] **Parser**: for-in 構文の解析
- [ ] **Evaluator**: for-in の評価
- [ ] **テスト**: for-in のテスト

### AST 設計

```go
type ForInStatement struct {
    Token    token.Token     // 'for' トークン
    Key      *Identifier     // イテレーション変数（インデックス/キー）
    Value    *Identifier     // イテレーション変数（値、オプション）
    Iterable Expression      // イテレート対象
    Body     *BlockStatement // ループ本体
}
```

### 制約

- for-in 変数は自動的に `const` として扱う（ループ内での再代入不可）
- マップのイテレーション順序は不定

---

## テスト方針

各 Step で以下のテストを作成：

1. 正常系テスト
2. 異常系テスト（エラーケース）
3. エッジケーステスト（空文字列、空配列、境界値など）
4. マルチバイト文字テスト（文字列関数）

---

## 仕様更新

各 Step 完了時に以下を更新：

- `docs/specification.md` - 言語仕様に追記
- `docs/backlog.md` - 完了した項目を削除

---

## 将来の拡張（Phase 5 以降の候補）

backlog.md から以下を将来の候補として記録：

- try/catch/finally の finally ブロック
- モジュールシステム（import）
- Unicode 対応（日本語変数名）
- 負のインデックス（Python スタイル `arr[-1]`）
- スライス操作（`arr[1:3]`）
