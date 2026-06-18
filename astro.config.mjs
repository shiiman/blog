// @ts-check
import { defineConfig } from 'astro/config'
import sitemap from '@astrojs/sitemap'
import { unified } from '@astrojs/markdown-remark'
import rehypeLinkCard from './src/lib/rehype-link-card.mjs'

// https://astro.build/config
export default defineConfig({
  site: 'https://shiimanblog.com',
  output: 'static',
  trailingSlash: 'always',
  build: { format: 'directory' },
  integrations: [sitemap()],
  markdown: {
    // Astro v6 で markdown.rehypePlugins は非推奨。unified() プロセッサ経由で渡す。
    // syntaxHighlight(shiki)/画像最適化等は別フィールドの既定として維持される。
    processor: unified({ rehypePlugins: [rehypeLinkCard] }),
  },
})
