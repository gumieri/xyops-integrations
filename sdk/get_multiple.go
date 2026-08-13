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

// GetMultiple batch-fetches several list types in one call.
func (s *GetMultipleService) Get(ctx context.Context, opts GetMultipleOpts) (*GetMultipleResponse, error) {
	params := map[string]string{}
	if len(opts.Lists) > 0 {
		params["lists"] = strings.Join(opts.Lists, ",")
	}
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
