# Sugu Language Project

## 概要
Sugu は「すぐ実行できる」JavaScript風のインタプリタ言語です。

## 技術スタック
- 実装言語: Go
- 種類: インタプリタ

## プロジェクト構成
```
sugu/
├── token/      # トークン定義
├── lexer/      # 字句解析器
├── ast/        # 抽象構文木
├── parser/     # 構文解析器
├── object/     # オブジェクトシステム
├── evaluator/  # 評価器
├── repl/       # 対話型実行環境
└── docs/       # 仕様書
```

## 開発ルール

### コーディング規約
- Go の標準的なスタイルに従う
- `go fmt` でフォーマット
- テストは `*_test.go` に記述

### テスト
- 各パッケージにテストを書く
- `go test ./...` で全テスト実行

### コミット
- 日本語でコミットメッセージを書く
- 機能単位で小さくコミット

## Sugu 言語の構文

### 変数宣言
```javascript
mut x = 10;      // 再代入可能
const PI = 3.14; // 再代入不可
```

### 関数定義
```javascript
func add(a, b) => {
    return a + b;
}
```

### コメント
```javascript
// 単一行コメント
//-- 複数行コメント --//
```

## 参照ドキュメント
- [docs/specification.md](docs/specification.md) - 言語仕様
- [docs/phase1.md](docs/phase1.md) - Phase 1 実装範囲
