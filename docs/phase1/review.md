# Phase 1 コードレビュー - 将来の改善点

## AST パッケージ

### 1. 位置情報の欠如
```go
// 現在: ASTノードに位置情報がない
type Identifier struct {
    Token token.Token
    Value string
}

// 将来的に: エラーメッセージを改善するために位置情報を追加
type Identifier struct {
    Token token.Token
    Value string
    Line  int  // 行番号
    Column int // 列番号
}
```

### 2. NumberLiteralの値が文字列型
```go
// 現在: Value string で保持
type NumberLiteral struct {
    Token token.Token
    Value string  // "42" や "3.14"
}

// 将来的に: 評価時に毎回パースする必要がある
// パーサーで型変換するか、evaluatorで適切に処理する必要がある
```

### 3. エラーハンドリングの仕組みがない
- ASTノードの構築時にエラーが発生する可能性があるが、現状では返す方法がない
- 将来的にはコンストラクタパターンやビルダーパターンの導入を検討

### 4. テストカバレッジの拡充
- 現在のテストはString()メソッドの動作確認が中心
- 以下のテストを追加すべき：
  - エッジケース（nil値、空の配列など）
  - ネストした構造のテスト
  - 複雑な式の組み合わせ

### 5. ForStatementのセミコロン表現
```go
// for文のString()で常にセミコロンが表示される
// "for (; ; ) { }" のような表示になるが、これは仕様通りか確認が必要
```

## Parser パッケージ

### 1. 変数シャドウイングの問題
```go
// parser.go:195-198
func (p *Parser) peekPrecedence() int {
    if p, ok := precedences[p.peekToken.Type]; ok {  // p がシャドウされている
        return p
    }
    return LOWEST
}
```
- レシーバーの `p` が `precedences` の値でシャドウされている
- 動作には問題ないが、可読性のため `prec` などに変更すべき

### 2. エラーメッセージに位置情報がない
```go
// parser.go:63-67
func (p *Parser) peekError(t token.TokenType) {
    msg := fmt.Sprintf("expected next token to be %s, got %s instead",
        t, p.peekToken.Type)
    // 行番号・列番号がないためデバッグが困難
}
```

### 3. for文の初期化式がセミコロンを2回消費
```go
// parser.go:414-420
if !p.peekTokenIs(token.SEMICOLON) {
    p.nextToken()
    stmt.Init = p.parseStatement()  // parseStatementがセミコロンを消費
} else {
    p.nextToken()
}
// その後 expectPeek(token.SEMICOLON) が必要ない可能性
```
- `parseStatement` がセミコロンを消費するため、for文内でのセミコロン処理に注意が必要

### 4. switch文のcaseブロックに明示的なブレース要求
```go
// 現在: case 1: { ... } の形式のみ対応
// JavaScript風: case 1: stmt; の形式は未対応
```
- Sugu言語仕様としてブレース必須なら問題なし
- 仕様との整合性確認が必要

### 5. 論理演算子のテスト不足
- `&&` と `||` は優先順位テーブルに定義されているが、専用のテストがない
- 短絡評価のテストも将来追加すべき

### 6. 代入式が未実装
```go
// 現在: x = 10; のような代入式は parseExpressionStatement でエラーになる可能性
// 将来的に代入演算子のサポートを検討
```

### 7. 配列・マップリテラルの未対応
- Phase 1 スコープ外だが、将来的にサポートが必要
- `parseExpression` の switch 文拡張で対応可能

## Lexer パッケージ

### 1. 文字列のエスケープシーケンス未対応
```go
// 現在: "hello\"world" は正しく解析できない
// 将来的に \n, \t, \" などのサポートが必要
```

### 2. Unicode非対応
```go
ch byte // 現在はASCIIのみ
// 日本語変数名などを使う場合は rune に変更が必要
```

### 3. エラー位置情報がない
- 行番号・列番号をTokenに持たせると、エラーメッセージが親切になる
- Phase 1では必須ではないが、将来的に検討

### 4. 負の数のリテラル
- `-10` は現在 `MINUS` + `NUMBER` として解析される
- 仕様上これで正しいが、明示的なテストがあると良い

### 5. テストのエッジケース追加
- 空入力
- 未終端の文字列 `"hello`
- 未終端の複数行コメント `//-- ...`
