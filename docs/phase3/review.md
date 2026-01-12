# Phase 3 レビュー

## PR #17: Phase 3: Error handling, type conversion, and file I/O

### 課題点

#### 1. try/catch に finally ブロックがない

**ファイル**: [ast/ast.go](ast/ast.go#L480-L499), [parser/parser.go](parser/parser.go#L737-L789)

多くの言語では `try/catch/finally` 構文をサポートしている。リソース解放などのクリーンアップ処理に `finally` は有用。現時点では必須ではないが、将来的な拡張として検討の余地がある。

#### 2. ファイルパスのサニタイズがない

**ファイル**: [evaluator/builtins.go](evaluator/builtins.go#L281-L341)

ファイル操作関数にパストラバーサル攻撃への対策がない。`../` などを使った悪意あるパスが渡された場合、意図しない場所のファイルにアクセスできてしまう可能性がある。

```javascript
readFile("../../../etc/passwd")  // 危険
```

**改善案**: 許可されたディレクトリ外へのアクセスを制限する、または相対パスを正規化してチェックする。

#### 3. throwValue が object.Object インターフェースを実装している

**ファイル**: [evaluator/evaluator.go](evaluator/evaluator.go#L803-L808)

`throwValue` は内部的な制御フロー用の型だが、`object.Object` インターフェースを実装している。これにより誤って外部に漏れる可能性がある。`breakValue` や `continueValue` と同様の問題。
