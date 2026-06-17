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
