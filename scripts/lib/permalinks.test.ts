import { describe, it, expect } from 'vitest'
import { linkToPath, buildPermalinkMap } from './permalinks'

describe('linkToPath', () => {
  it('絶対URLからパス部分(末尾スラッシュ付)を取り出す', () => {
    expect(linkToPath('https://shiimanblog.com/wordpress/conoha-wing/')).toBe('/wordpress/conoha-wing/')
  })
  it('クエリ・ハッシュを除去する', () => {
    expect(linkToPath('https://shiimanblog.com/profile/start-blog/?utm=x#h')).toBe('/profile/start-blog/')
  })
  it('末尾スラッシュが無ければ補う', () => {
    expect(linkToPath('https://shiimanblog.com/contact')).toBe('/contact/')
  })
})

describe('buildPermalinkMap', () => {
  it('id文字列キーで {slug, path} を作る', () => {
    const items = [
      { id: 6, slug: 'start-blog', link: 'https://shiimanblog.com/profile/start-blog/' },
      { id: 34, slug: 'conoha-wing', link: 'https://shiimanblog.com/wordpress/conoha-wing/' },
    ]
    expect(buildPermalinkMap(items)).toEqual({
      '6': { slug: 'start-blog', path: '/profile/start-blog/' },
      '34': { slug: 'conoha-wing', path: '/wordpress/conoha-wing/' },
    })
  })
})
