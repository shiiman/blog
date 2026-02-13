package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/shiimanblog/wp-cli/internal/config"
	"github.com/shiimanblog/wp-cli/internal/types"
	"github.com/shiimanblog/wp-cli/internal/wp"
)

// determineStatus はpublishフラグとFront Matterからステータスを決定する
func determineStatus(publishFlag bool, frontMatterStatus string) string {
	status := "draft"
	if publishFlag {
		status = "publish"
	}
	if frontMatterStatus != "" && !publishFlag {
		status = frontMatterStatus
	}
	return status
}

// showDryRunPreview はドライランモードで記事の内容をプレビュー表示する
func showDryRunPreview(title, slug, status, htmlContent string, categories, tags []int) {
	color.Yellow("=== ドライラン モード ===")
	color.White("タイトル: %s", title)
	color.White("スラッグ: %s", slug)
	color.White("ステータス: %s", status)
	if len(categories) > 0 {
		color.White("カテゴリ: %v", categories)
	}
	if len(tags) > 0 {
		color.White("タグ: %v", tags)
	}
	color.White("\n--- 本文（HTML）プレビュー ---")
	preview := htmlContent
	if len(preview) > 500 {
		preview = preview[:500] + "..."
	}
	color.White("%s", preview)
}

// showPageDryRunPreview は固定ページ用のドライランプレビューを表示する
func showPageDryRunPreview(title, slug, status, htmlContent string, parent, menuOrder int) {
	color.Yellow("=== ドライラン モード ===")
	color.White("タイトル: %s", title)
	color.White("スラッグ: %s", slug)
	color.White("ステータス: %s", status)
	color.White("親ページID: %d", parent)
	color.White("メニュー順序: %d", menuOrder)
	color.White("\n--- 本文（HTML）プレビュー ---")
	preview := htmlContent
	if len(preview) > 500 {
		preview = preview[:500] + "..."
	}
	color.White("%s", preview)
}

// setupClient は設定を読み込んでWordPressクライアントを生成する
func setupClient() (*wp.Client, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("設定エラー: %w", err)
	}
	return wp.NewClient(cfg), nil
}

// uploadEyecatchIfExists はアイキャッチ画像が存在する場合にアップロードする
// forceUpload が true の場合は、既にFeaturedMediaが設定されていても再アップロードする
func uploadEyecatchIfExists(ctx context.Context, client *wp.Client, articleDir string, currentFeaturedMediaID int, forceUpload bool) (int, error) {
	eyecatchPath := filepath.Join(articleDir, "assets", eyecatchFilename)

	// ファイル存在確認
	if _, err := os.Stat(eyecatchPath); err != nil {
		// ファイルが存在しない場合は現在のIDをそのまま返す
		return currentFeaturedMediaID, nil
	}

	// 既に設定されている場合は強制アップロードでない限りスキップ
	if currentFeaturedMediaID != 0 && !forceUpload {
		return currentFeaturedMediaID, nil
	}

	// アップロード実行
	if forceUpload {
		color.Cyan("アイキャッチ画像を強制的に再アップロード中...")
	} else {
		color.Cyan("アイキャッチ画像をアップロード中...")
	}

	imageData, err := os.ReadFile(eyecatchPath)
	if err != nil {
		return 0, fmt.Errorf("アイキャッチ画像の読み込みに失敗: %w", err)
	}

	media, err := client.UploadMedia(ctx, eyecatchFilename, imageData, "image/png")
	if err != nil {
		return 0, fmt.Errorf("アイキャッチ画像のアップロードに失敗: %w", err)
	}

	color.Green("アイキャッチ画像をアップロードしました！")
	color.White("  メディアID: %d", media.ID)
	color.White("  URL: %s", media.SourceURL)

	return media.ID, nil
}

// formatStatus はステータス文字列を色付きの日本語に変換する
func formatStatus(status string) string {
	switch status {
	case "publish":
		return color.GreenString("公開")
	case "draft":
		return color.YellowString("下書き")
	case "pending":
		return color.CyanString("保留中")
	case "private":
		return color.MagentaString("非公開")
	default:
		return status
	}
}

// truncate は文字列を指定された長さに切り詰める
func truncate(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen-3]) + "..."
}

// printPostResult は投稿作成/更新結果を表示する
func printPostResult(post *types.Post, isUpdate bool) {
	if isUpdate {
		color.Green("投稿が更新されました！")
	} else {
		color.Green("投稿が作成されました！")
	}
	color.White("  ID: %d", post.ID)
	color.White("  URL: %s", post.Link)
	color.White("  ステータス: %s", formatStatus(post.Status))
}

// printPageResult は固定ページ作成/更新結果を表示する
func printPageResult(page *types.Page, isUpdate bool) {
	if isUpdate {
		color.Green("固定ページが更新されました！")
	} else {
		color.Green("固定ページが作成されました！")
	}
	color.White("  ID: %d", page.ID)
	color.White("  URL: %s", page.Link)
	color.White("  ステータス: %s", formatStatus(page.Status))
}
