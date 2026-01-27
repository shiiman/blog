---
name: write-blog
description: 技術ブログ記事を執筆する。「記事を書いて」「ブログ記事作成」「記事執筆」「新しい記事」「ブログを書く」「技術記事」「記事を作って」などで起動。SEOを意識した構成で記事を生成し、drafts/ディレクトリに保存。
---

# Write Blog Skill

技術ブログ記事を執筆し、投稿用のMarkdownファイルを生成します。

## ワークフロー

### 1. テーマ確認

ユーザーに以下を確認:
- 記事のテーマ・タイトル案
- 対象読者（初心者/中級者/上級者）
- 記事の目的（How-to/解説/レビュー）

### 2. 構成案作成

- SEOキーワードリサーチ（WebSearchを使用）
- 記事構成（見出し構造）を提案
- ユーザーの承認を得る

### 3. 記事執筆

`.claude/agents/blog-writer.md` のテンプレートに従って執筆:
- 技術的に正確な内容
- コード例は動作確認可能な形で記載
- 適切な見出し階層（H2→H3→H4）

### 4. ファイル保存

保存先: `drafts/YYYY-MM-DD_slug/article.md`

```bash
# 保存先ディレクトリ作成
mkdir -p drafts/$(date +%Y-%m-%d)_slug
```

### 5. 次のステップ案内

記事作成後、以下を案内:
- `/publish-blog` で WordPressに投稿（下書き）
- 記事内容の確認・編集方法
- カテゴリ・タグIDの確認方法（`wp-cli categories`, `wp-cli tags`）

## Front Matter形式

```yaml
---
title: "記事タイトル"
excerpt: "記事の要約（120文字程度）"
categories: [1]  # カテゴリID
tags: [10, 20]   # タグID
slug: "url-slug"
status: draft
---
```

## 重要な注意事項

- 記事タイトルは SEO を意識（50文字以内推奨）
- Front Matter の categories/tags は ID で指定
- コードブロックには言語を明示（```python など）
- 画像は `drafts/YYYY-MM-DD_slug/assets/` に配置
