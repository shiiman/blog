---
name: eyecatch-create
description: ブログ記事のアイキャッチ画像を生成する。「アイキャッチ作成」「eyecatch」「サムネイル生成」「アイキャッチ画像」などで起動。記事の内容を分析し、統一デザインルールに基づいた高品質な画像を生成。
---

# Create Eyecatch Skill

ブログ記事の内容に基づいたアイキャッチ画像を生成します。

## ワークフロー

### 1. 対象記事の確認

ユーザーに対象となる記事のパスを確認:
- 例: `drafts/2026-01-27_my-article/article.md`
- 例: `posts/2026-01-27_my-article/article.md`

### 2. 記事の分析

記事ファイルを読み込み、以下を抽出:
- タイトル（Front MatterのtitleまたはH1）
- 主要キーワード（技術名、ツール名など）
- 記事の雰囲気（チュートリアル/解説/レビュー）

### 3. プロンプト生成

以下のデザインルールに厳密に従ってプロンプトを作成:

**【デザインルール】**
- **アスペクト比**: 16:9 (1280x720)
- **ベースカラー**: 純粋な白 (#FFFFFF) を背景色とする
- **アクセントカラー**: 鮮やかな青 (#007BFF) と濃いグレー (#424242)
- **テキスト**:
  - 記事のメインタイトルを中央に大きく配置
  - フォントは太字のサンセリフ体（日本語）
  - 文字色は濃いグレーまたは黒、わずかにドロップシャドウで視認性を確保
- **スタイル**: 「テクニカル・クリーン」
  - 清潔感があり、構造的なデザイン
  - shiimanblog.com のトンマナを再現

**【プロンプトテンプレート】**

```
A clean, professional blog eyecatch image with 16:9 aspect ratio.
Background: Pure white (#FFFFFF).
Main text: "[記事タイトル]" in bold sans-serif Japanese font, centered, dark gray color with subtle drop shadow.
Accent elements: Minimal geometric shapes or tech-related icons in corporate blue (#007BFF).
Style: Technical, minimalist, clean. No photographs, illustration-based.
```

### 4. 画像生成

`GenerateImage` ツールを使用:

```
GenerateImage(
  description: [生成したプロンプト],
  filename: "eyecatch.png"
)
```

### 5. 画像の配置

生成された画像を記事ディレクトリ内の `assets/` フォルダに移動:

```bash
mkdir -p [記事ディレクトリ]/assets
mv assets/eyecatch.png [記事ディレクトリ]/assets/eyecatch.png
```

### 6. 完了報告

ユーザーに以下を報告:
- 生成された画像のパス
- Front Matter の `eyecatch: ./assets/eyecatch.png` で記事に紐付くこと（画像は git push で記事ごとデプロイされるため、別途アップロード操作は不要）

## 注意事項

- このスキルの役割は `assets/eyecatch.png` を生成して記事ディレクトリに配置することのみです。画像は記事を `main` に git push すれば GitHub Actions のビルドでそのままデプロイされます
- 記事の Front Matter に `eyecatch: ./assets/eyecatch.png`（相対パス）を設定すればアイキャッチとして表示されます
- 画像が気に入らない場合は、再度このスキルを実行して再生成可能です
