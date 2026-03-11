package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type AlertSourceType struct {
	ID                   int    `json:"id"`
	Name                 string `json:"name"`
	DisplayName          string `json:"display_name"`
	Description          string `json:"description"`
	DefaultFieldMappings any    `json:"default_field_mappings"`
	WebhookSecretHeader  string `json:"webhook_secret_header"`
}

type AlertSource struct {
	ID              int       `json:"id"`
	UUID            string    `json:"uuid"`
	AlertSourceTypeID int    `json:"alert_source_type_id"`
	SourceTypeName  string    `json:"source_type_name,omitempty"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	WebhookSecret   string    `json:"webhook_secret"`
	WebhookURL      string    `json:"webhook_url,omitempty"`
	FieldMappings   any       `json:"field_mappings"`
	Settings        any       `json:"settings"`
	Enabled         bool      `json:"enabled"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	AlertSourceType *AlertSourceType `json:"alert_source_type,omitempty"`
}

type CreateAlertSourceRequest struct {
	SourceTypeName string `json:"source_type_name"`
	Name           string `json:"name"`
	Description    string `json:"description,omitempty"`
	WebhookSecret  string `json:"webhook_secret,omitempty"`
	FieldMappings  any    `json:"field_mappings,omitempty"`
	Settings       any    `json:"settings,omitempty"`
}

type UpdateAlertSourceRequest struct {
	Name          *string `json:"name,omitempty"`
	Description   *string `json:"description,omitempty"`
	WebhookSecret *string `json:"webhook_secret,omitempty"`
	FieldMappings any     `json:"field_mappings,omitempty"`
	Settings      any     `json:"settings,omitempty"`
	Enabled       *bool   `json:"enabled,omitempty"`
}

func (c *Client) GetAlertSourceTypes() ([]AlertSourceType, error) {
	body, err := c.doRequest(http.MethodGet, "/api/alert-source-types", nil)
	if err != nil {
		return nil, err
	}
	var types []AlertSourceType
	if err := json.Unmarshal(body, &types); err != nil {
		return nil, fmt.Errorf("failed to unmarshal alert source types: %w", err)
	}
	return types, nil
}

func (c *Client) GetAlertSource(uuid string) (*AlertSource, error) {
	body, err := c.doRequest(http.MethodGet, fmt.Sprintf("/api/alert-sources/%s", uuid), nil)
	if err != nil {
		return nil, err
	}
	var source AlertSource
	if err := json.Unmarshal(body, &source); err != nil {
		return nil, fmt.Errorf("failed to unmarshal alert source: %w", err)
	}
	if source.AlertSourceType != nil {
		source.SourceTypeName = source.AlertSourceType.Name
	}
	return &source, nil
}

func (c *Client) CreateAlertSource(req CreateAlertSourceRequest) (*AlertSource, error) {
	body, err := c.doRequest(http.MethodPost, "/api/alert-sources", req)
	if err != nil {
		return nil, err
	}
	var source AlertSource
	if err := json.Unmarshal(body, &source); err != nil {
		return nil, fmt.Errorf("failed to unmarshal alert source: %w", err)
	}
	if source.AlertSourceType != nil {
		source.SourceTypeName = source.AlertSourceType.Name
	}
	return &source, nil
}

func (c *Client) UpdateAlertSource(uuid string, req UpdateAlertSourceRequest) (*AlertSource, error) {
	body, err := c.doRequest(http.MethodPut, fmt.Sprintf("/api/alert-sources/%s", uuid), req)
	if err != nil {
		return nil, err
	}
	var source AlertSource
	if err := json.Unmarshal(body, &source); err != nil {
		return nil, fmt.Errorf("failed to unmarshal alert source: %w", err)
	}
	if source.AlertSourceType != nil {
		source.SourceTypeName = source.AlertSourceType.Name
	}
	return &source, nil
}

func (c *Client) DeleteAlertSource(uuid string) error {
	_, err := c.doRequest(http.MethodDelete, fmt.Sprintf("/api/alert-sources/%s", uuid), nil)
	return err
}
