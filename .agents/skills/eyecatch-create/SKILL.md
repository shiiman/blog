---
name: eyecatch-create
description: "ブログ記事の内容からアイキャッチ画像を生成する。『アイキャッチを作成』『画像を再生成』『/eyecatch-create』などで起動。google-genmedia-mcp の generate_image を使って assets/eyecatch.png を作成。"
---

# Create Eyecatch Skill

ブログ記事（`article.md`）の内容をもとに、`assets/eyecatch.png` を生成します。

## 前提条件

- 対象記事ファイル（`drafts/.../article.md` または `posts/.../article.md`）が存在すること
- `google-genmedia-mcp` が利用可能であること

## ワークフロー

### 1. 対象記事の確認

- ユーザーから記事パスを受け取る
- 未指定の場合は、作業文脈上の対象記事を確認してから進める

### 2. 記事内容の分析

- Front Matter の `title` / `excerpt` / `tags` を確認
- 本文から主要キーワードを 3〜6 個抽出
- 想定読者と記事のトーン（入門・実践・比較など）を把握

### 3. 生成プロンプトの作成

次のデザインルールを必ず反映する:

- アスペクト比: 16:9
- 背景: 白（`#FFFFFF`）
- アクセント: 青（`#007BFF`）+ 濃いグレー（`#424242`）
- タイトル文字: 中央配置、太字サンセリフ、日本語可読性重視
- テイスト: テクニカル・クリーン

不要要素を避けるため、必要に応じて `negative_prompt` を設定:
- 過度な装飾
- 読みにくい小さい文字
- 低コントラスト
- ノイズの多い背景

### 4. 画像生成

`mcp__google-genmedia-mcp__generate_image` を使用:

- `aspect_ratio`: `16:9`
- `number_of_images`: `1`
- `output_mime_type`: `image/png`
- モデルは指定がなければデフォルトを使用（必要時のみ明示）

### 5. ファイル保存

`generate_image` 実行後、リポジトリルート相対の `.google-genmedia-mcp/output/` を探索し、採用画像を記事ディレクトリの `assets/eyecatch.png` に保存する。

探索対象拡張子:

- `.png`
- `.jpg`
- `.jpeg`
- `.webp`

採用ルール:

- 更新時刻（mtime）が最も新しいファイルを 1 件採用する
- 同一mtimeで複数候補がある場合のみ拡張子優先順位で採用する
- 拡張子優先順位: `.png` > `.jpg` > `.jpeg` > `.webp`

```bash
mkdir -p <article_dir>/assets
cp <selected_image_from_.google-genmedia-mcp/output> <article_dir>/assets/eyecatch.png
```

### 6. 結果報告

ユーザーへ次を報告:
- 保存先パス
- 使用したプロンプト要約
- 必要なら再生成案（文言短縮、コントラスト強化など）

## 重要な注意事項

- 既存の `assets/eyecatch.png` は再生成時に上書きする
- 画像生成が失敗しても記事本文は変更しない
- 公開/更新時は `wp-cli` が `assets/eyecatch.png` を自動アップロードする
- 既存 `featured_media` を差し替える場合は `wp-cli update <file> --force-eyecatch` を使用する

## 失敗時の扱い

- `.google-genmedia-mcp/output/` が存在しない場合: 失敗として報告する
- 対象拡張子（`.png`, `.jpg`, `.jpeg`, `.webp`）のファイルが 0 件の場合: 失敗として報告する
- 画像コピー（`cp`）に失敗した場合: 失敗として報告する
- いずれの場合も `article.md` は変更しない
