package cmd

import (
	"fmt"
	"strconv"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "投稿を削除",
	Long: `指定されたIDの投稿を削除します。
デフォルトではゴミ箱に移動します。--forceで完全削除します。

例:
  wp-cli delete 123
  wp-cli delete 123 --force`,
	Args: cobra.ExactArgs(1),
	RunE: runDelete,
}

var deleteForce bool

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().BoolVar(&deleteForce, "force", false, "完全に削除（ゴミ箱をスキップ）")
}

func runDelete(cmd *cobra.Command, args []string) error {
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("無効な投稿ID: %s", args[0])
	}

	client, err := setupClient()
	if err != nil {
		return err
	}

	ctx := cmd.Context()

	if deleteForce {
		color.Cyan("投稿 %d を完全に削除中...", id)
	} else {
		color.Cyan("投稿 %d をゴミ箱に移動中...", id)
	}

	if err := client.DeletePost(ctx, id, deleteForce); err != nil {
		return fmt.Errorf("投稿の削除に失敗: %w", err)
	}

	if deleteForce {
		color.Green("投稿 %d を完全に削除しました。", id)
	} else {
		color.Green("投稿 %d をゴミ箱に移動しました。", id)
	}
	return nil
}
