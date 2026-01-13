# Blog Project - GEMINI.md

## Core Rules
- **ALWAYS respond to the user in Japanese.** (ユーザーへの応答は常に日本語で行ってください。)
- **ALWAYS generate implementation plans and task lists in Japanese.** (プランやタスクリストなどのアーティファクトも全て日本語で作成してください。)

## Project Overview
Technical blog management system for shiimanblog.com (WordPress on ConoHa).

## Directory Structure
- `posts/` - Posts imported from WordPress (format: `YYYY-MM-DD_slug/`)
- `pages/` - Pages imported from WordPress (format: `slug/`)
- `drafts/` - New article drafts
- `templates/` - Article templates
- `tools/wp-cli/` - Go-based WordPress management CLI tool
- `backlog/` - Historical article assets

## Custom Workflows
The following workflows are defined in `.agent/workflows/`:
- `/blog-import` - Import existing posts from WordPress
- `/blog-write` - Create a new blog post as a draft
- `/blog-eyecatch` - Generate an eyecatch image for the blog post
- `/blog-publish` - Publish a post to WordPress (default state is draft)
- `/blog-update` - Update an existing post on WordPress


## Standard Workflow
1. Use `/blog-import` to import existing posts.
2. Use `/blog-write` to create a new post in `drafts/`.
3. Review and edit the post content.
4. Use `/blog-publish` to upload the post to WordPress as a draft.
5. Use `/blog-update` to update existing posts.
6. Verify and finalize on the WordPress dashboard.

## CLI Tool (wp-cli)
How to build and run:
```bash
cd tools/wp-cli
go build -o wp-cli .
./wp-cli --help
```

## Security
- DO NOT read or expose the contents of the `.env` file.
- Application passwords are managed via environment variables.
- For safety, the default post status is "draft".

## Article Guidelines
- Create titles in Japanese (max 50 characters).
- Include SEO-optimized excerpts.
- Use proper heading hierarchy (H2, H3, H4).
- Add code blocks with language specifications.

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
- Commit messages should be in Japanese, concise, and one line only.
- Format: `<Summary of changes> #<IssueNumber>`
- Example: `updateコマンドにアイキャッチ自動アップロード機能を追加 #10`
