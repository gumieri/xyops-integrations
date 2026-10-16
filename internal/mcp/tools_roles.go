package mcp

import (
	"context"

	"github.com/gumieri/xyops-integrations/sdk"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func registerRoleTools(s *mcp.Server, c *sdk.Client) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_roles",
		Description: "List all roles",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, ListResult[sdk.Role], error) {
		return adaptListResult(c.Roles.List(ctx))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_role",
		Description: "Get a role by ID",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in struct {
		ID string `json:"id" jsonschema:"the role ID"`
	}) (*mcp.CallToolResult, *sdk.Role, error) {
		return adaptResult(c.Roles.Get(ctx, in.ID))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "create_role",
		Description: "Create a role",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in sdk.Role) (*mcp.CallToolResult, *sdk.Role, error) {
		return adaptResult(c.Roles.Create(ctx, &in))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "update_role",
		Description: "Update a role",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in sdk.Role) (*mcp.CallToolResult, *sdk.Role, error) {
		return adaptResult(c.Roles.Update(ctx, &in))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "delete_role",
		Description: "Delete a role by ID",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in struct {
		ID string `json:"id" jsonschema:"the role ID"`
	}) (*mcp.CallToolResult, StatusResult, error) {
		return adaptStatusResult("deleted", c.Roles.Delete(ctx, in.ID))
	})
}
