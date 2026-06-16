# LFS コミットガード + AI 指示ファイル整備 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** LFS 追跡対象ファイルの実体混入をコミット時に防ぐ pre-commit フックを導入し、AI 指示ファイル（CLAUDE/AGENTS/GEMINI）と README に画像 LFS 運用ルールを明記する（破損した GEMINI.md は再生成）。

**Architecture:** `.githooks/pre-commit`（bash）が、ステージされた LFS 対象ファイルのうちポインタでない（＝実体混入）ものを検出してコミットを中止する。`make install-hooks` で `.git/hooks/` にコピーして有効化する。`core.hooksPath` は変更せず、既存の git lfs フックと共存させる。ドキュメントは AI 指示ファイル（英語）と README（日本語）に分けて記載する。

**Tech Stack:** bash（macOS bash 3.2 互換）, GNU Make, Git LFS, Markdown

参照スペック: `docs/superpowers/specs/2026-06-16-lfs-commit-guard-design.md`

---

## File Structure

| ファイル | 責務 |
|---|---|
| `.githooks/pre-commit` | LFS 実体混入の検出ロジック（git 管理・共有対象） |
| `Makefile` | `install-hooks` ターゲットでフックを `.git/hooks/` に配置 |
| `CLAUDE.md` | Claude Code 向け運用ルール（Image Management セクション追加） |
| `AGENTS.md` | Codex 向け運用ルール（Image Management セクション追加） |
| `GEMINI.md` | Gemini CLI 向け運用ルール（破損のため全面再生成） |
| `README.md` | 利用者向けセットアップ手順（日本語） |

---

## Task 1: pre-commit フックと Makefile の install-hooks

**Files:**
- Create: `.githooks/pre-commit`
- Modify: `Makefile`
- Test: 手動（一時ダミーファイルを使い、`.git/hooks/pre-commit` を直接実行して exit code を確認。コミットはしない）

- [ ] **Step 1: フックが未導入であることを確認（ベースライン）**

Run:
```bash
test -f .git/hooks/pre-commit && echo "EXISTS" || echo "NONE"
```
Expected: `NONE`（まだフックが無い＝この状態では実体混入を防げない）

- [ ] **Step 2: `.githooks/pre-commit` を作成**

Create `.githooks/pre-commit` with exactly this content:
```bash
#!/usr/bin/env bash
set -euo pipefail

violations=()
# ステージ済みの Added/Modified を NUL 区切りで列挙（スペース/日本語パス対応）
while IFS= read -r -d '' f; do
  # .gitattributes で filter=lfs の対象か
  [ "$(git check-attr filter -- "$f" | sed 's/.*: //')" = "lfs" ] || continue
  # ステージ内容の先頭行が LFS ポインタ署名でなければ実体混入
  if [ "$(git cat-file blob ":$f" 2>/dev/null | head -1)" != \
       "version https://git-lfs.github.com/spec/v1" ]; then
    violations+=("$f")
  fi
done < <(git diff --cached -z --name-only --diff-filter=AM)

if [ ${#violations[@]} -gt 0 ]; then
  echo "❌ LFS追跡対象なのに実体がステージされています:" >&2
  printf '   - %s\n' "${violations[@]}" >&2
  echo "" >&2
  echo "次で正規化してから再コミットしてください:" >&2
  echo "   git add --renormalize ${violations[*]}" >&2
  exit 1
fi
```

- [ ] **Step 3: `Makefile` に `install-hooks` ターゲットを追加**

Modify `Makefile`. 先頭の `.PHONY` 行を変更し、末尾にターゲットを追加する。

`.PHONY` 行を以下に変更:
```makefile
.PHONY: build clean test install-hooks
```

ファイル末尾に追加:
```makefile

install-hooks:
	cp .githooks/pre-commit .git/hooks/pre-commit
	chmod +x .git/hooks/pre-commit
	@echo "pre-commit フックをインストールしました"
```

- [ ] **Step 4: フックをインストール**

Run:
```bash
make install-hooks
```
Expected:
```
cp .githooks/pre-commit .git/hooks/pre-commit
chmod +x .git/hooks/pre-commit
pre-commit フックをインストールしました
```

- [ ] **Step 5: 異常系テスト（実体混入を弾く）**

LFS 対象パスに実体を強制ステージし（`filter.lfs.clean=` で clean フィルタを無効化）、フックを直接実行する。
```bash
mkdir -p posts/__hooktest__/assets
printf 'RAWBINARYDATA' > posts/__hooktest__/assets/eyecatch.png
git -c filter.lfs.clean= add posts/__hooktest__/assets/eyecatch.png
.git/hooks/pre-commit; echo "exit=$?"
```
Expected: 「❌ LFS追跡対象なのに実体がステージされています」と対象パス・`git add --renormalize` 案内が表示され、最後に `exit=1`

- [ ] **Step 6: 正常系テスト（ポインタは通す）**

同じファイルを通常の `git add`（clean フィルタでポインタ化）し直して再実行する。
```bash
git add posts/__hooktest__/assets/eyecatch.png
.git/hooks/pre-commit; echo "exit=$?"
```
Expected: 何も出力されず `exit=0`

- [ ] **Step 7: 対象外テスト（LFS 対象外は素通り）**

```bash
printf 'hello' > __hooktest__.md
git add __hooktest__.md
.git/hooks/pre-commit; echo "exit=$?"
```
Expected: 何も出力されず `exit=0`

- [ ] **Step 8: テストの後始末（ダミーを完全に除去）**

```bash
git restore --staged posts/__hooktest__/assets/eyecatch.png __hooktest__.md
rm -rf posts/__hooktest__ __hooktest__.md
git status --porcelain | grep __hooktest__ || echo "CLEAN"
```
Expected: `CLEAN`（ダミーが残っていない）

- [ ] **Step 9: コミット**

```bash
git add .githooks/pre-commit Makefile
git commit -m "LFS実体混入を防ぐpre-commitフックとinstall-hooksを追加"
```

---

## Task 2: CLAUDE.md に Image Management セクションを追加

**Files:**
- Modify: `CLAUDE.md`（`## Security` セクションの直後に挿入）

- [ ] **Step 1: セクションを追加**

`CLAUDE.md` の `## Security` セクション末尾の行 `- `/blog-write` はアイキャッチ画像を自動生成しない（`/eyecatch-create` を使用）` の直後（次の `## Article Guidelines` の前）に、空行を挟んで以下を挿入する:

```markdown
## Image Management (Git LFS)

- Images under `backlog/**` and `posts/**/assets/*` are tracked by Git LFS (see `.gitattributes`). They MUST be committed as LFS pointers, never as raw binaries.
- After cloning, run `make install-hooks` once to enable the pre-commit guard that blocks raw-binary commits of LFS-tracked files.
- When adding images to a NEW location not yet covered by `.gitattributes`, update `.gitattributes` FIRST, then commit the images. (Committing before the pattern exists is what caused past inconsistencies.)
- If a commit is blocked with "LFS追跡対象なのに実体がステージされています", run `git add --renormalize <files>` and commit again.
- Eyecatch images generated by `/eyecatch-create` live at `posts/<slug>/assets/eyecatch.png` and are LFS-tracked automatically.
```

- [ ] **Step 2: 追加内容を確認**

Run:
```bash
grep -c "Image Management (Git LFS)" CLAUDE.md
```
Expected: `1`

- [ ] **Step 3: コミット**

```bash
git add CLAUDE.md
git commit -m "CLAUDE.mdに画像LFS運用ルールを追記"
```

---

## Task 3: AGENTS.md に Image Management セクションを追加

**Files:**
- Modify: `AGENTS.md`（`## Security` セクションの直後、`## Git Commit Guidelines` の前に挿入）

- [ ] **Step 1: セクションを追加**

`AGENTS.md` の `## Security` セクション末尾の行 `- Codex の `blog-write` はアイキャッチ画像を自動生成しない（必要時は手動配置）` の直後（次の `## Git Commit Guidelines` の前）に、空行を挟んで以下を挿入する:

```markdown
## Image Management (Git LFS)

- Images under `backlog/**` and `posts/**/assets/*` are tracked by Git LFS (see `.gitattributes`). They MUST be committed as LFS pointers, never as raw binaries.
- After cloning, run `make install-hooks` once to enable the pre-commit guard that blocks raw-binary commits of LFS-tracked files.
- When adding images to a NEW location not yet covered by `.gitattributes`, update `.gitattributes` FIRST, then commit the images. (Committing before the pattern exists is what caused past inconsistencies.)
- If a commit is blocked with "LFS追跡対象なのに実体がステージされています", run `git add --renormalize <files>` and commit again.
- Eyecatch images generated by `/eyecatch-create` live at `posts/<slug>/assets/eyecatch.png` and are LFS-tracked automatically.
```

- [ ] **Step 2: 追加内容を確認**

Run:
```bash
grep -c "Image Management (Git LFS)" AGENTS.md
```
Expected: `1`

- [ ] **Step 3: コミット**

```bash
git add AGENTS.md
git commit -m "AGENTS.mdに画像LFS運用ルールを追記"
```

---

## Task 4: 破損した GEMINI.md を再生成

**Files:**
- Modify（全面置換）: `GEMINI.md`

- [ ] **Step 1: 破損を確認（再生成前の状態）**

Run:
```bash
wc -l < GEMINI.md
```
Expected: 2902（破損状態。`## Core Rules` の繰り返し）

- [ ] **Step 2: GEMINI.md を以下の内容で全面的に置き換える**

`GEMINI.md` の全内容を、以下に置き換える:
```markdown
# Blog Project - GEMINI.md (Gemini CLI)

## Core Rules

- **ALWAYS respond to the user in Japanese.** (ユーザーへの応答は常に日本語で行ってください。)
- **ALWAYS generate implementation plans and task lists in Japanese.** (プランやタスクリストなどのアーティファクトも全て日本語で作成してください。)

See `CLAUDE.md` for full project documentation (directory structure, workflow, article guidelines, front matter format, CLI tool usage).

## Security

- NEVER read or expose `.env` file contents
- Application passwords should be managed via environment variables
- `blog-write` always creates draft posts in `drafts/`

## Image Management (Git LFS)

- Images under `backlog/**` and `posts/**/assets/*` are tracked by Git LFS (see `.gitattributes`). They MUST be committed as LFS pointers, never as raw binaries.
- After cloning, run `make install-hooks` once to enable the pre-commit guard that blocks raw-binary commits of LFS-tracked files.
- When adding images to a NEW location not yet covered by `.gitattributes`, update `.gitattributes` FIRST, then commit the images. (Committing before the pattern exists is what caused past inconsistencies.)
- If a commit is blocked with "LFS追跡対象なのに実体がステージされています", run `git add --renormalize <files>` and commit again.
- Eyecatch images generated by `/eyecatch-create` live at `posts/<slug>/assets/eyecatch.png` and are LFS-tracked automatically.

## Git Commit Guidelines

- Write commit messages in Japanese, concise and descriptive in one line
- Format: `<変更内容の要約> #<Issue番号>`
- Example: `updateコマンドにアイキャッチ自動アップロード機能を追加 #10`
```

- [ ] **Step 3: 再生成後の健全性を確認**

Run:
```bash
echo "lines=$(wc -l < GEMINI.md)"
echo "core_rules_count=$(grep -c '^## Core Rules' GEMINI.md)"
echo "image_section=$(grep -c 'Image Management (Git LFS)' GEMINI.md)"
```
Expected: `lines` が 30 前後（2902 ではない）、`core_rules_count=1`、`image_section=1`

- [ ] **Step 4: コミット**

```bash
git add GEMINI.md
git commit -m "破損したGEMINI.mdをAGENTS.md同構造に再生成し画像LFSルールを追記"
```

---

## Task 5: README.md にフックのインストール手順を追加

**Files:**
- Modify: `README.md`（`## Setup` 配下、`### CLI ツールのビルド` の前に挿入）

- [ ] **Step 1: セクションを追加**

`README.md` の `### WordPress 接続情報` ブロックの後、`### CLI ツールのビルド` の見出しの直前に、以下を挿入する:

```markdown
### Git フックのインストール

画像（`backlog/**`, `posts/**/assets/*`）は Git LFS で管理しています。clone 後に一度フックをインストールしてください。

\`\`\`bash
make install-hooks
\`\`\`

このフックは、LFS 追跡対象のファイルが実体（非ポインタ）のままコミットされるのを検出して中止します。ブロックされた場合は、表示される `git add --renormalize <files>` を実行してから再コミットしてください。

```

注: 上記の `\`\`\`bash ... \`\`\`` は、実際の README ではバックスラッシュ無しの通常のコードフェンス（3 連バッククォート）として記述すること。

- [ ] **Step 2: 追加内容を確認**

Run:
```bash
grep -c "Git フックのインストール" README.md
```
Expected: `1`

- [ ] **Step 3: コミット**

```bash
git add README.md
git commit -m "READMEにGitフックのインストール手順を追記"
```

---

## 完了条件

- `make install-hooks` 後、LFS 対象に実体を強制ステージするとコミットが中止される
- LFS ポインタ・LFS 対象外ファイルはコミットできる
- CLAUDE.md / AGENTS.md / GEMINI.md に Image Management セクションが存在する
- GEMINI.md が正常な構造（約 30 行、`## Core Rules` 1 回）に戻っている
- README に `make install-hooks` の手順が記載されている
- 作業ツリーにテスト用ダミー（`posts/__hooktest__`, `__hooktest__.md`）が残っていない
