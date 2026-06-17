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

export const resolveCategory = (slug: string): ResolvedTaxonomy => resolveTaxonomy(slug, categories)
export const resolveTag = (slug: string): ResolvedTaxonomy => resolveTaxonomy(slug, tags)
