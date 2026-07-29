package client

import "context"

// ConfigService (singleton)
func (s *ConfigService) Get(ctx context.Context) (map[string]any, error) {
	body, err := s.client.DoPOST(ctx, "admin_get_config/v1", nil)
	if err != nil {
		return nil, err
	}

	var resp ConfigResponse
	if err := unmarshalJSON(body, &resp); err != nil {
		return nil, err
	}
	return resp.Config, nil
}

func (s *ConfigService) Update(ctx context.Context, partial map[string]any) (map[string]any, error) {
	body, err := s.client.DoPOST(ctx, "admin_update_config/v1", partial)
	if err != nil {
		return nil, err
	}

	if err := checkSuccessCode(body); err != nil {
		return nil, err
	}

	return s.Get(ctx)
}

// ServersService
func (s *ServersService) List(ctx context.Context) ([]Satellite, error) {
	body, err := s.client.DoPOST(ctx, "get_servers/v1", nil)
	if err != nil {
		return nil, err
	}

	var resp ServerListResponse
	if err := unmarshalJSON(body, &resp); err != nil {
		return nil, err
	}
	return resp.Servers, nil
}

// SatelliteService (bootstrap)
func (s *SatelliteService) Config(ctx context.Context, token string) (*SatelliteConfig, error) {
	body, err := s.client.DoGETWithQuery(ctx, "satellite/config", map[string]string{"t": token})
	if err != nil {
		return nil, err
	}

	var resp SatelliteConfigResponse
	if err := unmarshalJSON(body, &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// Core returns raw tarball bytes for a satellite binary.
func (s *SatelliteService) Core(ctx context.Context, token, os, arch string) ([]byte, error) {
	return s.client.DoGETWithQuery(ctx, "satellite/core", map[string]string{"t": token, "os": os, "arch": arch})
}

func (s *SatelliteService) Token(ctx context.Context) (string, error) {
	body, err := s.client.DoPOST(ctx, "get_satellite_token/v1", nil)
	if err != nil {
		return "", err
	}

	var resp struct {
		Code  int    `json:"code"`
		Token string `json:"token"`
	}
	if err := unmarshalJSON(body, &resp); err != nil {
		return "", err
	}
	return resp.Token, nil
}

func (s *SatelliteService) Releases(ctx context.Context) (map[string]any, error) {
	body, err := s.client.DoGET(ctx, "get_satellite_releases/v1")
	if err != nil {
		return nil, err
	}

	var result map[string]any
	if err := unmarshalJSON(body, &result); err != nil {
		return nil, err
	}
	return result, nil
}
