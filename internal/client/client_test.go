package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewClient_WithToken(t *testing.T) {
	c, err := NewClient("http://localhost", "my-token", "", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Token != "my-token" {
		t.Errorf("expected token 'my-token', got '%s'", c.Token)
	}
}

func TestNewClient_WithUsernamePassword(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/auth/login" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(loginResponse{Token: "test-jwt-token"})
	}))
	defer server.Close()

	c, err := NewClient(server.URL, "", "admin", "password")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Token != "test-jwt-token" {
		t.Errorf("expected token 'test-jwt-token', got '%s'", c.Token)
	}
}

func TestNewClient_NoCredentials(t *testing.T) {
	_, err := NewClient("http://localhost", "", "", "")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestNewClient_LoginFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("invalid credentials"))
	}))
	defer server.Close()

	_, err := NewClient(server.URL, "", "admin", "wrong")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestIsNotFound(t *testing.T) {
	notFoundErr := &APIError{StatusCode: http.StatusNotFound, Message: "not found"}
	if !IsNotFound(notFoundErr) {
		t.Error("expected IsNotFound to return true for 404")
	}

	otherErr := &APIError{StatusCode: http.StatusInternalServerError, Message: "error"}
	if IsNotFound(otherErr) {
		t.Error("expected IsNotFound to return false for 500")
	}

	if IsNotFound(nil) {
		t.Error("expected IsNotFound to return false for nil")
	}
}

func TestAPIError_Error(t *testing.T) {
	err := &APIError{StatusCode: 404, Message: "not found"}
	expected := "API error (status 404): not found"
	if err.Error() != expected {
		t.Errorf("expected '%s', got '%s'", expected, err.Error())
	}
}
