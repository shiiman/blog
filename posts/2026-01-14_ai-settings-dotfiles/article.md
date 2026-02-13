---
title: "AIのグローバル設定をdotfilesで管理する方法：PC移行でも自分専用のAIアシスタントを持ち歩く"
id: 2213
slug: "ai-settings-dotfiles"
status: publish
categories: [18, 133]
tags: [101, 102]
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
├── AGENTS.md         # プロジェクト共通のAIエージェント指針（新設）
├── CLAUDE.md         # プロジェクト固有の指示書（AI用）
├── ai_setup.sh       # AI設定用の一括セットアップスクリプト
└── ai/
    ├── claude/
    │   ├── CLAUDE.md      # グローバルな指示（~/.claude/ へ）
    │   └── settings.json
    ├── codex/
    │   ├── config.toml
    │   └── skills/
    ├── cursor/
    │   ├── User/
    │   │   ├── settings.json
    │   │   ├── keybindings.json
    │   │   └── snippets/   # コードスニペット管理
    │   ├── extensions.json # 拡張機能リスト（自動インストール用）
    │   └── mcp.json
    └── antigravity/       # (旧gemini)名称変更
        ├── GEMINI.md      # グローバルな指示（~/.gemini/ へ）
        └── extensions.json
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

特に `setup.sh` では、`extensions.json` に記載された拡張機能IDを `jq` で読み取り、`cursor --install-extension` コマンドで自動インストールする仕組みも導入しました。これにより、プラグインまで含めた完全な再現が可能になります。

### 3. Codex (OpenAI)
Codex系のCLIツールや独自のスクリプトで使用する設定です。

- **管理ファイル**:
    - `config.toml`: APIキーの参照先やモデルパラメータ（temperatureなど）の設定。
    - `skills/`: 再利用可能なスキル（プロンプトテンプレート）の定義フォルダ。

- **シンボリックリンク先**:
    - `~/.codex/config.toml`
    - `~/.codex/skills/`

### 4. Antigravity (Gemini)
GoogleのAdvanced Agentic Codingエージェント用設定です。最近 `gemini` からより汎用的な `antigravity` へと名称変更されました。

- **管理ファイル**:
    - `GEMINI.md`: エージェントへの基本指示書。
    - `extensions.json`: Antigravity用の拡張機能リスト。

- **シンボリックリンク先**:
    - `~/.gemini/GEMINI.md`
    - `~/.antigravity/extensions/extensions.json`

Cursorと同様、CLIツールを使って拡張機能の自動インストールもサポートしています。

## セットアップスクリプトの実装例（ai_setup.sh）

これらの設定を一括で適用するためのスクリプト `ai_setup.sh` の例です。
既存のファイルがある場合はバックアップを取った上で、シンボリックリンクを作成するようにしています。

```bash
#!/bin/bash
set -e

DOTFILES_DIR=~/dotfiles
BACKUP_DIR=~/.ai_config_backup/$(date +%Y%m%d_%H%M%S)

# シンボリックリンク作成関数（バックアップ機能付き）
create_symlink() {
    local src=$1
    local dest=$2
    mkdir -p "$(dirname "$dest")"

    if [ -L "$dest" ]; then
        unlink "$dest"
    elif [ -e "$dest" ]; then
        mkdir -p "$BACKUP_DIR"
        cp -r "$dest" "$BACKUP_DIR/"
        rm -rf "$dest"
    fi

    ln -sf "$src" "$dest"
}

# Cursor設定：設定ファイルのリンクと拡張機能の自動インストール
setup_cursor() {
    echo "Cursor設定..."
    local cursor_user_dir=~/Library/Application\ Support/Cursor/User
    create_symlink "$DOTFILES_DIR/ai/cursor/User/settings.json" "$cursor_user_dir/settings.json"
    create_symlink "$DOTFILES_DIR/ai/cursor/User/keybindings.json" "$cursor_user_dir/keybindings.json"
    create_symlink "$DOTFILES_DIR/ai/cursor/mcp.json" ~/.cursor/mcp.json

    # 拡張機能の自動インストール
    if command -v jq &> /dev/null && command -v cursor &> /dev/null; then
        local ext_ids=$(jq -r '.[].identifier.id' "$DOTFILES_DIR/ai/cursor/extensions.json")
        for ext_id in $ext_ids; do
            cursor --install-extension "$ext_id" > /dev/null 2>&1 && echo "  ✓ $ext_id"
        done
    fi
}

# Antigravity設定
setup_antigravity() {
    echo "Antigravity設定..."
    create_symlink "$DOTFILES_DIR/ai/antigravity/GEMINI.md" ~/.gemini/GEMINI.md
    # 拡張機能のインストールも同様に行う
}

# メイン処理
setup_cursor
setup_antigravity
# ...他ツールの設定
```

## 設定ファイルの中身（GEMINI.md の例）

実際に管理している `GEMINI.md` の一部を紹介します。出力言語の指定、禁止事項に加えて、コミットやPRのガイドラインも追加されました。

```markdown
# GEMINI.md

## Output Language
- Always respond in Japanese.
- Write code comments and error messages in Japanese.

## Strict Prohibitions
1. **No hardcoded secrets** - 環境変数を使用すること。
2. **No destructive operations without confirmation** - 破壊的な操作は必ず確認。

## Commit Guidelines
- Write commit messages in Japanese, concise and descriptive in one line.
```

## プロジェクト固有のAIガイドライン（AGENTS.md）

今回のアップデートで最も特徴的なのが、リポジトリのルートに配置した `AGENTS.md` や `CLAUDE.md` です。

これらは AI アシスタント自身がそのプロジェクトでどう振る舞うべきかを定めた「プロジェクトの憲法」です。
dotfiles 自体の開発を AI に手伝わせる際、「Shellcheck を通すこと」「`.zshrc` は macOS の標準 Bash でも動くようにすること」といったルールをここに書いておくことで、AI が常にプロジェクトの文脈に沿った提案をしてくれるようになります。

## まとめ

エディタの設定やシェルの設定と同じように、**AIへの指示（プロンプト）も「設定ファイル」としてコード管理する**時代になりつつあります。

自分好みにチューニングされた「最強のAIアシスタント」を、Gitに乗せてどこへでも持ち運べるようにしてみてはいかがでしょうか。
PC移行のストレスがまた一つ減ること間違いなしです。
