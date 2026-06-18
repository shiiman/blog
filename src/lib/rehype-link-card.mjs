/**
 * rehype プラグイン: WordPress 移行で「画像 + 改行区切りテキスト」の塊に
 * なってしまった内部リンクのブログカードを、整った横型リンクカードへ変換する。
 *
 * 判定: <a> の子孫に src が `s2/favicons` を含む <img>（移行カードの目印）を持つもの。
 * 通常のリンク・画像リンクには一切手を加えない。
 *
 * 変換後:
 *   <a class="link-card" href title>
 *     <span class="link-card__thumb"><img(eyecatch)></span>
 *     <span class="link-card__body">
 *       <span class="link-card__title">タイトル</span>
 *       <span class="link-card__excerpt">要約</span>
 *       <span class="link-card__meta"><img class="link-card__favicon">ドメイン · 日付</span>
 *     </span>
 *   </a>
 */

const FAVICON_HINT = 's2/favicons'
const DATE_RE = /^\d{4}[./-]\d{1,2}[./-]\d{1,2}$/
const DOMAIN_RE = /^[\w-]+(\.[\w-]+)+$/

const el = (tagName, className, children) => ({
  type: 'element',
  tagName,
  properties: className ? { className: Array.isArray(className) ? className : [className] } : {},
  children: children || [],
})
const txt = (value) => ({ type: 'text', value })

/** 部分木内の <img> 要素をすべて集める */
function collectImgs(node, acc = []) {
  if (!node || typeof node !== 'object') return acc
  if (node.type === 'element' && node.tagName === 'img') acc.push(node)
  for (const c of node.children || []) collectImgs(c, acc)
  return acc
}

/** 部分木内のテキストノードの値を連結する */
function collectText(node, parts = []) {
  if (!node || typeof node !== 'object') return parts
  if (node.type === 'text') parts.push(node.value)
  for (const c of node.children || []) collectText(c, parts)
  return parts
}

const imgSrc = (img) => String((img.properties && img.properties.src) || '')
const isLinkCard = (a) => collectImgs(a).some((img) => imgSrc(img).includes(FAVICON_HINT))

/**
 * ブログカードの <a> 要素を整形カードへ書き換える。
 * 構造が想定外（行が少なすぎる等）の場合は false を返して変換しない。
 */
export function transformLinkCard(a) {
  const imgs = collectImgs(a)
  const favicon = imgs.find((img) => imgSrc(img).includes(FAVICON_HINT))
  const eyecatch = imgs.find((img) => img !== favicon)

  const lines = collectText(a)
    .join('\n')
    .split('\n')
    .map((s) => s.trim())
    .filter(Boolean)
  if (lines.length < 2) return false

  const rest = [...lines]
  let date = ''
  if (DATE_RE.test(rest[rest.length - 1])) date = rest.pop()
  let domain = ''
  if (rest.length && DOMAIN_RE.test(rest[rest.length - 1])) domain = rest.pop()
  const title = rest.shift() || String((a.properties && a.properties.title) || '')
  const excerpt = rest.join(' ')
  if (!title) return false

  if (favicon) {
    favicon.properties = { ...favicon.properties, className: ['link-card__favicon'], alt: '', loading: 'lazy' }
  }

  const meta = []
  if (favicon) meta.push(favicon)
  const metaText = [domain, date].filter(Boolean).join(' · ')
  if (metaText) meta.push(txt(metaText))

  const bodyChildren = [el('span', 'link-card__title', [txt(title)])]
  if (excerpt) bodyChildren.push(el('span', 'link-card__excerpt', [txt(excerpt)]))
  if (meta.length) bodyChildren.push(el('span', 'link-card__meta', meta))

  const children = []
  if (eyecatch) children.push(el('span', 'link-card__thumb', [eyecatch]))
  children.push(el('span', 'link-card__body', bodyChildren))

  a.children = children
  const existing = (a.properties && a.properties.className) || []
  a.properties = {
    ...a.properties,
    className: [...(Array.isArray(existing) ? existing : [existing]), 'link-card'],
  }
  return true
}

/** hast を走査して全ブログカードを変換する（テスト用に分離） */
export function transformTree(tree) {
  let count = 0
  const walk = (node) => {
    if (!node || typeof node !== 'object') return
    if (node.type === 'element' && node.tagName === 'a' && isLinkCard(node)) {
      if (transformLinkCard(node)) count++
      return // カード内部はこれ以上掘らない
    }
    for (const c of node.children || []) walk(c)
  }
  walk(tree)
  return count
}

/** rehype プラグイン本体 */
export default function rehypeLinkCard() {
  return (tree) => {
    transformTree(tree)
  }
}
