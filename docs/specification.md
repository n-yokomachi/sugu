# Sugu 言語仕様

> 「すぐ実行できる」JavaScript風のインタプリタ言語

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
| string | テンプレートリテラル | `` `Hello ${name}` `` |
| boolean | 真偽値 | `true`, `false` |
| null | 値がないことを表す | `null` |
| array | 配列（Phase 2以降） | `[1, 2, 3]` |
| object | オブジェクト（Phase 2以降） | `{ key: "value" }` |

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
| `+` | 加算 | `1 + 2` → `3` |
| `-` | 減算 | `5 - 3` → `2` |
| `*` | 乗算 | `2 * 4` → `8` |
| `/` | 除算 | `10 / 3` → `3.333...` |
| `%` | 剰余 | `10 % 3` → `1` |

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
| `||` | OR | `true || false` → `true` |
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
    case 1:
        outln("one");
        break;
    case 2:
        outln("two");
        break;
    default:
        outln("other");
}
```

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

### 注意
- 1行での省略記法は禁止（可読性のため）
- `{}` は必須

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

| 関数 | 説明 | 例 |
|---|---|---|
| `out(x)` | 出力（改行なし） | `out("Hello")` |
| `outln(x)` | 出力（改行あり） | `outln("Hello")` |
| `in()` | ユーザー入力を受け取る | `const name = in();` |
| `type(x)` | 型を文字列で返す | `type(42)` → `"number"` |
| `len(x)` | 長さを返す（Phase 2以降） | `len("abc")` → `3` |
