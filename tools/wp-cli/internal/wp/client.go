package wp

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/shiimanblog/wp-cli/internal/config"
	"github.com/shiimanblog/wp-cli/internal/types"
)

// maxResponseSize はレスポンスボディの最大サイズ（50MB）
const maxResponseSize = 50 * 1024 * 1024

// Client はWordPress REST APIクライアント
type Client struct {
	config     *config.Config
	httpClient *http.Client
	baseURL    string
}

// NewClient は新しいWordPressクライアントを作成する
func NewClient(cfg *config.Config) *Client {
	return &Client{
		config:     cfg,
		httpClient: &http.Client{Timeout: 30 * time.Second},
		baseURL:    strings.TrimSuffix(cfg.SiteURL, "/") + "/wp-json/wp/v2",
	}
}

// getAuthHeader は認証ヘッダーを生成する
func (c *Client) getAuthHeader() string {
	// スペースを除去したパスワードで認証
	password := strings.ReplaceAll(c.config.AppPassword, " ", "")
	credentials := c.config.Username + ":" + password
	encoded := base64.StdEncoding.EncodeToString([]byte(credentials))
	return "Basic " + encoded
}

// doRawRequest はHTTPリクエストを実行する共通処理
// 認証ヘッダー付与、レスポンス読み取り（サイズ制限付き）、エラーチェックを行う
func (c *Client) doRawRequest(ctx context.Context, method, endpoint string, body io.Reader, contentType string, extraHeaders map[string]string) ([]byte, error) {
	reqURL := c.baseURL + endpoint
	req, err := http.NewRequestWithContext(ctx, method, reqURL, body)
	if err != nil {
		return nil, fmt.Errorf("リクエストの作成に失敗: %w", err)
	}

	req.Header.Set("Authorization", c.getAuthHeader())
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("User-Agent", "wp-cli/1.0.0")

	for k, v := range extraHeaders {
		req.Header.Set(k, v)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("リクエストの実行に失敗: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	respBody, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseSize))
	if err != nil {
		return nil, fmt.Errorf("レスポンスの読み取りに失敗: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("APIエラー (status %d): %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// doRequest はJSON形式のHTTPリクエストを実行する
func (c *Client) doRequest(ctx context.Context, method, endpoint string, body interface{}) ([]byte, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("リクエストボディのJSON化に失敗: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	return c.doRawRequest(ctx, method, endpoint, reqBody, "application/json", nil)
}

// GetPosts は投稿一覧を取得する
func (c *Client) GetPosts(ctx context.Context, page, perPage int, status string) ([]types.Post, error) {
	endpoint := fmt.Sprintf("/posts?page=%d&per_page=%d&_embed", page, perPage)
	if status != "" {
		endpoint += "&status=" + url.QueryEscape(status)
	}

	body, err := c.doRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var posts []types.Post
	if err := json.Unmarshal(body, &posts); err != nil {
		return nil, fmt.Errorf("投稿一覧のパースに失敗: %w", err)
	}

	return posts, nil
}

// GetPost は特定の投稿を取得する
func (c *Client) GetPost(ctx context.Context, id int) (*types.Post, error) {
	endpoint := fmt.Sprintf("/posts/%d?context=edit", id)

	body, err := c.doRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var post types.Post
	if err := json.Unmarshal(body, &post); err != nil {
		return nil, fmt.Errorf("投稿のパースに失敗: %w", err)
	}

	return &post, nil
}

// CreatePost は新しい投稿を作成する
func (c *Client) CreatePost(ctx context.Context, req *types.CreatePostRequest) (*types.Post, error) {
	body, err := c.doRequest(ctx, "POST", "/posts", req)
	if err != nil {
		return nil, err
	}

	var post types.Post
	if err := json.Unmarshal(body, &post); err != nil {
		return nil, fmt.Errorf("投稿のパースに失敗: %w", err)
	}

	return &post, nil
}

// UpdatePost は既存の投稿を更新する
func (c *Client) UpdatePost(ctx context.Context, id int, req *types.UpdatePostRequest) (*types.Post, error) {
	endpoint := fmt.Sprintf("/posts/%d", id)

	body, err := c.doRequest(ctx, "POST", endpoint, req)
	if err != nil {
		return nil, err
	}

	var post types.Post
	if err := json.Unmarshal(body, &post); err != nil {
		return nil, fmt.Errorf("投稿のパースに失敗: %w", err)
	}

	return &post, nil
}

// DeletePost は投稿を削除する
func (c *Client) DeletePost(ctx context.Context, id int, force bool) error {
	endpoint := fmt.Sprintf("/posts/%d", id)
	if force {
		endpoint += "?force=true"
	}

	_, err := c.doRequest(ctx, "DELETE", endpoint, nil)
	return err
}

// GetPages は固定ページ一覧を取得する
func (c *Client) GetPages(ctx context.Context, page, perPage int, status string) ([]types.Page, error) {
	endpoint := fmt.Sprintf("/pages?page=%d&per_page=%d&_embed", page, perPage)
	if status != "" {
		endpoint += "&status=" + url.QueryEscape(status)
	}

	body, err := c.doRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var pages []types.Page
	if err := json.Unmarshal(body, &pages); err != nil {
		return nil, fmt.Errorf("固定ページ一覧のパースに失敗: %w", err)
	}

	return pages, nil
}

// GetPage は特定の固定ページを取得する
func (c *Client) GetPage(ctx context.Context, id int) (*types.Page, error) {
	endpoint := fmt.Sprintf("/pages/%d?context=edit", id)

	body, err := c.doRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var page types.Page
	if err := json.Unmarshal(body, &page); err != nil {
		return nil, fmt.Errorf("固定ページのパースに失敗: %w", err)
	}

	return &page, nil
}

// CreatePage は新しい固定ページを作成する
func (c *Client) CreatePage(ctx context.Context, req *types.CreatePageRequest) (*types.Page, error) {
	body, err := c.doRequest(ctx, "POST", "/pages", req)
	if err != nil {
		return nil, err
	}

	var page types.Page
	if err := json.Unmarshal(body, &page); err != nil {
		return nil, fmt.Errorf("固定ページのパースに失敗: %w", err)
	}

	return &page, nil
}

// UpdatePage は既存の固定ページを更新する
func (c *Client) UpdatePage(ctx context.Context, id int, req *types.UpdatePageRequest) (*types.Page, error) {
	endpoint := fmt.Sprintf("/pages/%d", id)

	body, err := c.doRequest(ctx, "POST", endpoint, req)
	if err != nil {
		return nil, err
	}

	var page types.Page
	if err := json.Unmarshal(body, &page); err != nil {
		return nil, fmt.Errorf("固定ページのパースに失敗: %w", err)
	}

	return &page, nil
}

// GetCategories はカテゴリ一覧を取得する
func (c *Client) GetCategories(ctx context.Context) ([]types.Category, error) {
	body, err := c.doRequest(ctx, "GET", "/categories?per_page=100", nil)
	if err != nil {
		return nil, err
	}

	var categories []types.Category
	if err := json.Unmarshal(body, &categories); err != nil {
		return nil, fmt.Errorf("カテゴリ一覧のパースに失敗: %w", err)
	}

	return categories, nil
}

// CreateCategory は新しいカテゴリを作成する
func (c *Client) CreateCategory(ctx context.Context, req *types.CreateCategoryRequest) (*types.Category, error) {
	body, err := c.doRequest(ctx, "POST", "/categories", req)
	if err != nil {
		return nil, err
	}

	var category types.Category
	if err := json.Unmarshal(body, &category); err != nil {
		return nil, fmt.Errorf("カテゴリのパースに失敗: %w", err)
	}

	return &category, nil
}

// UpdateCategory は既存のカテゴリを更新する
func (c *Client) UpdateCategory(ctx context.Context, id int, req *types.UpdateCategoryRequest) (*types.Category, error) {
	endpoint := fmt.Sprintf("/categories/%d", id)

	body, err := c.doRequest(ctx, "POST", endpoint, req)
	if err != nil {
		return nil, err
	}

	var category types.Category
	if err := json.Unmarshal(body, &category); err != nil {
		return nil, fmt.Errorf("カテゴリのパースに失敗: %w", err)
	}

	return &category, nil
}

// GetTags はタグ一覧を取得する
func (c *Client) GetTags(ctx context.Context) ([]types.Tag, error) {
	body, err := c.doRequest(ctx, "GET", "/tags?per_page=100", nil)
	if err != nil {
		return nil, err
	}

	var tags []types.Tag
	if err := json.Unmarshal(body, &tags); err != nil {
		return nil, fmt.Errorf("タグ一覧のパースに失敗: %w", err)
	}

	return tags, nil
}

// UploadMedia はメディアをアップロードする
func (c *Client) UploadMedia(ctx context.Context, filename string, data []byte, mimeType string) (*types.Media, error) {
	extraHeaders := map[string]string{
		"Content-Disposition": fmt.Sprintf(`attachment; filename="%s"`, filename),
	}

	respBody, err := c.doRawRequest(ctx, "POST", "/media", bytes.NewReader(data), mimeType, extraHeaders)
	if err != nil {
		return nil, err
	}

	var media types.Media
	if err := json.Unmarshal(respBody, &media); err != nil {
		return nil, fmt.Errorf("メディアのパースに失敗: %w", err)
	}

	return &media, nil
}
