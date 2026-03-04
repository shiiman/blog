package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var mediaCmd = &cobra.Command{
	Use:   "media",
	Short: "メディア管理",
	Long:  "WordPressメディアライブラリの管理コマンド",
}

var mediaUploadCmd = &cobra.Command{
	Use:   "upload <file> [file...]",
	Short: "メディアをWordPressにアップロード",
	Long: `ローカルファイルをWordPressメディアライブラリにアップロードします。
複数ファイルを同時にアップロード可能です。

例:
  wp-cli media upload image.png
  wp-cli media upload video.mp4 audio.mp3
  wp-cli media upload assets/*`,
	Args: cobra.MinimumNArgs(1),
	RunE: runMediaUpload,
}

const maxMediaSize = 100 * 1024 * 1024 // 100MB

func init() {
	mediaCmd.AddCommand(mediaUploadCmd)
	rootCmd.AddCommand(mediaCmd)
}

func runMediaUpload(cmd *cobra.Command, args []string) error {
	client, err := setupClient()
	if err != nil {
		return err
	}

	ctx := cmd.Context()

	for _, filePath := range args {
		// ファイル存在確認
		info, err := os.Stat(filePath)
		if err != nil {
			color.Red("ファイルが見つかりません: %s", filePath)
			continue
		}
		if info.IsDir() {
			color.Red("ディレクトリはスキップします: %s", filePath)
			continue
		}

		// サイズチェック
		if info.Size() > maxMediaSize {
			color.Red("ファイルサイズが上限（100MB）を超えています: %s (%d bytes)", filePath, info.Size())
			continue
		}

		// ファイル読み込み
		data, err := os.ReadFile(filePath)
		if err != nil {
			color.Red("ファイル読み込みに失敗: %s: %v", filePath, err)
			continue
		}

		filename := filepath.Base(filePath)
		mimeType := detectMediaMIMEType(filename)

		color.Cyan("アップロード中: %s (%s)", filename, mimeType)

		media, err := client.UploadMedia(ctx, filename, data, mimeType)
		if err != nil {
			color.Red("アップロードに失敗: %s: %v", filename, err)
			continue
		}

		color.Green("アップロード完了: %s", filename)
		color.White("  メディアID: %d", media.ID)
		color.White("  URL: %s", media.SourceURL)
		fmt.Println()
	}

	return nil
}

// detectMediaMIMEType はファイル拡張子からMIMEタイプを判定する（画像・動画・音声対応）
func detectMediaMIMEType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	// 画像
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".webp":
		return "image/webp"
	case ".gif":
		return "image/gif"
	case ".svg":
		return "image/svg+xml"
	// 動画
	case ".mp4":
		return "video/mp4"
	case ".webm":
		return "video/webm"
	case ".mov":
		return "video/quicktime"
	case ".avi":
		return "video/x-msvideo"
	// 音声
	case ".mp3":
		return "audio/mpeg"
	case ".wav":
		return "audio/wav"
	case ".ogg":
		return "audio/ogg"
	case ".m4a":
		return "audio/mp4"
	default:
		return "application/octet-stream"
	}
}
