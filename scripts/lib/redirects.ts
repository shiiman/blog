import type { TermDict } from './taxonomy'

export interface RedirectTerm {
  slug: string
  enSlug: string
}

export interface RedirectsInput {
  categories: RedirectTerm[]
  tags: RedirectTerm[]
}

const STATIC_RULES = [
  '/feed/             /rss.xml   301',
  '/comments/feed/   /rss.xml   301',
  '/*/feed/           /rss.xml   301',
]

/**
 * term辞書と記事で実際に使用中のslug集合から、
 * リダイレクト対象（enSlug保持 かつ 使用中）のtermを抽出する
 */
export function selectRedirectTerms(dict: TermDict, usedSlugs: Set<string>): RedirectTerm[] {
  const result: RedirectTerm[] = []
  for (const entry of Object.values(dict)) {
    if (entry.enSlug && usedSlugs.has(entry.slug)) {
      result.push({ slug: entry.slug, enSlug: entry.enSlug })
    }
  }
  return result
}

/** カテゴリ・タグの enSlug リダイレクトと静的ルールから _redirects テキストを生成する */
export function buildRedirects(input: RedirectsInput): string {
  const lines: string[] = []

  lines.push('# Feeds → RSS')
  lines.push(...STATIC_RULES)

  if (input.categories.length > 0) {
    lines.push('')
    lines.push('# Category（日本語slug → enSlug）')
    for (const { slug, enSlug } of input.categories) {
      lines.push(`/category/${slug}/ /category/${enSlug}/ 301`)
    }
  }

  if (input.tags.length > 0) {
    lines.push('')
    lines.push('# Tag（日本語slug → enSlug）')
    for (const { slug, enSlug } of input.tags) {
      lines.push(`/tag/${slug}/ /tag/${enSlug}/ 301`)
    }
  }

  return lines.join('\n') + '\n'
}
