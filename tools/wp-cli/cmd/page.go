package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/shiimanblog/wp-cli/internal/converter"
	"github.com/shiimanblog/wp-cli/internal/types"
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
	RunE: runPage,
}

var pagePublish bool
var pageDryRun bool

func init() {
	rootCmd.AddCommand(pageCmd)
	pageCmd.Flags().BoolVarP(&pagePublish, "publish", "p", false, "公開状態で投稿")
	pageCmd.Flags().BoolVar(&pageDryRun, "dry-run", false, "投稿せずに内容を確認")
}

func runPage(cmd *cobra.Command, args []string) error {
	filePath := args[0]

	// 記事ファイルを解析
	article, err := converter.ParseArticle(filePath)
	if err != nil {
		return fmt.Errorf("記事ファイルの解析に失敗: %w", err)
	}

	// ステータスの決定
	status := determineStatus(pagePublish, article.FrontMatter.Status)

	// MarkdownをHTMLに変換
	htmlContent := converter.MarkdownToHTML(article.Content)

	// ドライラン
	if pageDryRun {
		showPageDryRunPreview(
			article.FrontMatter.Title,
			article.FrontMatter.Slug,
			status,
			htmlContent,
			article.FrontMatter.Parent,
			article.FrontMatter.MenuOrder,
		)
		return nil
	}

	// 設定読み込みとクライアント生成
	client, err := setupClient()
	if err != nil {
		return err
	}

	ctx := cmd.Context()

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

	page, err := client.CreatePage(ctx, req)
	if err != nil {
		return fmt.Errorf("固定ページの作成に失敗: %w", err)
	}

	printPageResult(page, false)
	return nil
}
