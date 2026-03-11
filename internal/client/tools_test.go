package client

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestGetToolTypes(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/tool-types" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode([]ToolType{
			{ID: 1, Name: "zabbix", Description: "Zabbix monitoring"},
			{ID: 2, Name: "ssh", Description: "SSH access"},
		})
	})

	types, err := c.GetToolTypes()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(types) != 2 {
		t.Errorf("expected 2 types, got %d", len(types))
	}
}

func TestGetToolInstance(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/tools/42" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(ToolInstance{
			ID:   42,
			Name: "prod-zabbix",
			ToolType: &ToolType{Name: "zabbix"},
		})
	})

	instance, err := c.GetToolInstance(42)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if instance.Name != "prod-zabbix" {
		t.Errorf("expected name 'prod-zabbix', got '%s'", instance.Name)
	}
	if instance.ToolTypeName != "zabbix" {
		t.Errorf("expected tool type name 'zabbix', got '%s'", instance.ToolTypeName)
	}
}

func TestCreateToolInstance(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		json.NewEncoder(w).Encode(ToolInstance{ID: 10, Name: "new-tool"})
	})

	inst, err := c.CreateToolInstance(CreateToolInstanceRequest{
		ToolTypeID: 1,
		Name:       "new-tool",
		Enabled:    true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if inst.ID != 10 {
		t.Errorf("expected ID 10, got %d", inst.ID)
	}
}

func TestDeleteToolInstance(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/tools/5" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	})

	err := c.DeleteToolInstance(5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
