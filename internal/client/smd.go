package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// SMDNode represents a node as returned by the State Management Database (SMD)
type SMDNode struct {
	ID              string            `json:"ID"`
	State           string            `json:"State"`
	LastUpdate      time.Time         `json:"LastUpdate"`
	Labels          map[string]string `json:"Labels"`
}

type SMDClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewSMDClient(baseURL string) *SMDClient {
if baseURL == "" {
baseURL = "http://smd:27779"
}
return &SMDClient{
BaseURL: baseURL,
HTTPClient: &http.Client{
Timeout: 10 * time.Second,
},
}
}

func (c *SMDClient) ListNodes(ctx context.Context) ([]SMDNode, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/hms/v2/Inventory/Hardware/Nodes", c.BaseURL), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var nodes []SMDNode
	if err := json.NewDecoder(resp.Body).Decode(&nodes); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return nodes, nil
}