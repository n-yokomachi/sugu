# Lambda ビルドスキル

AWS Lambda 用の bootstrap バイナリをビルドし、デプロイ用の zip ファイルを作成します。

## 実行手順

1. **ビルドスクリプトを実行**
   ```powershell
   powershell -ExecutionPolicy Bypass -File scripts/build-lambda.ps1
   ```

2. **成果物の確認**
   ```bash
   ls -la dist/sugu-lambda.zip
   ```

## 成果物

| ファイル | 説明 |
|----------|------|
| `dist/bootstrap` | Lambda 用実行ファイル |
| `dist/sugu-lambda.zip` | デプロイ用 ZIP ファイル |

## 使用方法

```
/build-lambda
```

## デプロイ先設定

- ランタイム: `provided.al2023`
- アーキテクチャ: `x86_64`
- ハンドラー: `bootstrap`

## 注意事項

- ビルド前にテストが通ることを確認する（`go test ./...`）
- 成果物は AWS Lambda にアップロードして使用する
- Lambda 上で `main.sugu` ファイルを作成して実行コードを記述する
