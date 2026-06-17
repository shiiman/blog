import { describe, it, expect } from 'vitest'
import { formatJaDate, toUtcIso } from './date'

describe('formatJaDate', () => {
  it('Z形式(JST壁時計をZ誤記)を和文日付にする', () => {
    expect(formatJaDate('2021-08-30T19:30:00.000Z')).toBe('2021年8月30日')
  })
  it('+09:00形式(正しいJST)を和文日付にする', () => {
    expect(formatJaDate('2026-01-03T00:00:00+09:00')).toBe('2026年1月3日')
  })
  it('暦フィールドをそのまま読む(深夜でも繰り上がらない)', () => {
    expect(formatJaDate('2022-01-01T00:30:00.000Z')).toBe('2022年1月1日')
  })
  it('想定外の形式はエラー', () => {
    expect(() => formatJaDate('2021/08/30')).toThrow()
  })
})

describe('toUtcIso', () => {
  it('Z形式の暦時刻をJSTとみなし正しいUTC瞬時(-9h)にする', () => {
    expect(toUtcIso('2021-08-30T19:30:00.000Z')).toBe('2021-08-30T10:30:00.000Z')
  })
  it('+09:00形式も同じ規則で正しいUTC瞬時にする(日付跨ぎ)', () => {
    expect(toUtcIso('2026-01-03T00:00:00+09:00')).toBe('2026-01-02T15:00:00.000Z')
  })
  it('想定外の形式はエラー', () => {
    expect(() => toUtcIso('2021/08/30')).toThrow()
  })
})
