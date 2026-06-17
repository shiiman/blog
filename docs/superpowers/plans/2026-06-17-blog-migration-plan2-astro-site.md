# 計画2 Astroサイト構築 実装計画

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 計画1で Astro 可読化された記事85件・固定ページ4件を、軽量ミニマル（ライト/ダーク切替）の独自デザインで静的サイト化し、現行 WordPress の URL を保持してローカルで `npm run build` が通る状態にする。

**Architecture:** Astro + Content Collections（glob ローダーで `posts/`・`pages/` をルート据え置き参照）。`output: 'static'` / `trailingSlash: 'always'` / `build.format: 'directory'` で末尾スラッシュ URL を再現。記事 URL の primary カテゴリは稼働中 WordPress の REST から取得した canonical URL マップ（`data/permalinks.json`）で完全一致させる。日付は「JST 壁時計を Z 表記で保持」した値として解釈し、表示は JST 日付・機械可読は正しい UTC へ補正。画像は `astro:assets` でビルド時最適化。

**Tech Stack:** Astro（最新, Content Layer API）, `@astrojs/rss`, `@astrojs/sitemap`, `astro:assets`(sharp), Shiki(内蔵), TypeScript, vitest, gray-matter（既存スクリプト）。

**設計書:** `docs/superpowers/specs/2026-06-17-blog-migration-plan2-astro-site-design.md`

---

## 前提と共通ルール

- すべてのコマンドは**プロジェクトルート** `/Users/a12665/Documents/personal/blog` で実行する。
- ブランチは `feature/blog-migration-plan-2`（作成済み）。各タスク末尾でコミットする。
- 既存の純粋関数は `scripts/lib/*.ts`（vitest 済）。**サイト側のロジックは `src/lib/*.ts`** に新設し、純粋関数は `*.test.ts`（vitest, `scripts/**` ではなく `src/**` も対象に含める）でテストする。
- コミットメッセージは日本語1行（例 `feat: トップページの記事一覧を実装`）。`--no-verify` 禁止。
- 確認済みの権威的事実（現行 WP REST）:
  - 記事 `/{カテゴリslug}/{記事slug}/`、固定ページ `/{slug}/`、カテゴリ `/category/{slug}/`、タグ `/tag/{slug}/`、トップ `/page/2/` 実在。
  - フロントマター `date`/`modified` は**2形式が混在**: ①`2021-08-30T19:30:00.000Z`（80件・JST壁時計をZと誤記。REST の `date`=19:30, `date_gmt`=10:30 と照合済）②`'2026-01-03T00:00:00+09:00'`（5件・正しいJST, 全て00:00:00）。**両形式とも「文字列に書かれた暦時刻＝JST壁時計」**なので、tz指定子を無視して暦フィールドを読みJSTとして扱う（`new Date()` 経由は +09:00 で誤るため使わない）。pages 4件は全て `.000Z`。
  - 日本語カテゴリ/タグ slug は URL エンコード（例 節約=`%e7%af%80%e7%b4%84`）。記事 URL に出る primary カテゴリは全て英語 slug。

---

## ファイル構成（このフェーズで作成/変更）

**作成**
- `astro.config.mjs` — Astro 設定
- `src/env.d.ts` — Astro 型参照
- `src/content.config.ts` — Content Collections（zod スキーマ）
- `src/lib/date.ts` / `src/lib/date.test.ts` — 日付整形（純粋関数）
- `src/lib/taxonomy.ts` / `src/lib/taxonomy.test.ts` — カテゴリ/タグ slug 解決（純粋関数）
- `src/lib/posts.ts` — 公開記事の取得・URL 付与（Astro 依存）
- `src/lib/pagination.ts` / `src/lib/pagination.test.ts` — ページネーション計算（純粋関数）
- `src/styles/global.css` — テーマ変数・共通スタイル
- `src/layouts/BaseLayout.astro` / `PostLayout.astro` / `PageLayout.astro`
- `src/components/Header.astro` / `Footer.astro` / `ThemeToggle.astro` / `PostCard.astro` / `PostMeta.astro` / `TagList.astro` / `Pagination.astro` / `TableOfContents.astro` / `PrevNext.astro`
- `src/pages/index.astro` / `page/[page].astro` / `[...slug].astro` / `category/[category].astro` / `tag/[tag].astro` / `rss.xml.ts` / `404.astro`
- `scripts/export-permalinks.ts` — canonical URL 取得（I/O）
- `scripts/lib/permalinks.ts` / `scripts/lib/permalinks.test.ts` — link→path 変換（純粋関数）
- `scripts/lib/cleanup-tracking.ts` / `scripts/lib/cleanup-tracking.test.ts` — 本文トラッキング画像除去（純粋関数）
- `scripts/cleanup-tracking-images.ts` — 本文クリーニング実行（I/O）
- `public/robots.txt`
- `data/permalinks.json` — 生成物（記事 id→path）

**変更**
- `package.json` — astro 依存と scripts（dev/build/preview/export:permalinks/cleanup:tracking）
- `tsconfig.json` — Astro strict 継承、`src` を include
- `vitest.config.ts` — `src/**/*.test.ts` も対象に
- `data/categories.json` / `data/tags.json` — エンコード slug エントリに `enSlug` 追記
- `scripts/README.md` — 追加スクリプトの記載

---

## Task 1: Astro 導入と基本設定

**Files:**
- Modify: `package.json`
- Create: `astro.config.mjs`, `src/env.d.ts`
- Modify: `tsconfig.json`

- [ ] **Step 1: Astro と integration を導入**

Run:
```bash
npm install astro @astrojs/rss @astrojs/sitemap sharp
```
Expected: `package.json` の dependencies に astro 等が追加される。

- [ ] **Step 2: `package.json` に scripts を追加**

`scripts` を以下に置き換える（既存の移行系 scripts は残す）:
```json
  "scripts": {
    "dev": "astro dev",
    "build": "astro build",
    "preview": "astro preview",
    "astro": "astro",
    "test": "vitest run",
    "export:wp": "tsx scripts/export-wp-data.ts",
    "export:permalinks": "tsx scripts/export-permalinks.ts",
    "localize:images": "tsx scripts/localize-images.ts",
    "migrate:frontmatter": "tsx scripts/migrate-frontmatter.ts",
    "cleanup:tracking": "tsx scripts/cleanup-tracking-images.ts"
  },
```

- [ ] **Step 3: `astro.config.mjs` を作成**

```js
// @ts-check
import { defineConfig } from 'astro/config'
import sitemap from '@astrojs/sitemap'

// https://astro.build/config
export default defineConfig({
  site: 'https://shiimanblog.com',
  trailingSlash: 'always',
  build: { format: 'directory' },
  integrations: [sitemap()],
})
```

- [ ] **Step 4: `src/env.d.ts` を作成**

```ts
/// <reference path="../.astro/types.d.ts" />
/// <reference types="astro/client" />
```

- [ ] **Step 5: `tsconfig.json` を更新**

```json
{
  "extends": "astro/tsconfigs/strict",
  "compilerOptions": {
    "types": ["node"],
    "resolveJsonModule": true
  },
  "include": [".astro/types.d.ts", "src", "scripts"],
  "exclude": ["dist", "node_modules"]
}
```

- [ ] **Step 6: Astro が起動できることを確認**

Run:
```bash
npx astro sync && npx astro --version
```
Expected: バージョンが表示され、`.astro/` 型が生成される（この時点では content.config.ts 未作成なのでコレクション型は空でも可）。

- [ ] **Step 7: コミット**

```bash
git add package.json package-lock.json astro.config.mjs src/env.d.ts tsconfig.json
git commit -m "feat: Astroを導入し基本設定を追加"
```

---

## Task 2: Content Collections スキーマ

**Files:**
- Create: `src/content.config.ts`

- [ ] **Step 1: `src/content.config.ts` を作成**

```ts
import { defineCollection, z } from 'astro:content'
import { glob } from 'astro/loaders'

// YAML は未引用の ISO 日時(...Z)を Date に自動変換するため、文字列/Date の両方を受け、
// 保存値(JST壁時計を Z 表記)と同一の ISO 文字列へ正規化する。
// Date.toISOString() は元の瞬時を同じ文字列に戻すため Task4 の日付ロジックは不変。
const dateString = z
  .union([z.string(), z.date()])
  .transform((v) => (v instanceof Date ? v.toISOString() : v))

const posts = defineCollection({
  loader: glob({ pattern: '*/article.md', base: './posts' }),
  schema: ({ image }) =>
    z.object({
      title: z.string(),
      slug: z.string(),
      date: dateString,
      modified: dateString.optional(),
      excerpt: z.string().optional(),
      categories: z.array(z.string()).default([]),
      tags: z.array(z.string()).default([]),
      eyecatch: image().optional(),
      draft: z.boolean().default(false),
      id: z.number().optional(),
    }),
})

const pages = defineCollection({
  loader: glob({ pattern: '*/page.md', base: './pages' }),
  schema: z.object({
    title: z.string(),
    slug: z.string(),
    date: dateString,
    modified: dateString.optional(),
    excerpt: z.string().optional(),
    draft: z.boolean().default(false),
    id: z.number().optional(),
  }),
})

export const collections = { posts, pages }
```

- [ ] **Step 2: スキーマでコンテンツが読めることを確認**

Run:
```bash
npx astro sync
```
Expected: エラーなく完了（全 85 記事・4 ページがスキーマ検証を通過）。検証エラーが出た場合は該当フロントマターを確認（このフェーズでは内容変更せず、原因を記録して相談）。

- [ ] **Step 3: コミット**

```bash
git add src/content.config.ts
git commit -m "feat: Content Collectionsのスキーマを定義"
```

---

## Task 3: canonical URL マップ（permalinks.json）の取得

**Files:**
- Create: `scripts/lib/permalinks.ts`, `scripts/lib/permalinks.test.ts`, `scripts/export-permalinks.ts`
- Create（生成物）: `data/permalinks.json`

- [ ] **Step 1: 失敗するテストを書く** — `scripts/lib/permalinks.test.ts`

```ts
import { describe, it, expect } from 'vitest'
import { linkToPath, buildPermalinkMap } from './permalinks'

describe('linkToPath', () => {
  it('絶対URLからパス部分(末尾スラッシュ付)を取り出す', () => {
    expect(linkToPath('https://shiimanblog.com/wordpress/conoha-wing/')).toBe('/wordpress/conoha-wing/')
  })
  it('クエリ・ハッシュを除去する', () => {
    expect(linkToPath('https://shiimanblog.com/profile/start-blog/?utm=x#h')).toBe('/profile/start-blog/')
  })
  it('末尾スラッシュが無ければ補う', () => {
    expect(linkToPath('https://shiimanblog.com/contact')).toBe('/contact/')
  })
})

describe('buildPermalinkMap', () => {
  it('id文字列キーで {slug, path} を作る', () => {
    const items = [
      { id: 6, slug: 'start-blog', link: 'https://shiimanblog.com/profile/start-blog/' },
      { id: 34, slug: 'conoha-wing', link: 'https://shiimanblog.com/wordpress/conoha-wing/' },
    ]
    expect(buildPermalinkMap(items)).toEqual({
      '6': { slug: 'start-blog', path: '/profile/start-blog/' },
      '34': { slug: 'conoha-wing', path: '/wordpress/conoha-wing/' },
    })
  })
})
```

- [ ] **Step 2: テストが失敗することを確認**

Run: `npx vitest run scripts/lib/permalinks.test.ts`
Expected: FAIL（`permalinks` モジュール未作成）

- [ ] **Step 3: `scripts/lib/permalinks.ts` を実装**

```ts
export interface PermalinkEntry {
  slug: string
  path: string
}

export type PermalinkMap = Record<string, PermalinkEntry>

interface RestPost {
  id: number
  slug: string
  link: string
}

/** 絶対URLからパス部分を取り出し、末尾スラッシュを保証する（クエリ/ハッシュ除去） */
export function linkToPath(link: string): string {
  const { pathname } = new URL(link)
  return pathname.endsWith('/') ? pathname : `${pathname}/`
}

/** REST の記事配列から id→{slug,path} マップを構築する */
export function buildPermalinkMap(items: RestPost[]): PermalinkMap {
  const map: PermalinkMap = {}
  for (const p of items) {
    map[String(p.id)] = { slug: p.slug, path: linkToPath(p.link) }
  }
  return map
}
```

- [ ] **Step 4: テストが通ることを確認**

Run: `npx vitest run scripts/lib/permalinks.test.ts`
Expected: PASS

- [ ] **Step 5: `scripts/export-permalinks.ts`（I/O）を実装**

```ts
import { writeFile, mkdir } from 'node:fs/promises'
import { buildPermalinkMap } from './lib/permalinks'

const SITE = 'https://shiimanblog.com'
const DATA_DIR = 'data'

interface RestPost {
  id: number
  slug: string
  link: string
}

/** ページングしながら全記事の id/slug/link を取得する */
async function fetchAllPosts(): Promise<RestPost[]> {
  const items: RestPost[] = []
  for (let page = 1; ; page++) {
    const res = await fetch(`${SITE}/wp-json/wp/v2/posts?per_page=100&page=${page}&_fields=id,slug,link&status=publish`)
    if (res.status === 400) break // ページ超過
    if (!res.ok) throw new Error(`取得失敗 posts p${page}: ${res.status}`)
    const batch = (await res.json()) as unknown
    if (!Array.isArray(batch)) throw new Error(`想定外レスポンス posts p${page}`)
    if (batch.length === 0) break
    for (const item of batch) {
      const t = item as Record<string, unknown>
      if (typeof t.id !== 'number' || typeof t.slug !== 'string' || typeof t.link !== 'string') {
        throw new Error(`想定外の記事項目: ${JSON.stringify(item)}`)
      }
    }
    items.push(...(batch as RestPost[]))
    if (batch.length < 100) break
  }
  return items
}

async function main() {
  await mkdir(DATA_DIR, { recursive: true })
  const posts = await fetchAllPosts()
  const map = buildPermalinkMap(posts)
  await writeFile(`${DATA_DIR}/permalinks.json`, JSON.stringify(map, null, 2) + '\n')
  console.log(`permalinks: ${Object.keys(map).length}`)
}

main().catch((e) => {
  console.error(e)
  process.exit(1)
})
```

- [ ] **Step 6: 取得を実行（WP 稼働中・最優先）**

Run: `npm run export:permalinks`
Expected: `permalinks: 85`（前後で一致）。`data/permalinks.json` が生成される。

- [ ] **Step 7: 全記事 id がマップに存在することを確認**

Run:
```bash
node -e "const m=require('./data/permalinks.json');const fs=require('fs');const g=require('glob');" 2>/dev/null; \
grep -h '^id:' posts/*/article.md | sed 's/id: //' | sort -n | while read id; do \
  node -e "const m=require('./data/permalinks.json'); process.exit(m['$id']?0:1)" || echo "MISSING id=$id"; \
done; echo "checked"
```
Expected: `MISSING` が 1 件も出ず `checked` のみ表示。出た場合は WP 側で該当記事が publish か確認。

- [ ] **Step 8: コミット**

```bash
git add scripts/lib/permalinks.ts scripts/lib/permalinks.test.ts scripts/export-permalinks.ts data/permalinks.json
git commit -m "feat: canonical URLマップ(permalinks.json)の取得を追加"
```

---

## Task 4: 日付ユーティリティ

**Files:**
- Create: `src/lib/date.ts`, `src/lib/date.test.ts`
- Modify: `vitest.config.ts`

- [ ] **Step 1: vitest が `src` のテストも拾うようにする** — `vitest.config.ts`

```ts
import { defineConfig } from 'vitest/config'

export default defineConfig({
  test: {
    include: ['scripts/**/*.test.ts', 'src/**/*.test.ts'],
  },
})
```

- [ ] **Step 2: 失敗するテストを書く** — `src/lib/date.test.ts`

```ts
import { describe, it, expect } from 'vitest'
import { formatJaDate, toUtcIso } from './date'

describe('formatJaDate', () => {
  it('Z形式(JST壁時計をZ誤記)を和文日付にする', () => {
    expect(formatJaDate('2021-08-30T19:30:00.000Z')).toBe('2021年8月30日')
  })
  it('+09:00形式(正しいJST)を和文日付にする', () => {
    expect(formatJaDate('2026-01-03T00:00:00+09:00')).toBe('2026年1月3日')
  })
  it('暦フィールドをそのまま読む(深夜でも繰り上がらない)', () => {
    expect(formatJaDate('2022-01-01T00:30:00.000Z')).toBe('2022年1月1日')
  })
  it('想定外の形式はエラー', () => {
    expect(() => formatJaDate('2021/08/30')).toThrow()
  })
})

describe('toUtcIso', () => {
  it('Z形式の暦時刻をJSTとみなし正しいUTC瞬時(-9h)にする', () => {
    expect(toUtcIso('2021-08-30T19:30:00.000Z')).toBe('2021-08-30T10:30:00.000Z')
  })
  it('+09:00形式も同じ規則で正しいUTC瞬時にする(日付跨ぎ)', () => {
    expect(toUtcIso('2026-01-03T00:00:00+09:00')).toBe('2026-01-02T15:00:00.000Z')
  })
})
```

- [ ] **Step 3: テストが失敗することを確認**

Run: `npx vitest run src/lib/date.test.ts`
Expected: FAIL（`date` モジュール未作成）

- [ ] **Step 4: `src/lib/date.ts` を実装**

```ts
// 保存値は2形式が混在する:
//   - "2021-08-30T19:30:00.000Z"  … JST壁時計をZと誤記した値（移行由来, 80件）
//   - "2026-01-03T00:00:00+09:00" … 正しいJST（タイムゾーン付き, 5件）
// どちらも「文字列に書かれた暦時刻」がJST壁時計なので、tz指定子を無視して
// 暦フィールド(年月日時分秒)を読み、JST(=UTC+9)として扱う。
// （new Date() 経由だと +09:00 はオフセット解釈されて暦日がずれるため使わない）
const DATE_RE = /^(\d{4})-(\d{2})-(\d{2})T(\d{2}):(\d{2}):(\d{2})/
const JST_OFFSET_MS = 9 * 60 * 60 * 1000

interface DateParts {
  y: number
  mo: number
  d: number
  h: number
  mi: number
  s: number
}

function parseParts(value: string): DateParts {
  const m = DATE_RE.exec(value)
  if (!m) throw new Error(`想定外の日付形式: ${value}`)
  return { y: +m[1], mo: +m[2], d: +m[3], h: +m[4], mi: +m[5], s: +m[6] }
}

/** JST壁時計の暦日を和文で返す（例「2021年8月30日」） */
export function formatJaDate(value: string): string {
  const { y, mo, d } = parseParts(value)
  return `${y}年${mo}月${d}日`
}

/** JST壁時計を正しいUTC瞬時(-9h)へ補正したISO文字列を返す（RSS/sitemap/<time>用） */
export function toUtcIso(value: string): string {
  const { y, mo, d, h, mi, s } = parseParts(value)
  const utcMs = Date.UTC(y, mo - 1, d, h, mi, s) - JST_OFFSET_MS
  return new Date(utcMs).toISOString()
}
```

- [ ] **Step 5: テストが通ることを確認**

Run: `npx vitest run src/lib/date.test.ts`
Expected: PASS（3 件）

- [ ] **Step 6: コミット**

```bash
git add src/lib/date.ts src/lib/date.test.ts vitest.config.ts
git commit -m "feat: 日付整形ユーティリティ(JST表示/UTC補正)を追加"
```

---

## Task 5: カテゴリ/タグの英語 slug 付与と解決ユーティリティ

**Files:**
- Modify: `data/categories.json`, `data/tags.json`
- Create: `src/lib/taxonomy.ts`, `src/lib/taxonomy.test.ts`

- [ ] **Step 1: `data/categories.json` のエンコード slug に `enSlug` を追記**

対象は id `52`（節約）のみ。エントリを以下に変更:
```json
  "52": {
    "name": "節約",
    "slug": "%e7%af%80%e7%b4%84",
    "enSlug": "savings"
  },
```

- [ ] **Step 2: `data/tags.json` のエンコード slug に `enSlug` を追記**

以下の各エントリへ `enSlug` を追加する（`slug` は現状維持・追記のみ）:

| id | name | enSlug |
|---|---|---|
| 51 | メール | mail |
| 53 | マネーフォワードME | moneyforward-me |
| 54 | キーワードプランナー | keyword-planner |
| 56 | アフェリエイト | affiliate |
| 58 | バリューコマース | valuecommerce |
| 69 | つみたてNISA | tsumitate-nisa |
| 71 | マネリテ学園 | manerite-gakuen |
| 72 | 大河内 薫 | ohkouchi-kaoru |
| 74 | 税金 | tax |
| 86 | Google認証 | google-auth |
| 89 | アクセストレード | accesstrade |
| 90 | 新年 | new-year |
| 91 | 目標設定 | goal-setting |
| 95 | マンガ | manga |
| 96 | 異世界 | isekai |
| 102 | 楽天銀行 | rakuten-bank |
| 103 | 楽天証券 | rakuten-securities |
| 104 | イオン銀行 | aeon-bank |
| 105 | あおぞら銀行 | aozora-bank |
| 131 | 個人事業主 | sole-proprietor |
| 132 | フリーランス | freelance |

例（id 74）:
```json
  "74": {
    "name": "税金",
    "slug": "%e7%a8%8e%e9%87%91",
    "enSlug": "tax"
  },
```

- [ ] **Step 3: JSON が壊れていないことを確認**

Run:
```bash
node -e "require('./data/categories.json');require('./data/tags.json');console.log('ok')"
```
Expected: `ok`

- [ ] **Step 4: 失敗するテストを書く** — `src/lib/taxonomy.test.ts`

```ts
import { describe, it, expect } from 'vitest'
import { resolveTaxonomy, type TaxonomyMap } from './taxonomy'

const map: TaxonomyMap = {
  '14': { name: 'FIRE', slug: 'fire' },
  '52': { name: '節約', slug: '%e7%af%80%e7%b4%84', enSlug: 'savings' },
}

describe('resolveTaxonomy', () => {
  it('英語slugはそのままurlSlugになり、表示名はname', () => {
    expect(resolveTaxonomy('fire', map)).toEqual({ name: 'FIRE', urlSlug: 'fire' })
  })
  it('enSlugがあればurlSlugに使う(表示名は和名)', () => {
    expect(resolveTaxonomy('%e7%af%80%e7%b4%84', map)).toEqual({ name: '節約', urlSlug: 'savings' })
  })
  it('未知slugは防御的にそのまま返す', () => {
    expect(resolveTaxonomy('unknown', map)).toEqual({ name: 'unknown', urlSlug: 'unknown' })
  })
})
```

- [ ] **Step 5: テストが失敗することを確認**

Run: `npx vitest run src/lib/taxonomy.test.ts`
Expected: FAIL

- [ ] **Step 6: `src/lib/taxonomy.ts` を実装**

```ts
import categoriesData from '../../data/categories.json'
import tagsData from '../../data/tags.json'

export interface TaxonomyTerm {
  name: string
  slug: string
  enSlug?: string
}
export type TaxonomyMap = Record<string, TaxonomyTerm>

export interface ResolvedTaxonomy {
  name: string
  urlSlug: string
}

const categories = categoriesData as TaxonomyMap
const tags = tagsData as TaxonomyMap

/** 保存slug(エンコード含む)から表示名とURL用slug(enSlug優先)を解決する */
export function resolveTaxonomy(storedSlug: string, map: TaxonomyMap): ResolvedTaxonomy {
  for (const key of Object.keys(map)) {
    const term = map[key]
    if (term.slug === storedSlug) {
      return { name: term.name, urlSlug: term.enSlug ?? term.slug }
    }
  }
  return { name: storedSlug, urlSlug: storedSlug }
}

export const resolveCategory = (slug: string): ResolvedTaxonomy => resolveTaxonomy(slug, categories)
export const resolveTag = (slug: string): ResolvedTaxonomy => resolveTaxonomy(slug, tags)
```

- [ ] **Step 7: テストが通ることを確認**

Run: `npx vitest run src/lib/taxonomy.test.ts`
Expected: PASS（3 件）

- [ ] **Step 8: コミット**

```bash
git add data/categories.json data/tags.json src/lib/taxonomy.ts src/lib/taxonomy.test.ts
git commit -m "feat: カテゴリ/タグの英語slug付与と解決ユーティリティを追加"
```

---

## Task 6: ページネーション計算ユーティリティ

**Files:**
- Create: `src/lib/pagination.ts`, `src/lib/pagination.test.ts`

- [ ] **Step 1: 失敗するテストを書く** — `src/lib/pagination.test.ts`

```ts
import { describe, it, expect } from 'vitest'
import { totalPages, pageSlice, PAGE_SIZE } from './pagination'

describe('totalPages', () => {
  it('件数とサイズから総ページ数を求める', () => {
    expect(totalPages(85, 10)).toBe(9)
    expect(totalPages(0, 10)).toBe(1)
    expect(totalPages(10, 10)).toBe(1)
    expect(totalPages(11, 10)).toBe(2)
  })
})

describe('pageSlice', () => {
  it('1始まりのページ番号で該当範囲を切り出す', () => {
    const items = Array.from({ length: 25 }, (_, i) => i)
    expect(pageSlice(items, 1, 10)).toEqual(items.slice(0, 10))
    expect(pageSlice(items, 3, 10)).toEqual(items.slice(20, 25))
  })
})

describe('PAGE_SIZE', () => {
  it('WP既定に合わせ10', () => {
    expect(PAGE_SIZE).toBe(10)
  })
})
```

- [ ] **Step 2: テストが失敗することを確認**

Run: `npx vitest run src/lib/pagination.test.ts`
Expected: FAIL

- [ ] **Step 3: `src/lib/pagination.ts` を実装**

```ts
export const PAGE_SIZE = 10

/** 総ページ数（0件でも1ページ） */
export function totalPages(count: number, size: number): number {
  return Math.max(1, Math.ceil(count / size))
}

/** 1始まりのページ番号でスライスを返す */
export function pageSlice<T>(items: T[], page: number, size: number): T[] {
  const start = (page - 1) * size
  return items.slice(start, start + size)
}
```

- [ ] **Step 4: テストが通ることを確認**

Run: `npx vitest run src/lib/pagination.test.ts`
Expected: PASS

- [ ] **Step 5: コミット**

```bash
git add src/lib/pagination.ts src/lib/pagination.test.ts
git commit -m "feat: ページネーション計算ユーティリティを追加"
```

---

## Task 7: 公開記事取得ヘルパー（URL 付与）

**Files:**
- Create: `src/lib/posts.ts`

> このヘルパーは Astro の `getCollection` に依存するため単体 vitest は付けず、後続タスクの `npm run build` で検証する。

- [ ] **Step 1: `src/lib/posts.ts` を実装**

```ts
import { getCollection, type CollectionEntry } from 'astro:content'
import permalinksData from '../../data/permalinks.json'
import { resolveCategory, type ResolvedTaxonomy } from './taxonomy'

interface PermalinkEntry {
  slug: string
  path: string
}
const permalinks = permalinksData as Record<string, PermalinkEntry>

export interface PublishedPost {
  entry: CollectionEntry<'posts'>
  /** 例 /wordpress/conoha-wing/ */
  path: string
  /** 例 wordpress/conoha-wing（[...slug] の params 用） */
  routeParam: string
  primaryCategory: ResolvedTaxonomy
}

/** 公開記事を新しい順で取得し、canonical URL と primary カテゴリを付与する */
export async function getPublishedPosts(): Promise<PublishedPost[]> {
  const posts = await getCollection('posts', (p) => !p.data.draft)
  const enriched: PublishedPost[] = posts.map((entry) => {
    const id = entry.data.id
    const link = id !== undefined ? permalinks[String(id)] : undefined
    if (!link) {
      throw new Error(`permalinks.json に id=${id}（slug=${entry.data.slug}）のパスがありません`)
    }
    const routeParam = link.path.replace(/^\/+|\/+$/g, '')
    const primarySlug = routeParam.split('/')[0]
    return {
      entry,
      path: link.path,
      routeParam,
      primaryCategory: resolveCategory(primarySlug),
    }
  })
  enriched.sort((a, b) => (a.entry.data.date < b.entry.data.date ? 1 : -1))
  return enriched
}
```

- [ ] **Step 2: 型生成が通ることを確認**

Run: `npx astro sync`
Expected: エラーなく完了（`astro:content` の型が生成され `posts.ts` の import が解決可能になる）。本格的な型検証は後続タスクの `npm run build` で行う。

- [ ] **Step 3: コミット**

```bash
git add src/lib/posts.ts
git commit -m "feat: 公開記事取得ヘルパー(URL付与)を追加"
```

---

## Task 8: グローバルスタイルとテーマ（ライト/ダーク）

**Files:**
- Create: `src/styles/global.css`

- [ ] **Step 1: `src/styles/global.css` を作成**

```css
:root {
  --max-width: 720px;
  --font-sans: system-ui, -apple-system, "Segoe UI", "Hiragino Kaku Gothic ProN",
    "Noto Sans JP", Meiryo, sans-serif;
  --font-mono: "SF Mono", "Menlo", "Consolas", monospace;

  --bg: #ffffff;
  --bg-subtle: #f7f8fa;
  --border: #e9eaee;
  --text: #1a1a1a;
  --text-muted: #6b7280;
  --accent: #2563eb;
  --accent-hover: #1d4ed8;
}

[data-theme="dark"] {
  --bg: #0f1115;
  --bg-subtle: #161a21;
  --border: #242a33;
  --text: #e6e8eb;
  --text-muted: #9aa0aa;
  --accent: #60a5fa;
  --accent-hover: #93c5fd;
}

* { box-sizing: border-box; }

html { color-scheme: light dark; }

body {
  margin: 0;
  font-family: var(--font-sans);
  background: var(--bg);
  color: var(--text);
  line-height: 1.75;
  -webkit-font-smoothing: antialiased;
}

a { color: var(--accent); text-decoration: none; }
a:hover { color: var(--accent-hover); text-decoration: underline; }

.container { max-width: var(--max-width); margin: 0 auto; padding: 0 20px; }

img { max-width: 100%; height: auto; }

pre {
  padding: 16px;
  border-radius: 8px;
  overflow-x: auto;
  border: 1px solid var(--border);
}
:not(pre) > code {
  font-family: var(--font-mono);
  font-size: 0.9em;
  background: var(--bg-subtle);
  padding: 0.15em 0.4em;
  border-radius: 4px;
}

h1, h2, h3, h4 { line-height: 1.35; font-weight: 700; }

hr { border: none; border-top: 1px solid var(--border); margin: 2rem 0; }
```

- [ ] **Step 2: コミット**

```bash
git add src/styles/global.css
git commit -m "feat: グローバルスタイルとテーマ変数(ライト/ダーク)を追加"
```

---

## Task 9: BaseLayout・Header・Footer・ThemeToggle

**Files:**
- Create: `src/layouts/BaseLayout.astro`, `src/components/Header.astro`, `src/components/Footer.astro`, `src/components/ThemeToggle.astro`

- [ ] **Step 1: `src/components/ThemeToggle.astro` を作成**

```astro
---
// ライト/ダーク/システム追従を切り替えるボタン
---
<button id="theme-toggle" type="button" aria-label="テーマ切替" title="テーマ切替">
  <span aria-hidden="true">◐</span>
</button>

<style>
  #theme-toggle {
    background: none;
    border: 1px solid var(--border);
    color: var(--text);
    border-radius: 6px;
    width: 34px;
    height: 34px;
    cursor: pointer;
    font-size: 16px;
    line-height: 1;
  }
  #theme-toggle:hover { background: var(--bg-subtle); }
</style>

<script>
  // 'light' | 'dark' | 'system' の3状態をトグル
  const order = ['system', 'light', 'dark'] as const
  type Mode = (typeof order)[number]

  function apply(mode: Mode) {
    const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches
    const dark = mode === 'dark' || (mode === 'system' && prefersDark)
    document.documentElement.dataset.theme = dark ? 'dark' : 'light'
  }

  const btn = document.getElementById('theme-toggle')!
  btn.addEventListener('click', () => {
    const current = (localStorage.getItem('theme') as Mode) ?? 'system'
    const next = order[(order.indexOf(current) + 1) % order.length]
    localStorage.setItem('theme', next)
    apply(next)
  })
</script>
```

- [ ] **Step 2: `src/components/Header.astro` を作成**

```astro
---
import ThemeToggle from './ThemeToggle.astro'
---
<header class="site-header">
  <div class="container inner">
    <a class="brand" href="/">shiimanblog</a>
    <nav class="nav">
      <a href="/">記事</a>
      <a href="/profile/">プロフィール</a>
      <a href="/contact/">お問い合わせ</a>
      <ThemeToggle />
    </nav>
  </div>
</header>

<style>
  .site-header { border-bottom: 1px solid var(--border); }
  .inner { display: flex; align-items: center; justify-content: space-between; height: 60px; }
  .brand { font-weight: 700; font-size: 1.1rem; color: var(--text); }
  .nav { display: flex; align-items: center; gap: 18px; font-size: 0.9rem; }
  .nav a { color: var(--text-muted); }
  .nav a:hover { color: var(--accent); text-decoration: none; }
</style>
```

- [ ] **Step 3: `src/components/Footer.astro` を作成**

```astro
---
const year = new Date().getUTCFullYear()
---
<footer class="site-footer">
  <div class="container">
    <p>&copy; {year} shiimanblog</p>
  </div>
</footer>

<style>
  .site-footer {
    border-top: 1px solid var(--border);
    margin-top: 4rem;
    padding: 2rem 0;
    color: var(--text-muted);
    font-size: 0.85rem;
  }
</style>
```

- [ ] **Step 4: `src/layouts/BaseLayout.astro` を作成**

```astro
---
import '../styles/global.css'
import Header from '../components/Header.astro'
import Footer from '../components/Footer.astro'

interface Props {
  title: string
  description?: string
  ogImage?: string
  ogType?: 'website' | 'article'
}
const { title, description, ogImage, ogType = 'website' } = Astro.props
const canonical = new URL(Astro.url.pathname, Astro.site).href
const absoluteOg = ogImage ? new URL(ogImage, Astro.site).href : undefined
---
<!doctype html>
<html lang="ja">
  <head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <title>{title}</title>
    {description && <meta name="description" content={description} />}
    <link rel="canonical" href={canonical} />
    <meta property="og:type" content={ogType} />
    <meta property="og:title" content={title} />
    {description && <meta property="og:description" content={description} />}
    <meta property="og:url" content={canonical} />
    {absoluteOg && <meta property="og:image" content={absoluteOg} />}
    <meta name="twitter:card" content={absoluteOg ? 'summary_large_image' : 'summary'} />
    <link rel="alternate" type="application/rss+xml" title="shiimanblog" href="/rss.xml" />
    {/* FOUC回避: 保存テーマ/システム設定を即時適用 */}
    <script is:inline>
      (() => {
        const mode = localStorage.getItem('theme') ?? 'system'
        const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches
        const dark = mode === 'dark' || (mode === 'system' && prefersDark)
        document.documentElement.dataset.theme = dark ? 'dark' : 'light'
      })()
    </script>
  </head>
  <body>
    <Header />
    <main class="container">
      <slot />
    </main>
    <Footer />
  </body>
</html>
```

- [ ] **Step 5: コミット**

```bash
git add src/layouts/BaseLayout.astro src/components/Header.astro src/components/Footer.astro src/components/ThemeToggle.astro
git commit -m "feat: BaseLayoutとHeader/Footer/テーマ切替を追加"
```

---

## Task 10: PostCard・PostMeta・TagList・Pagination コンポーネント

**Files:**
- Create: `src/components/PostMeta.astro`, `src/components/TagList.astro`, `src/components/PostCard.astro`, `src/components/Pagination.astro`

- [ ] **Step 1: `src/components/PostMeta.astro` を作成**

```astro
---
import { formatJaDate, toUtcIso } from '../lib/date'
import type { ResolvedTaxonomy } from '../lib/taxonomy'

interface Props {
  category: ResolvedTaxonomy
  date: string
}
const { category, date } = Astro.props
---
<div class="post-meta">
  <a class="cat" href={`/category/${category.urlSlug}/`}>{category.name}</a>
  <time datetime={toUtcIso(date)}>{formatJaDate(date)}</time>
</div>

<style>
  .post-meta { display: flex; gap: 12px; align-items: center; font-size: 0.8rem; }
  .cat { color: var(--accent); font-weight: 600; text-transform: uppercase; letter-spacing: 0.02em; }
  time { color: var(--text-muted); }
</style>
```

- [ ] **Step 2: `src/components/TagList.astro` を作成**

```astro
---
import { resolveTag } from '../lib/taxonomy'

interface Props {
  tags: string[]
}
const { tags } = Astro.props
const resolved = tags.map((t) => resolveTag(t))
---
{resolved.length > 0 && (
  <ul class="tag-list">
    {resolved.map((t) => (
      <li><a href={`/tag/${t.urlSlug}/`}>#{t.name}</a></li>
    ))}
  </ul>
)}

<style>
  .tag-list { list-style: none; padding: 0; display: flex; flex-wrap: wrap; gap: 8px; }
  .tag-list a {
    font-size: 0.8rem;
    color: var(--text-muted);
    border: 1px solid var(--border);
    border-radius: 999px;
    padding: 2px 10px;
  }
  .tag-list a:hover { color: var(--accent); border-color: var(--accent); text-decoration: none; }
</style>
```

- [ ] **Step 3: `src/components/PostCard.astro` を作成（リスト+左サムネ）**

```astro
---
import { Image } from 'astro:assets'
import PostMeta from './PostMeta.astro'
import type { PublishedPost } from '../lib/posts'

interface Props {
  post: PublishedPost
}
const { post } = Astro.props
const { entry, path, primaryCategory } = post
const eyecatch = entry.data.eyecatch
---
<article class="card">
  {eyecatch && (
    <a href={path} class="thumb">
      <Image src={eyecatch} alt={entry.data.title} width={120} height={84} />
    </a>
  )}
  <div class="body">
    <PostMeta category={primaryCategory} date={entry.data.date} />
    <h2 class="title"><a href={path}>{entry.data.title}</a></h2>
  </div>
</article>

<style>
  .card { display: flex; gap: 14px; padding: 16px 0; border-bottom: 1px solid var(--border); }
  .thumb { flex: none; }
  .thumb :global(img) { width: 120px; height: 84px; object-fit: cover; border-radius: 6px; }
  .body { display: flex; flex-direction: column; gap: 6px; }
  .title { font-size: 1.05rem; margin: 0; }
  .title a { color: var(--text); }
  .title a:hover { color: var(--accent); text-decoration: none; }
</style>
```

- [ ] **Step 4: `src/components/Pagination.astro` を作成**

```astro
---
interface Props {
  current: number
  total: number
}
const { current, total } = Astro.props
const href = (page: number) => (page === 1 ? '/' : `/page/${page}/`)
const prev = current > 1 ? current - 1 : null
const next = current < total ? current + 1 : null
---
{total > 1 && (
  <nav class="pagination">
    {prev ? <a href={href(prev)}>← 前へ</a> : <span class="disabled">← 前へ</span>}
    <span class="status">{current} / {total}</span>
    {next ? <a href={href(next)}>次へ →</a> : <span class="disabled">次へ →</span>}
  </nav>
)}

<style>
  .pagination { display: flex; align-items: center; justify-content: space-between; margin-top: 2rem; font-size: 0.9rem; }
  .disabled { color: var(--text-muted); opacity: 0.5; }
  .status { color: var(--text-muted); }
</style>
```

- [ ] **Step 5: コミット**

```bash
git add src/components/PostMeta.astro src/components/TagList.astro src/components/PostCard.astro src/components/Pagination.astro
git commit -m "feat: 記事カード/メタ/タグ/ページネーションのコンポーネントを追加"
```

---

## Task 11: トップページとページネーション

**Files:**
- Create: `src/pages/index.astro`, `src/pages/page/[page].astro`

- [ ] **Step 1: `src/pages/index.astro` を作成**

```astro
---
import BaseLayout from '../layouts/BaseLayout.astro'
import PostCard from '../components/PostCard.astro'
import Pagination from '../components/Pagination.astro'
import { getPublishedPosts } from '../lib/posts'
import { PAGE_SIZE, pageSlice, totalPages } from '../lib/pagination'

const posts = await getPublishedPosts()
const total = totalPages(posts.length, PAGE_SIZE)
const items = pageSlice(posts, 1, PAGE_SIZE)
---
<BaseLayout title="shiimanblog" description="ゲーム開発エンジニアしーまんの技術・副業・投資ブログ">
  <div class="list">
    {items.map((post) => <PostCard post={post} />)}
  </div>
  <Pagination current={1} total={total} />
</BaseLayout>
```

- [ ] **Step 2: `src/pages/page/[page].astro` を作成（2ページ目以降）**

```astro
---
import BaseLayout from '../../layouts/BaseLayout.astro'
import PostCard from '../../components/PostCard.astro'
import Pagination from '../../components/Pagination.astro'
import { getPublishedPosts } from '../../lib/posts'
import { PAGE_SIZE, pageSlice, totalPages } from '../../lib/pagination'
import type { GetStaticPaths } from 'astro'

export const getStaticPaths = (async () => {
  const posts = await getPublishedPosts()
  const total = totalPages(posts.length, PAGE_SIZE)
  // 1ページ目は / が担当するため 2..total を生成
  const pages = []
  for (let p = 2; p <= total; p++) {
    pages.push({ params: { page: String(p) }, props: { page: p, total } })
  }
  return pages
}) satisfies GetStaticPaths

const { page, total } = Astro.props
const posts = await getPublishedPosts()
const items = pageSlice(posts, page, PAGE_SIZE)
---
<BaseLayout title={`記事一覧 (${page}/${total}) - shiimanblog`}>
  <div class="list">
    {items.map((post) => <PostCard post={post} />)}
  </div>
  <Pagination current={page} total={total} />
</BaseLayout>
```

- [ ] **Step 3: ビルドしてトップとページネーションを確認**

Run: `npm run build`
Expected: 成功。`dist/index.html` と `dist/page/2/index.html` が生成される。

Run:
```bash
ls dist/index.html dist/page/2/index.html && \
grep -c 'class="card"' dist/index.html
```
Expected: 両ファイルが存在し、`index.html` の `.card` が 10。

- [ ] **Step 4: コミット**

```bash
git add src/pages/index.astro src/pages/page/[page].astro
git commit -m "feat: トップページと2ページ目以降のページネーションを実装"
```

---

## Task 12: TableOfContents（右サイドバー追従目次）

**Files:**
- Create: `src/components/TableOfContents.astro`

- [ ] **Step 1: `src/components/TableOfContents.astro` を作成**

```astro
---
import type { MarkdownHeading } from 'astro'

interface Props {
  headings: MarkdownHeading[]
}
// h2/h3 のみ目次に採用
const items = Astro.props.headings.filter((h) => h.depth === 2 || h.depth === 3)
---
{items.length > 0 && (
  <nav class="toc" aria-label="目次">
    <p class="toc-title">目次</p>
    <ul>
      {items.map((h) => (
        <li class={`depth-${h.depth}`}><a href={`#${h.slug}`} data-toc={h.slug}>{h.text}</a></li>
      ))}
    </ul>
  </nav>
)}

<style>
  .toc { font-size: 0.82rem; }
  .toc-title { font-weight: 700; color: var(--text); margin: 0 0 8px; }
  .toc ul { list-style: none; margin: 0; padding: 0; border-left: 1px solid var(--border); }
  .toc li a { display: block; color: var(--text-muted); padding: 3px 0 3px 12px; margin-left: -1px; border-left: 2px solid transparent; }
  .toc li.depth-3 a { padding-left: 24px; }
  .toc li a:hover { color: var(--accent); text-decoration: none; }
  .toc li a.active { color: var(--accent); border-left-color: var(--accent); }
</style>

<script>
  // 現在地ハイライト（IntersectionObserver）
  const links = new Map<string, HTMLAnchorElement>()
  document.querySelectorAll<HTMLAnchorElement>('.toc a[data-toc]').forEach((a) => {
    links.set(a.dataset.toc!, a)
  })
  const targets = [...links.keys()]
    .map((id) => document.getElementById(id))
    .filter((el): el is HTMLElement => el !== null)

  const observer = new IntersectionObserver(
    (entries) => {
      for (const e of entries) {
        if (e.isIntersecting) {
          links.forEach((a) => a.classList.remove('active'))
          links.get(e.target.id)?.classList.add('active')
        }
      }
    },
    { rootMargin: '0px 0px -70% 0px', threshold: 0 },
  )
  targets.forEach((t) => observer.observe(t))
</script>
```

- [ ] **Step 2: コミット**

```bash
git add src/components/TableOfContents.astro
git commit -m "feat: 右サイドバー追従目次コンポーネントを追加"
```

---

## Task 13: PrevNext（前後記事）

**Files:**
- Create: `src/components/PrevNext.astro`

- [ ] **Step 1: `src/components/PrevNext.astro` を作成**

```astro
---
interface Link { path: string; title: string }
interface Props {
  prev: Link | null
  next: Link | null
}
const { prev, next } = Astro.props
---
<nav class="prevnext">
  <div class="col">
    {prev && (
      <a href={prev.path}>
        <span class="label">← 前の記事</span>
        <span class="t">{prev.title}</span>
      </a>
    )}
  </div>
  <div class="col right">
    {next && (
      <a href={next.path}>
        <span class="label">次の記事 →</span>
        <span class="t">{next.title}</span>
      </a>
    )}
  </div>
</nav>

<style>
  .prevnext { display: flex; gap: 16px; margin-top: 3rem; }
  .col { flex: 1; }
  .col.right { text-align: right; }
  .label { display: block; font-size: 0.75rem; color: var(--text-muted); margin-bottom: 4px; }
  .t { display: block; font-size: 0.9rem; font-weight: 600; }
</style>
```

> 前後関係は「公開記事を新しい順に並べた配列」での隣接で定義する（記事詳細タスクで配線）。

- [ ] **Step 2: コミット**

```bash
git add src/components/PrevNext.astro
git commit -m "feat: 前後記事ナビゲーションのコンポーネントを追加"
```

---

## Task 14: 記事/固定ページ詳細とレイアウト

**Files:**
- Create: `src/layouts/PostLayout.astro`, `src/layouts/PageLayout.astro`, `src/pages/[...slug].astro`

- [ ] **Step 1: `src/layouts/PostLayout.astro` を作成**

```astro
---
import { Image } from 'astro:assets'
import BaseLayout from './BaseLayout.astro'
import PostMeta from '../components/PostMeta.astro'
import TagList from '../components/TagList.astro'
import TableOfContents from '../components/TableOfContents.astro'
import PrevNext from '../components/PrevNext.astro'
import type { MarkdownHeading } from 'astro'
import type { CollectionEntry } from 'astro:content'
import type { ResolvedTaxonomy } from '../lib/taxonomy'

interface Link { path: string; title: string }
interface Props {
  entry: CollectionEntry<'posts'>
  primaryCategory: ResolvedTaxonomy
  headings: MarkdownHeading[]
  prev: Link | null
  next: Link | null
}
const { entry, primaryCategory, headings, prev, next } = Astro.props
const { title, excerpt, eyecatch, date, tags } = entry.data
const ogImage = eyecatch?.src
---
<BaseLayout title={`${title} - shiimanblog`} description={excerpt} ogImage={ogImage} ogType="article">
  <div class="post-grid">
    <article class="post">
      {eyecatch && <Image class="hero" src={eyecatch} alt={title} width={720} />}
      <h1>{title}</h1>
      <PostMeta category={primaryCategory} date={date} />
      <div class="content">
        <slot />
      </div>
      <TagList tags={tags} />
      <PrevNext prev={prev} next={next} />
    </article>
    <aside class="sidebar">
      <div class="toc-sticky">
        <TableOfContents headings={headings} />
      </div>
    </aside>
  </div>
</BaseLayout>

<style>
  .post-grid { display: grid; grid-template-columns: minmax(0, 1fr) 220px; gap: 40px; align-items: start; }
  .post { min-width: 0; }
  .hero { width: 100%; height: auto; border-radius: 10px; margin-bottom: 1.5rem; }
  .post h1 { font-size: 1.8rem; margin: 0 0 0.6rem; }
  .content { margin-top: 1.5rem; }
  .content :global(h2) { margin-top: 2.2rem; border-bottom: 1px solid var(--border); padding-bottom: 0.3rem; }
  .content :global(h3) { margin-top: 1.6rem; }
  .content :global(img) { border-radius: 8px; }
  .sidebar { position: relative; }
  .toc-sticky { position: sticky; top: 24px; }
  @media (max-width: 900px) {
    .post-grid { grid-template-columns: 1fr; }
    .sidebar { display: none; }
  }
</style>
```

> スマホ折りたたみは「900px 以下で右 TOC を非表示」とする（最小実装）。本文先頭のインライン目次が必要なら後続フェーズで検討。

- [ ] **Step 2: `src/layouts/PageLayout.astro` を作成**

```astro
---
import BaseLayout from './BaseLayout.astro'
import type { CollectionEntry } from 'astro:content'

interface Props {
  entry: CollectionEntry<'pages'>
}
const { entry } = Astro.props
const { title, excerpt } = entry.data
---
<BaseLayout title={`${title} - shiimanblog`} description={excerpt}>
  <article class="page">
    <h1>{title}</h1>
    <div class="content">
      <slot />
    </div>
  </article>
</BaseLayout>

<style>
  .page h1 { font-size: 1.8rem; margin: 0 0 1.5rem; }
  .content :global(h2) { margin-top: 2rem; }
  .content :global(img) { border-radius: 8px; }
</style>
```

- [ ] **Step 3: `src/pages/[...slug].astro` を作成（記事＋固定ページ）**

```astro
---
import { getCollection, render } from 'astro:content'
import PostLayout from '../layouts/PostLayout.astro'
import PageLayout from '../layouts/PageLayout.astro'
import { getPublishedPosts } from '../lib/posts'
import type { GetStaticPaths } from 'astro'

const RESERVED = new Set(['category', 'tag', 'page', 'rss.xml', '404'])

export const getStaticPaths = (async () => {
  const posts = await getPublishedPosts()

  // 記事: permalinks の routeParam を slug に。前後記事も配列順で確定
  const postPaths = posts.map((post, i) => ({
    params: { slug: post.routeParam },
    props: {
      kind: 'post' as const,
      entry: post.entry,
      primaryCategory: post.primaryCategory,
      prev: i > 0 ? { path: posts[i - 1].path, title: posts[i - 1].entry.data.title } : null,
      next: i < posts.length - 1 ? { path: posts[i + 1].path, title: posts[i + 1].entry.data.title } : null,
    },
  }))

  // 固定ページ: /{slug}/。予約語と衝突したらビルド失敗
  const pages = await getCollection('pages', (p) => !p.data.draft)
  const pagePaths = pages.map((entry) => {
    if (RESERVED.has(entry.data.slug)) {
      throw new Error(`固定ページの slug が予約語と衝突: ${entry.data.slug}`)
    }
    return { params: { slug: entry.data.slug }, props: { kind: 'page' as const, entry } }
  })

  return [...postPaths, ...pagePaths]
}) satisfies GetStaticPaths

const props = Astro.props
const { Content, headings } = await render(props.entry)
---
{props.kind === 'post' ? (
  <PostLayout
    entry={props.entry}
    primaryCategory={props.primaryCategory}
    headings={headings}
    prev={props.prev}
    next={props.next}
  >
    <Content />
  </PostLayout>
) : (
  <PageLayout entry={props.entry}>
    <Content />
  </PageLayout>
)}
```

- [ ] **Step 4: ビルドして記事・固定ページの URL を確認**

Run: `npm run build`
Expected: 成功。

Run:
```bash
test -f dist/profile/start-blog/index.html && echo "post-url OK"
test -f dist/wordpress/conoha-wing/index.html && echo "multi-cat-url OK"
test -f dist/contact/index.html && echo "page-url OK"
grep -q 'class="toc"' dist/profile/start-blog/index.html && echo "toc OK"
```
Expected: `post-url OK` / `multi-cat-url OK` / `page-url OK` / `toc OK` がすべて表示。

- [ ] **Step 5: コミット**

```bash
git add src/layouts/PostLayout.astro src/layouts/PageLayout.astro "src/pages/[...slug].astro"
git commit -m "feat: 記事/固定ページ詳細とレイアウト(右TOC/前後記事)を実装"
```

---

## Task 15: カテゴリ/タグ アーカイブ

**Files:**
- Create: `src/pages/category/[category].astro`, `src/pages/tag/[tag].astro`

- [ ] **Step 1: `src/pages/category/[category].astro` を作成**

```astro
---
import BaseLayout from '../../layouts/BaseLayout.astro'
import PostCard from '../../components/PostCard.astro'
import { getPublishedPosts, type PublishedPost } from '../../lib/posts'
import { resolveCategory } from '../../lib/taxonomy'
import type { GetStaticPaths } from 'astro'

export const getStaticPaths = (async () => {
  const posts = await getPublishedPosts()
  // urlSlug でグルーピング（記事の全カテゴリを対象）
  const groups = new Map<string, { name: string; posts: PublishedPost[] }>()
  for (const post of posts) {
    for (const stored of post.entry.data.categories) {
      const { name, urlSlug } = resolveCategory(stored)
      if (!groups.has(urlSlug)) groups.set(urlSlug, { name, posts: [] })
      groups.get(urlSlug)!.posts.push(post)
    }
  }
  return [...groups.entries()].map(([urlSlug, g]) => ({
    params: { category: urlSlug },
    props: { name: g.name, posts: g.posts },
  }))
}) satisfies GetStaticPaths

const { name, posts } = Astro.props
---
<BaseLayout title={`カテゴリ: ${name} - shiimanblog`} description={`${name}の記事一覧`}>
  <h1 class="archive-title">カテゴリ: {name}</h1>
  <div class="list">
    {posts.map((post) => <PostCard post={post} />)}
  </div>
</BaseLayout>

<style>
  .archive-title { font-size: 1.5rem; margin: 1rem 0 1.5rem; }
</style>
```

- [ ] **Step 2: `src/pages/tag/[tag].astro` を作成**

```astro
---
import BaseLayout from '../../layouts/BaseLayout.astro'
import PostCard from '../../components/PostCard.astro'
import { getPublishedPosts, type PublishedPost } from '../../lib/posts'
import { resolveTag } from '../../lib/taxonomy'
import type { GetStaticPaths } from 'astro'

export const getStaticPaths = (async () => {
  const posts = await getPublishedPosts()
  const groups = new Map<string, { name: string; posts: PublishedPost[] }>()
  for (const post of posts) {
    for (const stored of post.entry.data.tags) {
      const { name, urlSlug } = resolveTag(stored)
      if (!groups.has(urlSlug)) groups.set(urlSlug, { name, posts: [] })
      groups.get(urlSlug)!.posts.push(post)
    }
  }
  return [...groups.entries()].map(([urlSlug, g]) => ({
    params: { tag: urlSlug },
    props: { name: g.name, posts: g.posts },
  }))
}) satisfies GetStaticPaths

const { name, posts } = Astro.props
---
<BaseLayout title={`タグ: ${name} - shiimanblog`} description={`${name}の記事一覧`}>
  <h1 class="archive-title">タグ: {name}</h1>
  <div class="list">
    {posts.map((post) => <PostCard post={post} />)}
  </div>
</BaseLayout>

<style>
  .archive-title { font-size: 1.5rem; margin: 1rem 0 1.5rem; }
</style>
```

- [ ] **Step 3: ビルドしてアーカイブ URL を確認**

Run: `npm run build`
Expected: 成功。

Run:
```bash
test -f dist/category/wordpress/index.html && echo "cat OK"
test -f dist/category/savings/index.html && echo "encoded-cat-enSlug OK"
test -f dist/tag/tax/index.html && echo "tag-enSlug OK"
test -f dist/tag/aws/index.html && echo "tag OK"
```
Expected: `cat OK` / `encoded-cat-enSlug OK` / `tag-enSlug OK` / `tag OK` がすべて表示。

- [ ] **Step 4: コミット**

```bash
git add "src/pages/category/[category].astro" "src/pages/tag/[tag].astro"
git commit -m "feat: カテゴリ/タグのアーカイブページを実装"
```

---

## Task 16: RSS フィード

**Files:**
- Create: `src/pages/rss.xml.ts`

- [ ] **Step 1: `src/pages/rss.xml.ts` を作成**

```ts
import rss from '@astrojs/rss'
import type { APIContext } from 'astro'
import { getPublishedPosts } from '../lib/posts'
import { toUtcIso } from '../lib/date'

export async function GET(context: APIContext) {
  const posts = await getPublishedPosts()
  return rss({
    title: 'shiimanblog',
    description: 'ゲーム開発エンジニアしーまんの技術・副業・投資ブログ',
    site: context.site!,
    items: posts.map((p) => ({
      title: p.entry.data.title,
      link: p.path,
      pubDate: new Date(toUtcIso(p.entry.data.date)),
      description: p.entry.data.excerpt ?? '',
    })),
  })
}
```

- [ ] **Step 2: ビルドして RSS を確認**

Run: `npm run build`
Expected: 成功。

Run:
```bash
test -f dist/rss.xml && echo "rss OK" && grep -c '<item>' dist/rss.xml
```
Expected: `rss OK` と `85`。

- [ ] **Step 3: コミット**

```bash
git add src/pages/rss.xml.ts
git commit -m "feat: RSSフィード(/rss.xml)を実装"
```

---

## Task 17: サイトマップ・robots.txt・404

**Files:**
- Create: `public/robots.txt`, `src/pages/404.astro`
- （`@astrojs/sitemap` は Task 1 で導入済）

- [ ] **Step 1: `public/robots.txt` を作成**

```text
User-agent: *
Allow: /

Sitemap: https://shiimanblog.com/sitemap-index.xml
```

- [ ] **Step 2: `src/pages/404.astro` を作成**

```astro
---
import BaseLayout from '../layouts/BaseLayout.astro'
---
<BaseLayout title="ページが見つかりません - shiimanblog">
  <div class="notfound">
    <h1>404</h1>
    <p>お探しのページは見つかりませんでした。</p>
    <p><a href="/">トップへ戻る</a></p>
  </div>
</BaseLayout>

<style>
  .notfound { text-align: center; padding: 4rem 0; }
  .notfound h1 { font-size: 3rem; margin: 0 0 0.5rem; }
</style>
```

- [ ] **Step 3: ビルドして確認**

Run: `npm run build`
Expected: 成功。

Run:
```bash
test -f dist/404.html && echo "404 OK"
test -f dist/sitemap-index.xml && echo "sitemap OK"
test -f dist/robots.txt && echo "robots OK"
grep -q 'mahjong\|poker' dist/sitemap-0.xml 2>/dev/null && echo "WARN: future-sites in sitemap" || echo "future-sites excluded OK"
```
Expected: `404 OK` / `sitemap OK` / `robots OK` / `future-sites excluded OK`。

- [ ] **Step 4: コミット**

```bash
git add public/robots.txt src/pages/404.astro
git commit -m "feat: サイトマップ/robots.txt/404ページを追加"
```

---

## Task 18: 本文トラッキング画像のクリーニング

**Files:**
- Create: `scripts/lib/cleanup-tracking.ts`, `scripts/lib/cleanup-tracking.test.ts`, `scripts/cleanup-tracking-images.ts`
- Modify（生成）: `posts/*/article.md`（該当13ファイル）

> 対象は不可視のトラッキング画像のみ（a8.net / valuecommerce gifbanner / accesstrade rr）。リンクカードのファビコン・OGP 画像など可視要素は除去しない。

- [ ] **Step 1: 失敗するテストを書く** — `scripts/lib/cleanup-tracking.test.ts`

```ts
import { describe, it, expect } from 'vitest'
import { stripTrackingImages } from './cleanup-tracking'

describe('stripTrackingImages', () => {
  it('a8.netの計測gif画像記法を除去する', () => {
    const input = 'before\n![](https://px.a8.net/svt/ejp?a8mat=ABC)\nafter'
    expect(stripTrackingImages(input)).toBe('before\nafter')
  })
  it('valuecommerceのgifbannerを除去する', () => {
    const input = '![](https://ad.jp.ap.valuecommerce.com/servlet/gifbanner?sid=1&pid=2)\n本文'
    expect(stripTrackingImages(input)).toBe('本文')
  })
  it('accesstradeのrrを除去する', () => {
    const input = 'x ![](https://h.accesstrade.net/sp/rr?rk=01001aqe00lqea) y'
    expect(stripTrackingImages(input)).toBe('x  y')
  })
  it('プロトコル相対URLのトラッキング画像も除去しリンクは保持する', () => {
    const input = '[![](//ad.jp.ap.valuecommerce.com/servlet/gifbanner?sid=1&pid=2)バリューコマース](//ck.jp.ap.valuecommerce.com/servlet/referral?sid=1&pid=2)'
    expect(stripTrackingImages(input)).toBe('[バリューコマース](//ck.jp.ap.valuecommerce.com/servlet/referral?sid=1&pid=2)')
  })
  it('通常の画像やリンクは保持する', () => {
    const input = '![alt](./assets/foo.png)\n[a8リンク](https://px.a8.net/abc)'
    expect(stripTrackingImages(input)).toBe(input)
  })
})
```

- [ ] **Step 2: テストが失敗することを確認**

Run: `npx vitest run scripts/lib/cleanup-tracking.test.ts`
Expected: FAIL

- [ ] **Step 3: `scripts/lib/cleanup-tracking.ts` を実装**

```ts
// 不可視トラッキング/計測ピクセル画像のホスト（画像記法 ![]() のみ対象。
// リンク記法 [text](...) は対象外なのでアフィリエイトリンクは保持される）。
// valuecommerce は ad=計測gif / ck=referral計測pixel の双方が画像で混入する。
const TRACKING_HOSTS = [
  'px.a8.net',
  'ad.jp.ap.valuecommerce.com',
  'ck.jp.ap.valuecommerce.com',
  'h.accesstrade.net',
]

/** 本文から不可視トラッキング画像の Markdown 記法だけを除去する（http/https/プロトコル相対の全てに対応） */
export function stripTrackingImages(body: string): string {
  let out = body
  for (const host of TRACKING_HOSTS) {
    // ![...]((https?:)?//<host>...) を、直前の改行ごと（無ければ単体で）除去
    const escaped = host.replace(/[.]/g, '\\.')
    const re = new RegExp(`\\n?!\\[[^\\]]*\\]\\((?:https?:)?\\/\\/${escaped}[^)]*\\)`, 'g')
    out = out.replace(re, '')
  }
  return out
}
```

- [ ] **Step 4: テストが通ることを確認**

Run: `npx vitest run scripts/lib/cleanup-tracking.test.ts`
Expected: PASS（4 件）

> 注: テスト「a8.netの計測gif〜」は入力 `before\n![...]\nafter` に対し `\n![...]` を除去して `before\nafter` を期待する。実装の正規表現は先頭の改行ごと除去するため一致する。

- [ ] **Step 5: `scripts/cleanup-tracking-images.ts`（I/O）を実装**

```ts
import { readFile, writeFile } from 'node:fs/promises'
import { listContentDirs } from './lib/content-roots'
import { stripTrackingImages } from './lib/cleanup-tracking'

async function main() {
  const dirs = await listContentDirs()
  let changed = 0
  for (const { dir, file } of dirs) {
    const path = `${dir}/${file}`
    const raw = await readFile(path, 'utf8')
    const next = stripTrackingImages(raw)
    if (next !== raw) {
      await writeFile(path, next)
      changed++
      console.log(`cleaned: ${path}`)
    }
  }
  console.log(`changed files: ${changed}`)
}

main().catch((e) => {
  console.error(e)
  process.exit(1)
})
```

- [ ] **Step 6: クリーニングを実行**

Run: `npm run cleanup:tracking`
Expected: 実数を報告（画像記法 `![]()` に計測ピクセルを含むファイルのみが対象。アフィリエイトの「リンク」記法は保持される）。

検証（画像記法のトラッキングピクセルが残っていないこと。リンク記法は残ってよい）:
```bash
grep -rhoE '!\[[^]]*\]\((https?:)?//[^)]*\)' posts/ | grep -iE 'a8\.net|valuecommerce|accesstrade' | wc -l
```
Expected: `0`（画像記法の計測ピクセルは全除去）。なお `[text](url)` のアフィリエイトリンクは保持されるため URL 文字列自体は本文に残る（これは正常）。

- [ ] **Step 7: ビルドが引き続き通ることを確認**

Run: `npm run build`
Expected: 成功。

- [ ] **Step 8: コミット**

```bash
git add scripts/lib/cleanup-tracking.ts scripts/lib/cleanup-tracking.test.ts scripts/cleanup-tracking-images.ts posts/
git commit -m "feat: 本文の不可視トラッキング画像を除去"
```

---

## Task 19: 最終検証とドキュメント更新

**Files:**
- Modify: `scripts/README.md`

- [ ] **Step 1: 全テストが通ることを確認**

Run: `npm test`
Expected: PASS（date / taxonomy / pagination / permalinks / cleanup-tracking と既存 lib テスト）。

- [ ] **Step 2: クリーンビルドと完了条件の確認**

Run: `rm -rf dist && npm run build`
Expected: 成功。

Run:
```bash
# 記事85 + 固定ページ4 + アーカイブ + トップ/RSS/sitemap/404 の主要URLを確認
echo "posts(dirs with index.html under category paths):"; \
node -e "const m=require('./data/permalinks.json');const fs=require('fs');let ok=0,ng=0;for(const id in m){const p='dist'+decodeURIComponent(m[id].path)+'index.html';fs.existsSync(p)?ok++:(ng++,console.log('MISSING',p));}console.log('ok='+ok,'missing='+ng)"
for p in contact privacy-policy profile sitemap; do test -f "dist/$p/index.html" && echo "page $p OK" || echo "page $p MISSING"; done
test -f dist/index.html && test -f dist/rss.xml && test -f dist/sitemap-index.xml && test -f dist/404.html && echo "core OK"
```
Expected: `ok=85 missing=0`、各 `page ... OK`、`core OK`。

- [ ] **Step 3: URL 保持の照合（生成 URL = permalinks.json）**

Run:
```bash
node -e "const m=require('./data/permalinks.json');const fs=require('fs');let ng=0;for(const id in m){if(!fs.existsSync('dist'+decodeURIComponent(m[id].path)+'index.html')){ng++;console.log('URL不一致:',m[id].path)}}process.exit(ng?1:0)" && echo "URL保持 OK"
```
Expected: `URL保持 OK`

- [ ] **Step 4: `scripts/README.md` に追加スクリプトを追記**

「実行順」の表/節に以下を追記（計画2分）:
```markdown
## 計画2（Astroサイト構築）で追加のスクリプト

| ファイル | 役割 |
|---|---|
| `lib/permalinks.ts` | 記事 canonical URL の link→path 変換（純粋関数） |
| `export-permalinks.ts` | 全記事 canonical URL を `data/permalinks.json` に出力（I/O・WP稼働中のみ） |
| `lib/cleanup-tracking.ts` | 本文の不可視トラッキング画像除去（純粋関数） |
| `cleanup-tracking-images.ts` | 上記を全本文へ適用（I/O・冪等） |

- `npm run export:permalinks` — URL保持の要。**ConoHa 稼働中に必ず実行・コミット**。
- `npm run cleanup:tracking` — A8/valuecommerce/accesstrade の計測画像を除去（冪等）。
```

- [ ] **Step 5: コミット**

```bash
git add scripts/README.md
git commit -m "docs: 計画2で追加した移行スクリプトをREADMEに追記"
```

- [ ] **Step 6: ローカルプレビューで目視確認（任意・推奨）**

Run: `npm run preview`
Expected: 表示される URL（例 http://localhost:4321/ ）でトップ・記事・カテゴリ・タグ・404・テーマ切替（ライト/ダーク/システム）・右 TOC 追従を目視確認。確認後 Ctrl+C。

---

## 完了条件チェックリスト（設計書 §1 と対応）

- [ ] ローカルで `npm run build` が成功する（Task 19 Step 2）
- [ ] 全 85 記事 + 4 固定ページが生成・表示される（Task 19 Step 2）
- [ ] 生成 URL が現行 WordPress と一致（`data/permalinks.json` 照合, Task 19 Step 3）
- [ ] カテゴリ/タグ アーカイブ・ページネーション・RSS・サイトマップ・404 が機能（Task 11/15/16/17）
- [ ] `npm test`（vitest）が通る（Task 19 Step 1）

## スコープ外（後続）
- 検索 Pagefind / コメント giscus / 問い合わせ Functions → 計画3
- `_redirects`(301) / デプロイ / DNS / ConoHa 解約 → 計画4
- 本文中のリンクカード（ファビコン/OGP 画像）の体裁調整は本フェーズ対象外（表示はされる）
