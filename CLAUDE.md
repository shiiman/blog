# Blog Project - CLAUDE.md

## Core Rules
- **ALWAYS respond to the user in Japanese.** (ユーザーへの応答は常に日本語で行ってください。)
- **ALWAYS generate implementation plans and task lists in Japanese.** (プランやタスクリストなどのアーティファクトも全て日本語で作成してください。)

This file contains project-specific instructions for Claude Code.


## Project Overview

Technical blog management system for shiimanblog.com (WordPress on ConoHa).

## Directory Structure

- `posts/` - Imported posts from WordPress (format: `YYYY-MM-DD_slug/`)
- `pages/` - Imported pages from WordPress (format: `slug/`)
- `drafts/` - New article drafts
- `templates/` - Article templates
- `tools/wp-cli/` - Go CLI tool for WordPress management
- `backlog/` - Historical article assets

## Commands

- `/blog-write` - Write a new blog article
- `/blog-publish` - Publish article to WordPress (draft by default)
- `/blog-import` - Import existing articles from WordPress
- `/blog-update` - Update existing article on WordPress

## Workflow

1. Use `/blog-import` to import existing articles from WordPress
2. Use `/blog-write` to create new article in `drafts/`
3. Review and edit the article
4. Use `/blog-publish` to post to WordPress as draft
5. Use `/blog-update` to update existing articles
6. Finalize in WordPress dashboard

## CLI Tool (wp-cli)

Build and run:

```bash
cd tools/wp-cli
go build -o wp-cli .
./wp-cli --help
```

## Security

- NEVER read or expose `.env` file contents
- Application passwords should be managed via environment variables
- Default posting status is "draft" for safety

## Article Guidelines

- Write titles in Japanese (50 chars max)
- Include SEO-optimized excerpts
- Use proper heading hierarchy (H2, H3, H4)
- Add code blocks with language specification

## Front Matter Format

Posts use YAML front matter:

```yaml
---
id: 123          # WordPress post ID (for updates)
title: "Title"
slug: "url-slug"
status: draft    # draft | publish
categories: [1]  # Category IDs
tags: [10, 20]   # Tag IDs
---
```

## Git Commit Guidelines

- Write commit messages in Japanese, concise and descriptive in **one line**
- Format: `<変更内容の要約> #<Issue番号>`
- Example: `updateコマンドにアイキャッチ自動アップロード機能を追加 #10`
