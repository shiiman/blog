import { readFile, writeFile } from 'node:fs/promises'
import { glob } from 'node:fs/promises'
import matter from 'gray-matter'
import type { TermDict } from './lib/taxonomy'
import { categoryUrlPath } from '../src/lib/taxonomy'
import { selectRedirectTerms, buildRedirects } from './lib/redirects'

const CATEGORIES_JSON = 'data/categories.json'
const TAGS_JSON = 'data/tags.json'
const POSTS_DIR = 'posts'
const OUTPUT = 'public/_redirects'

async function collectUsedSlugs(): Promise<{ cats: Set<string>; tags: Set<string> }> {
  const cats = new Set<string>()
  const tags = new Set<string>()
  for await (const file of glob(`${POSTS_DIR}/*/article.md`)) {
    const raw = await readFile(file, 'utf-8')
    const { data } = matter(raw)
    if (data.draft === true) continue
    if (Array.isArray(data.categories)) {
      for (const s of data.categories) cats.add(String(s))
    }
    if (Array.isArray(data.tags)) {
      for (const s of data.tags) tags.add(String(s))
    }
  }
  return { cats, tags }
}

async function main() {
  const [catDict, tagDict] = await Promise.all([
    readFile(CATEGORIES_JSON, 'utf-8').then((s) => JSON.parse(s) as TermDict),
    readFile(TAGS_JSON, 'utf-8').then((s) => JSON.parse(s) as TermDict),
  ])

  const { cats, tags } = await collectUsedSlugs()

  const categories = selectRedirectTerms(catDict, cats).map(({ slug, enSlug }) => ({
    slug,
    enSlug: categoryUrlPath(enSlug),
  }))
  const tagTerms = selectRedirectTerms(tagDict, tags)

  const text = buildRedirects({ categories, tags: tagTerms })
  await writeFile(OUTPUT, text)

  console.log(`public/_redirects を生成しました（カテゴリ: ${categories.length}件、タグ: ${tagTerms.length}件）`)
}

main().catch((e) => {
  console.error(e)
  process.exit(1)
})
