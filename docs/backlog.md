# Sugu 言語 バックログ

各 Phase のレビューで発見された課題をまとめたもの。
Phase 実装計画（implementation.md）に移動した項目は削除する。

---

## Phase 2 未実装タスク（低優先度）

### Step 4: REPL 機能強化

> 詳細は後述の「REPL パッケージ」セクションを参照

- 複数行入力対応
- コマンド履歴機能（readline ライブラリ導入）
- 入力補完機能（オプション）
- エラーハンドリング改善

### Step 5: コード品質向上

> 詳細は後述の各パッケージセクションを参照

- Environment の const ガード強化（Object パッケージ）
- break/continue の内部型隠蔽（Evaluator パッケージ）
- 無限ループ検出（オプション）
- スタックオーバーフロー対策（オプション）

---

## Phase 3 レビュー課題

### エラーハンドリング関連（PR #17）

- **try/catch に finally ブロックがない**
  - リソース解放などのクリーンアップ処理に有用
  - 将来的な拡張として検討

- **ファイルパスのサニタイズがない**
  - パストラバーサル攻撃への対策が未実装
  - `../` などを使った悪意あるパスへのアクセスが可能
  - サンドボックス設計など大きな変更が必要

- **throwValue が object.Object インターフェースを実装している**
  - 内部的な制御フロー用の型だが外部に漏れる可能性
  - `breakValue` や `continueValue` と同様の問題

---

## Phase 2 レビュー課題

### 配列関連（PR #12）

- 負のインデックス未対応（Python スタイル `arr[-1]`）
  - 現在は負のインデックスは `null` を返す
  - 仕様として意図的であれば問題なし
- ~~配列要素への代入未実装（`arr[0] = 10`）~~ → **Phase 3 で実装済み**
- スライス操作未実装（`arr[1:3]`）
  - Python や Go スタイルのスライス構文は未対応

### マップ関連（PR #13）

- 空マップと空ブロックの曖昧さ（`{}`）
  - 現在の実装では `{}` は空マップとしてパースされる
  - ブロック文として解釈されるべきコンテキストでの挙動を確認する必要がある
- ~~マップ要素への代入未実装（`map["key"] = value`）~~ → **Phase 3 で実装済み**
- 浮動小数点キーのハッシュ問題（`0.0` vs `-0.0`）
  - `NaN` 同士も異なるハッシュになる可能性

### 組み込み関数関連（PR #14）

- `concat()` / `append()` 関数未実装
  - `+` 演算子での結合も未対応
- `contains()` / `includes()` 関数未実装
- `keys()` / `values()` の順序が不定
  - Go の map イテレーションと同様に順序が保証されない

### 位置情報関連（PR #15）

- 演算子エラーに位置情報が未適用
  - 型の不一致、ゼロ除算などのエラーには位置情報が含まれていない
  - これらのエラーは現在ASTノードを受け取っていないため、リファクタリングが必要
- 組み込み関数エラーに位置情報なし
  - 組み込み関数はASTノードにアクセスできないため、位置情報を持てない
  - 呼び出し元のCallExpressionから位置情報を渡す仕組みが必要
- タブ文字の列カウント（仕様として許容）
  - タブ文字も1列としてカウントされる
  - エディタによっては表示がずれる可能性があるが、多くの言語で同様の実装
- ASTノードへの位置情報の直接保持
  - 現在はToken経由でアクセスしているが、直接Line/Columnフィールドを持つ方が明確
  - 現在の実装でも動作に問題はない

---

## 位置情報関連（横断的課題）

> Phase 2 implementation.md Step 3 に移動済み

---

## Lexer パッケージ

### Unicode 対応
```go
// 現在
ch byte // ASCII のみ

// 改善後
ch rune // Unicode 対応
```
- 日本語変数名などを使用可能にする
- `readChar()`, `peekChar()`, `isLetter()` の変更が必要

---

## Parser パッケージ

### for 文のセミコロン処理の見直し
```go
// parser.go:414-420
// parseStatement がセミコロンを消費するため、二重消費の可能性
```

### switch 文の case ブロック形式
```go
// 現在: case 1: { ... } の形式のみ
// 検討: case 1: stmt; の形式もサポートするか
```
- 仕様との整合性確認が必要

---

## Object パッケージ

### Environment の const 再代入ガード
```go
// 現在: SetConst で設定した変数に Set で上書き可能
// 改善: Environment 側でもチェックを追加
```

---

## AST パッケージ

### NumberLiteral の値を数値型に変更
```go
// 現在
type NumberLiteral struct {
    Token token.Token
    Value string  // "42" や "3.14"
}

// 改善後
type NumberLiteral struct {
    Token token.Token
    Value float64  // 数値として保持
}
```
- パーサーで変換することで評価時の効率向上

### テストカバレッジの拡充
- エッジケース（nil値、空の配列など）
- ネストした構造のテスト
- 複雑な式の組み合わせ

### ForStatement の String() 表現
```go
// "for (; ; ) { }" のような表示
// 仕様との整合性確認が必要
```

---

## Evaluator パッケージ

### break/continue の内部型を隠蔽
```go
// 現在: object.Object を実装している
type breakValue struct{}
func (b *breakValue) Type() object.ObjectType { return "BREAK" }

// 改善: 内部実装として完全に隠蔽するか、object パッケージに移動
```

### 無限ループ検出
```go
// while (true) { } や for (;;) { } で無限ループ
// 対策: 実行ステップ数やタイムアウトによる保護
```

### 引数の数チェック（警告オプション）
```go
// 現在: 引数が多すぎても無視、少なすぎると NULL
// 改善: 警告を出すオプションを追加
```

### スタックオーバーフロー対策
```go
// 再帰関数の深すぎる呼び出し対策
// 呼び出し深度のトラッキングと制限
```

### 組み込み関数のテスト改善
```go
// out, outln, in のテストがない
// io.Writer/io.Reader を注入可能にしてテスト容易性を向上
```

### 文字列インデックスのマルチバイト文字対応
```go
// 現在: バイト単位でアクセス（ASCII のみ正しく動作）
return &object.String{Value: string(stringObject.Value[idx])}

// 改善: rune に変換してからアクセス
runes := []rune(stringObject.Value)
return &object.String{Value: string(runes[idx])}
```
- `"あいう"[1]` が正しく `"い"` を返すようになる
- Lexer の Unicode 対応と合わせて対応推奨

---

## REPL パッケージ

### 複数行入力対応
```go
// 括弧のバランスをチェックして継続入力を促す
// 例: { を入力後、} が来るまで入力を継続
```

### コマンド履歴機能
- 上下矢印キーで過去の入力を呼び出し
- readline ライブラリの導入を検討

### 入力補完機能
- キーワード補完
- 変数名補完
- readline + カスタム補完ハンドラ

### Scanner.Err() のチェック追加
```go
// repl.go
if !scanned {
    if err := scanner.Err(); err != nil {
        // エラー処理
    }
    return
}
```

### RunFile のエラーメッセージ改善
```go
// ファイルが存在しない場合とパーミッションエラーを区別
if os.IsNotExist(err) {
    return fmt.Errorf("file not found: %s", filename)
} else if os.IsPermission(err) {
    return fmt.Errorf("permission denied: %s", filename)
}
```

---

## 将来の機能拡張（優先度順）

### 重要レベル（スケールするプログラムに必要）

- **モジュールシステム** - `import` で他のファイルを読み込む
- **文字列操作関数** - `split`, `join`, `trim`, `replace` など
- **数学関数** - `abs`, `floor`, `ceil`, `random` など

### あると便利レベル

- **HTTP通信** - API呼び出し機能
- **JSON操作** - パース/シリアライズ関数
- **日時操作** - 現在時刻の取得、日付計算
- **正規表現** - パターンマッチング

### 上級レベル（本格的な開発に必要）

- **非同期処理** - async/await や goroutine 相当
- **パッケージマネージャ** - 依存関係管理
- **デバッガ** - ステップ実行、ブレークポイント
