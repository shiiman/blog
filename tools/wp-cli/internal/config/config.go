package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
)

// Config はアプリケーション設定を保持する
type Config struct {
	SiteURL     string
	Username    string
	AppPassword string
}

// Load は環境変数から設定を読み込む
func Load() (*Config, error) {
	// プロジェクトルートの.envを探す
	envPath := findEnvFile()
	if envPath != "" {
		// .envファイルの読み込みに失敗しても環境変数から読み込むため無視する
		_ = godotenv.Load(envPath)
	}

	siteURL := os.Getenv("WP_SITE_URL")
	username := os.Getenv("WP_USERNAME")
	appPassword := os.Getenv("WP_APP_PASSWORD")

	if siteURL == "" {
		return nil, fmt.Errorf("WP_SITE_URL が設定されていません")
	}
	if !strings.HasPrefix(siteURL, "https://") {
		return nil, fmt.Errorf("WP_SITE_URL は https:// で始まる必要があります（現在: %s）", siteURL)
	}
	if username == "" {
		return nil, fmt.Errorf("WP_USERNAME が設定されていません")
	}
	if appPassword == "" {
		return nil, fmt.Errorf("WP_APP_PASSWORD が設定されていません")
	}

	return &Config{
		SiteURL:     siteURL,
		Username:    username,
		AppPassword: appPassword,
	}, nil
}

// maxEnvSearchDepth は.envファイル探索の最大階層数
const maxEnvSearchDepth = 5

// findEnvFile はプロジェクトルートの.envファイルを探す
// .gitディレクトリが見つかった場合、または最大5階層まで探索する
func findEnvFile() string {
	// カレントディレクトリから上に向かって.envを探す
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}

	for depth := 0; depth < maxEnvSearchDepth; depth++ {
		envPath := filepath.Join(dir, ".env")
		if _, err := os.Stat(envPath); err == nil {
			return envPath
		}

		// .gitディレクトリが存在すればプロジェクトルートと見なして探索終了
		gitPath := filepath.Join(dir, ".git")
		if _, err := os.Stat(gitPath); err == nil {
			break
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return ""
}
