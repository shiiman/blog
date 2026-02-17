package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/fatih/color"
	"github.com/shiimanblog/wp-cli/internal/converter"
	"github.com/shiimanblog/wp-cli/internal/types"
	"github.com/shiimanblog/wp-cli/internal/wp"
	"github.com/spf13/cobra"
)

// importItem はバッチインポートの入力アイテムを表す
type importItem struct {
	id    int
	title string
}

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
	RunE: runImport,
}

var importLimit int
var importOutputDir string

func init() {
	rootCmd.AddCommand(importCmd)
	importCmd.Flags().IntVar(&importLimit, "limit", 100, "インポートする最大件数")
	importCmd.Flags().StringVarP(&importOutputDir, "output", "o", "", "出力ディレクトリ（デフォルト: posts/ または pages/）")
}

func runImport(cmd *cobra.Command, args []string) error {
	client, err := setupClient()
	if err != nil {
		return err
	}

	ctx := cmd.Context()
	itemType := args[0]

	switch itemType {
	case "posts":
		return importPosts(ctx, client)
	case "post":
		if len(args) < 2 {
			return fmt.Errorf("投稿IDを指定してください")
		}
		id, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("無効な投稿ID: %s", args[1])
		}
		if id <= 0 {
			return fmt.Errorf("投稿IDは正の整数を指定してください: %d", id)
		}
		return importPost(ctx, client, id)
	case "pages":
		return importPages(ctx, client)
	case "page":
		if len(args) < 2 {
			return fmt.Errorf("固定ページIDを指定してください")
		}
		id, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("無効な固定ページID: %s", args[1])
		}
		if id <= 0 {
			return fmt.Errorf("固定ページIDは正の整数を指定してください: %d", id)
		}
		return importPage(ctx, client, id)
	default:
		return fmt.Errorf("無効な引数です。'posts', 'post', 'pages', または 'page' を指定してください")
	}
}

// importResult は並行インポートの結果を保持する
type importResult struct {
	id      int
	article *types.Article
	dirPath string
	title   string
	err     error
}

const maxConcurrency = 10

// batchImport は複数アイテムを並行でインポートする共通関数
func batchImport(
	ctx context.Context,
	typeName string,
	items []importItem,
	fileName string,
	processItem func(ctx context.Context, id int) (article *types.Article, dirPath string, err error),
) error {
	if len(items) == 0 {
		color.Yellow("インポートする%sがありません。", typeName)
		return nil
	}

	results := make([]importResult, len(items))
	sem := make(chan struct{}, maxConcurrency)
	var wg sync.WaitGroup

	for i, item := range items {
		wg.Add(1)
		go func(idx int, it importItem) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			article, dirPath, err := processItem(ctx, it.id)
			if err != nil {
				results[idx] = importResult{id: it.id, err: err}
				return
			}
			results[idx] = importResult{
				id:      it.id,
				article: article,
				dirPath: dirPath,
				title:   it.title,
			}
		}(i, item)
	}
	wg.Wait()

	imported := 0
	for _, r := range results {
		if r.err != nil {
			color.Yellow("%s %d: %v", typeName, r.id, r.err)
			continue
		}
		if err := saveArticle(r.dirPath, fileName, r.article); err != nil {
			color.Yellow("%s %d の保存に失敗: %v", typeName, r.id, err)
			continue
		}
		color.White("  ✓ %s", r.title)
		imported++
	}

	color.Green("\n%d件の%sをインポートしました。", imported, typeName)
	return nil
}

func importPosts(ctx context.Context, client *wp.Client) error {
	color.Cyan("投稿をインポート中...")

	posts, err := client.GetPosts(ctx, 1, importLimit, "")
	if err != nil {
		return fmt.Errorf("投稿一覧の取得に失敗: %w", err)
	}

	outputDir := importOutputDir
	if outputDir == "" {
		outputDir = "posts"
	}

	items := make([]importItem, len(posts))
	for i, p := range posts {
		items[i] = importItem{id: p.ID, title: p.Title.Rendered}
	}

	return batchImport(ctx, "投稿", items, "article.md", func(ctx context.Context, id int) (*types.Article, string, error) {
		fullPost, err := client.GetPost(ctx, id)
		if err != nil {
			return nil, "", fmt.Errorf("詳細取得に失敗: %w", err)
		}
		article, err := converter.PostToArticle(fullPost)
		if err != nil {
			return nil, "", fmt.Errorf("変換に失敗: %w", err)
		}
		dirName := fullPost.Date.Format("2006-01-02") + "_" + sanitizeSlug(fullPost.Slug)
		return article, filepath.Join(outputDir, dirName), nil
	})
}

func importPost(ctx context.Context, client *wp.Client, id int) error {
	color.Cyan("投稿 %d をインポート中...", id)

	post, err := client.GetPost(ctx, id)
	if err != nil {
		return fmt.Errorf("投稿の取得に失敗: %w", err)
	}

	article, err := converter.PostToArticle(post)
	if err != nil {
		return fmt.Errorf("投稿の変換に失敗: %w", err)
	}

	outputDir := importOutputDir
	if outputDir == "" {
		outputDir = "posts"
	}

	// ディレクトリ名: YYYY-MM-DD_slug
	dirName := post.Date.Format("2006-01-02") + "_" + sanitizeSlug(post.Slug)
	dirPath := filepath.Join(outputDir, dirName)

	if err := saveArticle(dirPath, "article.md", article); err != nil {
		return fmt.Errorf("投稿の保存に失敗: %w", err)
	}

	color.Green("投稿をインポートしました: %s", filepath.Join(dirPath, "article.md"))
	return nil
}

func importPages(ctx context.Context, client *wp.Client) error {
	color.Cyan("固定ページをインポート中...")

	pages, err := client.GetPages(ctx, 1, importLimit, "")
	if err != nil {
		return fmt.Errorf("固定ページ一覧の取得に失敗: %w", err)
	}

	outputDir := importOutputDir
	if outputDir == "" {
		outputDir = "pages"
	}

	items := make([]importItem, len(pages))
	for i, p := range pages {
		items[i] = importItem{id: p.ID, title: p.Title.Rendered}
	}

	return batchImport(ctx, "固定ページ", items, "page.md", func(ctx context.Context, id int) (*types.Article, string, error) {
		fullPage, err := client.GetPage(ctx, id)
		if err != nil {
			return nil, "", fmt.Errorf("詳細取得に失敗: %w", err)
		}
		article, err := converter.PageToArticle(fullPage)
		if err != nil {
			return nil, "", fmt.Errorf("変換に失敗: %w", err)
		}
		return article, filepath.Join(outputDir, sanitizeSlug(fullPage.Slug)), nil
	})
}

func importPage(ctx context.Context, client *wp.Client, id int) error {
	color.Cyan("固定ページ %d をインポート中...", id)

	page, err := client.GetPage(ctx, id)
	if err != nil {
		return fmt.Errorf("固定ページの取得に失敗: %w", err)
	}

	article, err := converter.PageToArticle(page)
	if err != nil {
		return fmt.Errorf("固定ページの変換に失敗: %w", err)
	}

	outputDir := importOutputDir
	if outputDir == "" {
		outputDir = "pages"
	}

	// ディレクトリ名: slug
	dirPath := filepath.Join(outputDir, sanitizeSlug(page.Slug))

	if err := saveArticle(dirPath, "page.md", article); err != nil {
		return fmt.Errorf("固定ページの保存に失敗: %w", err)
	}

	color.Green("固定ページをインポートしました: %s", filepath.Join(dirPath, "page.md"))
	return nil
}

func saveArticle(dirPath, filename string, article *types.Article) error {
	// パストラバーサル対策: パスを正規化し、ベースディレクトリ外への書き込みを防止
	cleanPath := filepath.Clean(dirPath)
	absPath, err := filepath.Abs(cleanPath)
	if err != nil {
		return fmt.Errorf("パスの解決に失敗: %w", err)
	}

	// カレントディレクトリを基準にベースディレクトリ外への書き込みを防止
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("カレントディレクトリの取得に失敗: %w", err)
	}
	if !strings.HasPrefix(absPath, cwd+string(filepath.Separator)) && absPath != cwd {
		return fmt.Errorf("不正なパスです（ベースディレクトリ外への書き込み）: %s", dirPath)
	}

	// ディレクトリ作成
	if err := os.MkdirAll(cleanPath, 0755); err != nil {
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
