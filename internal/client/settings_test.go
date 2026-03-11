package client

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestGetSlackSettings(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/settings/slack" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(SlackSettings{AlertsChannel: "#alerts", Enabled: true})
	})

	settings, err := c.GetSlackSettings()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if settings.AlertsChannel != "#alerts" {
		t.Errorf("expected '#alerts', got '%s'", settings.AlertsChannel)
	}
}

func TestGetLLMSettings(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/settings/llm" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(LLMSettings{Provider: "anthropic", Model: "claude-sonnet-4-20250514"})
	})

	settings, err := c.GetLLMSettings()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if settings.Provider != "anthropic" {
		t.Errorf("expected 'anthropic', got '%s'", settings.Provider)
	}
}

func TestGetProxySettings(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/settings/proxy" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(ProxySettings{ProxyURL: "http://proxy:8080"})
	})

	settings, err := c.GetProxySettings()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if settings.ProxyURL != "http://proxy:8080" {
		t.Errorf("expected 'http://proxy:8080', got '%s'", settings.ProxyURL)
	}
}

func TestGetAggregationSettings(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/settings/aggregation" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(AggregationSettings{
			Enabled:                        true,
			CorrelationConfidenceThreshold: 0.70,
		})
	})

	settings, err := c.GetAggregationSettings()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !settings.Enabled {
		t.Error("expected enabled to be true")
	}
}
