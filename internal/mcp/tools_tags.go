package mcp

import (
	"context"

	"github.com/gumieri/xyops-integrations/internal/client"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func registerTagTools(s *mcp.Server, c *client.Client) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_tags",
		Description: "List all tags",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, []client.Tag, error) {
		return wrapErr(c.Tags.List(ctx))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_tag",
		Description: "Get a tag by ID",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in struct {
		ID string `json:"id" jsonschema:"the tag ID"`
	}) (*mcp.CallToolResult, *client.Tag, error) {
		return wrapErr(c.Tags.Get(ctx, in.ID))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "create_tag",
		Description: "Create a tag",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in client.Tag) (*mcp.CallToolResult, *client.Tag, error) {
		return wrapErr(c.Tags.Create(ctx, &in))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "update_tag",
		Description: "Update a tag",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in client.Tag) (*mcp.CallToolResult, *client.Tag, error) {
		return wrapErr(c.Tags.Update(ctx, &in))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "delete_tag",
		Description: "Delete a tag by ID",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in struct {
		ID string `json:"id" jsonschema:"the tag ID"`
	}) (*mcp.CallToolResult, string, error) {
		return wrapErr("deleted", c.Tags.Delete(ctx, in.ID))
	})
}
