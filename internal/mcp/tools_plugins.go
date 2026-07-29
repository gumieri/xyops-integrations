package mcp

import (
	"context"

	"github.com/gumieri/xyops-integrations/internal/client"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func registerPluginTools(s *mcp.Server, c *client.Client) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_plugins",
		Description: "List all plugins",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, []client.Plugin, error) {
		return wrapErr(c.Plugins.List(ctx))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_plugin",
		Description: "Get a plugin by ID",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in struct {
		ID string `json:"id" jsonschema:"the plugin ID"`
	}) (*mcp.CallToolResult, *client.Plugin, error) {
		return wrapErr(c.Plugins.Get(ctx, in.ID))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "create_plugin",
		Description: "Create a plugin",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in client.Plugin) (*mcp.CallToolResult, *client.Plugin, error) {
		return wrapErr(c.Plugins.Create(ctx, &in))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "update_plugin",
		Description: "Update a plugin",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in client.Plugin) (*mcp.CallToolResult, *client.Plugin, error) {
		return wrapErr(c.Plugins.Update(ctx, &in))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "delete_plugin",
		Description: "Delete a plugin by ID",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in struct {
		ID string `json:"id" jsonschema:"the plugin ID"`
	}) (*mcp.CallToolResult, string, error) {
		return wrapErr("deleted", c.Plugins.Delete(ctx, in.ID))
	})
}
