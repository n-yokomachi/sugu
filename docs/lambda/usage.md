# AWS Lambda での Sugu 実行

Sugu は AWS Lambda 上で実行できます。

## Lambda ランタイムのビルド

```powershell
# Windows
powershell -ExecutionPolicy Bypass -File scripts/build-lambda.ps1

# Linux/Mac
./scripts/build-lambda.sh
```

`dist/sugu-lambda.zip`（bootstrap を含む）が生成されます。

## デプロイ

1. AWS Lambda 関数を作成
   - ランタイム: `Amazon Linux 2023` (provided.al2023)
   - アーキテクチャ: `x86_64`

2. `sugu-lambda.zip` をアップロード

3. Lambda コンソールのコードエディタで `main.sugu` を作成

## main.sugu の例

```javascript
// テストイベントの JSON は event 変数でアクセス可能
outln("Hello, " + event["name"] + "!");
event["age"] + 10;
```

## テストイベント

```json
{
  "name": "Taro",
  "age": 25
}
```

## レスポンス

```json
{
  "output": "Hello, Taro!\n",
  "result": "35",
  "error": null
}
```

## 制限事項

- `in()` 関数は使用不可（Lambda は標準入力を持たない）
- 実行時間は Lambda のタイムアウト設定に依存

## 関連ドキュメント

- [Lambda 実装計画](implementation.md) - 実装の詳細
