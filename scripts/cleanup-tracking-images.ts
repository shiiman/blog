import { readFile, writeFile } from 'node:fs/promises'
import { listContentDirs } from './lib/content-roots'
import { stripTrackingImages } from './lib/cleanup-tracking'

async function main() {
  const dirs = await listContentDirs()
  let changed = 0
  for (const { dir, file } of dirs) {
    const path = `${dir}/${file}`
    const raw = await readFile(path, 'utf8')
    const next = stripTrackingImages(raw)
    if (next !== raw) {
      await writeFile(path, next)
      changed++
      console.log(`cleaned: ${path}`)
    }
  }
  console.log(`changed files: ${changed}`)
}

main().catch((e) => {
  console.error(e)
  process.exit(1)
})
