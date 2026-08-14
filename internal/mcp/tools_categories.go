package mcp

import (
	"context"

	"github.com/gumieri/xyops-integrations/sdk"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func registerCategoryTools(s *mcp.Server, c *sdk.Client) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_categories",
		Description: "List all categories",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, ListResult[sdk.Category], error) {
		return adaptListResult(c.Categories.List(ctx))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_category",
		Description: "Get a category by ID",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in struct {
		ID string `json:"id" jsonschema:"the category ID"`
	}) (*mcp.CallToolResult, *sdk.Category, error) {
		return adaptResult(c.Categories.Get(ctx, in.ID))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "create_category",
		Description: "Create a category",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in sdk.Category) (*mcp.CallToolResult, *sdk.Category, error) {
		return adaptResult(c.Categories.Create(ctx, &in))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "update_category",
		Description: "Update a category",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in sdk.Category) (*mcp.CallToolResult, *sdk.Category, error) {
		return adaptResult(c.Categories.Update(ctx, &in))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "delete_category",
		Description: "Delete a category by ID",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in struct {
		ID string `json:"id" jsonschema:"the category ID"`
	}) (*mcp.CallToolResult, string, error) {
		return adaptResult("deleted", c.Categories.Delete(ctx, in.ID))
	})
}
