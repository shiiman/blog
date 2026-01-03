package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/fatih/color"
	"github.com/shiimanblog/wp-cli/internal/config"
	"github.com/shiimanblog/wp-cli/internal/converter"
	"github.com/shiimanblog/wp-cli/internal/types"
	"github.com/shiimanblog/wp-cli/internal/wp"
	"github.com/spf13/cobra"
)

var importCmd = &cobra.Command{
	Use:   "import [posts|pages|post|page] [id]",
	Short: "WordPressから記事をインポート",
	Long: `WordPressから記事または固定ページをインポートし、ローカルにMarkdownファイルとして保存します。

例:
  wp-cli import posts           # 全投稿をインポート
  wp-cli import posts --limit=10  # 最新10件をインポート
  wp-cli import post 123        # ID=123の投稿をインポート
  wp-cli import pages           # 全固定ページをインポート
  wp-cli import page 45         # ID=45の固定ページをインポート`,
	Args: cobra.MinimumNArgs(1),
	Run:  runImport,
}

var importLimit int
var importOutputDir string

func init() {
	rootCmd.AddCommand(importCmd)
	importCmd.Flags().IntVar(&importLimit, "limit", 100, "インポートする最大件数")
	importCmd.Flags().StringVarP(&importOutputDir, "output", "o", "", "出力ディレクトリ（デフォルト: posts/ または pages/）")
}

func runImport(cmd *cobra.Command, args []string) {
	cfg, err := config.Load()
	if err != nil {
		color.Red("設定エラー: %v", err)
		os.Exit(1)
	}

	client := wp.NewClient(cfg)
	itemType := args[0]

	switch itemType {
	case "posts":
		importPosts(client)
	case "post":
		if len(args) < 2 {
			color.Red("投稿IDを指定してください。")
			os.Exit(1)
		}
		id, err := strconv.Atoi(args[1])
		if err != nil {
			color.Red("無効な投稿ID: %s", args[1])
			os.Exit(1)
		}
		importPost(client, id)
	case "pages":
		importPages(client)
	case "page":
		if len(args) < 2 {
			color.Red("固定ページIDを指定してください。")
			os.Exit(1)
		}
		id, err := strconv.Atoi(args[1])
		if err != nil {
			color.Red("無効な固定ページID: %s", args[1])
			os.Exit(1)
		}
		importPage(client, id)
	default:
		color.Red("無効な引数です。'posts', 'post', 'pages', または 'page' を指定してください。")
		os.Exit(1)
	}
}

func importPosts(client *wp.Client) {
	color.Cyan("投稿をインポート中...")

	posts, err := client.GetPosts(1, importLimit, "")
	if err != nil {
		color.Red("投稿一覧の取得に失敗: %v", err)
		os.Exit(1)
	}

	if len(posts) == 0 {
		color.Yellow("インポートする投稿がありません。")
		return
	}

	outputDir := importOutputDir
	if outputDir == "" {
		outputDir = "posts"
	}

	imported := 0
	for _, post := range posts {
		// 投稿を詳細取得（content.rawを取得するため）
		fullPost, err := client.GetPost(post.ID)
		if err != nil {
			color.Yellow("投稿 %d の詳細取得に失敗: %v", post.ID, err)
			continue
		}

		article, err := converter.PostToArticle(fullPost)
		if err != nil {
			color.Yellow("投稿 %d の変換に失敗: %v", post.ID, err)
			continue
		}

		// ディレクトリ名: YYYY-MM-DD_slug
		dirName := fullPost.Date.Format("2006-01-02") + "_" + fullPost.Slug
		dirPath := filepath.Join(outputDir, dirName)

		if err := saveArticle(dirPath, "article.md", article); err != nil {
			color.Yellow("投稿 %d の保存に失敗: %v", post.ID, err)
			continue
		}

		color.White("  ✓ %s", fullPost.Title.Rendered)
		imported++
	}

	color.Green("\n%d件の投稿をインポートしました。", imported)
}

func importPost(client *wp.Client, id int) {
	color.Cyan("投稿 %d をインポート中...", id)

	post, err := client.GetPost(id)
	if err != nil {
		color.Red("投稿の取得に失敗: %v", err)
		os.Exit(1)
	}

	article, err := converter.PostToArticle(post)
	if err != nil {
		color.Red("投稿の変換に失敗: %v", err)
		os.Exit(1)
	}

	outputDir := importOutputDir
	if outputDir == "" {
		outputDir = "posts"
	}

	// ディレクトリ名: YYYY-MM-DD_slug
	dirName := post.Date.Format("2006-01-02") + "_" + post.Slug
	dirPath := filepath.Join(outputDir, dirName)

	if err := saveArticle(dirPath, "article.md", article); err != nil {
		color.Red("投稿の保存に失敗: %v", err)
		os.Exit(1)
	}

	color.Green("投稿をインポートしました: %s", filepath.Join(dirPath, "article.md"))
}

func importPages(client *wp.Client) {
	color.Cyan("固定ページをインポート中...")

	pages, err := client.GetPages(1, importLimit, "")
	if err != nil {
		color.Red("固定ページ一覧の取得に失敗: %v", err)
		os.Exit(1)
	}

	if len(pages) == 0 {
		color.Yellow("インポートする固定ページがありません。")
		return
	}

	outputDir := importOutputDir
	if outputDir == "" {
		outputDir = "pages"
	}

	imported := 0
	for _, page := range pages {
		// 固定ページを詳細取得
		fullPage, err := client.GetPage(page.ID)
		if err != nil {
			color.Yellow("固定ページ %d の詳細取得に失敗: %v", page.ID, err)
			continue
		}

		article, err := converter.PageToArticle(fullPage)
		if err != nil {
			color.Yellow("固定ページ %d の変換に失敗: %v", page.ID, err)
			continue
		}

		// ディレクトリ名: slug
		dirPath := filepath.Join(outputDir, fullPage.Slug)

		if err := saveArticle(dirPath, "page.md", article); err != nil {
			color.Yellow("固定ページ %d の保存に失敗: %v", page.ID, err)
			continue
		}

		color.White("  ✓ %s", fullPage.Title.Rendered)
		imported++
	}

	color.Green("\n%d件の固定ページをインポートしました。", imported)
}

func importPage(client *wp.Client, id int) {
	color.Cyan("固定ページ %d をインポート中...", id)

	page, err := client.GetPage(id)
	if err != nil {
		color.Red("固定ページの取得に失敗: %v", err)
		os.Exit(1)
	}

	article, err := converter.PageToArticle(page)
	if err != nil {
		color.Red("固定ページの変換に失敗: %v", err)
		os.Exit(1)
	}

	outputDir := importOutputDir
	if outputDir == "" {
		outputDir = "pages"
	}

	// ディレクトリ名: slug
	dirPath := filepath.Join(outputDir, page.Slug)

	if err := saveArticle(dirPath, "page.md", article); err != nil {
		color.Red("固定ページの保存に失敗: %v", err)
		os.Exit(1)
	}

	color.Green("固定ページをインポートしました: %s", filepath.Join(dirPath, "page.md"))
}

func saveArticle(dirPath, filename string, article *types.Article) error {
	// ディレクトリ作成
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return fmt.Errorf("ディレクトリの作成に失敗: %w", err)
	}

	// assetsディレクトリも作成
	assetsPath := filepath.Join(dirPath, "assets")
	if err := os.MkdirAll(assetsPath, 0755); err != nil {
		return fmt.Errorf("assetsディレクトリの作成に失敗: %w", err)
	}

	// 記事ファイル生成
	content, err := converter.GenerateArticleFile(article)
	if err != nil {
		return err
	}

	filePath := filepath.Join(dirPath, filename)
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("ファイルの書き込みに失敗: %w", err)
	}

	return nil
}
