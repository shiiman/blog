package cmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/fatih/color"
	"github.com/shiimanblog/wp-cli/internal/wp"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list [posts|pages]",
	Short: "投稿または固定ページの一覧を表示",
	Long: `投稿または固定ページの一覧を表示します。

例:
  wp-cli list posts           # 投稿一覧
  wp-cli list posts --status=draft  # 下書きのみ
  wp-cli list pages           # 固定ページ一覧`,
	Args: cobra.ExactArgs(1),
	RunE: runList,
}

var listStatus string
var listLimit int

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringVar(&listStatus, "status", "", "フィルタするステータス (draft, publish, pending, private)")
	listCmd.Flags().IntVar(&listLimit, "limit", 20, "取得件数")
}

func runList(cmd *cobra.Command, args []string) error {
	normalizedStatus, err := normalizeAndValidateStatus(listStatus, true)
	if err != nil {
		return fmt.Errorf("list --status が不正です: %w", err)
	}
	listStatus = normalizedStatus

	client, err := setupClient()
	if err != nil {
		return err
	}

	ctx := cmd.Context()
	itemType := args[0]

	switch itemType {
	case "posts":
		return listPosts(ctx, client)
	case "pages":
		return listPages(ctx, client)
	default:
		return fmt.Errorf("無効な引数です。'posts' または 'pages' を指定してください")
	}
}

func listPosts(ctx context.Context, client *wp.Client) error {
	color.Cyan("投稿一覧を取得中...")

	posts, err := client.GetPosts(ctx, 1, listLimit, listStatus)
	if err != nil {
		return fmt.Errorf("投稿一覧の取得に失敗: %w", err)
	}

	if len(posts) == 0 {
		color.Yellow("投稿が見つかりませんでした。")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	_, _ = fmt.Fprintln(w, "ID\tタイトル\tステータス\tスラッグ\t日付")
	_, _ = fmt.Fprintln(w, "---\t---\t---\t---\t---")

	for _, post := range posts {
		status := formatStatus(post.Status)
		date := post.Date.Format("2006-01-02")
		title := truncate(post.Title.Rendered, 40)
		_, _ = fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\n", post.ID, title, status, post.Slug, date)
	}
	_ = w.Flush()

	color.Green("\n%d件の投稿を表示しました。", len(posts))
	return nil
}

func listPages(ctx context.Context, client *wp.Client) error {
	color.Cyan("固定ページ一覧を取得中...")

	pages, err := client.GetPages(ctx, 1, listLimit, listStatus)
	if err != nil {
		return fmt.Errorf("固定ページ一覧の取得に失敗: %w", err)
	}

	if len(pages) == 0 {
		color.Yellow("固定ページが見つかりませんでした。")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	_, _ = fmt.Fprintln(w, "ID\tタイトル\tステータス\tスラッグ\t親ID")
	_, _ = fmt.Fprintln(w, "---\t---\t---\t---\t---")

	for _, page := range pages {
		status := formatStatus(page.Status)
		title := truncate(page.Title.Rendered, 40)
		_, _ = fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%d\n", page.ID, title, status, page.Slug, page.Parent)
	}
	_ = w.Flush()

	color.Green("\n%d件の固定ページを表示しました。", len(pages))
	return nil
}
