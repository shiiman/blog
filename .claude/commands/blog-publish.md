# Blog Publish Command

記事をWordPressに投稿します。

## 使い方

```bash
/blog-publish
/blog-publish <file>
/blog-publish --publish
```

## オプション

- `<file>` - 投稿するMarkdownファイル
- `--publish` - 下書きではなく公開
- `--dry-run` - 投稿せずにプレビュー

## 実行例

```bash
# 最新の下書きを投稿
/blog-publish

# 特定ファイルを投稿
/blog-publish drafts/2025-01-03_conoha-setup/article.md

# 公開状態で投稿
/blog-publish --publish
```

## Claudeへの指示

1. 引数がなければ `drafts/` から記事を選択させる
2. 記事内容（タイトル・カテゴリ）を表示して確認
3. `./tools/wp-cli/wp-cli post` コマンドを実行
4. 投稿結果（ID・URL・ステータス）を報告

$ARGUMENTS
