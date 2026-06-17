// 不可視トラッキング画像のホスト/パターン（画像記法のみ対象）
const TRACKING_HOSTS = ['px.a8.net', 'ad.jp.ap.valuecommerce.com', 'h.accesstrade.net']

/** 本文から不可視トラッキング画像の Markdown 記法だけを除去する */
export function stripTrackingImages(body: string): string {
  let out = body
  for (const host of TRACKING_HOSTS) {
    // ![...](https?://<host>...) を除去する。
    // 直前に改行があればそちらを除去し、直前に改行がなく直後に改行があれば後ろを除去する
    const escaped = host.replace(/[.]/g, '\\.')
    const re = new RegExp(`(\\n)?!\\[[^\\]]*\\]\\(https?:\\/\\/${escaped}[^)]*\\)(\\n)?`, 'g')
    out = out.replace(re, (_match, pre, post) => {
      if (pre) return post ?? '' // 前の改行ごと除去（後ろの改行は後続テキストに残す）
      if (post) return '' // 前に改行なし・後ろの改行を除去
      return ''
    })
  }
  return out
}
