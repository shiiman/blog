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

var categoriesCmd = &cobra.Command{
	Use:   "categories",
	Short: "カテゴリ一覧を表示",
	Long: `WordPressのカテゴリ一覧を表示します。

例:
  wp-cli categories`,
	Run: runCategories,
}

var tagsCmd = &cobra.Command{
	Use:   "tags",
	Short: "タグ一覧を表示",
	Long: `WordPressのタグ一覧を表示します。

例:
  wp-cli tags`,
	Run: runTags,
}

func init() {
	rootCmd.AddCommand(categoriesCmd)
	rootCmd.AddCommand(tagsCmd)
}

func runCategories(cmd *cobra.Command, args []string) {
	cfg, err := config.Load()
	if err != nil {
		color.Red("設定エラー: %v", err)
		os.Exit(1)
	}

	client := wp.NewClient(cfg)

	color.Cyan("カテゴリ一覧を取得中...")

	categories, err := client.GetCategories()
	if err != nil {
		color.Red("カテゴリ一覧の取得に失敗: %v", err)
		os.Exit(1)
	}

	if len(categories) == 0 {
		color.Yellow("カテゴリが見つかりませんでした。")
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\t名前\tスラッグ\t記事数\t親ID")
	fmt.Fprintln(w, "---\t---\t---\t---\t---")

	for _, cat := range categories {
		fmt.Fprintf(w, "%d\t%s\t%s\t%d\t%d\n", cat.ID, cat.Name, cat.Slug, cat.Count, cat.Parent)
	}
	w.Flush()

	color.Green("\n%d件のカテゴリを表示しました。", len(categories))
}

func runTags(cmd *cobra.Command, args []string) {
	cfg, err := config.Load()
	if err != nil {
		color.Red("設定エラー: %v", err)
		os.Exit(1)
	}

	client := wp.NewClient(cfg)

	color.Cyan("タグ一覧を取得中...")

	tags, err := client.GetTags()
	if err != nil {
		color.Red("タグ一覧の取得に失敗: %v", err)
		os.Exit(1)
	}

	if len(tags) == 0 {
		color.Yellow("タグが見つかりませんでした。")
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\t名前\tスラッグ\t記事数")
	fmt.Fprintln(w, "---\t---\t---\t---")

	for _, tag := range tags {
		fmt.Fprintf(w, "%d\t%s\t%s\t%d\n", tag.ID, tag.Name, tag.Slug, tag.Count)
	}
	w.Flush()

	color.Green("\n%d件のタグを表示しました。", len(tags))
}
