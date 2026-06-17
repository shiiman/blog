import { describe, it, expect } from 'vitest'
import { resolveTaxonomy, type TaxonomyMap } from './taxonomy'

const map: TaxonomyMap = {
  '14': { name: 'FIRE', slug: 'fire' },
  '52': { name: '節約', slug: '%e7%af%80%e7%b4%84', enSlug: 'savings' },
}

describe('resolveTaxonomy', () => {
  it('英語slugはそのままurlSlugになり、表示名はname', () => {
    expect(resolveTaxonomy('fire', map)).toEqual({ name: 'FIRE', urlSlug: 'fire' })
  })
  it('enSlugがあればurlSlugに使う(表示名は和名)', () => {
    expect(resolveTaxonomy('%e7%af%80%e7%b4%84', map)).toEqual({ name: '節約', urlSlug: 'savings' })
  })
  it('未知slugは防御的にそのまま返す', () => {
    expect(resolveTaxonomy('unknown', map)).toEqual({ name: 'unknown', urlSlug: 'unknown' })
  })
})
