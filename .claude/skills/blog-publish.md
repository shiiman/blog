---
description: Markdown記事をWordPressに投稿する。「記事を投稿」「WordPressに公開」「ブログ投稿」「記事をアップ」「下書き保存」「記事を公開」「WPに投稿」などで起動。drafts/の記事をWordPress REST API経由で投稿。
allowed-tools:
  - Read
  - Bash
  - Glob
  - AskUserQuestion
---

# Publish Blog Skill

Markdown記事をWordPressに投稿します（デフォルト: 下書き）。

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

CLIツールで投稿（デフォルト: 下書き）:
```bash
# プロジェクトルートから実行
./tools/wp-cli/wp-cli post drafts/YYYY-MM-DD_slug/article.md

# 公開する場合
./tools/wp-cli/wp-cli post drafts/YYYY-MM-DD_slug/article.md --publish

# ドライランで確認
./tools/wp-cli/wp-cli post drafts/YYYY-MM-DD_slug/article.md --dry-run
```

### 4. 記事ディレクトリの移動（公開時）

記事を公開した場合（`--publish` 指定、または後から `update --publish` で公開した場合）、記事ディレクトリを `drafts/` から `posts/` に移動する:
```bash
mv drafts/YYYY-MM-DD_slug posts/
```

※ `post --publish` で直接公開した場合は CLI が自動移動するため不要。下書き保存後に `update --publish` で公開した場合は手動で移動すること。

### 5. 結果報告

投稿結果をユーザーに報告:
- タイトル
- 投稿ID
- URL
- ステータス（下書き/公開）
- 記事の配置先（`drafts/` or `posts/`）

## CLIコマンドリファレンス

```bash
# 投稿（下書き）
./tools/wp-cli/wp-cli post <file>

# 投稿（公開）
./tools/wp-cli/wp-cli post <file> --publish

# ドライラン
./tools/wp-cli/wp-cli post <file> --dry-run

# 固定ページ投稿
./tools/wp-cli/wp-cli page <file>
```

## 重要な注意事項

- デフォルトは下書き保存（安全のため）
- 公開には明示的な `--publish` オプションが必要
- カテゴリ・タグはIDで指定（`wp-cli categories` で確認可能）
