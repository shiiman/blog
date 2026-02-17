---
description: WordPress上の既存記事を更新する
---

1. 更新したい記事のMarkdownファイルパス（例: `posts/2025-01-03_my-article/article.md`）を確認します。
2. フロントマターに `id` が含まれていることを確認します。
3. `wp-cli` を使用して更新します。
// turbo
   ```bash
   ./tools/wp-cli/wp-cli update <article_path>
   ```
   ※ 既存のアイキャッチを差し替える場合は `--force-eyecatch` を使用してください。
