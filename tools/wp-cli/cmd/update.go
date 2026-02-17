package cmd

import (
	"context"
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
  wp-cli update drafts/article.md --id=123`,
	Args: cobra.ExactArgs(1),
	RunE: runUpdate,
}

var updateID int
var updateDryRun bool
var updatePage bool
var updateForceEyecatch bool

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.Flags().IntVar(&updateID, "id", 0, "更新する投稿/固定ページのID（Front Matterを上書き）")
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

	// ステータスはFront Matterを維持（updateは更新専用）
	status, err := normalizeAndValidateStatus(article.FrontMatter.Status, true)
	if err != nil {
		return fmt.Errorf("update の Front Matter status が不正です: %w", err)
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

	ctx := cmd.Context()

	// アイキャッチ画像のアップロード（投稿の場合のみ）
	featuredMediaID := article.FrontMatter.FeaturedMedia
	if !updatePage {
		articleDir := filepath.Dir(filePath)
		featuredMediaID, err = uploadEyecatchIfExists(ctx, client, articleDir, featuredMediaID, updateForceEyecatch)
		if err != nil {
			return err
		}
	}

	if updatePage {
		page, err := updatePageContent(ctx, client, id, article, htmlContent, status)
		if err != nil {
			return err
		}
		syncedPath, err := syncPageToLocal(filePath, page)
		if err != nil {
			return err
		}
		color.Green("  記事ファイルを同期しました: %s", syncedPath)
		return nil
	}

	post, err := updatePostContent(ctx, client, id, article, htmlContent, status, featuredMediaID)
	if err != nil {
		return err
	}

	syncedPath, err := syncPostToLocal(filePath, post, false)
	if err != nil {
		return err
	}

	color.Green("  記事ファイルを同期しました: %s", syncedPath)
	return nil
}

func updatePostContent(ctx context.Context, client *wp.Client, id int, article *types.Article, htmlContent, status string, featuredMediaID int) (*types.Post, error) {
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

	post, err := client.UpdatePost(ctx, id, req)
	if err != nil {
		return nil, fmt.Errorf("投稿の更新に失敗: %w", err)
	}

	printPostResult(post, true)
	return post, nil
}

func updatePageContent(ctx context.Context, client *wp.Client, id int, article *types.Article, htmlContent, status string) (*types.Page, error) {
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

	page, err := client.UpdatePage(ctx, id, req)
	if err != nil {
		return nil, fmt.Errorf("固定ページの更新に失敗: %w", err)
	}

	printPageResult(page, true)
	return page, nil
}
