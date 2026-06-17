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
    const routeParam = decodeURIComponent(link.path.replace(/^\/+|\/+$/g, ''))
    const primarySlug = routeParam.split('/')[0]
    // 記事カテゴリのうち最も具体的なもの(子カテゴリ)を優先。無ければURL先頭の親
    const resolved = entry.data.categories.map((c) => resolveCategory(c))
    const primaryCategory = resolved.find((c) => c.urlSlug.includes('/')) ?? resolveCategory(primarySlug)
    return {
      entry,
      path: link.path,
      routeParam,
      primaryCategory,
    }
  })
  enriched.sort((a, b) => (a.entry.data.date < b.entry.data.date ? 1 : -1))
  return enriched
}
