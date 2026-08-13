package mcp

import (
	"context"

	"github.com/gumieri/xyops-integrations/sdk"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func registerPluginTools(s *mcp.Server, c *sdk.Client) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_plugins",
		Description: "List all plugins",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, []sdk.Plugin, error) {
		return adaptResult(c.Plugins.List(ctx))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_plugin",
		Description: "Get a plugin by ID",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in struct {
		ID string `json:"id" jsonschema:"the plugin ID"`
	}) (*mcp.CallToolResult, *sdk.Plugin, error) {
		return adaptResult(c.Plugins.Get(ctx, in.ID))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "create_plugin",
		Description: "Create a plugin",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in sdk.Plugin) (*mcp.CallToolResult, *sdk.Plugin, error) {
		return adaptResult(c.Plugins.Create(ctx, &in))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "update_plugin",
		Description: "Update a plugin",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in sdk.Plugin) (*mcp.CallToolResult, *sdk.Plugin, error) {
		return adaptResult(c.Plugins.Update(ctx, &in))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "delete_plugin",
		Description: "Delete a plugin by ID",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in struct {
		ID string `json:"id" jsonschema:"the plugin ID"`
	}) (*mcp.CallToolResult, string, error) {
		return adaptResult("deleted", c.Plugins.Delete(ctx, in.ID))
	})
}
