---
description: 記事をWordPressに公開（デフォルトは下書き）
---

1. 公開したい記事のMarkdownファイルパス（例: `drafts/2025-01-03_my-article/article.md`）を確認します。
2. 記事ディレクトリを `drafts/` から `posts/` に移動します。
   ```bash
   mv drafts/2025-01-03_my-article posts/
   ```
3. `wp-cli` を使用して投稿します。
// turbo
   ```bash
   ./tools/wp-cli/wp-cli post posts/2025-01-03_my-article/article.md
   ```
   ※ 本番公開する場合は `--publish` フラグを付けてください。
   ```bash
   ./tools/wp-cli/wp-cli post posts/2025-01-03_my-article/article.md --publish
   ```
