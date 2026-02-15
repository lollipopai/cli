package httpclient

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var userAgent = "chp-cli/dev"

// SetUserAgent sets the User-Agent header for all requests.
func SetUserAgent(ua string) {
	userAgent = ua
}

// APIError is returned when an HTTP request fails.
type APIError struct {
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	return e.Message
}

// Client wraps net/http with localhost TLS skip and error parsing.
type Client struct {
	standard *http.Client
	insecure *http.Client
}

// New creates a Client with 30s timeout and a separate insecure transport for localhost.
func New() *Client {
	return &Client{
		standard: &http.Client{Timeout: 30 * time.Second},
		insecure: &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
	}
}

func (c *Client) clientFor(rawURL string) *http.Client {
	u, err := url.Parse(rawURL)
	if err != nil {
		return c.standard
	}
	host := u.Hostname()
	if host == "localhost" || host == "127.0.0.1" {
		return c.insecure
	}
	return c.standard
}

func (c *Client) do(req *http.Request) ([]byte, *http.Response, error) {
	req.Header.Set("User-Agent", userAgent)

	client := c.clientFor(req.URL.String())
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, &APIError{
			Message: fmt.Sprintf("Connection failed: %v\nCheck your network and base URL: %s", err, req.URL.String()),
		}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp, &APIError{
			Message: fmt.Sprintf("Failed to read response: %v", err),
		}
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		msg := parseErrorBody(body)
		return body, resp, &APIError{
			StatusCode: resp.StatusCode,
			Message:    fmt.Sprintf("HTTP %d: %s", resp.StatusCode, msg),
		}
	}

	return body, resp, nil
}

func parseErrorBody(body []byte) string {
	var obj map[string]any
	if err := json.Unmarshal(body, &obj); err == nil {
		if e, ok := obj["error"].(string); ok {
			return e
		}
		if m, ok := obj["msg"].(string); ok {
			return m
		}
	}
	return string(body)
}

// GetJSON performs a GET request and returns the response body.
func (c *Client) GetJSON(rawURL string, headers map[string]string) ([]byte, error) {
	req, err := http.NewRequest("GET", rawURL, nil)
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	body, _, err := c.do(req)
	return body, err
}

// PostJSON sends a JSON-encoded payload and returns the response body.
func (c *Client) PostJSON(rawURL string, payload any, headers map[string]string) ([]byte, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", rawURL, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	body, _, err := c.do(req)
	return body, err
}

// PostJSONRaw sends JSON and returns both body and the raw http.Response.
func (c *Client) PostJSONRaw(rawURL string, payload any, headers map[string]string) ([]byte, *http.Response, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, nil, err
	}
	req, err := http.NewRequest("POST", rawURL, bytes.NewReader(data))
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	return c.do(req)
}

// PostForm sends a form-encoded POST and returns the response body.
func (c *Client) PostForm(rawURL string, params url.Values, headers map[string]string) ([]byte, error) {
	req, err := http.NewRequest("POST", rawURL, strings.NewReader(params.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	body, _, err := c.do(req)
	return body, err
}
