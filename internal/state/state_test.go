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
	if s.LastSeenSHA != "" {
		t.Errorf("expected empty SHA, got %q", s.LastSeenSHA)
	}
}

func TestLoad_ValidFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.json")
	if err := os.WriteFile(path, []byte(`{"last_seen_sha":"abc123"}`), 0644); err != nil {
		t.Fatal(err)
	}

	s, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.LastSeenSHA != "abc123" {
		t.Errorf("expected %q, got %q", "abc123", s.LastSeenSHA)
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

func TestSave(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.json")
	s := State{LastSeenSHA: "xyz789"}

	if err := Save(path, s); err != nil {
		t.Fatalf("unexpected Save error: %v", err)
	}

	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected Load error: %v", err)
	}
	if loaded.LastSeenSHA != "xyz789" {
		t.Errorf("expected %q, got %q", "xyz789", loaded.LastSeenSHA)
	}
}

func TestSaveLoad_EmptySHA(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.json")
	s := State{LastSeenSHA: ""}

	if err := Save(path, s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if loaded.LastSeenSHA != "" {
		t.Errorf("expected empty SHA, got %q", loaded.LastSeenSHA)
	}
}
