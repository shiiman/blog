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

var postCmd = &cobra.Command{
	Use:   "post <file>",
	Short: "新しい投稿を作成",
	Long: `Markdownファイルから新しい投稿を作成します。
デフォルトでは下書きとして保存されます。

例:
  wp-cli post drafts/2025-01-03_my-article/article.md
  wp-cli post drafts/article.md --publish
  wp-cli post drafts/article.md --dry-run`,
	Args: cobra.ExactArgs(1),
	Run:  runPost,
}

var postPublish bool
var postDryRun bool

func init() {
	rootCmd.AddCommand(postCmd)
	postCmd.Flags().BoolVarP(&postPublish, "publish", "p", false, "公開状態で投稿")
	postCmd.Flags().BoolVar(&postDryRun, "dry-run", false, "投稿せずに内容を確認")
}

func runPost(cmd *cobra.Command, args []string) {
	filePath := args[0]

	// 記事ファイルを解析
	article, err := converter.ParseArticle(filePath)
	if err != nil {
		color.Red("記事ファイルの解析に失敗: %v", err)
		os.Exit(1)
	}

	// ステータスの決定
	status := "draft"
	if postPublish {
		status = "publish"
	}
	if article.FrontMatter.Status != "" && !postPublish {
		status = article.FrontMatter.Status
	}

	// MarkdownをHTMLに変換
	htmlContent := converter.MarkdownToHTML(article.Content)

	// ドライラン
	if postDryRun {
		color.Yellow("=== ドライラン モード ===")
		color.White("タイトル: %s", article.FrontMatter.Title)
		color.White("スラッグ: %s", article.FrontMatter.Slug)
		color.White("ステータス: %s", status)
		color.White("カテゴリ: %v", article.FrontMatter.Categories)
		color.White("タグ: %v", article.FrontMatter.Tags)
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

	// 投稿作成リクエスト
	req := &types.CreatePostRequest{
		Title:         article.FrontMatter.Title,
		Content:       htmlContent,
		Status:        status,
		Slug:          article.FrontMatter.Slug,
		Excerpt:       article.FrontMatter.Excerpt,
		Categories:    article.FrontMatter.Categories,
		Tags:          article.FrontMatter.Tags,
		FeaturedMedia: article.FrontMatter.FeaturedMedia,
	}

	color.Cyan("投稿を作成中...")

	post, err := client.CreatePost(req)
	if err != nil {
		color.Red("投稿の作成に失敗: %v", err)
		os.Exit(1)
	}

	color.Green("投稿が作成されました！")
	color.White("  ID: %d", post.ID)
	color.White("  URL: %s", post.Link)
	color.White("  ステータス: %s", formatStatus(post.Status))
}
