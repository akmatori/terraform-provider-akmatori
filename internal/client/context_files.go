package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"
)

type ContextFile struct {
	ID           int       `json:"id"`
	Filename     string    `json:"filename"`
	OriginalName string    `json:"original_name"`
	MimeType     string    `json:"mime_type"`
	Size         int64     `json:"size"`
	Description  string    `json:"description"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (c *Client) GetContextFile(id int) (*ContextFile, error) {
	body, err := c.doRequest(http.MethodGet, fmt.Sprintf("/api/context/%d", id), nil)
	if err != nil {
		return nil, err
	}
	var file ContextFile
	if err := json.Unmarshal(body, &file); err != nil {
		return nil, fmt.Errorf("failed to unmarshal context file: %w", err)
	}
	return &file, nil
}

func (c *Client) UploadContextFile(filename string, content []byte, description string) (*ContextFile, error) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}
	if _, err := part.Write(content); err != nil {
		return nil, fmt.Errorf("failed to write file content: %w", err)
	}

	if description != "" {
		if err := writer.WriteField("description", description); err != nil {
			return nil, fmt.Errorf("failed to write description field: %w", err)
		}
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, c.Host+"/api/context", &buf)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, &APIError{
			StatusCode: resp.StatusCode,
			Message:    string(respBody),
		}
	}

	var file ContextFile
	if err := json.Unmarshal(respBody, &file); err != nil {
		return nil, fmt.Errorf("failed to unmarshal context file: %w", err)
	}
	return &file, nil
}

func (c *Client) DeleteContextFile(id int) error {
	_, err := c.doRequest(http.MethodDelete, fmt.Sprintf("/api/context/%d", id), nil)
	return err
}
