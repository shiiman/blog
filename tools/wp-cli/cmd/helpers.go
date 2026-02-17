package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/shiimanblog/wp-cli/internal/config"
	"github.com/shiimanblog/wp-cli/internal/converter"
	"github.com/shiimanblog/wp-cli/internal/types"
	"github.com/shiimanblog/wp-cli/internal/wp"
)

// determinePostStatus はpostコマンド向けにステータスを決定する
// write=draft / publish=publish の責務分離に合わせ、
// postはデフォルトで公開、--draft指定時のみ下書きにする。
func determinePostStatus(draftFlag bool) string {
	if draftFlag {
		return "draft"
	}
	return "publish"
}

// determinePageStatus はpageコマンド向けにステータスを決定する
func determinePageStatus(publishFlag bool, frontMatterStatus string) string {
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

// syncPostToLocal は投稿結果をローカルMarkdownに同期し、必要に応じて公開ディレクトリへ移動する
func syncPostToLocal(filePath string, post *types.Post, moveOnPublish bool) (string, error) {
	article, err := converter.ParseArticle(filePath)
	if err != nil {
		return filePath, fmt.Errorf("フロントマター読み込みに失敗: %w", err)
	}

	article.FrontMatter.ID = post.ID
	article.FrontMatter.Status = post.Status
	article.FrontMatter.Date = post.Date.Format("2006-01-02T15:04:05")
	article.FrontMatter.Modified = post.Modified.Format("2006-01-02T15:04:05")
	article.FrontMatter.Slug = post.Slug
	article.FrontMatter.FeaturedMedia = post.FeaturedMedia

	currentPath := filePath
	if moveOnPublish && post.Status == "publish" && strings.Contains(filePath, "drafts/") {
		currentPath, err = moveToPublished(filePath, post)
		if err != nil {
			return filePath, err
		}
	}

	content, err := converter.GenerateArticleFile(article)
	if err != nil {
		return currentPath, fmt.Errorf("記事ファイル生成に失敗: %w", err)
	}
	if err := os.WriteFile(currentPath, []byte(content), 0600); err != nil {
		return currentPath, fmt.Errorf("記事ファイル保存に失敗: %w", err)
	}

	return currentPath, nil
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

// moveToPublished はdrafts/からposts/に記事を移動する
func moveToPublished(filePath string, post *types.Post) (string, error) {
	srcDir := filepath.Dir(filePath)
	newDirName := post.Date.Format("2006-01-02") + "_" + post.Slug
	destDir := filepath.Clean(filepath.Join("posts", newDirName))

	absDestDir, err := filepath.Abs(destDir)
	if err != nil {
		return filePath, fmt.Errorf("パスの解決に失敗: %w", err)
	}
	absPosts, err := filepath.Abs("posts")
	if err != nil {
		return filePath, fmt.Errorf("パスの解決に失敗: %w", err)
	}
	if !strings.HasPrefix(absDestDir, absPosts+string(filepath.Separator)) {
		return filePath, fmt.Errorf("不正なパスです（posts/ディレクトリ外への移動）: %s", destDir)
	}

	if _, err := os.Stat(destDir); err == nil {
		return filePath, fmt.Errorf("移動先が既に存在します: %s", destDir)
	}

	if err := os.Rename(srcDir, destDir); err != nil {
		return filePath, fmt.Errorf("記事の移動に失敗 (%s -> %s): %w", srcDir, destDir, err)
	}

	color.Green("  記事を移動しました: %s", destDir)
	fmt.Println()

	return filepath.Join(destDir, filepath.Base(filePath)), nil
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
