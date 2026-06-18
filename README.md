# shiimanblog.com ブログ管理システム

[Astro](https://astro.build/) + [Cloudflare Pages](https://pages.cloudflare.com/) で構築した技術ブログ（<https://shiimanblog.com>）のソース・記事管理リポジトリです。記事は Markdown で管理し、`main` に push すると GitHub Actions が自動でビルド・デプロイします。

> 旧 WordPress（ConoHa）から移行済みです。WordPress 連携ツール（`tools/wp-cli/`）と `scripts/` の一部は**移行専用**で、通常の記事公開には使いません（後述「移行用ツール」を参照）。

## 技術スタック

- **静的サイトジェネレータ**: Astro（`output: static`, `trailingSlash: always`）
- **ホスティング**: Cloudflare Pages（GitHub Actions から `wrangler pages deploy`）
- **記事**: Markdown + YAML front matter（`posts/<YYYY-MM-DD_slug>/article.md`）
- **検索**: Pagefind ／ **コメント**: giscus ／ **問い合わせ**: Cloudflare Pages Functions + Turnstile + Resend
- **アクセス解析/広告**: Google Analytics 4 / Google AdSense
- **画像**: Git LFS

## 執筆支援スキル

- 📝 `/blog-write` — 記事の下書きを `drafts/` に作成
- 🎨 `/eyecatch-create` — アイキャッチ画像を生成（Cursor/Antigravity。Claude Code は画像生成不可のため手動配置）

> 公開は WordPress への投稿ではなく **git push** です（下記「記事の執筆と公開」を参照）。

## セットアップ

### 前提条件

- **Node.js**（LTS 推奨）
- **Git LFS**（画像管理に使用）

### 手順

```bash
git clone git@github.com:shiiman/blog.git
cd blog
npm ci
make install-hooks   # LFS追跡ファイルの実体コミットを防ぐ pre-commit フック
```

`make install-hooks` のフックは、LFS 追跡対象のファイルが実体（非ポインタ）のままコミットされるのを検出して中止します。ブロックされた場合は、表示される `git add --renormalize <files>` を実行してから再コミットしてください。

## 記事の執筆と公開

```
1. /blog-write で drafts/ に下書きを作成（または手動で記事を作成）
2. 必要なら /eyecatch-create で assets/eyecatch.png を生成
3. 記事を posts/<YYYY-MM-DD_slug>/article.md に配置
4. main へ commit & push
5. GitHub Actions が test → build → wrangler pages deploy を実行し
   Cloudflare Pages に自動デプロイ（数分で本番反映）
```

- **公開 = git push** です。WordPress への publish/update 操作はありません。
- CI/CD は [`.github/workflows/deploy.yml`](.github/workflows/deploy.yml)。`main` への push と PR で起動します。
- **PR を作ると Cloudflare Pages がプレビューURL**（`https://<branch>.shiimanblog.pages.dev`）を発行します。マージ前の確認に使えます。
- カテゴリ/タグを追加・変更したら `npm run build:redirects` を実行して `public/_redirects` を再生成・コミットしてください。

## ディレクトリ構成

```
blog/
├── src/                    # Astro のソース（レイアウト・コンポーネント・ページ・lib）
│   ├── content.config.ts   # 記事/固定ページの front matter スキーマ（唯一の正）
│   ├── pages/              # ルーティング（記事・一覧・カテゴリ・タグ・RSS 等）
│   ├── layouts/ components/ lib/
├── posts/                  # 記事（YYYY-MM-DD_slug/article.md）
│   └── 2025-01-03_slug/
│       ├── article.md
│       └── assets/eyecatch.png   # アイキャッチ（Git LFS）
├── pages/                  # 固定ページ（slug/page.md）
├── drafts/                 # 新規下書き
├── public/                 # 静的配信物（_redirects, ads.txt, favicon 等）
├── functions/              # Cloudflare Pages Functions（問い合わせ等）
├── data/                   # categories.json / tags.json / permalinks.json 等（現役・URL維持に使用）
├── scripts/                # ビルド補助 + 移行用ワンショット（下記参照）
├── templates/              # 記事テンプレート
├── tools/wp-cli/           # 【移行専用】WordPress連携 Go CLI（通常は不使用）
├── .github/workflows/      # GitHub Actions（deploy.yml）
├── .claude/ .agents/ .cursor/ .agent/   # 各AIツールのスキル/エージェント定義
└── backlog/                # 過去の記事画像アセット（Git LFS）
```

## Front Matter 形式

スキーマの正は [`src/content.config.ts`](src/content.config.ts) です。

### 記事（`posts/<YYYY-MM-DD_slug>/article.md`）

```yaml
---
title: "記事タイトル"
slug: "url-slug"
date: 2026-01-03T12:00:00.000Z
excerpt: "記事の要約（メタディスクリプションに使用）"
categories: [savings]          # 文字列 slug の配列（data/categories.json 参照）
tags: [mail, freelance]        # 文字列 slug の配列（data/tags.json 参照）
eyecatch: ./assets/eyecatch.png  # アイキャッチ画像の相対パス（任意）
draft: false                   # true で非公開（ビルド対象外）
# modified: 2026-01-04T09:00:00.000Z   # 更新日時（任意）
# id: 123                              # 旧WordPress投稿ID（移行記事のURL維持用。新規は不要）
---
```

- `categories` / `tags` は**文字列 slug**です（旧 WordPress の数値IDではありません）。利用可能な slug は `data/categories.json` / `data/tags.json` を参照。
- `draft: true` の記事はビルドされません（旧 `status: draft|publish` は廃止）。
- アイキャッチは `eyecatch` に**相対パス**で指定します（旧 `featured_media`（メディアID）は廃止）。

### 固定ページ（`pages/<slug>/page.md`）

```yaml
---
title: "ページタイトル"
slug: "about"
date: 2026-01-03T12:00:00.000Z
draft: false
---
```

## ローカル開発

```bash
npm run dev        # 開発サーバー（http://localhost:4321）
npm run build      # astro build + Pagefind 検索インデックス生成（dist/）
npm run preview    # build 成果物のプレビュー（検索・giscus の確認）
npm run pages:dev  # 問い合わせ Functions の確認（http://localhost:8788/contact/）
npm test           # vitest（ロジックのユニットテスト）
```

問い合わせ Functions をローカルで動かすには秘密値が必要です。`.dev.vars`（gitignore 済み）を手動作成してください（下記「環境変数」）。

## リダイレクト（`public/_redirects`）

旧 WordPress の URL を新サイトへ 301 リダイレクトするための設定です。

```bash
npm run build:redirects   # public/_redirects を生成（カテゴリ/タグの日本語slug→enSlug + feed系）
npm run verify:redirects  # 旧 sitemap の全URLが新サイトでカバーされるか検証（要 npm run build）
```

- `public/_redirects` は git 管理対象です（生成後にコミット）。
- カテゴリ/タグを追加・変更したら再生成してコミットしてください。
- `data/old-sitemap-1.xml` は移行元の sitemap（検証に使用。削除しない）。

## 環境変数

`.env`（公開値）・`.dev.vars`（ローカル Functions 用の秘密値）はいずれも**手動作成**します（gitignore 済み）。本番値は Cloudflare Pages のダッシュボード、または `npm run setup:cf-secrets`（`.prd.vars` から投入）で設定します。

### 公開値（ビルド時に埋め込まれる）

| キー | 用途 |
| --- | --- |
| `PUBLIC_TURNSTILE_SITE_KEY` | Turnstile ウィジェット（問い合わせ） |
| `PUBLIC_GISCUS_REPO` | giscus 対象リポジトリ（例 `shiiman/blog`） |
| `PUBLIC_GISCUS_REPO_ID` | giscus リポジトリID |
| `PUBLIC_GISCUS_CATEGORY` | giscus カテゴリ名 |
| `PUBLIC_GISCUS_CATEGORY_ID` | giscus カテゴリID |

> `PUBLIC_*` は **ビルド時**に埋め込まれます。GitHub Actions では `build` ステップの env に渡してください（deploy ステップではなく）。

### 秘密値（Functions 実行時に参照）

| キー | 用途 |
| --- | --- |
| `TURNSTILE_SECRET_KEY` | Turnstile サーバー検証 |
| `RESEND_API_KEY` | Resend メール送信 |
| `CONTACT_FROM_EMAIL` | 送信元（例 `noreply@shiimanblog.com`） |
| `CONTACT_TO_EMAIL` | 受信先（問い合わせ通知先） |

- ローカル: `.dev.vars`／本番: Cloudflare ダッシュボード or `.prd.vars` + `npm run setup:cf-secrets`
- **秘密値はコード・ドキュメントにハードコードしない**こと。

## 画像管理（Git LFS）

- `backlog/**` と `posts/**/assets/*` の画像は Git LFS で管理します（`.gitattributes` 参照）。実体（生バイナリ）ではなく **LFS ポインタ**としてコミットしてください。
- clone 後に `make install-hooks` を一度実行してください。
- 新しい場所に画像を追加する場合は、先に `.gitattributes` のパターンを更新してからコミットします。
- コミットが「LFS追跡対象なのに実体がステージされています」でブロックされたら、`git add --renormalize <files>` を実行して再コミット。

## 移行用ツール（通常は使いません）

旧 WordPress からの移行時に使用したもので、**日常の記事公開には不要**です。

- `tools/wp-cli/` — WordPress REST API 連携の Go CLI（記事インポート等に使用した）。公開フローには関与しません。
- `scripts/` のワンショット — `export:wp` / `export:permalinks` / `localize:images` / `migrate:frontmatter` / `cleanup:tracking`。移行時の一度きりの処理です。
- 継続利用する `scripts/` は `build:redirects` / `verify:redirects` / `setup:cf-pages` / `setup:cf-secrets` のみ。

## トラブルシューティング

### ビルドで画像が読めない（CI）

```
[NoImageMetadata] Could not process image metadata ...
```

**原因:** Git LFS の画像がポインタファイルのまま取得されている。
**対処:** GitHub Actions の checkout で `lfs: true` を指定（`.github/workflows/deploy.yml` 設定済み）。ローカルでは `git lfs pull`。

### Turnstile/giscus が表示されない

**原因:** `PUBLIC_*` 環境変数が**ビルド時**に渡っていない。
**対処:** `npm run build` の env（GitHub Actions の build ステップ）に `PUBLIC_*` を設定。deploy ステップに渡しても埋め込まれません。

### リダイレクトが効かない（404）

**対処:** `public/_redirects` を確認。トレイリングスラッシュ有無（`/feed` と `/feed/` の両方）、パーセントエンコードは**大文字**（Cloudflare の実リクエストに合わせる）、子カテゴリは親込みパス（例 `/category/fire/savings/`）。`npm run build:redirects` で再生成。

### AdSense 広告がプレビューで出ない

プレビュードメイン（`*.pages.dev`）では AdSense が 403 を返し広告は表示されません。**本番ドメイン（`shiimanblog.com`）のみ**で配信されます。

## ライセンス

Private
