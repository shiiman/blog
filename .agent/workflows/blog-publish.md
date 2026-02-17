---
description: 記事をWordPressに公開（デフォルトは公開）
---

1. 公開したい記事のMarkdownファイルパス（例: `drafts/2025-01-03_my-article/article.md`）を確認します。
2. `wp-cli` を使用して公開します（`post` はデフォルトで公開）。
// turbo
   ```bash
   ./tools/wp-cli/wp-cli post drafts/2025-01-03_my-article/article.md
   ```
3. 下書き保存が必要な場合のみ `--draft` を指定します。
   ```bash
   ./tools/wp-cli/wp-cli post drafts/2025-01-03_my-article/article.md --draft
   ```
4. 公開時は CLI が `drafts/` から `posts/` への移動と Front Matter 同期を自動で行います。
