package cmd

import (
	"testing"
)

func TestDeterminePostStatus(t *testing.T) {
	tests := []struct {
		name      string
		draftFlag bool
		want      string
	}{
		{
			name:      "デフォルトはpublish",
			draftFlag: false,
			want:      "publish",
		},
		{
			name:      "draftフラグ指定時はdraft",
			draftFlag: true,
			want:      "draft",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := determinePostStatus(tt.draftFlag)
			if got != tt.want {
				t.Errorf("determinePostStatus(%v) = %q, want %q", tt.draftFlag, got, tt.want)
			}
		})
	}
}

func TestDeterminePageStatus(t *testing.T) {
	tests := []struct {
		name              string
		frontMatterStatus string
		want              string
		wantErr           bool
	}{
		{
			name:              "status未指定はdraft",
			frontMatterStatus: "",
			want:              "draft",
		},
		{
			name:              "status=publishを採用",
			frontMatterStatus: "publish",
			want:              "publish",
		},
		{
			name:              "status=pendingを採用",
			frontMatterStatus: "pending",
			want:              "pending",
		},
		{
			name:              "statusの前後空白と大文字を正規化",
			frontMatterStatus: " Publish ",
			want:              "publish",
		},
		{
			name:              "status不正値はエラー",
			frontMatterStatus: "publsih",
			wantErr:           true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := determinePageStatus(tt.frontMatterStatus)
			if tt.wantErr {
				if err == nil {
					t.Errorf("determinePageStatus(%q) はエラーになるべき", tt.frontMatterStatus)
				}
				return
			}
			if err != nil {
				t.Errorf("determinePageStatus(%q) が予期せずエラー: %v", tt.frontMatterStatus, err)
				return
			}
			if got != tt.want {
				t.Errorf("determinePageStatus(%q) = %q, want %q", tt.frontMatterStatus, got, tt.want)
			}
		})
	}
}

func TestNormalizeAndValidateStatus(t *testing.T) {
	tests := []struct {
		name       string
		status     string
		allowEmpty bool
		want       string
		wantErr    bool
	}{
		{
			name:       "draftは許可",
			status:     "draft",
			allowEmpty: false,
			want:       "draft",
		},
		{
			name:       "publishは許可",
			status:     "publish",
			allowEmpty: false,
			want:       "publish",
		},
		{
			name:       "pendingは許可",
			status:     "pending",
			allowEmpty: false,
			want:       "pending",
		},
		{
			name:       "privateは許可",
			status:     "private",
			allowEmpty: false,
			want:       "private",
		},
		{
			name:       "前後空白と大文字を正規化",
			status:     " Publish ",
			allowEmpty: false,
			want:       "publish",
		},
		{
			name:       "空文字はallowEmpty=trueで許可",
			status:     "",
			allowEmpty: true,
			want:       "",
		},
		{
			name:       "空文字はallowEmpty=falseでエラー",
			status:     "",
			allowEmpty: false,
			wantErr:    true,
		},
		{
			name:       "誤字はエラー",
			status:     "publsih",
			allowEmpty: false,
			wantErr:    true,
		},
		{
			name:       "未知statusはエラー",
			status:     "archived",
			allowEmpty: false,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := normalizeAndValidateStatus(tt.status, tt.allowEmpty)
			if tt.wantErr {
				if err == nil {
					t.Errorf("normalizeAndValidateStatus(%q, %v) はエラーになるべき", tt.status, tt.allowEmpty)
				}
				return
			}
			if err != nil {
				t.Errorf("normalizeAndValidateStatus(%q, %v) が予期せずエラー: %v", tt.status, tt.allowEmpty, err)
				return
			}
			if got != tt.want {
				t.Errorf("normalizeAndValidateStatus(%q, %v) = %q, want %q", tt.status, tt.allowEmpty, got, tt.want)
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
		{
			name:   "maxLen=3はそのまま切り詰め",
			input:  "hello",
			maxLen: 3,
			want:   "hel",
		},
		{
			name:   "maxLen=1はそのまま切り詰め",
			input:  "hello",
			maxLen: 1,
			want:   "h",
		},
		{
			name:   "maxLen=0は空文字列",
			input:  "hello",
			maxLen: 0,
			want:   "",
		},
		{
			name:   "maxLen=2はそのまま切り詰め",
			input:  "hello",
			maxLen: 2,
			want:   "he",
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

func TestSanitizeSlug(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{name: "通常のスラッグ", input: "my-article", want: "my-article"},
		{name: "スラッシュ除去", input: "../../etc/passwd", want: "etcpasswd"},
		{name: "バックスラッシュ除去", input: `my\article`, want: "myarticle"},
		{name: "ドット連続除去", input: "my..article", want: "myarticle"},
		{name: "複合パストラバーサル", input: "../../../tmp/evil", want: "tmpevil"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sanitizeSlug(tt.input)
			if got != tt.want {
				t.Errorf("sanitizeSlug(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
