# shiimanblog.com ブログ管理システム

ConoHa上で運営するWordPress技術ブログ（<https://shiimanblog.com>）の記事管理システムです。

## 機能

- 📝 **記事執筆** - Claude Codeでブログ記事を執筆
- 📤 **記事投稿** - CLIでWordPressに投稿（デフォルト: 下書き）
- 📥 **記事インポート** - WordPressから既存記事をMarkdownとして取得
- ✏️ **記事更新** - ローカルで編集した記事をWordPressに反映
- 📄 **固定ページ管理** - 投稿と同様に固定ページも操作可能

## セットアップ

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

### Claude Codeスキル（推奨）

```bash
# 記事を書く
/blog-write

# WordPressから既存記事をインポート
/blog-import posts

# 記事を投稿
/blog-publish

# 記事を更新
/blog-update
```

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

# カテゴリ・タグ一覧
./tools/wp-cli/wp-cli categories
./tools/wp-cli/wp-cli tags
```

## ディレクトリ構成

```
blog/
├── posts/                  # インポート済み投稿（YYYY-MM-DD_slug/）
├── pages/                  # インポート済み固定ページ（slug/）
├── drafts/                 # 新規下書き記事
├── templates/              # 記事テンプレート
├── tools/wp-cli/           # Go製CLIツール
├── .claude/
│   ├── agents/             # 記事執筆エージェント
│   ├── skills/             # スキル定義
│   └── commands/           # スラッシュコマンド
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

## ライセンス

Private
