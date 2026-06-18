// @ts-check
import { defineConfig } from 'astro/config'
import sitemap from '@astrojs/sitemap'
import rehypeLinkCard from './src/lib/rehype-link-card.mjs'

// https://astro.build/config
export default defineConfig({
  site: 'https://shiimanblog.com',
  output: 'static',
  trailingSlash: 'always',
  build: { format: 'directory' },
  integrations: [sitemap()],
  markdown: {
    rehypePlugins: [rehypeLinkCard],
  },
})
