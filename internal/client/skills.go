package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Skill struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	IsSystem    bool      `json:"is_system"`
	Enabled     bool      `json:"enabled"`
	Prompt      string    `json:"prompt"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateSkillRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	Category    *string `json:"category,omitempty"`
	Enabled     *bool   `json:"enabled,omitempty"`
	Prompt      *string `json:"prompt,omitempty"`
}

type UpdateSkillRequest struct {
	Description *string `json:"description,omitempty"`
	Category    *string `json:"category,omitempty"`
	Enabled     *bool   `json:"enabled,omitempty"`
	Prompt      *string `json:"prompt,omitempty"`
}

func (c *Client) GetSkill(name string) (*Skill, error) {
	body, err := c.doRequest(http.MethodGet, fmt.Sprintf("/api/skills/%s", name), nil)
	if err != nil {
		return nil, err
	}
	var skill Skill
	if err := json.Unmarshal(body, &skill); err != nil {
		return nil, fmt.Errorf("failed to unmarshal skill: %w", err)
	}
	return &skill, nil
}

func (c *Client) CreateSkill(req CreateSkillRequest) (*Skill, error) {
	body, err := c.doRequest(http.MethodPost, "/api/skills", req)
	if err != nil {
		return nil, err
	}
	var skill Skill
	if err := json.Unmarshal(body, &skill); err != nil {
		return nil, fmt.Errorf("failed to unmarshal skill: %w", err)
	}
	return &skill, nil
}

func (c *Client) UpdateSkill(name string, req UpdateSkillRequest) (*Skill, error) {
	body, err := c.doRequest(http.MethodPut, fmt.Sprintf("/api/skills/%s", name), req)
	if err != nil {
		return nil, err
	}
	var skill Skill
	if err := json.Unmarshal(body, &skill); err != nil {
		return nil, fmt.Errorf("failed to unmarshal skill: %w", err)
	}
	return &skill, nil
}

func (c *Client) DeleteSkill(name string) error {
	_, err := c.doRequest(http.MethodDelete, fmt.Sprintf("/api/skills/%s", name), nil)
	return err
}

func (c *Client) UpdateSkillTools(name string, toolInstanceIDs []int) error {
	req := map[string][]int{"tool_instance_ids": toolInstanceIDs}
	_, err := c.doRequest(http.MethodPut, fmt.Sprintf("/api/skills/%s/tools", name), req)
	return err
}

func (c *Client) GetSkillTools(name string) ([]int, error) {
	skill, err := c.GetSkill(name)
	if err != nil {
		return nil, err
	}
	// The skill response includes tools; we need to get their IDs
	// We'll fetch the full skill to get tool IDs
	body, err := c.doRequest(http.MethodGet, fmt.Sprintf("/api/skills/%s", name), nil)
	if err != nil {
		return nil, err
	}

	var fullSkill struct {
		Tools []struct {
			ID int `json:"id"`
		} `json:"tools"`
	}
	if err := json.Unmarshal(body, &fullSkill); err != nil {
		return nil, fmt.Errorf("failed to unmarshal skill tools: %w", err)
	}
	_ = skill

	ids := make([]int, len(fullSkill.Tools))
	for i, t := range fullSkill.Tools {
		ids[i] = t.ID
	}
	return ids, nil
}

type Script struct {
	Filename string `json:"filename"`
	Content  string `json:"content"`
}

func (c *Client) GetSkillScript(skillName, filename string) (*Script, error) {
	body, err := c.doRequest(http.MethodGet, fmt.Sprintf("/api/skills/%s/scripts/%s", skillName, filename), nil)
	if err != nil {
		return nil, err
	}
	var script Script
	if err := json.Unmarshal(body, &script); err != nil {
		return nil, fmt.Errorf("failed to unmarshal script: %w", err)
	}
	return &script, nil
}

func (c *Client) UpdateSkillScript(skillName, filename, content string) (*Script, error) {
	req := map[string]string{"content": content}
	body, err := c.doRequest(http.MethodPut, fmt.Sprintf("/api/skills/%s/scripts/%s", skillName, filename), req)
	if err != nil {
		return nil, err
	}
	var script Script
	if err := json.Unmarshal(body, &script); err != nil {
		return nil, fmt.Errorf("failed to unmarshal script: %w", err)
	}
	return &script, nil
}

func (c *Client) DeleteSkillScript(skillName, filename string) error {
	_, err := c.doRequest(http.MethodDelete, fmt.Sprintf("/api/skills/%s/scripts/%s", skillName, filename), nil)
	return err
}
