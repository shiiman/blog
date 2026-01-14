---
title: "AIのグローバル設定をdotfilesで管理する方法：PC移行でも自分専用のAIアシスタントを持ち歩く"
id: 2213
slug: "ai-settings-dotfiles"
status: publish
categories: [18, 133]
tags: [101, 102] # dotfiles, AI (Assuming IDs)
excerpt: "ClaudeやGeminiなどのAIアシスタントのカスタム指示（System Instructions）をdotfilesで管理し、複数のPCで環境を統一する方法を紹介します。シンボリックリンクを活用して、開発環境セットアップと一緒にAI設定も自動化しましょう。"
---

PCを新調したり、職場と自宅で別のマシンを使ったりするとき、開発環境のセットアップは `dotfiles` で自動化しているエンジニアも多いと思います。`.zshrc` や `.vimrc` は当然のようにGit管理されていますよね。

しかし、**AIアシスタントの設定（Custom Instructions / System Prompts）** はどうでしょうか？

「あのPCのClaudeでは日本語で答えるように設定したのに、こっちでは英語で返ってくる…」
「プロジェクトごとのルール設定を毎回コピペするのが面倒」

そんな経験はありませんか？
今回は、これらのAI設定ファイルも `dotfiles` で一元管理し、`setup.sh` 一発で自分専用のAIアシスタント環境を構築する方法を紹介します。

## なぜAI設定をdotfilesで管理するのか

エンジニアにとって、AIはもはや「単なるチャットボット」ではなく「ペアプログラミングの相手」です。
優秀なペアプロ相手には、自分の好みのコーディングスタイルや、プロジェクトのローカルルール（言語設定、コミットメッセージのフォーマットなど）を把握しておいてほしいものです。

AI設定をローカルファイル（`GEMINI.md`, `CLAUDE.md` など）として管理し、それを `dotfiles` に含めることで以下のメリットがあります。

1.  **環境の統一**: どのPCを使っても、同じ振る舞いをするAIアシスタントが利用できる。
2.  **バージョニング**: プロンプトの改善履歴をGitで追える。「先週のプロンプトの方が良かったな」という時に戻せる。
3.  **セットアップの自動化**: コマンド一発で配置完了。

## ディレクトリ構成の実例

私の `dotfiles` では、以下のように `ai` ディレクトリを掘って各AIの設定をまとめています。

```text
~/dotfiles
├── .zshrc
├── .vimrc
├── ai_setup.sh       # AI設定用のセットアップスクリプト
└── ai/
    ├── claude/
    │   ├── CLAUDE.md      # グローバルな指示書
    │   └── settings.json
    ├── codex/
    │   ├── config.toml    # API設定など
    │   └── skills/        # スキル定義ディレクトリ
    ├── cursor/
    │   ├── User/
    │   │   ├── settings.json     # エディタ設定
    │   │   └── keybindings.json  # キーバインド
    │   ├── extensions.json       # 拡張機能リスト
    │   └── mcp.json              # MCP設定
    └── gemini/
        └── GEMINI.md      # Gemini(Antigravity)用の設定
```

実際のコード構成はGitHubでも公開していますので、参考にしてください。
[https://github.com/shiiman/dotfiles](https://github.com/shiiman/dotfiles)


## 各AIアシスタントの管理詳細

私が実際に管理している4つの主要なAIツールについて、具体的な管理内容とディレクトリ構成を紹介します。

### 1. Claude Code
ClaudeのCLIツール「Claude Code」は、グローバル設定とシステムプロンプトを管理できます。

- **管理ファイル**:
    - `CLAUDE.md`: すべてのセッションに適用されるグローバルな振る舞いを定義。出力言語の指定や、禁止事項（機密情報のハードコード禁止など）を記述。
    - `settings.json`: テーマや通知設定などのツール自体の設定。

- **シンボリックリンク先**:
    - `~/.claude/CLAUDE.md`
    - `~/.claude/settings.json`

### 2. Cursor
AI搭載エディタ「Cursor」は、VS Code互換の設定ファイルに加えて、独自のAI設定を持っています。

- **管理ファイル**:
    - `User/settings.json`: エディタの基本設定（フォント、フォーマット設定など）。
    - `User/keybindings.json`: キーボードショートカット。
    - `mcp.json`: Model Context Protocol (MCP) サーバーの設定。
    - `extensions.json`: インストールすべき拡張機能リスト。

- **シンボリックリンク先**:
    - `~/Library/Application Support/Cursor/User/settings.json` など
    - `~/.cursor/mcp.json`

特に `setup.sh` では、`extensions.json` に記載された拡張機能IDを読み取り、`cursor --install-extension` コマンドで自動インストールする仕組みも実装しています。これでエディタの環境構築も自動化されます。

### 3. Codex (OpenAI)
Codex系のCLIツールや独自のスクリプトで使用する設定です。

- **管理ファイル**:
    - `config.toml`: APIキーの参照先やモデルパラメータ（temperatureなど）の設定。
    - `skills/`: 再利用可能なスキル（プロンプトテンプレート）の定義フォルダ。

- **シンボリックリンク先**:
    - `~/.codex/config.toml`
    - `~/.codex/skills/`

### 4. Gemini (Antigravity)
GoogleのAdvanced Agentic Codingエージェント用設定です。

- **管理ファイル**:
    - `GEMINI.md`: エージェントへの基本指示書。プロジェクト共通のルールやコンテキストを記載。

- **シンボリックリンク先**:
    - `~/.gemini/GEMINI.md`

## セットアップスクリプトの実装例（ai_setup.sh）

これらの設定を一括で適用するためのスクリプト `ai_setup.sh` の例です。
既存のファイルがある場合はバックアップを取った上で、シンボリックリンクを作成するようにしています。

```bash
#!/bin/bash
set -e

DOTFILES_DIR=~/dotfiles
# バックアップディレクトリ
BACKUP_DIR=~/.ai_config_backup/$(date +%Y%m%d_%H%M%S)

# シンボリックリンク作成関数
create_symlink() {
    local src=$1
    local dest=$2
    local dest_dir=$(dirname "$dest")

    mkdir -p "$dest_dir"

    if [ -L "$dest" ]; then
        unlink "$dest"
    elif [ -e "$dest" ]; then
        # 既存ファイルはバックアップ
        mkdir -p "$BACKUP_DIR"
        cp -r "$dest" "$BACKUP_DIR/"
        rm -rf "$dest"
    fi

    ln -sf "$src" "$dest"
    echo "  ✓ $dest -> $src"
}

# Claude設定
setup_claude() {
    echo "Claude Code設定..."
    create_symlink "$DOTFILES_DIR/ai/claude/settings.json" ~/.claude/settings.json
    create_symlink "$DOTFILES_DIR/ai/claude/CLAUDE.md" ~/.claude/CLAUDE.md
}

# Cursor設定
setup_cursor() {
    echo "Cursor設定..."
    # macOS標準パスへのリンク
    local cursor_user_dir=~/Library/Application\ Support/Cursor/User
    create_symlink "$DOTFILES_DIR/ai/cursor/User/settings.json" "$cursor_user_dir/settings.json"
    create_symlink "$DOTFILES_DIR/ai/cursor/mcp.json" ~/.cursor/mcp.json

    # 拡張機能のインストール処理（省略）
}

# メイン処理
setup_claude
setup_cursor
# setup_codex, setup_gemini...
```

## 設定ファイルの中身（GEMINI.md の例）

実際に管理している `GEMINI.md` の一部を紹介します。
ここでは「出力言語の指定」や「やってはいけないこと（禁止事項）」を明記しています。

```markdown
# GEMINI.md

## Output Language
- **Always respond in Japanese.** (常に日本語で応答すること)
- Technical terms and API names may remain in English.
- Write implementation plans and artifacts in Japanese.

## Strict Prohibitions
1. **No hardcoded secrets** - 環境変数を使用すること。
2. **No debugging code** - コミット前にデバッグ出力を削除すること。
```

これをグローバル設定として読み込ませることで、どのプロジェクトを開いていても、「日本語で答えて」といちいち指示する必要がなくなります。

## まとめ

エディタの設定やシェルの設定と同じように、**AIへの指示（プロンプト）も「設定ファイル」としてコード管理する**時代になりつつあります。

自分好みにチューニングされた「最強のAIアシスタント」を、Gitに乗せてどこへでも持ち運べるようにしてみてはいかがでしょうか。
PC移行のストレスがまた一つ減ること間違いなしです。
