package client

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestGetContextFile(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/context/1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(ContextFile{ID: 1, Filename: "runbook.md", Description: "Runbook"})
	})

	file, err := c.GetContextFile(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if file.Filename != "runbook.md" {
		t.Errorf("expected 'runbook.md', got '%s'", file.Filename)
	}
}

func TestDeleteContextFile(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/context/5" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	})

	err := c.DeleteContextFile(5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
