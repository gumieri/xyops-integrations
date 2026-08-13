package mcp

import (
	"github.com/gumieri/xyops-integrations/sdk"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Version is the MCP server implementation version.
const Version = "0.1.0"

// NewServer creates the xyops MCP server with all tools registered.
func NewServer(c *sdk.Client) (*mcp.Server, error) {
	s := mcp.NewServer(&mcp.Implementation{Name: "xyops", Version: Version}, nil)

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

// adaptResult adapts a (T, error) pair into the (result, output, error) shape
// expected by mcp.ToolHandlerFor.
func adaptResult[T any](v T, err error) (*mcp.CallToolResult, T, error) {
	return nil, v, err
}
