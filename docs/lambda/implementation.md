# AWS Lambda 対応 実装計画

## 概要

Sugu インタプリタを AWS Lambda で動作させるための実装。
既存コードには一切変更を加えず、Lambda 専用の実装を `lambda/` ディレクトリに追加する。

## 要件

- ランタイム: `provided.al2023` (x86_64)
- 入力: JSON で Sugu コードを受け取る
- 出力: 実行結果（出力文字列、エラー）を JSON で返す
- 既存コードへの変更: **なし**

## API 仕様

### リクエスト

```json
{
  "code": "mut x = 1 + 2; outln(x);"
}
```

### レスポンス（成功時）

```json
{
  "output": "3\n",
  "result": "3",
  "error": null
}
```

### レスポンス（エラー時）

```json
{
  "output": "",
  "result": null,
  "error": "line 1, column 5: undefined variable: foo"
}
```

---

## 実装タスク

### Step 1: プロジェクト設定
- [ ] `aws-lambda-go` SDK を `go.mod` に追加
- [ ] `lambda/` ディレクトリを作成

### Step 2: Lambda 用組み込み関数
- [ ] `lambda/builtins.go` を作成
- [ ] `out`, `outln` を出力キャプチャ版に置き換え
- [ ] `in` は Lambda 環境では使用不可としてエラーを返す
- [ ] その他の組み込み関数は既存のものを再利用

### Step 3: Lambda ハンドラー
- [ ] `lambda/main.go` を作成
- [ ] リクエストの JSON パース
- [ ] Sugu コードの実行（lexer → parser → evaluator）
- [ ] 出力のキャプチャとレスポンス生成

### Step 4: ビルドスクリプト
- [ ] Linux x86_64 向けクロスコンパイル設定
- [ ] ZIP パッケージ作成スクリプト

### Step 5: テスト
- [ ] `lambda/handler_test.go` を作成
- [ ] 正常系テスト
- [ ] エラー系テスト
- [ ] 出力キャプチャのテスト

---

## ディレクトリ構成

```
lambda/
├── main.go           # Lambda エントリーポイント
├── builtins.go       # Lambda 用組み込み関数
├── handler.go        # ハンドラーロジック
└── handler_test.go   # テスト
```

---

## 技術的な詳細

### 出力キャプチャの仕組み

Lambda 用の `out` / `outln` は `strings.Builder` に出力を書き込む。
実行後にビルダーの内容をレスポンスの `output` フィールドに設定する。

```go
type OutputCapture struct {
    builder strings.Builder
}

func (o *OutputCapture) Write(s string) {
    o.builder.WriteString(s)
}

func (o *OutputCapture) String() string {
    return o.builder.String()
}
```

### 組み込み関数の上書き

Environment に Lambda 用の組み込み関数を登録することで、
既存の evaluator をそのまま使いながら出力をキャプチャする。

```go
env := object.NewEnvironment()
// Lambda 用の out/outln を登録
env.Set("out", lambdaOut)
env.Set("outln", lambdaOutln)
// in は使用不可
env.Set("in", lambdaInDisabled)
```

---

## ビルドコマンド

```bash
# Linux x86_64 向けビルド
GOOS=linux GOARCH=amd64 go build -o bootstrap lambda/main.go

# ZIP パッケージ作成
zip sugu-lambda.zip bootstrap
```

---

## デプロイ

1. AWS Lambda 関数を作成（ランタイム: provided.al2023, アーキテクチャ: x86_64）
2. `sugu-lambda.zip` をアップロード
3. ハンドラー名: `bootstrap`（カスタムランタイムのデフォルト）

---

## 制限事項

- `in()` 関数は使用不可（Lambda は標準入力を持たない）
- 実行時間は Lambda のタイムアウト設定に依存
- ファイル I/O は `/tmp` ディレクトリのみ使用可能（将来の Phase 3 対応時）
