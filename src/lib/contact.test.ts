import { describe, it, expect } from 'vitest'
import { validateContact, isTurnstileSuccess, buildResendPayload } from './contact'

describe('validateContact', () => {
  const valid = { name: '山田', email: 'a@example.com', subject: '件名', message: '本文', token: 'tok' }

  it('正常な入力を受理し trim する', () => {
    const r = validateContact({ ...valid, name: '  山田  ' })
    expect(r.ok).toBe(true)
    if (r.ok) expect(r.value.name).toBe('山田')
  })

  it('必須項目が欠けると拒否する', () => {
    const r = validateContact({ ...valid, name: '' })
    expect(r.ok).toBe(false)
    if (!r.ok) expect(r.errors).toContain('name')
  })

  it('メール形式が不正なら拒否する', () => {
    const r = validateContact({ ...valid, email: 'invalid' })
    expect(r.ok).toBe(false)
    if (!r.ok) expect(r.errors).toContain('email')
  })

  it('Turnstile トークンが無いと拒否する', () => {
    const r = validateContact({ ...valid, token: '' })
    expect(r.ok).toBe(false)
    if (!r.ok) expect(r.errors).toContain('token')
  })

  it('文字列以外の型は拒否する', () => {
    const r = validateContact({ ...valid, message: 12345 })
    expect(r.ok).toBe(false)
    if (!r.ok) expect(r.errors).toContain('message')
  })

  it('最大長を超えると拒否する', () => {
    const r = validateContact({ ...valid, subject: 'x'.repeat(201) })
    expect(r.ok).toBe(false)
    if (!r.ok) expect(r.errors).toContain('subject')
  })
})

describe('isTurnstileSuccess', () => {
  it('success:true で true', () => {
    expect(isTurnstileSuccess({ success: true })).toBe(true)
  })
  it('success:false で false', () => {
    expect(isTurnstileSuccess({ success: false })).toBe(false)
  })
  it('null で false', () => {
    expect(isTurnstileSuccess(null)).toBe(false)
  })
})

describe('buildResendPayload', () => {
  it('from/to/reply_to/subject/text を整形する', () => {
    const v = { name: '山田', email: 'a@example.com', subject: '件名', message: '本文', token: 'tok' }
    const p = buildResendPayload(v, { from: 'noreply@shiimanblog.com', to: 'me@example.com' })
    expect(p.from).toBe('noreply@shiimanblog.com')
    expect(p.to).toEqual(['me@example.com'])
    expect(p.reply_to).toBe('a@example.com')
    expect(p.subject).toBe('[お問い合わせ] 件名')
    expect(p.text).toContain('山田')
    expect(p.text).toContain('本文')
  })
})
