---
id: 2244
title: Google生成メディアAPIをMCPで統合。画像・動画・音声をAIから直接生成する
slug: google-genmedia-mcp-setup
status: publish
excerpt: GoogleのImagen・Veo・Chirp・LyriaをMCPサーバー経由でClaude CodeやCodexから直接利用できるPython製ツールを自作しました。導入方法とデモを紹介します。
categories:
    - 133
featured_media: 2243
date: 2026-03-04T19:27:41
modified: 2026-03-04T19:27:41
---

## はじめに

Claude CodeやCodexで作業中に、「この説明に合う画像を生成して」「この画像を動画にして」「ナレーションとBGMも付けて」と指示するだけで、メディアファイルが出力される。そんな体験を実現するMCPサーバーを作りました。

[google-genmedia-mcp](https://github.com/shiiman/google-genmedia-mcp) は、GoogleのImagen・Veo・Chirp・Lyriaといった生成メディアAPIをMCP（Model Context Protocol）経由で提供するPython製サーバーです。

この記事では、なぜ作ったのか、何ができるのか、実際の生成物を交えて紹介します。

## なぜ自作したのか

### 公式Go実装の存在

GoogleCloudPlatformが公式で [`mcp-genmedia`](https://github.com/GoogleCloudPlatform/vertex-ai-creative-studio/tree/main/experiments/mcp-genmedia) というGo実装のMCPサーバーを公開しています。Imagen・Veo・Chirp・Lyriaに対応した本格的な実装です。

ただ、実際に使ってみると、自分のニーズとは合わない部分がいくつかありました。

### Pythonで再実装した理由

自作に踏み切った動機は4つあります。

1. **最新モデルへの対応** — 公式Go実装にまだ含まれていないモデル（Imagen 4、Veo 3.1など）を使いたかった。新モデルが出たらすぐに試したい
2. **設定のファイル化** — モデルのエイリアスやデフォルトパラメータを `config.yaml` で管理し、ビルドし直さずに自由に設定変更できるようにしたかった
3. **uvx一発インストール** — リポジトリのクローンもビルドも不要。`uvx` コマンド1行で導入可能にしたかった
4. **認証方式の拡充** — Vertex AI ADCだけでなく、API KeyやOAuth 2.0にも対応させたかった。個人開発者がGoogle Cloudプロジェクトなしで手軽に試せる環境を作りたかった

## できること — 8つのMCPツール

### ツール一覧

google-genmedia-mcpは8つのMCPツールを提供します。

| ツール | 概要 | 対応モデル |
| -------- | ------ | ----------- |
| `generate_image` | テキストから画像を生成 | Imagen 4 Ultra/Fast、Gemini 3.1 Flash 等 |
| `edit_image` | 画像編集（インペインティング、アウトペインティング、背景置換） | Imagen 3/4 |
| `generate_video` | テキストから動画を生成 | Veo 3.1 / 3.0 / 2.0 |
| `generate_video_from_image` | 画像から動画を生成（Image-to-Video） | Veo 3.1 / 3.0 / 2.0 |
| `generate_speech` | テキストを音声に変換（TTS） | Chirp 3 HD |
| `generate_music` | 音楽を生成（30秒、インストゥルメンタル） | Lyria |
| `combine_audio_video` | 動画と音声をffmpegで合成 | — |
| `server_info` | サーバー情報・利用可能ツール一覧を表示 | — |

`generate_image` はプロンプトやモデル指定に応じてImagen APIとGemini APIを自動的に切り替えます。参照画像を渡すとGeminiモードになるなど、ユーザーはAPIの違いを意識せずに使えます。

### 3つの認証方式

用途に応じて3つの認証方式から選べます。

| 認証方式 | 難易度 | 対応ツール | おすすめ用途 |
| --------- | -------- | ----------- | ------------ |
| **API Key** | 簡単 | 画像・動画・画像編集 | 個人で手軽に試したい場合 |
| **Vertex AI ADC** | 中程度 | 全ツール | GCPプロジェクトがある場合 |
| **OAuth 2.0** | やや複雑 | 全ツール | GCPなしで全機能を使いたい場合 |

API Key方式は [Google AI Studio](https://aistudio.google.com/) でキーを取得するだけで始められます。ただし、Chirp TTSとLyria音楽生成はCloud認証が必要なため、API Keyでは利用できません。

## デモ — 実際に生成してみた

ここからは、google-genmedia-mcpで実際に生成したメディアファイルを紹介します。

### Geminiで画像を生成（generate_image）

`generate_image` ツールでGemini 3.1 Flashを使い、マスコットキャラクターの画像を生成しました。

<figure><img src="https://shiimanblog.com/wp-content/uploads/2026/03/gemini-generated-image.jpg" alt="Gemini 3.1 Flashで生成した画像" width="400"></figure>

プロンプトにキャラクターの特徴やスタイルを指示するだけで、上記のような画像が生成されます。アスペクト比やモデルの切り替えもパラメータで指定可能です。

```yaml
# 使用例
prompt: "Pythonロゴのニット帽をかぶったかわいいロボットキャラクター..."
model: "Nano Banana 2"  # Gemini 3.1 Flash のエイリアス
aspect_ratio: "9:16"
```

### 画像から動画を生成（generate_video_from_image）

先ほど生成した画像を入力にして、`generate_video_from_image` ツールでVeoによるImage-to-Video変換を行いました。

<figure><video controls src="https://shiimanblog.com/wp-content/uploads/2026/03/veo-image-to-video.mp4" width="600"></video></figure>

静止画がアニメーションになり、キャラクターに動きが加わります。Veo 3.1ではプロンプトで動きの方向やスタイルを細かく指定できます。

```yaml
# 使用例
prompt: "キャラクターがゆっくりと手を振るアニメーション"
image_gcs_uri: "gs://bucket/generated-image.png"
model: "Veo 3.1"
duration_seconds: 5
```

### 生成動画＋生成音声＋生成BGMを合成（combine_audio_video）

最後に、複数のツールを組み合わせた一連のワークフローです。

1. `generate_video` でテキストから動画を生成
2. `generate_speech` でChirp 3 HDによるナレーション音声を生成
3. `generate_music` でLyriaによるBGMを生成
4. `combine_audio_video` で動画・音声・BGMを1つのファイルに合成

<figure><video controls src="https://shiimanblog.com/wp-content/uploads/2026/03/combined-audio-video.mp4" width="600"></video></figure>

テキスト指示だけで動画・ナレーション・BGMの生成から合成まで、AIとの対話で完結します。

## セットアップ

### 前提条件

- Python 3.11以上（`uvx` が動く環境）
- ffmpeg（`combine_audio_video` ツール使用時）
- Google AI StudioのAPIキー、またはVertex AI / OAuth環境

### Claude Codeへの登録

1行のコマンドで登録できます。

```bash
claude mcp add --scope user google-genmedia-mcp -- \
  uvx --from git+https://github.com/shiiman/google-genmedia-mcp google-genmedia-mcp
```

Codexの場合:

```bash
codex mcp add google-genmedia-mcp -- \
  uvx --from git+https://github.com/shiiman/google-genmedia-mcp google-genmedia-mcp
```

### config.yamlの設定

初回は設定ファイルを作成します。

```bash
mkdir -p ~/.google-genmedia-mcp
cp config.example.yaml ~/.google-genmedia-mcp/config.yaml
```

最小構成（API Key方式）:

```yaml
auth:
  method: "api_key"
  apiKey: "YOUR_API_KEY"

output:
  directory: "~/.google-genmedia-mcp/output"
```

モデルのエイリアスやデフォルト値も `config.yaml` で自由に変更できます。

```yaml
tools:
  generateImage:
    defaultModel: "Nano Banana 2"
    aspectRatio: "16:9"
    models:
      - id: "imagen-4.0-generate-001"
        aliases: ["Imagen 4"]
      - id: "gemini-3.1-flash-image-preview"
        aliases: ["Nano Banana 2"]

  generateVideo:
    defaultModel: "Veo 3.1"
    durationSeconds: 5
    models:
      - id: "veo-3.1-generate-preview"
        aliases: ["Veo 3.1"]
```

### 動作確認

登録後、Claude Codeで `server_info` ツールを呼ぶと、利用可能なツールとモデルの一覧が返ります。

## まとめ

- Google生成メディアAPI（Imagen / Veo / Chirp / Lyria）をMCPサーバー経由でAIクライアントから直接利用できるようにした
- uvx一発インストール + YAML設定管理で、手軽に導入・カスタマイズできる
- 画像生成から動画化、音声合成、合成まで一連のメディア制作がAIとの対話で完結する

興味のある方は、ぜひ試してみてください。

## 参考リンク

- [shiiman/google-genmedia-mcp](https://github.com/shiiman/google-genmedia-mcp) — 本記事で紹介したMCPサーバー
- [GoogleCloudPlatform/vertex-ai-creative-studio](https://github.com/GoogleCloudPlatform/vertex-ai-creative-studio/tree/main/experiments/mcp-genmedia) — Google公式のGo実装
- [Google AI Studio](https://aistudio.google.com/) — APIキー取得先
- [Model Context Protocol](https://modelcontextprotocol.io/) — MCP仕様ドキュメント