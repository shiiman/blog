package converter

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/shiimanblog/wp-cli/internal/types"
	"gopkg.in/yaml.v3"
)

// HTMLToMarkdown はHTMLをMarkdownに変換する
func HTMLToMarkdown(htmlContent string) (string, error) {
	converter := md.NewConverter("", true, nil)
	mdContent, err := converter.ConvertString(htmlContent)
	if err != nil {
		return "", fmt.Errorf("HTML→Markdown変換に失敗: %w", err)
	}
	return mdContent, nil
}

// MarkdownToHTML はMarkdownをHTMLに変換する
func MarkdownToHTML(mdContent string) string {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse([]byte(mdContent))

	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return string(markdown.Render(doc, renderer))
}

// ParseArticle はMarkdownファイルを解析してArticleを返す
func ParseArticle(filePath string) (*types.Article, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("ファイルの読み込みに失敗: %w", err)
	}

	fm, body, err := parseFrontMatter(string(content))
	if err != nil {
		return nil, err
	}

	return &types.Article{
		FrontMatter: *fm,
		Content:     body,
		FilePath:    filePath,
	}, nil
}

// parseFrontMatter はMarkdownからFront Matterと本文を分離する
func parseFrontMatter(content string) (*types.FrontMatter, string, error) {
	// Front Matterの正規表現
	re := regexp.MustCompile(`(?s)^---\n(.+?)\n---\n(.*)$`)
	matches := re.FindStringSubmatch(content)

	if len(matches) != 3 {
		return nil, content, nil // Front Matterがない場合
	}

	var fm types.FrontMatter
	if err := yaml.Unmarshal([]byte(matches[1]), &fm); err != nil {
		return nil, "", fmt.Errorf("Front Matterのパースに失敗: %w", err)
	}

	return &fm, strings.TrimSpace(matches[2]), nil
}

// GenerateArticleFile はArticleからMarkdownファイルの内容を生成する
func GenerateArticleFile(article *types.Article) (string, error) {
	fmBytes, err := yaml.Marshal(article.FrontMatter)
	if err != nil {
		return "", fmt.Errorf("Front Matterの生成に失敗: %w", err)
	}

	return fmt.Sprintf("---\n%s---\n\n%s", string(fmBytes), article.Content), nil
}

// PostToArticle はWordPress投稿をArticleに変換する
func PostToArticle(post *types.Post) (*types.Article, error) {
	// HTMLをMarkdownに変換
	mdContent, err := HTMLToMarkdown(post.Content.Rendered)
	if err != nil {
		return nil, err
	}

	fm := types.FrontMatter{
		ID:            post.ID,
		Title:         post.Title.Rendered,
		Slug:          post.Slug,
		Status:        post.Status,
		Categories:    post.Categories,
		Tags:          post.Tags,
		FeaturedMedia: post.FeaturedMedia,
		Date:          post.Date.Time.Format("2006-01-02T15:04:05"),
		Modified:      post.Modified.Time.Format("2006-01-02T15:04:05"),
	}

	// Excerptがある場合はHTMLタグを除去
	if post.Excerpt.Rendered != "" {
		excerpt, _ := HTMLToMarkdown(post.Excerpt.Rendered)
		fm.Excerpt = strings.TrimSpace(excerpt)
	}

	return &types.Article{
		FrontMatter: fm,
		Content:     mdContent,
	}, nil
}

// PageToArticle はWordPress固定ページをArticleに変換する
func PageToArticle(page *types.Page) (*types.Article, error) {
	// HTMLをMarkdownに変換
	mdContent, err := HTMLToMarkdown(page.Content.Rendered)
	if err != nil {
		return nil, err
	}

	fm := types.FrontMatter{
		ID:        page.ID,
		Title:     page.Title.Rendered,
		Slug:      page.Slug,
		Status:    page.Status,
		Parent:    page.Parent,
		MenuOrder: page.MenuOrder,
		Date:      page.Date.Time.Format("2006-01-02T15:04:05"),
		Modified:  page.Modified.Time.Format("2006-01-02T15:04:05"),
	}

	// Excerptがある場合はHTMLタグを除去
	if page.Excerpt.Rendered != "" {
		excerpt, _ := HTMLToMarkdown(page.Excerpt.Rendered)
		fm.Excerpt = strings.TrimSpace(excerpt)
	}

	return &types.Article{
		FrontMatter: fm,
		Content:     mdContent,
	}, nil
}
