---
description: WordPressから既存記事をインポートする
---

1. `wp-cli` がビルドされているか確認します。
// turbo
2. `tools/wp-cli/wp-cli` が存在しない場合は、以下のコマンドでビルドします。
   ```bash
   cd tools/wp-cli && go build -o wp-cli .
   ```
3. 記事をインポートします。
// turbo
   ```bash
   ./tools/wp-cli/wp-cli import posts --limit=10
   ```
   ※ 必要に応じて `--limit` 引数を調整してください。
