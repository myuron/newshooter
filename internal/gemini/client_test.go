package gemini

import (
	"context"
	"fmt"
	"os"
	"testing"
)

// TestSummarize_Integration は実際のGemini APIを使う統合テスト。
// GEMINI_API_KEY が未設定の場合はスキップする。
func TestSummarize_Integration(t *testing.T) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		t.Skip("GEMINI_API_KEY not set, skipping integration test")
	}

	changelog := `## v1.0.0
- Added new feature X
- Fixed bug Y
- Improved performance Z`

	summary, err := Summarize(context.Background(), apiKey, changelog)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if summary == "" {
		t.Error("expected non-empty summary")
	}
	t.Logf("Summary:\n%s", summary)
}

// TestSummarize_Within2000Chars は大量のCHANGELOGを渡しても
// 要約が2000文字以内に収まることを確認する統合テスト。
func TestSummarize_Within2000Chars(t *testing.T) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		t.Skip("GEMINI_API_KEY not set, skipping integration test")
	}

	// 大量の変更履歴を生成
	var changelog string
	for i := 1; i <= 50; i++ {
		changelog += fmt.Sprintf("## v%d.0.0\n", i)
		for j := 1; j <= 20; j++ {
			changelog += fmt.Sprintf("- Feature %d.%d: Added comprehensive support for advanced functionality including multiple subsystems and integration points\n", i, j)
		}
		changelog += "\n"
	}

	summary, err := Summarize(context.Background(), apiKey, changelog)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	charCount := len([]rune(summary))
	t.Logf("Summary length: %d chars", charCount)
	t.Logf("Summary:\n%s", summary)

	if charCount > 2000 {
		t.Errorf("summary exceeds 2000 characters: got %d", charCount)
	}
}
