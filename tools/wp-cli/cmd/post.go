package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/shiimanblog/wp-cli/internal/config"
	"github.com/shiimanblog/wp-cli/internal/converter"
	"github.com/shiimanblog/wp-cli/internal/types"
	"github.com/shiimanblog/wp-cli/internal/wp"
	"github.com/spf13/cobra"
)

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
			FeaturedMedia: article.FrontMatter.FeaturedMedia,
		}

		color.Cyan("投稿 %d を更新中...", article.FrontMatter.ID)

		post, err = client.UpdatePost(article.FrontMatter.ID, req)
		if err != nil {
			color.Red("投稿の更新に失敗: %v", err)
			os.Exit(1)
		}

		color.Green("投稿が更新されました！")
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
			FeaturedMedia: article.FrontMatter.FeaturedMedia,
		}

		color.Cyan("投稿を作成中...")

		post, err = client.CreatePost(req)
		if err != nil {
			color.Red("投稿の作成に失敗: %v", err)
			os.Exit(1)
		}

		color.Green("投稿が作成されました！")
	}

	color.White("  ID: %d", post.ID)
	color.White("  URL: %s", post.Link)
	color.White("  ステータス: %s", formatStatus(post.Status))

	// 公開時にdrafts/からposts/に移動
	if status == "publish" && strings.Contains(filePath, "drafts/") {
		moveToPublished(filePath, post)
	}
}

// moveToPublished はdrafts/からposts/に記事を移動する
func moveToPublished(filePath string, post *types.Post) {
	// ディレクトリパスを取得
	srcDir := filepath.Dir(filePath)

	// 新しいディレクトリ名: YYYY-MM-DD_slug
	newDirName := post.Date.Time.Format("2006-01-02") + "_" + post.Slug
	destDir := filepath.Join("posts", newDirName)

	// 移動先が既に存在する場合はスキップ
	if _, err := os.Stat(destDir); err == nil {
		color.Yellow("  移動先が既に存在します: %s", destDir)
		return
	}

	// ディレクトリを移動
	if err := os.Rename(srcDir, destDir); err != nil {
		color.Yellow("  記事の移動に失敗: %v", err)
		return
	}

	// フロントマターを更新（IDとstatusを反映）
	articlePath := filepath.Join(destDir, filepath.Base(filePath))
	article, err := converter.ParseArticle(articlePath)
	if err != nil {
		color.Yellow("  フロントマター更新に失敗: %v", err)
		return
	}

	article.FrontMatter.ID = post.ID
	article.FrontMatter.Status = "publish"
	article.FrontMatter.Date = post.Date.Time.Format("2006-01-02T15:04:05")

	content, err := converter.GenerateArticleFile(article)
	if err != nil {
		color.Yellow("  記事ファイル生成に失敗: %v", err)
		return
	}

	if err := os.WriteFile(articlePath, []byte(content), 0644); err != nil {
		color.Yellow("  記事ファイル保存に失敗: %v", err)
		return
	}

	color.Green("  記事を移動しました: %s", destDir)
	fmt.Println()
}
