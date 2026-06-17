import { describe, it, expect } from 'vitest'
import { totalPages, pageSlice, PAGE_SIZE } from './pagination'

describe('totalPages', () => {
  it('件数とサイズから総ページ数を求める', () => {
    expect(totalPages(85, 10)).toBe(9)
    expect(totalPages(0, 10)).toBe(1)
    expect(totalPages(10, 10)).toBe(1)
    expect(totalPages(11, 10)).toBe(2)
  })
})

describe('pageSlice', () => {
  it('1始まりのページ番号で該当範囲を切り出す', () => {
    const items = Array.from({ length: 25 }, (_, i) => i)
    expect(pageSlice(items, 1, 10)).toEqual(items.slice(0, 10))
    expect(pageSlice(items, 3, 10)).toEqual(items.slice(20, 25))
  })
})

describe('PAGE_SIZE', () => {
  it('WP既定に合わせ10', () => {
    expect(PAGE_SIZE).toBe(10)
  })
})
