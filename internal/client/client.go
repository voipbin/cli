package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// API is the interface for all REST operations. Commands depend on this.
type API interface {
	List(ctx context.Context, path string, params url.Values) ([]map[string]interface{}, string, error)
	Get(ctx context.Context, path string) (map[string]interface{}, error)
	Post(ctx context.Context, path string, body interface{}) (map[string]interface{}, error)
	Put(ctx context.Context, path string, body interface{}) (map[string]interface{}, error)
	Delete(ctx context.Context, path string) (map[string]interface{}, error)
	RawGet(ctx context.Context, path string) (*http.Response, error)
}

// Client implements API using net/http.
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

type accessKeyTransport struct {
	accessKey string
}

func (t *accessKeyTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	newURL := *req.URL
	query := newURL.Query()
	query.Set("accesskey", t.accessKey)
	newURL.RawQuery = query.Encode()

	newReq := req.Clone(req.Context())
	newReq.URL = &newURL
	return http.DefaultTransport.RoundTrip(newReq)
}

// New creates a Client with access key authentication, 30s timeout,
// and cross-domain redirect protection.
func New(baseURL, accessKey string) *Client {
	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout:   30 * time.Second,
			Transport: &accessKeyTransport{accessKey: accessKey},
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				if len(via) > 0 && req.URL.Host != via[0].URL.Host {
					return http.ErrUseLastResponse
				}
				return nil
			},
		},
	}
}

// listEnvelope is the standard API list response shape.
type listEnvelope struct {
	Result        []map[string]interface{} `json:"result"`
	NextPageToken string                   `json:"next_page_token"`
}

// do executes an HTTP request and returns the response.
func (c *Client) do(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
	u := c.BaseURL + path

	var reqBody io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, u, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return c.HTTPClient.Do(req)
}

// checkStatus returns an error if the status code is not in the success range.
func checkStatus(resp *http.Response, body []byte) error {
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}
	return fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
}

// doJSON executes a request, checks status, and parses JSON.
// Returns (map, nil) on success; (nil, error) on failure. Never (nil, nil).
func (c *Client) doJSON(ctx context.Context, method, path string, body interface{}) (map[string]interface{}, error) {
	resp, err := c.do(ctx, method, path, body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if err := checkStatus(resp, data); err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return map[string]interface{}{}, nil
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return result, nil
}

// List sends GET and parses the standard list envelope {"result": [...], "next_page_token": "..."}.
func (c *Client) List(ctx context.Context, path string, params url.Values) ([]map[string]interface{}, string, error) {
	fullPath := path
	if len(params) > 0 {
		sep := "?"
		if strings.Contains(path, "?") {
			sep = "&"
		}
		fullPath = path + sep + params.Encode()
	}

	resp, err := c.do(ctx, http.MethodGet, fullPath, nil)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read response body: %w", err)
	}

	if err := checkStatus(resp, data); err != nil {
		return nil, "", err
	}

	var envelope listEnvelope
	if err := json.Unmarshal(data, &envelope); err != nil {
		return nil, "", fmt.Errorf("failed to parse list response: %w", err)
	}

	items := envelope.Result
	if items == nil {
		items = []map[string]interface{}{}
	}

	return items, envelope.NextPageToken, nil
}

// Get sends GET and returns a single JSON object.
func (c *Client) Get(ctx context.Context, path string) (map[string]interface{}, error) {
	return c.doJSON(ctx, http.MethodGet, path, nil)
}

// Post sends POST with optional JSON body. Treats 200 and 201 as success.
func (c *Client) Post(ctx context.Context, path string, body interface{}) (map[string]interface{}, error) {
	return c.doJSON(ctx, http.MethodPost, path, body)
}

// Put sends PUT with JSON body.
func (c *Client) Put(ctx context.Context, path string, body interface{}) (map[string]interface{}, error) {
	return c.doJSON(ctx, http.MethodPut, path, body)
}

// Delete sends DELETE and returns the parsed response body, if any.
func (c *Client) Delete(ctx context.Context, path string) (map[string]interface{}, error) {
	return c.doJSON(ctx, http.MethodDelete, path, nil)
}

// RawGet sends GET and returns the raw http.Response without JSON parsing.
// Caller is responsible for closing resp.Body.
func (c *Client) RawGet(ctx context.Context, path string) (*http.Response, error) {
	return c.do(ctx, http.MethodGet, path, nil)
}
