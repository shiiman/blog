package cmd

import (
	"fmt"
	"strconv"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "投稿または固定ページを削除",
	Long: `指定されたIDの投稿または固定ページを削除します。
デフォルトではゴミ箱に移動します。--forceで完全削除します。
--page フラグを指定すると固定ページを削除します。

例:
  wp-cli delete 123
  wp-cli delete 123 --force
  wp-cli delete 456 --page
  wp-cli delete 456 --page --force`,
	Args: cobra.ExactArgs(1),
	RunE: runDelete,
}

var deleteForce bool
var deletePage bool

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().BoolVar(&deleteForce, "force", false, "完全に削除（ゴミ箱をスキップ）")
	deleteCmd.Flags().BoolVar(&deletePage, "page", false, "固定ページを削除")
}

func runDelete(cmd *cobra.Command, args []string) error {
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("無効なID: %s", args[0])
	}
	if id <= 0 {
		return fmt.Errorf("IDは正の整数を指定してください: %d", id)
	}

	client, err := setupClient()
	if err != nil {
		return err
	}

	ctx := cmd.Context()

	// 対象の種別を決定
	target := "投稿"
	if deletePage {
		target = "固定ページ"
	}

	if deleteForce {
		color.Cyan("%s %d を完全に削除中...", target, id)
	} else {
		color.Cyan("%s %d をゴミ箱に移動中...", target, id)
	}

	if deletePage {
		err = client.DeletePage(ctx, id, deleteForce)
	} else {
		err = client.DeletePost(ctx, id, deleteForce)
	}
	if err != nil {
		return fmt.Errorf("%sの削除に失敗: %w", target, err)
	}

	if deleteForce {
		color.Green("%s %d を完全に削除しました。", target, id)
	} else {
		color.Green("%s %d をゴミ箱に移動しました。", target, id)
	}
	return nil
}
