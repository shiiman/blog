# 計画1: データ移行レイヤー Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** WordPress 依存のコンテンツ（数値IDのタクソノミー・リモート画像）を、Astro が読める形（ローカル画像＋slug ベースのフロントマター）へ完全変換し、`posts/` と `pages/` を「サイト構築可能な状態」にする。

**Architecture:** WordPress REST API（公開エンドポイント・認証不要）から Node スクリプトでタクソノミー／アイキャッチ URL をエクスポートし、`data/*.json` に保存。続いて純粋関数（テスト可能）で本文の画像 URL ローカル化とフロントマター変換を行い、I/O を担うオーケストレーションスクリプトが全記事・全ページへ適用する。`mahjong`/`poker` は事前に `future-sites/` へ退避してスコープ外にする。

**Tech Stack:** Node.js (ESM) + TypeScript / tsx（実行）/ vitest（テスト）/ gray-matter（フロントマター）/ Git LFS（画像）/ WordPress REST API。

> 関連仕様書: `docs/superpowers/specs/2026-06-17-blog-migration-astro-cloudflare-design.md`
> この計画は全4計画の1つ目。後続: 計画2(Astro構築)・計画3(動的機能)・計画4(デプロイ/切替)。

---

## File Structure

| ファイル | 役割 |
|---|---|
| `package.json` | Node プロジェクト定義・スクリプト・依存 |
| `tsconfig.json` | TypeScript 設定（ESM） |
| `vitest.config.ts` | テスト設定 |
| `scripts/lib/taxonomy.ts` | REST 配列→ID→{name,slug} マップ／ID配列→slug配列（純粋関数） |
| `scripts/lib/taxonomy.test.ts` | 上記のテスト |
| `scripts/lib/frontmatter.ts` | WP フロントマター→Astro フロントマター変換（純粋関数） |
| `scripts/lib/frontmatter.test.ts` | 上記のテスト |
| `scripts/lib/images.ts` | 本文の画像 URL 抽出・正規化・ローカルパス置換（純粋関数） |
| `scripts/lib/images.test.ts` | 上記のテスト |
| `scripts/lib/content-roots.ts` | 処理対象ルート（posts/article.md・pages/page.md）の共有定義 |
| `scripts/export-wp-data.ts` | REST から `data/*.json`・`data/old-sitemap.xml` を生成（I/O） |
| `scripts/localize-images.ts` | 画像DL＋本文書き換え＋アイキャッチDL（I/O） |
| `scripts/migrate-frontmatter.ts` | `data/*.json` を用いて全フロントマターを変換（I/O） |
| `data/categories.json` 他 | エクスポート成果物（コミット対象） |
| `.gitattributes` | `pages/**/assets/*` を LFS 追跡対象に追加 |
| `.gitignore` | `node_modules/` を追加 |

---

## Task 1: Node プロジェクト初期化

**Files:**
- Create: `package.json`
- Create: `tsconfig.json`
- Create: `vitest.config.ts`
- Modify: `.gitignore`
- Modify: `.gitattributes`

- [ ] **Step 1: `package.json` を作成**

```json
{
  "name": "shiimanblog",
  "private": true,
  "type": "module",
  "scripts": {
    "test": "vitest run",
    "export:wp": "tsx scripts/export-wp-data.ts",
    "localize:images": "tsx scripts/localize-images.ts",
    "migrate:frontmatter": "tsx scripts/migrate-frontmatter.ts"
  },
  "devDependencies": {
    "gray-matter": "^4.0.3",
    "tsx": "^4.19.0",
    "typescript": "^5.6.3",
    "vitest": "^2.1.8"
  }
}
```

- [ ] **Step 2: `tsconfig.json` を作成**

```json
{
  "compilerOptions": {
    "target": "ES2022",
    "module": "ESNext",
    "moduleResolution": "Bundler",
    "strict": true,
    "esModuleInterop": true,
    "skipLibCheck": true,
    "types": ["node"],
    "noEmit": true
  },
  "include": ["scripts"]
}
```

- [ ] **Step 3: `vitest.config.ts` を作成**

```ts
import { defineConfig } from 'vitest/config'

export default defineConfig({
  test: {
    include: ['scripts/**/*.test.ts'],
  },
})
```

- [ ] **Step 4: `.gitignore` に `node_modules/` を追記**

`.gitignore` の末尾に以下を追加する:

```
node_modules/
```

- [ ] **Step 5: `.gitattributes` に pages 配下の画像 LFS 追跡を追加**

`posts/**/assets/*` の各行の下に、対応する `pages/**/assets/*` 行を追加する（CLAUDE.md のルール「画像追加前に .gitattributes を先に更新」に従う）:

```
pages/**/assets/*.png filter=lfs diff=lfs merge=lfs -text
pages/**/assets/*.jpg filter=lfs diff=lfs merge=lfs -text
pages/**/assets/*.jpeg filter=lfs diff=lfs merge=lfs -text
pages/**/assets/*.gif filter=lfs diff=lfs merge=lfs -text
pages/**/assets/*.webp filter=lfs diff=lfs merge=lfs -text
```

- [ ] **Step 6: 依存をインストール**

Run: `npm install`
Expected: `node_modules/` が生成され、エラーなく完了する。

- [ ] **Step 7: テストランナーが利用可能なことを確認**

Run: `npx vitest --version`
Expected: vitest のバージョン番号が表示される（この時点ではテストファイルが無いため `npm test` は実行しない。最初のテストは Task 3 で追加する）。

- [ ] **Step 8: コミット**

```bash
git add package.json package-lock.json tsconfig.json vitest.config.ts .gitignore .gitattributes
git commit -m "build: 移行スクリプト用のNode/TypeScript/vitest環境を追加"
```

---

## Task 2: mahjong / poker を future-sites/ へ退避

**Files:**
- Move: `pages/mahjong` → `future-sites/mahjong`
- Move: `pages/poker` → `future-sites/poker`

- [ ] **Step 1: 退避先を作成して `git mv` で移動**

Run:
```bash
mkdir -p future-sites
git mv pages/mahjong future-sites/mahjong
git mv pages/poker future-sites/poker
```

- [ ] **Step 2: 結果を確認**

Run: `ls pages && echo '---' && ls future-sites`
Expected: `pages` に `contact / privacy-policy / profile / sitemap` の4件のみ。`future-sites` に `mahjong / poker`。

- [ ] **Step 3: コミット**

```bash
git add -A
git commit -m "chore: 今回の移行対象外のmahjong/pokerをfuture-sitesへ退避"
```

---

## Task 3: タクソノミー判定ロジック（純粋関数）

**Files:**
- Create: `scripts/lib/taxonomy.ts`
- Test: `scripts/lib/taxonomy.test.ts`

- [ ] **Step 1: 失敗するテストを書く**

`scripts/lib/taxonomy.test.ts`:

```ts
import { describe, it, expect } from 'vitest'
import { buildTaxonomyMap, mapIdsToSlugs } from './taxonomy'

describe('buildTaxonomyMap', () => {
  it('REST配列をID文字列キーの{name,slug}マップへ変換する', () => {
    const map = buildTaxonomyMap([
      { id: 2, name: '技術', slug: 'tech' },
      { id: 10, name: 'ブログ', slug: 'blog' },
    ])
    expect(map).toEqual({
      '2': { name: '技術', slug: 'tech' },
      '10': { name: 'ブログ', slug: 'blog' },
    })
  })
})

describe('mapIdsToSlugs', () => {
  const map = { '2': { name: '技術', slug: 'tech' }, '10': { name: 'ブログ', slug: 'blog' } }

  it('ID配列をslug配列へ変換する', () => {
    expect(mapIdsToSlugs([2, 10], map)).toEqual(['tech', 'blog'])
  })

  it('空配列は空配列を返す', () => {
    expect(mapIdsToSlugs([], map)).toEqual([])
  })

  it('未知のIDはエラーを投げる', () => {
    expect(() => mapIdsToSlugs([999], map)).toThrow('未知のタクソノミーID: 999')
  })
})
```

- [ ] **Step 2: テストが失敗することを確認**

Run: `npx vitest run scripts/lib/taxonomy.test.ts`
Expected: FAIL（`taxonomy` モジュールが存在しない）

- [ ] **Step 3: 実装する**

`scripts/lib/taxonomy.ts`:

```ts
export interface TaxonomyTerm {
  name: string
  slug: string
}

export type TaxonomyMap = Record<string, TaxonomyTerm>

interface RestTerm {
  id: number
  name: string
  slug: string
}

/** WordPress REST のカテゴリ/タグ配列を ID 文字列キーのマップへ変換する */
export function buildTaxonomyMap(items: RestTerm[]): TaxonomyMap {
  const map: TaxonomyMap = {}
  for (const item of items) {
    map[String(item.id)] = { name: item.name, slug: item.slug }
  }
  return map
}

/** 数値ID配列を slug 配列へ変換する。未知IDはエラー */
export function mapIdsToSlugs(ids: number[], map: TaxonomyMap): string[] {
  return ids.map((id) => {
    const term = map[String(id)]
    if (!term) throw new Error(`未知のタクソノミーID: ${id}`)
    return term.slug
  })
}
```

- [ ] **Step 4: テストが通ることを確認**

Run: `npx vitest run scripts/lib/taxonomy.test.ts`
Expected: PASS（4 tests）

- [ ] **Step 5: コミット**

```bash
git add scripts/lib/taxonomy.ts scripts/lib/taxonomy.test.ts
git commit -m "feat: タクソノミーID→slug変換の純粋関数を追加"
```

---

## Task 4: フロントマター変換ロジック（純粋関数）

**Files:**
- Create: `scripts/lib/frontmatter.ts`
- Test: `scripts/lib/frontmatter.test.ts`

- [ ] **Step 1: 失敗するテストを書く**

`scripts/lib/frontmatter.test.ts`:

```ts
import { describe, it, expect } from 'vitest'
import { transformFrontmatter } from './frontmatter'

const catMap = { '2': { name: '技術', slug: 'tech' } }
const tagMap = { '6': { name: 'ブログ', slug: 'blog' }, '7': { name: 'SEO', slug: 'seo' } }

describe('transformFrontmatter', () => {
  const base = {
    id: 6,
    title: 'タイトル',
    slug: 'start-blog',
    status: 'publish',
    date: '2021-08-30T19:30:00',
    modified: '2021-08-30T19:22:28',
    excerpt: '概要',
    categories: [2],
    tags: [6, 7],
    featured_media: 13,
  }

  it('カテゴリ/タグをslugへ変換しdraftを設定する', () => {
    const out = transformFrontmatter(base, catMap, tagMap)
    expect(out.categories).toEqual(['tech'])
    expect(out.tags).toEqual(['blog', 'seo'])
    expect(out.draft).toBe(false)
    expect(out.id).toBe(6)
    expect(out.title).toBe('タイトル')
  })

  it('status=draft は draft:true になる', () => {
    expect(transformFrontmatter({ ...base, status: 'draft' }, catMap, tagMap).draft).toBe(true)
  })

  it('status と featured_media は出力に残さない', () => {
    const out = transformFrontmatter(base, catMap, tagMap) as Record<string, unknown>
    expect(out.status).toBeUndefined()
    expect(out.featured_media).toBeUndefined()
  })

  it('eyecatchFile が渡されると相対パスを設定する', () => {
    const out = transformFrontmatter(base, catMap, tagMap, 'eyecatch.png')
    expect(out.eyecatch).toBe('./assets/eyecatch.png')
  })

  it('eyecatchFile 無しなら eyecatch は未設定', () => {
    expect(transformFrontmatter(base, catMap, tagMap).eyecatch).toBeUndefined()
  })

  it('categories/tags が無い（固定ページ）場合は空配列', () => {
    const page = { title: 'プロフィール', slug: 'profile', status: 'publish', date: '2021-08-31T02:29:25' }
    const out = transformFrontmatter(page, catMap, tagMap)
    expect(out.categories).toEqual([])
    expect(out.tags).toEqual([])
  })
})
```

- [ ] **Step 2: テストが失敗することを確認**

Run: `npx vitest run scripts/lib/frontmatter.test.ts`
Expected: FAIL（`frontmatter` モジュールが存在しない）

- [ ] **Step 3: 実装する**

`scripts/lib/frontmatter.ts`:

```ts
import type { TaxonomyMap } from './taxonomy'
import { mapIdsToSlugs } from './taxonomy'

export interface WpFrontmatter {
  id?: number
  title: string
  slug: string
  status: string
  date: string
  modified?: string
  excerpt?: string
  categories?: number[]
  tags?: number[]
  featured_media?: number
}

export interface AstroFrontmatter {
  id?: number
  title: string
  slug: string
  date: string
  modified?: string
  excerpt?: string
  categories: string[]
  tags: string[]
  eyecatch?: string
  draft: boolean
}

/** WordPress 前提のフロントマターを Astro 用スキーマへ変換する（純粋関数） */
export function transformFrontmatter(
  fm: WpFrontmatter,
  catMap: TaxonomyMap,
  tagMap: TaxonomyMap,
  eyecatchFile?: string,
): AstroFrontmatter {
  const out: AstroFrontmatter = {
    title: fm.title,
    slug: fm.slug,
    date: fm.date,
    categories: mapIdsToSlugs(fm.categories ?? [], catMap),
    tags: mapIdsToSlugs(fm.tags ?? [], tagMap),
    draft: fm.status !== 'publish',
  }
  if (fm.id !== undefined) out.id = fm.id
  if (fm.modified !== undefined) out.modified = fm.modified
  if (fm.excerpt !== undefined) out.excerpt = fm.excerpt
  if (eyecatchFile) out.eyecatch = `./assets/${eyecatchFile}`
  return out
}
```

- [ ] **Step 4: テストが通ることを確認**

Run: `npx vitest run scripts/lib/frontmatter.test.ts`
Expected: PASS（6 tests）

- [ ] **Step 5: コミット**

```bash
git add scripts/lib/frontmatter.ts scripts/lib/frontmatter.test.ts
git commit -m "feat: WPフロントマター→Astroスキーマ変換の純粋関数を追加"
```

---

## Task 5: 画像URL処理ロジック（純粋関数）

**Files:**
- Create: `scripts/lib/images.ts`
- Test: `scripts/lib/images.test.ts`

- [ ] **Step 1: 失敗するテストを書く**

`scripts/lib/images.test.ts`:

```ts
import { describe, it, expect } from 'vitest'
import {
  toOriginalUrl, extractImageUrls, normalizeLightbox, stripTrackingPixels, rewriteImageUrls,
} from './images'

describe('toOriginalUrl', () => {
  it('-WxH サイズ接尾辞を除去する', () => {
    expect(toOriginalUrl('https://shiimanblog.com/wp-content/uploads/2021/08/4-1-1024x545.jpg'))
      .toBe('https://shiimanblog.com/wp-content/uploads/2021/08/4-1.jpg')
  })
  it('接尾辞が無ければそのまま', () => {
    expect(toOriginalUrl('https://shiimanblog.com/wp-content/uploads/2021/09/3.png'))
      .toBe('https://shiimanblog.com/wp-content/uploads/2021/09/3.png')
  })
})

describe('extractImageUrls', () => {
  it('wp-content の画像URLを重複なく抽出する', () => {
    const md = '![](https://shiimanblog.com/wp-content/uploads/2021/09/3.png) と ![](https://shiimanblog.com/wp-content/uploads/2021/09/3.png)'
    expect(extractImageUrls(md)).toEqual(['https://shiimanblog.com/wp-content/uploads/2021/09/3.png'])
  })
  it('外部ドメインの画像は無視する', () => {
    expect(extractImageUrls('![](https://www19.a8.net/0.gif)')).toEqual([])
  })
})

describe('normalizeLightbox', () => {
  it('[![alt](thumb)](full) を ![alt](full) に正規化する', () => {
    const md = '[![](https://shiimanblog.com/wp-content/uploads/2021/08/4-1-1024x545.jpg)](https://shiimanblog.com/wp-content/uploads/2021/08/4-1.jpg)'
    expect(normalizeLightbox(md)).toBe('![](https://shiimanblog.com/wp-content/uploads/2021/08/4-1.jpg)')
  })
})

describe('stripTrackingPixels', () => {
  it('a8.net の計測gif画像を除去する', () => {
    expect(stripTrackingPixels('text ![](https://www19.a8.net/0.gif?a8mat=XXX) end')).toBe('text  end')
  })
})

describe('rewriteImageUrls', () => {
  it('原本URLへのマップでローカルパスへ置換する（サイズ違いも原本扱い）', () => {
    const map = { 'https://shiimanblog.com/wp-content/uploads/2021/08/4-1.jpg': './assets/4-1.jpg' }
    const md = '![](https://shiimanblog.com/wp-content/uploads/2021/08/4-1-1024x545.jpg)'
    expect(rewriteImageUrls(md, map)).toBe('![](./assets/4-1.jpg)')
  })
})
```

- [ ] **Step 2: テストが失敗することを確認**

Run: `npx vitest run scripts/lib/images.test.ts`
Expected: FAIL（`images` モジュールが存在しない）

- [ ] **Step 3: 実装する**

`scripts/lib/images.ts`:

```ts
const WP_UPLOAD_RE =
  /https?:\/\/shiimanblog\.com\/wp-content\/uploads\/[^\s)"'\\]+\.(?:png|jpe?g|gif|webp)/gi

/** -1024x545 のようなサイズ接尾辞を除去し原本URLにする */
export function toOriginalUrl(url: string): string {
  return url.replace(/-\d+x\d+(\.(?:png|jpe?g|gif|webp))$/i, '$1')
}

/** 本文から wp-content 画像URLを重複なく抽出する */
export function extractImageUrls(markdown: string): string[] {
  const found = markdown.match(WP_UPLOAD_RE) ?? []
  return [...new Set(found)]
}

/** WordPress ライトボックス記法 [![alt](thumb)](full) を ![alt](full) に正規化する */
export function normalizeLightbox(markdown: string): string {
  return markdown.replace(
    /\[!\[([^\]]*)\]\(([^)]+)\)\]\(([^)]+)\)/g,
    (_m, alt: string, _thumb: string, full: string) => `![${alt}](${full})`,
  )
}

/** a8.net 等の計測用gif画像記法を除去する */
export function stripTrackingPixels(markdown: string): string {
  return markdown.replace(
    /!\[[^\]]*\]\(https?:\/\/[^\s)]*a8\.net\/[^\s)]*\.gif[^)]*\)/gi,
    '',
  )
}

/** 本文中の wp-content URL（サイズ違い含む）をローカル相対パスへ置換する */
export function rewriteImageUrls(markdown: string, urlToLocal: Record<string, string>): string {
  return markdown.replace(WP_UPLOAD_RE, (url) => urlToLocal[toOriginalUrl(url)] ?? url)
}
```

- [ ] **Step 4: テストが通ることを確認**

Run: `npx vitest run scripts/lib/images.test.ts`
Expected: PASS（全テスト）

- [ ] **Step 5: コミット**

```bash
git add scripts/lib/images.ts scripts/lib/images.test.ts
git commit -m "feat: 本文画像URLの抽出・正規化・ローカル置換の純粋関数を追加"
```

---

## Task 6: 処理対象ルートの共有定義

**Files:**
- Create: `scripts/lib/content-roots.ts`

- [ ] **Step 1: 実装する**（純粋な定数のみ。テスト不要）

`scripts/lib/content-roots.ts`:

```ts
import { readdir } from 'node:fs/promises'

/** 変換対象のコンテンツルート。future-sites/ は含めない（スコープ外） */
export const CONTENT_ROOTS = [
  { base: 'posts', file: 'article.md' },
  { base: 'pages', file: 'page.md' },
] as const

/** 各ルート配下のディレクトリ名を列挙して {dir,file} を返す */
export async function listContentDirs(): Promise<{ dir: string; file: string }[]> {
  const result: { dir: string; file: string }[] = []
  for (const root of CONTENT_ROOTS) {
    const entries = await readdir(root.base, { withFileTypes: true })
    for (const e of entries) {
      if (e.isDirectory()) result.push({ dir: `${root.base}/${e.name}`, file: root.file })
    }
  }
  return result
}
```

- [ ] **Step 2: コミット**

```bash
git add scripts/lib/content-roots.ts
git commit -m "feat: 変換対象コンテンツルートの共有定義を追加"
```

---

## Task 7: WP データエクスポートスクリプト（REST）

**Files:**
- Create: `scripts/export-wp-data.ts`
- Output: `data/categories.json` / `data/tags.json` / `data/featured-media.json` / `data/old-sitemap.xml`

- [ ] **Step 1: スクリプトを実装する**

`scripts/export-wp-data.ts`:

```ts
import { writeFile, readFile, mkdir } from 'node:fs/promises'
import matter from 'gray-matter'
import { buildTaxonomyMap } from './lib/taxonomy'
import { listContentDirs } from './lib/content-roots'

const SITE = 'https://shiimanblog.com'
const DATA_DIR = 'data'

/** ページングしながら REST エンドポイントの全件を取得する */
async function fetchAll(endpoint: string): Promise<{ id: number; name: string; slug: string }[]> {
  const items: { id: number; name: string; slug: string }[] = []
  for (let page = 1; ; page++) {
    const res = await fetch(`${SITE}/wp-json/wp/v2/${endpoint}?per_page=100&page=${page}`)
    if (res.status === 400) break // ページ超過時にWPは400を返す
    if (!res.ok) throw new Error(`取得失敗 ${endpoint} p${page}: ${res.status}`)
    const batch = (await res.json()) as { id: number; name: string; slug: string }[]
    if (batch.length === 0) break
    items.push(...batch)
    if (batch.length < 100) break
  }
  return items
}

/** 全記事・全ページのフロントマターから featured_media ID を集める */
async function collectFeaturedMediaIds(): Promise<number[]> {
  const dirs = await listContentDirs()
  const ids = new Set<number>()
  for (const { dir, file } of dirs) {
    const raw = await readFile(`${dir}/${file}`, 'utf8')
    const { data } = matter(raw)
    if (typeof data.featured_media === 'number' && data.featured_media > 0) ids.add(data.featured_media)
  }
  return [...ids]
}

/** 候補パスを順に試して最初に成功した旧サイトマップを保存する */
async function saveOldSitemap(): Promise<void> {
  for (const path of ['/wp-sitemap.xml', '/sitemap_index.xml', '/sitemap.xml']) {
    const res = await fetch(`${SITE}${path}`)
    if (res.ok) {
      await writeFile(`${DATA_DIR}/old-sitemap.xml`, await res.text())
      console.log(`old-sitemap: ${path}`)
      return
    }
  }
  console.warn('旧サイトマップが取得できませんでした。手動で確認してください。')
}

async function main() {
  await mkdir(DATA_DIR, { recursive: true })

  const categories = await fetchAll('categories')
  await writeFile(`${DATA_DIR}/categories.json`, JSON.stringify(buildTaxonomyMap(categories), null, 2) + '\n')
  console.log(`categories: ${categories.length}`)

  const tags = await fetchAll('tags')
  await writeFile(`${DATA_DIR}/tags.json`, JSON.stringify(buildTaxonomyMap(tags), null, 2) + '\n')
  console.log(`tags: ${tags.length}`)

  const ids = await collectFeaturedMediaIds()
  const media: Record<string, string> = {}
  for (const id of ids) {
    const res = await fetch(`${SITE}/wp-json/wp/v2/media/${id}`)
    if (!res.ok) {
      console.warn(`media ${id}: ${res.status}（スキップ）`)
      continue
    }
    const m = (await res.json()) as { source_url: string }
    media[String(id)] = m.source_url
  }
  await writeFile(`${DATA_DIR}/featured-media.json`, JSON.stringify(media, null, 2) + '\n')
  console.log(`featured-media: ${Object.keys(media).length}/${ids.length}`)

  await saveOldSitemap()
}

main().catch((e) => {
  console.error(e)
  process.exit(1)
})
```

- [ ] **Step 2: エクスポートを実行する**

Run: `npm run export:wp`
Expected: ログに `categories: 24`（前後）・`tags: ~100`・`featured-media: 80/80`（前後）が表示され、`data/` に4ファイルが生成される。

- [ ] **Step 3: 出力の妥当性を確認する**

Run:
```bash
node -e "const c=require('./data/categories.json'),t=require('./data/tags.json'),m=require('./data/featured-media.json'); console.log('cat',Object.keys(c).length,'tag',Object.keys(t).length,'media',Object.keys(m).length); console.log('sample cat', c['2']); console.log('sample media', Object.values(m)[0])"
```
Expected: カテゴリ・タグの件数が前述の規模と整合し、`sample cat` に `{name, slug}`、`sample media` に `https://shiimanblog.com/wp-content/...` の URL が表示される。

> Note: 全フロントマターのカテゴリ/タグIDがマップに揃っているかの厳密な突き合わせは、Task 9 で `mapIdsToSlugs` が未知IDを検出して停止すること、および Task 10 の検証スクリプトで保証する。ここでは件数確認まででよい。

- [ ] **Step 4: コミット**

```bash
git add data/categories.json data/tags.json data/featured-media.json data/old-sitemap.xml scripts/export-wp-data.ts
git commit -m "feat: WP REST APIからタクソノミー・アイキャッチ・旧サイトマップをエクスポート"
```

---

## Task 8: 画像ローカル化スクリプト

**Files:**
- Create: `scripts/localize-images.ts`
- Modify: `posts/*/article.md`・`pages/*/page.md`（本文の画像URL）
- Create: `posts/*/assets/*`・`pages/*/assets/*`（DLした画像・LFS）

- [ ] **Step 1: スクリプトを実装する**

`scripts/localize-images.ts`:

```ts
import { readFile, writeFile, mkdir } from 'node:fs/promises'
import { basename, extname } from 'node:path'
import matter from 'gray-matter'
import {
  normalizeLightbox, stripTrackingPixels, extractImageUrls, toOriginalUrl, rewriteImageUrls,
} from './lib/images'
import { listContentDirs } from './lib/content-roots'

/** URL からファイルをダウンロードして dest へ保存する */
async function download(url: string, dest: string): Promise<void> {
  const res = await fetch(url)
  if (!res.ok) throw new Error(`DL失敗 ${url}: ${res.status}`)
  await writeFile(dest, Buffer.from(await res.arrayBuffer()))
}

async function processOne(dir: string, file: string, featured: Record<string, string>): Promise<void> {
  const raw = await readFile(`${dir}/${file}`, 'utf8')
  const { data, content } = matter(raw)
  let body = stripTrackingPixels(normalizeLightbox(content))

  const urls = [...new Set(extractImageUrls(body).map(toOriginalUrl))]
  const map: Record<string, string> = {}
  if (urls.length > 0) await mkdir(`${dir}/assets`, { recursive: true })
  for (const url of urls) {
    const name = basename(new URL(url).pathname)
    await download(url, `${dir}/assets/${name}`)
    map[url] = `./assets/${name}`
  }
  body = rewriteImageUrls(body, map)

  const fmId = typeof data.featured_media === 'number' ? data.featured_media : 0
  const eyeUrl = featured[String(fmId)]
  if (eyeUrl) {
    const ext = extname(new URL(eyeUrl).pathname) || '.png'
    await mkdir(`${dir}/assets`, { recursive: true })
    await download(toOriginalUrl(eyeUrl), `${dir}/assets/eyecatch${ext}`)
  }

  await writeFile(`${dir}/${file}`, matter.stringify(body, data))
}

async function main() {
  const featured = JSON.parse(await readFile('data/featured-media.json', 'utf8')) as Record<string, string>
  for (const { dir, file } of await listContentDirs()) {
    console.log(`localize: ${dir}`)
    await processOne(dir, file, featured)
  }
}

main().catch((e) => {
  console.error(e)
  process.exit(1)
})
```

- [ ] **Step 2: 作業ツリーがクリーンなことを確認する（一括書き換え前の前提）**

Run: `git status --short`
Expected: 出力が空（直前までのコミット済み状態）。本スクリプトは本文を一括書き換えするため、差分をレビュー可能にする目的で適用前はクリーンにしておく。スクリプトは冪等（ローカル化済み本文は wp-content URL を含まないため再DLされない）。

- [ ] **Step 3: 全件にローカル化を適用する**

Run: `npm run localize:images`
Expected: 各ディレクトリについて `localize: posts/...` が出力され、エラーなく完了する。DL 失敗があれば該当 URL が例外表示される（その場合は Step 6 の対処へ）。

- [ ] **Step 4: 本文に wp-content URL が残っていないことを確認する**

Run: `grep -rIl 'shiimanblog.com/wp-content' posts pages || echo 'OK: 残存なし'`
Expected: `OK: 残存なし`（何も該当しない）。残った場合は対象ファイルを確認し、正規表現の想定外パターン（クエリ付き等）を Step 6 で対処。

- [ ] **Step 5: DL 画像が LFS ポインタとしてステージされることを確認する**

Run:
```bash
git add -A
git status --short | grep assets | head
git diff --cached -- 'posts/**/assets/*' | head -5
```
Expected: `git diff --cached` の中身が `version https://git-lfs.github.com/spec/v1` で始まる LFS ポインタであること（実体バイナリでないこと）。実体が混入していたら `git add --renormalize posts pages` を実行して再ステージ（CLAUDE.md のガード手順）。

- [ ] **Step 6: （DL失敗時のみ）対処してから再実行**

Note: 404 等で取得できない画像があった場合は、URL を `data/missing-images.txt` に控え、(a) 別サイズが取得可能ならそれを、(b) 取得不能なら本文の該当 `![](...)` を手動で除去、のいずれかで対処してから `npm run localize:images` を再実行する（冪等：既にローカル化済みの本文は wp-content URL を含まないため再DLされない）。

- [ ] **Step 7: コミット**

```bash
git add -A
git commit -m "feat: 本文画像とアイキャッチをローカル化しWP絶対URLを撤去"
```

---

## Task 9: フロントマター変換スクリプトの適用

**Files:**
- Create: `scripts/migrate-frontmatter.ts`
- Modify: `posts/*/article.md`・`pages/*/page.md`（フロントマター）

- [ ] **Step 1: スクリプトを実装する**

`scripts/migrate-frontmatter.ts`:

```ts
import { readFile, writeFile, access } from 'node:fs/promises'
import matter from 'gray-matter'
import { transformFrontmatter, type WpFrontmatter } from './lib/frontmatter'
import { listContentDirs } from './lib/content-roots'
import type { TaxonomyMap } from './lib/taxonomy'

/** assets 配下の eyecatch.<ext> を探してファイル名を返す */
async function findEyecatch(dir: string): Promise<string | undefined> {
  for (const ext of ['png', 'jpg', 'jpeg', 'gif', 'webp']) {
    try {
      await access(`${dir}/assets/eyecatch.${ext}`)
      return `eyecatch.${ext}`
    } catch {
      // 次の拡張子へ
    }
  }
  return undefined
}

async function main() {
  const catMap = JSON.parse(await readFile('data/categories.json', 'utf8')) as TaxonomyMap
  const tagMap = JSON.parse(await readFile('data/tags.json', 'utf8')) as TaxonomyMap

  for (const { dir, file } of await listContentDirs()) {
    const raw = await readFile(`${dir}/${file}`, 'utf8')
    const { data, content } = matter(raw)
    const eyecatch = await findEyecatch(dir)
    const next = transformFrontmatter(data as WpFrontmatter, catMap, tagMap, eyecatch)
    await writeFile(`${dir}/${file}`, matter.stringify(content, next as unknown as Record<string, unknown>))
    console.log(`frontmatter: ${dir}`)
  }
}

main().catch((e) => {
  console.error(e)
  process.exit(1)
})
```

- [ ] **Step 2: 全件にフロントマター変換を適用する**

Run: `npm run migrate:frontmatter`
Expected: 各ディレクトリで `frontmatter: ...` が出力される。未知のタクソノミーIDがあれば `未知のタクソノミーID: N` で停止する（その場合は `data/*.json` の不足を調査）。

- [ ] **Step 3: 変換結果を目視確認する**

Run:
```bash
head -15 posts/2021-08-30_start-blog/article.md
echo '--- page ---'
head -12 pages/profile/page.md
```
Expected: 記事側 `categories`/`tags` が slug 文字列配列、`draft: false`、`eyecatch: ./assets/eyecatch.*`（該当時）になり、`status`/`featured_media` が消えている。ページ側は `categories: []`/`tags: []`/`draft: false`。

- [ ] **Step 4: 残存する WP 固有フィールドが無いことを確認する**

Run: `grep -rIl -E '^(status|featured_media):' posts pages || echo 'OK: WP固有フィールドなし'`
Expected: `OK: WP固有フィールドなし`

- [ ] **Step 5: コミット**

```bash
git add -A
git commit -m "feat: 全記事・固定ページのフロントマターをAstroスキーマへ変換"
```

---

## Task 10: 受け入れ検証（データ移行レイヤー完了ゲート）

**Files:**
- Create: `scripts/verify-migration.ts`

- [ ] **Step 1: 検証スクリプトを実装する**

`scripts/verify-migration.ts`:

```ts
import { readFile } from 'node:fs/promises'
import matter from 'gray-matter'
import { listContentDirs } from './lib/content-roots'
import type { TaxonomyMap } from './lib/taxonomy'

async function main() {
  const catMap = JSON.parse(await readFile('data/categories.json', 'utf8')) as TaxonomyMap
  const tagMap = JSON.parse(await readFile('data/tags.json', 'utf8')) as TaxonomyMap
  const catSlugs = new Set(Object.values(catMap).map((t) => t.slug))
  const tagSlugs = new Set(Object.values(tagMap).map((t) => t.slug))

  const errors: string[] = []
  for (const { dir, file } of await listContentDirs()) {
    const raw = await readFile(`${dir}/${file}`, 'utf8')
    const { data, content } = matter(raw)

    if (/shiimanblog\.com\/wp-content/.test(content)) errors.push(`${dir}: 本文にwp-content URLが残存`)
    if ('status' in data) errors.push(`${dir}: status フィールドが残存`)
    if ('featured_media' in data) errors.push(`${dir}: featured_media フィールドが残存`)
    if (typeof data.draft !== 'boolean') errors.push(`${dir}: draft が boolean でない`)
    for (const c of (data.categories ?? []) as string[]) {
      if (!catSlugs.has(c)) errors.push(`${dir}: 未知カテゴリslug ${c}`)
    }
    for (const t of (data.tags ?? []) as string[]) {
      if (!tagSlugs.has(t)) errors.push(`${dir}: 未知タグslug ${t}`)
    }
  }

  if (errors.length > 0) {
    console.error(`検証NG: ${errors.length}件`)
    for (const e of errors) console.error(' - ' + e)
    process.exit(1)
  }
  console.log('検証OK: データ移行レイヤーは健全です')
}

main().catch((e) => {
  console.error(e)
  process.exit(1)
})
```

- [ ] **Step 2: 検証を実行する**

Run: `npx tsx scripts/verify-migration.ts`
Expected: `検証OK: データ移行レイヤーは健全です`。NG が出たら該当ディレクトリを修正し、再実行する。

- [ ] **Step 3: 全ユニットテストが通ることを確認する**

Run: `npm test`
Expected: taxonomy / frontmatter / images の全テストが PASS。

- [ ] **Step 4: LFS の最終確認**

Run: `git lfs ls-files | wc -l && echo '--- 実体混入チェック ---' && git lfs status`
Expected: ローカル化した画像が LFS 管理下にあり、ポインタとしてコミットされている。

- [ ] **Step 5: コミット**

```bash
git add scripts/verify-migration.ts
git commit -m "test: データ移行レイヤーの受け入れ検証スクリプトを追加"
```

---

## 完了の定義（この計画のゴール）

- `mahjong`/`poker` が `future-sites/` へ退避され、`pages/` は移行対象4件のみ
- `data/categories.json`・`data/tags.json`・`data/featured-media.json`・`data/old-sitemap.xml` が揃っている
- 全記事・全ページの本文に `shiimanblog.com/wp-content` 参照が無い（画像はローカル＋LFS）
- 全フロントマターが Astro スキーマ（slug 配列・`draft`・`eyecatch`）に変換済み、`status`/`featured_media` が消えている
- `npx tsx scripts/verify-migration.ts` と `npm test` がともに成功する

→ この状態で **計画2（Astro サイト構築）** に進める。
