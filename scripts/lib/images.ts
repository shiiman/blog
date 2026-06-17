const WP_UPLOAD_RE =
  /https?:\/\/shiimanblog\.com\/wp-content\/uploads\/[^\s)"'\\]+\.(?:png|jpe?g|gif|webp)/gi

/** -1024x545 のようなサイズ接尾辞を除去し原本URLにする */
export function toOriginalUrl(url: string): string {
  return url.replace(/-\d+x\d+(\.(?:png|jpe?g|gif|webp))$/i, '$1')
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

/** 本文中の wp-content URL（サイズ違い含む）をローカル相対パスへ置換する */
export function rewriteImageUrls(markdown: string, urlToLocal: Record<string, string>): string {
  return markdown.replace(WP_UPLOAD_RE, (url) => urlToLocal[toOriginalUrl(url)] ?? url)
}
