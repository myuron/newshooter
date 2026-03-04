package discord

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSend_Success(t *testing.T) {
	var received webhookPayload

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected Content-Type application/json, got %q", r.Header.Get("Content-Type"))
		}
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Error(err)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	if err := Send(srv.URL, "Test Title", "Test Description"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(received.Embeds) != 1 {
		t.Fatalf("expected 1 embed, got %d", len(received.Embeds))
	}
	if received.Embeds[0].Title != "Test Title" {
		t.Errorf("expected title %q, got %q", "Test Title", received.Embeds[0].Title)
	}
	if received.Embeds[0].Description != "Test Description" {
		t.Errorf("expected description %q, got %q", "Test Description", received.Embeds[0].Description)
	}
}

func TestSend_TruncatesLongDescription(t *testing.T) {
	var received webhookPayload

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Error(err)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	longDesc := strings.Repeat("a", 5000)
	if err := Send(srv.URL, "Title", longDesc); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := received.Embeds[0].Description
	if len(got) != maxDescriptionLen {
		t.Errorf("expected len %d, got %d", maxDescriptionLen, len(got))
	}
	if !strings.HasSuffix(got, "...") {
		t.Errorf("expected description to end with '...', got %q", got[len(got)-10:])
	}
}

func TestSend_ShortDescriptionNotTruncated(t *testing.T) {
	var received webhookPayload

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Error(err)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	desc := "short description"
	if err := Send(srv.URL, "Title", desc); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if received.Embeds[0].Description != desc {
		t.Errorf("expected %q, got %q", desc, received.Embeds[0].Description)
	}
}

func TestSend_NonOKStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("bad request"))
	}))
	defer srv.Close()

	err := Send(srv.URL, "Title", "Description")
	if err == nil {
		t.Fatal("expected error for non-2xx status")
	}
}
