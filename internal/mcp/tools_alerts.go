package mcp

import (
	"context"

	"github.com/gumieri/xyops-integrations/internal/client"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func registerAlertTools(s *mcp.Server, c *client.Client) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_alerts",
		Description: "List all alerts",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, []client.Alert, error) {
		return wrapErr(c.Alerts.List(ctx))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_alert",
		Description: "Get an alert by ID",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in struct {
		ID string `json:"id" jsonschema:"the alert ID"`
	}) (*mcp.CallToolResult, *client.Alert, error) {
		return wrapErr(c.Alerts.Get(ctx, in.ID))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "create_alert",
		Description: "Create an alert",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in client.Alert) (*mcp.CallToolResult, *client.Alert, error) {
		return wrapErr(c.Alerts.Create(ctx, &in))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "update_alert",
		Description: "Update an alert",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in client.Alert) (*mcp.CallToolResult, *client.Alert, error) {
		return wrapErr(c.Alerts.Update(ctx, &in))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "delete_alert",
		Description: "Delete an alert by ID",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in struct {
		ID string `json:"id" jsonschema:"the alert ID"`
	}) (*mcp.CallToolResult, string, error) {
		return wrapErr("deleted", c.Alerts.Delete(ctx, in.ID))
	})
}
