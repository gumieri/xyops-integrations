package mcp

import (
	"context"

	"github.com/gumieri/xyops-integrations/sdk"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func registerAlertTools(s *mcp.Server, c *sdk.Client) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_alerts",
		Description: "List all alerts",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, []sdk.Alert, error) {
		return adaptResult(c.Alerts.List(ctx))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_alert",
		Description: "Get an alert by ID",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in struct {
		ID string `json:"id" jsonschema:"the alert ID"`
	}) (*mcp.CallToolResult, *sdk.Alert, error) {
		return adaptResult(c.Alerts.Get(ctx, in.ID))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "create_alert",
		Description: "Create an alert",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in sdk.Alert) (*mcp.CallToolResult, *sdk.Alert, error) {
		return adaptResult(c.Alerts.Create(ctx, &in))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "update_alert",
		Description: "Update an alert",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in sdk.Alert) (*mcp.CallToolResult, *sdk.Alert, error) {
		return adaptResult(c.Alerts.Update(ctx, &in))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "delete_alert",
		Description: "Delete an alert by ID",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in struct {
		ID string `json:"id" jsonschema:"the alert ID"`
	}) (*mcp.CallToolResult, string, error) {
		return adaptResult("deleted", c.Alerts.Delete(ctx, in.ID))
	})
}
