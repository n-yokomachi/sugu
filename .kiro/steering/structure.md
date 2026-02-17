# プロジェクト構成

## 構成方針

レイヤードアーキテクチャ。インタプリタのパイプライン各段階（字句解析 → 構文解析 → 評価）をパッケージとして分離。

## ディレクトリパターン

### パイプライン層
**パターン**: `<stage>/` — パイプラインの各段階を独立パッケージとして配置

| パッケージ | 役割 | 入力 → 出力 |
|---|---|---|
| `token/` | トークン定義 | — |
| `lexer/` | 字句解析 | ソース文字列 → Token 列 |
| `ast/` | AST ノード定義 | — |
| `parser/` | 構文解析 | Token 列 → AST |
| `object/` | オブジェクト・環境定義 | — |
| `evaluator/` | 評価器・組み込み関数 | AST → Object |

### アプリケーション層
| パッケージ | 役割 |
|---|---|
| `repl/` | REPL とファイル実行（`repl.go` + `runner.go`）|
| `lambda/` | AWS Lambda ハンドラ・Lambda 用組み込み関数 |
| `main.go` | エントリーポイント（REPL / ファイル実行の切り替え）|

### ドキュメント
| ディレクトリ | 内容 |
|---|---|
| `docs/` | 言語仕様・フェーズ計画 |

## 命名規約

- **ファイル**: スネークケース（Go 標準）— `lexer.go`, `lexer_test.go`
- **パッケージ**: 小文字単語 — `token`, `lexer`, `parser`
- **型**: パスカルケース — `TokenType`, `NumberLiteral`, `Environment`
- **メソッド**: パスカルケース（エクスポート）— `Eval()`, `ParseProgram()`
- **テスト**: `Test` プレフィックス — `TestNextToken`, `TestEvalInfixExpression`

## インポート構成

```go
import (
    "標準ライブラリ"

    "sugu/token"    // プロジェクト内パッケージ
    "sugu/ast"
)
```

**依存方向**: `token` ← `lexer` ← `parser` → `ast`, `evaluator` → `ast` + `object`
- 循環依存なし。各層は下位層のみに依存

## コード構成の原則

- 各パッケージは単一責任を持つ（token は定義のみ、lexer は解析のみ）
- テストは同一パッケージ内に `*_test.go` として配置
- 組み込み関数は `evaluator/builtins.go` に集約（Lambda 用は `lambda/builtins.go`）
- 仕様に対応する機能がフェーズ単位で段階的に追加される

---
_パターンを記録し、ファイルツリーは網羅しない。既存パターンに従うコードはステアリング更新不要_
