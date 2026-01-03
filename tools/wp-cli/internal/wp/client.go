package wp

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/shiimanblog/wp-cli/internal/config"
	"github.com/shiimanblog/wp-cli/internal/types"
)

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
		httpClient: &http.Client{},
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

// doRequest はHTTPリクエストを実行する
func (c *Client) doRequest(method, endpoint string, body interface{}) ([]byte, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("リクエストボディのJSON化に失敗: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	url := c.baseURL + endpoint
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("リクエストの作成に失敗: %w", err)
	}

	req.Header.Set("Authorization", c.getAuthHeader())
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("リクエストの実行に失敗: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("レスポンスの読み取りに失敗: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("APIエラー (status %d): %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// GetPosts は投稿一覧を取得する
func (c *Client) GetPosts(page, perPage int, status string) ([]types.Post, error) {
	endpoint := fmt.Sprintf("/posts?page=%d&per_page=%d&_embed", page, perPage)
	if status != "" {
		endpoint += "&status=" + status
	}

	body, err := c.doRequest("GET", endpoint, nil)
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
func (c *Client) GetPost(id int) (*types.Post, error) {
	endpoint := fmt.Sprintf("/posts/%d?context=edit", id)

	body, err := c.doRequest("GET", endpoint, nil)
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
func (c *Client) CreatePost(req *types.CreatePostRequest) (*types.Post, error) {
	body, err := c.doRequest("POST", "/posts", req)
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
func (c *Client) UpdatePost(id int, req *types.UpdatePostRequest) (*types.Post, error) {
	endpoint := fmt.Sprintf("/posts/%d", id)

	body, err := c.doRequest("POST", endpoint, req)
	if err != nil {
		return nil, err
	}

	var post types.Post
	if err := json.Unmarshal(body, &post); err != nil {
		return nil, fmt.Errorf("投稿のパースに失敗: %w", err)
	}

	return &post, nil
}

// GetPages は固定ページ一覧を取得する
func (c *Client) GetPages(page, perPage int, status string) ([]types.Page, error) {
	endpoint := fmt.Sprintf("/pages?page=%d&per_page=%d&_embed", page, perPage)
	if status != "" {
		endpoint += "&status=" + status
	}

	body, err := c.doRequest("GET", endpoint, nil)
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
func (c *Client) GetPage(id int) (*types.Page, error) {
	endpoint := fmt.Sprintf("/pages/%d?context=edit", id)

	body, err := c.doRequest("GET", endpoint, nil)
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
func (c *Client) CreatePage(req *types.CreatePageRequest) (*types.Page, error) {
	body, err := c.doRequest("POST", "/pages", req)
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
func (c *Client) UpdatePage(id int, req *types.UpdatePageRequest) (*types.Page, error) {
	endpoint := fmt.Sprintf("/pages/%d", id)

	body, err := c.doRequest("POST", endpoint, req)
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
func (c *Client) GetCategories() ([]types.Category, error) {
	body, err := c.doRequest("GET", "/categories?per_page=100", nil)
	if err != nil {
		return nil, err
	}

	var categories []types.Category
	if err := json.Unmarshal(body, &categories); err != nil {
		return nil, fmt.Errorf("カテゴリ一覧のパースに失敗: %w", err)
	}

	return categories, nil
}

// GetTags はタグ一覧を取得する
func (c *Client) GetTags() ([]types.Tag, error) {
	body, err := c.doRequest("GET", "/tags?per_page=100", nil)
	if err != nil {
		return nil, err
	}

	var tags []types.Tag
	if err := json.Unmarshal(body, &tags); err != nil {
		return nil, fmt.Errorf("タグ一覧のパースに失敗: %w", err)
	}

	return tags, nil
}
