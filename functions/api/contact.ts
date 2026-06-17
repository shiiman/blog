import {
  validateContact,
  isTurnstileSuccess,
  buildResendPayload,
  type ContactInput,
  type TurnstileVerifyResponse,
} from '../../src/lib/contact'

interface Env {
  TURNSTILE_SECRET_KEY: string
  RESEND_API_KEY: string
  CONTACT_FROM_EMAIL: string
  CONTACT_TO_EMAIL: string
}

function json(data: unknown, status: number): Response {
  return new Response(JSON.stringify(data), {
    status,
    headers: { 'Content-Type': 'application/json' },
  })
}

export const onRequestPost = async (context: {
  request: Request
  env: Env
}): Promise<Response> => {
  const { request, env } = context

  let body: unknown
  try {
    body = await request.json()
  } catch {
    return json({ ok: false, message: '不正なリクエストです。' }, 400)
  }

  const result = validateContact(body as ContactInput)
  if (!result.ok) {
    return json({ ok: false, message: '入力内容をご確認ください。' }, 400)
  }

  // Turnstile 検証
  const ts = (await fetch(
    'https://challenges.cloudflare.com/turnstile/v0/siteverify',
    {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        secret: env.TURNSTILE_SECRET_KEY,
        response: result.value.token,
      }),
    }
  )
    .then((r) => r.json())
    .catch(() => null)) as TurnstileVerifyResponse | null

  if (!isTurnstileSuccess(ts)) {
    return json({ ok: false, message: '認証に失敗しました。再度お試しください。' }, 403)
  }

  // Resend 送信
  const payload = buildResendPayload(result.value, {
    from: env.CONTACT_FROM_EMAIL,
    to: env.CONTACT_TO_EMAIL,
  })
  const sendRes = await fetch('https://api.resend.com/emails', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${env.RESEND_API_KEY}`,
    },
    body: JSON.stringify(payload),
  })

  if (!sendRes.ok) {
    // 本番調査用に Resend のステータスとレスポンスボディを記録
    const errBody = await sendRes.text().catch(() => '')
    console.error('Resend error:', sendRes.status, errBody)
    return json({ ok: false, message: '送信に失敗しました。時間をおいて再度お試しください。' }, 502)
  }

  return json({ ok: true, message: 'お問い合わせを送信しました。ありがとうございます。' }, 200)
}
