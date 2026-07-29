package client

import "context"

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
	body, err := s.client.DoPOST(ctx, "update_alert/v1", alert)
	if err != nil {
		return nil, err
	}

	if err := checkSuccessCode(body); err != nil {
		return nil, err
	}

	return s.Get(ctx, alert.ID)
}

func (s *AlertsService) Delete(ctx context.Context, id string) error {
	_, err := s.client.DoPOST(ctx, "delete_alert/v1", map[string]string{"id": id})
	return err
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
	body, err := s.client.DoPOST(ctx, "update_role/v1", role)
	if err != nil {
		return nil, err
	}

	if err := checkSuccessCode(body); err != nil {
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

	body, err := s.client.DoPOST(ctx, "update_api_key/v1", payload)
	if err != nil {
		return nil, err
	}

	if err := checkSuccessCode(body); err != nil {
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
	body, err := s.client.DoPOST(ctx, "update_web_hook/v1", hook)
	if err != nil {
		return nil, err
	}

	if err := checkSuccessCode(body); err != nil {
		return nil, err
	}

	return s.Get(ctx, hook.ID)
}

func (s *WebHooksService) Delete(ctx context.Context, id string) error {
	_, err := s.client.DoPOST(ctx, "delete_web_hook/v1", map[string]string{"id": id})
	return err
}
