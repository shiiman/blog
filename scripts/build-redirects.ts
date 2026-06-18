import { readFile, writeFile } from 'node:fs/promises'
import { glob } from 'node:fs/promises'
import { existsSync } from 'node:fs'
import matter from 'gray-matter'
import type { TermDict } from './lib/taxonomy'
import { categoryUrlPath } from '../src/lib/taxonomy'
import { selectRedirectTerms, buildRedirects, appendManualRedirects } from './lib/redirects'

const CATEGORIES_JSON = 'data/categories.json'
const TAGS_JSON = 'data/tags.json'
const MANUAL_REDIRECTS = 'data/manual-redirects.txt'
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

  let text = buildRedirects({ categories, tags: tagTerms })

  // 自動生成できない手動リダイレクト（個別記事の slug 変更等）を末尾へ結合する。
  // 再生成しても消えないよう、定義は data/manual-redirects.txt に分離している。
  let manualCount = 0
  if (existsSync(MANUAL_REDIRECTS)) {
    const manual = await readFile(MANUAL_REDIRECTS, 'utf-8')
    manualCount = manual
      .split('\n')
      .filter((l) => l.trim() && !l.trim().startsWith('#')).length
    text = appendManualRedirects(text, manual)
  }
  await writeFile(OUTPUT, text)

  console.log(
    `public/_redirects を生成しました（カテゴリ: ${categories.length}件、タグ: ${tagTerms.length}件、手動: ${manualCount}件）`,
  )
}

main().catch((e) => {
  console.error(e)
  process.exit(1)
})
