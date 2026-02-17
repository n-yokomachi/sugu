# 技術スタック

## アーキテクチャ

クラシカルなインタプリタパイプライン: Lexer → Parser (Pratt Parser) → Tree-Walking Evaluator

## コア技術

- **言語**: Go 1.25+
- **ランタイム**: Go native binary / AWS Lambda
- **依存**: `github.com/aws/aws-lambda-go`（Lambda 用）

## 言語設計

### 型システム
- 動的型付け。すべての数値は `float64` で統一（Number 型）
- オブジェクト型: Number, String, Boolean, Null, Array, Map, Function, Builtin, Error, Return

### 構文の特徴
- 変数宣言: `mut`（可変） / `const`（不変）
- 関数定義: `func name(params) => { ... }` （アロー構文が必須）
- 制御構文: if/else, switch/case, while, for, break, continue
- エラーハンドリング: try/catch/throw
- コメント: `//` （単一行）、`//-- ... --//`（複数行）

### 組み込み関数
- 入出力: `out`, `outln`, `in`
- 型操作: `type`, `len`, `int`, `float`, `string`, `bool`
- コレクション: `push`, `pop`, `first`, `last`, `rest`, `keys`, `values`
- ファイル: `readFile`, `writeFile`, `appendFile`, `fileExists`

## 開発標準

### コード品質
- `go fmt` によるフォーマット
- 各パッケージに `*_test.go` でテストを配置

### テスト
- テーブル駆動テスト（Go 標準パターン）
- `go test ./...` で全テスト実行

### 共通コマンド
```bash
# 実行（REPL）: go run main.go
# 実行（ファイル）: go run main.go <filename>.sugu
# テスト: go test ./...
# フォーマット: go fmt ./...
# ビルド: go build -o sugu main.go
```

## 重要な技術的決定

- **数値型の統一**: 整数と浮動小数点を `float64` に統一。表示時に整数判定して小数点を省略
- **不変データ操作**: `push`, `pop`, `rest` は新しい配列を返す（元を変更しない）
- **Pratt Parser**: 演算子の優先度を柔軟に管理するため Pratt パーサーを採用
- **Environment チェーン**: 関数のクロージャを実現するため環境をチェーン構造で管理
- **エラー位置情報**: Token に行番号・列番号を保持し、エラーメッセージに含める

---
_標準とパターンを記録し、個々の依存関係は網羅しない_
