package mcp

import (
	"context"

	"github.com/gumieri/xyops-integrations/internal/client"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func registerCategoryTools(s *mcp.Server, c *client.Client) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_categories",
		Description: "List all categories",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, []client.Category, error) {
		return wrapErr(c.Categories.List(ctx))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_category",
		Description: "Get a category by ID",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in struct {
		ID string `json:"id" jsonschema:"the category ID"`
	}) (*mcp.CallToolResult, *client.Category, error) {
		return wrapErr(c.Categories.Get(ctx, in.ID))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "create_category",
		Description: "Create a category",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in client.Category) (*mcp.CallToolResult, *client.Category, error) {
		return wrapErr(c.Categories.Create(ctx, &in))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "update_category",
		Description: "Update a category",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in client.Category) (*mcp.CallToolResult, *client.Category, error) {
		return wrapErr(c.Categories.Update(ctx, &in))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "delete_category",
		Description: "Delete a category by ID",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in struct {
		ID string `json:"id" jsonschema:"the category ID"`
	}) (*mcp.CallToolResult, string, error) {
		return wrapErr("deleted", c.Categories.Delete(ctx, in.ID))
	})
}
