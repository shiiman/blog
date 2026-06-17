# ブログ移行設計書: WordPress → Astro + Cloudflare Pages

- 作成日: 2026-06-17
- 対象サイト: shiimanblog.com（現行: WordPress on ConoHa WING / Cocoon テーマ）
- 目的: WordPress を廃止し、Astro で生成した静的サイトを Cloudflare Pages で完全無料運用へ移行する

---

## 1. ゴールと前提

### 達成したいこと
- カスタムドメイン `shiimanblog.com` を維持し、既存 URL を保持（SEO 資産・被リンクの引き継ぎ）
- 完全無料運用（ホスティング・CDN・DNS・動的機能すべて無料枠内）
- 記事・固定ページ・画像をすべて自前管理し、WordPress 依存を撤去
- プロフィール・アプリ LP・規約ページなどを自由に作れる構成
- AI（Claude Code スキル）で記事作成・編集しやすい Markdown ベース構成
- サイト内検索・問い合わせフォーム・記事コメントを無料で提供

### 確定した方針（ブレインストーミングでの決定事項）
- WordPress（ConoHa）は**廃止**。本文中の `wp-content` 画像は全て**ローカル化**する
- ドメインは**継続**し、URL を保持して**301 リダイレクト**で旧 URL を救済する
- 動的 3 機能（検索・問い合わせ・コメント）はすべて**無料**で実装する
- デザインは**軽量・ミニマルを新規作成**（既存 Cocoon の見た目は踏襲しない）
- 既存 WordPress コメントは**移行しない**（giscus で新規コメントのみ受け付ける）
- 固定ページ `mahjong` / `poker` は今回の移行**対象外**。将来**別ドメインの別サイト**として移行するため、`pages/` から `future-sites/` へ退避させ、今回のビルド・リダイレクト・サイトマップのいずれにも含めない

### 現状の規模（移行対象）
- 記事: 86 公開 + 1 下書き（`posts/YYYY-MM-DD_slug/article.md`）
- 固定ページ: 今回の移行対象は 4 件（contact, privacy-policy, profile, sitemap）。mahjong / poker は対象外（後述）
- 本文が参照する `wp-content` 画像 URL: 698 個（サムネ除く原本は約 367 個）
- アイキャッチ（`featured_media`）を持つ記事: 80 件
- 使用カテゴリ ID: 24 種 / タグ ID: 約 100 種（数値 ID → 名前の変換が必須）
- 既存ローカル画像実体: 11 ファイルのみ（大半の画像は WordPress 上にある）

---

## 2. プラットフォーム選定

**採用: Cloudflare Pages + Astro**

| 観点 | Cloudflare Pages（採用） | GitHub Pages（退避先） | Netlify |
|---|---|---|---|
| 無料帯域 | 実質無制限 | 100GB/月（ソフト上限） | 100GB/月 |
| 動的機能 | Functions 同梱 | なし（外部必須） | Functions 有 |
| リダイレクト | `_redirects` ネイティブ | 非対応 | 対応 |
| DNS/CDN/ドメイン | 同一基盤で完結 | 別管理 | 別管理 |
| ビルド | 内蔵（500回/月） | GitHub Actions | 300分/月 |

**決め手**: (a) 帯域無制限で将来無料の確度が最も高い、(b) 問い合わせフォームを Pages Functions + Turnstile で同一基盤に閉じられる、(c) DNS・リダイレクト・ドメイン管理が 1 画面で完結。

**ロックインの小ささ**: Astro の出力は素の HTML/CSS/JS。Cloudflare の方針が変わっても同じビルド成果物を GitHub Pages へ退避できるため、プラットフォーム選定自体が低リスク（可逆）。GitHub Pages を恒久的な退避先として位置付ける。note は集客・告知用のサブ導線として併用。

---

## 3. 技術スタック（すべて無料枠内）

| 役割 | 採用技術 | 補足 |
|---|---|---|
| サイト生成 | Astro（最新）+ Content Collections | Markdown → 静的 HTML |
| ホスティング | Cloudflare Pages | 帯域無制限 |
| 画像管理 | Git LFS（既存運用を踏襲） | 記事と co-located の `assets/` |
| 検索 | Pagefind | ビルド後に静的インデックス生成、サーバー不要 |
| コメント | giscus | GitHub Discussions 連携 |
| 問い合わせ | Pages Functions + Turnstile + Resend（無料枠） | 同一基盤で完結 |
| RSS | @astrojs/rss | `/feed/` からリダイレクト |
| サイトマップ | @astrojs/sitemap | `sitemap.xml` 生成 |
| デザイン | 軽量・ミニマルの独自 CSS | フレームワーク非依存で軽量維持 |

---

## 4. リポジトリ構成

**既存リポジトリを発展させる**（新規リポジトリは作らない）。`posts/` `pages/` はルートに残したまま Astro の glob ローダーで参照し、移設の手間とパス破壊を回避する。`assets/` は記事と co-located のまま Astro の画像最適化対象にする。

```
blog/
├─ src/
│  ├─ content.config.ts        # glob loader + zod スキーマ(posts/ pages/ を参照)
│  ├─ layouts/                 # Base / Post / Page
│  ├─ components/              # Header, Footer, PostCard, Comments(giscus), Search
│  └─ pages/                   # ルーティング(URL 保持)
│     ├─ index.astro           # トップ(記事一覧 / ページネーション)
│     ├─ [...slug].astro       # 記事・固定ページ詳細
│     ├─ category/[category].astro
│     ├─ tag/[tag].astro
│     ├─ rss.xml.ts
│     └─ 404.astro
├─ functions/api/contact.ts    # Cloudflare Pages Functions(問い合わせ受信)
├─ public/
│  ├─ _redirects               # 301 リダイレクト台帳
│  └─ robots.txt
├─ scripts/                    # 移行用ワンショット
│  ├─ localize-images.ts       # 画像 DL + 本文リンク書き換え
│  ├─ migrate-frontmatter.ts   # ID→slug 変換・スキーマ変換
│  └─ build-redirects.ts       # 旧 URL→新 URL マップ生成
├─ data/
│  ├─ categories.json          # ID→{name,slug} マップ(wp-cli 出力)
│  └─ tags.json
├─ posts/                      # 既存資産(ルート維持・glob 参照)
├─ pages/                      # 移行対象の固定ページ(contact, privacy-policy, profile, sitemap)
├─ future-sites/               # 今回スコープ外。将来別ドメインで移行する退避置き場
│  ├─ mahjong/                 # pages/mahjong から退避(glob 非対象)
│  └─ poker/                   # pages/poker から退避(glob 非対象)
├─ tools/wp-cli/               # 移行エクスポートに使用 → 完了後 archive
├─ astro.config.mjs
└─ package.json
```

---

## 5. コンテンツスキーマ（フロントマター変換）

`src/content.config.ts` の zod スキーマでバリデーションする。

| 現在（WordPress 前提） | 変換後（Astro） | 変換方法 |
|---|---|---|
| `categories: [2]`（数値 ID） | `categories: ["wordpress"]`（slug） | `data/categories.json` で変換 |
| `tags: [6,7]`（数値 ID） | `tags: ["blog","seo"]`（slug） | `data/tags.json` で変換 |
| `featured_media: 13`（ID） | `eyecatch: ./assets/eyecatch.png` | メディアを DL しローカル化 |
| `status: publish/draft` | `draft: false/true` | Astro 標準の下書き制御 |
| `id` | 保持（optional） | リダイレクト・参照用に残す |
| `title`/`slug`/`date`/`modified`/`excerpt` | そのまま | 変更なし |

---

## 6. 移行パイプライン（順序厳守＝不可逆対策）

### Phase 0: 解約前エクスポート（取り返しのつかないものを先に確保）
※ 既存 `wp-cli` の `categories`/`tags` はテーブル出力のみ・`media` は ID→URL 取得手段が無いため、エクスポートは **WordPress REST API（`/wp-json/wp/v2/...`、公開コンテンツは認証不要）を Node スクリプトから取得**する方式に統一する。
1. REST `/wp-json/wp/v2/categories`・`/tags`（ページング）→ `data/categories.json` / `data/tags.json`（ID→{name,slug}）
2. REST `/wp-json/wp/v2/media/<id>` で `featured_media` ID → `source_url` を取得 → `data/featured-media.json`
3. 旧サイトマップ（`/wp-sitemap.xml` 等）を取得し `data/old-sitemap.xml` としてリダイレクト元台帳に保存
4. （必要なら）既存コメントを WXR 等でバックアップ ※本設計では移行しない方針

### Phase 1: コンテンツ変換（WordPress は稼働したまま）
4.5. `pages/mahjong` / `pages/poker` を `future-sites/` へ `git mv` で退避（今回のスコープから除外）
5. `localize-images`: 本文の `wp-content` URL を全 DL → `posts/<slug>/assets/`（LFS）に保存。`-WxH` サムネは原本に集約。本文リンクをローカル相対パスへ書き換え。ライトボックス記法 `[![](u)](u)` を `![](u)` に正規化。A8 計測 gif 等のノイズを除去
6. `featured_media` → `eyecatch.png` を DL しフロントマターへ
7. `migrate-frontmatter`: カテゴリ/タグの ID→slug 変換、`draft` 化、不要フィールド整理

### Phase 2: サイト構築（Astro）
8. Astro 雛形 + Content Collections + レイアウト/ルーティング（URL 保持）
9. Pagefind / giscus / RSS / sitemap / 問い合わせ Functions を組み込み
10. `public/_redirects` 生成（旧 URL→新 URL、`/feed/`→`/rss.xml`）

### Phase 3: デプロイ＆切替
11. Cloudflare Pages へデプロイ → プレビューで全 URL・全画像・フォーム・検索・コメントを検証
12. DNS を Cloudflare へ切替 → 本番反映
13. Search Console へサイトマップ再送信・主要 URL のインデックス確認
14. 数日〜数週モニタ後に **ConoHa 解約**（不可逆。ここで初めて実施）

---

## 7. 運用フロー（移行後の AI 執筆ループ）

- **旧**: `/blog-write` → drafts → `/eyecatch-create` → `/blog-publish`（wp-cli → WordPress）
- **新**: `/blog-write` → `posts/<date>_<slug>/index.md`（`draft: true`）→ `/eyecatch-create` → **git commit & push → Cloudflare Pages が自動ビルド & デプロイ**
- 公開は `draft: false` にして push。PR ごとに Cloudflare がプレビュー URL を自動発行
- 既存スキル `blog-publish` / `blog-update` を「git 操作 + ローカルビルド確認」に改修し、wp-cli 依存を撤去
- `tools/wp-cli/` は移行のエクスポート用途で残し、移行完了後に archive

---

## 8. 動的 3 機能の無料実装

### 検索: Pagefind
- ビルド後に静的全文検索インデックスを生成。数十 KB の JS をクライアントで読み込むだけでサーバー不要・完全無料

### コメント: giscus
- GitHub Discussions と連携。記事 slug でスレッドをマッピング。既存 WordPress コメントは移行しない

### 問い合わせ: Pages Functions + Turnstile + Resend
- `functions/api/contact.ts` で Turnstile トークンを検証 → Resend（無料枠）でメール送信
- Cloudflare Pages Functions の無料枠内で完結。スパム対策に Turnstile を必須化

---

## 9. SEO / リダイレクト

- `public/_redirects` に旧パーマリンク → 新 URL の 301 を記述
- `/feed/` → `/rss.xml` のリダイレクト
- canonical / OGP / `sitemap.xml`（@astrojs/sitemap）/ `robots.txt` を整備
- Search Console へ再申請し、主要 URL のインデックス状況を監視

---

## 10. リスクと対策

| # | リスク | 対策 |
|---|---|---|
| 1 | 「無料」≠ ドメイン無料 | ドメイン更新費（年 1,000〜1,500 円）は不可避。解約時にドメインを Cloudflare Registrar 等へ移管 or ConoHa でドメインのみ継続を決定 |
| 2 | 画像ローカル化が大量（約 367 原本） | `localize-images` で自動化。Cloudflare Pages 上限（1 ファイル 25MiB / 合計 20,000 ファイル）を事前確認 |
| 3 | URL 保持の漏れ | Phase 0 で旧 URL 台帳を作成し `_redirects` に網羅。リンク切れを切替前に検証 |
| 4 | RSS 購読者の喪失 | `/rss.xml` 生成 + `/feed/` リダイレクト |
| 5 | カテゴリ/タグ ID 依存 | 解約前に `wp-cli` でタクソノミーをエクスポート（解約後は取得不能） |
| 6 | 既存コメントの喪失 | 移行しない方針を明示。必要なら Phase 0 でバックアップのみ確保 |
| 7 | WP 特有記法の混入 | ライトボックス・A8 計測 gif・ショートコードを変換/除去 |
| 8 | 問い合わせフォームのスパム | Turnstile を必須化 |
| 9 | 解約順序の誤り（不可逆） | エクスポート → 本番稼働確認 → 解約 の順を厳守 |
| 10 | ビルド時画像最適化の負荷 | 事前最適化 or ビルド時最適化のどちらに寄せるか構築時に決定 |

---

## 11. スコープ外（YAGNI）

- 固定ページ `mahjong` / `poker`（`future-sites/` へ退避し、別ドメインの別サイトとして将来移行）
- 既存 WordPress コメントの移行（割り切り）
- 既存 Cocoon デザインの忠実な再現（軽量ミニマルを新規作成）
- 会員機能・有料コンテンツ・DB を伴う動的機能（無料運用の前提を崩すため不採用）
- 動画配信・大量ファイル配信（無料枠を逸脱するため不採用）

---

## 12. オープンな決定事項（実装計画フェーズで確定）

- ドメインの移管先（Cloudflare Registrar / ConoHa 継続 / 他社）
- ビルド時画像最適化の方針（事前 vs ビルド時）
- 問い合わせメール送信先と Resend アカウントの準備
- giscus 用 GitHub リポジトリ / Discussions の準備
- カテゴリ/タグの最終的な slug 命名（`data/*.json` 確定時）
- 旧 URL `shiimanblog.com/mahjong/` `shiimanblog.com/poker/` の切替後の扱い（新サイトには存在しないため、404 のまま放置 / 将来の別ドメインが立ち次第リダイレクトを追加、のいずれか）
