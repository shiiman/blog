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
Front Matter の status が未指定の場合は下書きとして保存されます。
status の許可値: draft, publish, pending, private

例:
  wp-cli page drafts/about/page.md
  # Front Matter で公開したい場合:
  # status: publish  # draft | publish | pending | private
  wp-cli page drafts/page.md --dry-run`,
	Args: cobra.ExactArgs(1),
	RunE: runPage,
}

var pageDryRun bool

func init() {
	rootCmd.AddCommand(pageCmd)
	pageCmd.Flags().BoolVar(&pageDryRun, "dry-run", false, "投稿せずに内容を確認")
}

func runPage(cmd *cobra.Command, args []string) error {
	filePath := args[0]

	// 記事ファイルを解析
	article, err := converter.ParseArticle(filePath)
	if err != nil {
		return fmt.Errorf("記事ファイルの解析に失敗: %w", err)
	}

	// ステータスの決定（Front Matter優先、未指定時はdraft）
	status, err := determinePageStatus(article.FrontMatter.Status)
	if err != nil {
		return fmt.Errorf("page の status が不正です: %w", err)
	}

	// MarkdownをHTMLに変換
	htmlContent := converter.MarkdownToHTML(article.Content)

	// ドライラン
	if pageDryRun {
		showDryRunPreview(dryRunInfo{
			Title:          article.FrontMatter.Title,
			Slug:           article.FrontMatter.Slug,
			Status:         status,
			HTML:           htmlContent,
			Parent:         article.FrontMatter.Parent,
			MenuOrder:      article.FrontMatter.MenuOrder,
			ShowPageFields: true,
		})
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
