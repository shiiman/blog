# Blog Project - CLAUDE.md

## Core Rules

- **ALWAYS respond to the user in Japanese.** (ユーザーへの応答は常に日本語で行ってください。)
- **ALWAYS generate implementation plans and task lists in Japanese.** (プランやタスクリストなどのアーティファクトも全て日本語で作成してください。)

This file contains project-specific instructions for Claude Code.

## Project Overview

Technical blog for shiimanblog.com, built with **Astro** and deployed to **Cloudflare Pages**.
Articles are managed as Markdown files; pushing to `main` triggers GitHub Actions to build and deploy automatically.

> Migrated from WordPress (ConoHa). `tools/wp-cli/` and the migration scripts in `scripts/` are **migration-only** and are NOT part of the publishing flow.

## Directory Structure

- `src/` - Astro source (`content.config.ts` is the source of truth for front matter; `pages/`, `layouts/`, `components/`, `lib/`)
- `posts/` - Articles (format: `YYYY-MM-DD_slug/article.md`)
- `pages/` - Static pages (format: `slug/page.md`)
- `drafts/` - New article drafts
- `public/` - Static assets served as-is (`_redirects`, `ads.txt`, favicon)
- `functions/` - Cloudflare Pages Functions (contact form, etc.)
- `data/` - `categories.json` / `tags.json` / `permalinks.json` (in use; needed for redirects & canonical URLs)
- `templates/` - Article templates
- `tools/wp-cli/` - **Migration-only** Go CLI for WordPress (not used for publishing)
- `backlog/` - Historical article assets

## Skills

- `/blog-write` - Write a new blog article draft into `drafts/`
- `/eyecatch-create` - Generate eyecatch image (Cursor/Antigravity; Claude Code cannot generate images)

For Codex, project-local skills are defined in `.agents/skills/` (`blog-write`).

## Workflow (publishing = git push)

1. Use `/blog-write` to create a new draft in `drafts/`
2. Review and edit the article
3. (Optional) Use `/eyecatch-create` to generate `assets/eyecatch.png`
4. Place the article at `posts/<YYYY-MM-DD_slug>/article.md`
5. Commit and push to `main` → GitHub Actions runs test → build → `wrangler pages deploy` and deploys to Cloudflare Pages

There is NO publish/update step to WordPress. Opening a PR produces a Cloudflare Pages preview URL (`https://<branch>.shiimanblog.pages.dev`).

## Security

- NEVER read or expose `.env`, `.dev.vars`, or `.prd.vars` contents
- Cloudflare Pages secrets must be set in the dashboard or via `npm run setup:cf-secrets` — never hardcode them
- `/blog-write` always creates drafts in `drafts/`
- `/blog-write` does NOT auto-generate the eyecatch image (use `/eyecatch-create`)

## Image Management (Git LFS)

- Images under `backlog/**` and `posts/**/assets/*` are tracked by Git LFS (see `.gitattributes`). They MUST be committed as LFS pointers, never as raw binaries.
- After cloning, run `make install-hooks` once to enable the pre-commit guard that blocks raw-binary commits of LFS-tracked files.
- When adding images to a NEW location not yet covered by `.gitattributes`, update `.gitattributes` FIRST, then commit the images.
- If a commit is blocked with "LFS追跡対象なのに実体がステージされています", run `git add --renormalize <files>` and commit again.
- Eyecatch images live at `posts/<YYYY-MM-DD_slug>/assets/eyecatch.png` and are LFS-tracked automatically.

## Article Guidelines

- Write titles in Japanese (70 chars max)
- Include SEO-optimized excerpts
- Use proper heading hierarchy (H2, H3, H4)
- Add code blocks with language specification
- Slugs should use hyphen-separated format (e.g., `my-article-title`)

## Front Matter Format

Schema source of truth: `src/content.config.ts`.

```yaml
---
title: "Title"
slug: "url-slug"
date: 2026-01-03T12:00:00.000Z
excerpt: "Post summary"
categories: [savings]          # string slugs (see data/categories.json)
tags: [mail, freelance]        # string slugs (see data/tags.json)
eyecatch: ./assets/eyecatch.png  # relative path (optional)
draft: false                   # true = excluded from build
# modified: 2026-01-04T09:00:00.000Z   # optional
# id: 123                              # legacy WordPress post ID (only for migrated posts' URL mapping)
---
```

- `categories` / `tags` are **string slugs**, NOT numeric IDs.
- `draft: true` excludes the post from the build (the old `status: draft|publish` is gone).
- Eyecatch is a **relative path** in `eyecatch` (the old `featured_media` media ID is gone).

## Cloudflare Pages / Deploy

- `npm run build:redirects` — generates `public/_redirects` (commit after running; regenerate when categories/tags change)
- `npm run verify:redirects` — verifies all old URLs are covered (requires `npm run build` first)
- `data/old-sitemap-1.xml` — old WordPress sitemap captured during migration (do not delete)
- CI/CD: `.github/workflows/deploy.yml`. `PUBLIC_*` env vars are embedded at **build time** (set them on the build step).
- Cloudflare Pages secrets must be entered in the dashboard or via `npm run setup:cf-secrets` — never hardcode them

## Git Commit Guidelines

- Write commit messages in Japanese, concise and descriptive in **one line**
- Format: `<変更内容の要約> #<Issue番号>`
- Example: `カテゴリページにインフィード広告を追加 #10`
