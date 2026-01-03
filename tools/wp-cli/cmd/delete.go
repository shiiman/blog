package cmd

import (
	"os"
	"strconv"

	"github.com/fatih/color"
	"github.com/shiimanblog/wp-cli/internal/config"
	"github.com/shiimanblog/wp-cli/internal/wp"
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
	Run:  runDelete,
}

var deleteForce bool

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().BoolVar(&deleteForce, "force", false, "完全に削除（ゴミ箱をスキップ）")
}

func runDelete(cmd *cobra.Command, args []string) {
	id, err := strconv.Atoi(args[0])
	if err != nil {
		color.Red("無効な投稿ID: %s", args[0])
		os.Exit(1)
	}

	cfg, err := config.Load()
	if err != nil {
		color.Red("設定エラー: %v", err)
		os.Exit(1)
	}

	client := wp.NewClient(cfg)

	if deleteForce {
		color.Cyan("投稿 %d を完全に削除中...", id)
	} else {
		color.Cyan("投稿 %d をゴミ箱に移動中...", id)
	}

	if err := client.DeletePost(id, deleteForce); err != nil {
		color.Red("投稿の削除に失敗: %v", err)
		os.Exit(1)
	}

	if deleteForce {
		color.Green("投稿 %d を完全に削除しました。", id)
	} else {
		color.Green("投稿 %d をゴミ箱に移動しました。", id)
	}
}
