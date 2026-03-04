package github

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchFile_Success(t *testing.T) {
	content := "# CHANGELOG\n## v1.0.0\n- Initial release"
	encoded := base64.StdEncoding.EncodeToString([]byte(content))

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer testtoken" {
			t.Errorf("expected auth header 'Bearer testtoken', got %q", r.Header.Get("Authorization"))
		}
		if err := json.NewEncoder(w).Encode(map[string]string{
			"sha":     "deadbeef",
			"content": encoded,
		}); err != nil {
			t.Error(err)
		}
	}))
	defer srv.Close()

	old := apiBaseURL
	apiBaseURL = srv.URL
	defer func() { apiBaseURL = old }()

	fc, err := FetchFile("testtoken", "owner", "repo", "CHANGELOG.md")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fc.SHA != "deadbeef" {
		t.Errorf("expected SHA %q, got %q", "deadbeef", fc.SHA)
	}
	if fc.Content != content {
		t.Errorf("expected content %q, got %q", content, fc.Content)
	}
}

func TestFetchFile_NoToken(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "" {
			t.Errorf("expected no auth header, got %q", r.Header.Get("Authorization"))
		}
		encoded := base64.StdEncoding.EncodeToString([]byte("content"))
		if err := json.NewEncoder(w).Encode(map[string]string{"sha": "abc", "content": encoded}); err != nil {
			t.Error(err)
		}
	}))
	defer srv.Close()

	old := apiBaseURL
	apiBaseURL = srv.URL
	defer func() { apiBaseURL = old }()

	_, err := FetchFile("", "owner", "repo", "CHANGELOG.md")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestFetchFile_NonOKStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte("not found"))
	}))
	defer srv.Close()

	old := apiBaseURL
	apiBaseURL = srv.URL
	defer func() { apiBaseURL = old }()

	_, err := FetchFile("", "owner", "repo", "CHANGELOG.md")
	if err == nil {
		t.Fatal("expected error for non-OK status")
	}
}

func TestFetchFile_InvalidBase64(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewEncoder(w).Encode(map[string]string{
			"sha":     "abc",
			"content": "!!!invalid base64!!!",
		}); err != nil {
			t.Error(err)
		}
	}))
	defer srv.Close()

	old := apiBaseURL
	apiBaseURL = srv.URL
	defer func() { apiBaseURL = old }()

	_, err := FetchFile("", "owner", "repo", "CHANGELOG.md")
	if err == nil {
		t.Fatal("expected error for invalid base64")
	}
}

func TestFetchFile_Base64WithNewlines(t *testing.T) {
	content := "hello world"
	raw := base64.StdEncoding.EncodeToString([]byte(content))
	// GitHub APIが返すように改行を挿入
	withNewlines := raw[:4] + "\n" + raw[4:]

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewEncoder(w).Encode(map[string]string{
			"sha":     "abc",
			"content": withNewlines,
		}); err != nil {
			t.Error(err)
		}
	}))
	defer srv.Close()

	old := apiBaseURL
	apiBaseURL = srv.URL
	defer func() { apiBaseURL = old }()

	fc, err := FetchFile("", "owner", "repo", "CHANGELOG.md")
	if err != nil {
		t.Fatalf("unexpected error for base64 with newlines: %v", err)
	}
	if fc.Content != content {
		t.Errorf("expected %q, got %q", content, fc.Content)
	}
}
