package types

import (
	"strings"
	"time"
)

// WPTime はWordPress REST APIの日付形式をパースするカスタム型
type WPTime struct {
	time.Time
}

// UnmarshalJSON はWordPressの日付形式をパースする
func (t *WPTime) UnmarshalJSON(data []byte) error {
	// 引用符を除去
	s := strings.Trim(string(data), `"`)
	if s == "" || s == "null" {
		return nil
	}

	// WordPress REST APIの日付形式: "2023-01-05T19:30:00"
	parsed, err := time.Parse("2006-01-02T15:04:05", s)
	if err != nil {
		// RFC3339形式も試す
		parsed, err = time.Parse(time.RFC3339, s)
		if err != nil {
			return err
		}
	}
	t.Time = parsed
	return nil
}

// Post はWordPress投稿を表す
type Post struct {
	ID            int      `json:"id"`
	Date          WPTime   `json:"date"`
	DateGMT       WPTime   `json:"date_gmt"`
	Modified      WPTime   `json:"modified"`
	ModifiedGMT   WPTime   `json:"modified_gmt"`
	Slug          string   `json:"slug"`
	Status        string   `json:"status"`
	Title         Rendered `json:"title"`
	Content       Rendered `json:"content"`
	Excerpt       Rendered `json:"excerpt"`
	Author        int      `json:"author"`
	FeaturedMedia int      `json:"featured_media"`
	Categories    []int    `json:"categories"`
	Tags          []int    `json:"tags"`
	Link          string   `json:"link"`
}

// Page はWordPress固定ページを表す
type Page struct {
	ID          int      `json:"id"`
	Date        WPTime   `json:"date"`
	DateGMT     WPTime   `json:"date_gmt"`
	Modified    WPTime   `json:"modified"`
	ModifiedGMT WPTime   `json:"modified_gmt"`
	Slug        string   `json:"slug"`
	Status      string   `json:"status"`
	Title       Rendered `json:"title"`
	Content     Rendered `json:"content"`
	Excerpt     Rendered `json:"excerpt"`
	Author      int      `json:"author"`
	Parent      int      `json:"parent"`
	MenuOrder   int      `json:"menu_order"`
	Link        string   `json:"link"`
}

// Rendered はWordPressのレンダリング済みコンテンツを表す
type Rendered struct {
	Rendered string `json:"rendered"`
	Raw      string `json:"raw,omitempty"`
}

// Category はWordPressカテゴリを表す
type Category struct {
	ID          int    `json:"id"`
	Count       int    `json:"count"`
	Description string `json:"description"`
	Link        string `json:"link"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Parent      int    `json:"parent"`
}

// Tag はWordPressタグを表す
type Tag struct {
	ID          int    `json:"id"`
	Count       int    `json:"count"`
	Description string `json:"description"`
	Link        string `json:"link"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
}

// Media はWordPressメディアを表す
type Media struct {
	ID        int      `json:"id"`
	Date      WPTime   `json:"date"`
	Slug      string   `json:"slug"`
	Status    string   `json:"status"`
	Title     Rendered `json:"title"`
	Author    int      `json:"author"`
	MediaType string   `json:"media_type"`
	MimeType  string   `json:"mime_type"`
	SourceURL string   `json:"source_url"`
	Link      string   `json:"link"`
}

// CreatePostRequest は投稿作成リクエストを表す
type CreatePostRequest struct {
	Title         string `json:"title"`
	Content       string `json:"content"`
	Status        string `json:"status,omitempty"`
	Slug          string `json:"slug,omitempty"`
	Excerpt       string `json:"excerpt,omitempty"`
	Categories    []int  `json:"categories,omitempty"`
	Tags          []int  `json:"tags,omitempty"`
	FeaturedMedia int    `json:"featured_media,omitempty"`
}

// UpdatePostRequest は投稿更新リクエストを表す
type UpdatePostRequest struct {
	Title         string `json:"title,omitempty"`
	Content       string `json:"content,omitempty"`
	Status        string `json:"status,omitempty"`
	Slug          string `json:"slug,omitempty"`
	Excerpt       string `json:"excerpt,omitempty"`
	Categories    []int  `json:"categories,omitempty"`
	Tags          []int  `json:"tags,omitempty"`
	FeaturedMedia int    `json:"featured_media,omitempty"`
}

// CreatePageRequest は固定ページ作成リクエストを表す
type CreatePageRequest struct {
	Title     string `json:"title"`
	Content   string `json:"content"`
	Status    string `json:"status,omitempty"`
	Slug      string `json:"slug,omitempty"`
	Excerpt   string `json:"excerpt,omitempty"`
	Parent    int    `json:"parent,omitempty"`
	MenuOrder int    `json:"menu_order,omitempty"`
}

// UpdatePageRequest は固定ページ更新リクエストを表す
type UpdatePageRequest struct {
	Title     string `json:"title,omitempty"`
	Content   string `json:"content,omitempty"`
	Status    string `json:"status,omitempty"`
	Slug      string `json:"slug,omitempty"`
	Excerpt   string `json:"excerpt,omitempty"`
	Parent    int    `json:"parent,omitempty"`
	MenuOrder int    `json:"menu_order,omitempty"`
}

// CreateCategoryRequest はカテゴリ作成リクエストを表す
type CreateCategoryRequest struct {
	Name        string `json:"name"`
	Slug        string `json:"slug,omitempty"`
	Description string `json:"description,omitempty"`
	Parent      int    `json:"parent,omitempty"`
}

// UpdateCategoryRequest はカテゴリ更新リクエストを表す
type UpdateCategoryRequest struct {
	Name        string `json:"name,omitempty"`
	Slug        string `json:"slug,omitempty"`
	Description string `json:"description,omitempty"`
	Parent      int    `json:"parent,omitempty"`
}

// FrontMatter はMarkdownファイルのフロントマターを表す
type FrontMatter struct {
	ID            int    `yaml:"id,omitempty"`
	Title         string `yaml:"title"`
	Slug          string `yaml:"slug,omitempty"`
	Status        string `yaml:"status,omitempty"`
	Excerpt       string `yaml:"excerpt,omitempty"`
	Categories    []int  `yaml:"categories,omitempty"`
	Tags          []int  `yaml:"tags,omitempty"`
	FeaturedMedia int    `yaml:"featured_media,omitempty"`
	Date          string `yaml:"date,omitempty"`
	Modified      string `yaml:"modified,omitempty"`
	Parent        int    `yaml:"parent,omitempty"`
	MenuOrder     int    `yaml:"menu_order,omitempty"`
}

// Article はMarkdown記事を表す
type Article struct {
	FrontMatter FrontMatter
	Content     string
	FilePath    string
}
