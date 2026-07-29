package mcp

import (
	"context"

	"github.com/gumieri/xyops-integrations/internal/client"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func registerServerTools(s *mcp.Server, c *client.Client) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_servers",
		Description: "List all satellite servers",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, []client.Satellite, error) {
		return wrapErr(c.Servers.List(ctx))
	})
}
