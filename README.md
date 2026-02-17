# shiimanblog.com ブログ管理システム

ConoHa上で運営するWordPress技術ブログ（<https://shiimanblog.com>）の記事管理システムです。

## 機能

- 📝 **記事執筆** - Claude Code/Cursor/Codex/Antigravityでブログ記事を執筆
- 🎨 **アイキャッチ生成** - Cursor/AntigravityでAI画像生成（記事内容に合わせて自動生成）
- 📤 **記事投稿** - CLIでWordPressに投稿（デフォルト: 下書き）
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

### 3. CLIツールのビルド

```bash
cd tools/wp-cli
go build -o wp-cli .
```

## 使い方

### スキル/ワークフロー（推奨）

#### 記事管理（全ツール対応）

```bash
# 記事を書く
/blog-write

# 記事を投稿
/blog-publish

# 記事を更新
/blog-update
```

> **Note**: Codexのプロジェクトローカルスキルは `.agents/skills/` に配置します。

#### アイキャッチ画像生成（Cursor/Antigravity のみ）

```bash
# アイキャッチ画像を生成
/eyecatch-create
```

> **Note**: Claude Codeには画像生成機能がないため、アイキャッチは手動配置が必要です。

### CLIコマンド

```bash
# 投稿一覧
./tools/wp-cli/wp-cli list posts

# 記事インポート
./tools/wp-cli/wp-cli import posts
./tools/wp-cli/wp-cli import post 123

# 新規投稿（下書き）
./tools/wp-cli/wp-cli post drafts/2025-01-03_article/article.md

# 新規投稿（公開）
./tools/wp-cli/wp-cli post drafts/article.md --publish

# 記事更新
./tools/wp-cli/wp-cli update posts/2025-01-03_slug/article.md

# 固定ページ更新
./tools/wp-cli/wp-cli update pages/poker/page.md --page --publish

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
│   ├── poker/              # ポーカータイマー
│   └── mahjong/            # 麻雀点数計算
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

# 2. アイキャッチ画像を生成（AIが記事内容を分析して自動生成）
/eyecatch-create

# 3. 投稿（アイキャッチは自動でアップロード・設定）
./tools/wp-cli/wp-cli post drafts/2026-01-03_my-article/article.md --publish
```

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

# 3. 投稿（アイキャッチは自動でアップロード・設定）
./tools/wp-cli/wp-cli post drafts/2026-01-03_my-article/article.md --publish
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

## ライセンス

Private
