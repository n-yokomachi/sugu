# Research & Design Decisions

## Summary
- **Feature**: phase5-syntax-sugar-and-array-utils
- **Discovery Scope**: Extension
- **Key Findings**:
  - 既存の Lexer は `+` の次の文字を peekChar で確認する仕組みが既にあり、`++` / `+=` のトークン追加は低リスク
  - 後置インクリメントは Pratt Parser の中置演算子パースとして優先順位テーブルに追加可能
  - スライス式は既存の `parseIndexExpression` を拡張し、COLON トークン検出で分岐する設計が最適

## Research Log

### 後置 vs 前置インクリメント/デクリメントの設計選択
- **Context**: `++` / `--` を前置・後置どちらで（または両方）サポートするか
- **Sources Consulted**: JavaScript / Go の仕様、既存の Sugu パーサー構造
- **Findings**:
  - JavaScript は前置・後置両方をサポートするが、Go は後置のみ（文 (statement) として扱い、式の値は返さない）
  - Sugu の Pratt Parser では後置演算子を infix 関数として登録可能（左辺の式を受け取り、ポスト処理する）
  - 前置インクリメントは既存の PrefixExpression の拡張で可能だが、副作用のある前置式はセマンティクスが複雑になる
- **Implications**: 後置のみをサポートし、JavaScript の後置セマンティクス（操作前の値を返す）を採用する。前置は Phase 5 のスコープ外とする

### 複合代入演算子のパーサー統合
- **Context**: `+=` 等を既存の代入処理にどう統合するか
- **Sources Consulted**: 既存の `parseExpression` 内の代入処理（parser.go L254-265）
- **Findings**:
  - 現在の実装では `token.ASSIGN` を検出した後に `AssignExpression` / `IndexAssignExpression` を分岐生成
  - 複合代入は同じ分岐に新トークン（`PLUS_ASSIGN` 等）を追加し、`CompoundAssignExpression` を生成する方式が最もシンプル
  - 既存の `AssignExpression` に `Operator` フィールドを追加する案もあるが、ゼロ値（空文字列）が通常代入を意味する暗黙ルールが生まれるため、別ノードが望ましい
- **Implications**: `CompoundAssignExpression` を新規 AST ノードとして追加し、`Operator`, `Name`, `Value` フィールドを持たせる

### スライス式のパース戦略
- **Context**: `arr[1:3]` を既存の `parseIndexExpression` にどう統合するか
- **Sources Consulted**: 既存の `parseIndexExpression`（parser.go L804-815）、Python / Go のスライス構文
- **Findings**:
  - 現在は `[` の後で `parseExpression(LOWEST)` を呼び、`]` を期待する
  - COLON トークンは既存のトークンセットに含まれている（マップリテラルで使用）
  - `parseExpression` が式を返した後に peekToken が COLON なら SliceExpression に切り替える設計が可能
  - `[:end]` の場合は `[` の直後に COLON が来るので、最初の式がない場合も処理する必要がある
- **Implications**: `parseIndexExpression` 内で COLON を検出したら `SliceExpression` ノードを生成。Low/High は nil 許容（省略時）

### 負インデックスの既存挙動との整合性
- **Context**: 現在 `arr[-1]` は `null` を返す。これを末尾アクセスに変更する際の影響
- **Sources Consulted**: evalArrayIndexExpression、evalStringIndexExpression、evalArrayIndexAssignment
- **Findings**:
  - 読み取り: `idx < 0 || idx > max` → NULL（3箇所：配列、文字列、マップ）
  - 代入: `idx < 0 || idx >= length` → エラー
  - 負インデックスを `length + idx` に変換するロジックを各所に追加する
  - 変換後もまだ負の場合は、従来通り NULL（読み取り）/ エラー（代入）を返す
- **Implications**: 既存テストで `arr[-1]` → `null` を期待しているものがあれば修正が必要。破壊的変更だが、既存テストを確認したところ負インデックスの明示的テストは存在しない

## Design Decisions

### Decision: 後置インクリメントのみサポート
- **Context**: `++` / `--` の前置・後置サポート範囲
- **Alternatives Considered**:
  1. 前置・後置両方 — JavaScript 互換だが複雑
  2. 後置のみ — Go スタイルでシンプル
  3. 後置のみだが値を返さない — Go と同じく文としてのみ使用
- **Selected Approach**: 後置のみ、値を返す（JavaScript スタイル）
- **Rationale**: Sugu は JavaScript 風の構文を目指しており、`for` ループでの `i++` が主要ユースケース。値を返すことで `arr[i++]` のような式も可能になる
- **Trade-offs**: 前置が使えないが、`--x` は `x -= 1` で代替可能
- **Follow-up**: for ループの update 式での動作テストを重点的に行う

### Decision: 別 AST ノード vs 既存ノード拡張
- **Context**: 複合代入をどの AST ノードで表現するか
- **Alternatives Considered**:
  1. `AssignExpression` に `Operator` フィールド追加 — 最小変更
  2. `CompoundAssignExpression` 新規追加 — 責任分離が明確
- **Selected Approach**: `CompoundAssignExpression` を新規追加
- **Rationale**: 既存の `AssignExpression` の型安全性を維持し、評価器での分岐も明確になる
- **Trade-offs**: AST ノードが1つ増えるが、保守性が向上

### Decision: スライスの範囲外処理
- **Context**: `arr[0:100]` のように範囲が配列長を超える場合の挙動
- **Alternatives Considered**:
  1. エラーを返す — 厳密だが使いにくい
  2. 自動クランプ — Python スタイルで利用可能な範囲に調整
- **Selected Approach**: 自動クランプ（Python スタイル）
- **Rationale**: Sugu の既存挙動（範囲外インデックスは null を返す）と一貫性を保ちつつ、スライスでは部分結果を返す方が実用的
- **Trade-offs**: バグの隠蔽リスクがあるが、スクリプト言語としての使いやすさを優先

## Risks & Mitigations
- **破壊的変更（負インデックス）**: 既存コードで `arr[-1]` → `null` を期待している場合に挙動が変わる。ただし既存テストに負インデックスの明示的テストはないためリスクは低い
- **後置演算子の副作用**: `i++` が式の値を返しつつ変数を更新するため、評価順序に注意が必要。複雑な式内での使用（`a[i++] + b[i++]`）は未定義とし、ドキュメントで注意喚起する
- **スライスとインデックス代入の組み合わせ**: `arr[1:3] = [...]` のようなスライス代入は Phase 5 ではサポートしない。読み取りのみ
