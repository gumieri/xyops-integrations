package mcp

import (
	"context"

	"github.com/gumieri/xyops-integrations/sdk"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func registerEventTools(s *mcp.Server, c *sdk.Client) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_events",
		Description: "List all events",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, ListResult[sdk.Event], error) {
		return adaptListResult(c.Events.List(ctx))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_event",
		Description: "Get an event by ID",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in struct {
		ID string `json:"id" jsonschema:"the event ID"`
	}) (*mcp.CallToolResult, *sdk.Event, error) {
		return adaptResult(c.Events.Get(ctx, in.ID))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "create_event",
		Description: "Create an event",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in sdk.Event) (*mcp.CallToolResult, *sdk.Event, error) {
		return adaptResult(c.Events.Create(ctx, &in))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "update_event",
		Description: "Update an event",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in sdk.Event) (*mcp.CallToolResult, *sdk.Event, error) {
		return adaptResult(c.Events.Update(ctx, &in))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "delete_event",
		Description: "Delete an event by ID. Optionally delete its jobs.",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in struct {
		ID         string `json:"id" jsonschema:"the event ID"`
		DeleteJobs bool   `json:"delete_jobs,omitempty" jsonschema:"delete the event's jobs too"`
	}) (*mcp.CallToolResult, StatusResult, error) {
		opts := &sdk.DeleteEventOpts{DeleteJobs: in.DeleteJobs}
		return adaptStatusResult("deleted", c.Events.Delete(ctx, in.ID, opts))
	})
}
