import { describe, it, expect } from 'vitest'
import { stripTrackingImages } from './cleanup-tracking'

describe('stripTrackingImages', () => {
  it('a8.netの計測gif画像記法を除去する', () => {
    const input = 'before\n![](https://px.a8.net/svt/ejp?a8mat=ABC)\nafter'
    expect(stripTrackingImages(input)).toBe('before\nafter')
  })
  it('valuecommerceのgifbannerを除去する', () => {
    const input = '![](https://ad.jp.ap.valuecommerce.com/servlet/gifbanner?sid=1&pid=2)\n本文'
    expect(stripTrackingImages(input)).toBe('本文')
  })
  it('accesstradeのrrを除去する', () => {
    const input = 'x ![](https://h.accesstrade.net/sp/rr?rk=01001aqe00lqea) y'
    expect(stripTrackingImages(input)).toBe('x  y')
  })
  it('通常の画像やリンクは保持する', () => {
    const input = '![alt](./assets/foo.png)\n[a8リンク](https://px.a8.net/abc)'
    expect(stripTrackingImages(input)).toBe(input)
  })
  it('プロトコル相対URLのトラッキング画像も除去しリンクは保持する', () => {
    const input = '[![](//ad.jp.ap.valuecommerce.com/servlet/gifbanner?sid=1&pid=2)バリューコマース](//ck.jp.ap.valuecommerce.com/servlet/referral?sid=1&pid=2)'
    expect(stripTrackingImages(input)).toBe('[バリューコマース](//ck.jp.ap.valuecommerce.com/servlet/referral?sid=1&pid=2)')
  })
})
