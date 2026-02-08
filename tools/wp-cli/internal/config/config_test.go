package config

import (
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

func TestLoadConfig_全環境変数未設定(t *testing.T) {
	isolateFromEnvFile(t)
	clearWPEnvVars(t)

	_, err := Load()
	if err == nil {
		t.Error("全環境変数未設定でエラーが発生するべき")
	}
}
