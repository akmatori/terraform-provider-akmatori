package client

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Slack settings
type SlackSettings struct {
	ID            int    `json:"id"`
	BotToken      string `json:"bot_token"`
	SigningSecret string `json:"signing_secret"`
	AppToken      string `json:"app_token"`
	AlertsChannel string `json:"alerts_channel"`
	Enabled       bool   `json:"enabled"`
}

type UpdateSlackSettingsRequest struct {
	BotToken      *string `json:"bot_token,omitempty"`
	SigningSecret *string `json:"signing_secret,omitempty"`
	AppToken      *string `json:"app_token,omitempty"`
	AlertsChannel *string `json:"alerts_channel,omitempty"`
	Enabled       *bool   `json:"enabled,omitempty"`
}

func (c *Client) GetSlackSettings() (*SlackSettings, error) {
	body, err := c.doRequest(http.MethodGet, "/api/settings/slack", nil)
	if err != nil {
		return nil, err
	}
	var settings SlackSettings
	if err := json.Unmarshal(body, &settings); err != nil {
		return nil, fmt.Errorf("failed to unmarshal slack settings: %w", err)
	}
	return &settings, nil
}

func (c *Client) UpdateSlackSettings(req UpdateSlackSettingsRequest) (*SlackSettings, error) {
	body, err := c.doRequest(http.MethodPut, "/api/settings/slack", req)
	if err != nil {
		return nil, err
	}
	var settings SlackSettings
	if err := json.Unmarshal(body, &settings); err != nil {
		return nil, fmt.Errorf("failed to unmarshal slack settings: %w", err)
	}
	return &settings, nil
}

// LLM settings
type LLMSettings struct {
	ID            int    `json:"id"`
	Provider      string `json:"provider"`
	APIKey        string `json:"api_key"`
	Model         string `json:"model"`
	ThinkingLevel string `json:"thinking_level"`
	BaseURL       string `json:"base_url"`
}

type UpdateLLMSettingsRequest struct {
	Provider      *string `json:"provider,omitempty"`
	APIKey        *string `json:"api_key,omitempty"`
	Model         *string `json:"model,omitempty"`
	ThinkingLevel *string `json:"thinking_level,omitempty"`
	BaseURL       *string `json:"base_url,omitempty"`
}

func (c *Client) GetLLMSettings() (*LLMSettings, error) {
	body, err := c.doRequest(http.MethodGet, "/api/settings/llm", nil)
	if err != nil {
		return nil, err
	}
	var settings LLMSettings
	if err := json.Unmarshal(body, &settings); err != nil {
		return nil, fmt.Errorf("failed to unmarshal LLM settings: %w", err)
	}
	return &settings, nil
}

func (c *Client) UpdateLLMSettings(req UpdateLLMSettingsRequest) (*LLMSettings, error) {
	body, err := c.doRequest(http.MethodPut, "/api/settings/llm", req)
	if err != nil {
		return nil, err
	}
	var settings LLMSettings
	if err := json.Unmarshal(body, &settings); err != nil {
		return nil, fmt.Errorf("failed to unmarshal LLM settings: %w", err)
	}
	return &settings, nil
}

// Proxy settings
type ProxySettings struct {
	ProxyURL string `json:"proxy_url"`
	NoProxy  string `json:"no_proxy"`
	Services struct {
		OpenAI struct {
			Enabled bool `json:"enabled"`
		} `json:"openai"`
		Slack struct {
			Enabled bool `json:"enabled"`
		} `json:"slack"`
		Zabbix struct {
			Enabled bool `json:"enabled"`
		} `json:"zabbix"`
	} `json:"services"`
}

type UpdateProxySettingsRequest struct {
	ProxyURL string `json:"proxy_url"`
	NoProxy  string `json:"no_proxy"`
	Services struct {
		OpenAI struct {
			Enabled bool `json:"enabled"`
		} `json:"openai"`
		Slack struct {
			Enabled bool `json:"enabled"`
		} `json:"slack"`
		Zabbix struct {
			Enabled bool `json:"enabled"`
		} `json:"zabbix"`
	} `json:"services"`
}

func (c *Client) GetProxySettings() (*ProxySettings, error) {
	body, err := c.doRequest(http.MethodGet, "/api/settings/proxy", nil)
	if err != nil {
		return nil, err
	}
	var settings ProxySettings
	if err := json.Unmarshal(body, &settings); err != nil {
		return nil, fmt.Errorf("failed to unmarshal proxy settings: %w", err)
	}
	return &settings, nil
}

func (c *Client) UpdateProxySettings(req UpdateProxySettingsRequest) (*ProxySettings, error) {
	body, err := c.doRequest(http.MethodPut, "/api/settings/proxy", req)
	if err != nil {
		return nil, err
	}
	var settings ProxySettings
	if err := json.Unmarshal(body, &settings); err != nil {
		return nil, fmt.Errorf("failed to unmarshal proxy settings: %w", err)
	}
	return &settings, nil
}

// Aggregation settings
type AggregationSettings struct {
	ID                             int     `json:"id"`
	Enabled                        bool    `json:"enabled"`
	CorrelationConfidenceThreshold float64 `json:"correlation_confidence_threshold"`
	MergeConfidenceThreshold       float64 `json:"merge_confidence_threshold"`
	RecorrelationEnabled           bool    `json:"recorrelation_enabled"`
	RecorrelationIntervalMinutes   int     `json:"recorrelation_interval_minutes"`
	MaxIncidentsToAnalyze          int     `json:"max_incidents_to_analyze"`
	ObservingDurationMinutes       int     `json:"observing_duration_minutes"`
	CorrelatorTimeoutSeconds       int     `json:"correlator_timeout_seconds"`
	MergeAnalyzerTimeoutSeconds    int     `json:"merge_analyzer_timeout_seconds"`
}

func (c *Client) GetAggregationSettings() (*AggregationSettings, error) {
	body, err := c.doRequest(http.MethodGet, "/api/settings/aggregation", nil)
	if err != nil {
		return nil, err
	}
	var settings AggregationSettings
	if err := json.Unmarshal(body, &settings); err != nil {
		return nil, fmt.Errorf("failed to unmarshal aggregation settings: %w", err)
	}
	return &settings, nil
}

func (c *Client) UpdateAggregationSettings(req AggregationSettings) (*AggregationSettings, error) {
	body, err := c.doRequest(http.MethodPut, "/api/settings/aggregation", req)
	if err != nil {
		return nil, err
	}
	var settings AggregationSettings
	if err := json.Unmarshal(body, &settings); err != nil {
		return nil, fmt.Errorf("failed to unmarshal aggregation settings: %w", err)
	}
	return &settings, nil
}
