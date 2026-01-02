# リリースビルドスキル

macOSとWindows向けのバイナリをビルドし、リリース用のzipファイルを作成します。

## 実行手順

1. **ビルドディレクトリの確認**
   ```bash
   ls -la build/
   ```

2. **macOS (Apple Silicon) 向けビルド**
   ```bash
   GOOS=darwin GOARCH=arm64 go build -o build/darwin-arm64/sugu .
   ```

3. **Windows (64bit) 向けビルド**
   ```bash
   GOOS=windows GOARCH=amd64 go build -o build/windows-amd64/sugu.exe .
   ```

4. **zipファイル作成**
   ```bash
   cd build && zip -j sugu-darwin-arm64.zip darwin-arm64/sugu
   cd build && zip -j sugu-windows-amd64.zip windows-amd64/sugu.exe
   ```

5. **成果物の確認**
   ```bash
   ls -la build/*.zip
   ```

## 成果物

| ファイル | 対象 |
|----------|------|
| `build/sugu-darwin-arm64.zip` | macOS Apple Silicon |
| `build/sugu-windows-amd64.zip` | Windows 64bit |

## 使用方法

```
/build-release
```

## 注意事項

- ビルド前にテストが通ることを確認する（`go test ./...`）
- バージョンタグを作成してからビルドすることを推奨
- 成果物はGitHub Releasesにアップロードして使用する
