---
id: 1903
title: Claude CodeでWordPressブログをMarkdown管理する方法
slug: claude-code-wordpress-markdown
status: publish
excerpt: WordPressブログの記事をMarkdownでローカル管理し、Claude Codeと自作CLIツールで執筆・投稿を自動化するワークフローを紹介します。
categories:
    - 18
date: 2026-01-03T11:44:46
---

## はじめに

WordPressでブログを運営していると、こんな悩みはありませんか？

- ブラウザのエディタで長文を書くのがつらい
- 記事のバージョン管理ができない
- オフラインで執筆できない
- 複数の記事を横断的に編集しづらい

私も同じ悩みを抱えていました。そこで、**Markdown + Git + Claude Code**で記事を管理するシステムを構築しました。

この記事では、WordPressブログをMarkdownで管理し、Claude Codeで執筆を効率化するワークフローを紹介します。

## システム構成

### 全体のアーキテクチャ

```
┌─────────────────────────────────────────────────────────┐
│  ローカル環境                                              │
│  ┌───────────┐    ┌───────────┐    ┌───────────┐       │
│  │  Markdown │───▶│  wp-cli   │───▶│ WordPress │       │
│  │  (Git管理) │◀───│  (Go製)   │◀───│ REST API  │       │
│  └───────────┘    └───────────┘    └───────────┘       │
│        ▲                                                │
│        │                                                │
│  ┌─────┴─────┐                                          │
│  │Claude Code│                                          │
│  │ (執筆支援) │                                          │
│  └───────────┘                                          │
└─────────────────────────────────────────────────────────┘
```

### 使用技術

| 技術 | 用途 |
|------|------|
| Go | CLIツール実装 |
| WordPress REST API | 記事の取得・投稿 |
| Claude Code | 執筆支援・ワークフロー自動化 |
| Git | 記事のバージョン管理 |

### ディレクトリ構成

```
blog/
├── posts/           # インポート済み投稿
│   └── YYYY-MM-DD_slug/
│       ├── article.md
│       └── assets/
│           └── eyecatch.png  # アイキャッチ画像
├── pages/           # 固定ページ
├── drafts/          # 下書き
├── tools/wp-cli/    # Go製CLIツール
└── .claude/         # Claude Code設定
    ├── commands/    # スラッシュコマンド
    ├── skills/      # スキル定義
    └── agents/      # エージェント定義
```

## Go製CLIツール（wp-cli）の機能

WordPress REST APIを操作するためのCLIツールをGoで実装しました。

### 主要コマンド

```bash
# 記事一覧の表示
./wp-cli list posts

# 全記事をインポート
./wp-cli import posts

# 特定の記事をインポート
./wp-cli import post 123

# 記事を投稿（下書き）
./wp-cli post drafts/2026-01-03_my-article/article.md

# 記事を公開
./wp-cli post drafts/2026-01-03_my-article/article.md --publish

# 既存記事を更新
./wp-cli update posts/2026-01-03_my-article/article.md
```

### 記事のフォーマット

記事はYAMLフロントマター付きのMarkdownで管理します。

```markdown
---
id: 123              # WordPress投稿ID（更新時に使用）
title: "記事タイトル"
slug: "url-slug"
status: draft        # draft | publish
categories: [1, 2]   # カテゴリID
tags: [10, 20]       # タグID
---

## 見出し

本文をMarkdownで記述...
```

### HTML ↔ Markdown の自動変換

インポート時にWordPressのHTMLコンテンツを自動的にMarkdownに変換します。

```go
// html-to-markdownライブラリを使用
converter := md.NewConverter("", true, nil)
mdContent, err := converter.ConvertString(htmlContent)
```

投稿時は逆にMarkdownをHTMLに変換してAPIに送信します。

### アイキャッチ画像のアップロード

記事ディレクトリの`assets/eyecatch.png`に画像を配置すると、投稿時に自動でWordPressにアップロードしてアイキャッチに設定します。

```bash
# アイキャッチ付きで投稿（画像は自動アップロード）
./wp-cli post drafts/2026-01-03_my-article/article.md --publish
```

アイキャッチ画像は[Gemini](https://gemini.google.com/)などの画像生成AIで作成し、手動で配置できます。

## Claude Codeとの連携

Claude Codeのスキル機能を使って、記事執筆のワークフローを自動化しています。

### 利用可能なスキル

| スキル | 説明 |
|--------|------|
| `/blog-write` | 新規記事の執筆（SEOリサーチ含む） |
| `/blog-import` | WordPressから記事をインポート |
| `/blog-publish` | 記事をWordPressに投稿 |
| `/blog-update` | 既存記事を更新 |

### 記事執筆の流れ

`/blog-write`を実行すると、Claude Codeが以下のステップで記事執筆を支援します。

1. **テーマの確認** - 何について書くか確認
2. **SEOリサーチ** - WebSearchでキーワード調査
3. **構成提案** - 見出し構造を提案して承認を得る
4. **本文執筆** - テンプレートに沿って記事を執筆
5. **保存** - `drafts/`ディレクトリに保存

### スキルの定義例

`.claude/skills/blog-write.md`でスキルを定義しています。

```markdown
# Blog Write Command

技術ブログ記事を執筆します。

## Claudeへの指示

1. テーマがない場合はユーザーにテーマを聞く
2. WebSearchでSEOキーワードをリサーチ
3. 見出し構造を提案して承認を得る
4. 記事を執筆
5. drafts/YYYY-MM-DD_slug/article.md に保存
```

## 実際のワークフロー

### 1. 既存記事のインポート

まず、WordPressにある既存記事をローカルに取り込みます。

```bash
# Claude Codeで実行
/blog-import posts

# または直接CLIを実行
./tools/wp-cli/wp-cli import posts
```

79件の投稿が`posts/`ディレクトリにMarkdown形式で保存されます。

### 2. 新規記事の作成

Claude Codeで`/blog-write`を実行し、対話的に記事を作成します。

```
User: /blog-write

Claude: ブログ記事を書きます。テーマは何ですか？

User: Dockerの基本的な使い方

Claude: SEOキーワードをリサーチします...
        以下の見出し構成を提案します：
        ## はじめに
        ## Dockerとは
        ...
        この構成でよろしいですか？

User: はい

Claude: 記事を執筆しました。
        drafts/2026-01-03_docker-basics/article.md に保存しました。
```

### 3. WordPressへの投稿

記事が完成したら`/blog-publish`で投稿します。

```
User: /blog-publish

Claude: どの記事を投稿しますか？
        1. drafts/2026-01-03_docker-basics/article.md

User: 1

Claude: 下書きとして投稿しますか？それとも公開しますか？

User: 下書き

Claude: 投稿しました！
        URL: https://shiimanblog.com/?p=1850
```

## メリットと今後の展望

### このシステムのメリット

1. **Git管理によるバージョン履歴**
   - 記事の変更履歴を追跡できる
   - 複数人での共同編集も可能

2. **オフライン編集対応**
   - ネット環境がなくても執筆可能
   - お気に入りのエディタで編集できる

3. **AIアシストによる執筆効率化**
   - SEOリサーチを自動化
   - 構成提案で迷わない
   - 一貫したフォーマットで記事を作成

4. **双方向同期**
   - WordPress → ローカル（インポート）
   - ローカル → WordPress（投稿・更新）

### 今後の展望

- 予約投稿機能の追加

## まとめ

この記事では、WordPressブログをMarkdownで管理するシステムを紹介しました。

- **Go製CLIツール**でWordPress REST APIを操作
- **Claude Code**のスキル機能で執筆ワークフローを自動化
- **Git**で記事をバージョン管理

ブラウザに縛られない快適な執筆環境を実現できました。同じような悩みを持つ方の参考になれば幸いです。

## 参考リンク

- [WordPress REST API Handbook](https://developer.wordpress.org/rest-api/)
- [Claude Code 公式ドキュメント](https://docs.anthropic.com/en/docs/claude-code)
- [html-to-markdown (Go)](https://github.com/JohannesKaufmann/html-to-markdown)