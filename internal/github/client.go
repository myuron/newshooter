package github

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

var apiBaseURL = "https://api.github.com"

type FileContent struct {
	SHA     string
	Content string
}

type contentsResponse struct {
	SHA     string `json:"sha"`
	Content string `json:"content"`
}

func FetchFile(token, owner, repo, path string) (*FileContent, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/contents/%s", apiBaseURL, owner, repo, path)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("GitHub API returned %d: %s", resp.StatusCode, string(body))
	}

	var cr contentsResponse
	if err := json.NewDecoder(resp.Body).Decode(&cr); err != nil {
		return nil, err
	}

	decoded, err := base64.StdEncoding.DecodeString(strings.ReplaceAll(cr.Content, "\n", ""))
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 content: %w", err)
	}

	return &FileContent{
		SHA:     cr.SHA,
		Content: string(decoded),
	}, nil
}

type Release struct {
	TagName     string `json:"tag_name"`
	Name        string `json:"name"`
	Body        string `json:"body"`
	PublishedAt string `json:"published_at"`
	Prerelease  bool   `json:"prerelease"`
	Draft       bool   `json:"draft"`
}

// FetchLatestRelease returns the most recent non-prerelease, non-draft release.
func FetchLatestRelease(token, owner, repo string) (*Release, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/releases?per_page=20", apiBaseURL, owner, repo)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("GitHub API returned %d: %s", resp.StatusCode, string(body))
	}

	var releases []Release
	if err := json.NewDecoder(resp.Body).Decode(&releases); err != nil {
		return nil, err
	}

	for i := range releases {
		r := &releases[i]
		if !r.Draft && !r.Prerelease && len(r.Body) > 30 {
			return r, nil
		}
	}

	return nil, fmt.Errorf("no stable release with content found")
}
