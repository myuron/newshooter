package gemini

import (
	"context"
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
