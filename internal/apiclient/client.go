package apiclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"sync"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	DefaultBaseURL = "https://api.statuspage.io/v1"
)

type ClientOption func(*Client)

func WithRateLimitInterval(d time.Duration) ClientOption {
	return func(c *Client) {
		c.rateLimitInterval = d
	}
}

func WithDebug() ClientOption {
	return func(c *Client) {
		c.debug = true
	}
}

type APIError struct {
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("statuspage API error (HTTP %d): %s", e.StatusCode, e.Message)
}

type Client struct {
	baseURL           string
	apiKey            string
	httpClient        *http.Client
	mu                sync.Mutex
	lastReq           time.Time
	rateLimitInterval time.Duration
	debug             bool
}

func NewClient(apiKey string, opts ...ClientOption) *Client {
	c := &Client{
		baseURL:           DefaultBaseURL,
		apiKey:            apiKey,
		httpClient:        &http.Client{Timeout: 30 * time.Second},
		rateLimitInterval: time.Second,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *Client) rateLimit() {
	if c.rateLimitInterval <= 0 {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	elapsed := time.Since(c.lastReq)
	if elapsed < c.rateLimitInterval {
		time.Sleep(c.rateLimitInterval - elapsed)
	}
	c.lastReq = time.Now()
}

func (c *Client) doRequest(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
	c.rateLimit()

	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshaling request body: %w", err)
		}
		reqBody = bytes.NewReader(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "OAuth "+c.apiKey)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if c.debug {
		b, _ := httputil.DumpRequest(req, true)
		additionalField := map[string]any{
			"req": fmt.Sprintf("---request begin---\n%s\n---request end---\n", b),
		}
		tflog.Debug(ctx, "statuspage API request", additionalField)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if c.debug {
		b, _ := httputil.DumpResponse(resp, true)
		additionalField := map[string]any{
			"req": fmt.Sprintf("---response begin---\n%s\n---response end---\n", b),
		}
		tflog.Debug(ctx, "statuspage API response", additionalField)
	}

	if resp.StatusCode == 420 || resp.StatusCode == 429 {
		resp.Body.Close() //nolint:errcheck
		time.Sleep(2 * time.Second)
		return c.doRequest(ctx, method, path, body)
	}

	return resp, nil
}

func (c *Client) checkResponse(resp *http.Response) error {
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}
	defer resp.Body.Close() //nolint:errcheck
	body, _ := io.ReadAll(resp.Body)
	var errResp struct {
		Error   string `json:"error"`
		Message string `json:"message"`
	}
	json.Unmarshal(body, &errResp) //nolint:errcheck
	msg := errResp.Error
	if msg == "" {
		msg = errResp.Message
	}
	if msg == "" {
		msg = string(body)
	}
	return &APIError{StatusCode: resp.StatusCode, Message: msg}
}

func (c *Client) Get(ctx context.Context, path string, result interface{}) error {
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close() //nolint:errcheck

	if err := c.checkResponse(resp); err != nil {
		return err
	}

	return json.NewDecoder(resp.Body).Decode(result)
}

func (c *Client) Post(ctx context.Context, path string, body, result interface{}) error {
	resp, err := c.doRequest(ctx, http.MethodPost, path, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close() //nolint:errcheck

	if err := c.checkResponse(resp); err != nil {
		return err
	}

	if result != nil {
		return json.NewDecoder(resp.Body).Decode(result)
	}
	return nil
}

func (c *Client) Patch(ctx context.Context, path string, body, result interface{}) error {
	resp, err := c.doRequest(ctx, http.MethodPatch, path, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close() //nolint:errcheck

	if err := c.checkResponse(resp); err != nil {
		return err
	}

	if result != nil {
		return json.NewDecoder(resp.Body).Decode(result)
	}
	return nil
}

func (c *Client) Put(ctx context.Context, path string, body, result interface{}) error {
	resp, err := c.doRequest(ctx, http.MethodPut, path, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close() //nolint:errcheck

	if err := c.checkResponse(resp); err != nil {
		return err
	}

	if result != nil {
		return json.NewDecoder(resp.Body).Decode(result)
	}
	return nil
}

func (c *Client) Delete(ctx context.Context, path string) error {
	resp, err := c.doRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close() //nolint:errcheck

	return c.checkResponse(resp)
}
