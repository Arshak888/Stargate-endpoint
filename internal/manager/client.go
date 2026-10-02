package manager

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Arshak888/Stargate-endpoint/internal/protocol"
)

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
	Credential string
}

func NewClient(baseURL, credential string) *Client {
	return &Client{
		BaseURL: strings.TrimRight(baseURL, "/"),
		HTTPClient: &http.Client{Timeout: 15 * time.Second},
		Credential: credential,
	}
}

func (c *Client) Bootstrap(req protocol.BootstrapRequest) (protocol.BootstrapResponse, error) {
	var out protocol.BootstrapResponse
	if err := c.post("/api/v1/enrollment/bootstrap", req, "", &out); err != nil {
		return out, err
	}
	return out, nil
}

func (c *Client) Heartbeat(req protocol.HeartbeatRequest) (protocol.HeartbeatResponse, error) {
	var out protocol.HeartbeatResponse
	if err := c.post("/api/v1/endpoint/heartbeat", req, c.Credential, &out); err != nil {
		return out, err
	}
	return out, nil
}

func (c *Client) post(path string, body interface{}, credential string, out interface{}) error {
	data, err := json.Marshal(body)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, c.BaseURL+path, bytes.NewReader(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if credential != "" {
		req.Header.Set("Authorization", "Bearer "+credential)
	}
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("manager returned HTTP %d", resp.StatusCode)
	}
	return json.NewDecoder(resp.Body).Decode(out)
}
