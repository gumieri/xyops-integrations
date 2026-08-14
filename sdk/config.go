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

// Get returns a server's record, host data, and online status via
// get_server/v1 (three sibling keys).
func (s *ServersService) Get(ctx context.Context, id string) (*Server, map[string]any, bool, error) {
	body, err := s.client.DoGETWithQuery(ctx, "get_server/v1", map[string]string{"id": id})
	if err != nil {
		return nil, nil, false, err
	}

	var resp struct {
		Code   int            `json:"code"`
		Server *Server        `json:"server"`
		Data   map[string]any `json:"data"`
		Online bool           `json:"online"`
	}
	if err := unmarshalJSON(body, &resp); err != nil {
		return nil, nil, false, err
	}
	return resp.Server, resp.Data, resp.Online, nil
}

// Update shallow-merges server properties via update_server/v1, then
// re-fetches the server.
func (s *ServersService) Update(ctx context.Context, server *Server) (*Server, error) {
	_, err := s.client.DoPOST(ctx, "update_server/v1", server)
	if err != nil {
		return nil, err
	}

	srv, _, _, err := s.Get(ctx, server.ID)
	return srv, err
}

// UpdateData shallow-merges (or replaces, when replace is true) a server's
// user data object and returns the complete resulting object.
func (s *ServersService) UpdateData(ctx context.Context, id string, data map[string]any, replace bool) (map[string]any, error) {
	payload := map[string]any{"id": id, "data": data}
	if replace {
		payload["replace"] = true
	}
	body, err := s.client.DoPOST(ctx, "update_server_data/v1", payload)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Code int            `json:"code"`
		Data map[string]any `json:"data"`
	}
	if err := unmarshalJSON(body, &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// Delete uninstalls and optionally deletes a server (admin only). Offline
// servers require history=true. Runs in the background.
func (s *ServersService) Delete(ctx context.Context, id string, history bool) error {
	payload := map[string]any{"id": id}
	if history {
		payload["history"] = true
	}
	_, err := s.client.DoPOST(ctx, "delete_server/v1", payload)
	return err
}

// Watch takes a snapshot once per minute for the given duration (seconds).
// A duration of 0 cancels an existing watch.
func (s *ServersService) Watch(ctx context.Context, id string, duration int) error {
	_, err := s.client.DoPOST(ctx, "watch_server/v1", map[string]any{"id": id, "duration": duration})
	return err
}

// Active returns servers currently connected to the conductor via
// get_active_servers/v1 (rows + list envelope).
func (s *ServersService) Active(ctx context.Context) ([]Server, error) {
	body, err := s.client.DoGET(ctx, "get_active_servers/v1")
	if err != nil {
		return nil, err
	}

	var resp struct {
		Code int      `json:"code"`
		Rows []Server `json:"rows"`
	}
	if err := unmarshalJSON(body, &resp); err != nil {
		return nil, err
	}
	return resp.Rows, nil
}

// Summaries returns aggregate server property counts keyed by field ID
// (os_platform, os_distro, os_release, os_arch, cpu_virt, cpu_brand, cpu_cores).
func (s *ServersService) Summaries(ctx context.Context) (map[string]map[string]int, error) {
	body, err := s.client.DoGET(ctx, "get_server_summaries/v1")
	if err != nil {
		return nil, err
	}

	var resp struct {
		Code      int                       `json:"code"`
		Summaries map[string]map[string]int `json:"summaries"`
	}
	if err := unmarshalJSON(body, &resp); err != nil {
		return nil, err
	}
	return resp.Summaries, nil
}

// CreateSnapshot takes an immediate snapshot of a server.
func (s *ServersService) CreateSnapshot(ctx context.Context, serverID string) (string, error) {
	body, err := s.client.DoPOST(ctx, "create_snapshot/v1", map[string]string{"server": serverID})
	if err != nil {
		return "", err
	}

	var resp struct {
		Code int    `json:"code"`
		ID   string `json:"id"`
	}
	if err := unmarshalJSON(body, &resp); err != nil {
		return "", err
	}
	return resp.ID, nil
}

// DeleteSnapshot permanently deletes a snapshot.
func (s *ServersService) DeleteSnapshot(ctx context.Context, snapshotID string) error {
	_, err := s.client.DoPOST(ctx, "delete_snapshot/v1", map[string]string{"id": snapshotID})
	return err
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
