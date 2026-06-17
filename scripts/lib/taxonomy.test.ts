import { describe, it, expect } from 'vitest'
import { buildTaxonomyMap, mapIdsToSlugs } from './taxonomy'

describe('buildTaxonomyMap', () => {
  it('REST配列をID文字列キーの{name,slug}マップへ変換する', () => {
    const map = buildTaxonomyMap([
      { id: 2, name: '技術', slug: 'tech' },
      { id: 10, name: 'ブログ', slug: 'blog' },
    ])
    expect(map).toEqual({
      '2': { name: '技術', slug: 'tech' },
      '10': { name: 'ブログ', slug: 'blog' },
    })
  })
})

describe('mapIdsToSlugs', () => {
  const map = { '2': { name: '技術', slug: 'tech' }, '10': { name: 'ブログ', slug: 'blog' } }

  it('ID配列をslug配列へ変換する', () => {
    expect(mapIdsToSlugs([2, 10], map)).toEqual(['tech', 'blog'])
  })

  it('空配列は空配列を返す', () => {
    expect(mapIdsToSlugs([], map)).toEqual([])
  })

  it('未知のIDはエラーを投げる', () => {
    expect(() => mapIdsToSlugs([999], map)).toThrow('未知のタクソノミーID: 999')
  })
})
