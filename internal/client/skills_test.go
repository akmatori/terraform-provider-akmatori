package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestClient(t *testing.T, handler http.HandlerFunc) *Client {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	return &Client{
		Host:       server.URL,
		Token:      "test-token",
		HTTPClient: server.Client(),
	}
}

func TestGetSkill(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/api/skills/test-skill" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Error("missing or incorrect authorization header")
		}
		json.NewEncoder(w).Encode(Skill{ID: 1, Name: "test-skill", Description: "A test skill"})
	})

	skill, err := c.GetSkill("test-skill")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if skill.Name != "test-skill" {
		t.Errorf("expected name 'test-skill', got '%s'", skill.Name)
	}
}

func TestCreateSkill(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/api/skills" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		var req CreateSkillRequest
		json.NewDecoder(r.Body).Decode(&req)
		if req.Name != "new-skill" {
			t.Errorf("expected name 'new-skill', got '%s'", req.Name)
		}
		json.NewEncoder(w).Encode(Skill{ID: 2, Name: "new-skill"})
	})

	skill, err := c.CreateSkill(CreateSkillRequest{Name: "new-skill"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if skill.ID != 2 {
		t.Errorf("expected ID 2, got %d", skill.ID)
	}
}

func TestDeleteSkill(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete || r.URL.Path != "/api/skills/old-skill" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	})

	err := c.DeleteSkill("old-skill")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetSkill_NotFound(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("not found"))
	})

	_, err := c.GetSkill("missing")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !IsNotFound(err) {
		t.Errorf("expected not found error, got: %v", err)
	}
}
