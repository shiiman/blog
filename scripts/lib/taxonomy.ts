export interface TaxonomyTerm {
  name: string
  slug: string
}

export type TaxonomyMap = Record<string, TaxonomyTerm>

interface RestTerm {
  id: number
  name: string
  slug: string
}

/** WordPress REST のカテゴリ/タグ配列を ID 文字列キーのマップへ変換する */
export function buildTaxonomyMap(items: RestTerm[]): TaxonomyMap {
  const map: TaxonomyMap = {}
  for (const item of items) {
    map[String(item.id)] = { name: item.name, slug: item.slug }
  }
  return map
}

/** 数値ID配列を slug 配列へ変換する。未知IDはエラー */
export function mapIdsToSlugs(ids: number[], map: TaxonomyMap): string[] {
  return ids.map((id) => {
    const term = map[String(id)]
    if (!term) throw new Error(`未知のタクソノミーID: ${id}`)
    return term.slug
  })
}
