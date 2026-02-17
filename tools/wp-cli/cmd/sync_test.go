package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/shiimanblog/wp-cli/internal/types"
)

// テスト用のヘルパー: drafts/ディレクトリ構造を持つ一時記事ファイルを作成する
func writeTempArticleInDrafts(t *testing.T, slug, content string) (projectRoot, filePath string) {
	t.Helper()
	projectRoot = t.TempDir()
	draftsDir := filepath.Join(projectRoot, "drafts", slug)
	if err := os.MkdirAll(draftsDir, 0755); err != nil {
		t.Fatalf("draftsディレクトリの作成に失敗: %v", err)
	}
	filePath = filepath.Join(draftsDir, "article.md")
	if err := os.WriteFile(filePath, []byte(content), 0600); err != nil {
		t.Fatalf("一時記事ファイルの作成に失敗: %v", err)
	}
	return projectRoot, filePath
}

// テスト用のWPTime生成ヘルパー
func mustWPTime(t *testing.T, value string) types.WPTime {
	t.Helper()
	parsed, err := time.Parse("2006-01-02T15:04:05", value)
	if err != nil {
		t.Fatalf("WPTime のパースに失敗: %v", err)
	}
	return types.WPTime{Time: parsed}
}

func TestSyncPostToLocal(t *testing.T) {
	baseArticle := `---
title: "テスト記事"
slug: "test-article"
status: draft
---

テスト本文
`

	tests := []struct {
		name          string
		article       string
		post          *types.Post
		moveOnPublish bool
		wantErr       bool
		wantMoved     bool // drafts/ → posts/ への移動が行われるか
		checkFn       func(t *testing.T, resultPath string) // 追加の検証
	}{
		{
			name:    "正常系: draft投稿のフロントマターを同期",
			article: baseArticle,
			post: &types.Post{
				ID:            42,
				Status:        "draft",
				Date:          mustWPTime(t, "2026-01-15T10:00:00"),
				Modified:      mustWPTime(t, "2026-01-15T12:00:00"),
				Slug:          "test-article",
				FeaturedMedia: 100,
			},
			moveOnPublish: false,
			wantErr:       false,
			wantMoved:     false,
			checkFn: func(t *testing.T, resultPath string) {
				content, err := os.ReadFile(resultPath)
				if err != nil {
					t.Fatalf("ファイル読み込みに失敗: %v", err)
				}
				s := string(content)
				if !strings.Contains(s, "id: 42") {
					t.Error("IDが同期されていない")
				}
				if !strings.Contains(s, "status: draft") {
					t.Error("statusが同期されていない")
				}
				if !strings.Contains(s, "featured_media: 100") {
					t.Error("featured_mediaが同期されていない")
				}
				if !strings.Contains(s, "slug: test-article") {
					t.Error("slugが同期されていない")
				}
			},
		},
		{
			name:    "正常系: publish投稿でmoveOnPublish=trueだがdrafts/外なので移動しない",
			article: baseArticle,
			post: &types.Post{
				ID:       1,
				Status:   "publish",
				Date:     mustWPTime(t, "2026-02-01T09:00:00"),
				Modified: mustWPTime(t, "2026-02-01T09:00:00"),
				Slug:     "test-article",
			},
			moveOnPublish: true,
			wantErr:       false,
			wantMoved:     false, // writeTempArticle はdrafts/内に作らない
		},
		{
			name:    "正常系: publish投稿でmoveOnPublish=false",
			article: baseArticle,
			post: &types.Post{
				ID:       2,
				Status:   "publish",
				Date:     mustWPTime(t, "2026-02-01T09:00:00"),
				Modified: mustWPTime(t, "2026-02-01T09:00:00"),
				Slug:     "test-article",
			},
			moveOnPublish: false,
			wantErr:       false,
			wantMoved:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath := writeTempArticle(t, tt.article)

			resultPath, err := syncPostToLocal(filePath, tt.post, tt.moveOnPublish)
			if (err != nil) != tt.wantErr {
				t.Fatalf("syncPostToLocal() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				return
			}

			// ファイルが存在するか確認
			if _, statErr := os.Stat(resultPath); statErr != nil {
				t.Fatalf("結果ファイルが存在しない: %s", resultPath)
			}

			if tt.checkFn != nil {
				tt.checkFn(t, resultPath)
			}
		})
	}
}

func TestSyncPostToLocal_MoveFromDrafts(t *testing.T) {
	article := `---
title: "テスト記事"
slug: "test-article"
status: draft
---

テスト本文
`
	post := &types.Post{
		ID:       42,
		Status:   "publish",
		Date:     mustWPTime(t, "2026-01-15T10:00:00"),
		Modified: mustWPTime(t, "2026-01-15T12:00:00"),
		Slug:     "test-article",
	}

	projectRoot, filePath := writeTempArticleInDrafts(t, "test-article", article)

	// posts/ ディレクトリを事前に作成
	postsDir := filepath.Join(projectRoot, "posts")
	if err := os.MkdirAll(postsDir, 0755); err != nil {
		t.Fatalf("postsディレクトリの作成に失敗: %v", err)
	}

	resultPath, err := syncPostToLocal(filePath, post, true)
	if err != nil {
		t.Fatalf("syncPostToLocal() が予期せずエラー: %v", err)
	}

	// 移動先パスの確認
	expectedDirName := "2026-01-15_test-article"
	if !strings.Contains(resultPath, filepath.Join("posts", expectedDirName)) {
		t.Errorf("移動先パスが正しくない: got=%s, want contains %s", resultPath, expectedDirName)
	}

	// 元のファイルは存在しないはず
	if _, statErr := os.Stat(filePath); !os.IsNotExist(statErr) {
		t.Error("元のファイルがまだ存在している")
	}

	// 新しいファイルが存在するはず
	if _, statErr := os.Stat(resultPath); statErr != nil {
		t.Fatalf("移動先のファイルが存在しない: %s", resultPath)
	}

	// フロントマターが更新されているか確認
	content, err := os.ReadFile(resultPath)
	if err != nil {
		t.Fatalf("ファイル読み込みに失敗: %v", err)
	}
	s := string(content)
	if !strings.Contains(s, "id: 42") {
		t.Error("IDが同期されていない")
	}
	if !strings.Contains(s, "status: publish") {
		t.Error("statusが同期されていない")
	}
}

func TestSyncPostToLocal_InvalidFile(t *testing.T) {
	post := &types.Post{
		ID:       1,
		Status:   "draft",
		Date:     mustWPTime(t, "2026-01-01T00:00:00"),
		Modified: mustWPTime(t, "2026-01-01T00:00:00"),
		Slug:     "test",
	}

	_, err := syncPostToLocal("/nonexistent/path/article.md", post, false)
	if err == nil {
		t.Fatal("存在しないファイルの場合はエラーになるべき")
	}
	if !strings.Contains(err.Error(), "フロントマター読み込みに失敗") {
		t.Errorf("期待したエラーメッセージでない: %v", err)
	}
}

func TestSyncPageToLocal(t *testing.T) {
	tests := []struct {
		name    string
		article string
		page    *types.Page
		wantErr bool
		checkFn func(t *testing.T, resultPath string)
	}{
		{
			name: "正常系: 固定ページのフロントマターを同期",
			article: `---
title: "テストページ"
slug: "test-page"
status: draft
---

テストページ本文
`,
			page: &types.Page{
				ID:        99,
				Status:    "publish",
				Date:      mustWPTime(t, "2026-02-10T14:00:00"),
				Modified:  mustWPTime(t, "2026-02-10T15:00:00"),
				Slug:      "test-page",
				Parent:    5,
				MenuOrder: 3,
			},
			wantErr: false,
			checkFn: func(t *testing.T, resultPath string) {
				content, err := os.ReadFile(resultPath)
				if err != nil {
					t.Fatalf("ファイル読み込みに失敗: %v", err)
				}
				s := string(content)
				if !strings.Contains(s, "id: 99") {
					t.Error("IDが同期されていない")
				}
				if !strings.Contains(s, "status: publish") {
					t.Error("statusが同期されていない")
				}
				if !strings.Contains(s, "slug: test-page") {
					t.Error("slugが同期されていない")
				}
				if !strings.Contains(s, "parent: 5") {
					t.Error("parentが同期されていない")
				}
				if !strings.Contains(s, "menu_order: 3") {
					t.Error("menu_orderが同期されていない")
				}
			},
		},
		{
			name: "正常系: parent=0, menu_order=0のページ",
			article: `---
title: "ルートページ"
slug: "root-page"
---

本文
`,
			page: &types.Page{
				ID:        50,
				Status:    "draft",
				Date:      mustWPTime(t, "2026-01-01T00:00:00"),
				Modified:  mustWPTime(t, "2026-01-01T00:00:00"),
				Slug:      "root-page",
				Parent:    0,
				MenuOrder: 0,
			},
			wantErr: false,
			checkFn: func(t *testing.T, resultPath string) {
				content, err := os.ReadFile(resultPath)
				if err != nil {
					t.Fatalf("ファイル読み込みに失敗: %v", err)
				}
				s := string(content)
				if !strings.Contains(s, "id: 50") {
					t.Error("IDが同期されていない")
				}
				// parent=0, menu_order=0 は omitempty により出力されない
				if strings.Contains(s, "parent:") {
					t.Error("parent=0 は出力されるべきでない (omitempty)")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath := writeTempArticle(t, tt.article)

			resultPath, err := syncPageToLocal(filePath, tt.page)
			if (err != nil) != tt.wantErr {
				t.Fatalf("syncPageToLocal() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				return
			}

			if _, statErr := os.Stat(resultPath); statErr != nil {
				t.Fatalf("結果ファイルが存在しない: %s", resultPath)
			}

			if tt.checkFn != nil {
				tt.checkFn(t, resultPath)
			}
		})
	}
}

func TestSyncPageToLocal_InvalidFile(t *testing.T) {
	page := &types.Page{
		ID:       1,
		Status:   "draft",
		Date:     mustWPTime(t, "2026-01-01T00:00:00"),
		Modified: mustWPTime(t, "2026-01-01T00:00:00"),
		Slug:     "test",
	}

	_, err := syncPageToLocal("/nonexistent/path/page.md", page)
	if err == nil {
		t.Fatal("存在しないファイルの場合はエラーになるべき")
	}
	if !strings.Contains(err.Error(), "フロントマター読み込みに失敗") {
		t.Errorf("期待したエラーメッセージでない: %v", err)
	}
}

func TestMoveToPublished(t *testing.T) {
	tests := []struct {
		name      string
		setupFn   func(t *testing.T) (filePath string, cleanup func())
		post      *types.Post
		wantErr   bool
		errContains string
		checkFn   func(t *testing.T, resultPath, projectRoot string)
	}{
		{
			name: "正常系: drafts/からposts/へ移動",
			setupFn: func(t *testing.T) (string, func()) {
				root := t.TempDir()
				slug := "my-article"
				draftsDir := filepath.Join(root, "drafts", slug)
				if err := os.MkdirAll(draftsDir, 0755); err != nil {
					t.Fatal(err)
				}
				postsDir := filepath.Join(root, "posts")
				if err := os.MkdirAll(postsDir, 0755); err != nil {
					t.Fatal(err)
				}
				fp := filepath.Join(draftsDir, "article.md")
				if err := os.WriteFile(fp, []byte("test"), 0600); err != nil {
					t.Fatal(err)
				}
				return fp, func() {}
			},
			post: &types.Post{
				Date: mustWPTime(t, "2026-03-01T10:00:00"),
				Slug: "my-article",
			},
			wantErr: false,
			checkFn: func(t *testing.T, resultPath, _ string) {
				if !strings.Contains(resultPath, "posts") {
					t.Error("移動先パスにposts/が含まれていない")
				}
				if !strings.Contains(resultPath, "2026-03-01_my-article") {
					t.Error("移動先ディレクトリ名が正しくない")
				}
				if _, err := os.Stat(resultPath); err != nil {
					t.Errorf("移動先のファイルが存在しない: %v", err)
				}
			},
		},
		{
			name: "異常系: draftsディレクトリが見つからない",
			setupFn: func(t *testing.T) (string, func()) {
				root := t.TempDir()
				dir := filepath.Join(root, "somewhere")
				if err := os.MkdirAll(dir, 0755); err != nil {
					t.Fatal(err)
				}
				fp := filepath.Join(dir, "article.md")
				if err := os.WriteFile(fp, []byte("test"), 0600); err != nil {
					t.Fatal(err)
				}
				return fp, func() {}
			},
			post: &types.Post{
				Date: mustWPTime(t, "2026-01-01T00:00:00"),
				Slug: "test",
			},
			wantErr:     true,
			errContains: "draftsディレクトリが見つかりません",
		},
		{
			name: "異常系: 移動先が既に存在する",
			setupFn: func(t *testing.T) (string, func()) {
				root := t.TempDir()
				slug := "existing-article"
				draftsDir := filepath.Join(root, "drafts", slug)
				if err := os.MkdirAll(draftsDir, 0755); err != nil {
					t.Fatal(err)
				}
				fp := filepath.Join(draftsDir, "article.md")
				if err := os.WriteFile(fp, []byte("test"), 0600); err != nil {
					t.Fatal(err)
				}
				// 移動先ディレクトリを事前に作成
				destDir := filepath.Join(root, "posts", "2026-01-01_existing-article")
				if err := os.MkdirAll(destDir, 0755); err != nil {
					t.Fatal(err)
				}
				return fp, func() {}
			},
			post: &types.Post{
				Date: mustWPTime(t, "2026-01-01T00:00:00"),
				Slug: "existing-article",
			},
			wantErr:     true,
			errContains: "移動先が既に存在します",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath, cleanup := tt.setupFn(t)
			defer cleanup()

			resultPath, err := moveToPublished(filePath, tt.post)
			if (err != nil) != tt.wantErr {
				t.Fatalf("moveToPublished() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("エラーメッセージが期待と異なる: got=%v, want contains=%q", err, tt.errContains)
				}
				return
			}

			if tt.checkFn != nil {
				tt.checkFn(t, resultPath, "")
			}
		})
	}
}

func TestMoveToPublished_PathTraversal(t *testing.T) {
	root := t.TempDir()
	slug := "evil-article"
	draftsDir := filepath.Join(root, "drafts", slug)
	if err := os.MkdirAll(draftsDir, 0755); err != nil {
		t.Fatal(err)
	}
	postsDir := filepath.Join(root, "posts")
	if err := os.MkdirAll(postsDir, 0755); err != nil {
		t.Fatal(err)
	}
	fp := filepath.Join(draftsDir, "article.md")
	if err := os.WriteFile(fp, []byte("test"), 0600); err != nil {
		t.Fatal(err)
	}

	// パストラバーサルを含むスラッグ
	post := &types.Post{
		Date: mustWPTime(t, "2026-01-01T00:00:00"),
		Slug: "../../etc/evil",
	}

	// sanitizeSlug が適用されるので、パストラバーサルは防がれるはず
	resultPath, err := moveToPublished(fp, post)
	if err != nil {
		// エラーでも可（パストラバーサル防止で弾かれる場合）
		return
	}
	// エラーにならない場合でも、posts/ディレクトリ内にあることを確認
	if !strings.Contains(resultPath, "posts") {
		t.Errorf("パストラバーサルが防がれていない: %s", resultPath)
	}
}

func TestUploadEyecatchIfExists(t *testing.T) {
	tests := []struct {
		name                   string
		setupFn                func(t *testing.T) string // articleDir を返す
		currentFeaturedMediaID int
		forceUpload            bool
		wantID                 int
		wantErr                bool
		errContains            string
	}{
		{
			name: "正常系: アイキャッチ画像が存在しない場合は現在のIDを返す",
			setupFn: func(t *testing.T) string {
				return t.TempDir() // assets/eyecatch.png なし
			},
			currentFeaturedMediaID: 50,
			forceUpload:            false,
			wantID:                 50,
			wantErr:                false,
		},
		{
			name: "正常系: アイキャッチ画像が存在しない場合（ID=0）",
			setupFn: func(t *testing.T) string {
				return t.TempDir()
			},
			currentFeaturedMediaID: 0,
			forceUpload:            false,
			wantID:                 0,
			wantErr:                false,
		},
		{
			name: "正常系: 既にFeaturedMediaが設定済みで強制アップロードなし",
			setupFn: func(t *testing.T) string {
				dir := t.TempDir()
				assetsDir := filepath.Join(dir, "assets")
				if err := os.MkdirAll(assetsDir, 0755); err != nil {
					t.Fatal(err)
				}
				// 小さなダミー画像ファイルを作成
				if err := os.WriteFile(filepath.Join(assetsDir, "eyecatch.png"), []byte("fake-png-data"), 0600); err != nil {
					t.Fatal(err)
				}
				return dir
			},
			currentFeaturedMediaID: 100,
			forceUpload:            false,
			wantID:                 100, // 既存IDがそのまま返る
			wantErr:                false,
		},
		{
			name: "異常系: アイキャッチ画像のサイズが上限を超えている",
			setupFn: func(t *testing.T) string {
				dir := t.TempDir()
				assetsDir := filepath.Join(dir, "assets")
				if err := os.MkdirAll(assetsDir, 0755); err != nil {
					t.Fatal(err)
				}
				// 20MB超のダミーファイルを作成
				largeData := make([]byte, 21*1024*1024)
				if err := os.WriteFile(filepath.Join(assetsDir, "eyecatch.png"), largeData, 0600); err != nil {
					t.Fatal(err)
				}
				return dir
			},
			currentFeaturedMediaID: 0,
			forceUpload:            false,
			wantErr:                true,
			errContains:            "サイズが上限",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			articleDir := tt.setupFn(t)

			ctx := t.Context()
			gotID, err := uploadEyecatchIfExists(ctx, nil, articleDir, tt.currentFeaturedMediaID, tt.forceUpload)
			if (err != nil) != tt.wantErr {
				t.Fatalf("uploadEyecatchIfExists() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("エラーメッセージが期待と異なる: got=%v, want contains=%q", err, tt.errContains)
				}
				return
			}
			if gotID != tt.wantID {
				t.Errorf("uploadEyecatchIfExists() = %d, want %d", gotID, tt.wantID)
			}
		})
	}
}

func TestSaveArticle(t *testing.T) {
	tests := []struct {
		name        string
		// setupFn はCWDを基準にdirPathを返す（saveArticleがCWD基準でパス検証するため）
		setupFn     func(t *testing.T) (dirPath, filename string)
		article     *types.Article
		wantErr     bool
		errContains string
		checkFn     func(t *testing.T, dirPath, filename string)
	}{
		{
			name: "正常系: ディレクトリとファイルを作成",
			setupFn: func(t *testing.T) (string, string) {
				// saveArticleはCWD基準でパス検証するので、CWD内の相対パスを使う
				dirPath := filepath.Join("posts", "2026-01-01_test")
				return dirPath, "article.md"
			},
			article: &types.Article{
				FrontMatter: types.FrontMatter{
					ID:     1,
					Title:  "テスト記事",
					Slug:   "test",
					Status: "draft",
				},
				Content: "テスト本文です",
			},
			wantErr: false,
			checkFn: func(t *testing.T, dirPath, filename string) {
				fp := filepath.Join(dirPath, filename)
				content, err := os.ReadFile(fp)
				if err != nil {
					t.Fatalf("ファイル読み込みに失敗: %v", err)
				}
				s := string(content)
				if !strings.Contains(s, "title: テスト記事") {
					t.Error("タイトルが保存されていない")
				}
				if !strings.Contains(s, "テスト本文です") {
					t.Error("本文が保存されていない")
				}

				// assetsディレクトリが作成されているか確認
				assetsDir := filepath.Join(dirPath, "assets")
				if _, err := os.Stat(assetsDir); err != nil {
					t.Errorf("assetsディレクトリが作成されていない: %v", err)
				}
			},
		},
		{
			name: "正常系: page.md ファイルを作成",
			setupFn: func(t *testing.T) (string, string) {
				dirPath := filepath.Join("pages", "about")
				return dirPath, "page.md"
			},
			article: &types.Article{
				FrontMatter: types.FrontMatter{
					ID:     10,
					Title:  "概要ページ",
					Slug:   "about",
					Status: "publish",
				},
				Content: "概要ページの本文",
			},
			wantErr: false,
			checkFn: func(t *testing.T, dirPath, filename string) {
				fp := filepath.Join(dirPath, filename)
				if _, err := os.Stat(fp); err != nil {
					t.Errorf("ファイルが作成されていない: %v", err)
				}
			},
		},
		{
			name: "異常系: パストラバーサル（ベースディレクトリ外）",
			setupFn: func(t *testing.T) (string, string) {
				return "/tmp/../../../etc/evil", "article.md"
			},
			article: &types.Article{
				FrontMatter: types.FrontMatter{
					Title: "evil",
				},
				Content: "evil content",
			},
			wantErr:     true,
			errContains: "不正なパス",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// saveArticle はCWD基準でパストラバーサルチェックするため、
			// CWDをtmpディレクトリに変更する
			tmpRoot := t.TempDir()
			t.Chdir(tmpRoot)

			dirPath, filename := tt.setupFn(t)

			err := saveArticle(dirPath, filename, tt.article)
			if (err != nil) != tt.wantErr {
				t.Fatalf("saveArticle() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				if tt.errContains != "" && err != nil && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("エラーメッセージが期待と異なる: got=%v, want contains=%q", err, tt.errContains)
				}
				return
			}

			if tt.checkFn != nil {
				tt.checkFn(t, dirPath, filename)
			}
		})
	}
}

func TestSaveArticle_ExistingDirectory(t *testing.T) {
	// saveArticle はCWD基準でパストラバーサルチェックするため、tmpにCWDを変更
	tmpRoot := t.TempDir()
	t.Chdir(tmpRoot)

	dirPath := filepath.Join("posts", "existing")

	// 1回目: 新規作成
	article := &types.Article{
		FrontMatter: types.FrontMatter{
			Title:  "最初の記事",
			Slug:   "existing",
			Status: "draft",
		},
		Content: "最初の本文",
	}
	if err := saveArticle(dirPath, "article.md", article); err != nil {
		t.Fatalf("1回目のsaveArticle()がエラー: %v", err)
	}

	// 2回目: 上書き更新
	article.FrontMatter.Title = "更新された記事"
	article.Content = "更新された本文"
	if err := saveArticle(dirPath, "article.md", article); err != nil {
		t.Fatalf("2回目のsaveArticle()がエラー: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(dirPath, "article.md"))
	if err != nil {
		t.Fatalf("ファイル読み込みに失敗: %v", err)
	}
	s := string(content)
	if !strings.Contains(s, "更新された記事") {
		t.Error("タイトルが更新されていない")
	}
	if !strings.Contains(s, "更新された本文") {
		t.Error("本文が更新されていない")
	}
}
