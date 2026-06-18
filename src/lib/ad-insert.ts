/** 記事一覧の各要素（記事 or インフィード広告） */
export type AdListItem<T> = { kind: 'post'; post: T } | { kind: 'ad'; key: number }

/**
 * 記事リストに every 件ごとにインフィード広告マーカーを挿入する。
 * 末尾（最後の記事の直後）には挿入しない。
 */
export function insertInfeedAds<T>(posts: T[], every: number): AdListItem<T>[] {
  const result: AdListItem<T>[] = []
  let adKey = 0
  posts.forEach((post, i) => {
    result.push({ kind: 'post', post })
    const pos = i + 1
    if (pos % every === 0 && pos < posts.length) {
      result.push({ kind: 'ad', key: adKey++ })
    }
  })
  return result
}
