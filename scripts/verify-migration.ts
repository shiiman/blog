import { readFile, access } from 'node:fs/promises'
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
    if (typeof data.eyecatch === 'string') {
      const eyecatchPath = `${dir}/${data.eyecatch.replace(/^\.\//, '')}`
      try {
        await access(eyecatchPath)
      } catch {
        errors.push(`${dir}: eyecatch ファイルが存在しない (${data.eyecatch})`)
      }
    }
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
