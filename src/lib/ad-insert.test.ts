import { describe, it, expect } from 'vitest'
import { insertInfeedAds } from './ad-insert'

// テスト用に記事を文字列で代用
const posts = (n: number) => Array.from({ length: n }, (_, i) => `p${i}`)

describe('insertInfeedAds', () => {
  it('空配列なら空を返す', () => {
    expect(insertInfeedAds([], 5)).toEqual([])
  })

  it('every 未満なら広告を挿入しない', () => {
    const result = insertInfeedAds(posts(3), 5)
    expect(result).toEqual([
      { kind: 'post', post: 'p0' },
      { kind: 'post', post: 'p1' },
      { kind: 'post', post: 'p2' },
    ])
  })

  it('ちょうど every 件なら末尾に広告を挿入しない', () => {
    const result = insertInfeedAds(posts(5), 5)
    expect(result.filter((x) => x.kind === 'ad')).toHaveLength(0)
    expect(result).toHaveLength(5)
  })

  it('every+1 件なら 5件目の後に広告を1つ挿入する', () => {
    const result = insertInfeedAds(posts(6), 5)
    expect(result).toEqual([
      { kind: 'post', post: 'p0' },
      { kind: 'post', post: 'p1' },
      { kind: 'post', post: 'p2' },
      { kind: 'post', post: 'p3' },
      { kind: 'post', post: 'p4' },
      { kind: 'ad', key: 0 },
      { kind: 'post', post: 'p5' },
    ])
  })

  it('12件なら5件ごとに広告を2つ挿入し、末尾には挿入しない', () => {
    const result = insertInfeedAds(posts(12), 5)
    const ads = result.filter((x) => x.kind === 'ad')
    expect(ads).toEqual([
      { kind: 'ad', key: 0 },
      { kind: 'ad', key: 1 },
    ])
    // 末尾は post（広告で終わらない）
    expect(result[result.length - 1]).toEqual({ kind: 'post', post: 'p11' })
  })
})
