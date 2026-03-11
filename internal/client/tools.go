package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type ToolType struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Schema      any       `json:"schema"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ToolInstance struct {
	ID           int       `json:"id"`
	ToolTypeID   int       `json:"tool_type_id"`
	Name         string    `json:"name"`
	Settings     any       `json:"settings"`
	Enabled      bool      `json:"enabled"`
	ToolTypeName string    `json:"tool_type_name,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	ToolType     *ToolType `json:"tool_type,omitempty"`
}

type CreateToolInstanceRequest struct {
	ToolTypeID int    `json:"tool_type_id"`
	Name       string `json:"name"`
	Settings   any    `json:"settings,omitempty"`
	Enabled    bool   `json:"enabled"`
}

type UpdateToolInstanceRequest struct {
	Name     string `json:"name"`
	Settings any    `json:"settings,omitempty"`
	Enabled  bool   `json:"enabled"`
}

func (c *Client) GetToolTypes() ([]ToolType, error) {
	body, err := c.doRequest(http.MethodGet, "/api/tool-types", nil)
	if err != nil {
		return nil, err
	}
	var types []ToolType
	if err := json.Unmarshal(body, &types); err != nil {
		return nil, fmt.Errorf("failed to unmarshal tool types: %w", err)
	}
	return types, nil
}

func (c *Client) GetToolInstance(id int) (*ToolInstance, error) {
	body, err := c.doRequest(http.MethodGet, fmt.Sprintf("/api/tools/%d", id), nil)
	if err != nil {
		return nil, err
	}
	var instance ToolInstance
	if err := json.Unmarshal(body, &instance); err != nil {
		return nil, fmt.Errorf("failed to unmarshal tool instance: %w", err)
	}
	if instance.ToolType != nil {
		instance.ToolTypeName = instance.ToolType.Name
	}
	return &instance, nil
}

func (c *Client) CreateToolInstance(req CreateToolInstanceRequest) (*ToolInstance, error) {
	body, err := c.doRequest(http.MethodPost, "/api/tools", req)
	if err != nil {
		return nil, err
	}
	var instance ToolInstance
	if err := json.Unmarshal(body, &instance); err != nil {
		return nil, fmt.Errorf("failed to unmarshal tool instance: %w", err)
	}
	if instance.ToolType != nil {
		instance.ToolTypeName = instance.ToolType.Name
	}
	return &instance, nil
}

func (c *Client) UpdateToolInstance(id int, req UpdateToolInstanceRequest) (*ToolInstance, error) {
	body, err := c.doRequest(http.MethodPut, fmt.Sprintf("/api/tools/%d", id), req)
	if err != nil {
		return nil, err
	}
	var instance ToolInstance
	if err := json.Unmarshal(body, &instance); err != nil {
		return nil, fmt.Errorf("failed to unmarshal tool instance: %w", err)
	}
	if instance.ToolType != nil {
		instance.ToolTypeName = instance.ToolType.Name
	}
	return &instance, nil
}

func (c *Client) DeleteToolInstance(id int) error {
	_, err := c.doRequest(http.MethodDelete, fmt.Sprintf("/api/tools/%d", id), nil)
	return err
}

// SSH Key types
type SSHKey struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	IsDefault bool   `json:"is_default"`
}

func (c *Client) GetSSHKeys(toolID int) ([]SSHKey, error) {
	body, err := c.doRequest(http.MethodGet, fmt.Sprintf("/api/tools/%d/ssh-keys", toolID), nil)
	if err != nil {
		return nil, err
	}
	var keys []SSHKey
	if err := json.Unmarshal(body, &keys); err != nil {
		return nil, fmt.Errorf("failed to unmarshal SSH keys: %w", err)
	}
	return keys, nil
}

func (c *Client) CreateSSHKey(toolID int, name, privateKey string, isDefault bool) (*SSHKey, error) {
	req := map[string]interface{}{
		"name":        name,
		"private_key": privateKey,
		"is_default":  isDefault,
	}
	body, err := c.doRequest(http.MethodPost, fmt.Sprintf("/api/tools/%d/ssh-keys", toolID), req)
	if err != nil {
		return nil, err
	}
	var key SSHKey
	if err := json.Unmarshal(body, &key); err != nil {
		return nil, fmt.Errorf("failed to unmarshal SSH key: %w", err)
	}
	return &key, nil
}

func (c *Client) UpdateSSHKey(toolID int, keyID string, name *string, isDefault *bool) (*SSHKey, error) {
	req := map[string]interface{}{}
	if name != nil {
		req["name"] = *name
	}
	if isDefault != nil {
		req["is_default"] = *isDefault
	}
	body, err := c.doRequest(http.MethodPut, fmt.Sprintf("/api/tools/%d/ssh-keys/%s", toolID, keyID), req)
	if err != nil {
		return nil, err
	}
	var key SSHKey
	if err := json.Unmarshal(body, &key); err != nil {
		return nil, fmt.Errorf("failed to unmarshal SSH key: %w", err)
	}
	return &key, nil
}

func (c *Client) DeleteSSHKey(toolID int, keyID string) error {
	_, err := c.doRequest(http.MethodDelete, fmt.Sprintf("/api/tools/%d/ssh-keys/%s", toolID, keyID), nil)
	return err
}
