package wp

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/shiimanblog/wp-cli/internal/config"
	"github.com/shiimanblog/wp-cli/internal/types"
)

// テスト用のクライアントとモックサーバーを作成するヘルパー
func setupTestClient(t *testing.T, handler http.HandlerFunc) (*Client, *httptest.Server) {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	cfg := &config.Config{
		SiteURL:     server.URL,
		Username:    "testuser",
		AppPassword: "test pass word",
	}
	client := NewClient(cfg)
	// baseURL をテストサーバーのURLに上書き（/wp-json/wp/v2 サフィックスなし）
	client.baseURL = server.URL

	return client, server
}

func TestDoRawRequest_正常レスポンス(t *testing.T) {
	client, _ := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	ctx := context.Background()
	body, err := client.doRawRequest(ctx, "GET", "/test", nil, "application/json", nil)
	if err != nil {
		t.Fatalf("予期しないエラー: %v", err)
	}

	if string(body) != `{"status":"ok"}` {
		t.Errorf("レスポンスボディ = %q, want %q", string(body), `{"status":"ok"}`)
	}
}

func TestDoRawRequest_エラーレスポンス(t *testing.T) {
	client, _ := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"code":"not_found"}`))
	})

	ctx := context.Background()
	_, err := client.doRawRequest(ctx, "GET", "/test", nil, "application/json", nil)
	if err == nil {
		t.Fatal("エラーが発生するべき")
	}

	if !strings.Contains(err.Error(), "status 404") {
		t.Errorf("エラーメッセージにステータスコードが含まれるべき: %v", err)
	}
}

func TestDoRawRequest_コンテキストキャンセル(t *testing.T) {
	client, _ := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // キャンセル済みのコンテキスト

	_, err := client.doRawRequest(ctx, "GET", "/test", nil, "application/json", nil)
	if err == nil {
		t.Fatal("キャンセル済みコンテキストでエラーが発生するべき")
	}
}

func TestGetPosts_正常(t *testing.T) {
	client, _ := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("メソッド = %s, want GET", r.Method)
		}
		if !strings.Contains(r.URL.String(), "/posts") {
			t.Errorf("URL = %s, /posts を含むべき", r.URL.String())
		}

		posts := []types.Post{
			{ID: 1, Slug: "test-post", Status: "publish"},
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(posts)
	})

	ctx := context.Background()
	posts, err := client.GetPosts(ctx, 1, 10, "")
	if err != nil {
		t.Fatalf("予期しないエラー: %v", err)
	}

	if len(posts) != 1 {
		t.Fatalf("投稿数 = %d, want 1", len(posts))
	}
	if posts[0].ID != 1 {
		t.Errorf("投稿ID = %d, want 1", posts[0].ID)
	}
}

func TestGetPost_正常(t *testing.T) {
	client, _ := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		post := types.Post{ID: 42, Slug: "my-post", Status: "draft"}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(post)
	})

	ctx := context.Background()
	post, err := client.GetPost(ctx, 42)
	if err != nil {
		t.Fatalf("予期しないエラー: %v", err)
	}

	if post.ID != 42 {
		t.Errorf("投稿ID = %d, want 42", post.ID)
	}
}

func TestCreatePost_正常(t *testing.T) {
	client, _ := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("メソッド = %s, want POST", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Content-Type = %s, want application/json", r.Header.Get("Content-Type"))
		}

		post := types.Post{ID: 100, Slug: "new-post", Status: "draft"}
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(post)
	})

	ctx := context.Background()
	req := &types.CreatePostRequest{
		Title:   "テスト投稿",
		Content: "<p>テスト</p>",
		Status:  "draft",
	}
	post, err := client.CreatePost(ctx, req)
	if err != nil {
		t.Fatalf("予期しないエラー: %v", err)
	}

	if post.ID != 100 {
		t.Errorf("投稿ID = %d, want 100", post.ID)
	}
}

func TestUploadMedia_正常(t *testing.T) {
	client, _ := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		media := types.Media{ID: 200, SourceURL: "https://example.com/image.png"}
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(media)
	})

	ctx := context.Background()
	media, err := client.UploadMedia(ctx, "test.png", []byte("fake-image-data"), "image/png")
	if err != nil {
		t.Fatalf("予期しないエラー: %v", err)
	}

	if media.ID != 200 {
		t.Errorf("メディアID = %d, want 200", media.ID)
	}
}

func TestUploadMedia_ヘッダー確認(t *testing.T) {
	client, _ := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		// Content-Type の確認
		if r.Header.Get("Content-Type") != "image/png" {
			t.Errorf("Content-Type = %s, want image/png", r.Header.Get("Content-Type"))
		}

		// Content-Disposition の確認
		disposition := r.Header.Get("Content-Disposition")
		if !strings.Contains(disposition, "test.png") {
			t.Errorf("Content-Disposition = %s, test.png を含むべき", disposition)
		}

		// Authorization ヘッダーの確認
		auth := r.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Basic ") {
			t.Errorf("Authorization = %s, Basic で始まるべき", auth)
		}

		media := types.Media{ID: 1}
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(media)
	})

	ctx := context.Background()
	_, err := client.UploadMedia(ctx, "test.png", []byte("data"), "image/png")
	if err != nil {
		t.Fatalf("予期しないエラー: %v", err)
	}
}

func TestGetAuthHeader_パスワードスペース除去(t *testing.T) {
	cfg := &config.Config{
		SiteURL:     "https://example.com",
		Username:    "admin",
		AppPassword: "abcd efgh ijkl mnop",
	}
	client := NewClient(cfg)

	header := client.getAuthHeader()
	if !strings.HasPrefix(header, "Basic ") {
		t.Fatalf("ヘッダー = %s, Basic で始まるべき", header)
	}

	// Base64デコードして確認
	encoded := strings.TrimPrefix(header, "Basic ")
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		t.Fatalf("Base64デコードに失敗: %v", err)
	}

	expected := "admin:abcdefghijklmnop"
	if string(decoded) != expected {
		t.Errorf("デコード結果 = %q, want %q", string(decoded), expected)
	}
}
