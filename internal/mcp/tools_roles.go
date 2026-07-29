package mcp

import (
	"context"

	"github.com/gumieri/xyops-integrations/internal/client"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func registerRoleTools(s *mcp.Server, c *client.Client) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_roles",
		Description: "List all roles",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, []client.Role, error) {
		return wrapErr(c.Roles.List(ctx))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_role",
		Description: "Get a role by ID",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in struct {
		ID string `json:"id" jsonschema:"the role ID"`
	}) (*mcp.CallToolResult, *client.Role, error) {
		return wrapErr(c.Roles.Get(ctx, in.ID))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "create_role",
		Description: "Create a role",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in client.Role) (*mcp.CallToolResult, *client.Role, error) {
		return wrapErr(c.Roles.Create(ctx, &in))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "update_role",
		Description: "Update a role",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in client.Role) (*mcp.CallToolResult, *client.Role, error) {
		return wrapErr(c.Roles.Update(ctx, &in))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "delete_role",
		Description: "Delete a role by ID",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in struct {
		ID string `json:"id" jsonschema:"the role ID"`
	}) (*mcp.CallToolResult, string, error) {
		return wrapErr("deleted", c.Roles.Delete(ctx, in.ID))
	})
}
