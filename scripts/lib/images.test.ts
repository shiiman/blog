import { describe, it, expect } from 'vitest'
import {
  toOriginalUrl, extractImageUrls, normalizeLightbox, stripTrackingPixels, rewriteImageUrls,
} from './images'

describe('toOriginalUrl', () => {
  it('-WxH サイズ接尾辞を除去する', () => {
    expect(toOriginalUrl('https://shiimanblog.com/wp-content/uploads/2021/08/4-1-1024x545.jpg'))
      .toBe('https://shiimanblog.com/wp-content/uploads/2021/08/4-1.jpg')
  })
  it('接尾辞が無ければそのまま', () => {
    expect(toOriginalUrl('https://shiimanblog.com/wp-content/uploads/2021/09/3.png'))
      .toBe('https://shiimanblog.com/wp-content/uploads/2021/09/3.png')
  })
})

describe('extractImageUrls', () => {
  it('wp-content の画像URLを重複なく抽出する', () => {
    const md = '![](https://shiimanblog.com/wp-content/uploads/2021/09/3.png) と ![](https://shiimanblog.com/wp-content/uploads/2021/09/3.png)'
    expect(extractImageUrls(md)).toEqual(['https://shiimanblog.com/wp-content/uploads/2021/09/3.png'])
  })
  it('外部ドメインの画像は無視する', () => {
    expect(extractImageUrls('![](https://www19.a8.net/0.gif)')).toEqual([])
  })
})

describe('normalizeLightbox', () => {
  it('[![alt](thumb)](full) を ![alt](full) に正規化する', () => {
    const md = '[![](https://shiimanblog.com/wp-content/uploads/2021/08/4-1-1024x545.jpg)](https://shiimanblog.com/wp-content/uploads/2021/08/4-1.jpg)'
    expect(normalizeLightbox(md)).toBe('![](https://shiimanblog.com/wp-content/uploads/2021/08/4-1.jpg)')
  })
})

describe('stripTrackingPixels', () => {
  it('a8.net の計測gif画像を除去する', () => {
    expect(stripTrackingPixels('text ![](https://www19.a8.net/0.gif?a8mat=XXX) end')).toBe('text  end')
  })
})

describe('rewriteImageUrls', () => {
  it('原本URLへのマップでローカルパスへ置換する（サイズ違いも原本扱い）', () => {
    const map = { 'https://shiimanblog.com/wp-content/uploads/2021/08/4-1.jpg': './assets/4-1.jpg' }
    const md = '![](https://shiimanblog.com/wp-content/uploads/2021/08/4-1-1024x545.jpg)'
    expect(rewriteImageUrls(md, map)).toBe('![](./assets/4-1.jpg)')
  })
})
