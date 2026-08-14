package mcp

import (
	"context"

	"github.com/gumieri/xyops-integrations/sdk"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func registerChannelTools(s *mcp.Server, c *sdk.Client) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_channels",
		Description: "List all channels",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, ListResult[sdk.Channel], error) {
		return adaptListResult(c.Channels.List(ctx))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_channel",
		Description: "Get a channel by ID",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in struct {
		ID string `json:"id" jsonschema:"the channel ID"`
	}) (*mcp.CallToolResult, *sdk.Channel, error) {
		return adaptResult(c.Channels.Get(ctx, in.ID))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "create_channel",
		Description: "Create a channel",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in sdk.Channel) (*mcp.CallToolResult, *sdk.Channel, error) {
		return adaptResult(c.Channels.Create(ctx, &in))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "update_channel",
		Description: "Update a channel",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in sdk.Channel) (*mcp.CallToolResult, *sdk.Channel, error) {
		return adaptResult(c.Channels.Update(ctx, &in))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "delete_channel",
		Description: "Delete a channel by ID",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in struct {
		ID string `json:"id" jsonschema:"the channel ID"`
	}) (*mcp.CallToolResult, string, error) {
		return adaptResult("deleted", c.Channels.Delete(ctx, in.ID))
	})
}
