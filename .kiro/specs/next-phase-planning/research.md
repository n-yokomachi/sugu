# Research & Design Decisions

## Summary
- **Feature**: `next-phase-planning` (Phase 4 実装)
- **Discovery Scope**: Extension（既存インタプリタの機能拡張）
- **Key Findings**:
  - 剰余演算子 `%` は既に `math.Mod` で実装済み。Req 1 は仕様更新とテスト追加のみ
  - 組み込み関数は `evaluator/builtins.go` の `builtins` マップに追加するだけで拡張可能
  - for-in ループは Token → AST → Parser → Evaluator の全層に変更が必要な唯一の要件

## Research Log

### 剰余演算子の現状
- **Context**: Req 1 が `math.Mod` への移行を要求しているが、既に実装済みか確認
- **Sources Consulted**: `evaluator/evaluator.go:409-413`
- **Findings**: `case "%"` は既に `math.Mod(leftVal, rightVal)` を使用。ゼロ除算チェックも実装済み
- **Implications**: コード変更不要。テスト追加と仕様ドキュメント更新のみ

### delete 関数の設計
- **Context**: マップからキーを削除する `delete` 関数。Go の `delete(map, key)` と同様の操作
- **Sources Consulted**: `evaluator/builtins.go`, `object/object.go`, `object/environment.go`
- **Findings**:
  - `object.Map.Pairs` は `map[HashKey]HashPair` 型。Go の `delete()` で直接キーを削除可能
  - const チェックは組み込み関数からは直接行えない（Environment にアクセスできない）
  - 現在の組み込み関数は `args ...object.Object` のみ受け取り、環境情報がない
- **Implications**: const チェックを実装するには、(a) 組み込み関数に Environment を渡す仕組みが必要、または (b) delete をインプレース変更として設計し、呼び出し元（Evaluator）で const チェックを行う。(b) が既存パターンに合致

### for-in の構文解析設計
- **Context**: `for (item in arr)` を `for (mut i = 0; ...)` と区別してパースする方法
- **Sources Consulted**: `parser/parser.go:479-521`, `token/token.go`
- **Findings**:
  - 現在の `parseForStatement` は `(` の後に `parseStatement` を呼ぶ
  - for-in を区別するには: `for` `(` の後に IDENT が来て、その次が `in` または `,` なら for-in
  - Lookahead が必要：最大 2 トークン先読み（既に `peekToken` で 1 つ先読み可能）
- **Implications**: `parseForStatement` 内で条件分岐し、`in` キーワードを検出したら `parseForInStatement` に委譲

### 文字列操作のマルチバイト対応
- **Context**: `substring`, `indexOf` が rune 単位で動作する必要がある
- **Sources Consulted**: Go 標準ライブラリ `strings`, `unicode/utf8`
- **Findings**:
  - `strings.Index` はバイト位置を返す → rune 位置への変換が必要
  - `[]rune(str)[start:end]` でスライスすれば rune 単位の substring が実現可能
- **Implications**: indexOf は内部で `[]rune` 変換して rune 単位のインデックスを計算する

## Architecture Pattern Evaluation

| Option | Description | Strengths | Risks / Limitations | Notes |
|--------|-------------|-----------|---------------------|-------|
| 組み込み関数追加 | builtins マップにエントリ追加 | 変更箇所が最小、既存パターン踏襲 | builtins.go が肥大化する | Req 2, 3, 4 に適用 |
| パイプライン全層変更 | Token→AST→Parser→Evaluator | 言語構文の正式な拡張 | 変更範囲が広い | Req 5 (for-in) に必要 |

## Design Decisions

### Decision: delete の const チェック方式
- **Context**: delete は組み込み関数だが、const 変数への操作を防ぐ必要がある
- **Alternatives Considered**:
  1. 組み込み関数に Environment を渡す仕組みを追加
  2. CallExpression の評価時に特殊処理を挿入
- **Selected Approach**: delete をインプレース変更として扱い、Map オブジェクトを直接変更する。const チェックは既存の `evalIndexAssignExpression` と同様、呼び出し元の評価器レベルでは行わず、delete 関数が Map を直接変更する設計とする（push/pop と同様の不変操作パターンは採用しない）
- **Rationale**: delete は「既存のマップからキーを除去する」操作であり、新しいマップを返すのは不自然。const チェックは将来の課題として記録
- **Trade-offs**: const マップに対する delete が防げないが、既存の配列代入（`arr[0] = 10`）と同様の制約
- **Follow-up**: Environment を組み込み関数に渡す仕組みの設計を将来検討

### Decision: for-in のパーサー設計
- **Context**: 既存の `for` 文と `for-in` 文を同じ `for` キーワードで開始する
- **Selected Approach**: `parseForStatement` 内で、`(` の後の最初の IDENT の次のトークンが `in` または `,`（2変数形式）なら for-in として分岐
- **Rationale**: 既存のパーサー構造を最小限の変更で拡張できる

## Risks & Mitigations
- **builtins.go の肥大化** — 関数カテゴリごとにコメントで区切り、将来的にはファイル分割を検討
- **for-in と for の構文曖昧性** — パーサーの先読みで確実に区別可能（`in` キーワードの存在で判定）
- **マルチバイト文字のパフォーマンス** — `[]rune` 変換のコストは短い文字列では無視可能。長大な文字列処理は将来の最適化対象
