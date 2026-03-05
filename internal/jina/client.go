package jina

import (
	"fmt"
	"io"
	"net/http"
)

// FetchMarkdown fetches the given URL via Jina Reader API and returns its content as Markdown.
func FetchMarkdown(url string) (string, error) {
	resp, err := http.Get("https://r.jina.ai/" + url)
	if err != nil {
		return "", fmt.Errorf("jina reader request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("jina reader returned %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read jina reader response: %w", err)
	}

	return string(body), nil
}
