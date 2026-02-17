package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRootCommand_RejectsPublishFlagForPost(t *testing.T) {
	rootCmd.SetArgs([]string{"post", "dummy.md", "--publish"})
	_, err := rootCmd.ExecuteC()
	if err == nil {
		t.Fatal("post --publish は未知フラグエラーになるべき")
	}
	if !strings.Contains(err.Error(), "unknown flag: --publish") {
		t.Fatalf("期待したエラーではありません: %v", err)
	}
}

func TestRootCommand_RejectsPublishFlagForPage(t *testing.T) {
	rootCmd.SetArgs([]string{"page", "dummy.md", "--publish"})
	_, err := rootCmd.ExecuteC()
	if err == nil {
		t.Fatal("page --publish は未知フラグエラーになるべき")
	}
	if !strings.Contains(err.Error(), "unknown flag: --publish") {
		t.Fatalf("期待したエラーではありません: %v", err)
	}
}

func TestRunPage_InvalidStatusFailsBeforeAPI(t *testing.T) {
	filePath := writeTempArticle(t, `---
title: "test"
slug: "test"
status: "publsih"
---

body
`)

	err := runPage(pageCmd, []string{filePath})
	if err == nil {
		t.Fatal("不正statusはエラーになるべき")
	}
	if !strings.Contains(err.Error(), "page の status が不正です") {
		t.Fatalf("期待したエラーではありません: %v", err)
	}
}

func TestRunUpdate_InvalidStatusFailsBeforeAPI(t *testing.T) {
	filePath := writeTempArticle(t, `---
id: 1
title: "test"
slug: "test"
status: "publsih"
---

body
`)

	err := runUpdate(updateCmd, []string{filePath})
	if err == nil {
		t.Fatal("不正statusはエラーになるべき")
	}
	if !strings.Contains(err.Error(), "update の Front Matter status が不正です") {
		t.Fatalf("期待したエラーではありません: %v", err)
	}
}

func TestRunList_InvalidStatusFailsImmediately(t *testing.T) {
	prev := listStatus
	listStatus = "publsih"
	t.Cleanup(func() {
		listStatus = prev
	})

	err := runList(listCmd, []string{"posts"})
	if err == nil {
		t.Fatal("不正statusはエラーになるべき")
	}
	if !strings.Contains(err.Error(), "list --status が不正です") {
		t.Fatalf("期待したエラーではありません: %v", err)
	}
}

func TestRunList_DoesNotMutateGlobalListStatus(t *testing.T) {
	prev := listStatus
	original := " Publish "
	listStatus = original
	t.Cleanup(func() {
		listStatus = prev
	})

	// setupClientを確実に失敗させる（ネットワークに到達しない）
	t.Setenv("WP_SITE_URL", "")
	t.Setenv("WP_USERNAME", "")
	t.Setenv("WP_APP_PASSWORD", "")

	err := runList(listCmd, []string{"posts"})
	if err == nil {
		t.Fatal("設定エラーになるべき")
	}
	if !strings.Contains(err.Error(), "設定エラー") {
		t.Fatalf("期待した設定エラーではありません: %v", err)
	}
	if listStatus != original {
		t.Fatalf("listStatus が変更されています: got=%q want=%q", listStatus, original)
	}
}

func writeTempArticle(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	filePath := filepath.Join(dir, "article.md")
	if err := os.WriteFile(filePath, []byte(content), 0600); err != nil {
		t.Fatalf("一時記事ファイルの作成に失敗: %v", err)
	}
	return filePath
}
