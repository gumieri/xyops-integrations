package mcp

import (
	"context"

	"github.com/gumieri/xyops-integrations/sdk"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func registerServerTools(s *mcp.Server, c *sdk.Client) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_servers",
		Description: "List all satellite servers",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, ListResult[sdk.Server], error) {
		return adaptListResult(c.Servers.List(ctx))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_masters",
		Description: "List all conductor servers",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, map[string]sdk.Conductor, error) {
		return adaptResult(c.Servers.Masters(ctx))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_server",
		Description: "Get a server's record, host data, and online status by ID",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in struct {
		ID string `json:"id" jsonschema:"the server ID"`
	}) (*mcp.CallToolResult, struct {
		Server *sdk.Server    `json:"server"`
		Data   map[string]any `json:"data"`
		Online bool           `json:"online"`
	}, error) {
		server, data, online, err := c.Servers.Get(ctx, in.ID)
		return adaptResult(struct {
			Server *sdk.Server    `json:"server"`
			Data   map[string]any `json:"data"`
			Online bool           `json:"online"`
		}{Server: server, Data: data, Online: online}, err)
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "update_server",
		Description: "Update a server's hostname or group",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in sdk.Server) (*mcp.CallToolResult, *sdk.Server, error) {
		return adaptResult(c.Servers.Update(ctx, &in))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "update_server_data",
		Description: "Set or replace a server's custom data",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in struct {
		ID      string         `json:"id" jsonschema:"the server ID"`
		Data    map[string]any `json:"data" jsonschema:"the data to set"`
		Replace bool           `json:"replace,omitempty" jsonschema:"replace the data object instead of merging"`
	}) (*mcp.CallToolResult, map[string]any, error) {
		return adaptResult(c.Servers.UpdateData(ctx, in.ID, in.Data, in.Replace))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "delete_server",
		Description: "Delete a server",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in struct {
		ID      string `json:"id" jsonschema:"the server ID"`
		History bool   `json:"history,omitempty" jsonschema:"also delete the server's history"`
	}) (*mcp.CallToolResult, string, error) {
		return adaptResult("deleted", c.Servers.Delete(ctx, in.ID, in.History))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "watch_server",
		Description: "Watch a server for incoming jobs",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in struct {
		ID       string `json:"id" jsonschema:"the server ID"`
		Duration int    `json:"duration" jsonschema:"watch duration in seconds"`
	}) (*mcp.CallToolResult, string, error) {
		return adaptResult("watching", c.Servers.Watch(ctx, in.ID, in.Duration))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_active_servers",
		Description: "List servers currently watching for jobs",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, ListResult[sdk.Server], error) {
		return adaptListResult(c.Servers.Active(ctx))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_server_summaries",
		Description: "Get job count summaries per server",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, map[string]map[string]int, error) {
		return adaptResult(c.Servers.Summaries(ctx))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "create_server_snapshot",
		Description: "Create a server snapshot (trigger an external sync)",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in struct {
		ServerID string `json:"server_id" jsonschema:"the server ID"`
	}) (*mcp.CallToolResult, string, error) {
		return adaptResult(c.Servers.CreateSnapshot(ctx, in.ServerID))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "delete_server_snapshot",
		Description: "Delete a server snapshot",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in struct {
		SnapshotID string `json:"snapshot_id" jsonschema:"the snapshot ID"`
	}) (*mcp.CallToolResult, string, error) {
		return adaptResult("deleted", c.Servers.DeleteSnapshot(ctx, in.SnapshotID))
	})
}
