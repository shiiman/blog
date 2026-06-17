export interface PermalinkEntry {
  slug: string
  path: string
}

export type PermalinkMap = Record<string, PermalinkEntry>

interface RestPost {
  id: number
  slug: string
  link: string
}

/** 絶対URLからパス部分を取り出し、末尾スラッシュを保証する（クエリ/ハッシュ除去） */
export function linkToPath(link: string): string {
  const { pathname } = new URL(link)
  return pathname.endsWith('/') ? pathname : `${pathname}/`
}

/** REST の記事配列から id→{slug,path} マップを構築する */
export function buildPermalinkMap(items: RestPost[]): PermalinkMap {
  const map: PermalinkMap = {}
  for (const p of items) {
    map[String(p.id)] = { slug: p.slug, path: linkToPath(p.link) }
  }
  return map
}
