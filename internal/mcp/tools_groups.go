package mcp

import (
	"context"

	"github.com/gumieri/xyops-integrations/sdk"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func registerGroupTools(s *mcp.Server, c *sdk.Client) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_groups",
		Description: "List all groups",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, ListResult[sdk.Group], error) {
		return adaptListResult(c.Groups.List(ctx))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_group",
		Description: "Get a group by ID",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in struct {
		ID string `json:"id" jsonschema:"the group ID"`
	}) (*mcp.CallToolResult, *sdk.Group, error) {
		return adaptResult(c.Groups.Get(ctx, in.ID))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "create_group",
		Description: "Create a group",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in sdk.Group) (*mcp.CallToolResult, *sdk.Group, error) {
		return adaptResult(c.Groups.Create(ctx, &in))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "update_group",
		Description: "Update a group",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in sdk.Group) (*mcp.CallToolResult, *sdk.Group, error) {
		return adaptResult(c.Groups.Update(ctx, &in))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "delete_group",
		Description: "Delete a group by ID",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in struct {
		ID string `json:"id" jsonschema:"the group ID"`
	}) (*mcp.CallToolResult, StatusResult, error) {
		return adaptStatusResult("deleted", c.Groups.Delete(ctx, in.ID))
	})
}
