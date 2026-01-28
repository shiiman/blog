---
name: blog-import
description: WordPressから記事をインポートする。「記事をインポート」「既存記事を取得」「WordPressから取得」「記事をダウンロード」「投稿をローカルに」「記事を同期」などで起動。WordPress REST API経由で記事を取得しMarkdownとして保存。
---

# Import Blog Skill

WordPressから記事または固定ページをインポートし、ローカルにMarkdownファイルとして保存します。

## 前提条件

- `.env` ファイルにWordPress API認証情報が設定されていること
- `tools/wp-cli/wp-cli` がビルドされていること

## ワークフロー

### 1. インポート対象の確認

ユーザーに確認:
- 投稿をインポートするか、固定ページをインポートするか
- 全件か、特定の記事（ID指定）か
- 件数制限が必要か

### 2. 記事一覧の確認（任意）

```bash
# 投稿一覧を確認
./tools/wp-cli/wp-cli list posts

# 固定ページ一覧を確認
./tools/wp-cli/wp-cli list pages

# 下書きのみ表示
./tools/wp-cli/wp-cli list posts --status=draft
```

### 3. インポート実行

```bash
# 全投稿をインポート（最大100件）
./tools/wp-cli/wp-cli import posts

# 件数を制限
./tools/wp-cli/wp-cli import posts --limit=10

# 特定の投稿をインポート
./tools/wp-cli/wp-cli import post 123

# 固定ページをインポート
./tools/wp-cli/wp-cli import pages
./tools/wp-cli/wp-cli import page 45
```

### 4. 結果報告

インポート結果をユーザーに報告:
- インポートした記事数
- 保存先ディレクトリ
- 次のステップ（編集・更新方法）

## 保存先

- 投稿: `posts/YYYY-MM-DD_slug/article.md`
- 固定ページ: `pages/slug/page.md`

## CLIコマンドリファレンス

```bash
# 投稿インポート
./tools/wp-cli/wp-cli import posts [--limit=N]
./tools/wp-cli/wp-cli import post <id>

# 固定ページインポート
./tools/wp-cli/wp-cli import pages [--limit=N]
./tools/wp-cli/wp-cli import page <id>

# 出力先を指定
./tools/wp-cli/wp-cli import posts --output=my-posts/
```

## 重要な注意事項

- HTMLはMarkdownに自動変換される
- 画像URLはそのまま保持（ローカルダウンロードは手動）
- Front Matterに元のWordPress IDが保存される（更新時に使用）
