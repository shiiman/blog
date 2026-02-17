package converter

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/shiimanblog/wp-cli/internal/types"
)

// --- TestParseFrontMatter ---

func TestParseFrontMatter_正常なFrontMatter(t *testing.T) {
	content := `---
title: "テスト記事"
slug: "test-article"
status: draft
categories:
  - 1
  - 2
tags:
  - 10
  - 20
---

これは本文です。`

	fm, body, err := parseFrontMatter(content)
	if err != nil {
		t.Fatalf("エラーが発生: %v", err)
	}

	if fm.Title != "テスト記事" {
		t.Errorf("Title = %q, want %q", fm.Title, "テスト記事")
	}
	if fm.Slug != "test-article" {
		t.Errorf("Slug = %q, want %q", fm.Slug, "test-article")
	}
	if fm.Status != "draft" {
		t.Errorf("Status = %q, want %q", fm.Status, "draft")
	}
	if len(fm.Categories) != 2 || fm.Categories[0] != 1 || fm.Categories[1] != 2 {
		t.Errorf("Categories = %v, want [1, 2]", fm.Categories)
	}
	if len(fm.Tags) != 2 || fm.Tags[0] != 10 || fm.Tags[1] != 20 {
		t.Errorf("Tags = %v, want [10, 20]", fm.Tags)
	}
	if body != "これは本文です。" {
		t.Errorf("Body = %q, want %q", body, "これは本文です。")
	}
}

func TestParseFrontMatter_FrontMatterなし(t *testing.T) {
	content := "これはFront Matterのない本文です。"

	fm, body, err := parseFrontMatter(content)
	if err != nil {
		t.Fatalf("エラーが発生: %v", err)
	}

	// Front Matterがない場合、空のFrontMatterが返される
	if fm.Title != "" {
		t.Errorf("Title = %q, want empty", fm.Title)
	}
	if body != content {
		t.Errorf("Body = %q, want %q", body, content)
	}
}

func TestParseFrontMatter_不正なYAML(t *testing.T) {
	content := `---
title: [invalid yaml
  : broken
---

本文`

	_, _, err := parseFrontMatter(content)
	if err == nil {
		t.Error("不正なYAMLでエラーが発生するべき")
	}
}

func TestParseFrontMatter_IDとメタデータ付き(t *testing.T) {
	content := `---
id: 123
title: "更新記事"
slug: "update-article"
status: publish
featured_media: 456
date: "2024-01-01T00:00:00"
modified: "2024-06-01T12:00:00"
---

更新された記事の本文`

	fm, body, err := parseFrontMatter(content)
	if err != nil {
		t.Fatalf("エラーが発生: %v", err)
	}

	if fm.ID != 123 {
		t.Errorf("ID = %d, want 123", fm.ID)
	}
	if fm.FeaturedMedia != 456 {
		t.Errorf("FeaturedMedia = %d, want 456", fm.FeaturedMedia)
	}
	if fm.Date != "2024-01-01T00:00:00" {
		t.Errorf("Date = %q, want %q", fm.Date, "2024-01-01T00:00:00")
	}
	if fm.Modified != "2024-06-01T12:00:00" {
		t.Errorf("Modified = %q, want %q", fm.Modified, "2024-06-01T12:00:00")
	}
	if body != "更新された記事の本文" {
		t.Errorf("Body = %q, want %q", body, "更新された記事の本文")
	}
}

// --- TestHTMLToMarkdown ---

func TestHTMLToMarkdown_基本的なHTML(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		contains string
	}{
		{
			name:     "段落タグ",
			html:     "<p>こんにちは世界</p>",
			contains: "こんにちは世界",
		},
		{
			name:     "見出しタグ",
			html:     "<h2>見出し2</h2>",
			contains: "## 見出し2",
		},
		{
			name:     "リンク",
			html:     `<a href="https://example.com">リンクテキスト</a>`,
			contains: "[リンクテキスト](https://example.com)",
		},
		{
			name:     "太字",
			html:     "<strong>太字テキスト</strong>",
			contains: "**太字テキスト**",
		},
		{
			name:     "リスト",
			html:     "<ul><li>項目1</li><li>項目2</li></ul>",
			contains: "項目1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := HTMLToMarkdown(tt.html)
			if err != nil {
				t.Fatalf("エラーが発生: %v", err)
			}
			if !strings.Contains(result, tt.contains) {
				t.Errorf("結果 %q に %q が含まれていない", result, tt.contains)
			}
		})
	}
}

func TestHTMLToMarkdown_空文字列(t *testing.T) {
	result, err := HTMLToMarkdown("")
	if err != nil {
		t.Fatalf("エラーが発生: %v", err)
	}
	if result != "" {
		t.Errorf("空HTMLの結果 = %q, want empty", result)
	}
}

// --- TestMarkdownToHTML ---

func TestMarkdownToHTML_基本的なMarkdown(t *testing.T) {
	tests := []struct {
		name     string
		md       string
		contains string
	}{
		{
			name:     "見出し",
			md:       "## テスト見出し",
			contains: "<h2",
		},
		{
			name:     "段落",
			md:       "これはテスト段落です。",
			contains: "<p>これはテスト段落です。</p>",
		},
		{
			name:     "コードブロック",
			md:       "```go\nfmt.Println(\"hello\")\n```",
			contains: "<code",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MarkdownToHTML(tt.md)
			if !strings.Contains(result, tt.contains) {
				t.Errorf("結果 %q に %q が含まれていない", result, tt.contains)
			}
		})
	}
}

func TestMarkdownToHTML_StyleタグはHTMLブロックでラップ(t *testing.T) {
	md := "<style>.test { color: red; }</style>\n\n## 見出し"
	result := MarkdownToHTML(md)

	if !strings.Contains(result, "<!-- wp:html -->") {
		t.Errorf("styleタグを含むHTMLにwp:htmlブロックが含まれていない: %q", result)
	}
}

func TestMarkdownToHTML_ScriptタグはHTMLブロックでラップ(t *testing.T) {
	md := "<script>console.log('test')</script>"
	result := MarkdownToHTML(md)

	if !strings.Contains(result, "<!-- wp:html -->") {
		t.Errorf("scriptタグを含むHTMLにwp:htmlブロックが含まれていない: %q", result)
	}
}

// --- TestParseArticle ---

func TestParseArticle_正常なファイル(t *testing.T) {
	// 一時ファイルを作成
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test-article.md")

	content := `---
id: 42
title: "テスト記事タイトル"
slug: "test-article"
status: draft
categories:
  - 1
tags:
  - 10
---

## テスト本文

これはテスト記事の本文です。`

	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatalf("テストファイルの作成に失敗: %v", err)
	}

	article, err := ParseArticle(filePath)
	if err != nil {
		t.Fatalf("エラーが発生: %v", err)
	}

	if article.FrontMatter.ID != 42 {
		t.Errorf("ID = %d, want 42", article.FrontMatter.ID)
	}
	if article.FrontMatter.Title != "テスト記事タイトル" {
		t.Errorf("Title = %q, want %q", article.FrontMatter.Title, "テスト記事タイトル")
	}
	if article.FrontMatter.Slug != "test-article" {
		t.Errorf("Slug = %q, want %q", article.FrontMatter.Slug, "test-article")
	}
	if article.FilePath != filePath {
		t.Errorf("FilePath = %q, want %q", article.FilePath, filePath)
	}
	if !strings.Contains(article.Content, "テスト本文") {
		t.Errorf("Content に 'テスト本文' が含まれていない: %q", article.Content)
	}
}

func TestParseArticle_存在しないファイル(t *testing.T) {
	_, err := ParseArticle("/nonexistent/path/test.md")
	if err == nil {
		t.Error("存在しないファイルでエラーが発生するべき")
	}
}

// --- TestGenerateArticleFile ---

func TestGenerateArticleFile_正常な記事(t *testing.T) {
	article := &types.Article{
		FrontMatter: types.FrontMatter{
			ID:     100,
			Title:  "生成テスト",
			Slug:   "generate-test",
			Status: "draft",
		},
		Content: "## 本文\n\nテスト内容です。",
	}

	result, err := GenerateArticleFile(article)
	if err != nil {
		t.Fatalf("エラーが発生: %v", err)
	}

	if !strings.Contains(result, "---") {
		t.Error("結果にFront Matterの区切り線が含まれていない")
	}
	if !strings.Contains(result, "title: 生成テスト") {
		t.Errorf("結果にタイトルが含まれていない: %q", result)
	}
	if !strings.Contains(result, "テスト内容です。") {
		t.Errorf("結果に本文が含まれていない: %q", result)
	}
}

// --- TestWrapInGutenbergBlocks ---

func TestWrapInGutenbergBlocks_通常のHTML(t *testing.T) {
	html := "<p>通常の段落</p>"
	result := wrapInGutenbergBlocks(html)

	// 通常のHTMLはそのまま返される
	if result != html {
		t.Errorf("通常のHTMLが変更された: %q", result)
	}
}

func TestWrapInGutenbergBlocks_Styleタグ含む(t *testing.T) {
	html := "<style>.test{}</style><p>テスト</p>"
	result := wrapInGutenbergBlocks(html)

	if !strings.Contains(result, "<!-- wp:html -->") {
		t.Error("styleタグを含むHTMLにwp:htmlが含まれていない")
	}
	if !strings.Contains(result, "<!-- /wp:html -->") {
		t.Error("styleタグを含むHTMLに/wp:htmlが含まれていない")
	}
}

func TestWrapInGutenbergBlocks_Scriptタグ含む(t *testing.T) {
	html := "<script>alert('test')</script>"
	result := wrapInGutenbergBlocks(html)

	if !strings.Contains(result, "<!-- wp:html -->") {
		t.Error("scriptタグを含むHTMLにwp:htmlが含まれていない")
	}
}

// --- TestProcessHtmlAndMarkdownMixed ---

func TestProcessHtmlAndMarkdownMixed_ScriptタグとMarkdown混在(t *testing.T) {
	content := "<script>console.log('test')</script>\n\n## 見出し\n\nMarkdown本文"
	result := processHtmlAndMarkdownMixed(content)

	if !strings.Contains(result, "<!-- wp:html -->") {
		t.Error("wp:htmlブロックが含まれていない")
	}
	if !strings.Contains(result, "<script>") {
		t.Error("scriptタグが保持されていない")
	}
}

func TestProcessHtmlAndMarkdownMixed_Scriptタグのみ(t *testing.T) {
	content := "<script>console.log('only script')</script>"
	result := processHtmlAndMarkdownMixed(content)

	if !strings.Contains(result, "<!-- wp:html -->") {
		t.Error("wp:htmlブロックが含まれていない")
	}
}

func TestProcessHtmlAndMarkdownMixed_Scriptなしのコンテンツ(t *testing.T) {
	content := "<div>HTMLブロックのみ</div>"
	result := processHtmlAndMarkdownMixed(content)

	// scriptタグがないのでwrapInGutenbergBlocksに委譲される
	if strings.Contains(result, "<!-- wp:html -->") {
		t.Error("scriptなしのHTMLにwp:htmlブロックが含まれるべきでない")
	}
}

// --- TestIsMarkdownPrefix ---

func TestIsMarkdownPrefix(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{name: "見出し", input: "## 見出し", want: true},
		{name: "リスト（ハイフン）", input: "- 項目", want: true},
		{name: "リスト（アスタリスク）", input: "* 項目", want: true},
		{name: "番号付きリスト", input: "1. 項目", want: true},
		{name: "引用", input: "> 引用文", want: true},
		{name: "コードブロック", input: "```go\ncode\n```", want: true},
		{name: "HTMLタグ", input: "<div>test</div>", want: false},
		{name: "通常テキスト", input: "plain text", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isMarkdownPrefix(tt.input)
			if got != tt.want {
				t.Errorf("isMarkdownPrefix(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

// --- TestPostToArticle ---

func TestPostToArticle_正常な投稿(t *testing.T) {
	post := &types.Post{
		ID:            1,
		Slug:          "test-post",
		Status:        "publish",
		Date:          types.WPTime{Time: time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)},
		Modified:      types.WPTime{Time: time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)},
		Title:         types.Rendered{Rendered: "テスト投稿"},
		Content:       types.Rendered{Rendered: "<p>テスト本文</p>"},
		Excerpt:       types.Rendered{Rendered: "<p>テスト概要</p>"},
		Categories:    []int{1, 2},
		Tags:          []int{10},
		FeaturedMedia: 100,
	}

	article, err := PostToArticle(post)
	if err != nil {
		t.Fatalf("エラーが発生: %v", err)
	}

	if article.FrontMatter.ID != 1 {
		t.Errorf("ID = %d, want 1", article.FrontMatter.ID)
	}
	if article.FrontMatter.Title != "テスト投稿" {
		t.Errorf("Title = %q, want %q", article.FrontMatter.Title, "テスト投稿")
	}
	if article.FrontMatter.FeaturedMedia != 100 {
		t.Errorf("FeaturedMedia = %d, want 100", article.FrontMatter.FeaturedMedia)
	}
	if !strings.Contains(article.FrontMatter.Excerpt, "テスト概要") {
		t.Errorf("Excerpt にテスト概要が含まれていない: %q", article.FrontMatter.Excerpt)
	}
}

func TestPostToArticle_Excerptなし(t *testing.T) {
	post := &types.Post{
		ID:       2,
		Slug:     "no-excerpt",
		Status:   "draft",
		Date:     types.WPTime{Time: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
		Modified: types.WPTime{Time: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
		Title:    types.Rendered{Rendered: "概要なし投稿"},
		Content:  types.Rendered{Rendered: "<p>本文のみ</p>"},
		Excerpt:  types.Rendered{Rendered: ""},
	}

	article, err := PostToArticle(post)
	if err != nil {
		t.Fatalf("エラーが発生: %v", err)
	}

	if article.FrontMatter.Excerpt != "" {
		t.Errorf("Excerpt = %q, want empty", article.FrontMatter.Excerpt)
	}
}

// --- TestPageToArticle ---

func TestPageToArticle_正常な固定ページ(t *testing.T) {
	page := &types.Page{
		ID:        10,
		Slug:      "about",
		Status:    "publish",
		Date:      types.WPTime{Time: time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)},
		Modified:  types.WPTime{Time: time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC)},
		Title:     types.Rendered{Rendered: "会社概要"},
		Content:   types.Rendered{Rendered: "<h2>概要</h2><p>テスト</p>"},
		Excerpt:   types.Rendered{Rendered: "<p>概要ページ</p>"},
		Parent:    0,
		MenuOrder: 1,
	}

	article, err := PageToArticle(page)
	if err != nil {
		t.Fatalf("エラーが発生: %v", err)
	}

	if article.FrontMatter.ID != 10 {
		t.Errorf("ID = %d, want 10", article.FrontMatter.ID)
	}
	if article.FrontMatter.Title != "会社概要" {
		t.Errorf("Title = %q, want %q", article.FrontMatter.Title, "会社概要")
	}
	if article.FrontMatter.MenuOrder != 1 {
		t.Errorf("MenuOrder = %d, want 1", article.FrontMatter.MenuOrder)
	}
	if !strings.Contains(article.FrontMatter.Excerpt, "概要ページ") {
		t.Errorf("Excerpt に概要ページが含まれていない: %q", article.FrontMatter.Excerpt)
	}
}

