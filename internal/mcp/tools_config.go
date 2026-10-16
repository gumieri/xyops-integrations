package mcp

import (
	"context"

	"github.com/gumieri/xyops-integrations/sdk"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func registerConfigTools(s *mcp.Server, c *sdk.Client) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_config",
		Description: "Get the xyops configuration",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, map[string]any, error) {
		return adaptResult(c.Config.Get(ctx))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "update_config",
		Description: "Update xyops configuration (partial merge)",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in struct {
		Patch map[string]any `json:"patch" jsonschema:"configuration keys to update"`
	}) (*mcp.CallToolResult, map[string]any, error) {
		return adaptResult(c.Config.Update(ctx, in.Patch))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_health",
		Description: "Check the conductor's health (ping)",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, HealthResult, error) {
		return adaptHealthResult(c.Health(ctx))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_multiple",
		Description: "Batch-fetch several list types in one call",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in struct {
		Lists       []string `json:"lists,omitempty" jsonschema:"list types to fetch (events, jobs, alerts, groups, monitors, plugins, channels, buckets, secrets, ...)"`
		Jobs        bool     `json:"jobs,omitempty" jsonschema:"include jobs"`
		Alerts      bool     `json:"alerts,omitempty" jsonschema:"include alerts"`
		State       bool     `json:"state,omitempty" jsonschema:"include state"`
		Stats       bool     `json:"stats,omitempty" jsonschema:"include stats"`
		Servers     bool     `json:"servers,omitempty" jsonschema:"include servers"`
		ServerCache bool     `json:"server_cache,omitempty" jsonschema:"include the server cache"`
	}) (*mcp.CallToolResult, *sdk.GetMultipleResponse, error) {
		opts := &sdk.GetMultipleOpts{
			Lists:       in.Lists,
			Jobs:        in.Jobs,
			Alerts:      in.Alerts,
			State:       in.State,
			Stats:       in.Stats,
			Servers:     in.Servers,
			ServerCache: in.ServerCache,
		}
		return adaptResult(c.GetMultiple.Get(ctx, *opts))
	})
}

// HealthResult is the object-shaped output of get_health: the conductor's
// ping is a bare boolean, but the MCP spec requires StructuredContent to be
// a JSON object.
type HealthResult struct {
	Healthy bool `json:"healthy"`
}

// adaptHealthResult adapts the (bool, error) Health pair into the
// (result, output, error) shape expected by mcp.ToolHandlerFor.
func adaptHealthResult(ok bool, err error) (*mcp.CallToolResult, HealthResult, error) {
	if err != nil {
		return nil, HealthResult{}, err
	}
	return nil, HealthResult{Healthy: ok}, nil
}
