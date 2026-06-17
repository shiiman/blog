export const PAGE_SIZE = 10

/** 総ページ数（0件でも1ページ） */
export function totalPages(count: number, size: number): number {
  return Math.max(1, Math.ceil(count / size))
}

/** 1始まりのページ番号でスライスを返す */
export function pageSlice<T>(items: T[], page: number, size: number): T[] {
  const start = (page - 1) * size
  return items.slice(start, start + size)
}
