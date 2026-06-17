import { readFile, writeFile, mkdir } from 'node:fs/promises'
import { extname } from 'node:path'
import matter from 'gray-matter'
import {
  normalizeLightbox, stripTrackingPixels, extractImageUrls, toOriginalUrl, rewriteImageUrls, buildLocalNames,
} from './lib/images'
import { listContentDirs } from './lib/content-roots'

/** DL失敗したURLを蓄積するリスト（status / 記事ディレクトリ / URL のタブ区切り） */
const missingUrls: string[] = []

/** URL からファイルをダウンロードして dest へ保存する。失敗時は missingUrls に追記してスキップ */
async function download(url: string, dest: string, dir: string): Promise<boolean> {
  try {
    const res = await fetch(url)
    if (!res.ok) {
      console.warn(`  ⚠ DL失敗 (${res.status}) ${dir}: ${url}`)
      missingUrls.push(`${res.status}\t${dir}\t${url}`)
      return false
    }
    await writeFile(dest, Buffer.from(await res.arrayBuffer()))
    return true
  } catch (e) {
    const msg = e instanceof Error ? e.message : String(e)
    console.warn(`  ⚠ DLエラー ${dir}: ${url} — ${msg}`)
    missingUrls.push(`ERR\t${dir}\t${url}`)
    return false
  }
}

async function processOne(dir: string, file: string, featured: Record<string, string>): Promise<void> {
  const raw = await readFile(`${dir}/${file}`, 'utf8')
  const { data, content } = matter(raw)
  let body = stripTrackingPixels(normalizeLightbox(content))

  const urls = [...new Set(extractImageUrls(body).map(toOriginalUrl))]
  const names = buildLocalNames(urls)
  const map: Record<string, string> = {}
  if (urls.length > 0) await mkdir(`${dir}/assets`, { recursive: true })
  for (const url of urls) {
    const name = names[url]
    const ok = await download(url, `${dir}/assets/${name}`, dir)
    if (ok) map[url] = `./assets/${name}`
    // DL失敗時は map に登録しない → 本文の URL はそのまま残る
  }
  body = rewriteImageUrls(body, map)

  const fmId = typeof data.featured_media === 'number' ? data.featured_media : 0
  const eyeUrl = featured[String(fmId)]
  if (eyeUrl) {
    const ext = extname(new URL(eyeUrl).pathname) || '.png'
    await mkdir(`${dir}/assets`, { recursive: true })
    await download(toOriginalUrl(eyeUrl), `${dir}/assets/eyecatch${ext}`, dir)
  }

  await writeFile(`${dir}/${file}`, matter.stringify(body, data))
}

async function main() {
  const featured = JSON.parse(await readFile('data/featured-media.json', 'utf8')) as Record<string, string>
  for (const { dir, file } of await listContentDirs()) {
    console.log(`localize: ${dir}`)
    await processOne(dir, file, featured)
  }

  // DL失敗URLを出力
  if (missingUrls.length > 0) {
    await mkdir('data', { recursive: true })
    await writeFile('data/missing-images.txt', missingUrls.join('\n') + '\n')
    console.warn(`\n⚠ ${missingUrls.length}件のDL失敗 → data/missing-images.txt に記録しました`)
  } else {
    console.log('\n✓ DL失敗なし')
  }
}

main().catch((e) => {
  console.error(e)
  process.exit(1)
})
