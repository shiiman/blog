import { describe, it, expect } from 'vitest'
// @ts-expect-error - plain JS module
import { transformTree, transformLinkCard } from './rehype-link-card.mjs'

const el = (tagName: string, properties: any, children: any[] = []) => ({ type: 'element', tagName, properties, children })
const txt = (value: string) => ({ type: 'text', value })

/** 移行ブログカード相当の <a>（img + <br>区切りテキスト + favicon + ドメイン + 日付） */
function makeCardAnchor() {
  return el('a', { href: '/fire/nisa-202110/', title: 'タイトルA' }, [
    el('img', { src: '/_astro/eyecatch.webp', alt: '' }),
    el('br', {}),
    txt('\n'),
    el('br', {}),
    txt('\n【つみたてNISA】15ヶ月目の運用実績公開！| 2021年10月'),
    el('br', {}),
    txt('\n'),
    el('br', {}),
    txt('\n楽天証券でつみたてNISAを初めて、15ヶ月目の運用実績を紹介します。'),
    el('br', {}),
    el('img', { src: 'https://www.google.com/s2/favicons?domain=https://shiimanblog.com', alt: '' }),
    el('br', {}),
    txt('\nshiimanblog.com'),
    el('br', {}),
    txt('\n2021.10.06'),
  ])
}

const findClass = (node: any, cls: string): any => {
  if (node?.type === 'element' && (node.properties?.className || []).includes(cls)) return node
  for (const c of node?.children || []) {
    const r = findClass(c, cls)
    if (r) return r
  }
  return null
}
const textOf = (node: any): string =>
  node?.type === 'text' ? node.value : (node?.children || []).map(textOf).join('')

describe('transformLinkCard', () => {
  it('移行ブログカードを横型カード構造へ変換する', () => {
    const a = makeCardAnchor()
    expect(transformLinkCard(a)).toBe(true)
    expect(a.properties.className).toContain('link-card')
    expect(findClass(a, 'link-card__thumb')).toBeTruthy()
    expect(textOf(findClass(a, 'link-card__title'))).toBe('【つみたてNISA】15ヶ月目の運用実績公開！| 2021年10月')
    expect(textOf(findClass(a, 'link-card__excerpt'))).toContain('楽天証券でつみたてNISAを初めて')
    expect(textOf(findClass(a, 'link-card__meta'))).toBe('shiimanblog.com · 2021.10.06')
  })

  it('favicon に link-card__favicon クラスを付与する', () => {
    const a = makeCardAnchor()
    transformLinkCard(a)
    const fav = findClass(a, 'link-card__favicon')
    expect(fav?.tagName).toBe('img')
    expect(fav.properties.src).toContain('s2/favicons')
  })
})

describe('transformTree', () => {
  it('ツリー内のブログカードのみ変換し、通常リンクは触らない', () => {
    const normal = el('a', { href: '/foo/' }, [txt('普通のリンク')])
    const tree = el('div', {}, [el('p', {}, [makeCardAnchor()]), el('p', {}, [normal])])
    const count = transformTree(tree)
    expect(count).toBe(1)
    expect(normal.properties.className).toBeUndefined()
    expect(textOf(normal)).toBe('普通のリンク')
  })

  it('favicon を含まないリンクはカード化しない', () => {
    const a = el('a', { href: '/x/' }, [el('img', { src: '/img.png' }), txt('画像リンク')])
    const tree = el('div', {}, [a])
    expect(transformTree(tree)).toBe(0)
    expect((a.properties.className || []).includes('link-card')).toBe(false)
  })
})
