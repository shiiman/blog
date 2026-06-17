import { writeFile, mkdir } from 'node:fs/promises'
import { buildPermalinkMap, type RestPost } from './lib/permalinks'

const SITE = 'https://shiimanblog.com'
const DATA_DIR = 'data'

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
