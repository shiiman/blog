import { getPublishedPosts } from './posts'
import { resolveCategory, categoryParent } from './taxonomy'

export interface CatNode {
  name: string
  urlSlug: string
  count: number
  children: CatNode[]
}

/** カテゴリを親子階層で取得（トップレベル・子ともに記事数の降順） */
export async function getCategoryTree(): Promise<CatNode[]> {
  const posts = await getPublishedPosts()
  // enSlug（urlSlug の末尾セグメント）で記事数を集計
  const map = new Map<string, { name: string; enSlug: string; urlSlug: string; count: number }>()
  for (const post of posts) {
    for (const stored of post.entry.data.categories) {
      const { name, urlSlug } = resolveCategory(stored)
      const enSlug = urlSlug.split('/').pop()!
      const g = map.get(enSlug) ?? { name, enSlug, urlSlug, count: 0 }
      g.count++
      map.set(enSlug, g)
    }
  }
  const all = [...map.values()]
  const top = all.filter((c) => !categoryParent(c.enSlug)).sort((a, b) => b.count - a.count)
  return top.map((p) => ({
    name: p.name,
    urlSlug: p.urlSlug,
    count: p.count,
    children: all
      .filter((c) => categoryParent(c.enSlug) === p.enSlug)
      .sort((a, b) => b.count - a.count)
      .map((c) => ({ name: c.name, urlSlug: c.urlSlug, count: c.count, children: [] })),
  }))
}
