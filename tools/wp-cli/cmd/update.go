package cmd

import (
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/shiimanblog/wp-cli/internal/config"
	"github.com/shiimanblog/wp-cli/internal/converter"
	"github.com/shiimanblog/wp-cli/internal/types"
	"github.com/shiimanblog/wp-cli/internal/wp"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update <file>",
	Short: "既存の投稿または固定ページを更新",
	Long: `Markdownファイルの内容で既存の投稿または固定ページを更新します。
投稿IDはFront Matterの 'id' フィールドから取得されます。

例:
  wp-cli update posts/2025-01-03_my-article/article.md
  wp-cli update pages/about/page.md
  wp-cli update drafts/article.md --id=123
  wp-cli update posts/article.md --publish`,
	Args: cobra.ExactArgs(1),
	Run:  runUpdate,
}

var updateID int
var updatePublish bool
var updateDryRun bool
var updatePage bool

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.Flags().IntVar(&updateID, "id", 0, "更新する投稿/固定ページのID（Front Matterを上書き）")
	updateCmd.Flags().BoolVarP(&updatePublish, "publish", "p", false, "公開状態に変更")
	updateCmd.Flags().BoolVar(&updateDryRun, "dry-run", false, "更新せずに内容を確認")
	updateCmd.Flags().BoolVar(&updatePage, "page", false, "固定ページとして更新")
}

func runUpdate(cmd *cobra.Command, args []string) {
	filePath := args[0]

	// 記事ファイルを解析
	article, err := converter.ParseArticle(filePath)
	if err != nil {
		color.Red("記事ファイルの解析に失敗: %v", err)
		os.Exit(1)
	}

	// IDの決定
	id := article.FrontMatter.ID
	if updateID > 0 {
		id = updateID
	}
	if id == 0 {
		color.Red("投稿IDが指定されていません。Front Matterの 'id' フィールドまたは --id オプションで指定してください。")
		os.Exit(1)
	}

	// ステータスの決定
	status := article.FrontMatter.Status
	if updatePublish {
		status = "publish"
	}

	// MarkdownをHTMLに変換
	htmlContent := converter.MarkdownToHTML(article.Content)

	// ドライラン
	if updateDryRun {
		color.Yellow("=== ドライラン モード ===")
		color.White("更新対象ID: %d", id)
		color.White("タイトル: %s", article.FrontMatter.Title)
		color.White("スラッグ: %s", article.FrontMatter.Slug)
		color.White("ステータス: %s", status)
		if !updatePage {
			color.White("カテゴリ: %v", article.FrontMatter.Categories)
			color.White("タグ: %v", article.FrontMatter.Tags)
		}
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

	// アイキャッチ画像のアップロード（投稿の場合のみ）
	featuredMediaID := article.FrontMatter.FeaturedMedia
	if !updatePage {
		articleDir := filepath.Dir(filePath)
		eyecatchPath := filepath.Join(articleDir, "assets", eyecatchFilename)

		if _, err := os.Stat(eyecatchPath); err == nil && featuredMediaID == 0 {
			// アイキャッチ画像が存在し、まだ設定されていない場合
			color.Cyan("アイキャッチ画像をアップロード中...")

			imageData, err := os.ReadFile(eyecatchPath)
			if err != nil {
				color.Red("アイキャッチ画像の読み込みに失敗: %v", err)
				os.Exit(1)
			}

			media, err := client.UploadMedia(eyecatchFilename, imageData, "image/png")
			if err != nil {
				color.Red("アイキャッチ画像のアップロードに失敗: %v", err)
				os.Exit(1)
			}

			featuredMediaID = media.ID
			color.Green("アイキャッチ画像をアップロードしました！")
			color.White("  メディアID: %d", media.ID)
			color.White("  URL: %s", media.SourceURL)
		}
	}

	if updatePage {
		updatePageContent(client, id, article, htmlContent, status)
	} else {
		updatePostContent(client, id, article, htmlContent, status, featuredMediaID)
	}
}

func updatePostContent(client *wp.Client, id int, article *types.Article, htmlContent, status string, featuredMediaID int) {
	req := &types.UpdatePostRequest{
		Title:         article.FrontMatter.Title,
		Content:       htmlContent,
		Slug:          article.FrontMatter.Slug,
		Excerpt:       article.FrontMatter.Excerpt,
		Categories:    article.FrontMatter.Categories,
		Tags:          article.FrontMatter.Tags,
		FeaturedMedia: featuredMediaID,
	}
	if status != "" {
		req.Status = status
	}

	color.Cyan("投稿 %d を更新中...", id)

	post, err := client.UpdatePost(id, req)
	if err != nil {
		color.Red("投稿の更新に失敗: %v", err)
		os.Exit(1)
	}

	color.Green("投稿が更新されました！")
	color.White("  ID: %d", post.ID)
	color.White("  URL: %s", post.Link)
	color.White("  ステータス: %s", formatStatus(post.Status))
}

func updatePageContent(client *wp.Client, id int, article *types.Article, htmlContent, status string) {
	req := &types.UpdatePageRequest{
		Title:     article.FrontMatter.Title,
		Content:   htmlContent,
		Slug:      article.FrontMatter.Slug,
		Excerpt:   article.FrontMatter.Excerpt,
		Parent:    article.FrontMatter.Parent,
		MenuOrder: article.FrontMatter.MenuOrder,
	}
	if status != "" {
		req.Status = status
	}

	color.Cyan("固定ページ %d を更新中...", id)

	page, err := client.UpdatePage(id, req)
	if err != nil {
		color.Red("固定ページの更新に失敗: %v", err)
		os.Exit(1)
	}

	color.Green("固定ページが更新されました！")
	color.White("  ID: %d", page.ID)
	color.White("  URL: %s", page.Link)
	color.White("  ステータス: %s", formatStatus(page.Status))
}
