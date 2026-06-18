# wp-cli（移行専用ツール）

> ⚠️ **このツールは旧 WordPress（ConoHa）からの移行時に使用したもので、現在の記事公開フローには使いません。**
>
> 現在の公開は「記事を `posts/<YYYY-MM-DD_slug>/article.md` に置いて `main` に git push → GitHub Actions が Cloudflare Pages へ自動デプロイ」です。WordPress への publish/update 操作は廃止されています。

WordPress REST API と連携する Go 製 CLI です。移行時に記事・固定ページ・メディアのインポート等に使用しました。

## ビルド

```bash
cd tools/wp-cli
go build -o wp-cli .
./wp-cli --help
```

## 用途（移行時のみ）

WordPress からの記事・ページ・カテゴリ/タグのインポートなど。実行には `.env` に WordPress 接続情報（`WP_SITE_URL` / `WP_USERNAME` / `WP_APP_PASSWORD`）が必要です。日常の記事公開では不要です。
