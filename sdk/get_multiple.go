package sdk

import (
	"context"
	"strings"
)

// GetMultipleOpts configures get_multiple/v1.
type GetMultipleOpts struct {
	Lists       []string
	Jobs        bool
	Alerts      bool
	State       bool
	Stats       bool
	Servers     bool
	ServerCache bool
}

// defaultGetMultipleLists is the list set sent when opts.Lists is empty.
// v1.0.90+ rejects get_multiple without a `lists` parameter (HTTP 500,
// code=admin: "Missing required request property: lists"); v1.0.49 returned
// every list when the parameter was absent. Sending the standard set keeps
// empty-opts behavior working across both versions.
var defaultGetMultipleLists = []string{
	"tags", "groups", "events", "categories", "plugins", "monitors",
	"alerts", "roles", "api_keys", "web_hooks", "channels", "buckets", "secrets",
}

// GetMultiple batch-fetches several list types in one call.
func (s *GetMultipleService) Get(ctx context.Context, opts GetMultipleOpts) (*GetMultipleResponse, error) {
	params := map[string]string{}
	lists := opts.Lists
	if len(lists) == 0 {
		lists = defaultGetMultipleLists
	}
	params["lists"] = strings.Join(lists, ",")
	if opts.Jobs {
		params["jobs"] = "1"
	}
	if opts.Alerts {
		params["alerts"] = "1"
	}
	if opts.State {
		params["state"] = "1"
	}
	if opts.Stats {
		params["stats"] = "1"
	}
	if opts.Servers {
		params["servers"] = "1"
	}
	if opts.ServerCache {
		params["serverCache"] = "1"
	}

	body, err := s.client.DoGETWithQuery(ctx, "get_multiple/v1", params)
	if err != nil {
		return nil, err
	}

	var resp GetMultipleResponse
	if err := unmarshalJSON(body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
