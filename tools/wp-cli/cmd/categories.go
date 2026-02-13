package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/fatih/color"
	"github.com/shiimanblog/wp-cli/internal/types"
	"github.com/spf13/cobra"
)

var categoriesCmd = &cobra.Command{
	Use:   "categories",
	Short: "カテゴリ一覧を表示",
	Long: `WordPressのカテゴリ一覧を表示します。

例:
  wp-cli categories`,
	RunE: runCategories,
}

var createCategoryCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "新しいカテゴリを作成",
	Args:  cobra.ExactArgs(1),
	RunE:  runCreateCategory,
}

var updateCategoryCmd = &cobra.Command{
	Use:   "update [id]",
	Short: "既存のカテゴリを更新",
	Args:  cobra.ExactArgs(1),
	RunE:  runUpdateCategory,
}

var tagsCmd = &cobra.Command{
	Use:   "tags",
	Short: "タグ一覧を表示",
	Long: `WordPressのタグ一覧を表示します。

例:
  wp-cli tags`,
	RunE: runTags,
}

var parentID int
var categoryName string

func init() {
	categoriesCmd.AddCommand(createCategoryCmd)
	createCategoryCmd.Flags().IntVarP(&parentID, "parent", "p", 0, "親カテゴリのID")

	categoriesCmd.AddCommand(updateCategoryCmd)
	updateCategoryCmd.Flags().IntVarP(&parentID, "parent", "p", 0, "親カテゴリのID")
	updateCategoryCmd.Flags().StringVarP(&categoryName, "name", "n", "", "カテゴリ名")

	rootCmd.AddCommand(categoriesCmd)
	rootCmd.AddCommand(tagsCmd)
}

func runCategories(cmd *cobra.Command, args []string) error {
	client, err := setupClient()
	if err != nil {
		return err
	}

	ctx := cmd.Context()

	color.Cyan("カテゴリ一覧を取得中...")

	categories, err := client.GetCategories(ctx)
	if err != nil {
		return fmt.Errorf("カテゴリ一覧の取得に失敗: %w", err)
	}

	if len(categories) == 0 {
		color.Yellow("カテゴリが見つかりませんでした。")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	_, _ = fmt.Fprintln(w, "ID\t名前\tスラッグ\t記事数\t親ID")
	_, _ = fmt.Fprintln(w, "---\t---\t---\t---\t---")

	for _, cat := range categories {
		_, _ = fmt.Fprintf(w, "%d\t%s\t%s\t%d\t%d\n", cat.ID, cat.Name, cat.Slug, cat.Count, cat.Parent)
	}
	_ = w.Flush()

	color.Green("\n%d件のカテゴリを表示しました。", len(categories))
	return nil
}

func runCreateCategory(cmd *cobra.Command, args []string) error {
	name := args[0]
	client, err := setupClient()
	if err != nil {
		return err
	}

	ctx := cmd.Context()

	color.Cyan("カテゴリ '%s' を作成中...", name)

	req := &types.CreateCategoryRequest{
		Name:   name,
		Parent: parentID,
	}

	category, err := client.CreateCategory(ctx, req)
	if err != nil {
		return fmt.Errorf("カテゴリの作成に失敗: %w", err)
	}

	color.Green("カテゴリが作成されました！")
	fmt.Printf("  ID: %d\n", category.ID)
	fmt.Printf("  名前: %s\n", category.Name)
	fmt.Printf("  スラッグ: %s\n", category.Slug)
	return nil
}

func runUpdateCategory(cmd *cobra.Command, args []string) error {
	idStr := args[0]
	var id int
	_, err := fmt.Sscanf(idStr, "%d", &id)
	if err != nil {
		return fmt.Errorf("ID形式エラー: %w", err)
	}

	client, err := setupClient()
	if err != nil {
		return err
	}

	ctx := cmd.Context()

	color.Cyan("カテゴリID %d を更新中...", id)

	req := &types.UpdateCategoryRequest{
		Name:   categoryName,
		Parent: parentID,
	}

	category, err := client.UpdateCategory(ctx, id, req)
	if err != nil {
		return fmt.Errorf("カテゴリの更新に失敗: %w", err)
	}

	color.Green("カテゴリが更新されました！")
	fmt.Printf("  ID: %d\n", category.ID)
	fmt.Printf("  名前: %s\n", category.Name)
	fmt.Printf("  スラッグ: %s\n", category.Slug)
	fmt.Printf("  親ID: %d\n", category.Parent)
	return nil
}

func runTags(cmd *cobra.Command, args []string) error {
	client, err := setupClient()
	if err != nil {
		return err
	}

	ctx := cmd.Context()

	color.Cyan("タグ一覧を取得中...")

	tags, err := client.GetTags(ctx)
	if err != nil {
		return fmt.Errorf("タグ一覧の取得に失敗: %w", err)
	}

	if len(tags) == 0 {
		color.Yellow("タグが見つかりませんでした。")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	_, _ = fmt.Fprintln(w, "ID\t名前\tスラッグ\t記事数")
	_, _ = fmt.Fprintln(w, "---\t---\t---\t---")

	for _, tag := range tags {
		_, _ = fmt.Fprintf(w, "%d\t%s\t%s\t%d\n", tag.ID, tag.Name, tag.Slug, tag.Count)
	}
	_ = w.Flush()

	color.Green("\n%d件のタグを表示しました。", len(tags))
	return nil
}
