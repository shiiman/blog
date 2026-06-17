# 計画3 動的3機能（検索・コメント・問い合わせ）Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 計画2の静的Astroサイトに、無料枠の動的3機能（Pagefind検索・giscusコメント・Pages Functions問い合わせ）を、既存のミニマル＋テーマ連動デザインを踏襲して追加する。

**Architecture:** `output: 'static'` を維持しアダプタは導入しない。検索はビルド後に `pagefind --site dist` でインデックス生成し、ヘッダーのオーバーレイUIが Pagefind JS API を叩く。コメントは giscus を `PostLayout` に組み込み、`MutationObserver` で `data-theme` を監視して `postMessage` でテーマ連動。問い合わせは `functions/api/contact.ts`（Pages Functions）が `src/lib/contact.ts` の純粋ロジック（vitestでテスト）を使い、入力検証→Turnstile検証→Resend送信を行う。フォームは `[...slug].astro` で contact ページにのみ差し込む。

**Tech Stack:** Astro v6（既存）/ Pagefind（npm CLI）/ giscus（client.js）/ Cloudflare Pages Functions / Turnstile / Resend / wrangler（ローカル検証）/ vitest（既存）。

## Global Constraints

- `astro.config.mjs` は変更しない（`output: 'static'` 維持、アダプタ追加なし）。
- シークレットはハードコードしない。`.env` / `.dev.vars` / Cloudflare 環境変数経由。`.env` の内容は読まない・出力しない。
- `.env*` ファイルはツールの権限制約で読み書きしない。必要キー名は `README.md` に明記し、`.env` / `.dev.vars` はユーザーが手動作成する。
- 公開値は `PUBLIC_` プレフィックス（ビルド埋め込み）、秘密値は Functions のみ参照。
- 既存テスト（vitest 43件）を壊さない。各タスク完了時に `npm test` が緑。
- デザインは既存の `src/styles/global.css` の CSS 変数（`--bg` `--bg-subtle` `--border` `--text` `--text-muted` `--accent` `--accent-hover` `--max-width` `--font-sans`）でライト/ダーク両対応。
- ブランチ `feature/blog-migration-plan-3` で作業。タスクごとにコミット。コミットメッセージは日本語一行 + `Co-Authored-By: Claude Opus 4.8 (1M context) <noreply@anthropic.com>` トレーラー。`--no-verify` 禁止。
- Node v26 / npm 11。vitest は設定ファイル無し（`*.test.ts` 自動検出）。

---

## Task 1: 検索（Pagefind ビルド統合 + ヘッダーオーバーレイUI）

**Files:**
- Modify: `package.json`（`build` スクリプト、`pagefind` を devDependencies へ）
- Create: `src/components/SearchOverlay.astro`
- Modify: `src/components/Header.astro`（検索ボタン追加 + SearchOverlay 配置）

**Interfaces:**
- Consumes: なし（Pagefind が生成する `/pagefind/pagefind.js` をクライアントで動的 import）
- Produces: ヘッダーの検索ボタン `#search-open`、オーバーレイ `#search-overlay`。後続タスクからの依存なし。

**Notes:**
- `dist/` は既に `.gitignore` 済み（`dist/pagefind/` の追記は不要）。
- `import('/pagefind/pagefind.js')` は **必ず `/* @vite-ignore */`** を付ける（ビルド時に Vite が解決しようとして失敗するのを防ぐ。pagefind.js はビルド後に生成される）。
- `astro dev` ではインデックスが無いため検索は動かない。検証は `npm run build && npm run preview`。

- [ ] **Step 1: pagefind を devDependencies に追加**

Run:
```bash
npm install -D pagefind
```
Expected: `package.json` の devDependencies に `pagefind` が追加され、インストール成功。

- [ ] **Step 2: build スクリプトを変更**

`package.json` の `scripts.build` を次に変更:
```json
"build": "astro build && pagefind --site dist",
```

- [ ] **Step 3: SearchOverlay.astro を作成**

Create `src/components/SearchOverlay.astro`:
```astro
---
// Pagefind を使うオーバーレイ検索。Header から呼び出す。
---
<div id="search-overlay" class="search-overlay" hidden role="dialog" aria-modal="true" aria-label="サイト内検索">
  <div class="search-modal">
    <input
      id="search-input"
      type="search"
      autocomplete="off"
      placeholder="記事を検索..."
      aria-label="検索キーワード"
    />
    <div id="search-results" class="search-results"></div>
  </div>
</div>

<style>
  .search-overlay {
    position: fixed; inset: 0; z-index: 100;
    background: rgba(0, 0, 0, 0.5);
    display: flex; align-items: flex-start; justify-content: center;
    padding: 10vh 20px 20px;
  }
  .search-overlay[hidden] { display: none; }
  .search-modal {
    width: 100%; max-width: var(--max-width);
    background: var(--bg); color: var(--text);
    border: 1px solid var(--border); border-radius: 12px;
    padding: 16px; max-height: 70vh; overflow-y: auto;
  }
  #search-input {
    width: 100%; font-size: 1rem; padding: 10px 12px;
    background: var(--bg-subtle); color: var(--text);
    border: 1px solid var(--border); border-radius: 8px;
    font-family: var(--font-sans);
  }
  #search-input:focus { outline: 2px solid var(--accent); }
  .search-results { margin-top: 12px; display: flex; flex-direction: column; gap: 8px; }
  .search-results :global(a) {
    display: block; padding: 10px 12px; border-radius: 8px;
    border: 1px solid var(--border); color: var(--text);
  }
  .search-results :global(a:hover) { background: var(--bg-subtle); text-decoration: none; }
  .search-results :global(strong) { display: block; color: var(--accent); }
  .search-results :global(span) { display: block; font-size: 0.85rem; color: var(--text-muted); }
  .search-results :global(p) { color: var(--text-muted); font-size: 0.9rem; }
</style>

<script>
  let pagefind: any = null
  const overlay = document.getElementById('search-overlay') as HTMLElement
  const input = document.getElementById('search-input') as HTMLInputElement
  const results = document.getElementById('search-results') as HTMLElement
  const openBtn = document.getElementById('search-open')

  async function loadPagefind() {
    if (pagefind) return pagefind
    try {
      pagefind = await import(/* @vite-ignore */ '/pagefind/pagefind.js')
    } catch {
      results.innerHTML = '<p>検索はビルド後（npm run build）に有効になります。</p>'
    }
    return pagefind
  }

  function open() {
    overlay.hidden = false
    input.focus()
    loadPagefind()
  }
  function close() {
    overlay.hidden = true
    input.value = ''
    results.innerHTML = ''
  }

  openBtn?.addEventListener('click', open)
  overlay.addEventListener('click', (e) => { if (e.target === overlay) close() })
  document.addEventListener('keydown', (e) => {
    if (e.key === 'Escape' && !overlay.hidden) close()
  })

  let timer: ReturnType<typeof setTimeout>
  input.addEventListener('input', () => {
    clearTimeout(timer)
    timer = setTimeout(runSearch, 200)
  })

  async function runSearch() {
    const q = input.value.trim()
    if (!q) { results.innerHTML = ''; return }
    const pf = await loadPagefind()
    if (!pf) return
    const search = await pf.search(q)
    const items = await Promise.all(search.results.slice(0, 10).map((r: any) => r.data()))
    if (items.length === 0) {
      results.innerHTML = '<p>該当する記事がありません。</p>'
      return
    }
    results.innerHTML = items
      .map((d: any) => `<a href="${d.url}"><strong>${d.meta?.title ?? ''}</strong><span>${d.excerpt}</span></a>`)
      .join('')
  }
</script>
```

- [ ] **Step 4: Header.astro に検索ボタンとオーバーレイを追加**

`src/components/Header.astro` を次のように改修（`ThemeToggle` の前に検索ボタン、`</header>` の後に `SearchOverlay` を配置）:
```astro
---
import ThemeToggle from './ThemeToggle.astro'
import SearchOverlay from './SearchOverlay.astro'
---
<header class="site-header">
  <div class="container inner">
    <a class="brand" href="/">shiimanblog</a>
    <nav class="nav">
      <a href="/">記事</a>
      <a href="/profile/">プロフィール</a>
      <a href="/contact/">お問い合わせ</a>
      <button id="search-open" type="button" class="icon-btn" aria-label="検索" title="検索">🔍</button>
      <ThemeToggle />
    </nav>
  </div>
</header>
<SearchOverlay />

<style>
  .site-header { border-bottom: 1px solid var(--border); }
  .inner { display: flex; align-items: center; justify-content: space-between; height: 60px; }
  .brand { font-weight: 700; font-size: 1.1rem; color: var(--text); }
  .nav { display: flex; align-items: center; gap: 18px; font-size: 0.9rem; }
  .nav a { color: var(--text-muted); }
  .nav a:hover { color: var(--accent); text-decoration: none; }
  .icon-btn {
    background: none; border: 1px solid var(--border); color: var(--text);
    border-radius: 6px; width: 34px; height: 34px; cursor: pointer;
    font-size: 14px; line-height: 1;
  }
  .icon-btn:hover { background: var(--bg-subtle); }
</style>
```

- [ ] **Step 5: ビルドしてインデックス生成を確認**

Run:
```bash
npm run build && test -f dist/pagefind/pagefind.js && echo "PAGEFIND_OK"
```
Expected: ビルド成功し `PAGEFIND_OK` が出力される（`dist/pagefind/pagefind.js` が生成されている）。

- [ ] **Step 6: preview で検索動作を手動確認**

Run:
```bash
npm run preview
```
Expected: ブラウザでヘッダーの🔍をクリック→オーバーレイ表示→キーワード入力で記事がヒット→Escで閉じる。ライト/ダーク両方で表示崩れなし。確認後 Ctrl+C で停止。

- [ ] **Step 7: 既存テストが緑であることを確認**

Run:
```bash
npm test
```
Expected: 全テストパス（既存43件）。

- [ ] **Step 8: Commit**

```bash
git add package.json package-lock.json src/components/SearchOverlay.astro src/components/Header.astro
git commit -m "feat: Pagefindによるサイト内検索(ヘッダーオーバーレイ)を追加" -m "Co-Authored-By: Claude Opus 4.8 (1M context) <noreply@anthropic.com>"
```

---

## Task 2: コメント（giscus）

**Files:**
- Create: `src/components/Comments.astro`
- Modify: `src/layouts/PostLayout.astro`（`PrevNext` の後に `Comments` を配置）

**Interfaces:**
- Consumes: 環境変数 `PUBLIC_GISCUS_REPO` / `PUBLIC_GISCUS_REPO_ID` / `PUBLIC_GISCUS_CATEGORY` / `PUBLIC_GISCUS_CATEGORY_ID`（`import.meta.env` 経由、いずれも公開値）。
- Produces: `Comments.astro`（props 無し）。記事ページのみで使用。

**Notes:**
- giscus はリポジトリ `shiiman/blog` の Discussions に紐づく（ユーザー手動準備。値は README に記載のキーで `.env` から注入）。
- 4つの env がすべて設定されている時のみコメント欄を描画（未設定時は何も出さず、ビルドは成功する）。
- 初期テーマは描画時の `document.documentElement.dataset.theme`（`light`/`dark`）を giscus に渡す。テーマ切替は `MutationObserver` で `data-theme` 変化を検知し `postMessage`（`ThemeToggle` は無改修）。
- `define:vars` を使うとスクリプトは `is:inline` 扱い（バンドルされない）になる。

- [ ] **Step 1: Comments.astro を作成**

Create `src/components/Comments.astro`:
```astro
---
// giscus（GitHub Discussions）コメント欄。PUBLIC_GISCUS_* が揃った時のみ描画。
const giscus = {
  repo: import.meta.env.PUBLIC_GISCUS_REPO,
  repoId: import.meta.env.PUBLIC_GISCUS_REPO_ID,
  category: import.meta.env.PUBLIC_GISCUS_CATEGORY,
  categoryId: import.meta.env.PUBLIC_GISCUS_CATEGORY_ID,
}
const enabled = Boolean(giscus.repo && giscus.repoId && giscus.category && giscus.categoryId)
---
{enabled && (
  <section class="comments" aria-label="コメント">
    <h2>コメント</h2>
    <div id="giscus-container"></div>
  </section>
  <script define:vars={{ giscus }}>
    const theme = document.documentElement.dataset.theme === 'dark' ? 'dark' : 'light'
    const s = document.createElement('script')
    s.src = 'https://giscus.app/client.js'
    s.async = true
    s.crossOrigin = 'anonymous'
    s.setAttribute('data-repo', giscus.repo)
    s.setAttribute('data-repo-id', giscus.repoId)
    s.setAttribute('data-category', giscus.category)
    s.setAttribute('data-category-id', giscus.categoryId)
    s.setAttribute('data-mapping', 'pathname')
    s.setAttribute('data-strict', '0')
    s.setAttribute('data-reactions-enabled', '1')
    s.setAttribute('data-emit-metadata', '0')
    s.setAttribute('data-input-position', 'top')
    s.setAttribute('data-theme', theme)
    s.setAttribute('data-lang', 'ja')
    s.setAttribute('data-loading', 'lazy')
    document.getElementById('giscus-container').appendChild(s)

    // テーマ連動: data-theme の変化を giscus iframe に伝える
    const observer = new MutationObserver(() => {
      const t = document.documentElement.dataset.theme === 'dark' ? 'dark' : 'light'
      const iframe = document.querySelector('iframe.giscus-frame')
      if (iframe) {
        iframe.contentWindow.postMessage(
          { giscus: { setConfig: { theme: t } } },
          'https://giscus.app'
        )
      }
    })
    observer.observe(document.documentElement, { attributes: true, attributeFilter: ['data-theme'] })
  </script>
)}

<style>
  .comments { margin-top: 3rem; padding-top: 1.5rem; border-top: 1px solid var(--border); }
  .comments h2 { font-size: 1.2rem; margin: 0 0 1rem; }
</style>
```

- [ ] **Step 2: PostLayout.astro に Comments を組み込む**

`src/layouts/PostLayout.astro` を改修。`import` 行に追加:
```astro
import Comments from '../components/Comments.astro'
```
`<PrevNext prev={prev} next={next} />` の直後（`</article>` の前）に追加:
```astro
      <PrevNext prev={prev} next={next} />
      <Comments />
```

- [ ] **Step 3: ビルドが通ることを確認（env 未設定時はコメント欄が出ない）**

Run:
```bash
npm run build && echo "BUILD_OK"
```
Expected: `BUILD_OK`。env 未設定なら giscus 欄は出力されない（エラーにならない）。

- [ ] **Step 4: env を設定して手動確認（ユーザー準備後）**

`.env` に `PUBLIC_GISCUS_*` 4値を設定し:
```bash
npm run build && npm run preview
```
Expected: 記事ページ末尾に giscus コメント欄が表示され、ヘッダーのテーマ切替（ライト⇄ダーク）に追従する。確認後 Ctrl+C。
（値未取得の段階ではこのステップはスキップ可。Step 3 でビルド成功を確認していれば次へ進む。）

- [ ] **Step 5: 既存テストが緑であることを確認**

Run:
```bash
npm test
```
Expected: 全テストパス。

- [ ] **Step 6: Commit**

```bash
git add src/components/Comments.astro src/layouts/PostLayout.astro
git commit -m "feat: giscusコメント欄を記事ページに追加(テーマ連動)" -m "Co-Authored-By: Claude Opus 4.8 (1M context) <noreply@anthropic.com>"
```

---

## Task 3: 問い合わせ純粋ロジック（src/lib/contact.ts）+ vitest

**Files:**
- Create: `src/lib/contact.ts`
- Test: `src/lib/contact.test.ts`

**Interfaces:**
- Produces:
  - `validateContact(input: ContactInput): ValidationResult`
  - `isTurnstileSuccess(res: TurnstileVerifyResponse | null): boolean`
  - `buildResendPayload(v: ValidatedContact, opts: { from: string; to: string }): ResendPayload`
  - 型: `ContactInput`, `ValidatedContact`, `ValidationResult`, `TurnstileVerifyResponse`, `ResendPayload`
- Consumes: なし（純粋ロジック。`astro:` 系 import は含めない＝Functions から再利用するため）。

**Notes:** Task 4 の `functions/api/contact.ts` がこのモジュールを `../../src/lib/contact` として import する。

- [ ] **Step 1: 失敗するテストを書く**

Create `src/lib/contact.test.ts`:
```ts
import { describe, it, expect } from 'vitest'
import { validateContact, isTurnstileSuccess, buildResendPayload } from './contact'

describe('validateContact', () => {
  const valid = { name: '山田', email: 'a@example.com', subject: '件名', message: '本文', token: 'tok' }

  it('正常な入力を受理し trim する', () => {
    const r = validateContact({ ...valid, name: '  山田  ' })
    expect(r.ok).toBe(true)
    if (r.ok) expect(r.value.name).toBe('山田')
  })

  it('必須項目が欠けると拒否する', () => {
    const r = validateContact({ ...valid, name: '' })
    expect(r.ok).toBe(false)
    if (!r.ok) expect(r.errors).toContain('name')
  })

  it('メール形式が不正なら拒否する', () => {
    const r = validateContact({ ...valid, email: 'invalid' })
    expect(r.ok).toBe(false)
    if (!r.ok) expect(r.errors).toContain('email')
  })

  it('Turnstile トークンが無いと拒否する', () => {
    const r = validateContact({ ...valid, token: '' })
    expect(r.ok).toBe(false)
    if (!r.ok) expect(r.errors).toContain('token')
  })

  it('文字列以外の型は拒否する', () => {
    const r = validateContact({ ...valid, message: 12345 })
    expect(r.ok).toBe(false)
    if (!r.ok) expect(r.errors).toContain('message')
  })

  it('最大長を超えると拒否する', () => {
    const r = validateContact({ ...valid, subject: 'x'.repeat(201) })
    expect(r.ok).toBe(false)
    if (!r.ok) expect(r.errors).toContain('subject')
  })
})

describe('isTurnstileSuccess', () => {
  it('success:true で true', () => {
    expect(isTurnstileSuccess({ success: true })).toBe(true)
  })
  it('success:false で false', () => {
    expect(isTurnstileSuccess({ success: false })).toBe(false)
  })
  it('null で false', () => {
    expect(isTurnstileSuccess(null)).toBe(false)
  })
})

describe('buildResendPayload', () => {
  it('from/to/reply_to/subject/text を整形する', () => {
    const v = { name: '山田', email: 'a@example.com', subject: '件名', message: '本文', token: 'tok' }
    const p = buildResendPayload(v, { from: 'noreply@shiimanblog.com', to: 'me@example.com' })
    expect(p.from).toBe('noreply@shiimanblog.com')
    expect(p.to).toEqual(['me@example.com'])
    expect(p.reply_to).toBe('a@example.com')
    expect(p.subject).toBe('[お問い合わせ] 件名')
    expect(p.text).toContain('山田')
    expect(p.text).toContain('本文')
  })
})
```

- [ ] **Step 2: テストを実行して失敗を確認**

Run:
```bash
npx vitest run src/lib/contact.test.ts
```
Expected: FAIL（`./contact` が見つからない）。

- [ ] **Step 3: contact.ts を実装**

Create `src/lib/contact.ts`:
```ts
// 問い合わせの純粋ロジック（Functions から再利用するため astro: 系 import を含めない）

export interface ContactInput {
  name?: unknown
  email?: unknown
  subject?: unknown
  message?: unknown
  token?: unknown
}

export interface ValidatedContact {
  name: string
  email: string
  subject: string
  message: string
  token: string
}

export type ValidationResult =
  | { ok: true; value: ValidatedContact }
  | { ok: false; errors: string[] }

export interface TurnstileVerifyResponse {
  success: boolean
  [key: string]: unknown
}

export interface ResendPayload {
  from: string
  to: string[]
  reply_to: string
  subject: string
  text: string
}

const EMAIL_RE = /^[^\s@]+@[^\s@]+\.[^\s@]+$/

function asString(v: unknown): string {
  return typeof v === 'string' ? v.trim() : ''
}

export function validateContact(input: ContactInput): ValidationResult {
  const errors: string[] = []
  const name = asString(input.name)
  const email = asString(input.email)
  const subject = asString(input.subject)
  const message = asString(input.message)
  const token = typeof input.token === 'string' ? input.token : ''

  if (!name || name.length > 100) errors.push('name')
  if (!email || email.length > 200 || !EMAIL_RE.test(email)) errors.push('email')
  if (!subject || subject.length > 200) errors.push('subject')
  if (!message || message.length > 5000) errors.push('message')
  if (!token) errors.push('token')

  if (errors.length > 0) return { ok: false, errors }
  return { ok: true, value: { name, email, subject, message, token } }
}

export function isTurnstileSuccess(res: TurnstileVerifyResponse | null): boolean {
  return Boolean(res && res.success === true)
}

export function buildResendPayload(
  v: ValidatedContact,
  opts: { from: string; to: string }
): ResendPayload {
  return {
    from: opts.from,
    to: [opts.to],
    reply_to: v.email,
    subject: `[お問い合わせ] ${v.subject}`,
    text: `お名前: ${v.name}\nメールアドレス: ${v.email}\n\n${v.message}`,
  }
}
```

- [ ] **Step 4: テストを実行して成功を確認**

Run:
```bash
npx vitest run src/lib/contact.test.ts
```
Expected: PASS（全ケース）。

- [ ] **Step 5: 全テストが緑であることを確認**

Run:
```bash
npm test
```
Expected: 全テストパス（既存43件 + 新規）。

- [ ] **Step 6: Commit**

```bash
git add src/lib/contact.ts src/lib/contact.test.ts
git commit -m "feat: 問い合わせの検証/整形ロジックとテストを追加" -m "Co-Authored-By: Claude Opus 4.8 (1M context) <noreply@anthropic.com>"
```

---

## Task 4: 問い合わせ Functions + フォームUI + 結線 + ローカル検証

**Files:**
- Create: `functions/api/contact.ts`
- Create: `src/components/ContactForm.astro`
- Modify: `src/layouts/PageLayout.astro`（本文後に `after-content` named slot を追加）
- Modify: `src/pages/[...slug].astro`（contact ページに `ContactForm` を差し込む）
- Modify: `pages/contact/page.md`（本文を案内文に書き換え）
- Modify: `package.json`（`wrangler` を devDependencies へ、`pages:dev` スクリプト）
- Modify: `.gitignore`（`.dev.vars` を追加）

**Interfaces:**
- Consumes: `src/lib/contact.ts`（`validateContact` / `isTurnstileSuccess` / `buildResendPayload`）、環境変数 `PUBLIC_TURNSTILE_SITE_KEY`（クライアント）、`TURNSTILE_SECRET_KEY` / `RESEND_API_KEY` / `CONTACT_FROM_EMAIL` / `CONTACT_TO_EMAIL`（Functions）。
- Produces: `POST /api/contact`（`{ ok: boolean, message: string }` を返す）。`ContactForm.astro`（props 無し）。

**Notes:**
- `functions/api/contact.ts` は Cloudflare Pages Functions 規約（`onRequestPost` を export）。`PagesFunction` 型は使わず、`env` を自前 interface で型付けして追加の型パッケージを避ける。
- `[...slug].astro` の `getStaticPaths` は変更しない（contact は既に `pages` コレクションから生成される）。描画分岐でのみ `ContactForm` を差し込み、新ルートは作らない。

- [ ] **Step 1: functions/api/contact.ts を作成**

Create `functions/api/contact.ts`:
```ts
import {
  validateContact,
  isTurnstileSuccess,
  buildResendPayload,
  type ContactInput,
  type TurnstileVerifyResponse,
} from '../../src/lib/contact'

interface Env {
  TURNSTILE_SECRET_KEY: string
  RESEND_API_KEY: string
  CONTACT_FROM_EMAIL: string
  CONTACT_TO_EMAIL: string
}

function json(data: unknown, status: number): Response {
  return new Response(JSON.stringify(data), {
    status,
    headers: { 'Content-Type': 'application/json' },
  })
}

export const onRequestPost = async (context: {
  request: Request
  env: Env
}): Promise<Response> => {
  const { request, env } = context

  let body: unknown
  try {
    body = await request.json()
  } catch {
    return json({ ok: false, message: '不正なリクエストです。' }, 400)
  }

  const result = validateContact(body as ContactInput)
  if (!result.ok) {
    return json({ ok: false, message: '入力内容をご確認ください。' }, 400)
  }

  // Turnstile 検証
  const ts = (await fetch(
    'https://challenges.cloudflare.com/turnstile/v0/siteverify',
    {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        secret: env.TURNSTILE_SECRET_KEY,
        response: result.value.token,
      }),
    }
  )
    .then((r) => r.json())
    .catch(() => null)) as TurnstileVerifyResponse | null

  if (!isTurnstileSuccess(ts)) {
    return json({ ok: false, message: '認証に失敗しました。再度お試しください。' }, 403)
  }

  // Resend 送信
  const payload = buildResendPayload(result.value, {
    from: env.CONTACT_FROM_EMAIL,
    to: env.CONTACT_TO_EMAIL,
  })
  const sendRes = await fetch('https://api.resend.com/emails', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${env.RESEND_API_KEY}`,
    },
    body: JSON.stringify(payload),
  })

  if (!sendRes.ok) {
    return json({ ok: false, message: '送信に失敗しました。時間をおいて再度お試しください。' }, 502)
  }

  return json({ ok: true, message: 'お問い合わせを送信しました。ありがとうございます。' }, 200)
}
```

- [ ] **Step 2: ContactForm.astro を作成**

Create `src/components/ContactForm.astro`:
```astro
---
const siteKey = import.meta.env.PUBLIC_TURNSTILE_SITE_KEY ?? ''
---
<form id="contact-form" class="contact-form">
  <label>お名前
    <input name="name" type="text" required maxlength="100" autocomplete="name" />
  </label>
  <label>メールアドレス
    <input name="email" type="email" required maxlength="200" autocomplete="email" />
  </label>
  <label>題名
    <input name="subject" type="text" required maxlength="200" />
  </label>
  <label>メッセージ
    <textarea name="message" required maxlength="5000" rows="6"></textarea>
  </label>
  <div class="cf-turnstile" data-sitekey={siteKey}></div>
  <button type="submit">送信する</button>
  <p id="contact-status" class="status" role="status" aria-live="polite"></p>
</form>

<script src="https://challenges.cloudflare.com/turnstile/v0/api.js" async defer is:inline></script>
<script>
  const form = document.getElementById('contact-form') as HTMLFormElement | null
  const status = document.getElementById('contact-status') as HTMLElement | null

  form?.addEventListener('submit', async (e) => {
    e.preventDefault()
    if (!status) return
    status.textContent = '送信中...'
    const fd = new FormData(form)
    const payload = {
      name: fd.get('name'),
      email: fd.get('email'),
      subject: fd.get('subject'),
      message: fd.get('message'),
      token: fd.get('cf-turnstile-response'),
    }
    try {
      const res = await fetch('/api/contact', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload),
      })
      const data = await res.json()
      status.textContent = data.message ?? (res.ok ? '送信しました。' : '送信に失敗しました。')
      if (res.ok && data.ok) {
        form.reset()
        // @ts-ignore Turnstile グローバル
        window.turnstile?.reset()
      }
    } catch {
      status.textContent = '送信に失敗しました。時間をおいて再度お試しください。'
    }
  })
</script>

<style>
  .contact-form { display: flex; flex-direction: column; gap: 14px; margin-top: 1.5rem; }
  .contact-form label { display: flex; flex-direction: column; gap: 6px; font-size: 0.9rem; color: var(--text-muted); }
  .contact-form input, .contact-form textarea {
    font-family: var(--font-sans); font-size: 1rem; padding: 10px 12px;
    background: var(--bg-subtle); color: var(--text);
    border: 1px solid var(--border); border-radius: 8px;
  }
  .contact-form input:focus, .contact-form textarea:focus { outline: 2px solid var(--accent); }
  .contact-form button {
    align-self: flex-start; background: var(--accent); color: #fff;
    border: none; border-radius: 8px; padding: 10px 20px; font-size: 1rem; cursor: pointer;
  }
  .contact-form button:hover { background: var(--accent-hover); }
  .status { font-size: 0.9rem; color: var(--text-muted); min-height: 1.2em; }
</style>
```

- [ ] **Step 3: PageLayout.astro に after-content スロットを追加**

`src/layouts/PageLayout.astro` の `<div class="content"><slot /></div>` の直後に named slot を追加:
```astro
    <div class="content">
      <slot />
    </div>
    <slot name="after-content" />
  </article>
```

- [ ] **Step 4: [...slug].astro で contact に ContactForm を差し込む**

`src/pages/[...slug].astro` の import に追加:
```astro
import ContactForm from '../components/ContactForm.astro'
```
`PageLayout` の描画部分を改修（`Content` の後に contact 限定で `ContactForm` を `after-content` スロットへ）:
```astro
  <PageLayout entry={props.entry}>
    <Content />
    {props.entry.data.slug === 'contact' && <ContactForm slot="after-content" />}
  </PageLayout>
```

- [ ] **Step 5: pages/contact/page.md の本文を案内文に書き換え**

`pages/contact/page.md` のフロントマター（`---` で囲まれた部分）は維持し、本文（13行目以降）を次に置き換え:
```markdown
お問い合わせは下記フォームよりお送りください。内容を確認のうえ、必要に応じてご入力いただいたメールアドレス宛に返信いたします。返信には数日いただく場合があります。
```

- [ ] **Step 6: wrangler を devDependencies に追加し pages:dev スクリプトを追加**

Run:
```bash
npm install -D wrangler
```
`package.json` の `scripts` に追加:
```json
"pages:dev": "wrangler pages dev dist",
```

- [ ] **Step 7: .gitignore に .dev.vars を追加**

`.gitignore` の「# 環境変数」セクションに追記:
```
.dev.vars
```

- [ ] **Step 8: ビルドが通ることを確認**

Run:
```bash
npm run build && echo "BUILD_OK"
```
Expected: `BUILD_OK`。`/contact/` ページがフォーム付きで生成される（`functions/` は astro build の対象外）。

- [ ] **Step 9: 問い合わせロジックのテストが緑であることを確認**

Run:
```bash
npm test
```
Expected: 全テストパス。

- [ ] **Step 10: wrangler でローカル検証（手動）**

ユーザーが `.dev.vars` に `TURNSTILE_SECRET_KEY` / `RESEND_API_KEY` / `CONTACT_FROM_EMAIL` / `CONTACT_TO_EMAIL`、`.env` に `PUBLIC_TURNSTILE_SITE_KEY` を設定後:
```bash
npm run build && npm run pages:dev
```
Expected: `http://localhost:8788/contact/` でフォーム送信→`functions/api/contact.ts` が検証→Turnstile→Resend のフローを実行し、成功/失敗メッセージが表示される。確認後 Ctrl+C。
（キー未取得の段階では Step 8/9 のビルド・テスト緑をもって完了とし、本ステップは値投入後に実施。）

- [ ] **Step 11: Commit**

```bash
git add functions/api/contact.ts src/components/ContactForm.astro src/layouts/PageLayout.astro src/pages/[...slug].astro pages/contact/page.md package.json package-lock.json .gitignore
git commit -m "feat: 問い合わせフォームとPages Functions(Turnstile+Resend)を追加" -m "Co-Authored-By: Claude Opus 4.8 (1M context) <noreply@anthropic.com>"
```

---

## Task 5: 環境変数ドキュメント（README）と最終統合検証

**Files:**
- Modify: `README.md`（計画3の環境変数セクションを追記）

**Interfaces:**
- Consumes: Task 1〜4 の成果物。
- Produces: README の環境変数表（ユーザーが `.env` / `.dev.vars` / Cloudflare に設定する際の参照）。

- [ ] **Step 1: README.md に環境変数セクションを追記**

`README.md` の適切な位置（末尾、または既存の設定/セットアップ節の後）に次を追記:
```markdown
## 環境変数（計画3: 動的機能）

`.env`（公開値）と `.dev.vars`（秘密値・ローカル wrangler 用）はいずれも **手動作成**します（gitignore 済み）。本番値の Cloudflare 投入は計画4。

### 公開値（`.env`、ビルドに埋め込まれる）
| キー | 用途 |
|---|---|
| `PUBLIC_TURNSTILE_SITE_KEY` | Turnstile ウィジェット（問い合わせ） |
| `PUBLIC_GISCUS_REPO` | giscus 対象リポジトリ（例 `shiiman/blog`） |
| `PUBLIC_GISCUS_REPO_ID` | giscus リポジトリID |
| `PUBLIC_GISCUS_CATEGORY` | giscus カテゴリ名 |
| `PUBLIC_GISCUS_CATEGORY_ID` | giscus カテゴリID |

### 秘密値（`.dev.vars`、Functions のみ参照）
| キー | 用途 |
|---|---|
| `TURNSTILE_SECRET_KEY` | Turnstile サーバー検証 |
| `RESEND_API_KEY` | Resend メール送信 |
| `CONTACT_FROM_EMAIL` | 送信元（例 `noreply@shiimanblog.com`） |
| `CONTACT_TO_EMAIL` | 受信先（問い合わせ通知先） |

### ローカル検証
```bash
npm run build      # astro build + pagefind インデックス生成
npm run preview    # 検索・コメント（giscus）の確認
npm run pages:dev  # 問い合わせ Functions の確認（http://localhost:8788/contact/）
```
```

- [ ] **Step 2: 全体の通し検証**

Run:
```bash
npm test && npm run build && test -f dist/pagefind/pagefind.js && echo "ALL_OK"
```
Expected: テスト全緑 → ビルド成功 → Pagefind インデックス存在 → `ALL_OK`。

- [ ] **Step 3: Commit**

```bash
git add README.md
git commit -m "docs: 計画3の環境変数とローカル検証手順をREADMEに追記" -m "Co-Authored-By: Claude Opus 4.8 (1M context) <noreply@anthropic.com>"
```

---

## 完了条件（設計書 §1 と対応）

1. `npm run build`（astro build + Pagefind インデックス生成）が成功する → Task 1 Step 5 / Task 5 Step 2
2. 検索オーバーレイで記事を全文検索できる → Task 1 Step 6
3. 記事ページに giscus が表示され、テーマ連動する → Task 2 Step 4
4. `/contact/` が実フォーム化され、`wrangler pages dev` で検証→Turnstile→Resend が動く → Task 4 Step 10
5. `npm test`（vitest）が緑（既存43件 + 問い合わせロジック）→ 各タスクの test ステップ

## 計画4へ送るもの（このプランのスコープ外）

- Resend の独自ドメイン DNS 認証、本番シークレットの Cloudflare 投入、giscus 本番疎通
- `public/_redirects`（301、`/feed/`→`/rss.xml` 含む）、Cloudflare Pages デプロイ、DNS 切替、ConoHa 解約
