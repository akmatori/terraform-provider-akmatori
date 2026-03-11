package client

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestGetAlertSourceTypes(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/alert-source-types" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode([]AlertSourceType{
			{ID: 1, Name: "alertmanager", DisplayName: "Prometheus Alertmanager"},
		})
	})

	types, err := c.GetAlertSourceTypes()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(types) != 1 {
		t.Errorf("expected 1 type, got %d", len(types))
	}
	if types[0].Name != "alertmanager" {
		t.Errorf("expected 'alertmanager', got '%s'", types[0].Name)
	}
}

func TestGetAlertSource(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/alert-sources/test-uuid" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(AlertSource{
			UUID: "test-uuid",
			Name: "prod-alertmanager",
			AlertSourceType: &AlertSourceType{Name: "alertmanager"},
		})
	})

	source, err := c.GetAlertSource("test-uuid")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if source.Name != "prod-alertmanager" {
		t.Errorf("expected 'prod-alertmanager', got '%s'", source.Name)
	}
	if source.SourceTypeName != "alertmanager" {
		t.Errorf("expected source type 'alertmanager', got '%s'", source.SourceTypeName)
	}
}

func TestCreateAlertSource(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		json.NewEncoder(w).Encode(AlertSource{UUID: "new-uuid", Name: "new-source"})
	})

	source, err := c.CreateAlertSource(CreateAlertSourceRequest{
		SourceTypeName: "alertmanager",
		Name:           "new-source",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if source.UUID != "new-uuid" {
		t.Errorf("expected UUID 'new-uuid', got '%s'", source.UUID)
	}
}

func TestDeleteAlertSource(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/alert-sources/del-uuid" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	})

	err := c.DeleteAlertSource("del-uuid")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
