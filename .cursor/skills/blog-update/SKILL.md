---
name: blog-update
description: 既存記事の修正をWordPressに反映する。「記事を更新」「修正を反映」「記事を同期」「変更をアップロード」「記事を上書き」「更新を投稿」などで起動。ローカルで編集したMarkdown記事をWordPressに反映。
---

# Update Blog Skill

ローカルで編集したMarkdown記事の内容をWordPressに反映します。

## 前提条件

- `.env` ファイルにWordPress API認証情報が設定されていること
- `tools/wp-cli/wp-cli` がビルドされていること
- 記事のFront Matterに `id` フィールドが設定されていること

## ワークフロー

### 1. 更新対象の確認

更新する記事を確認:
```bash
# インポート済み投稿一覧
ls -la posts/

# インポート済み固定ページ一覧
ls -la pages/
```

### 2. 記事内容の確認

記事のFront Matterを確認（特にid）:
```bash
head -20 posts/YYYY-MM-DD_slug/article.md
```

### 3. 更新実行

```bash
# 投稿を更新（IDはFront Matterから取得）
./tools/wp-cli/wp-cli update posts/YYYY-MM-DD_slug/article.md

# 固定ページを更新
./tools/wp-cli/wp-cli update pages/slug/page.md --page

# IDを明示的に指定
./tools/wp-cli/wp-cli update drafts/article.md --id=123

# 同時に公開状態に変更
./tools/wp-cli/wp-cli update posts/article.md --publish

# ドライランで確認
./tools/wp-cli/wp-cli update posts/article.md --dry-run
```

### 4. 結果報告

更新結果をユーザーに報告:
- 更新した投稿ID
- URL
- 新しいステータス

## CLIコマンドリファレンス

```bash
# 投稿更新
./tools/wp-cli/wp-cli update <file>

# 固定ページ更新
./tools/wp-cli/wp-cli update <file> --page

# ID指定
./tools/wp-cli/wp-cli update <file> --id=123

# 公開状態に変更
./tools/wp-cli/wp-cli update <file> --publish

# ドライラン
./tools/wp-cli/wp-cli update <file> --dry-run
```

## Front Matter形式

更新時に必要なフィールド:

```yaml
---
id: 123  # WordPress投稿ID（必須）
title: "更新後のタイトル"
slug: "url-slug"
status: draft  # draft | publish
categories: [1, 5]
tags: [10, 20]
---
```

## 重要な注意事項

- `id` フィールドがないと更新できない
- 更新前にドライラン（`--dry-run`）で確認を推奨
- ステータスを変更しない場合は `--publish` を付けない
