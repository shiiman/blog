---
description: 新しいブログ記事を下書きとして作成する
---

1. 記事のテーマやタイトルを決定します。
2. `drafts/` ディレクトリに新しいディレクトリを作成します（形式: `YYYY-MM-DD_slug`）。
3. その中に `article.md` を作成し、適切なフロントマターを設定します。
4. 記事の内容を執筆します。

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
