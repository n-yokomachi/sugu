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

## Phase 2 レビュー課題

### 配列関連（PR #12）

- 負のインデックス未対応（Python スタイル `arr[-1]`）
- 配列要素への代入未実装（`arr[0] = 10`）
- スライス操作未実装（`arr[1:3]`）

### マップ関連（PR #13）

- 空マップと空ブロックの曖昧さ（`{}`）
- マップ要素への代入未実装（`map["key"] = value`）
- 浮動小数点キーのハッシュ問題（`0.0` vs `-0.0`）
- キー削除機能未実装（`delete(map, key)`）

### 組み込み関数関連（PR #14）

- `concat()` / `append()` 関数未実装
- `contains()` / `includes()` 関数未実装
- `keys()` / `values()` の順序が不定

### 位置情報関連（PR #15）

- 演算子エラーに位置情報が未適用
- 組み込み関数エラーに位置情報なし
- タブ文字の列カウント（仕様として許容）

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

### 剰余演算子の浮動小数点対応
```go
// 現在
case "%":
    return &object.Number{Value: float64(int64(leftVal) % int64(rightVal))}

// 改善: math.Mod を使用
case "%":
    return &object.Number{Value: math.Mod(leftVal, rightVal)}
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
