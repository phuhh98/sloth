package host

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	baseURL    string
	token      string
	httpClient *http.Client
}

type PluginStatus struct {
	PluginName               string           `json:"pluginName"`
	PluginVersion            string           `json:"pluginVersion"`
	CompatibleSchemaVersions []string         `json:"compatibleSchemaVersions"`
	Components               []map[string]any `json:"components"`
	TotalComponents          int              `json:"totalComponents"`
}

type ContractSchemaResponse struct {
	SchemaVersion string         `json:"schemaVersion"`
	SchemaURL     string         `json:"schemaUrl,omitempty"`
	Schema        map[string]any `json:"schema,omitempty"`
}

type IngestResponse struct {
	Created []string `json:"created"`
	Updated []string `json:"updated"`
	Failed  []struct {
		Name   string `json:"name"`
		Reason string `json:"reason"`
	} `json:"failed"`
	TotalReceived int `json:"totalReceived"`
}

func NewClient(baseURL string, token string) *Client {
	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		token:   strings.TrimSpace(token),
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) PluginStatus() (*PluginStatus, error) {
	resp, err := c.doRequest(http.MethodGet, "/sloth/inspection/plugin-status", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return decode[PluginStatus](resp.Body)
}

func (c *Client) ContractSchema(schemaVersion string, inline bool) (*ContractSchemaResponse, error) {
	query := url.Values{}
	if strings.TrimSpace(schemaVersion) != "" {
		query.Set("schemaVersion", schemaVersion)
	}
	query.Set("inline", fmt.Sprintf("%t", inline))
	path := "/sloth/inspection/contract-schema"
	if encoded := query.Encode(); encoded != "" {
		path = path + "?" + encoded
	}

	resp, err := c.doRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return decode[ContractSchemaResponse](resp.Body)
}

func (c *Client) IngestContracts(contracts []map[string]any) (*IngestResponse, error) {
	payload := map[string]any{"contracts": contracts}
	resp, err := c.doRequest(http.MethodPost, "/sloth/contracts/ingest", payload)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return decode[IngestResponse](resp.Body)
}

func (c *Client) doRequest(method string, path string, payload any) (*http.Response, error) {
	var body io.Reader
	if payload != nil {
		raw, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("marshal request payload: %w", err)
		}
		body = bytes.NewBuffer(raw)
	}

	req, err := http.NewRequest(method, c.baseURL+path, body)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("perform request %s %s: %w", method, path, err)
	}

	if resp.StatusCode >= 400 {
		raw, _ := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		return nil, fmt.Errorf("request failed %s %s: status=%d body=%s", method, path, resp.StatusCode, strings.TrimSpace(string(raw)))
	}

	return resp, nil
}

func decode[T any](r io.Reader) (*T, error) {
	v := new(T)
	if err := json.NewDecoder(r).Decode(v); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return v, nil
}
