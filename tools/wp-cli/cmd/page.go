package cmd

import (
	"os"

	"github.com/fatih/color"
	"github.com/shiimanblog/wp-cli/internal/config"
	"github.com/shiimanblog/wp-cli/internal/converter"
	"github.com/shiimanblog/wp-cli/internal/types"
	"github.com/shiimanblog/wp-cli/internal/wp"
	"github.com/spf13/cobra"
)

var pageCmd = &cobra.Command{
	Use:   "page <file>",
	Short: "新しい固定ページを作成",
	Long: `Markdownファイルから新しい固定ページを作成します。
デフォルトでは下書きとして保存されます。

例:
  wp-cli page drafts/about/page.md
  wp-cli page drafts/page.md --publish
  wp-cli page drafts/page.md --dry-run`,
	Args: cobra.ExactArgs(1),
	Run:  runPage,
}

var pagePublish bool
var pageDryRun bool

func init() {
	rootCmd.AddCommand(pageCmd)
	pageCmd.Flags().BoolVarP(&pagePublish, "publish", "p", false, "公開状態で投稿")
	pageCmd.Flags().BoolVar(&pageDryRun, "dry-run", false, "投稿せずに内容を確認")
}

func runPage(cmd *cobra.Command, args []string) {
	filePath := args[0]

	// 記事ファイルを解析
	article, err := converter.ParseArticle(filePath)
	if err != nil {
		color.Red("記事ファイルの解析に失敗: %v", err)
		os.Exit(1)
	}

	// ステータスの決定
	status := "draft"
	if pagePublish {
		status = "publish"
	}
	if article.FrontMatter.Status != "" && !pagePublish {
		status = article.FrontMatter.Status
	}

	// MarkdownをHTMLに変換
	htmlContent := converter.MarkdownToHTML(article.Content)

	// ドライラン
	if pageDryRun {
		color.Yellow("=== ドライラン モード ===")
		color.White("タイトル: %s", article.FrontMatter.Title)
		color.White("スラッグ: %s", article.FrontMatter.Slug)
		color.White("ステータス: %s", status)
		color.White("親ページID: %d", article.FrontMatter.Parent)
		color.White("メニュー順序: %d", article.FrontMatter.MenuOrder)
		color.White("\n--- 本文（HTML）プレビュー ---")
		preview := htmlContent
		if len(preview) > 500 {
			preview = preview[:500] + "..."
		}
		color.White("%s", preview)
		return
	}

	// 設定読み込み
	cfg, err := config.Load()
	if err != nil {
		color.Red("設定エラー: %v", err)
		os.Exit(1)
	}

	client := wp.NewClient(cfg)

	// 固定ページ作成リクエスト
	req := &types.CreatePageRequest{
		Title:     article.FrontMatter.Title,
		Content:   htmlContent,
		Status:    status,
		Slug:      article.FrontMatter.Slug,
		Excerpt:   article.FrontMatter.Excerpt,
		Parent:    article.FrontMatter.Parent,
		MenuOrder: article.FrontMatter.MenuOrder,
	}

	color.Cyan("固定ページを作成中...")

	page, err := client.CreatePage(req)
	if err != nil {
		color.Red("固定ページの作成に失敗: %v", err)
		os.Exit(1)
	}

	color.Green("固定ページが作成されました！")
	color.White("  ID: %d", page.ID)
	color.White("  URL: %s", page.Link)
	color.White("  ステータス: %s", formatStatus(page.Status))
}
