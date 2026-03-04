package state

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad_FileNotExist(t *testing.T) {
	s, err := Load(filepath.Join(t.TempDir(), "nonexistent.json"))
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(s.Repos) != 0 {
		t.Errorf("expected empty repos, got %v", s.Repos)
	}
}

func TestLoad_ValidFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.json")
	data := `{"repos":{"anthropics/claude-code":{"last_seen_sha":"abc123"}}}`
	if err := os.WriteFile(path, []byte(data), 0644); err != nil {
		t.Fatal(err)
	}

	s, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.SHA("anthropics/claude-code") != "abc123" {
		t.Errorf("expected %q, got %q", "abc123", s.SHA("anthropics/claude-code"))
	}
}

func TestLoad_LegacyFormat(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.json")
	if err := os.WriteFile(path, []byte(`{"last_seen_sha":"abc123"}`), 0644); err != nil {
		t.Fatal(err)
	}

	s, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.SHA("anthropics/claude-code") != "abc123" {
		t.Errorf("expected legacy SHA migrated, got %q", s.SHA("anthropics/claude-code"))
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.json")
	if err := os.WriteFile(path, []byte(`not json`), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestSaveLoad(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.json")
	var s State
	s.Repos = make(map[string]RepoState)
	s.SetSHA("openai/codex", "xyz789")

	if err := Save(path, s); err != nil {
		t.Fatalf("unexpected Save error: %v", err)
	}

	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected Load error: %v", err)
	}
	if loaded.SHA("openai/codex") != "xyz789" {
		t.Errorf("expected %q, got %q", "xyz789", loaded.SHA("openai/codex"))
	}
}

func TestSHA_Missing(t *testing.T) {
	s := State{Repos: make(map[string]RepoState)}
	if s.SHA("nonexistent/repo") != "" {
		t.Errorf("expected empty SHA for missing repo")
	}
}
