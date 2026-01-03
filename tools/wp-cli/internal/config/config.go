package config

import (
	"fmt"
	"os"
	"path/filepath"

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
		if err := godotenv.Load(envPath); err != nil {
			// .envファイルが見つからない場合は環境変数から読み込む
		}
	}

	siteURL := os.Getenv("WP_SITE_URL")
	username := os.Getenv("WP_USERNAME")
	appPassword := os.Getenv("WP_APP_PASSWORD")

	if siteURL == "" {
		return nil, fmt.Errorf("WP_SITE_URL が設定されていません")
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

// findEnvFile はプロジェクトルートの.envファイルを探す
func findEnvFile() string {
	// カレントディレクトリから上に向かって.envを探す
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}

	for {
		envPath := filepath.Join(dir, ".env")
		if _, err := os.Stat(envPath); err == nil {
			return envPath
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return ""
}
