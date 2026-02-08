package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/shiimanblog/wp-cli/internal/converter"
	"github.com/shiimanblog/wp-cli/internal/types"
	"github.com/spf13/cobra"
)

const eyecatchFilename = "eyecatch.png"

var postCmd = &cobra.Command{
	Use:   "post <file>",
	Short: "投稿を作成または更新",
	Long: `Markdownファイルから投稿を作成または更新します。
フロントマターにIDがある場合は既存記事を更新します。
デフォルトでは下書きとして保存されます。

例:
  wp-cli post drafts/2025-01-03_my-article/article.md
  wp-cli post drafts/article.md --publish
  wp-cli post drafts/article.md --dry-run`,
	Args: cobra.ExactArgs(1),
	RunE: runPost,
}

var postPublish bool
var postDryRun bool

func init() {
	rootCmd.AddCommand(postCmd)
	postCmd.Flags().BoolVarP(&postPublish, "publish", "p", false, "公開状態で投稿")
	postCmd.Flags().BoolVar(&postDryRun, "dry-run", false, "投稿せずに内容を確認")
}

func runPost(cmd *cobra.Command, args []string) error {
	filePath := args[0]

	// 記事ファイルを解析
	article, err := converter.ParseArticle(filePath)
	if err != nil {
		return fmt.Errorf("記事ファイルの解析に失敗: %w", err)
	}

	// ステータスの決定
	status := determineStatus(postPublish, article.FrontMatter.Status)

	// MarkdownをHTMLに変換
	htmlContent := converter.MarkdownToHTML(article.Content)

	// ドライラン
	if postDryRun {
		showDryRunPreview(
			article.FrontMatter.Title,
			article.FrontMatter.Slug,
			status,
			htmlContent,
			article.FrontMatter.Categories,
			article.FrontMatter.Tags,
		)
		return nil
	}

	// 設定読み込みとクライアント生成
	client, err := setupClient()
	if err != nil {
		return err
	}

	// アイキャッチ画像のアップロード
	articleDir := filepath.Dir(filePath)
	featuredMediaID, err := uploadEyecatchIfExists(client, articleDir, article.FrontMatter.FeaturedMedia, false)
	if err != nil {
		return err
	}

	var post *types.Post

	// IDがある場合は更新、ない場合は新規作成
	if article.FrontMatter.ID > 0 {
		// 既存記事を更新
		req := &types.UpdatePostRequest{
			Title:         article.FrontMatter.Title,
			Content:       htmlContent,
			Status:        status,
			Slug:          article.FrontMatter.Slug,
			Excerpt:       article.FrontMatter.Excerpt,
			Categories:    article.FrontMatter.Categories,
			Tags:          article.FrontMatter.Tags,
			FeaturedMedia: featuredMediaID,
		}

		color.Cyan("投稿 %d を更新中...", article.FrontMatter.ID)

		post, err = client.UpdatePost(article.FrontMatter.ID, req)
		if err != nil {
			return fmt.Errorf("投稿の更新に失敗: %w", err)
		}

	} else {
		// 新規作成
		req := &types.CreatePostRequest{
			Title:         article.FrontMatter.Title,
			Content:       htmlContent,
			Status:        status,
			Slug:          article.FrontMatter.Slug,
			Excerpt:       article.FrontMatter.Excerpt,
			Categories:    article.FrontMatter.Categories,
			Tags:          article.FrontMatter.Tags,
			FeaturedMedia: featuredMediaID,
		}

		color.Cyan("投稿を作成中...")

		post, err = client.CreatePost(req)
		if err != nil {
			return fmt.Errorf("投稿の作成に失敗: %w", err)
		}
	}

	printPostResult(post, article.FrontMatter.ID > 0)

	// 公開時にdrafts/からposts/に移動
	if status == "publish" && strings.Contains(filePath, "drafts/") {
		moveToPublished(filePath, post)
	}
	return nil
}

// moveToPublished はdrafts/からposts/に記事を移動する
func moveToPublished(filePath string, post *types.Post) {
	// ディレクトリパスを取得
	srcDir := filepath.Dir(filePath)

	// 新しいディレクトリ名: YYYY-MM-DD_slug
	newDirName := post.Date.Format("2006-01-02") + "_" + post.Slug
	destDir := filepath.Clean(filepath.Join("posts", newDirName))

	// パストラバーサル対策: 移動先がposts/配下であることを検証
	absDestDir, err := filepath.Abs(destDir)
	if err != nil {
		color.Yellow("  パスの解決に失敗: %v", err)
		return
	}
	absPosts, err := filepath.Abs("posts")
	if err != nil {
		color.Yellow("  パスの解決に失敗: %v", err)
		return
	}
	if !strings.HasPrefix(absDestDir, absPosts+string(filepath.Separator)) {
		color.Yellow("  不正なパスです（posts/ディレクトリ外への移動）: %s", destDir)
		return
	}

	// 移動先が既に存在する場合はスキップ
	if _, err := os.Stat(destDir); err == nil {
		color.Yellow("  移動先が既に存在します: %s", destDir)
		return
	}

	// 移動前に記事をパースして検証
	article, err := converter.ParseArticle(filePath)
	if err != nil {
		color.Yellow("  フロントマター読み込みに失敗: %v", err)
		return
	}

	// ディレクトリを移動
	if err := os.Rename(srcDir, destDir); err != nil {
		color.Yellow("  記事の移動に失敗 (%s -> %s): %v", srcDir, destDir, err)
		return
	}

	// フロントマターを更新（ID、status、featured_mediaを反映）
	articlePath := filepath.Join(destDir, filepath.Base(filePath))
	article.FrontMatter.ID = post.ID
	article.FrontMatter.Status = "publish"
	article.FrontMatter.Date = post.Date.Format("2006-01-02T15:04:05")
	article.FrontMatter.FeaturedMedia = post.FeaturedMedia

	content, err := converter.GenerateArticleFile(article)
	if err != nil {
		color.Yellow("  記事ファイル生成に失敗: %v", err)
		return
	}

	if err := os.WriteFile(articlePath, []byte(content), 0600); err != nil {
		color.Yellow("  記事ファイル保存に失敗: %v", err)
		return
	}

	color.Green("  記事を移動しました: %s", destDir)
	fmt.Println()
}
