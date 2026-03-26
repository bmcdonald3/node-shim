package client

import (
"context"
"fmt"
"net/http"
"time"
)

type MetadataClient struct {
BaseURL    string
HTTPClient *http.Client
}

func NewMetadataClient(baseURL string) *MetadataClient {
if baseURL == "" {
baseURL = "http://metadata-service:8080"
}
return &MetadataClient{
BaseURL: baseURL,
HTTPClient: &http.Client{
Timeout: 10 * time.Second,
},
}
}

func (c *MetadataClient) UpdateProfile(ctx context.Context, xname, profile string) error {
req, err := http.NewRequestWithContext(ctx, "PUT", fmt.Sprintf("%s/admin/profiles/%s/nodes/%s", c.BaseURL, profile, xname), nil)
if err != nil {
return fmt.Errorf("failed to create request: %w", err)
}

resp, err := c.HTTPClient.Do(req)
if err != nil {
return fmt.Errorf("failed to execute request: %w", err)
}
defer resp.Body.Close()

if resp.StatusCode != http.StatusOK {
return fmt.Errorf("metadata-service returned status code: %d", resp.StatusCode)
}

return nil
}