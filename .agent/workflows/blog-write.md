---
description: 新しいブログ記事を下書きとして作成する
---

1. 記事のテーマやタイトルを決定します。
2. `drafts/` ディレクトリに新しいディレクトリを作成します（形式: `YYYY-MM-DD_slug`）。
3. その中に `article.md` を作成し、適切なフロントマターを設定します。
4. 記事の内容を執筆します。
5. 同じフロー内で記事内容を分析し、`assets/eyecatch.png` を自動生成します。
6. アイキャッチ自動生成に失敗した場合でも `article.md` はそのまま保存し、警告を出して `/eyecatch-create` での再生成を案内します。
7. 公開は `/blog-publish` で行います（`blog-write` は下書き作成専用）。

例:
```yaml
---
title: "タイトル"
slug: "url-slug"
status: draft
categories: [1]
tags: [10, 20]
---
```
