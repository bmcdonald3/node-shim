package client

import (
"context"
"fmt"
"net/http"
"time"
)

type BootClient struct {
BaseURL    string
HTTPClient *http.Client
}

func NewBootClient(baseURL string) *BootClient {
if baseURL == "" {
baseURL = "http://boot-service:8080"
}
return &BootClient{
BaseURL: baseURL,
HTTPClient: &http.Client{
Timeout: 10 * time.Second,
},
}
}

func (c *BootClient) UpdateProfile(ctx context.Context, xname, profile string) error {
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
return fmt.Errorf("boot-service returned status code: %d", resp.StatusCode)
}

return nil
}