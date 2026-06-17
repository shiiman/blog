import { writeFile, readFile, mkdir } from 'node:fs/promises'
import matter from 'gray-matter'
import { buildTaxonomyMap } from './lib/taxonomy'
import { listContentDirs } from './lib/content-roots'

const SITE = 'https://shiimanblog.com'
const DATA_DIR = 'data'

type Term = { id: number; name: string; slug: string }

/** REST レスポンス項目が想定のタクソノミー形だと実行時に検証する */
function assertTerm(item: unknown, endpoint: string): asserts item is Term {
  const t = item as Record<string, unknown>
  if (typeof t?.id !== 'number' || typeof t?.name !== 'string' || typeof t?.slug !== 'string') {
    throw new Error(`想定外のタクソノミー項目 (${endpoint}): ${JSON.stringify(item)}`)
  }
}

/** ページングしながら REST エンドポイントの全件を取得する */
async function fetchAll(endpoint: string): Promise<Term[]> {
  const items: Term[] = []
  for (let page = 1; ; page++) {
    const res = await fetch(`${SITE}/wp-json/wp/v2/${endpoint}?per_page=100&page=${page}`)
    if (res.status === 400) break // ページ超過時にWPは400を返す
    if (!res.ok) throw new Error(`取得失敗 ${endpoint} p${page}: ${res.status}`)
    const batch = (await res.json()) as unknown
    if (!Array.isArray(batch)) throw new Error(`想定外レスポンス (${endpoint} p${page}): 配列ではありません`)
    if (batch.length === 0) break
    for (const item of batch) assertTerm(item, endpoint)
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
    const m = (await res.json()) as { source_url?: unknown }
    if (typeof m.source_url !== 'string') {
      console.warn(`media ${id}: source_url が不正（スキップ）`)
      continue
    }
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
