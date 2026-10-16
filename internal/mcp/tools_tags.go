package mcp

import (
	"context"

	"github.com/gumieri/xyops-integrations/sdk"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func registerTagTools(s *mcp.Server, c *sdk.Client) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_tags",
		Description: "List all tags",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, ListResult[sdk.Tag], error) {
		return adaptListResult(c.Tags.List(ctx))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_tag",
		Description: "Get a tag by ID",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in struct {
		ID string `json:"id" jsonschema:"the tag ID"`
	}) (*mcp.CallToolResult, *sdk.Tag, error) {
		return adaptResult(c.Tags.Get(ctx, in.ID))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "create_tag",
		Description: "Create a tag",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in sdk.Tag) (*mcp.CallToolResult, *sdk.Tag, error) {
		return adaptResult(c.Tags.Create(ctx, &in))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "update_tag",
		Description: "Update a tag",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in sdk.Tag) (*mcp.CallToolResult, *sdk.Tag, error) {
		return adaptResult(c.Tags.Update(ctx, &in))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "delete_tag",
		Description: "Delete a tag by ID",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in struct {
		ID string `json:"id" jsonschema:"the tag ID"`
	}) (*mcp.CallToolResult, StatusResult, error) {
		return adaptStatusResult("deleted", c.Tags.Delete(ctx, in.ID))
	})
}
