package mcp

import (
	"context"

	"github.com/gumieri/xyops-integrations/sdk"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func registerAPIKeyTools(s *mcp.Server, c *sdk.Client) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_api_keys",
		Description: "List all API keys",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, ListResult[sdk.APIKey], error) {
		return adaptListResult(c.APIKeys.List(ctx))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_api_key",
		Description: "Get an API key by ID",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in struct {
		ID string `json:"id" jsonschema:"the API key ID"`
	}) (*mcp.CallToolResult, *sdk.APIKey, error) {
		return adaptResult(c.APIKeys.Get(ctx, in.ID))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "create_api_key",
		Description: "Create an API key",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in sdk.APIKey) (*mcp.CallToolResult, *sdk.APIKey, error) {
		return adaptResult(c.APIKeys.Create(ctx, &in))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "update_api_key",
		Description: "Update an API key",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in sdk.APIKey) (*mcp.CallToolResult, *sdk.APIKey, error) {
		return adaptResult(c.APIKeys.Update(ctx, &in))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "delete_api_key",
		Description: "Delete an API key by ID",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in struct {
		ID string `json:"id" jsonschema:"the API key ID"`
	}) (*mcp.CallToolResult, StatusResult, error) {
		return adaptStatusResult("deleted", c.APIKeys.Delete(ctx, in.ID))
	})
}
