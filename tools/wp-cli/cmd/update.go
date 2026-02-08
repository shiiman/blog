package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/fatih/color"
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
	RunE: runUpdate,
}

var updateID int
var updatePublish bool
var updateDryRun bool
var updatePage bool
var updateForceEyecatch bool

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.Flags().IntVar(&updateID, "id", 0, "更新する投稿/固定ページのID（Front Matterを上書き）")
	updateCmd.Flags().BoolVarP(&updatePublish, "publish", "p", false, "公開状態に変更")
	updateCmd.Flags().BoolVar(&updateDryRun, "dry-run", false, "更新せずに内容を確認")
	updateCmd.Flags().BoolVar(&updatePage, "page", false, "固定ページとして更新")
	updateCmd.Flags().BoolVar(&updateForceEyecatch, "force-eyecatch", false, "アイキャッチ画像を強制的に再アップロード")
}

func runUpdate(cmd *cobra.Command, args []string) error {
	filePath := args[0]

	// 記事ファイルを解析
	article, err := converter.ParseArticle(filePath)
	if err != nil {
		return fmt.Errorf("記事ファイルの解析に失敗: %w", err)
	}

	// IDの決定
	id := article.FrontMatter.ID
	if updateID > 0 {
		id = updateID
	}
	if id == 0 {
		return fmt.Errorf("投稿IDが指定されていません。Front Matterの 'id' フィールドまたは --id オプションで指定してください")
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
		return nil
	}

	// 設定読み込みとクライアント生成
	client, err := setupClient()
	if err != nil {
		return err
	}

	// アイキャッチ画像のアップロード（投稿の場合のみ）
	featuredMediaID := article.FrontMatter.FeaturedMedia
	if !updatePage {
		articleDir := filepath.Dir(filePath)
		featuredMediaID, err = uploadEyecatchIfExists(client, articleDir, featuredMediaID, updateForceEyecatch)
		if err != nil {
			return err
		}
	}

	if updatePage {
		return updatePageContent(client, id, article, htmlContent, status)
	}
	return updatePostContent(client, id, article, htmlContent, status, featuredMediaID)
}

func updatePostContent(client *wp.Client, id int, article *types.Article, htmlContent, status string, featuredMediaID int) error {
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
		return fmt.Errorf("投稿の更新に失敗: %w", err)
	}

	printPostResult(post, true)
	return nil
}

func updatePageContent(client *wp.Client, id int, article *types.Article, htmlContent, status string) error {
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
		return fmt.Errorf("固定ページの更新に失敗: %w", err)
	}

	printPageResult(page, true)
	return nil
}
