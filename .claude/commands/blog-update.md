# Blog Update Command

既存記事の修正をWordPressに反映します。

## 使い方

```bash
/blog-update
/blog-update <file>
/blog-update --id=123
```

## オプション

- `<file>` - 更新するMarkdownファイル
- `--id=N` - WordPress投稿IDを明示的に指定
- `--publish` - 公開状態に変更
- `--page` - 固定ページとして更新
- `--dry-run` - 更新せずにプレビュー

## 実行例

```bash
# インポート済み記事を更新
/blog-update posts/2025-01-03_my-article/article.md

# 固定ページを更新
/blog-update pages/about/page.md --page

# 公開状態に変更
/blog-update posts/article.md --publish

# ドライラン
/blog-update posts/article.md --dry-run
```

## Claudeへの指示

1. 引数がなければ `posts/` または `pages/` から記事を選択させる
2. Front Matterのidを確認
3. `./tools/wp-cli/wp-cli update` コマンドを実行
4. 更新結果（ID・URL・ステータス）を報告

$ARGUMENTS
