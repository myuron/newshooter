package gemini

import (
	"context"
	"fmt"

	"google.golang.org/genai"
)

const prompt = `以下はソフトウェアのCHANGELOGです。最新バージョンの変更内容を日本語で箇条書きに要約してください。
技術的な内容はそのまま保持し、簡潔にまとめてください。

CHANGELOG:
%s`

func Summarize(ctx context.Context, apiKey, changelog string) (string, error) {
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return "", fmt.Errorf("failed to create Gemini client: %w", err)
	}

	result, err := client.Models.GenerateContent(ctx, "gemini-2.5-flash", genai.Text(fmt.Sprintf(prompt, changelog)), nil)
	if err != nil {
		return "", fmt.Errorf("gemini API error: %w", err)
	}

	if len(result.Candidates) == 0 || len(result.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("gemini returned no content")
	}

	return result.Candidates[0].Content.Parts[0].Text, nil
}

const urlPrompt = `以下のURLのページ内容を日本語で箇条書きに要約してください。
技術的な内容はそのまま保持し、簡潔にまとめてください。

URL: %s`

func SummarizeURL(ctx context.Context, apiKey, url string) (string, error) {
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return "", fmt.Errorf("failed to create Gemini client: %w", err)
	}

	config := &genai.GenerateContentConfig{
		Tools: []*genai.Tool{
			{URLContext: &genai.URLContext{}},
		},
	}

	result, err := client.Models.GenerateContent(ctx, "gemini-2.5-flash", genai.Text(fmt.Sprintf(urlPrompt, url)), config)
	if err != nil {
		return "", fmt.Errorf("gemini API error: %w", err)
	}

	text := result.Text()
	if text == "" {
		reason := ""
		if len(result.Candidates) > 0 {
			reason = fmt.Sprintf(" (finish_reason=%s)", result.Candidates[0].FinishReason)
		}
		if result.PromptFeedback != nil {
			reason += fmt.Sprintf(" (block_reason=%s)", result.PromptFeedback.BlockReason)
		}
		return "", fmt.Errorf("gemini returned no content%s", reason)
	}

	return text, nil
}
