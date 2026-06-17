// shiimanblog.com の wp-content/uploads 画像URL。末尾にサイズ変更クエリ(?w=.. 等)が
// 付く場合も1トークンとして取り込む。
// 注: /g フラグ付きだが match()/replace() でのみ使用するため lastIndex の持続問題は起きない。
//     将来 exec()/test() をループで使う場合は都度 new RegExp すること。
const WP_UPLOAD_RE =
  /https?:\/\/shiimanblog\.com\/wp-content\/uploads\/[^\s)"'\\]+\.(?:png|jpe?g|gif|webp)(?:\?[^\s)"'\\]*)?/gi

/** クエリ文字列と -1024x545 のようなサイズ接尾辞を除去し原本URLにする */
export function toOriginalUrl(url: string): string {
  const noQuery = url.replace(/\?.*$/, '')
  return noQuery.replace(/-\d+x\d+(\.(?:png|jpe?g|gif|webp))$/i, '$1')
}

/** 本文から wp-content 画像URLを重複なく抽出する */
export function extractImageUrls(markdown: string): string[] {
  const found = markdown.match(WP_UPLOAD_RE) ?? []
  return [...new Set(found)]
}

/** WordPress ライトボックス記法 [![alt](thumb)](full) を ![alt](full) に正規化する */
export function normalizeLightbox(markdown: string): string {
  return markdown.replace(
    /\[!\[([^\]]*)\]\(([^)]+)\)\]\(([^)]+)\)/g,
    (_m, alt: string, _thumb: string, full: string) => `![${alt}](${full})`,
  )
}

/** a8.net 等の計測用gif画像記法を除去する */
export function stripTrackingPixels(markdown: string): string {
  return markdown.replace(
    /!\[[^\]]*\]\(https?:\/\/[^\s)]*a8\.net\/[^\s)]*\.gif[^)]*\)/gi,
    '',
  )
}

/** 本文中の wp-content URL（サイズ違い・クエリ付き含む）をローカル相対パスへ置換する */
export function rewriteImageUrls(markdown: string, urlToLocal: Record<string, string>): string {
  return markdown.replace(WP_UPLOAD_RE, (url) => urlToLocal[toOriginalUrl(url)] ?? url)
}

/** URL パスの末尾ファイル名を返す（node:path 非依存） */
function fileBase(pathname: string): string {
  const parts = pathname.split('/')
  return parts[parts.length - 1] ?? ''
}

/**
 * 原本URL配列から「原本URL→ローカルファイル名」のマップを作る。
 * 通常はファイル名(basename)をそのまま使うが、同一記事内で異なる原本URLが
 * 同名になる場合のみ uploads 以降のパスをフラット化して一意化する
 * （例: uploads/2021/08/photo.jpg → 2021-08-photo.jpg）。
 * 衝突が無ければ basename のままなので、既存のローカル化済みデータと出力が一致する。
 */
export function buildLocalNames(originalUrls: string[]): Record<string, string> {
  const counts = new Map<string, number>()
  for (const u of originalUrls) {
    const n = fileBase(new URL(u).pathname)
    counts.set(n, (counts.get(n) ?? 0) + 1)
  }
  const map: Record<string, string> = {}
  for (const u of originalUrls) {
    const pathname = new URL(u).pathname
    const n = fileBase(pathname)
    if ((counts.get(n) ?? 0) > 1) {
      const rel = pathname.replace(/^.*\/wp-content\/uploads\//, '')
      map[u] = rel.replace(/\//g, '-')
    } else {
      map[u] = n
    }
  }
  return map
}
