# Sugu 言語仕様

## 基本情報

| 項目 | 内容 |
|---|---|
| 名前 | Sugu |
| 種類 | インタプリタ |
| 構文 | JavaScript風 |
| 型システム | 動的型付け |
| 実装言語 | Go |

---

## データ型

| 型 | 説明 | 例 |
|---|---|---|
| number | 数値（整数・小数を統一） | `42`, `3.14`, `-10` |
| string | 文字列（ダブルクォート） | `"hello"` |
| boolean | 真偽値 | `true`, `false` |
| null | 値がないことを表す | `null` |
| array | 配列 | `[1, 2, 3]` |
| map | マップ（連想配列） | `{"key": "value"}` |
| function | 関数 | `func(x) => { return x; }` |

## 変数宣言

| キーワード | 意味 | 例 |
|---|---|---|
| `mut` | 再代入可能（mutable） | `mut x = 10;` |
| `const` | 再代入不可（constant） | `const PI = 3.14;` |

```javascript
mut count = 0;
count = count + 1;  // OK

const name = "Sugu";
name = "Other";     // エラー！
```

## 演算子

### 算術演算子

| 演算子 | 意味 | 例 |
|---|---|---|
| `+` | 加算 / 文字列結合 | `1 + 2` → `3`, `"a" + "b"` → `"ab"` |
| `-` | 減算 / 単項マイナス | `5 - 3` → `2`, `-10` |
| `*` | 乗算 | `2 * 4` → `8` |
| `/` | 除算 | `10 / 3` → `3.333...` |
| `%` | 剰余（浮動小数点対応） | `10 % 3` → `1`, `5.5 % 2.0` → `1.5` |

> 注: `/` および `%` でゼロ除算を行うとランタイムエラー（`division by zero`）になります。

### 比較演算子

| 演算子 | 意味 | 例 |
|---|---|---|
| `==` | 等しい | `1 == 1` → `true` |
| `!=` | 等しくない | `1 != 2` → `true` |
| `<` | より小さい | `1 < 2` → `true` |
| `>` | より大きい | `2 > 1` → `true` |
| `<=` | 以下 | `1 <= 1` → `true` |
| `>=` | 以上 | `2 >= 1` → `true` |

### 論理演算子

| 演算子 | 意味 | 例 |
|---|---|---|
| `&&` | AND | `true && false` → `false` |
| `\|\|` | OR | `true \|\| false` → `true` |
| `!` | NOT | `!true` → `false` |

## 制御構文

### 条件分岐

```javascript
if (x > 0) {
    outln("positive");
} else if (x < 0) {
    outln("negative");
} else {
    outln("zero");
}
```

### switch文

```javascript
switch (value) {
    case 1: {
        outln("one");
    }
    case 2: {
        outln("two");
    }
    default: {
        outln("other");
    }
}
```

> 注: 各caseにはブレース `{}` が必要です。フォールスルーはありません。

### ループ

```javascript
// while
while (x > 0) {
    x = x - 1;
}

// for
for (mut i = 0; i < 10; i = i + 1) {
    outln(i);
}
```

### for-in ループ

配列やマップを簡潔にイテレートできます：

```javascript
// 配列: 要素をイテレート
for (item in [1, 2, 3]) {
    outln(item);  // 1, 2, 3
}

// 配列: インデックスと要素
for (i, item in ["a", "b", "c"]) {
    outln(string(i) + ": " + item);  // 0: a, 1: b, 2: c
}

// マップ: キーをイテレート
for (key in {"x": 1, "y": 2}) {
    outln(key);
}

// マップ: キーと値
for (key, value in {"x": 1, "y": 2}) {
    outln(key + " = " + string(value));
}
```

> 注: イテレーション変数は `const` として扱われ、ループ本体内での再代入はできません。`break` と `continue` も使用可能です。`in` は予約語ではなく、for 文のコンテキストでのみ特別に解釈されます。

### ループ制御

| キーワード | 意味 |
|---|---|
| `break` | ループを抜ける |
| `continue` | 次のイテレーションへ |

```javascript
for (mut i = 0; i < 10; i = i + 1) {
    if (i == 5) {
        break;     // ループ終了
    }
    if (i % 2 == 0) {
        continue;  // 偶数はスキップ
    }
    outln(i);      // 1, 3 が出力される
}
```

### エラーハンドリング

`try`/`catch` 文を使って実行時エラーをキャッチできます：

```javascript
try {
    const result = riskyOperation();
} catch (e) {
    outln("Error: " + e);
}
```

`throw` 文でエラーを投げることができます：

```javascript
func divide(a, b) => {
    if (b == 0) {
        throw "Division by zero";
    }
    return a / b;
}

try {
    const result = divide(10, 0);
} catch (e) {
    outln(e);  // "Division by zero"
}
```

> 注: キャッチされなかった `throw` はプログラムを終了させます。

## 関数

### 関数定義

```javascript
func add(a, b) => {
    return a + b;
}

func greet(name) => {
    outln("Hello, " + name);
}
```

### 関数呼び出し

```javascript
const result = add(1, 2);
greet("Sugu");
```

### 無名関数

```javascript
const double = func(x) => { return x * 2; };
outln(double(5));  // 10
```

### 注意
- 1行での省略記法は禁止（可読性のため）
- `{}` は必須

## 配列

### 配列リテラル

```javascript
const arr = [1, 2, 3, 4, 5];
const mixed = [1, "two", true, null];
const nested = [[1, 2], [3, 4]];
```

### インデックスアクセス

```javascript
const arr = [10, 20, 30];
outln(arr[0]);  // 10
outln(arr[2]);  // 30
outln(arr[5]);  // null（範囲外）
```

### 要素の変更

`mut` で宣言した配列は要素を変更できます：

```javascript
mut arr = [1, 2, 3];
arr[0] = 10;      // arr は [10, 2, 3]
arr[2] = 30;      // arr は [10, 2, 30]
```

> 注: 範囲外のインデックスへの代入はエラーになります。

## マップ

### マップリテラル

```javascript
const person = {
    "name": "Alice",
    "age": 30,
    "active": true
};
```

### キーアクセス

```javascript
outln(person["name"]);  // Alice
outln(person["age"]);   // 30
outln(person["foo"]);   // null（存在しないキー）
```

### 要素の変更・追加

`mut` で宣言したマップは要素を変更・追加できます：

```javascript
mut map = {"a": 1};
map["a"] = 10;    // 既存キーの変更: {"a": 10}
map["b"] = 2;     // 新規キーの追加: {"a": 10, "b": 2}
```

### キーの型

マップのキーには以下の型が使用できます：
- string: `{"name": "value"}`
- number: `{42: "value"}`
- boolean: `{true: "value"}`

## 文字列

### 文字列リテラル

```javascript
const str = "Hello, World!";
```

### エスケープシーケンス

| シーケンス | 意味 |
|---|---|
| `\n` | 改行 |
| `\t` | タブ |
| `\\` | バックスラッシュ |
| `\"` | ダブルクォート |

### インデックスアクセス

```javascript
const str = "Hello";
outln(str[0]);  // H
outln(str[4]);  // o

// マルチバイト文字も正しく扱える
const jp = "あいう";
outln(jp[1]);   // い
```

## コメント

```javascript
// 単一行コメント

//--
これは
複数行コメント
--//

mut x = 10; //-- インラインでも使える --//
```

## 組み込み関数

### 入出力

| 関数 | 説明 | 例 |
|---|---|---|
| `out(x, ...)` | 出力（改行なし） | `out("Hello")` |
| `outln(x, ...)` | 出力（改行あり） | `outln("Hello")` |
| `in()` | ユーザー入力を受け取る | `const name = in();` |

### 型と長さ

| 関数 | 説明 | 例 |
|---|---|---|
| `type(x)` | 型を文字列で返す | `type(42)` → `"NUMBER"` |
| `len(x)` | 長さを返す（文字数） | `len("あいう")` → `3` |

### 型変換

| 関数 | 説明 | 例 |
|---|---|---|
| `int(x)` | 整数に変換（小数点以下切り捨て） | `int(3.7)` → `3` |
| `float(x)` | 浮動小数点数に変換 | `float("3.14")` → `3.14` |
| `string(x)` | 文字列に変換 | `string(42)` → `"42"` |
| `bool(x)` | 真偽値に変換 | `bool(0)` → `false` |

#### 変換ルール

**int()**
- 数値: 小数点以下切り捨て（0方向への切り捨て）
  - `int(3.7)` → `3`
  - `int(-3.7)` → `-3`（floor ではなく trunc）
- 文字列: 数値文字列を変換 (`int("42")` → `42`)
- 真偽値: `true` → `1`, `false` → `0`

**float()**
- 数値: そのまま返す
- 文字列: 数値文字列を変換 (`float("3.14")` → `3.14`)
- 真偽値: `true` → `1.0`, `false` → `0.0`

**string()**
- すべての型を文字列表現に変換

**bool()**
- `false` になる値: `0`, `""`, `null`, `[]`, `{}`
- それ以外は `true`

### 配列操作

| 関数 | 説明 | 例 |
|---|---|---|
| `push(arr, x)` | 末尾に追加した新配列を返す | `push([1,2], 3)` → `[1,2,3]` |
| `pop(arr)` | 末尾を除いた新配列を返す | `pop([1,2,3])` → `[1,2]` |
| `first(arr)` | 最初の要素を返す | `first([1,2,3])` → `1` |
| `last(arr)` | 最後の要素を返す | `last([1,2,3])` → `3` |
| `rest(arr)` | 最初を除いた新配列を返す | `rest([1,2,3])` → `[2,3]` |

### マップ操作

| 関数 | 説明 | 例 |
|---|---|---|
| `keys(map)` | キーの配列を返す | `keys({"a":1})` → `["a"]` |
| `values(map)` | 値の配列を返す | `values({"a":1})` → `[1]` |
| `delete(map, key)` | キーを削除（成功: `true`、不在: `false`） | `delete(m, "a")` → `true` |

> 注: `keys()` と `values()` の順序は保証されません。

### 文字列操作

| 関数 | 説明 | 例 |
|---|---|---|
| `split(str, sep)` | 区切り文字で分割した配列を返す | `split("a,b,c", ",")` → `["a","b","c"]` |
| `join(arr, sep)` | 配列要素を区切り文字で結合 | `join(["a","b"], ",")` → `"a,b"` |
| `trim(str)` | 前後の空白を除去 | `trim("  hi  ")` → `"hi"` |
| `replace(str, old, new)` | すべての一致箇所を置換 | `replace("aaa", "a", "b")` → `"bbb"` |
| `substring(str, start, end)` | rune 単位で部分文字列を返す | `substring("こんにちは", 0, 3)` → `"こんに"` |
| `indexOf(str, substr)` | 部分文字列の rune 位置を返す | `indexOf("hello", "ell")` → `1` |
| `toUpper(str)` | 全文字を大文字に変換 | `toUpper("hello")` → `"HELLO"` |
| `toLower(str)` | 全文字を小文字に変換 | `toLower("HELLO")` → `"hello"` |

> 注: `substring` と `indexOf` は rune 単位で動作するため、マルチバイト文字（日本語等）でも正しいインデックスを返します。`indexOf` は見つからない場合 `-1` を返します。

### 数学関数

| 関数 | 説明 | 例 |
|---|---|---|
| `abs(x)` | 絶対値 | `abs(-5)` → `5` |
| `floor(x)` | 負の無限大方向への切り捨て | `floor(-3.7)` → `-4` |
| `ceil(x)` | 正の無限大方向への切り上げ | `ceil(3.2)` → `4` |
| `round(x)` | 四捨五入 | `round(3.5)` → `4` |
| `sqrt(x)` | 平方根 | `sqrt(9)` → `3` |
| `pow(x, y)` | x の y 乗 | `pow(2, 3)` → `8` |
| `min(a, b, ...)` | 最小値（可変長引数） | `min(3, 1, 2)` → `1` |
| `max(a, b, ...)` | 最大値（可変長引数） | `max(3, 1, 2)` → `3` |
| `random()` | 0 以上 1 未満の乱数 | `random()` → `0.42...` |

> 注: `floor` は負の無限大方向への切り捨てで、`int` は 0 方向への切り捨て（trunc）です。例: `floor(-3.7)` → `-4`, `int(-3.7)` → `-3`。`sqrt` に負の数を渡すとエラーになります。

### ファイル操作

| 関数 | 説明 | 例 |
|---|---|---|
| `readFile(path)` | ファイル内容を文字列で返す | `readFile("data.txt")` |
| `writeFile(path, content)` | ファイルに書き込む | `writeFile("out.txt", "hello")` |
| `appendFile(path, content)` | ファイルに追記 | `appendFile("log.txt", "entry")` |
| `fileExists(path)` | ファイルの存在確認 | `fileExists("config.json")` |

```javascript
// ファイルの読み書き例
writeFile("hello.txt", "Hello, World!");
const content = readFile("hello.txt");
outln(content);  // "Hello, World!"

// ファイルの存在確認
if (fileExists("config.json")) {
    const config = readFile("config.json");
}

// エラーハンドリング
try {
    const data = readFile("missing.txt");
} catch (e) {
    outln("File not found: " + e);
}
```

## エラーメッセージ

エラーメッセージには行番号と列番号が含まれます：

```
line 5, column 10: expected next token to be ), got EOF instead
```

## 予約語

以下のキーワードは変数名として使用できません：

```
mut, const, func, return,
if, else, switch, case, default,
while, for, break, continue,
try, catch, throw,
true, false, null
```
