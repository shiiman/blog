package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/fatih/color"
	"github.com/shiimanblog/wp-cli/internal/config"
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
	Run:  runList,
}

var listStatus string
var listLimit int

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringVar(&listStatus, "status", "", "フィルタするステータス (draft, publish, pending, private)")
	listCmd.Flags().IntVar(&listLimit, "limit", 20, "取得件数")
}

func runList(cmd *cobra.Command, args []string) {
	cfg, err := config.Load()
	if err != nil {
		color.Red("設定エラー: %v", err)
		os.Exit(1)
	}

	client := wp.NewClient(cfg)
	itemType := args[0]

	switch itemType {
	case "posts":
		listPosts(client)
	case "pages":
		listPages(client)
	default:
		color.Red("無効な引数です。'posts' または 'pages' を指定してください。")
		os.Exit(1)
	}
}

func listPosts(client *wp.Client) {
	color.Cyan("投稿一覧を取得中...")

	posts, err := client.GetPosts(1, listLimit, listStatus)
	if err != nil {
		color.Red("投稿一覧の取得に失敗: %v", err)
		os.Exit(1)
	}

	if len(posts) == 0 {
		color.Yellow("投稿が見つかりませんでした。")
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tタイトル\tステータス\tスラッグ\t日付")
	fmt.Fprintln(w, "---\t---\t---\t---\t---")

	for _, post := range posts {
		status := formatStatus(post.Status)
		date := post.Date.Format("2006-01-02")
		title := truncate(post.Title.Rendered, 40)
		fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\n", post.ID, title, status, post.Slug, date)
	}
	w.Flush()

	color.Green("\n%d件の投稿を表示しました。", len(posts))
}

func listPages(client *wp.Client) {
	color.Cyan("固定ページ一覧を取得中...")

	pages, err := client.GetPages(1, listLimit, listStatus)
	if err != nil {
		color.Red("固定ページ一覧の取得に失敗: %v", err)
		os.Exit(1)
	}

	if len(pages) == 0 {
		color.Yellow("固定ページが見つかりませんでした。")
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tタイトル\tステータス\tスラッグ\t親ID")
	fmt.Fprintln(w, "---\t---\t---\t---\t---")

	for _, page := range pages {
		status := formatStatus(page.Status)
		title := truncate(page.Title.Rendered, 40)
		fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%d\n", page.ID, title, status, page.Slug, page.Parent)
	}
	w.Flush()

	color.Green("\n%d件の固定ページを表示しました。", len(pages))
}

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

func truncate(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen-3]) + "..."
}
