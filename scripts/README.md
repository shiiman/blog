# 移行スクリプト（計画1: データ移行レイヤー）

WordPress(ConoHa) → Astro/Cloudflare 移行のためのワンショットスクリプト群です。
WordPress 前提の記事を Astro が読める形（ローカル画像＋Git LFS、slug 配列・`draft`・`eyecatch` のフロントマター）へ変換します。

> **重要: すべてプロジェクトルートから実行してください。** スクリプトは `posts/`・`pages/`・`data/` を
> カレントディレクトリからの相対パスで参照します（`npm run ...` はルートで実行されます）。

## 実行順

1. `npm run export:wp` — WordPress REST API からカテゴリ・タグ・アイキャッチURL・旧サイトマップを `data/` に出力（読み取り専用・認証不要）
2. `npm run localize:images` — 本文の wp-content 画像をローカルDL（LFS）し、本文リンクを相対パスへ書き換え。失敗URLは `data/missing-images.txt` に記録
3. `npm run migrate:frontmatter` — `data/*.json` を使い全フロントマターを Astro スキーマへ変換
4. `npx tsx scripts/verify-migration.ts` — 受け入れ検証（wp-content 残存・WP固有フィールド・draft型・slug有効性・eyecatch実在）

いずれのスクリプトも冪等です（localize はローカル化済み本文を再DLしません）。

## 構成

| ファイル | 役割 |
|---|---|
| `lib/taxonomy.ts` | カテゴリ/タグ ID → slug 変換（純粋関数） |
| `lib/frontmatter.ts` | WP → Astro フロントマター変換（純粋関数） |
| `lib/images.ts` | 本文画像URLの抽出・正規化・ローカル名生成・置換（純粋関数） |
| `lib/content-roots.ts` | 変換対象（`posts/`・`pages/`）の列挙。`future-sites/` は対象外 |
| `export-wp-data.ts` | REST エクスポート（I/O） |
| `localize-images.ts` | 画像DL＋本文書き換え（I/O） |
| `migrate-frontmatter.ts` | フロントマター一括変換（I/O） |
| `verify-migration.ts` | 受け入れ検証 |

`lib/` 配下の純粋関数は `*.test.ts` で vitest テスト済み（`npm test`）。

## 計画2（Astroサイト構築）で追加のスクリプト

| ファイル | 役割 |
|---|---|
| `lib/permalinks.ts` | 記事 canonical URL の link→path 変換（純粋関数） |
| `export-permalinks.ts` | 全記事 canonical URL を `data/permalinks.json` に出力（I/O・WP稼働中のみ） |
| `lib/cleanup-tracking.ts` | 本文の不可視トラッキング画像除去（純粋関数） |
| `cleanup-tracking-images.ts` | 上記を全本文へ適用（I/O・冪等） |

- `npm run export:permalinks` — URL保持の要。**ConoHa 稼働中に必ず実行・コミット**。
- `npm run cleanup:tracking` — A8/valuecommerce/accesstrade の計測画像を除去（冪等）。
