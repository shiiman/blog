// 問い合わせの純粋ロジック（Functions から再利用するため astro: 系 import を含めない）

export interface ContactInput {
  name?: unknown
  email?: unknown
  subject?: unknown
  message?: unknown
  token?: unknown
}

export interface ValidatedContact {
  name: string
  email: string
  subject: string
  message: string
  token: string
}

export type ValidationResult =
  | { ok: true; value: ValidatedContact }
  | { ok: false; errors: string[] }

export interface TurnstileVerifyResponse {
  success: boolean
  [key: string]: unknown
}

export interface ResendPayload {
  from: string
  to: string[]
  reply_to: string
  subject: string
  text: string
}

const EMAIL_RE = /^[^\s@]+@[^\s@]+\.[^\s@]+$/

function asString(v: unknown): string {
  return typeof v === 'string' ? v.trim() : ''
}

export function validateContact(input: ContactInput): ValidationResult {
  const errors: string[] = []
  const name = asString(input.name)
  const email = asString(input.email)
  const subject = asString(input.subject)
  const message = asString(input.message)
  const token = typeof input.token === 'string' ? input.token : ''

  if (!name || name.length > 100) errors.push('name')
  if (!email || email.length > 200 || !EMAIL_RE.test(email)) errors.push('email')
  if (!subject || subject.length > 200) errors.push('subject')
  if (!message || message.length > 5000) errors.push('message')
  if (!token) errors.push('token')

  if (errors.length > 0) return { ok: false, errors }
  return { ok: true, value: { name, email, subject, message, token } }
}

export function isTurnstileSuccess(res: TurnstileVerifyResponse | null): boolean {
  return Boolean(res && res.success === true)
}

export function buildResendPayload(
  v: ValidatedContact,
  opts: { from: string; to: string }
): ResendPayload {
  return {
    from: opts.from,
    to: [opts.to],
    reply_to: v.email,
    subject: `[お問い合わせ] ${v.subject}`,
    text: `お名前: ${v.name}\nメールアドレス: ${v.email}\n\n${v.message}`,
  }
}
