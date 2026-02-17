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

// supportedEyecatchExts はサポートするアイキャッチ画像の拡張子（優先順）
var supportedEyecatchExts = []string{".png", ".jpg", ".jpeg", ".webp", ".gif"}

const allowedStatusHint = "draft/publish/pending/private"

var allowedStatuses = map[string]struct{}{
	"draft":   {},
	"publish": {},
	"pending": {},
	"private": {},
}

// determinePostStatus はpostコマンド向けにステータスを決定する
// write=draft / publish=publish の責務分離に合わせ、
// postはデフォルトで公開、--draft指定時のみ下書きにする。
func determinePostStatus(draftFlag bool) string {
	if draftFlag {
		return "draft"
	}
	return "publish"
}

// normalizeAndValidateStatus はstatus値を正規化して許可値か検証する。
func normalizeAndValidateStatus(status string, allowEmpty bool) (string, error) {
	normalized := strings.ToLower(strings.TrimSpace(status))
	if normalized == "" {
		if allowEmpty {
			return "", nil
		}
		return "", fmt.Errorf("status は %s のいずれかを指定してください", allowedStatusHint)
	}

	if _, ok := allowedStatuses[normalized]; !ok {
		return "", fmt.Errorf("status は %s のいずれかを指定してください (入力値: %q)", allowedStatusHint, status)
	}

	return normalized, nil
}

// determinePageStatus はpageコマンド向けにステータスを決定する
// Front Matterで明示されたstatusを優先し、未指定時はdraftにする。
func determinePageStatus(frontMatterStatus string) (string, error) {
	if strings.TrimSpace(frontMatterStatus) == "" {
		return "draft", nil
	}
	return normalizeAndValidateStatus(frontMatterStatus, false)
}

// dryRunInfo はドライランプレビューの表示情報を保持する
type dryRunInfo struct {
	UpdateID       int    // 更新対象ID（0の場合は表示しない）
	Title          string
	Slug           string
	Status         string
	HTML           string
	Categories     []int
	Tags           []int
	Parent         int
	MenuOrder      int
	ShowPageFields bool // 親ページ/メニュー順序を表示するか
}

// showDryRunPreview はドライランモードで記事の内容をプレビュー表示する
func showDryRunPreview(info dryRunInfo) {
	color.Yellow("=== ドライラン モード ===")
	if info.UpdateID > 0 {
		color.White("更新対象ID: %d", info.UpdateID)
	}
	color.White("タイトル: %s", info.Title)
	color.White("スラッグ: %s", info.Slug)
	color.White("ステータス: %s", info.Status)
	if len(info.Categories) > 0 {
		color.White("カテゴリ: %v", info.Categories)
	}
	if len(info.Tags) > 0 {
		color.White("タグ: %v", info.Tags)
	}
	if info.ShowPageFields {
		color.White("親ページID: %d", info.Parent)
		color.White("メニュー順序: %d", info.MenuOrder)
	}
	color.White("\n--- 本文（HTML）プレビュー ---")
	preview := info.HTML
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

// writeArticleFile はArticleをファイルに書き込む共通ヘルパー
func writeArticleFile(filePath string, article *types.Article) error {
	content, err := converter.GenerateArticleFile(article)
	if err != nil {
		return fmt.Errorf("記事ファイル生成に失敗: %w", err)
	}
	if err := os.WriteFile(filePath, []byte(content), 0600); err != nil {
		return fmt.Errorf("記事ファイル保存に失敗: %w", err)
	}
	return nil
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

	if err := writeArticleFile(currentPath, article); err != nil {
		return currentPath, err
	}
	return currentPath, nil
}

// syncPageToLocal は固定ページの更新結果をローカルMarkdownに同期する
func syncPageToLocal(filePath string, page *types.Page) (string, error) {
	article, err := converter.ParseArticle(filePath)
	if err != nil {
		return filePath, fmt.Errorf("フロントマター読み込みに失敗: %w", err)
	}

	article.FrontMatter.ID = page.ID
	article.FrontMatter.Status = page.Status
	article.FrontMatter.Date = page.Date.Format("2006-01-02T15:04:05")
	article.FrontMatter.Modified = page.Modified.Format("2006-01-02T15:04:05")
	article.FrontMatter.Slug = page.Slug
	article.FrontMatter.Parent = page.Parent
	article.FrontMatter.MenuOrder = page.MenuOrder

	if err := writeArticleFile(filePath, article); err != nil {
		return filePath, err
	}
	return filePath, nil
}

// findEyecatchFile はassetsディレクトリからアイキャッチ画像ファイルを検索する
func findEyecatchFile(assetsDir string) string {
	for _, ext := range supportedEyecatchExts {
		path := filepath.Join(assetsDir, "eyecatch"+ext)
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}
	return ""
}

// detectMIMEType はファイル拡張子からMIMEタイプを判定する
func detectMIMEType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".webp":
		return "image/webp"
	case ".gif":
		return "image/gif"
	default:
		return "application/octet-stream"
	}
}

// uploadEyecatchIfExists はアイキャッチ画像が存在する場合にアップロードする
// forceUpload が true の場合は、既にFeaturedMediaが設定されていても再アップロードする
func uploadEyecatchIfExists(ctx context.Context, client *wp.Client, articleDir string, currentFeaturedMediaID int, forceUpload bool) (int, error) {
	assetsDir := filepath.Join(articleDir, "assets")
	eyecatchPath := findEyecatchFile(assetsDir)

	if eyecatchPath == "" {
		// アイキャッチ画像が存在しない場合は現在のIDをそのまま返す
		return currentFeaturedMediaID, nil
	}

	// サイズチェック
	const maxEyecatchSize = 20 * 1024 * 1024 // 20MB
	info, err := os.Stat(eyecatchPath)
	if err != nil {
		return currentFeaturedMediaID, nil
	}
	if info.Size() > maxEyecatchSize {
		return 0, fmt.Errorf("アイキャッチ画像のサイズが上限（20MB）を超えています: %d bytes", info.Size())
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

	filename := filepath.Base(eyecatchPath)
	mimeType := detectMIMEType(filename)

	media, err := client.UploadMedia(ctx, filename, imageData, mimeType)
	if err != nil {
		return 0, fmt.Errorf("アイキャッチ画像のアップロードに失敗: %w", err)
	}

	color.Green("アイキャッチ画像をアップロードしました！")
	color.White("  メディアID: %d", media.ID)
	color.White("  URL: %s", media.SourceURL)

	return media.ID, nil
}

// sanitizeSlug はスラッグをサニタイズしてパストラバーサルを防止する
func sanitizeSlug(slug string) string {
	slug = strings.ReplaceAll(slug, "/", "")
	slug = strings.ReplaceAll(slug, "\\", "")
	slug = strings.ReplaceAll(slug, "..", "")
	return slug
}

// moveToPublished はdrafts/からposts/に記事を移動する
// ソースファイルのパスからプロジェクトルートを導出し、CWDに依存しない
func moveToPublished(filePath string, post *types.Post) (string, error) {
	srcDir := filepath.Dir(filePath)
	newDirName := post.Date.Format("2006-01-02") + "_" + sanitizeSlug(post.Slug)

	// ソースファイルの絶対パスからプロジェクトルートを導出
	absSrcDir, err := filepath.Abs(srcDir)
	if err != nil {
		return filePath, fmt.Errorf("パスの解決に失敗: %w", err)
	}

	sep := string(filepath.Separator)
	draftsMarker := sep + "drafts" + sep
	idx := strings.LastIndex(absSrcDir+sep, draftsMarker)
	if idx == -1 {
		return filePath, fmt.Errorf("draftsディレクトリが見つかりません: %s", absSrcDir)
	}
	projectRoot := absSrcDir[:idx]

	postsDir := filepath.Join(projectRoot, "posts")
	destDir := filepath.Join(postsDir, newDirName)

	// パストラバーサル防止
	if !strings.HasPrefix(destDir, postsDir+sep) {
		return filePath, fmt.Errorf("不正なパスです（posts/ディレクトリ外への移動）: %s", destDir)
	}

	if _, err := os.Stat(destDir); err == nil {
		return filePath, fmt.Errorf("移動先が既に存在します: %s", destDir)
	}

	if err := os.Rename(absSrcDir, destDir); err != nil {
		return filePath, fmt.Errorf("記事の移動に失敗 (%s -> %s): %w", absSrcDir, destDir, err)
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
	if maxLen <= 0 {
		return ""
	}
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	// "..." を付加する余裕がない場合はそのまま切り詰め
	if maxLen <= 3 {
		return string(runes[:maxLen])
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
