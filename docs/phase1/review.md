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

### 1. エラーメッセージに位置情報がない
```go
// parser.go:63-67
func (p *Parser) peekError(t token.TokenType) {
    msg := fmt.Sprintf("expected next token to be %s, got %s instead",
        t, p.peekToken.Type)
    // 行番号・列番号がないためデバッグが困難
}
```

### 2. for文の初期化式がセミコロンを2回消費
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

### 3. switch文のcaseブロックに明示的なブレース要求
```go
// 現在: case 1: { ... } の形式のみ対応
// JavaScript風: case 1: stmt; の形式は未対応
```
- Sugu言語仕様としてブレース必須なら問題なし
- 仕様との整合性確認が必要

### 4. 配列・マップリテラルの未対応
- Phase 1 スコープ外だが、将来的にサポートが必要
- `parseExpression` の switch 文拡張で対応可能

## Object パッケージ

### 1. const変数への再代入チェックがEnvironment側にない
- `SetConst`で設定した変数に`Set`で上書きできてしまう
- Evaluator側でチェックする設計だが、Environment側でもガードがあると安全

## Lexer パッケージ

### 1. Unicode非対応
```go
ch byte // 現在はASCIIのみ
// 日本語変数名などを使う場合は rune に変更が必要
```

### 2. エラー位置情報がない
- 行番号・列番号をTokenに持たせると、エラーメッセージが親切になる
- Phase 1では必須ではないが、将来的に検討


## Evaluator パッケージ

### 1. break/continue の内部型が object.Object を実装
```go
// evaluator.go:616-624
type breakValue struct{}
func (b *breakValue) Type() object.ObjectType { return "BREAK" }
func (b *breakValue) Inspect() string         { return "break" }
```
- これらは内部制御用であり、外部に公開されるべきではない
- objectパッケージに移動するか、完全に内部実装として隠蔽すべき

### 2. 無限ループの検出がない
```go
// while (true) { } や for (;;) { } で無限ループになる
// 将来的に: 実行ステップ数やタイムアウトによる保護を検討
```

### 3. 引数の数チェックが関数呼び出し時にない
```go
// evaluator.go:541-550
// 引数が多すぎる場合はエラーにならず、余分な引数は無視される
// 引数が少ない場合はNULLで埋められる
for paramIdx, param := range fn.Parameters {
    if paramIdx < len(args) {
        env.Set(param.Value, args[paramIdx])
    } else {
        env.Set(param.Value, NULL)
    }
}
```
- 仕様としてはこれで正しいが、警告を出すオプションがあると便利

### 4. スタックオーバーフロー対策がない
```go
// 再帰関数で深すぎる呼び出しでスタックオーバーフローになる可能性
// 将来的に: 呼び出し深度のトラッキングと制限を検討
```

### 5. 組み込み関数のテストが限定的
```go
// TestBuiltinFunctions では type() のみテスト
// out, outln, in のテストがない（標準入出力を使うため難しい）
// 将来的に: io.Writer/io.Reader を注入可能にしてテスト容易性を向上
```

### 6. エラーメッセージに位置情報がない
```go
// newError() は単にメッセージを返すのみ
// 将来的に: ASTノードの位置情報を含めてエラーを生成すべき
func newError(node ast.Node, format string, a ...interface{}) *object.Error
```

### 7. 剰余演算子が整数のみ対応
```go
// evaluator.go:320-323
case "%":
    return &object.Number{Value: float64(int64(leftVal) % int64(rightVal))}
```
- float64をint64にキャストしているため、小数点以下が失われる
- `math.Mod` を使用して浮動小数点剰余に対応すべきか検討

## REPL パッケージ

### 1. 複数行入力に未対応
```go
// repl.go:26-28
scanned := scanner.Scan()
line := scanner.Text()
```
- 1行ずつしか入力できないため、複数行にわたる関数定義などが入力しづらい
- 将来的に: 括弧のバランスをチェックして継続入力を促す機能を検討

### 2. コマンド履歴機能がない
- 上下矢印キーで過去の入力を呼び出せない
- 将来的に: readline ライブラリの導入を検討

### 3. 入力補完機能がない
- キーワードや変数名の補完ができない
- 将来的に: readline + カスタム補完ハンドラの実装を検討

### 4. Scanner.Err() のチェックがない
```go
// repl.go:26-28
scanned := scanner.Scan()
if !scanned {
    return  // scanner.Err() で入力エラーの原因を確認すべき
}
```

### 5. RunFile のファイル存在チェックのエラーメッセージ
```go
// runner.go:17-19
if err != nil {
    return fmt.Errorf("failed to read file: %w", err)
}
```
- ファイルが存在しない場合とパーミッションエラーの区別ができない
- 将来的に: より詳細なエラーメッセージを検討
