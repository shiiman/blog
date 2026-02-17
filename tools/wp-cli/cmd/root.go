package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "wp-cli",
	Short: "WordPress記事管理CLI",
	Long: `WordPress記事管理CLI for shiimanblog.com

このCLIツールを使用して、WordPressの記事・固定ページを管理できます。

機能:
  - 記事・固定ページの一覧表示
  - 記事・固定ページのインポート（WordPress → ローカルMarkdown）
  - 記事・固定ページの投稿（ローカルMarkdown → WordPress）
  - 記事・固定ページの更新
  - カテゴリ・タグの一覧表示

使用前に .env ファイルにWordPress APIの認証情報を設定してください。`,
	Version: version,
}

// version はビルド時に -ldflags で注入される（例: -ldflags "-X github.com/shiimanblog/wp-cli/cmd.version=1.2.0"）
var version = "dev"

// Execute はルートコマンドを実行する
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	// グローバルフラグを追加
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "詳細出力を有効にする")
}
