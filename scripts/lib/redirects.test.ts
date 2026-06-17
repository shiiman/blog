import { describe, it, expect } from 'vitest'
import { selectRedirectTerms, buildRedirects } from './redirects'
import type { TermDict } from './taxonomy'

describe('selectRedirectTerms', () => {
  const dict: TermDict = {
    '52': { name: '節約', slug: '%e7%af%80%e7%b4%84', enSlug: 'savings' },
    '3': { name: 'WordPress', slug: 'wordpress' },
    '90': { name: '新年', slug: '%e6%96%b0%e5%b9%b4', enSlug: 'new-year' },
  }

  it('enSlug を持ち公開記事で使用中の term を抽出する', () => {
    const used = new Set(['%e7%af%80%e7%b4%84', 'wordpress'])
    expect(selectRedirectTerms(dict, used)).toEqual([{ slug: '%e7%af%80%e7%b4%84', enSlug: 'savings' }])
  })

  it('enSlug が無い term（新旧一致）は除外する', () => {
    const used = new Set(['wordpress'])
    expect(selectRedirectTerms(dict, used)).toEqual([])
  })

  it('enSlug を持つが公開記事で未使用の term は除外する', () => {
    const used = new Set(['%e7%af%80%e7%b4%84'])
    // 90:新年 は used に無いので除外、52:節約 のみ
    expect(selectRedirectTerms(dict, used)).toEqual([{ slug: '%e7%af%80%e7%b4%84', enSlug: 'savings' }])
  })
})

describe('buildRedirects', () => {
  it('静的ルール（/feed 系 → /rss.xml）を含む（スラッシュあり・なし両方）', () => {
    const text = buildRedirects({ categories: [], tags: [] })
    expect(text).toContain('/feed ')
    expect(text).toContain('/feed/')
    expect(text).toContain('/comments/feed ')
    expect(text).toContain('/comments/feed/')
    expect(text).toContain('/*/feed ')
    expect(text).toContain('/*/feed/')
    expect(text).toContain('/rss.xml')
  })

  it('カテゴリの enSlug リダイレクト行を生成する', () => {
    const text = buildRedirects({
      categories: [{ slug: '%e7%af%80%e7%b4%84', enSlug: 'savings' }],
      tags: [],
    })
    expect(text).toContain('/category/%e7%af%80%e7%b4%84/ /category/savings/ 301')
  })

  it('タグの enSlug リダイレクト行を生成する', () => {
    const text = buildRedirects({
      categories: [],
      tags: [{ slug: '%e3%83%a1%e3%83%bc%e3%83%ab', enSlug: 'mail' }],
    })
    expect(text).toContain('/tag/%e3%83%a1%e3%83%bc%e3%83%ab/ /tag/mail/ 301')
  })

  it('各リダイレクト行が <from> <to> <status> 形式で出力される', () => {
    const text = buildRedirects({
      categories: [{ slug: '%e7%af%80%e7%b4%84', enSlug: 'savings' }],
      tags: [],
    })
    const ruleLines = text.split('\n').filter((l) => l.trim() && !l.startsWith('#'))
    for (const line of ruleLines) {
      const parts = line.trim().split(/\s+/)
      expect(parts).toHaveLength(3)
      expect(parts[2]).toBe('301')
    }
  })

  it('カテゴリ・タグが空でも静的ルールだけで成立する', () => {
    const text = buildRedirects({ categories: [], tags: [] })
    const ruleLines = text.split('\n').filter((l) => l.trim() && !l.startsWith('#'))
    expect(ruleLines.length).toBeGreaterThanOrEqual(3)
  })
})
