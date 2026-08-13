package mcp

import (
	"context"

	"github.com/gumieri/xyops-integrations/sdk"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func registerWebHookTools(s *mcp.Server, c *sdk.Client) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_webhooks",
		Description: "List all webhooks",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, []sdk.WebHook, error) {
		return adaptResult(c.WebHooks.List(ctx))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_webhook",
		Description: "Get a webhook by ID",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in struct {
		ID string `json:"id" jsonschema:"the webhook ID"`
	}) (*mcp.CallToolResult, *sdk.WebHook, error) {
		return adaptResult(c.WebHooks.Get(ctx, in.ID))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "create_webhook",
		Description: "Create a webhook",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in sdk.WebHook) (*mcp.CallToolResult, *sdk.WebHook, error) {
		return adaptResult(c.WebHooks.Create(ctx, &in))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "update_webhook",
		Description: "Update a webhook",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in sdk.WebHook) (*mcp.CallToolResult, *sdk.WebHook, error) {
		return adaptResult(c.WebHooks.Update(ctx, &in))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "delete_webhook",
		Description: "Delete a webhook by ID",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in struct {
		ID string `json:"id" jsonschema:"the webhook ID"`
	}) (*mcp.CallToolResult, string, error) {
		return adaptResult("deleted", c.WebHooks.Delete(ctx, in.ID))
	})
}
