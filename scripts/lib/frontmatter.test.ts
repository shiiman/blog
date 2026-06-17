import { describe, it, expect } from 'vitest'
import { transformFrontmatter } from './frontmatter'

const catMap = { '2': { name: '技術', slug: 'tech' } }
const tagMap = { '6': { name: 'ブログ', slug: 'blog' }, '7': { name: 'SEO', slug: 'seo' } }

describe('transformFrontmatter', () => {
  const base = {
    id: 6,
    title: 'タイトル',
    slug: 'start-blog',
    status: 'publish',
    date: '2021-08-30T19:30:00',
    modified: '2021-08-30T19:22:28',
    excerpt: '概要',
    categories: [2],
    tags: [6, 7],
    featured_media: 13,
  }

  it('カテゴリ/タグをslugへ変換しdraftを設定する', () => {
    const out = transformFrontmatter(base, catMap, tagMap)
    expect(out.categories).toEqual(['tech'])
    expect(out.tags).toEqual(['blog', 'seo'])
    expect(out.draft).toBe(false)
    expect(out.id).toBe(6)
    expect(out.title).toBe('タイトル')
  })

  it('status=draft は draft:true になる', () => {
    expect(transformFrontmatter({ ...base, status: 'draft' }, catMap, tagMap).draft).toBe(true)
  })

  it('status と featured_media は出力に残さない', () => {
    const out = transformFrontmatter(base, catMap, tagMap) as unknown as Record<string, unknown>
    expect(out.status).toBeUndefined()
    expect(out.featured_media).toBeUndefined()
  })

  it('eyecatchFile が渡されると相対パスを設定する', () => {
    const out = transformFrontmatter(base, catMap, tagMap, 'eyecatch.png')
    expect(out.eyecatch).toBe('./assets/eyecatch.png')
  })

  it('eyecatchFile 無しなら eyecatch は未設定', () => {
    expect(transformFrontmatter(base, catMap, tagMap).eyecatch).toBeUndefined()
  })

  it('categories/tags が無い（固定ページ）場合は空配列', () => {
    const page = { title: 'プロフィール', slug: 'profile', status: 'publish', date: '2021-08-31T02:29:25' }
    const out = transformFrontmatter(page, catMap, tagMap)
    expect(out.categories).toEqual([])
    expect(out.tags).toEqual([])
  })
})
