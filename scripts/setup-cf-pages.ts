/**
 * Cloudflare Pages プロジェクトを API で作成する（初回のみ）
 * 実行: npm run setup:cf-pages
 */
const token = process.env.CLOUDFLARE_API_TOKEN
const accountId = process.env.CLOUDFLARE_ACCOUNT_ID

if (!token) throw new Error('CLOUDFLARE_API_TOKEN が未設定です')
if (!accountId) throw new Error('CLOUDFLARE_ACCOUNT_ID が未設定です')

const res = await fetch(
  `https://api.cloudflare.com/client/v4/accounts/${accountId}/pages/projects`,
  {
    method: 'POST',
    headers: {
      Authorization: `Bearer ${token}`,
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      name: 'shiimanblog',
      production_branch: 'main',
    }),
  },
)

const json = (await res.json()) as { success: boolean; errors: { message: string }[]; result?: { subdomain: string } }

if (!json.success) {
  const msg = json.errors.map((e) => e.message).join(', ')
  // 既に存在する場合はスキップ
  if (msg.includes('already exists') || msg.includes('taken')) {
    console.log('✅ プロジェクト "shiimanblog" は既に存在します（スキップ）')
    process.exit(0)
  }
  throw new Error(`Cloudflare API エラー: ${msg}`)
}

console.log(`✅ Pages プロジェクト作成完了`)
console.log(`   プレビューURL: https://shiimanblog.pages.dev`)
