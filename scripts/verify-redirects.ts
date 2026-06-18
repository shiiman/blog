/**
 * 旧サイトの全URLが新サイトで 200 または 301 で解決することを検証する。
 *
 * 使い方:
 *   npx tsx scripts/verify-redirects.ts
 *
 * 事前条件:
 *   - npm run build 済み（dist/ が存在する）
 *   - data/old-sitemap-1.xml が存在する（Phase 0 で取得済み）
 *   - public/_redirects が生成済み（npm run build:redirects 実行済み）
 */
import { readFile } from 'node:fs/promises'
import { existsSync } from 'node:fs'

const OLD_SITEMAP = 'data/old-sitemap-1.xml'
const DIST_DIR = 'dist'
const REDIRECTS_FILE = 'public/_redirects'

/** XMLのsitemap-1.xmlから <loc> URLを抽出する */
function extractUrls(xml: string): string[] {
  const urls: string[] = []
  for (const m of xml.matchAll(/<loc>(https?:\/\/[^<]+)<\/loc>/g)) {
    urls.push(m[1].trim())
  }
  return urls
}

/** URLパスを dist の静的ファイルパスへ変換して存在確認 */
function existsInDist(urlPath: string): boolean {
  const clean = urlPath.replace(/^\//, '').replace(/\/$/, '')
  if (clean === '') return existsSync(`${DIST_DIR}/index.html`)
  return (
    existsSync(`${DIST_DIR}/${clean}/index.html`) ||
    existsSync(`${DIST_DIR}/${clean}.html`)
  )
}

interface ParsedRule {
  from: string
  to: string
}

/** _redirects のルール行をパースする */
function parseRedirectRules(text: string): ParsedRule[] {
  return text
    .split('\n')
    .map((l) => l.trim())
    .filter((l) => l && !l.startsWith('#'))
    .map((l) => l.split(/\s+/))
    .filter((parts) => parts.length >= 2)
    .map(([from, to]) => ({ from, to }))
}

/**
 * URLパスが _redirects ルールにマッチするか確認する。
 * 完全一致 と splat(*) パターンのみ対応（Cloudflare _redirects の主要パターン）
 */
function matchesRedirects(urlPath: string, rules: ParsedRule[]): boolean {
  for (const { from } of rules) {
    if (from === urlPath) return true
    if (from.includes('*')) {
      // /*/foo/ の splat をRegexへ変換
      const pattern = '^' + from.replace(/[.+?^${}()|[\]\\]/g, '\\$&').replace(/\*/g, '.*') + '$'
      if (new RegExp(pattern).test(urlPath)) return true
    }
  }
  return false
}

async function main() {
  if (!existsSync(OLD_SITEMAP)) {
    console.error(`${OLD_SITEMAP} が見つかりません。Phase 0 でコミット済みのはずです（git checkout を確認してください）。`)
    process.exit(1)
  }
  if (!existsSync(DIST_DIR)) {
    console.error(`${DIST_DIR} が見つかりません。npm run build を先に実行してください。`)
    process.exit(1)
  }

  const xml = await readFile(OLD_SITEMAP, 'utf-8')
  const urls = extractUrls(xml)

  let rules: ParsedRule[] = []
  if (existsSync(REDIRECTS_FILE)) {
    const rText = await readFile(REDIRECTS_FILE, 'utf-8')
    rules = parseRedirectRules(rText)
  }

  const missing: string[] = []
  for (const url of urls) {
    const urlPath = new URL(url).pathname.replace(/%[0-9a-f]{2}/gi, (m) => m.toUpperCase())
    if (!existsInDist(urlPath) && !matchesRedirects(urlPath, rules)) {
      missing.push(url)
    }
  }

  console.log(`旧URL総数: ${urls.length}`)
  if (missing.length === 0) {
    console.log('✅ 全URL: dist または _redirects でカバーされています')
  } else {
    console.log(`❌ 未カバーURL: ${missing.length}件`)
    for (const u of missing) console.log(`  ${u}`)
    process.exit(1)
  }
}

main().catch((e) => {
  console.error(e)
  process.exit(1)
})
