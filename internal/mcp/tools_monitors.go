package mcp

import (
	"context"

	"github.com/gumieri/xyops-integrations/sdk"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func registerMonitorTools(s *mcp.Server, c *sdk.Client) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_monitors",
		Description: "List all monitors",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, []sdk.Monitor, error) {
		return adaptResult(c.Monitors.List(ctx))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_monitor",
		Description: "Get a monitor by ID",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in struct {
		ID string `json:"id" jsonschema:"the monitor ID"`
	}) (*mcp.CallToolResult, *sdk.Monitor, error) {
		return adaptResult(c.Monitors.Get(ctx, in.ID))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "create_monitor",
		Description: "Create a monitor",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in sdk.Monitor) (*mcp.CallToolResult, *sdk.Monitor, error) {
		return adaptResult(c.Monitors.Create(ctx, &in))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "update_monitor",
		Description: "Update a monitor",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in sdk.Monitor) (*mcp.CallToolResult, *sdk.Monitor, error) {
		return adaptResult(c.Monitors.Update(ctx, &in))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "delete_monitor",
		Description: "Delete a monitor by ID",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in struct {
		ID string `json:"id" jsonschema:"the monitor ID"`
	}) (*mcp.CallToolResult, string, error) {
		return adaptResult("deleted", c.Monitors.Delete(ctx, in.ID))
	})
}
