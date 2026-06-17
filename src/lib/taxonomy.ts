import categoriesData from '../../data/categories.json'
import tagsData from '../../data/tags.json'

export interface TaxonomyTerm {
  name: string
  slug: string
  enSlug?: string
}
export type TaxonomyMap = Record<string, TaxonomyTerm>

export interface ResolvedTaxonomy {
  name: string
  urlSlug: string
}

const categories = categoriesData as TaxonomyMap
const tags = tagsData as TaxonomyMap

/** 保存slug(エンコード含む)から表示名とURL用slug(enSlug優先)を解決する */
export function resolveTaxonomy(storedSlug: string, map: TaxonomyMap): ResolvedTaxonomy {
  for (const key of Object.keys(map)) {
    const term = map[key]
    if (term.slug === storedSlug) {
      return { name: term.name, urlSlug: term.enSlug ?? term.slug }
    }
  }
  return { name: storedSlug, urlSlug: storedSlug }
}

// カテゴリの親子関係（WordPress時の階層URLを復元）。子enSlug -> 親enSlug
const CATEGORY_PARENT: Record<string, string> = {
  mac: 'engineering',
  aws: 'engineering',
  gcp: 'engineering',
  terraform: 'engineering',
  docker: 'engineering',
  'git-engineering': 'engineering',
  golang: 'engineering',
  'slack-engineering': 'engineering',
  gas: 'engineering',
  snowflake: 'engineering',
  investment: 'fire',
  savings: 'fire',
  initialization: 'wordpress',
  plugin: 'wordpress',
  cocoon: 'wordpress',
  advertisement: 'wordpress',
  analytics: 'wordpress',
}

/** カテゴリのURL用パス。子カテゴリは「親/子」の階層パスを返す */
export function categoryUrlPath(enSlug: string): string {
  const parent = CATEGORY_PARENT[enSlug]
  return parent ? `${parent}/${enSlug}` : enSlug
}

/** カテゴリの親enSlug（子の場合）。親カテゴリ/未知は undefined */
export function categoryParent(enSlug: string): string | undefined {
  return CATEGORY_PARENT[enSlug]
}

export function resolveCategory(slug: string): ResolvedTaxonomy {
  const base = resolveTaxonomy(slug, categories)
  return { name: base.name, urlSlug: categoryUrlPath(base.urlSlug) }
}
export const resolveTag = (slug: string): ResolvedTaxonomy => resolveTaxonomy(slug, tags)
