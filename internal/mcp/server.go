package mcp

import (
	"github.com/gumieri/xyops-integrations/internal/client"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// NewServer creates the xyops MCP server with all tools registered.
func NewServer(c *client.Client) (*mcp.Server, error) {
	s := mcp.NewServer(&mcp.Implementation{Name: "xyops", Version: "0.1.0"}, nil)

	registerTagTools(s, c)
	registerGroupTools(s, c)
	registerCategoryTools(s, c)
	registerPluginTools(s, c)
	registerEventTools(s, c)
	registerMonitorTools(s, c)
	registerAlertTools(s, c)
	registerRoleTools(s, c)
	registerAPIKeyTools(s, c)
	registerWebHookTools(s, c)
	registerServerTools(s, c)
	registerConfigTools(s, c)
	registerJobTools(s, c)

	return s, nil
}

// wrapErr adapts a (T, error) pair into the (result, output, error) shape
// expected by mcp.ToolHandlerFor.
func wrapErr[T any](v T, err error) (*mcp.CallToolResult, T, error) {
	return nil, v, err
}
