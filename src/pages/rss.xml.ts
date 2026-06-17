import rss from '@astrojs/rss'
import type { APIContext } from 'astro'
import { getPublishedPosts } from '../lib/posts'
import { toUtcIso } from '../lib/date'

export async function GET(context: APIContext) {
  const posts = await getPublishedPosts()
  return rss({
    title: 'shiimanblog',
    description: 'ゲーム開発エンジニアしーまんの技術・副業・投資ブログ',
    site: context.site!,
    items: posts.map((p) => ({
      title: p.entry.data.title,
      link: p.path,
      pubDate: new Date(toUtcIso(p.entry.data.date)),
      description: p.entry.data.excerpt ?? '',
    })),
  })
}
