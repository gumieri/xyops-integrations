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
}
