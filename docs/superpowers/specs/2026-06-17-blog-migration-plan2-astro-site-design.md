# 計画2 設計書: Astroサイト構築

- 作成日: 2026-06-17
- 親設計書: [2026-06-17-blog-migration-astro-cloudflare-design.md](./2026-06-17-blog-migration-astro-cloudflare-design.md)
- 位置づけ: 移行プロジェクトの「計画2」。計画1（データ移行レイヤー, マージ済）を前提に、`posts/`・`pages/` を読む Astro 静的サイトを構築する
- 対象サイト: shiimanblog.com（現行: WordPress on ConoHa / Cocoon）

---

## 1. ゴールと完了条件

### ゴール
計画1で Astro 可読状態になった記事（85件）・固定ページ（4件）を、軽量・ミニマルな独自デザインの静的サイトとして生成する。現行 WordPress の URL を保持し、RSS・サイトマップ・SEO メタを備える。

### 完了条件（このフェーズの Done）
1. ローカルで `npm run build` が成功する
2. 全 85 記事 + 4 固定ページが生成・表示される
3. 生成 URL が現行 WordPress と一致する（`data/permalinks.json` と照合）
4. カテゴリ／タグ アーカイブ、トップのページネーション、RSS、サイトマップ、404 が機能する
5. `npm test`（vitest）が通る（日付・データ整備ロジックの単体テスト含む）

### スコープ外（後続フェーズ）
- 検索（Pagefind）・コメント（giscus）・問い合わせ（Pages Functions + Turnstile + Resend） → **計画3**
- `public/_redirects`（301）・Cloudflare Pages デプロイ・DNS 切替・ConoHa 解約 → **計画4**

---

## 2. 確定事項（計画2ブレインストーミングでの決定）

| # | 論点 | 決定 |
|---|---|---|
| 1 | URL保持（primaryカテゴリの確定） | 稼働中 WordPress の REST API から全記事の canonical URL を取得し `data/permalinks.json` に保存。ルーティングはこれを参照（解約後は取得不能な資産） |
| 2 | 日付のタイムゾーン整形 | 表示は **JST・日付のみ・和文**（例「2021年8月30日」）。機械可読時刻（`<time datetime>` / RSS / sitemap）は **正しい UTC 瞬時に補正**して出力 |
| 3 | カテゴリ/タグの slug と表示名 | 表示名 = 日本語名（`name`）、URL slug = **英語に再命名**（例 節約→`savings`、税金→`tax`）。`data/*.json` に `enSlug` を追記 |
| 4 | 画像最適化 | **Astro ビルド時最適化**（`astro:assets` / sharp）。co-located `assets/` を入力に webp 化・サイズ最適化 |
| 5 | デザイン方向 | **モダン・ミニマル**（ライト基調・サンセリフ・青アクセント #2563eb 系・ヘアライン区切り・広い余白） |
| 6 | トップ記事一覧 | **リスト + 左サムネ**（各行に小さめアイキャッチ） |
| 7 | 記事詳細の目次 | **右サイドバー追従目次**（現在地ハイライト、スマホは折りたたみ） |
| 8 | カラーテーマ | **ライト / ダーク / システム追従の切替**を提供（Header にトグル、localStorage 永続化、FOUC 回避） |

### 調査で確認済みの事実（権威: 現行 WP REST API）
- パーマリンク: 記事 `/{カテゴリslug}/{記事slug}/`（`/%category%/%postname%/`）、固定ページ `/{slug}/`、カテゴリ `/category/{slug}/`、タグ `/tag/{slug}/`
- 85 記事中 72 記事が複数カテゴリを持つ。URL に入る primary カテゴリは「フロントマター配列の先頭」とは一致せず、おおむね「カテゴリ ID 最小」（WordPress 既定）。例外（Yoast 上書き）を取りこぼさないため canonical URL を直接取得する
- トップのページネーション `/page/2/` は現行で存在（200）
- 日本語カテゴリ/タグの slug は URL エンコード（例 節約=`%e7%af%80%e7%b4%84`）。記事 URL に出る primary カテゴリは全て英語 slug のため記事 URL には影響しない（アーカイブ URL のみ顕在化）

---

## 3. 技術スタック / Astro 設定

| 役割 | 採用 | 備考 |
|---|---|---|
| サイト生成 | Astro（最新）+ Content Collections | `output: 'static'` |
| 画像最適化 | `astro:assets`（sharp） | ビルド時 webp 化 |
| RSS | `@astrojs/rss` | `/rss.xml` |
| サイトマップ | `@astrojs/sitemap` | `sitemap.xml`、future-sites 除外 |
| コード装飾 | Astro 内蔵 Shiki | 追加依存なし |
| フォント | system-ui スタック | Web フォント配信なし（高速・軽量） |

`astro.config.mjs` の主要設定:
- `site: 'https://shiimanblog.com'`
- `trailingSlash: 'always'`、`build.format: 'directory'`（末尾スラッシュ URL を完全一致で再現）
- integrations: `sitemap()`（RSS は `rss.xml.ts` で実装）
- 静的アダプタ不要（Functions は計画3で `@astrojs/cloudflare` を導入）

---

## 4. コンテンツモデル（`src/content.config.ts`）

- **glob ローダー 2 本**
  - `posts`: `posts/*/article.md`
  - `pages`: `pages/*/page.md`
  - （`posts/`・`pages/` はルート据え置き。移設しない）
- **zod スキーマ**（フロントマター検証）
  - 共通: `title: string`, `slug: string`, `date: string`, `modified: string?`, `excerpt: string?`, `draft: boolean`, `id: number?`
  - posts 追加: `categories: string[]`(slug), `tags: string[]`(slug), `eyecatch: image()?`
  - `date`/`modified` は文字列で受け、日付ユーティリティ（§7）が「JST 壁時計」として解釈する
- `draft: true` はビルド対象から除外（`getCollection` のフィルタ）

---

## 5. ルーティング / URL 保持

| 種別 | URL | 実装ファイル | 補足 |
|---|---|---|---|
| トップ | `/`, `/page/2/`… | `src/pages/index.astro`, `src/pages/page/[page].astro` | 10 件/ページ（WP 既定） |
| 記事 | `/{primaryカテゴリ}/{slug}/` | `src/pages/[...slug].astro` | `data/permalinks.json` から経路生成 |
| 固定ページ | `/{slug}/` | `src/pages/[...slug].astro` | contact / privacy-policy / profile / sitemap |
| カテゴリ | `/category/{enSlug}/` | `src/pages/category/[category].astro` | アーカイブ |
| タグ | `/tag/{enSlug}/` | `src/pages/tag/[tag].astro` | アーカイブ |
| RSS | `/rss.xml` | `src/pages/rss.xml.ts` | `/feed/`→`/rss.xml` の 301 は計画4 |
| 404 | — | `src/pages/404.astro` | |

- 記事・固定ページは単一の `[...slug].astro` が `getStaticPaths` で全パスを生成
  - 記事の `params.slug` = `data/permalinks.json` から得た `{primaryカテゴリ}/{記事slug}`（先頭/末尾スラッシュは Astro 規約に合わせて正規化）
  - 固定ページの `params.slug` = `{slug}`
  - 衝突防止: 記事 URL は必ずカテゴリ階層配下、固定ページはトップ直下のため衝突しない（`category`/`tag`/`page` 予約語とページ slug が衝突しないことをビルド時に検証）
- カテゴリ/タグ アーカイブの slug は `enSlug`（§6）。旧エンコード URL の 301 は計画4

---

## 6. データ整備（移行レイヤーの拡張）

### 6.1 canonical URL の取得（`scripts/export-permalinks.ts` 新規）
- REST `/wp-json/wp/v2/posts`（ページング, `_fields=id,slug,link`）で全記事の canonical URL を取得
- `data/permalinks.json` を生成: `{ "<id>": { "slug": "...", "path": "/wordpress/conoha-wing/" } }`
- ルーティングはこのマップを参照（Yoast の主カテゴリ上書きも含め完全一致）

### 6.2 カテゴリ/タグの英語 slug（`data/categories.json` / `data/tags.json` 拡張）
- 各エントリに `enSlug` を追記（英語の無い既存英語 slug はそのまま流用）
- 日本語名（`name`）は表示名として使用
- 原案は実装フェーズで一覧化し、確定後にコミット（命名規則: 小文字・ハイフン区切り・既存英語 slug と重複させない）

### 6.3 日付ユーティリティ（`src/lib/date.ts` 新規）
- 入力（例 `2021-08-30T19:30:00.000Z`）を **JST 壁時計**として解釈
- `formatDisplay()` → 「2021年8月30日」（JST 日付のみ）
- `toUtcInstant()` → 正しい UTC 瞬時（壁時計 − 9h）を `Date`/ISO で返す（`<time datetime>` / RSS pubDate / sitemap lastmod 用）
- vitest で両関数を検証

---

## 7. レイアウト / コンポーネント（`src/`）

- **layouts**
  - `BaseLayout.astro`: html/head、メタ（canonical/OGP/Twitter Card）、共通 CSS（テーマ変数）、テーマ初期化インラインスクリプト、Header/Footer
  - `PostLayout.astro`: アイキャッチ → タイトル → メタ（カテゴリ・JST 日付）→ 本文（+ 右 TOC）→ タグ → 前後記事
  - `PageLayout.astro`: 固定ページ用（タイトル + 本文）
- **components**
  - `Header.astro`（サイト名 + ナビ + テーマトグル）
  - `Footer.astro`
  - `PostCard.astro`（**リスト + 左サムネ**）
  - `Pagination.astro`
  - `TableOfContents.astro`（**右サイドバー追従 + 現在地ハイライト**、スマホ折りたたみ。`IntersectionObserver` で現在地検出）
  - `PostMeta.astro`（カテゴリバッジ + 日付）、`TagList.astro`、`PrevNext.astro`、`ThemeToggle.astro`
- **スタイル**: `src/styles/global.css`（CSS 変数でライト/ダークのトークン定義）

### カラーテーマ（ライト/ダーク）
- `:root`（ライト）と `[data-theme="dark"]` で CSS 変数を切替
- 初期テーマ: `localStorage` の保存値 → なければ `prefers-color-scheme`
- FOUC 回避: `<head>` 先頭のインラインスクリプトで `data-theme` を即時付与
- トグル: ライト / ダーク / システム追従（3 状態）を `localStorage` に永続化

---

## 8. RSS / サイトマップ / SEO

- `src/pages/rss.xml.ts`: 全記事（`draft` 除く）を新しい順、`pubDate` は §6.3 の正しい UTC 瞬時
- `@astrojs/sitemap`: `sitemap.xml` 生成。future-sites は glob 非対象のため自然に除外（必要なら filter で明示除外）
- 各ページ: canonical、OGP（og:title / og:description / og:image=アイキャッチ）、Twitter Card、`public/robots.txt`

---

## 9. 画像最適化

- `astro:assets` の `<Image>` / `<Picture>` でビルド時に webp 化・サイズ最適化
- 入力は co-located の `posts/<slug>/assets/`（アイキャッチ・本文画像とも）
- 本文 Markdown 内の画像は、Content Collections の image 最適化に乗せる（必要なら remark/コンポーネントで対応。詳細は実装計画で確定）

---

## 10. リポジトリ構成（追加分）

```
blog/
├─ src/
│  ├─ content.config.ts
│  ├─ layouts/        Base / Post / Page
│  ├─ components/     Header, Footer, PostCard, Pagination, TableOfContents,
│  │                  PostMeta, TagList, PrevNext, ThemeToggle
│  ├─ lib/            date.ts ほか
│  ├─ styles/         global.css
│  └─ pages/          index, page/[page], [...slug], category/[category],
│                     tag/[tag], rss.xml.ts, 404
├─ public/            robots.txt（_redirects は計画4）
├─ scripts/           export-permalinks.ts（新規）
├─ data/              permalinks.json（新規）, categories.json/tags.json（enSlug 追記）
├─ astro.config.mjs   （新規）
└─ package.json       （astro 等を追加）
```

---

## 11. リスクと対策

| # | リスク | 対策 |
|---|---|---|
| 1 | canonical URL 取得は WP 稼働中のみ可能 | 計画2 着手時に `export-permalinks.ts` を最優先で実行・コミット |
| 2 | ビルド時画像最適化の負荷（約 367 原本） | ローカルビルド時間を計測。問題あれば対象を絞る/事前最適化に切替を検討 |
| 3 | URL 衝突（固定ページ slug と予約語 category/tag/page） | `getStaticPaths` 生成時に重複検出してビルドを失敗させる |
| 4 | 本文中の残存ノイズ（A8 計測 gif 等の外部画像参照） | 表示崩れを確認し、残存があれば本文クリーニングを実装計画に追加 |
| 5 | 日付補正の取り違え | `date.ts` を vitest で固定値検証（表示=JST 日付、機械可読=UTC 瞬時） |

---

## 12. このフェーズで残すオープン事項（実装計画で確定）
- カテゴリ/タグの `enSlug` 最終命名（原案 → レビュー → 確定）
- 本文画像を `astro:assets` 最適化に載せる具体手段（remark プラグイン or カスタムコンポーネント）
- ナビゲーション項目の確定（記事一覧 / カテゴリ / プロフィール / お問い合わせ 等）
