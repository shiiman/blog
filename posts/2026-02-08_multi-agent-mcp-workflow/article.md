---
id: 2231
title: "マルチエージェントMCPを作ってみた"
excerpt: "Zenn記事とmulti-agent-shogunを土台に、multi-agent-mcpとmulti-flow/multi-issue-flowを実装して見えた差分を整理。特に複数AI CLI対応、復旧機構、Issue/PR連携の価値を実例で解説します。"
categories: [1]
tags: [10, 20]
slug: "multi-agent-mcp-workflow"
status: publish
featured_media: 2232
---

## はじめに

最近、AIエージェントを並列で動かして実装を進める「マルチエージェント開発」がかなり現実的になってきました。

私も以下の先行記事・実装を参考にしつつ、MCPサーバーと実運用フローを自作しました。

- [Claude Code × MCPで複数AIエージェントを作って並列開発する](https://zenn.dev/shio_shoppaize/articles/5fee11d03a11a1)
- [MCPサーバーを作って、Claude Codeでマルチエージェントの並列開発を実現する（v1.1.0）](https://zenn.dev/shio_shoppaize/articles/8870bbf7c14c22)
- [Codex CLIでMCPマルチエージェント並列開発を実現する](https://zenn.dev/shio_shoppaize/articles/dc85db324bb3f0)
- [yohey-w/multi-agent-shogun](https://github.com/yohey-w/multi-agent-shogun)

今回の記事では、次の2点を中心にまとめます。

1. 自作した `multi-agent-mcp` と `multi-flow / multi-issue-flow` の構成
2. 参考実装にはない追加機能（差分）の調査結果

## この記事の前提

### 参考実装側

`multi-agent-shogun` のREADMEでは、最小構成のMCPツールセット（Worker作成、タスク送信、状態確認、出力確認など）で並列開発を回す思想が示されています。まず「動くマルチエージェント」の基礎を作るには非常にわかりやすい構成です。

### 今回の実装側

今回作った実装は次の3つです。

- MCP本体: [`shiiman/multi-agent-mcp`](https://github.com/shiiman/multi-agent-mcp)
- Claude向けWorkflow: [`shiiman/claude-code-plugins`](https://github.com/shiiman/claude-code-plugins) の `multi-flow` / `multi-issue-flow`
- Codex向けWorkflow: [`shiiman/dotfiles`](https://github.com/shiiman/dotfiles) の同名グローバルスキル

狙いは「試作レベルを超えて、日常開発で壊れにくく回せる運用基盤」にすることでした。

## 参考にはない機能を調査して整理

比較しやすいように、差分を機能カテゴリごとにまとめます。

| 機能カテゴリ | 参考実装（基礎） | 今回の追加実装 | 実運用での価値 |
|---|---|---|---|
| ロール分離 | Worker中心 | `owner / admin / worker` の3層ロール + 権限チェック | 操作責務を分離し、誤操作を減らせる |
| タスク運用 | 単純なタスク送信 | タスク作成・キュー投入・依存関係・優先度 | 大きめタスクを分割して管理しやすい |
| 自動割り当て | 手動寄り | 空きWorkerへの自動割り当て、バッチ作成 | 並列化の立ち上げコストを削減 |
| 障害復旧 | 基本操作中心 | ヘルスチェック、`attempt_recovery`、`full_recovery` | セッション異常時の復旧が速い |
| Git連携 | 最小限 | `git worktree` 管理、`gtr` 連携、完了タスクのマージ補助 | ブランチ衝突を減らしつつ並列実装 |
| 可観測性 | 状態確認 | ダッシュボード、進捗報告、未読管理、tmux出力取得 | 現在地が把握しやすくレビューしやすい |
| ナレッジ管理 | なし/限定的 | メモリ保存・検索、グローバルメモリ、アーカイブ | ノウハウ再利用で2回目以降が速い |
| コスト/モデル管理 | なし/限定的 | コスト集計、警告閾値、`standard/performance` 切替 | 速度とコストの両立を調整しやすい |
| CI前後フロー | 実装中心 | `multi-flow`（軽量）/`multi-issue-flow`（Issue〜PR） | プロジェクト規模に応じて運用を選べる |

### 補足: 追加実装のポイント

#### 運用の堅牢化

MCP側に「監視・復旧・権限」を持たせたことで、単に並列で速いだけでなく、止まりにくい運用に寄せられました。

#### 開発フローの2モード化

`multi-flow` は「Issueなしでサッと並列実装」、`multi-issue-flow` は「Issue起点でPRまで完走」と使い分けできます。

#### チーム向けの管理機能

メモリ/アーカイブ/コスト管理は、個人開発よりチーム開発で効いてきます。とくに「誰が何をやったか」と「次回再利用できる知見」が残るのが大きいです。

## 特におすすめ: 複数種類のAIエージェント対応

今回の実装で個人的に一番推しているのが、利用するAI CLIを切り替えられる設計です。

- Claude系で強いタスク
- Codex系で強いタスク
- Gemini系で強いタスク

を、同じマルチエージェント基盤で扱えるようにしました。

### なぜ効くのか

#### 1. タスク適性に合わせてモデルを振り分けられる

設計議論・実装・調査で得意分野が異なるため、単一モデル固定より精度が上がるケースが多いです。

#### 2. ベンダーロックを弱められる

1社の仕様変更や障害があっても、運用自体を止めずに回しやすくなります。

#### 3. 将来の機能進化を取り込みやすい

新しいモデル・新しいクライアントが出ても、フロー全体の作り直しを避けやすいです。

## 実際の使い分けイメージ

### 軽量で回したい場合（Issueなし）

```bash
# 例: 小〜中規模改修を並列実装
multi-flow "この機能を3タスクに分割して並列実装"
```

### ガバナンス込みで回したい場合（Issue〜PR）

```bash
# 例: Issue起点で並列実装し、PRまで作成
multi-issue-flow "要件XをIssue作成からPR作成まで実行"
```

### MCP側の実行イメージ

```text
1. ownerがworkspace初期化
2. adminがタスク分解・worker割り当て
3. workerが並列実装と進捗報告
4. 異常時はhealthcheck -> recovery
5. 完了後にレビュー・統合・クリーンアップ
```

## 公式機能の進化と今後方針（2026-02-08時点）

ここは現時点の動向も重要なので、具体日付で整理します。

### Claude側

Claude Code公式ドキュメントでは、`/agents` コマンドでエージェント作成・管理を行い、`.claude/agents/` でプロジェクト固有エージェントを定義できる導線が提供されています。

- 公式: [Claude Code - Agent Teams](https://docs.anthropic.com/ja/docs/claude-code/sub-agents)

私の体感でも、公式機能のほうが安定する場面は増えています。今後は自作フローに公式のAgent Teamsをどう組み込むかが次のテーマです。

### Codex側

OpenAI公式の「Introducing the Codex app」は **2026年2月2日公開** で、複数エージェント実行・worktree・スキル活用を含む運用のしやすさが一段上がっています。

- 公式: [Introducing the Codex app](https://openai.com/ja-JP/index/introducing-the-codex-app/)

Codex側もUI/UX改善が続く前提で、MCPの価値は「バックエンドのオーケストレーション」に寄っていくと見ています。

## まとめ

参考実装を土台にしつつ、今回の実装では次の方向に広げました。

1. 並列実装の速度だけでなく、運用の堅牢性を強化
2. Issueあり/なしの2種類のWorkflowを用意
3. 複数AI CLI対応で将来の選択肢を確保

マルチエージェントは「一部の人の実験」から「日常開発の選択肢」に移りつつあります。
まずは小さい改修からでも、ぜひ並列開発を体験してみてください。

## 参考リンク

- [Claude Code × MCPで複数AIエージェントを作って並列開発する](https://zenn.dev/shio_shoppaize/articles/5fee11d03a11a1)
- [MCPサーバーを作って、Claude Codeでマルチエージェントの並列開発を実現する（v1.1.0）](https://zenn.dev/shio_shoppaize/articles/8870bbf7c14c22)
- [Codex CLIでMCPマルチエージェント並列開発を実現する](https://zenn.dev/shio_shoppaize/articles/dc85db324bb3f0)
- [yohey-w/multi-agent-shogun](https://github.com/yohey-w/multi-agent-shogun)
- [shiiman/multi-agent-mcp](https://github.com/shiiman/multi-agent-mcp)
- [shiiman/claude-code-plugins](https://github.com/shiiman/claude-code-plugins)
- [shiiman/dotfiles](https://github.com/shiiman/dotfiles)
- [Claude Code - Agent Teams](https://docs.anthropic.com/ja/docs/claude-code/sub-agents)
- [Introducing the Codex app](https://openai.com/ja-JP/index/introducing-the-codex-app/)
