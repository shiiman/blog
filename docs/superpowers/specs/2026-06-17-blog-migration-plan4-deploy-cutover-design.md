# 計画4 設計書: 本番デプロイ・切替・旧環境解約

- 作成日: 2026-06-17
- 親設計書: [2026-06-17-blog-migration-astro-cloudflare-design.md](./2026-06-17-blog-migration-astro-cloudflare-design.md)（§6 移行パイプライン / §9 SEO / §10 リスク / §12 オープン事項）
- 前提設計書: [計画2](./2026-06-17-blog-migration-plan2-astro-site-design.md)・[計画3](./2026-06-17-blog-migration-plan3-dynamic-features-design.md)
- 位置づけ: 移行プロジェクトの「計画4（最終フェーズ）」。計画1（データ移行）・計画2（Astro静的サイト）・計画3（検索/コメント/問い合わせ）はすべて `main` にマージ済み。本計画で **本番デプロイ → DNS切替 → 旧環境（ConoHa）解約** を完遂する
- 対象サイト: shiimanblog.com

---

## 1. ゴールと完了条件

### ゴール
`main` で生成される静的サイト＋`functions/` を Cloudflare Pages で本番公開し、`shiimanblog.com` を WordPress(ConoHa) から完全移行する。旧URLを301で救済し、SEO資産を維持したまま ConoHa を解約する。

### 完了条件（このフェーズの Done）
1. `public/_redirects` が生成され、旧パーマリンク（カテゴリ/タグの日本語slug・`/feed/` 系）を301で新URLへ救済する
2. 旧 `sitemap-1.xml` の全URLが、新サイトで `200` か `301` で解決することを検証済み（漏れゼロ）
3. Cloudflare Pages へ GitHub 連携で自動ビルド・デプロイされ、`shiimanblog.com` が新サイトを配信する
4. 本番で 全URL/画像/検索(Pagefind)/コメント(giscus)/問い合わせ(Turnstile+Resend) が動作する
5. Resend が独自ドメイン `shiimanblog.com` で認証され、送信元が `noreply@shiimanblog.com` 等に切替済み
6. Search Console に新サイトマップを再送信し、主要URLのインデックス移行が進行している
7. 約2週間のモニタ後、ConoHa WING サーバー契約を解約済み（不可逆・最後に実施）

### スコープ外（YAGNI）
- `?p=ID` 形式の plain パーマリンクのリダイレクト（旧サイトは投稿名パーマリンク運用のため被リンクに残る可能性が低い）
- 旧 `sitemap.xml` 自体のリダイレクト（Search Console へ新サイトマップを再送信するため不要）
- `/wp-admin/` `/wp-login.php` `/xmlrpc.php` 等 WordPress 管理系のリダイレクト（404のままで可）
- 旧 `/mahjong/` `/poker/` の救済（404のまま放置。将来別ドメイン稼働時に301を追加）
- 既存 WordPress コメントの移行（計画3で決定済み・giscusで新規のみ）

---

## 2. 確定事項（計画4ブレインストーミングでの決定）

| # | 論点 | 決定 |
|---|---|---|
| 1 | ドメイン/DNS運用方針 | **Cloudflare Registrar へ移管**。ドメイン・DNS・CDN・リダイレクトを Cloudflare に集約。ConoHa サーバー解約と独立してドメインを保持 |
| 2 | リダイレクト網羅範囲 | **推奨セット**: カテゴリ/タグの日本語slug→enSlug ＋ `/feed/` 系→`/rss.xml` ＋ 旧 `sitemap-1.xml` 全URLで網羅検証し漏れを個別補完。`?p=ID`・旧 `sitemap.xml` 自体は除外 |
| 3 | 旧 `/mahjong/`・`/poker/` | **404のまま放置**（将来別ドメイン稼働時に301を追加） |
| 4 | 切替の安全設計 | **「NS切替（DNS管理移行）」と「サイト切替（カスタムドメイン割当）」を分離**。先にNSのみCloudflareへ移し（Aレコードは旧WPのまま）、Resend認証後にカスタムドメインをPagesへ割当てて初めて新サイトへ切替 |
| 5 | 解約前モニタ期間 | **約2週間**。KPI（主要URL200維持・301正常・SCインデックス移行開始・問い合わせ到達）で最終判断。判断はその時点で必ず確認（不可逆） |
| 6 | Astro アダプタ | **不要**（計画3で確定）。静的出力＋`functions/`（Pages Functions）は自動検出され、アダプタなしで動作 |

### 役割分担
- **私（Claude）**: コード・スクリプト・テスト・検証（Phase 0/A/D/F の自動化部分）、各Phaseの手順書/チェックリスト作成と伴走、トラブル時の調査
- **shiiman さん**: Cloudflare/ConoHa/Resend/Search Console のダッシュボード操作、DNS変更、ドメイン移管、解約などの外部操作
- **進め方**: 外部操作が多いため、**各Phaseで完了を確認してから次へ**進む。完了報告を早とちりして勝手に次へ進まない。シークレットはハードコードせず環境変数経由（`.env`/`.dev.vars` の内容は読まない・出力しない）

---

## 3. コード成果物

### 3.1 `_redirects` 生成スクリプト
既存の移行スクリプト群（`scripts/`）と同じ tsx パターンで実装する。`astro:content` には依存させず、`posts/*/article.md` を `gray-matter` で直接読むことで単体実行を可能にする。

- **`scripts/lib/redirects.ts`（純粋関数・vitest対象）**
  - `buildRedirects(input)`: 「実在taxonomy（enSlug保持のみ）＋静的ルール」から `_redirects` テキストを生成する純粋関数
- **`scripts/build-redirects.ts`（I/O）**
  1. 公開記事（`draft: false`）のフロントマターから、**実際に使われている** category/tag の保存slugを集計
  2. `data/categories.json`・`data/tags.json` から「**enSlug を持ち、かつ実在する**」term のみ抽出（リダイレクト先が404にならないことを保証）
  3. `buildRedirects()` で `_redirects` テキストを生成
  4. **`public/_redirects`** に書き出す（git管理。Astroが `dist/_redirects` へコピー）
- `package.json` に `build:redirects` スクリプトを追加。`public/` 配下はソース管理対象のため、**生成→コミット**を基本とする（生成タイミングをビルドに組み込むかは実装計画で確定）

### 3.2 生成される `_redirects` の内容

```
# Feeds → RSS
/feed/            /rss.xml   301
/comments/feed/   /rss.xml   301
/*/feed/          /rss.xml   301

# Category（日本語エンコードslug → enSlug。enSlug保持かつ実在のもののみ自動生成）
/category/%e7%af%80%e7%b4%84/   /category/savings/   301
# … 以下 enSlug 保持カテゴリ分

# Tag（同上）
/tag/%e3%83%a1%e3%83%bc%e3%83%ab/   /tag/mail/   301
# … 以下 enSlug 保持タグ分
```

- **記事URL** は新サイトが `permalinks.json` の `path` をそのまま再現するため**新旧完全一致 → 301不要**（`post-1834` のエンコード混じりURLも一致）
- **英語slugのカテゴリ/タグ** は新旧一致のため出力しない
- **enSlug を持つが公開記事で未使用** の term は新サイトにアーカイブが存在しないため出力しない（404回避）
- **mahjong/poker** はルールを書かない（=404のまま）
- Cloudflare `_redirects` の構文（splat `*` の中間使用可否 / placeholder `:name` / status既定値 / ルール上限2,100）は**実装計画フェーズで公式仕様を精査し、実デプロイ（プレビューURL）で挙動を検証**する

### 3.3 テスト（vitest・既存テスト群に追加）
`scripts/lib/redirects.test.ts`:
- enSlug 有り → リダイレクト行を生成する
- enSlug 無し（新旧一致）→ 生成しない
- 公開記事で未使用の taxonomy → 除外する
- 静的ルール（`/feed/` 系）が含まれる
- 出力フォーマットが `<from> <to> <status>` で正しい

### 3.4 網羅検証スクリプト `scripts/verify-redirects.ts`
- 入力: 旧 `sitemap-1.xml`（実URL一覧）。**WordPress稼働中に取得**し `data/old-urls.xml`（仮）として保存（解約後は取得不能＝Phase 0で確保）
- 各旧URLが「新 `dist` に静的ファイルとして存在」または「`_redirects` ルールにマッチ」のどちらかを満たすか照合し、**漏れURL（404になるもの）を一覧出力**
- デプロイ後の実HTTP（200/301）検証は runbook の Phase F で実施

---

## 4. 本番切替 runbook（順序厳守＝不可逆対策）

**安全設計**: DNS管理の移行（NS切替）とサイトの切替（カスタムドメイン割当）を分離する。NS切替時点ではAレコードを旧WPのまま維持し、サイトは無変化。Resend認証完了後にカスタムドメインをPagesへ割り当てて初めて新サイトへ切り替わる。

| Phase | 内容 | 担当 |
|---|---|---|
| **0. 事前確保**（WP稼働中・不可逆対策） | 旧 `sitemap-1.xml`（実URL一覧）取得→`data/`保存。既存エクスポート（permalinks/categories/tags/featured-media/画像実体）の完全性を再確認 | 私+確認 |
| **A. コード準備** | `redirects.ts`/`build-redirects.ts`/`verify-redirects.ts` 実装＋テスト → `public/_redirects` 生成・コミット → ローカル網羅検証（漏れゼロ）→ `npm run build`・`npm test` 通過 | 私 |
| **B. Pages作成・GitHub連携** | Cloudflare Pages で `shiiman/blog` 連携。ビルド=`npm run build`、出力=`dist`。初回ビルドで**プレビューURL(`*.pages.dev`)** 発行。`functions/` は自動検出（アダプタ不要） | shiiman+手順 |
| **C. 本番シークレット投入** | Pages環境変数(Production)に `PUBLIC_*`(giscus4値+Turnstile site key) と秘密値(`TURNSTILE_SECRET_KEY`/`RESEND_API_KEY`/`CONTACT_FROM_EMAIL`/`CONTACT_TO_EMAIL`)。**ダッシュボードで直接入力**（私は値を読まない/出力しない）→ 再デプロイ | shiiman |
| **D. プレビュー検証** | `*.pages.dev` で全URL/画像/検索/giscus/問い合わせ/リダイレクトを検証（本番ドメインは無変化） | 私+shiiman |
| **E1. NS切替** | Cloudflareにゾーン追加→既存DNSレコード(WP向けA/MX/TXT等)をインポート→**NSをConoHa→Cloudflareへ変更**。※サイトは依然WordPress（DNS管理のみ移行） | shiiman |
| **E2. Resend独自ドメイン認証** | Resendに`shiimanblog.com`登録→指定のSPF/DKIM/DMARCをCloudflare DNSに追加→認証完了→`CONTACT_FROM_EMAIL`を`noreply@shiimanblog.com`へ切替・再デプロイ | shiiman |
| **E3. Turnstile/giscus 本番設定** | Turnstileの許可ホストに`shiimanblog.com`追加。giscus(`shiiman/blog`)の本番ドメイン動作確認 | shiiman |
| **E4. サイト切替** | Pagesにカスタムドメイン`shiimanblog.com`(+`www`)を割当→Cloudflareが自動でレコード設定→**この瞬間に新サイトへ切替**。HTTPS(TLS証明書)発行完了を確認 | shiiman |
| **F. 本番検証** | 全URL200 / 画像 / 検索 / giscus実投稿 / 問い合わせ実送信到達 / リダイレクト301 / `verify-redirects`の旧URL全件200or301 | 私+shiiman |
| **G. Search Console** | 新サイトマップ(`sitemap-index.xml`)再送信、主要URLのインデックス・カバレッジ監視 | shiiman |
| **H. ドメイン移管→モニタ→解約** | ドメインをCloudflare Registrarへ移管(ConoHa側で移管ロック解除・AuthCode取得)。約2週間のモニタ後に**ConoHa WINGサーバー契約のみ解約**（不可逆・最後） | shiiman |

---

## 5. 環境変数 / シークレット（本番投入）

計画3で定義済みの変数を Cloudflare Pages のプロジェクト環境変数（Production）へ投入する。**実値はダッシュボードで直接入力**し、コード・ドキュメントにはハードコードしない。

| 変数 | 種別 | 用途 |
|---|---|---|
| `PUBLIC_TURNSTILE_SITE_KEY` | 公開 | Turnstile ウィジェット |
| `PUBLIC_GISCUS_REPO` / `PUBLIC_GISCUS_REPO_ID` / `PUBLIC_GISCUS_CATEGORY` / `PUBLIC_GISCUS_CATEGORY_ID` | 公開 | giscus 設定 |
| `TURNSTILE_SECRET_KEY` | 秘密 | Turnstile サーバー検証 |
| `RESEND_API_KEY` | 秘密 | Resend メール送信 |
| `CONTACT_FROM_EMAIL` | 設定 | 送信元（認証後 `noreply@shiimanblog.com` 等） |
| `CONTACT_TO_EMAIL` | 設定 | 受信先（問い合わせ通知先） |

- 公開値(`PUBLIC_*`)はビルド時に埋め込まれるため、**Production環境にも設定したうえで再ビルド**が必要
- 秘密値は Functions 実行時に参照（暗号化扱い）
- README に「本番（Cloudflare）への投入は計画4」と既に明記済み。投入後の確認手順を README に追記する

---

## 6. リスクと対策

| # | リスク | 対策 |
|---|---|---|
| 1 | シークレット未投入で問い合わせ/コメントが動かない | Phase C で全キー投入を確認 → Phase D のプレビューで実動作確認 |
| 2 | NS切替時の既存DNSレコード取りこぼし（メール等が停止） | E1 で ConoHa の全DNSレコードを書き出してから Cloudflare にインポート。MX/SPF/DKIM/TXT を漏れなく移行 |
| 3 | カスタムドメイン割当時の一時的TLS未発行 | E4 で証明書発行完了を待ってから検証。切替は低トラフィック時間帯に実施 |
| 4 | リダイレクト漏れ（旧URLが404） | Phase 0 で旧 `sitemap-1.xml` を確保、`verify-redirects` でローカル照合（漏れゼロ）、Phase F で実HTTP再検証 |
| 5 | Resend独自ドメイン認証の不備（メール不達/迷惑メール判定） | E2 で SPF/DKIM/DMARC を正しく設定し、認証完了後に実送信テスト（到達・ヘッダ確認） |
| 6 | 解約後にエクスポート漏れが発覚（不可逆） | Phase 0（事前確保）＋ Phase H（解約前チェックリスト）の二重確認。約2週間のモニタで担保 |
| 7 | `_redirects` のCloudflare構文差異（splat等） | 実装計画で公式仕様を精査、プレビューURLで実挙動を検証してから本番反映 |
| 8 | ドメイン移管の遅延（移管ロック/AuthCode/60日制限等） | 移管はサイト切替（E4）とは独立に進行可能。NSが既にCloudflareなら移管完了前でも本番稼働に支障なし |

---

## 7. オープン事項（実装計画フェーズで確定する細部）
- Cloudflare `_redirects` の splat/placeholder の正確な構文と、`/*/feed/` のような中間splatの可否（実デプロイで検証）
- `build:redirects` の実行タイミング（手動生成＋コミット / `prebuild` フック）
- 旧 `sitemap-1.xml` の保存先ファイル名と取得スクリプト（`export-wp-data.ts` への追記 or 新規）
- Search Console のプロパティ形式（ドメインプロパティ / URLプレフィックス）と既存登録の引き継ぎ
- `www` サブドメインの正規化方針（`www`→apex の統一 or 併存）
