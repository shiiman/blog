import { defineCollection, z } from 'astro:content'
import { glob } from 'astro/loaders'

// YAMLは未引用のISO日時(...Z)をDateに自動変換するため、文字列/Dateの両方を受け、
// 保存値(JST壁時計をZ表記)と同一のISO文字列へ正規化する。
// Date.toISOString()は元の瞬時を同じ文字列に戻すため、後続の日付ロジックは不変。
const dateString = z
  .union([z.string(), z.date()])
  .transform((v) => (v instanceof Date ? v.toISOString() : v))

const posts = defineCollection({
  loader: glob({ pattern: '*/article.md', base: './posts' }),
  schema: ({ image }) =>
    z.object({
      title: z.string(),
      slug: z.string(),
      date: dateString,
      modified: dateString.optional(),
      excerpt: z.string().optional(),
      categories: z.array(z.string()).default([]),
      tags: z.array(z.string()).default([]),
      eyecatch: image().optional(),
      draft: z.boolean().default(false),
      id: z.number().optional(),
    }),
})

const pages = defineCollection({
  loader: glob({ pattern: '*/page.md', base: './pages' }),
  schema: z.object({
    title: z.string(),
    slug: z.string(),
    date: dateString,
    modified: dateString.optional(),
    excerpt: z.string().optional(),
    draft: z.boolean().default(false),
    id: z.number().optional(),
  }),
})

export const collections = { posts, pages }
