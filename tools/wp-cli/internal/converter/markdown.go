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
	trimmed := strings.TrimSpace(mdContent)

	// HTMLタグで始まる場合は、HTMLブロックとMarkdown部分を分離して処理
	if strings.HasPrefix(trimmed, "<style") ||
		strings.HasPrefix(trimmed, "<script") ||
		strings.HasPrefix(trimmed, "<div") ||
		strings.HasPrefix(trimmed, "<!--") {
		return processHtmlAndMarkdownMixed(mdContent)
	}

	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse([]byte(mdContent))

	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	htmlContent := string(markdown.Render(doc, renderer))

	// Gutenbergブロック形式でラップ（script/styleタグを含むHTMLを保護）
	return wrapInGutenbergBlocks(htmlContent)
}

// processHtmlAndMarkdownMixed はHTML部分とMarkdown部分が混在するコンテンツを処理する
func processHtmlAndMarkdownMixed(content string) string {
	// </script>の後にMarkdownがあるかチェック
	scriptEndIdx := strings.LastIndex(content, "</script>")
	if scriptEndIdx != -1 {
		afterScript := content[scriptEndIdx+len("</script>"):]
		trimmedAfter := strings.TrimSpace(afterScript)

		// </script>の後にMarkdownコンテンツがある場合
		if len(trimmedAfter) > 0 && (strings.HasPrefix(trimmedAfter, "#") || strings.HasPrefix(trimmedAfter, "-") || strings.HasPrefix(trimmedAfter, "1.")) {
			htmlPart := content[:scriptEndIdx+len("</script>")]
			mdPart := afterScript

			// Markdown部分を変換
			extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
			p := parser.NewWithExtensions(extensions)
			doc := p.Parse([]byte(mdPart))

			htmlFlags := html.CommonFlags | html.HrefTargetBlank
			opts := html.RendererOptions{Flags: htmlFlags}
			renderer := html.NewRenderer(opts)
			convertedMd := string(markdown.Render(doc, renderer))

			// HTMLブロックでラップしたHTML部分 + 変換済みMarkdown部分
			return "<!-- wp:html -->\n" + htmlPart + "\n<!-- /wp:html -->\n\n" + convertedMd
		}
	}

	// Markdown部分がない場合はそのままHTMLブロックでラップ
	return wrapInGutenbergBlocks(content)
}

// wrapInGutenbergBlocks はHTMLコンテンツをGutenbergブロック形式でラップする
// script/styleタグを含む部分はwp:htmlブロックで保護し、
// それ以外の部分はwp:freeformブロックでラップする
func wrapInGutenbergBlocks(htmlContent string) string {
	// <style>タグと<script>タグを含むかチェック
	hasStyle := strings.Contains(htmlContent, "<style")
	hasScript := strings.Contains(htmlContent, "<script")

	if hasStyle || hasScript {
		// script/styleを含む場合はwp:htmlブロックでラップ
		return "<!-- wp:html -->\n" + htmlContent + "\n<!-- /wp:html -->"
	}

	// 通常のHTMLはそのまま返す（WordPressが適切に処理）
	return htmlContent
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
		return &types.FrontMatter{}, content, nil // Front Matterがない場合
	}

	var fm types.FrontMatter
	if err := yaml.Unmarshal([]byte(matches[1]), &fm); err != nil {
		return nil, "", fmt.Errorf("front matterのパースに失敗: %w", err)
	}

	return &fm, strings.TrimSpace(matches[2]), nil
}

// GenerateArticleFile はArticleからMarkdownファイルの内容を生成する
func GenerateArticleFile(article *types.Article) (string, error) {
	fmBytes, err := yaml.Marshal(article.FrontMatter)
	if err != nil {
		return "", fmt.Errorf("front matterの生成に失敗: %w", err)
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
		Date:          post.Date.Format("2006-01-02T15:04:05"),
		Modified:      post.Modified.Format("2006-01-02T15:04:05"),
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
		Date:      page.Date.Format("2006-01-02T15:04:05"),
		Modified:  page.Modified.Format("2006-01-02T15:04:05"),
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
