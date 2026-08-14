package sdk

import (
	"context"
	"fmt"
	"strings"
)

// AlertsService
func (s *AlertsService) List(ctx context.Context) ([]Alert, error) {
	body, err := s.client.DoGET(ctx, "get_alerts/v1")
	if err != nil {
		return nil, err
	}

	var resp AlertListResponse
	if err := unmarshalJSON(body, &resp); err != nil {
		return nil, err
	}
	return resp.Rows, nil
}

func (s *AlertsService) Get(ctx context.Context, id string) (*Alert, error) {
	body, err := s.client.DoGETWithQuery(ctx, "get_alert/v1", map[string]string{"id": id})
	if err != nil {
		return nil, err
	}

	var resp AlertResponse
	if err := unmarshalJSON(body, &resp); err != nil {
		return nil, err
	}
	return resp.Alert, nil
}

func (s *AlertsService) Create(ctx context.Context, alert *Alert) (*Alert, error) {
	if alert.ID != "" {
		return s.Update(ctx, alert)
	}

	body, err := s.client.DoPOST(ctx, "create_alert/v1", alert)
	if err != nil {
		return nil, err
	}

	var resp AlertResponse
	if err := unmarshalJSON(body, &resp); err != nil {
		return nil, err
	}
	return resp.Alert, nil
}

func (s *AlertsService) Update(ctx context.Context, alert *Alert) (*Alert, error) {
	_, err := s.client.DoPOST(ctx, "update_alert/v1", alert)
	if err != nil {
		return nil, err
	}

	return s.Get(ctx, alert.ID)
}

func (s *AlertsService) Delete(ctx context.Context, id string) error {
	_, err := s.client.DoPOST(ctx, "delete_alert/v1", map[string]string{"id": id})
	return err
}

// Invocations fetches alert invocations by ID via get_alert_invocations/v1.
// Each element is either a populated invocation or an error object
// {"err": "Alert invocation not found"}.
func (s *AlertsService) Invocations(ctx context.Context, ids []string) ([]map[string]any, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	body, err := s.client.DoGETWithQuery(ctx, "get_alert_invocations/v1", map[string]string{"ids": strings.Join(ids, ",")})
	if err != nil {
		return nil, err
	}

	var resp struct {
		Code   int              `json:"code"`
		Alerts []map[string]any `json:"alerts"`
	}
	if err := unmarshalJSON(body, &resp); err != nil {
		return nil, err
	}
	return resp.Alerts, nil
}

// ManageInvocationTickets replaces the tickets attached to an alert invocation
// (complete replacement array).
func (s *AlertsService) ManageInvocationTickets(ctx context.Context, invocationID string, tickets []string) error {
	if tickets == nil {
		tickets = []string{}
	}
	_, err := s.client.DoPOST(ctx, "manage_alert_invocation_tickets/v1", map[string]any{"id": invocationID, "tickets": tickets})
	return err
}

// DeleteInvocation permanently deletes an alert invocation.
func (s *AlertsService) DeleteInvocation(ctx context.Context, invocationID string) error {
	_, err := s.client.DoPOST(ctx, "delete_alert_invocation/v1", map[string]string{"id": invocationID})
	return err
}

// Test evaluates an alert expression against live server data via test_alert/v1.
// Returns whether the alert would fire and the rendered message.
func (s *AlertsService) Test(ctx context.Context, server, expression, message string) (result bool, rendered string, err error) {
	body, err := s.client.DoPOST(ctx, "test_alert/v1", map[string]string{
		"server":     server,
		"expression": expression,
		"message":    message,
	})
	if err != nil {
		return false, "", err
	}

	var resp struct {
		Code    int    `json:"code"`
		Result  bool   `json:"result"`
		Message string `json:"message"`
	}
	if err := unmarshalJSON(body, &resp); err != nil {
		return false, "", err
	}
	return resp.Result, resp.Message, nil
}

// RolesService
func (s *RolesService) List(ctx context.Context) ([]Role, error) {
	body, err := s.client.DoGET(ctx, "get_roles/v1")
	if err != nil {
		return nil, err
	}

	var resp RoleListResponse
	if err := unmarshalJSON(body, &resp); err != nil {
		return nil, err
	}
	return resp.Rows, nil
}

func (s *RolesService) Get(ctx context.Context, id string) (*Role, error) {
	body, err := s.client.DoGETWithQuery(ctx, "get_role/v1", map[string]string{"id": id})
	if err != nil {
		return nil, err
	}

	var resp RoleResponse
	if err := unmarshalJSON(body, &resp); err != nil {
		return nil, err
	}
	return resp.Role, nil
}

func (s *RolesService) Create(ctx context.Context, role *Role) (*Role, error) {
	if role.ID != "" {
		return s.Update(ctx, role)
	}

	body, err := s.client.DoPOST(ctx, "create_role/v1", role)
	if err != nil {
		return nil, err
	}

	var resp RoleResponse
	if err := unmarshalJSON(body, &resp); err != nil {
		return nil, err
	}
	return resp.Role, nil
}

func (s *RolesService) Update(ctx context.Context, role *Role) (*Role, error) {
	_, err := s.client.DoPOST(ctx, "update_role/v1", role)
	if err != nil {
		return nil, err
	}

	return s.Get(ctx, role.ID)
}

func (s *RolesService) Delete(ctx context.Context, id string) error {
	_, err := s.client.DoPOST(ctx, "delete_role/v1", map[string]string{"id": id})
	return err
}

// APIKeysService
func (s *APIKeysService) List(ctx context.Context) ([]APIKey, error) {
	body, err := s.client.DoGET(ctx, "get_api_keys/v1")
	if err != nil {
		return nil, err
	}

	var resp APIKeyListResponse
	if err := unmarshalJSON(body, &resp); err != nil {
		return nil, err
	}
	return resp.Rows, nil
}

// Get uses POST body with {"id":"..."} — the API rejects GET for api keys.
func (s *APIKeysService) Get(ctx context.Context, id string) (*APIKey, error) {
	body, err := s.client.DoPOST(ctx, "get_api_key/v1", map[string]string{"id": id})
	if err != nil {
		return nil, err
	}

	var resp APIKeyResponse
	if err := unmarshalJSON(body, &resp); err != nil {
		return nil, err
	}
	return resp.APIKey, nil
}

func (s *APIKeysService) Create(ctx context.Context, key *APIKey) (*APIKey, error) {
	if key.ID != "" {
		return s.Update(ctx, key)
	}

	body, err := s.client.DoPOST(ctx, "create_api_key/v1", key)
	if err != nil {
		return nil, err
	}

	var resp APIKeyResponse
	if err := unmarshalJSON(body, &resp); err != nil {
		return nil, err
	}
	// The API returns the plain key value only once, on create, in the
	// plain_key property. Surface it on the returned key's Key field.
	if resp.APIKey != nil && resp.PlainKey != "" {
		resp.APIKey.Key = resp.PlainKey
	}
	return resp.APIKey, nil
}

func (s *APIKeysService) Update(ctx context.Context, key *APIKey) (*APIKey, error) {
	// Key is immutable and write-only — NOT sent on update.
	payload := &APIKeyUpdatePayload{
		ID:          key.ID,
		Title:       key.Title,
		Description: key.Description,
		Active:      key.Active,
		Privileges:  key.Privileges,
		Roles:       key.Roles,
		Expires:     key.Expires,
	}

	_, err := s.client.DoPOST(ctx, "update_api_key/v1", payload)
	if err != nil {
		return nil, err
	}

	return s.Get(ctx, key.ID)
}

func (s *APIKeysService) Delete(ctx context.Context, id string) error {
	_, err := s.client.DoPOST(ctx, "delete_api_key/v1", map[string]string{"id": id})
	return err
}

// WebHooksService
func (s *WebHooksService) List(ctx context.Context) ([]WebHook, error) {
	body, err := s.client.DoGET(ctx, "get_web_hooks/v1")
	if err != nil {
		return nil, err
	}

	var resp WebHookListResponse
	if err := unmarshalJSON(body, &resp); err != nil {
		return nil, err
	}
	return resp.Rows, nil
}

func (s *WebHooksService) Get(ctx context.Context, id string) (*WebHook, error) {
	body, err := s.client.DoGETWithQuery(ctx, "get_web_hook/v1", map[string]string{"id": id})
	if err != nil {
		return nil, err
	}

	var resp WebHookResponse
	if err := unmarshalJSON(body, &resp); err != nil {
		return nil, err
	}
	return resp.WebHook, nil
}

func (s *WebHooksService) Create(ctx context.Context, hook *WebHook) (*WebHook, error) {
	if hook.ID != "" {
		return s.Update(ctx, hook)
	}

	body, err := s.client.DoPOST(ctx, "create_web_hook/v1", hook)
	if err != nil {
		return nil, err
	}

	var resp WebHookResponse
	if err := unmarshalJSON(body, &resp); err != nil {
		return nil, err
	}
	return resp.WebHook, nil
}

func (s *WebHooksService) Update(ctx context.Context, hook *WebHook) (*WebHook, error) {
	_, err := s.client.DoPOST(ctx, "update_web_hook/v1", hook)
	if err != nil {
		return nil, err
	}

	return s.Get(ctx, hook.ID)
}

func (s *WebHooksService) Delete(ctx context.Context, id string) error {
	_, err := s.client.DoPOST(ctx, "delete_web_hook/v1", map[string]string{"id": id})
	return err
}

// Test sends a test request through a webhook (or a temporary copy with the
// given overrides) via test_web_hook/v1. Nothing is saved. title, method, and
// url are required even for existing hooks.
func (s *WebHooksService) Test(ctx context.Context, hook *WebHook) (code any, description, details string, err error) {
	if hook == nil {
		return nil, "", "", fmt.Errorf("test web hook: hook is required")
	}
	payload := map[string]any{
		"id":     hook.ID,
		"title":  hook.Title,
		"method": hook.Method,
		"url":    hook.URL,
	}
	if len(hook.Headers) > 0 {
		payload["headers"] = hook.Headers
	}
	if hook.Body != "" {
		payload["body"] = hook.Body
	}
	payload["timeout"] = hook.Timeout
	payload["retries"] = hook.Retries
	payload["follow"] = hook.Follow
	payload["ssl_cert_bypass"] = hook.SSLCertBypass
	payload["max_per_day"] = hook.MaxPerDay

	body, err := s.client.DoPOST(ctx, "test_web_hook/v1", payload)
	if err != nil {
		return nil, "", "", err
	}

	var resp struct {
		Code   int `json:"code"`
		Result struct {
			Code        any    `json:"code"`
			Description string `json:"description"`
			Details     string `json:"details"`
		} `json:"result"`
	}
	if err := unmarshalJSON(body, &resp); err != nil {
		return nil, "", "", err
	}
	return resp.Result.Code, resp.Result.Description, resp.Result.Details, nil
}
