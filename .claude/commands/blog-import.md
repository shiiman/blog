# Blog Import Command

WordPressから記事をインポートします。

## 使い方

```bash
/blog-import
/blog-import posts
/blog-import post <id>
/blog-import pages
```

## オプション

- `posts` - 全投稿をインポート
- `post <id>` - 特定の投稿をインポート
- `pages` - 全固定ページをインポート
- `page <id>` - 特定の固定ページをインポート
- `--limit=N` - インポート件数を制限

## 実行例

```bash
# 全投稿をインポート
/blog-import posts

# 最新10件をインポート
/blog-import posts --limit=10

# 特定の投稿をインポート
/blog-import post 123

# 固定ページをインポート
/blog-import pages
```

## Claudeへの指示

1. 引数がなければ何をインポートするか確認
2. 必要に応じて `./tools/wp-cli/wp-cli list` で一覧を表示
3. `./tools/wp-cli/wp-cli import` コマンドを実行
4. インポート結果と保存先を報告

$ARGUMENTS
