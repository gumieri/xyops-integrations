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

const (
	defaultTimeout = 30 * time.Second
	basePath       = "/api/app/"
)

// Client is the xyops API client.
type Client struct {
	baseURL    *url.URL
	apiKey     string
	httpClient *http.Client

	// Service fields — set by New after construction
	Tags        *TagsService
	Groups      *GroupsService
	Categories  *CategoriesService
	Plugins     *PluginsService
	Events      *EventsService
	Monitors    *MonitorsService
	Alerts      *AlertsService
	Roles       *RolesService
	APIKeys     *APIKeysService
	WebHooks    *WebHooksService
	Config      *ConfigService
	Servers     *ServersService
	Satellite   *SatelliteService
	Jobs        *JobsService
	GetMultiple *GetMultipleService
}

// New creates a new xyops API client.
func New(baseURL, apiKey string) (*Client, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL %q: %w", baseURL, err)
	}

	c := &Client{
		baseURL:    u,
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: defaultTimeout},
	}

	c.initServices()
	return c, nil
}

// NewWithHTTPClient creates a client with a custom HTTP client (for testing).
func NewWithHTTPClient(baseURL, apiKey string, httpClient *http.Client) (*Client, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL %q: %w", baseURL, err)
	}

	c := &Client{
		baseURL:    u,
		apiKey:     apiKey,
		httpClient: httpClient,
	}

	c.initServices()
	return c, nil
}

func (c *Client) initServices() {
	c.Tags = &TagsService{client: c}
	c.Groups = &GroupsService{client: c}
	c.Categories = &CategoriesService{client: c}
	c.Plugins = &PluginsService{client: c}
	c.Events = &EventsService{client: c}
	c.Monitors = &MonitorsService{client: c}
	c.Alerts = &AlertsService{client: c}
	c.Roles = &RolesService{client: c}
	c.APIKeys = &APIKeysService{client: c}
	c.WebHooks = &WebHooksService{client: c}
	c.Config = &ConfigService{client: c}
	c.Servers = &ServersService{client: c}
	c.Satellite = &SatelliteService{client: c}
	c.Jobs = &JobsService{client: c}
	c.GetMultiple = &GetMultipleService{client: c}
}

// DoRequest performs an HTTP request with auth headers and returns raw bytes.
func (c *Client) DoRequest(ctx context.Context, method, path string, body any) ([]byte, error) {
	reqURL := c.resolveURL(path)

	var reqBody io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, reqURL, reqBody)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("X-API-Key", c.apiKey)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	// Check HTTP status
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var apiErr struct {
			Code    int    `json:"code"`
			Message string `json:"message,omitempty"`
		}
		if err := json.Unmarshal(respBody, &apiErr); err == nil && (apiErr.Code != 0 || apiErr.Message != "") {
			return nil, NewAPIError(apiErr.Code, apiErr.Message, resp.StatusCode)
		}
		return nil, NewAPIError(0, string(respBody), resp.StatusCode)
	}

	// Check code field even on successful HTTP responses
	var bodyCode struct {
		Code int `json:"code"`
	}
	if err := json.Unmarshal(respBody, &bodyCode); err == nil && bodyCode.Code != 0 {
		var msg struct {
			Message string `json:"message"`
		}
		json.Unmarshal(respBody, &msg)
		return nil, NewAPIError(bodyCode.Code, msg.Message, resp.StatusCode)
	}

	return respBody, nil
}

// DoGET performs a GET request.
func (c *Client) DoGET(ctx context.Context, path string) ([]byte, error) {
	return c.DoRequest(ctx, http.MethodGet, path, nil)
}

// DoPOST performs a POST request with an optional JSON body.
func (c *Client) DoPOST(ctx context.Context, path string, body any) ([]byte, error) {
	return c.DoRequest(ctx, http.MethodPost, path, body)
}

// DoGETWithQuery performs a GET request with query parameters.
func (c *Client) DoGETWithQuery(ctx context.Context, path string, params map[string]string) ([]byte, error) {
	if len(params) == 0 {
		return c.DoGET(ctx, path)
	}

	var parts []string
	for k, v := range params {
		parts = append(parts, url.QueryEscape(k)+"="+url.QueryEscape(v))
	}
	return c.DoGET(ctx, path+"?"+strings.Join(parts, "&"))
}

// Health checks if the xyops instance is reachable.
func (c *Client) Health(ctx context.Context) (bool, error) {
	body, err := c.DoGET(ctx, "/health")
	if err != nil {
		// Instance not reachable or not primary
		return false, nil
	}

	resp := struct {
		Status  string `json:"status,omitempty"`
		Primary bool   `json:"primary,omitempty"`
	}{}
	if err := json.Unmarshal(body, &resp); err == nil {
		return resp.Status == "ok" || resp.Primary, nil
	}

	return true, nil
}

func (c *Client) resolveURL(path string) string {
	// If path is absolute (/...), append to base URL
	if strings.HasPrefix(path, "/") {
		return c.baseURL.String() + path
	}
	// Otherwise, prepend /api/app/
	return c.baseURL.String() + basePath + path
}
