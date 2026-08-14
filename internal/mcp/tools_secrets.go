package mcp

import (
	"context"

	"github.com/gumieri/xyops-integrations/sdk"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func registerSecretTools(s *mcp.Server, c *sdk.Client) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_secrets",
		Description: "List all secrets",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, ListResult[sdk.Secret], error) {
		return adaptListResult(c.Secrets.List(ctx))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_secret",
		Description: "Get a secret by ID (metadata only)",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in struct {
		ID string `json:"id" jsonschema:"the secret ID"`
	}) (*mcp.CallToolResult, *sdk.Secret, error) {
		return adaptResult(c.Secrets.Get(ctx, in.ID))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "decrypt_secret",
		Description: "Decrypt a secret's field values (admin only)",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in struct {
		ID string `json:"id" jsonschema:"the secret ID"`
	}) (*mcp.CallToolResult, ListResult[sdk.SecretField], error) {
		return adaptListResult(c.Secrets.Decrypt(ctx, in.ID))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "create_secret",
		Description: "Create a secret (admin only). Fields are encrypted and write-only.",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in sdk.Secret) (*mcp.CallToolResult, *sdk.Secret, error) {
		return adaptResult(c.Secrets.Create(ctx, &in))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "update_secret",
		Description: "Update a secret (admin only). Fields are encrypted and write-only.",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in sdk.Secret) (*mcp.CallToolResult, *sdk.Secret, error) {
		return adaptResult(c.Secrets.Update(ctx, &in))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "delete_secret",
		Description: "Delete a secret by ID",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in struct {
		ID string `json:"id" jsonschema:"the secret ID"`
	}) (*mcp.CallToolResult, string, error) {
		return adaptResult("deleted", c.Secrets.Delete(ctx, in.ID))
	})
}
