package cmd

import (
	"fmt"
	"path/filepath"

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

フロントマターに id フィールドがある場合は既存記事を更新し、
id がない場合は新規記事として作成します。
デフォルトでは公開状態で投稿されます。

アイキャッチ画像（eyecatch.png）が記事ディレクトリに存在する場合、
自動的にWordPressへアップロードされます。

例:
  wp-cli post drafts/2025-01-03_my-article/article.md        # 新規作成（公開）
  wp-cli post drafts/2025-01-03_my-article/article.md --draft # 新規作成（下書き）
  wp-cli post drafts/2025-01-03_my-article/article.md --dry-run # 投稿せずに確認

フロントマター例（新規作成）:
  ---
  title: "記事タイトル"
  slug: "article-slug"
  categories: [1]
  tags: [10, 20]
  ---

フロントマター例（既存記事の更新）:
  ---
  id: 123
  title: "記事タイトル"
  slug: "article-slug"
  ---`,
	Args: cobra.ExactArgs(1),
	RunE: runPost,
}

var postDraft bool
var postDryRun bool

func init() {
	rootCmd.AddCommand(postCmd)
	postCmd.Flags().BoolVar(&postDraft, "draft", false, "下書き状態で投稿")
	postCmd.Flags().BoolVar(&postDryRun, "dry-run", false, "投稿せずに内容を確認")
}

func runPost(cmd *cobra.Command, args []string) error {
	filePath := args[0]

	// 記事ファイルを解析
	article, err := converter.ParseArticle(filePath)
	if err != nil {
		return fmt.Errorf("記事ファイルの解析に失敗: %w", err)
	}

	// ステータスの決定（postコマンドはデフォルト公開）
	status := determinePostStatus(postDraft)

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

	ctx := cmd.Context()

	// アイキャッチ画像のアップロード
	articleDir := filepath.Dir(filePath)
	featuredMediaID, err := uploadEyecatchIfExists(ctx, client, articleDir, article.FrontMatter.FeaturedMedia, false)
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

		post, err = client.UpdatePost(ctx, article.FrontMatter.ID, req)
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

		post, err = client.CreatePost(ctx, req)
		if err != nil {
			return fmt.Errorf("投稿の作成に失敗: %w", err)
		}
	}

	printPostResult(post, article.FrontMatter.ID > 0)

	// WordPress結果をローカル記事に同期し、公開時はdrafts->postsへ移動
	syncedPath, err := syncPostToLocal(filePath, post, true)
	if err != nil {
		return err
	}

	color.Green("  記事ファイルを同期しました: %s", syncedPath)
	return nil
}
