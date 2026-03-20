package slack

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	webhookURL string
	httpClient *http.Client
}

var instance *Client

func Set(c *Client) { instance = c }
func Get() *Client  { return instance }

func NewClient(webhookURL, _ string) *Client {
	return &Client{
		webhookURL: webhookURL,
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}
}

type slackPayload struct {
	Text string `json:"text"`
}

func (c *Client) SendAlert(ctx context.Context, title, message string) error {
	if c == nil || c.webhookURL == "" {
		return nil
	}
	payload := slackPayload{Text: fmt.Sprintf(":rotating_light: *%s*\n%s", title, message)}
	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.webhookURL, bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("slack: unexpected status %d", resp.StatusCode)
	}
	return nil
}
