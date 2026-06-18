# Blog Project - AGENTS.md (Codex)

## Core Rules

- **ALWAYS respond to the user in Japanese.** (ユーザーへの応答は常に日本語で行ってください。)
- **ALWAYS generate implementation plans and task lists in Japanese.** (プランやタスクリストなどのアーティファクトも全て日本語で作成してください。)

See `CLAUDE.md` for full project documentation (directory structure, workflow, article guidelines, front matter format).

## Project Overview

Astro blog deployed to Cloudflare Pages. Publishing = push to `main` → GitHub Actions builds and deploys. There is NO WordPress publish/update step (migrated away from WordPress/ConoHa).

## Codex Skills

Project-local skills are defined in `.agents/skills/`:

- `blog-write` - Write a new blog article draft into `drafts/`

## Security

- NEVER read or expose `.env`, `.dev.vars`, or `.prd.vars` contents
- Cloudflare Pages secrets must be set in the dashboard or via `npm run setup:cf-secrets` — never hardcode them
- `blog-write` always creates drafts in `drafts/`
- Codex's `blog-write` does NOT auto-generate the eyecatch image (place it manually if needed)

## Image Management (Git LFS)

- Images under `backlog/**` and `posts/**/assets/*` are tracked by Git LFS (see `.gitattributes`). They MUST be committed as LFS pointers, never as raw binaries.
- After cloning, run `make install-hooks` once to enable the pre-commit guard that blocks raw-binary commits of LFS-tracked files.
- When adding images to a NEW location not yet covered by `.gitattributes`, update `.gitattributes` FIRST, then commit the images.
- If a commit is blocked with "LFS追跡対象なのに実体がステージされています", run `git add --renormalize <files>` and commit again.
- Eyecatch images live at `posts/<YYYY-MM-DD_slug>/assets/eyecatch.png` and are LFS-tracked automatically.

## Cloudflare Pages / Deploy

- `npm run build:redirects` — generates `public/_redirects` (commit after running)
- `npm run verify:redirects` — verifies all old URLs are covered (requires `npm run build` first)
- `data/old-sitemap-1.xml` — old WordPress sitemap captured during migration (do not delete)
- CI/CD: `.github/workflows/deploy.yml`. `PUBLIC_*` env vars are embedded at build time.

## Git Commit Guidelines

- Write commit messages in Japanese, concise and descriptive in one line
- Format: `<変更内容の要約> #<Issue番号>`
- Example: `カテゴリページにインフィード広告を追加 #10`
