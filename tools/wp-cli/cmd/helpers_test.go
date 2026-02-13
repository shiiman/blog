package cmd

import (
	"testing"
)

func TestDetermineStatus(t *testing.T) {
	tests := []struct {
		name              string
		publishFlag       bool
		frontMatterStatus string
		want              string
	}{
		{
			name:              "デフォルトはdraft",
			publishFlag:       false,
			frontMatterStatus: "",
			want:              "draft",
		},
		{
			name:              "publishフラグが優先",
			publishFlag:       true,
			frontMatterStatus: "draft",
			want:              "publish",
		},
		{
			name:              "publishフラグなしでフロントマターステータスを使用",
			publishFlag:       false,
			frontMatterStatus: "pending",
			want:              "pending",
		},
		{
			name:              "publishフラグありでフロントマターステータスは無視",
			publishFlag:       true,
			frontMatterStatus: "pending",
			want:              "publish",
		},
		{
			name:              "フロントマターがpublishの場合",
			publishFlag:       false,
			frontMatterStatus: "publish",
			want:              "publish",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := determineStatus(tt.publishFlag, tt.frontMatterStatus)
			if got != tt.want {
				t.Errorf("determineStatus(%v, %q) = %q, want %q", tt.publishFlag, tt.frontMatterStatus, got, tt.want)
			}
		})
	}
}

func TestTruncate(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		maxLen int
		want   string
	}{
		{
			name:   "短い文字列はそのまま",
			input:  "hello",
			maxLen: 10,
			want:   "hello",
		},
		{
			name:   "ちょうどの長さ",
			input:  "hello",
			maxLen: 5,
			want:   "hello",
		},
		{
			name:   "長い文字列は切り詰め",
			input:  "hello world",
			maxLen: 8,
			want:   "hello...",
		},
		{
			name:   "日本語の切り詰め",
			input:  "これはテストの文字列です",
			maxLen: 9,
			want:   "これはテスト...",
		},
		{
			name:   "空文字列",
			input:  "",
			maxLen: 10,
			want:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := truncate(tt.input, tt.maxLen)
			if got != tt.want {
				t.Errorf("truncate(%q, %d) = %q, want %q", tt.input, tt.maxLen, got, tt.want)
			}
		})
	}
}

func TestFormatStatus(t *testing.T) {
	tests := []struct {
		name   string
		status string
	}{
		{name: "publish", status: "publish"},
		{name: "draft", status: "draft"},
		{name: "pending", status: "pending"},
		{name: "private", status: "private"},
		{name: "unknown", status: "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatStatus(tt.status)
			if got == "" {
				t.Errorf("formatStatus(%q) は空であるべきでない", tt.status)
			}
			// unknownステータスはそのまま返される
			if tt.status == "unknown" && got != "unknown" {
				t.Errorf("formatStatus(%q) = %q, want %q", tt.status, got, "unknown")
			}
		})
	}
}
