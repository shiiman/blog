---
name: blog-publish
description: "Markdown記事をWordPressに公開する。「記事を投稿」「WordPressに公開」「ブログ投稿」「記事をアップ」「記事を公開」「WPに投稿」などで起動。drafts/の記事をWordPress REST API経由で公開。"
---


# Publish Blog Skill

Markdown記事をWordPressに公開します（デフォルト: 公開）。

## 前提条件

- `.env` ファイルにWordPress API認証情報が設定されていること
- `tools/wp-cli/wp-cli` がビルドされていること

## ワークフロー

### 1. 記事選択

下書き記事一覧を表示:
```bash
ls -la drafts/
```

ユーザーに投稿する記事を確認。

### 2. 記事内容確認

記事内容をプレビュー:
```bash
cat drafts/YYYY-MM-DD_slug/article.md
```

タイトル、カテゴリ、タグを確認。

### 3. 投稿実行

CLIツールで公開（デフォルト: 公開）:
```bash
# プロジェクトルートから実行
./tools/wp-cli/wp-cli post drafts/YYYY-MM-DD_slug/article.md

# ドライランで確認
./tools/wp-cli/wp-cli post drafts/YYYY-MM-DD_slug/article.md --dry-run

# 下書きで投稿したい場合
./tools/wp-cli/wp-cli post drafts/YYYY-MM-DD_slug/article.md --draft
```

### 4. 記事ディレクトリの移動（公開時）

`post` 実行後、公開された記事は CLI が自動で `drafts/` から `posts/` に移動し、Front Matter の `id/status/date/modified/featured_media` を同期します。

### 5. 結果報告

投稿結果をユーザーに報告:
- タイトル
- 投稿ID
- URL
- ステータス（下書き/公開）
- 記事の配置先（`drafts/` or `posts/`）

## CLIコマンドリファレンス

```bash
# 投稿（公開・デフォルト）
./tools/wp-cli/wp-cli post <file>

# 投稿（下書き）
./tools/wp-cli/wp-cli post <file> --draft

# ドライラン
./tools/wp-cli/wp-cli post <file> --dry-run

# 固定ページ投稿
./tools/wp-cli/wp-cli page <file>
```

## 重要な注意事項

- デフォルトは公開
- 下書きで投稿したい場合のみ `--draft` を指定
- カテゴリ・タグはIDで指定（`wp-cli categories` で確認可能）
