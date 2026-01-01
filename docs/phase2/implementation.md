# Phase 2 実装計画

## 概要

Phase 2 では、Phase 1 で構築した基本的なインタプリタに以下の機能を追加する：
- 配列とマップのサポート
- エラーメッセージの改善（位置情報）
- REPL の機能強化
- コード品質の向上

---

## Step 1: 配列リテラルのサポート

### 目標
`[1, 2, 3]` 形式の配列リテラルを使用可能にする。

### 実装タスク

#### 1.1 Token の追加
- `token/token.go` に必要なトークンを確認（`LBRACKET`, `RBRACKET` は既存）

#### 1.2 AST ノードの追加
```go
// ast/ast.go
type ArrayLiteral struct {
    Token    token.Token // '[' トークン
    Elements []Expression
}
```

#### 1.3 Parser の拡張
- `parseExpression` に配列リテラルのパースを追加
- `parseExpressionList` ヘルパー関数の実装

#### 1.4 Object の追加
```go
// object/object.go
type Array struct {
    Elements []Object
}
```

#### 1.5 Evaluator の拡張
- `evalArrayLiteral` の実装
- インデックスアクセス `arr[0]` の実装

#### 1.6 テストの作成
- 配列リテラルのパーステスト
- 配列の評価テスト
- インデックスアクセスのテスト

---

## Step 2: マップリテラルのサポート

### 目標
`{"key": "value"}` 形式のマップリテラルを使用可能にする。

### 実装タスク

#### 2.1 AST ノードの追加
```go
// ast/ast.go
type MapLiteral struct {
    Token token.Token // '{' トークン
    Pairs map[Expression]Expression
}
```

#### 2.2 Parser の拡張
- マップリテラルと BlockStatement の区別
- `parseMapLiteral` の実装

#### 2.3 Object の追加
```go
// object/object.go
type Map struct {
    Pairs map[HashKey]HashPair
}

type HashKey struct {
    Type  ObjectType
    Value uint64
}

type HashPair struct {
    Key   Object
    Value Object
}
```

#### 2.4 Evaluator の拡張
- `evalMapLiteral` の実装
- キーアクセス `map["key"]` の実装

#### 2.5 テストの作成
- マップリテラルのパーステスト
- マップの評価テスト
- キーアクセスのテスト

---

## Step 3: 位置情報の追加

### 目標
エラーメッセージに行番号・列番号を含める。

### 実装タスク

#### 3.1 Token に位置情報を追加
```go
// token/token.go
type Token struct {
    Type    TokenType
    Literal string
    Line    int
    Column  int
}
```

#### 3.2 Lexer の修正
- 行番号・列番号のトラッキング
- `NewToken` で位置情報を設定

#### 3.3 Parser のエラーメッセージ改善
```go
func (p *Parser) peekError(t token.TokenType) {
    msg := fmt.Sprintf("line %d, column %d: expected %s, got %s",
        p.peekToken.Line, p.peekToken.Column, t, p.peekToken.Type)
    p.errors = append(p.errors, msg)
}
```

#### 3.4 Evaluator のエラーメッセージ改善
- AST ノードから位置情報を取得
- `newError` に位置情報を含める

#### 3.5 テストの更新
- 既存テストの期待値を更新
- 位置情報の正確性テスト

---

## Step 4: REPL の機能強化

### 目標
REPL の使い勝手を向上させる。

### 実装タスク

#### 4.1 複数行入力のサポート
- 括弧のバランスチェック
- 継続プロンプト `... ` の表示

#### 4.2 readline ライブラリの導入
- `github.com/chzyer/readline` の導入
- コマンド履歴の実装
- 上下矢印キーでの履歴ナビゲーション

#### 4.3 入力補完（オプション）
- キーワード補完
- 変数名補完

#### 4.4 エラーハンドリングの改善
- `Scanner.Err()` のチェック
- より詳細なファイルエラーメッセージ

---

## Step 5: コード品質の向上

### 目標
コードの保守性と安全性を向上させる。

### 実装タスク

#### 5.1 Environment の const ガード
```go
func (e *Environment) Set(name string, val Object) Object {
    if _, isConst := e.constVars[name]; isConst {
        return &Error{Message: "cannot reassign const variable: " + name}
    }
    e.store[name] = val
    return val
}
```

#### 5.2 break/continue の内部型隠蔽
- `object.Object` インターフェースから分離
- evaluator 内部でのみ使用

#### 5.3 剰余演算子の浮動小数点対応
```go
case "%":
    return &object.Number{Value: math.Mod(leftVal, rightVal)}
```

#### 5.4 無限ループ検出（オプション）
- 実行ステップカウンター
- 最大ステップ数の設定

#### 5.5 スタックオーバーフロー対策（オプション）
- 呼び出し深度のトラッキング
- 最大深度の設定

---

## Step 6: 組み込み関数の拡張

### 目標
配列・マップ操作のための組み込み関数を追加する。

### 実装タスク

#### 6.1 配列用組み込み関数
- `len(arr)` - 配列の長さ
- `push(arr, elem)` - 要素の追加
- `first(arr)` - 最初の要素
- `last(arr)` - 最後の要素
- `rest(arr)` - 最初の要素を除いた配列

#### 6.2 マップ用組み込み関数
- `keys(map)` - キーの配列
- `values(map)` - 値の配列

#### 6.3 テストの作成
- 各組み込み関数のテスト

---

## 優先順位

| 優先度 | Step | 内容 |
|--------|------|------|
| 高 | Step 1 | 配列リテラル |
| 高 | Step 2 | マップリテラル |
| 中 | Step 3 | 位置情報 |
| 中 | Step 6 | 組み込み関数の拡張 |
| 低 | Step 4 | REPL 機能強化 |
| 低 | Step 5 | コード品質向上 |

---

## 完了条件

- [ ] 配列リテラル `[1, 2, 3]` が動作する
- [ ] マップリテラル `{"a": 1}` が動作する
- [ ] インデックスアクセス `arr[0]`, `map["key"]` が動作する
- [ ] エラーメッセージに行番号が含まれる
- [ ] `len()`, `push()` などの組み込み関数が動作する
- [ ] 全テストが成功する
