# 計画3 設計書: 動的3機能の無料実装（検索・コメント・問い合わせ）

- 作成日: 2026-06-17
- 親設計書: [2026-06-17-blog-migration-astro-cloudflare-design.md](./2026-06-17-blog-migration-astro-cloudflare-design.md)（§8）
- 前提設計書: [2026-06-17-blog-migration-plan2-astro-site-design.md](./2026-06-17-blog-migration-plan2-astro-site-design.md)
- 位置づけ: 移行プロジェクトの「計画3」。計画1（データ移行レイヤー）・計画2（Astro静的サイト）はマージ済み。本計画で**検索・コメント・問い合わせ**の動的3機能をすべて無料枠で実装する
- 対象サイト: shiimanblog.com

---

## 1. ゴールと完了条件

### ゴール
計画2で完成した静的サイトに、無料枠の動的3機能を追加する。デザインは既存のモダン・ミニマル＋テーマ連動（ライト/ダーク/システム）を踏襲する。

### 完了条件（このフェーズの Done）
1. `npm run build`（`astro build` + Pagefindインデックス生成）がローカルで成功する
2. **検索**: ヘッダーのオーバーレイ検索が表示され、記事を全文検索できる（`npm run preview` で確認）
3. **コメント**: 記事ページに giscus が表示され、テーマ切替に連動する
4. **問い合わせ**: `/contact/` が実フォーム化され、`wrangler pages dev` でローカル検証（入力検証・Turnstile・Resend送信フロー）が動作する
5. `npm test`（vitest）が通る（既存43件を維持＋問い合わせ純粋ロジックのテストを追加）

### スコープ外（後続フェーズ＝計画4）
- 本番デプロイ・DNS切替・ConoHa解約
- ResendのDNSドメイン認証（独自ドメイン）・本番シークレット投入
- `public/_redirects`（301、`/feed/`→`/rss.xml` 含む）
- giscus の本番動作確認（appインストール後の本番疎通）

---

## 2. 確定事項（計画3ブレインストーミングでの決定）

| # | 論点 | 決定 |
|---|---|---|
| 1 | 検索UIの配置 | **ヘッダーにアイコン → クリックでオーバーレイ（モーダル）検索**。全ページから利用可・URL追加なし。Pagefind の JS API を自前の軽量UIで駆動し、既存CSS変数でテーマ連動 |
| 2 | giscus 保管先 | **`shiiman/blog`（既にpublic）で Discussions を有効化**して使用。`Announcements` 形式カテゴリにコメントをホスティング |
| 3 | giscus 設定値の管理 | repo/repoId/category/categoryId（いずれも公開値）を **`PUBLIC_` 環境変数**で管理 |
| 4 | Resend 送信元 | **独自ドメイン `shiimanblog.com` を認証**（`noreply@shiimanblog.com` 等）。DNS認証は計画4で実施し、計画3は環境変数プレースホルダで実装・ローカル検証 |
| 5 | 問い合わせのテスト | **純粋ロジックの vitest のみ**（入力検証・Turnstile/Resendレスポンス処理）。外部呼び出しはモック。wrangler 統合テストは自動化しない（手動ローカル検証） |
| 6 | Astro アダプタ | **不要**。静的サイトでは `@astrojs/cloudflare` は不要（SSR専用）。`functions/` 配下の Pages Functions は静的出力と独立に動作する |

### 調査で確認済みの事実（権威: 公式ドキュメント, 2026-06-17 時点）
- **アダプタ**: 「If you're using Astro as a static site builder, you don't need an adapter.」`@astrojs/cloudflare` はオンデマンドレンダリング（SSR・server islands・actions・sessions）専用。本計画は静的出力のままで可
- **Pagefind**: `pagefind --site dist` でビルド成果物をインデックス化。検索は `const pagefind = await import("/pagefind/pagefind.js"); const search = await pagefind.search("...")` の JS API で実行でき、公式UI（pagefind-ui）に依存せず自前UIを構築可能
- **giscus**: テーマは `<script>` の `data-theme` 属性指定に加え、`iframe.contentWindow.postMessage({ giscus: { setConfig: { theme } } }, 'https://giscus.app')` で動的切替が可能。リポジトリは **public 必須**（`shiiman/blog` は条件を満たす）

---

## 3. 全体方針

- **`astro.config.mjs` は変更しない**（`output: 'static'` 維持）。アダプタも追加しない
- Pages Functions は `functions/` ディレクトリ方式（`onRequestPost` 等を export）で実装。Astro のルーティングとは独立
- **シークレットは環境変数経由**。`.env` は読まない・出力しない（CLAUDE.md 準拠）
  - 公開値（クライアント埋め込み）= Astro の `PUBLIC_` プレフィックス
  - 秘密値（Functions のみ参照）= Cloudflare Pages 環境変数（暗号化）。ローカルは `.dev.vars`（gitignore）
  - `*.example` ファイルにはキー名のみ記載（実値は書かない）
- `main` から新ブランチを切り、**タスクごとにコミット**

---

## 4. 機能1: 検索（Pagefind）

### ビルド統合
- `package.json` の `build` を `astro build && pagefind --site dist` に変更
- `pagefind` を devDependencies に追加（npm パッケージ版CLI）
- 生成物 `dist/pagefind/` は `.gitignore` に追加（ビルドの出力でありソース管理不要）

### UI（ヘッダーオーバーレイ）
- `src/components/SearchOverlay.astro`（新規）: モーダル本体（検索入力＋結果リスト）。初期状態は非表示
- `src/components/Header.astro`（改修）: ナビに検索アイコンボタンを追加。クリックでオーバーレイを開く
- クライアントスクリプト: 初回オープン時に `/pagefind/pagefind.js` を動的 import → `pagefind.search(query)` → 各結果の `.data()`（`url`/`meta.title`/`excerpt`）を描画
- 操作性: `Esc` で閉じる、背景クリックで閉じる、開いたら入力にフォーカス、デバウンス入力
- スタイル: 既存 `src/styles/global.css` の CSS 変数（`--bg`/`--text`/`--text-muted`/`--border`/`--accent`）でライト/ダーク両対応

### 開発時の扱い
- `astro dev` では `dist/pagefind/` が存在せず検索は動作しない（既知の挙動）
- 検証は `npm run build && npm run preview`（または `wrangler pages dev dist`）で行う
- import 失敗時はUIに「検索はビルド後に有効になります」等のフォールバック表示

---

## 5. 機能2: コメント（giscus）

### 事前準備（ユーザー手作業 / §12 に手順）
1. `shiiman/blog` の Discussions を有効化
2. `Announcements` 形式のカテゴリを作成（例: `Comments`）
3. https://github.com/apps/giscus を `shiiman/blog` にインストール
4. https://giscus.app で repo/repoId/category/categoryId を取得 → `.env`（ローカル）と Cloudflare 環境変数に設定

### 実装
- `src/components/Comments.astro`（新規）: giscus の `<script>` を埋め込む。`mapping="pathname"`（記事URLでスレッド対応）、`reactions-enabled`、`input-position="top"`、`lang="ja"`
- 設定値は `PUBLIC_GISCUS_REPO` / `PUBLIC_GISCUS_REPO_ID` / `PUBLIC_GISCUS_CATEGORY` / `PUBLIC_GISCUS_CATEGORY_ID` から注入
- `src/layouts/PostLayout.astro`（改修）: `PrevNext` の後に `<Comments />` を配置（記事ページのみ。固定ページには出さない）
- 初期テーマ: 描画時の `data-theme`（システム追従時は解決後の light/dark）を giscus の `data-theme` に渡す

### テーマ連動
- 既存テーマ切替ロジック（`ThemeToggle` / `BaseLayout` の初期化スクリプト）が `data-theme` を変更する箇所に、giscus iframe への `postMessage` 送信を追加
- `iframe.contentWindow.postMessage({ giscus: { setConfig: { theme: 'light' | 'dark' } } }, 'https://giscus.app')`
- giscus 未ロード/iframe 不在時は no-op（記事ページ以外でエラーにならないようガード）

---

## 6. 機能3: 問い合わせ（Pages Functions + Turnstile + Resend）

### フォーム（クライアント）
- `pages/contact/page.md`（改修）: 本文をWordPressフォーム残骸から簡潔な案内文に書き換え（送信先・返信目安など）
- `src/components/ContactForm.astro`（新規）: 氏名 / メール / 題名 / 本文の入力＋ **Turnstile ウィジェット**（`PUBLIC_TURNSTILE_SITE_KEY`）＋送信ボタン。送信は `fetch('/api/contact', { method:'POST' })`、結果（成功/失敗）をUIに反映。ファイル添付はスコープ外（YAGNI）
- `src/pages/[...slug].astro`（改修）: 固定ページ描画時に `entry.data.slug === 'contact'` のとき本文後に `<ContactForm />` を差し込む（**専用ルートを新設せず、`[...slug]` との重複ルートを回避**）
- スタイルは既存CSS変数でテーマ連動

### サーバー（Pages Functions）
- `functions/api/contact.ts`（新規）: `onRequestPost`
  1. リクエストボディ（JSON）をパースし**入力検証**（必須項目・メール形式・最大長）
  2. **Turnstile 検証**: `https://challenges.cloudflare.com/turnstile/v0/siteverify` に `TURNSTILE_SECRET_KEY` とトークンを POST、`success` を確認
  3. **Resend 送信**: `https://api.resend.com/emails` に `RESEND_API_KEY` で POST（`from`=`CONTACT_FROM_EMAIL`, `to`=`CONTACT_TO_EMAIL`, `reply_to`=送信者メール, 件名/本文を整形）
  4. 成功/失敗を JSON（`{ ok: boolean, message }`）で返す。失敗時は適切なステータスコード
- 検証・整形の**純粋ロジックは `src/lib/contact.ts` に切り出し**、vitest でテスト可能にする（`onRequestPost` 本体は I/O のみ）。`functions/api/contact.ts` は相対 import（`../../src/lib/contact.ts`）で参照し、wrangler のバンドルで解決する（`astro:` 系 import を含まない純粋TSとする）

### ローカル検証
- `npm run build && npx wrangler pages dev dist`
- シークレットは `.dev.vars`（gitignore）で注入。Resend 実キーがあれば実送信、無ければ検証・整形・エラーハンドリングまで確認
- 本番疎通（独自ドメイン認証後の到達確認）は計画4

---

## 7. 環境変数 / シークレット管理

| 変数 | 種別 | 参照側 | 用途 |
|---|---|---|---|
| `PUBLIC_TURNSTILE_SITE_KEY` | 公開 | クライアント | Turnstile ウィジェット表示 |
| `PUBLIC_GISCUS_REPO` | 公開 | クライアント | giscus 対象リポジトリ |
| `PUBLIC_GISCUS_REPO_ID` | 公開 | クライアント | giscus リポジトリID |
| `PUBLIC_GISCUS_CATEGORY` | 公開 | クライアント | giscus カテゴリ名 |
| `PUBLIC_GISCUS_CATEGORY_ID` | 公開 | クライアント | giscus カテゴリID |
| `TURNSTILE_SECRET_KEY` | 秘密 | Functions | Turnstile サーバー検証 |
| `RESEND_API_KEY` | 秘密 | Functions | Resend メール送信 |
| `CONTACT_FROM_EMAIL` | 設定 | Functions | 送信元（例 noreply@shiimanblog.com） |
| `CONTACT_TO_EMAIL` | 設定 | Functions | 受信先（問い合わせ通知先） |

- 管理方針:
  - ローカル: `PUBLIC_*` は `.env`、秘密値は `.dev.vars`（いずれも gitignore）。**両ファイルともユーザーが手動作成**する
  - 本番: Cloudflare Pages のプロジェクト環境変数（秘密値は暗号化扱い）。投入は計画4
  - **必要なキー名の一覧は `README.md` に明記**（`.env*` はツールの権限制約で読み書きしないため、`.env.example`/`.dev.vars.example` は作らず README で代替）
- `.env` の内容は読まない・出力しない（CLAUDE.md 準拠）

---

## 8. リポジトリ構成（追加・変更分）

```
blog/
├─ functions/
│  └─ api/
│     └─ contact.ts            # 新規: Pages Functions（onRequestPost）
├─ src/
│  ├─ components/
│  │  ├─ Header.astro          # 改修: 検索アイコン追加
│  │  ├─ SearchOverlay.astro   # 新規: Pagefind オーバーレイ検索UI
│  │  ├─ Comments.astro        # 新規: giscus 埋め込み
│  │  └─ ContactForm.astro     # 新規: 問い合わせフォーム
│  ├─ layouts/
│  │  └─ PostLayout.astro      # 改修: Comments 組み込み
│  ├─ lib/
│  │  ├─ contact.ts            # 新規: 入力検証/整形の純粋ロジック
│  │  └─ contact.test.ts       # 新規: vitest
│  └─ pages/
│     └─ [...slug].astro       # 改修: contact に ContactForm 差し込み
├─ pages/contact/page.md       # 改修: 案内文へ
├─ .env.example                # 新規: PUBLIC_* キー名
├─ .dev.vars.example           # 新規: 秘密値キー名
├─ .gitignore                  # 追記: dist/pagefind/, .dev.vars
├─ package.json                # 改修: build スクリプト, pagefind 追加
└─ astro.config.mjs            # 変更なし
```

---

## 9. テスト

- vitest（既存43件を維持）に **`src/lib/contact.test.ts`** を追加
  - 入力検証（必須/メール形式/最大長/欠落）
  - Turnstile レスポンス処理（success / failure / ネットワーク異常）→ fetch をモック
  - Resend ペイロード整形（from/to/reply_to/件名/本文）
- `functions/api/contact.ts` 本体（I/O）は純粋ロジックに委譲し、ユニットテスト対象を `src/lib/contact.ts` に集約
- 検索・コメントは静的/外部依存のため自動テスト対象外（手動ローカル検証）

---

## 10. ローカル検証フロー（完了判定）

1. `npm install`（pagefind 追加後）
2. `npm test` → 全件パス
3. `npm run build` → `astro build` 成功 + `dist/pagefind/` 生成
4. `npm run preview` → 検索オーバーレイで記事ヒット、記事ページに giscus 表示・テーマ連動
5. `npx wrangler pages dev dist`（`.dev.vars` 注入）→ `/contact/` フォーム送信で `functions/api/contact.ts` が検証→Turnstile→Resend のフローを実行

---

## 11. リスクと対策

| # | リスク | 対策 |
|---|---|---|
| 1 | dev時に Pagefind index が無く検索が動かない | 仕様として割り切り。検証は build→preview。UIにフォールバック表示 |
| 2 | giscus 設定値未取得のまま実装が進む | 公開値を `PUBLIC_*` 環境変数化し、未設定時はコメント欄を非表示/プレースホルダ。値投入は事前準備（§12）で実施 |
| 3 | Resend 無料枠・ドメイン未認証 | 計画3はローカル検証主体。独自ドメイン認証と本番キーは計画4。無料枠（月間送信上限）は運用前提に収まる想定 |
| 4 | Turnstile 秘密鍵の漏洩 | Functions のみ参照。クライアントには site key のみ。秘密値は `.dev.vars`/Cloudflare 環境変数で管理しコミットしない |
| 5 | テーマ切替と giscus iframe の同期漏れ | iframe 不在時 no-op のガード。初期テーマも描画時に渡す。手動でライト/ダーク往復を確認 |
| 6 | `[...slug]` と contact 専用ルートの衝突 | 専用ルートを作らず、`[...slug].astro` 内で contact 時に ContactForm を差し込む方式で衝突を回避 |
| 7 | Pages Functions のローカル挙動差異 | `wrangler pages dev` で本番に近い形で検証。本番疎通は計画4で再確認 |

---

## 12. 事前準備（ユーザー手作業）

実装と並行して、以下を準備する（実値の投入はローカル `.env`/`.dev.vars`、本番は計画4）。

1. **giscus**: `shiiman/blog` で Discussions 有効化 → `Announcements` 形式カテゴリ作成 → giscus app インストール → giscus.app で 4 値取得
2. **Cloudflare Turnstile**: サイト登録 → site key / secret key 取得
3. **Resend**: アカウント作成 → API キー発行（独自ドメイン認証は計画4） → 送信元/受信先メール決定

> いずれも未取得でも実装は進行可能（プレースホルダで構築 → 値投入後にローカル検証）。

---

## 13. スコープ外（YAGNI）

- 既存 WordPress コメントの移行（giscus で新規のみ）
- 問い合わせのファイル添付
- 検索の高度なファセット/フィルタUI（Pagefind 標準機能で十分）
- 本番デプロイ・DNS・リダイレクト・解約（すべて計画4）
