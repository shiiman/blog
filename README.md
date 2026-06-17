# shiimanblog.com ブログ管理システム

ConoHa上で運営するWordPress技術ブログ（<https://shiimanblog.com>）の記事管理システムです。

## 機能

- 📝 **記事執筆** - Claude Code/Cursor/Codex/Antigravityでブログ記事を執筆
- 🎨 **アイキャッチ生成** - Cursor/AntigravityでAI画像生成（記事内容に合わせて自動生成）
- 📤 **記事投稿** - CLIでWordPressに公開（デフォルト: 公開）
- 📥 **記事インポート** - WordPressから既存記事をMarkdownとして取得
- ✏️ **記事更新** - ローカルで編集した記事をWordPressに反映
- 📄 **固定ページ管理** - 投稿と同様に固定ページも操作可能

## セットアップ

### 前提条件

- **Go 1.24以上** - CLIツールのビルドに必要です。[Go公式サイト](https://go.dev/dl/)からインストールしてください。

### 1. WordPress アプリケーションパスワードの発行

1. <https://shiimanblog.com/wp-admin/> にログイン
2. **ユーザー → プロフィール** に移動
3. ページ下部「アプリケーションパスワード」セクション
4. 名前に `Claude Blog CLI` を入力
5. 「新しいアプリケーションパスワードを追加」をクリック
6. 生成されたパスワードをコピー

### 2. 環境変数の設定

```bash
cp .env.example .env
```

`.env` を編集:

```bash
WP_SITE_URL=https://shiimanblog.com
WP_USERNAME=your-username
WP_APP_PASSWORD=xxxx xxxx xxxx xxxx xxxx xxxx
```

### 3. Git フックのインストール

画像（`backlog/**`, `posts/**/assets/*`）は Git LFS で管理しています。clone 後に一度フックをインストールしてください。

```bash
make install-hooks
```

このフックは、LFS 追跡対象のファイルが実体（非ポインタ）のままコミットされるのを検出して中止します。ブロックされた場合は、表示される `git add --renormalize <files>` を実行してから再コミットしてください。

### 4. CLIツールのビルド

```bash
cd tools/wp-cli
go build -o wp-cli .
```

## 使い方

### 変更点（運用ルール統一）

- `blog-write` は下書き作成専用（`drafts/` + `status: draft`）
  - Gemini/Cursor は `blog-write` 内で `assets/eyecatch.png` も自動生成
- `blog-publish` は公開専用（`wp-cli post` はデフォルト公開）
- `blog-update` は既存記事更新専用

### スキル/ワークフロー（推奨）

#### 記事管理（全ツール対応）

```bash
# 記事を書く
/blog-write

# 記事を公開
/blog-publish

# 記事を更新
/blog-update
```

> **Note**: Codexのプロジェクトローカルスキルは `.agents/skills/` に配置します。

#### アイキャッチ画像再生成（Cursor/Antigravity のみ）

```bash
# アイキャッチ画像を再生成（必要時のみ）
/eyecatch-create
```

> **Note**: Claude Codeには画像生成機能がないため、アイキャッチは手動配置が必要です。

### CLIコマンド

```bash
# 投稿一覧
./tools/wp-cli/wp-cli list posts
./tools/wp-cli/wp-cli list posts --status=draft  # status: draft|publish|pending|private

# 記事インポート
./tools/wp-cli/wp-cli import posts
./tools/wp-cli/wp-cli import post 123

# 新規投稿（公開・デフォルト）
./tools/wp-cli/wp-cli post drafts/2025-01-03_article/article.md

# 新規投稿（下書き）
./tools/wp-cli/wp-cli post drafts/article.md --draft

# 記事更新
./tools/wp-cli/wp-cli update posts/2025-01-03_slug/article.md

# 固定ページ更新
./tools/wp-cli/wp-cli update pages/poker/page.md --page

# 固定ページ作成（公開する場合はFront Matterで status: publish を指定）
./tools/wp-cli/wp-cli page drafts/about/page.md

# カテゴリ・タグ一覧
./tools/wp-cli/wp-cli categories
./tools/wp-cli/wp-cli tags
```

## ディレクトリ構成

```
blog/
├── posts/                  # インポート済み投稿（YYYY-MM-DD_slug/）
│   └── 2025-01-03_slug/
│       ├── article.md
│       └── assets/         # 記事用アセット
│           └── eyecatch.png
├── pages/                  # インポート済み固定ページ（slug/）
│   ├── contact/            # お問い合わせ
│   ├── mahjong/            # 麻雀点数計算
│   ├── poker/              # ポーカータイマー
│   ├── privacy-policy/     # プライバシーポリシー
│   ├── profile/            # プロフィール
│   └── sitemap/            # サイトマップ
├── drafts/                 # 新規下書き記事
├── templates/              # 記事テンプレート
├── tools/wp-cli/           # Go製CLIツール
├── .agents/
│   └── skills/             # Codexプロジェクトローカルスキル定義
├── .cursor/
│   └── skills/             # Cursorスキル定義
├── .claude/
│   ├── agents/             # 記事執筆エージェント
│   └── skills/             # Claude Codeスキル定義
├── .agent/
│   └── workflows/          # Antigravity (Gemini) ワークフロー
└── backlog/                # 過去の記事画像アセット
```

## Front Matter形式

### 投稿

```yaml
---
id: 123                    # WordPress投稿ID（更新時に使用）
title: "記事タイトル"
slug: "url-slug"
status: draft              # draft | publish
excerpt: "記事の要約"
categories: [1, 5]         # カテゴリID
tags: [10, 20]             # タグID
featured_media: 456        # アイキャッチ画像のメディアID
---
```

### 固定ページ

```yaml
---
id: 45
title: "ページタイトル"
slug: "about"
status: publish
parent: 0
menu_order: 0
---
```

## アイキャッチ画像

記事ディレクトリの `assets/eyecatch.png` に画像を配置すると、投稿時に自動でWordPressにアップロードしてアイキャッチに設定します。

### 画像生成ワークフロー

#### Cursor/Antigravity（AI自動生成）

```bash
# 1. 記事を作成
/blog-write

# 2. 必要ならアイキャッチ画像を再生成
/eyecatch-create

# 3. 公開（アイキャッチは自動でアップロード・設定）
./tools/wp-cli/wp-cli post drafts/2026-01-03_my-article/article.md
```

`/blog-write` 実行時に、記事本文の作成と同時に `assets/eyecatch.png` が自動生成されます。生成に失敗した場合でも `article.md` は保存されるため、必要に応じて `/eyecatch-create` で再生成してください。

**デザインルール:**
- アスペクト比: 16:9 (1280x720)
- ベースカラー: 白 (#FFFFFF)
- アクセントカラー: 青 (#007BFF) + 濃いグレー (#424242)
- スタイル: テクニカル・クリーン

#### Claude Code/手動配置

```bash
# 1. 記事を作成
/blog-write

# 2. アイキャッチ画像を手動で配置
# drafts/2026-01-03_my-article/assets/eyecatch.png

# 3. 公開（アイキャッチは自動でアップロード・設定）
./tools/wp-cli/wp-cli post drafts/2026-01-03_my-article/article.md
```

> **Note**: Claude Codeには画像生成機能がないため、外部ツールで画像を作成して手動配置してください。

## 固定ページ

### ポーカータイマー (`pages/poker/`)

ポーカートーナメント用のブラインドタイマー。

- フルスクリーン対応（PC/モバイル）
- ブラインドレベル・プライズ配分のカスタマイズ
- レベルアップ時のチャイム音（iOS/iPad対応）
- 残り10秒で時間表示が赤く点滅

### 麻雀点数計算 (`pages/mahjong/`)

麻雀の点数計算ツール。

## トラブルシューティング

### wp-cli ビルドエラー

```
go build: command not found
```

**対処法:** Go がインストールされていません。[Go公式サイト](https://go.dev/dl/)からインストールしてください（Go 1.24以上が必要）。

```
go: module requires Go >= 1.24
```

**対処法:** Go のバージョンが古いです。`go version` で確認し、1.24以上にアップデートしてください。

### Goツールチェーンの不整合エラー

```
compile: version "go1.25.5" does not match go tool version "go1.25.7"
```

**原因:** `GOROOT` が古いGoに固定されている状態で、`GOTOOLCHAIN=auto` が新しいツールチェーンを選択すると発生します。

**対処法:**
1. 現在のシェルで `unset GOROOT` を実行
2. `mise` 利用時は `mise settings set go_set_goroot false` を実行（恒久対策）
3. 新しいシェルを開き直して `go test ./...` を再実行

### .env 未設定時のエラー

```
Error: WP_SITE_URL is not set
```

**対処法:**
1. `.env.example` をコピーして `.env` を作成: `cp .env.example .env`
2. `.env` に WordPress の接続情報を設定してください（セットアップ手順を参照）

### WordPress API 認証失敗

```
Error: 401 Unauthorized
```

**対処法:**
- `WP_USERNAME` と `WP_APP_PASSWORD` が正しいか確認
- アプリケーションパスワードにスペースが含まれていることを確認（`xxxx xxxx xxxx xxxx` 形式）
- WordPress管理画面でアプリケーションパスワードが有効か確認
- パスワードを再発行して `.env` を更新

### WordPress API 接続エラー

```
Error: connection refused / timeout
```

**対処法:**
- `WP_SITE_URL` が正しいか確認（末尾のスラッシュは不要）
- サイトが稼働中か確認: `curl -I https://shiimanblog.com`
- WordPress REST API が有効か確認: `curl https://shiimanblog.com/wp-json/wp/v2/posts`

## 計画4: 本番デプロイ・DNS切替（Cloudflare Pages）

### リダイレクト生成・検証

```bash
# public/_redirects を生成（カテゴリ/タグの日本語slug → enSlug + feed系）
npm run build:redirects

# 旧sitemap-1.xmlの全URLが新サイトでカバーされているか検証（要 npm run build）
npm run verify:redirects
```

- `public/_redirects` は git 管理対象（生成後コミット）
- カテゴリ/タグを追加・変更した場合は再生成してコミットする

### Cloudflare Pages 環境変数（本番投入）

Cloudflare Pages のダッシュボードで直接入力する（コード・ドキュメントに実値をハードコードしない）。

#### 公開値（ビルド時埋め込み、Production 環境に設定後に再ビルド必要）

| キー | 用途 |
| --- | --- |
| `PUBLIC_TURNSTILE_SITE_KEY` | Turnstile ウィジェット（問い合わせ） |
| `PUBLIC_GISCUS_REPO` | giscus 対象リポジトリ（例 `shiiman/blog`） |
| `PUBLIC_GISCUS_REPO_ID` | giscus リポジトリID |
| `PUBLIC_GISCUS_CATEGORY` | giscus カテゴリ名 |
| `PUBLIC_GISCUS_CATEGORY_ID` | giscus カテゴリID |

#### 秘密値（Functions 実行時参照、暗号化扱い）

| キー | 用途 |
| --- | --- |
| `TURNSTILE_SECRET_KEY` | Turnstile サーバー検証 |
| `RESEND_API_KEY` | Resend メール送信 |
| `CONTACT_FROM_EMAIL` | 送信元（Resend 独自ドメイン認証後 `noreply@shiimanblog.com`） |
| `CONTACT_TO_EMAIL` | 受信先（問い合わせ通知先） |

### 切替 runbook（概要）

| Phase | 内容 | 担当 |
| --- | --- | --- |
| 0 | 旧 sitemap-1.xml 取得 → `data/` 保存（完了済み） | 完了 |
| A | `_redirects` 生成・テスト・ローカル検証（完了済み） | 完了 |
| B | Cloudflare Pages 作成・GitHub 連携 → プレビューURL発行 | shiiman |
| C | 本番環境変数を Cloudflare ダッシュボードで投入 → 再ビルド | shiiman |
| D | プレビューURL で全機能・リダイレクト検証 | shiiman+Claude |
| E1 | NS を ConoHa → Cloudflare へ切替（Aレコードは旧WPのまま） | shiiman |
| E2 | Resend で独自ドメイン認証（SPF/DKIM/DMARC） | shiiman |
| E3 | Turnstile 許可ホストに `shiimanblog.com` を追加 | shiiman |
| E4 | Pages にカスタムドメイン割当 → **この瞬間に新サイトへ切替** | shiiman |
| F | 本番検証（全URL/リダイレクト/問い合わせ実送信） | shiiman+Claude |
| G | Search Console に新サイトマップ再送信 | shiiman |
| H | ドメイン移管 → 約2週間モニタ → ConoHa 解約（不可逆・最後） | shiiman |

## 環境変数（計画3: 動的機能）

`.env`（公開値）と `.dev.vars`（秘密値・ローカル wrangler 用）はいずれも **手動作成**します（gitignore 済み）。本番値の投入は上記「計画4」を参照。

### 公開値（`.env`、ビルドに埋め込まれる）

| キー | 用途 |
| --- | --- |
| `PUBLIC_TURNSTILE_SITE_KEY` | Turnstile ウィジェット（問い合わせ） |
| `PUBLIC_GISCUS_REPO` | giscus 対象リポジトリ（例 `shiiman/blog`） |
| `PUBLIC_GISCUS_REPO_ID` | giscus リポジトリID |
| `PUBLIC_GISCUS_CATEGORY` | giscus カテゴリ名 |
| `PUBLIC_GISCUS_CATEGORY_ID` | giscus カテゴリID |

### 秘密値（`.dev.vars`、Functions のみ参照）

| キー | 用途 |
| --- | --- |
| `TURNSTILE_SECRET_KEY` | Turnstile サーバー検証 |
| `RESEND_API_KEY` | Resend メール送信 |
| `CONTACT_FROM_EMAIL` | 送信元（例 `noreply@shiimanblog.com`） |
| `CONTACT_TO_EMAIL` | 受信先（問い合わせ通知先） |

### ローカル検証
```bash
npm run build      # astro build + pagefind インデックス生成
npm run preview    # 検索・コメント（giscus）の確認
npm run pages:dev  # 問い合わせ Functions の確認（http://localhost:8788/contact/）
```

> 初回の `npm run pages:dev` で wrangler の依存（workerd）のビルド承認を求められた場合は `npm approve-builds` を実行してください（`npm run build` / `npm test` には不要）。

## ライセンス

Private
