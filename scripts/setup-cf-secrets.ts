/**
 * Cloudflare Pages 本番環境のシークレットを API で一括投入する
 * 実行: npm run setup:cf-secrets
 *
 * 事前準備: .prd.vars.example をコピーして .prd.vars を作成し実値を入力する
 */
const token = process.env.CLOUDFLARE_API_TOKEN
const accountId = process.env.CLOUDFLARE_ACCOUNT_ID

if (!token) throw new Error('CLOUDFLARE_API_TOKEN が未設定です（.env を確認）')
if (!accountId) throw new Error('CLOUDFLARE_ACCOUNT_ID が未設定です（.env を確認）')

const REQUIRED = [
  'TURNSTILE_SECRET_KEY',
  'RESEND_API_KEY',
  'CONTACT_FROM_EMAIL',
  'CONTACT_TO_EMAIL',
] as const

// .prd.vars から読み込まれた値を確認
const missing = REQUIRED.filter((k) => !process.env[k])
if (missing.length > 0) {
  throw new Error(`以下のキーが未設定です（.prd.vars を確認）:\n  ${missing.join('\n  ')}`)
}

const envVars: Record<string, { type: 'secret_text'; value: string }> = {}
for (const key of REQUIRED) {
  envVars[key] = { type: 'secret_text', value: process.env[key]! }
}

const res = await fetch(
  `https://api.cloudflare.com/client/v4/accounts/${accountId}/pages/projects/shiimanblog`,
  {
    method: 'PATCH',
    headers: {
      Authorization: `Bearer ${token}`,
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      deployment_configs: {
        production: { env_vars: envVars },
      },
    }),
  },
)

const json = (await res.json()) as { success: boolean; errors: { message: string }[] }

if (!json.success) {
  const msg = json.errors.map((e) => e.message).join(', ')
  throw new Error(`Cloudflare API エラー: ${msg}`)
}

console.log('✅ 本番シークレットを Cloudflare Pages に投入しました')
console.log('   再デプロイで反映されます（main への push またはダッシュボードから手動デプロイ）')
