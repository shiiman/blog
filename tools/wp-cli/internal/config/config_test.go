package config

import (
	"os"
	"path/filepath"
	"testing"
)

// テスト用に環境変数を設定するヘルパー
// t.Setenv はテスト終了時に自動的に元の値に復元する
func setEnvVars(t *testing.T, vars map[string]string) {
	t.Helper()
	for k, v := range vars {
		t.Setenv(k, v)
	}
}

func clearWPEnvVars(t *testing.T) {
	t.Helper()
	keys := []string{"WP_SITE_URL", "WP_USERNAME", "WP_APP_PASSWORD"}
	for _, k := range keys {
		t.Setenv(k, "")
	}
}

// isolateFromEnvFile はカレントディレクトリを一時ディレクトリに変更し、
// .envファイルの自動読み込みを防止する
func isolateFromEnvFile(t *testing.T) {
	t.Helper()
	t.Chdir(t.TempDir())
}

func TestLoadConfig_全環境変数が設定済み(t *testing.T) {
	isolateFromEnvFile(t)
	clearWPEnvVars(t)

	setEnvVars(t, map[string]string{
		"WP_SITE_URL":     "https://example.com",
		"WP_USERNAME":     "admin",
		"WP_APP_PASSWORD": "secret-password",
	})

	cfg, err := Load()
	if err != nil {
		t.Fatalf("エラーが発生: %v", err)
	}

	if cfg.SiteURL != "https://example.com" {
		t.Errorf("SiteURL = %q, want %q", cfg.SiteURL, "https://example.com")
	}
	if cfg.Username != "admin" {
		t.Errorf("Username = %q, want %q", cfg.Username, "admin")
	}
	if cfg.AppPassword != "secret-password" {
		t.Errorf("AppPassword = %q, want %q", cfg.AppPassword, "secret-password")
	}
}

func TestLoadConfig_SiteURL未設定(t *testing.T) {
	isolateFromEnvFile(t)
	clearWPEnvVars(t)

	setEnvVars(t, map[string]string{
		"WP_USERNAME":     "admin",
		"WP_APP_PASSWORD": "secret",
	})

	_, err := Load()
	if err == nil {
		t.Error("WP_SITE_URL未設定でエラーが発生するべき")
	}
}

func TestLoadConfig_Username未設定(t *testing.T) {
	isolateFromEnvFile(t)
	clearWPEnvVars(t)

	setEnvVars(t, map[string]string{
		"WP_SITE_URL":     "https://example.com",
		"WP_APP_PASSWORD": "secret",
	})

	_, err := Load()
	if err == nil {
		t.Error("WP_USERNAME未設定でエラーが発生するべき")
	}
}

func TestLoadConfig_AppPassword未設定(t *testing.T) {
	isolateFromEnvFile(t)
	clearWPEnvVars(t)

	setEnvVars(t, map[string]string{
		"WP_SITE_URL": "https://example.com",
		"WP_USERNAME": "admin",
	})

	_, err := Load()
	if err == nil {
		t.Error("WP_APP_PASSWORD未設定でエラーが発生するべき")
	}
}

func TestLoadConfig_HTTPのSiteURLはエラー(t *testing.T) {
	isolateFromEnvFile(t)
	clearWPEnvVars(t)

	setEnvVars(t, map[string]string{
		"WP_SITE_URL":     "http://example.com",
		"WP_USERNAME":     "admin",
		"WP_APP_PASSWORD": "secret",
	})

	_, err := Load()
	if err == nil {
		t.Error("http:// のSiteURLでエラーが発生するべき")
	}
}

func TestLoadConfig_全環境変数未設定(t *testing.T) {
	isolateFromEnvFile(t)
	clearWPEnvVars(t)

	_, err := Load()
	if err == nil {
		t.Error("全環境変数未設定でエラーが発生するべき")
	}
}

func TestFindEnvFile_gitディレクトリで停止(t *testing.T) {
	// 一時ディレクトリ構造を作成: root/.git/ + root/sub/
	root := t.TempDir()
	gitDir := filepath.Join(root, ".git")
	subDir := filepath.Join(root, "sub")
	if err := os.MkdirAll(gitDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatal(err)
	}

	// root/.envを作成
	envPath := filepath.Join(root, ".env")
	if err := os.WriteFile(envPath, []byte("TEST=1"), 0644); err != nil {
		t.Fatal(err)
	}

	// subDirに移動して探索
	t.Chdir(subDir)

	result := findEnvFile()
	if result != envPath {
		t.Errorf("findEnvFile() = %q, want %q", result, envPath)
	}
}

func TestFindEnvFile_深さ制限で停止(t *testing.T) {
	// maxEnvSearchDepth+1階層の深いディレクトリを作成（.gitなし）
	root := t.TempDir()
	deepDir := root
	for i := 0; i < maxEnvSearchDepth+1; i++ {
		deepDir = filepath.Join(deepDir, "sub")
	}
	if err := os.MkdirAll(deepDir, 0755); err != nil {
		t.Fatal(err)
	}

	// rootに.envを作成（深すぎて見つからないはず）
	envPath := filepath.Join(root, ".env")
	if err := os.WriteFile(envPath, []byte("TEST=1"), 0644); err != nil {
		t.Fatal(err)
	}

	t.Chdir(deepDir)

	result := findEnvFile()
	if result != "" {
		t.Errorf("深すぎるディレクトリで.envが見つかるべきでない: %q", result)
	}
}
