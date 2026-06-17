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
