import { readdir } from 'node:fs/promises'

/** 変換対象のコンテンツルート。 */
export const CONTENT_ROOTS = [
  { base: 'posts', file: 'article.md' },
  { base: 'pages', file: 'page.md' },
] as const

/** 各ルート配下のディレクトリ名を列挙して {dir,file} を返す */
export async function listContentDirs(): Promise<{ dir: string; file: string }[]> {
  const result: { dir: string; file: string }[] = []
  for (const root of CONTENT_ROOTS) {
    const entries = await readdir(root.base, { withFileTypes: true })
    for (const e of entries) {
      if (e.isDirectory()) result.push({ dir: `${root.base}/${e.name}`, file: root.file })
    }
  }
  return result
}
