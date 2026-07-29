package mcp

import (
	"context"

	"github.com/gumieri/xyops-integrations/internal/client"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func registerAPIKeyTools(s *mcp.Server, c *client.Client) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_api_keys",
		Description: "List all API keys",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, []client.APIKey, error) {
		return wrapErr(c.APIKeys.List(ctx))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_api_key",
		Description: "Get an API key by ID",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in struct {
		ID string `json:"id" jsonschema:"the API key ID"`
	}) (*mcp.CallToolResult, *client.APIKey, error) {
		return wrapErr(c.APIKeys.Get(ctx, in.ID))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "create_api_key",
		Description: "Create an API key",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in client.APIKey) (*mcp.CallToolResult, *client.APIKey, error) {
		return wrapErr(c.APIKeys.Create(ctx, &in))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "update_api_key",
		Description: "Update an API key",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in client.APIKey) (*mcp.CallToolResult, *client.APIKey, error) {
		return wrapErr(c.APIKeys.Update(ctx, &in))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "delete_api_key",
		Description: "Delete an API key by ID",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in struct {
		ID string `json:"id" jsonschema:"the API key ID"`
	}) (*mcp.CallToolResult, string, error) {
		return wrapErr("deleted", c.APIKeys.Delete(ctx, in.ID))
	})
}
