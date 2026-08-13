package sdk

import (
	"context"
	"sort"
)

// ConfigService
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
	_, err := s.client.DoPOST(ctx, "admin_update_config/v1", partial)
	if err != nil {
		return nil, err
	}

	return s.Get(ctx)
}

// ServersService
// List returns the connected worker servers. get_servers/v1 returns the
// servers object keyed by server ID, so the map is flattened and sorted by ID.
func (s *ServersService) List(ctx context.Context) ([]Server, error) {
	resp, err := s.listResponse(ctx)
	if err != nil {
		return nil, err
	}
	servers := make([]Server, 0, len(resp.Servers))
	for _, srv := range resp.Servers {
		servers = append(servers, srv)
	}
	sort.Slice(servers, func(i, j int) bool { return servers[i].ID < servers[j].ID })
	return servers, nil
}

// Masters returns the conductor/peer servers from get_servers/v1.
func (s *ServersService) Masters(ctx context.Context) (map[string]Conductor, error) {
	resp, err := s.listResponse(ctx)
	if err != nil {
		return nil, err
	}
	return resp.Masters, nil
}

// listResponse fetches and parses the ServerListResponse, shared by List
// and Masters to avoid duplicate HTTP calls.
func (s *ServersService) listResponse(ctx context.Context) (*ServerListResponse, error) {
	body, err := s.client.DoGET(ctx, "get_servers/v1")
	if err != nil {
		return nil, err
	}

	var resp ServerListResponse
	if err := unmarshalJSON(body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// SatelliteService (bootstrap)
func (s *SatelliteService) Config(ctx context.Context, token string) (*SatelliteConfig, error) {
	body, err := s.client.DoGETWithQuery(ctx, "satellite/config", map[string]string{"t": token})
	if err != nil {
		return nil, err
	}

	// This endpoint serves the raw config.json document, not a {code, data}
	// envelope.
	var cfg SatelliteConfig
	if err := unmarshalJSON(body, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
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
