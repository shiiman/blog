---
description: 新しいブログ記事を下書きとして作成する
---

1. 記事のテーマやタイトルを決定します。
2. `drafts/` ディレクトリに新しいディレクトリを作成します（形式: `YYYY-MM-DD_slug`）。
3. その中に `article.md` を作成し、適切なフロントマターを設定します。
4. 記事の内容を執筆します。
5. 同じフロー内で記事内容を分析し、`assets/eyecatch.png` を自動生成します。
6. アイキャッチ自動生成に失敗した場合でも `article.md` はそのまま保存し、警告を出して `/eyecatch-create` での再生成を案内します。
7. `blog-write` は下書き作成専用です。公開は記事を `posts/<YYYY-MM-DD_slug>/article.md` に移動し、`main` に git push すると GitHub Actions が自動でビルド・デプロイします。

例:
```yaml
---
title: "タイトル"
slug: "url-slug"
date: 2026-01-03T12:00:00.000Z
draft: true            # true で非公開。公開時は false
categories: [savings]  # 文字列slugの配列（data/categories.json 参照）
tags: [mail]           # 文字列slugの配列（data/tags.json 参照）
eyecatch: ./assets/eyecatch.png  # 任意（相対パス）
---
```
