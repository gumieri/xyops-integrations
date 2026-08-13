package sdk

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
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
	c, err := newClient(baseURL, apiKey, &http.Client{Timeout: defaultTimeout})
	if err != nil {
		return nil, err
	}
	return c, nil
}

// NewWithHTTPClient creates a client with a custom HTTP client (for testing).
func NewWithHTTPClient(baseURL, apiKey string, httpClient *http.Client) (*Client, error) {
	if httpClient == nil {
		return nil, fmt.Errorf("invalid base URL %q: nil http client", baseURL)
	}
	return newClient(baseURL, apiKey, httpClient)
}

func newClient(baseURL, apiKey string, httpClient *http.Client) (*Client, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL %q: %w", baseURL, err)
	}
	// url.Parse lowercases the scheme per RFC 3986.
	if u.Scheme != "http" && u.Scheme != "https" {
		return nil, fmt.Errorf("invalid base URL %q: scheme must be http or https", baseURL)
	}
	if u.Hostname() == "" {
		return nil, fmt.Errorf("invalid base URL %q: missing host", baseURL)
	}
	if u.RawQuery != "" || u.Fragment != "" || u.User != nil {
		return nil, fmt.Errorf("invalid base URL %q: must not contain query, fragment, or userinfo", baseURL)
	}
	// Trim trailing slashes so resolveURL never produces a double slash
	// before /api/app/. A path after the host is preserved and used as-is.
	u.Path = strings.TrimRight(u.Path, "/")

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
// Follows redirects automatically (http.Client default). On HTTP 3xx, 4xx, 5xx,
// the response body is parsed for code/message and an APIError is returned.
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
	defer func() {
		_ = resp.Body.Close()
	}()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	if err := checkErrorResponse(resp.StatusCode, respBody); err != nil {
		return nil, err
	}
	return respBody, nil
}

// checkErrorResponse inspects a response for API-reported errors.
// Per the API spec, errors use:
//   - code: numeric (0 for success, non-zero for error) OR string
//     (e.g., "session" for auth errors)
//   - description: human-readable error message (not "message")
//
// A "message" field is also accepted as a fallback for the description.
// Errors are detected on non-2xx HTTP status, on non-zero numeric codes,
// on the "session" string code, and on any string code carrying a
// description (except "0" for success). The description field only appears
// in the top-level response envelope on errors; successful responses with
// code 0 never include a description at the envelope level (nested objects
// like Job may have a description field, which is unrelated).
func checkErrorResponse(statusCode int, respBody []byte) error {
	var rawErr struct {
		Code        any    `json:"code"`        // int or string
		Description string `json:"description"` // error message
		Message     string `json:"message"`     // legacy alias
	}
	if err := json.Unmarshal(respBody, &rawErr); err == nil {
		desc := rawErr.Description
		if desc == "" {
			desc = rawErr.Message
		}
		code := normalizeCode(rawErr.Code)
		// Per the spec, a description only ever appears on error responses,
		// so any non-empty description or non-success code is an error.
		if codeString(code) == "session" || isErrorCode(code) || desc != "" {
			return NewAPIError(code, desc, statusCode)
		}
	}
	if statusCode < 200 || statusCode >= 300 {
		const maxRawMsg = 200
		msg := string(respBody)
		if len(msg) > maxRawMsg {
			msg = msg[:maxRawMsg]
		}
		return NewAPIError(nil, msg, statusCode)
	}
	return nil
}

// normalizeCode converts a decoded JSON code value into an int for numeric
// codes (JSON numbers decode to float64 in an interface{}) and keeps string
// codes as-is. A raw int is also accepted (Go's default json decoder never
// produces int for an untyped interface{}, but we handle it for robustness).
func normalizeCode(code any) any {
	switch v := code.(type) {
	case float64:
		return int(v)
	case int:
		return v
	default:
		return code
	}
}

// isErrorCode reports whether a code value indicates an error. Numeric codes
// are errors when non-zero; string codes are errors when neither "0" nor
// "session" (any other string code, e.g. "test", is an error).
func isErrorCode(code any) bool {
	switch v := code.(type) {
	case int:
		return v != 0
	case string:
		return v != "session" && v != "0"
	default:
		return false
	}
}

// codeString returns the string form of a code value, if any.
func codeString(code any) string {
	if s, ok := code.(string); ok {
		return s
	}
	return ""
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

	q := url.Values{}
	for k, v := range params {
		q.Set(k, v)
	}
	return c.DoGET(ctx, path+"?"+q.Encode())
}

// Health checks if the xyops instance is reachable and healthy via ping/v1.
// Returns (true, nil) when the API responds with code 0, (false, nil)
// when the instance is reachable but not the primary (HTTP 503), and (false, err)
// for network failures or non-503 API errors.
func (c *Client) Health(ctx context.Context) (bool, error) {
	body, err := c.DoGET(ctx, "ping/v1")
	if err != nil {
		var apiErr *APIError
		if errors.As(err, &apiErr) && apiErr.HTTPStatus == http.StatusServiceUnavailable {
			// Reachable but not the primary instance.
			return false, nil
		}
		return false, err
	}

	var resp struct {
		Code int `json:"code"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return false, err
	}
	return resp.Code == 0, nil
}

func (c *Client) resolveURL(path string) string {
	// If path is absolute (/...), use it directly against the base URL.
	// Plain concatenation is intentional: url.JoinPath would percent-encode
	// the '?' in query strings that callers append (e.g. "get_tag/v1?id=t1").
	if strings.HasPrefix(path, "/") {
		return c.baseURL.String() + path
	}
	// Otherwise, prepend /api/app/
	return c.baseURL.String() + basePath + path
}
