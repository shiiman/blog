// 保存値は2形式が混在する:
//   - "2021-08-30T19:30:00.000Z"  … JST壁時計をZと誤記した値（移行由来, 80件）
//   - "2026-01-03T00:00:00+09:00" … 正しいJST（タイムゾーン付き, 5件）
// どちらも「文字列に書かれた暦時刻」がJST壁時計なので、tz指定子を無視して
// 暦フィールド(年月日時分秒)を読み、JST(=UTC+9)として扱う。
// （new Date() 経由だと +09:00 はオフセット解釈されて暦日がずれるため使わない）
// フロントマター(スキーマ検証済み)の値のみ入力される前提。先頭の暦フィールドだけ読む
const DATE_RE = /^(\d{4})-(\d{2})-(\d{2})T(\d{2}):(\d{2}):(\d{2})/
const JST_OFFSET_MS = 9 * 60 * 60 * 1000

interface DateParts {
  y: number
  mo: number
  d: number
  h: number
  mi: number
  s: number
}

function parseParts(value: string): DateParts {
  const m = DATE_RE.exec(value)
  if (!m) throw new Error(`想定外の日付形式: ${value}`)
  return { y: +m[1], mo: +m[2], d: +m[3], h: +m[4], mi: +m[5], s: +m[6] }
}

/** JST壁時計の暦日を和文で返す（例「2021年8月30日」） */
export function formatJaDate(value: string): string {
  const { y, mo, d } = parseParts(value)
  return `${y}年${mo}月${d}日`
}

/** JST壁時計を正しいUTC瞬時(-9h)へ補正したISO文字列を返す（RSS/sitemap/<time>用） */
export function toUtcIso(value: string): string {
  const { y, mo, d, h, mi, s } = parseParts(value)
  const utcMs = Date.UTC(y, mo - 1, d, h, mi, s) - JST_OFFSET_MS
  return new Date(utcMs).toISOString()
}
