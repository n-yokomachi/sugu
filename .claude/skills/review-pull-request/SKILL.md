---
name: review-pull-request
description: Pull Requestのレビューを行い、結果をreview.mdに記録する。PRレビュー、コードレビュー、レビュー依頼された場合に使用。
allowed-tools: Bash, Read, Edit, Write
---

# Pull Requestレビュースキル

Pull Requestのレビューを行い、結果をドキュメントに記録します。

## 実行手順

1. **PR情報の取得**
   ```bash
   gh pr view <PR番号> --json title,body,files,url
   ```

2. **変更内容の確認**
   ```bash
   gh pr diff <PR番号>
   ```

3. **レビュー実施**
   - 変更されたファイルを確認
   - コードの品質をチェック
   - テストが適切に書かれているか確認
   - 仕様との整合性を確認

4. **レビュー結果の記録**
   - PRの概要から対象Phaseを判定
   - `docs/phase<N>/review.md` に結果を記録
   - 既存のreview.mdの形式に従う
   - 新しいPRのレビューは既存の内容に追記する形で記録

5. **GitHubへのレビューコメント投稿（オプション）**
   ```bash
   gh pr review <PR番号> --approve -b "レビューコメント"
   ```

## 使用方法

```
/review-pull-request <PR番号>
```

## 注意事項

- レビューは客観的かつ建設的に行う
- 良い点も必ず記載する
- 改善提案は具体的に記載する
- Phase番号はPRのタイトルや本文から自動判定する
- review.mdが既に存在する場合は、既存の形式を踏襲して新しいセクションを追加する
- フォーマットは既存のreview.mdに合わせる（将来の改善点スタイル）
