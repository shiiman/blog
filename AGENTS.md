# Blog Project - AGENTS.md (Codex)

## Core Rules

- **ALWAYS respond to the user in Japanese.** (ユーザーへの応答は常に日本語で行ってください。)
- **ALWAYS generate implementation plans and task lists in Japanese.** (プランやタスクリストなどのアーティファクトも全て日本語で作成してください。)

See `CLAUDE.md` for full project documentation (directory structure, workflow, article guidelines, front matter format, CLI tool usage).

## Codex Skills

Project-local skills are defined in `.agents/skills/`:

- `blog-write` - Write a new blog article
- `blog-publish` - Publish article to WordPress
- `blog-update` - Update existing article on WordPress

## Security

- NEVER read or expose `.env` file contents
- Application passwords should be managed via environment variables
- `blog-write` always creates draft posts in `drafts/`
- Codex の `blog-write` はアイキャッチ画像を自動生成しない（必要時は手動配置）

## Git Commit Guidelines

- Write commit messages in Japanese, concise and descriptive in one line
- Format: `<変更内容の要約> #<Issue番号>`
- Example: `updateコマンドにアイキャッチ自動アップロード機能を追加 #10`
