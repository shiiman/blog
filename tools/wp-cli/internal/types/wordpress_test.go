package types

import (
	"encoding/json"
	"testing"
	"time"
)

func TestWPTimeUnmarshalJSON_WordPress形式(t *testing.T) {
	// WordPress REST APIの標準形式: "2023-01-05T19:30:00"
	input := `"2023-01-05T19:30:00"`
	var wt WPTime

	err := wt.UnmarshalJSON([]byte(input))
	if err != nil {
		t.Fatalf("エラーが発生: %v", err)
	}

	expected := time.Date(2023, 1, 5, 19, 30, 0, 0, time.UTC)
	if !wt.Equal(expected) {
		t.Errorf("Time = %v, want %v", wt.Time, expected)
	}
}

func TestWPTimeUnmarshalJSON_RFC3339形式(t *testing.T) {
	// RFC3339形式: "2023-01-05T19:30:00Z"
	input := `"2023-01-05T19:30:00Z"`
	var wt WPTime

	err := wt.UnmarshalJSON([]byte(input))
	if err != nil {
		t.Fatalf("エラーが発生: %v", err)
	}

	expected := time.Date(2023, 1, 5, 19, 30, 0, 0, time.UTC)
	if !wt.Equal(expected) {
		t.Errorf("Time = %v, want %v", wt.Time, expected)
	}
}

func TestWPTimeUnmarshalJSON_RFC3339タイムゾーン付き(t *testing.T) {
	input := `"2023-06-15T10:00:00+09:00"`
	var wt WPTime

	err := wt.UnmarshalJSON([]byte(input))
	if err != nil {
		t.Fatalf("エラーが発生: %v", err)
	}

	// +09:00 で解析されるはず
	if wt.IsZero() {
		t.Error("Time がゼロ値のまま")
	}
	if wt.Hour() != 10 {
		t.Errorf("Hour = %d, want 10", wt.Hour())
	}
}

func TestWPTimeUnmarshalJSON_空文字列(t *testing.T) {
	input := `""`
	var wt WPTime

	err := wt.UnmarshalJSON([]byte(input))
	if err != nil {
		t.Fatalf("空文字列でエラーが発生: %v", err)
	}

	if !wt.IsZero() {
		t.Errorf("空文字列の場合、Time はゼロ値であるべき: %v", wt.Time)
	}
}

func TestWPTimeUnmarshalJSON_null(t *testing.T) {
	input := `"null"`
	var wt WPTime

	err := wt.UnmarshalJSON([]byte(input))
	if err != nil {
		t.Fatalf("nullでエラーが発生: %v", err)
	}

	if !wt.IsZero() {
		t.Errorf("nullの場合、Time はゼロ値であるべき: %v", wt.Time)
	}
}

func TestWPTimeUnmarshalJSON_不正な形式(t *testing.T) {
	input := `"not-a-date"`
	var wt WPTime

	err := wt.UnmarshalJSON([]byte(input))
	if err == nil {
		t.Error("不正な日付形式でエラーが発生するべき")
	}
}

func TestWPTimeUnmarshalJSON_JSON構造体内での使用(t *testing.T) {
	// 実際のWordPress APIレスポンスに近い構造でテスト
	jsonData := `{
		"id": 1,
		"date": "2024-03-15T09:00:00",
		"slug": "test-post",
		"title": {"rendered": "テスト"},
		"content": {"rendered": "<p>本文</p>"},
		"excerpt": {"rendered": ""},
		"categories": [1],
		"tags": [],
		"status": "publish"
	}`

	var post Post
	err := json.Unmarshal([]byte(jsonData), &post)
	if err != nil {
		t.Fatalf("JSONパースエラー: %v", err)
	}

	if post.ID != 1 {
		t.Errorf("ID = %d, want 1", post.ID)
	}
	if post.Date.Year() != 2024 || post.Date.Month() != 3 || post.Date.Day() != 15 {
		t.Errorf("Date = %v, want 2024-03-15", post.Date.Time)
	}
	if post.Title.Rendered != "テスト" {
		t.Errorf("Title = %q, want %q", post.Title.Rendered, "テスト")
	}
	if post.Status != "publish" {
		t.Errorf("Status = %q, want %q", post.Status, "publish")
	}
}

func TestWPTimeUnmarshalJSON_Page構造体(t *testing.T) {
	jsonData := `{
		"id": 10,
		"date": "2024-01-01T00:00:00",
		"slug": "about",
		"status": "publish",
		"title": {"rendered": "About"},
		"content": {"rendered": "<p>About page</p>"},
		"excerpt": {"rendered": ""},
		"parent": 0,
		"menu_order": 1
	}`

	var page Page
	err := json.Unmarshal([]byte(jsonData), &page)
	if err != nil {
		t.Fatalf("JSONパースエラー: %v", err)
	}

	if page.ID != 10 {
		t.Errorf("ID = %d, want 10", page.ID)
	}
	if page.MenuOrder != 1 {
		t.Errorf("MenuOrder = %d, want 1", page.MenuOrder)
	}
}
