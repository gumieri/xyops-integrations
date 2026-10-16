package mcp

import (
	"context"

	"github.com/gumieri/xyops-integrations/sdk"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func registerBucketTools(s *mcp.Server, c *sdk.Client) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_buckets",
		Description: "List all buckets",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, ListResult[sdk.Bucket], error) {
		return adaptListResult(c.Buckets.List(ctx))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_bucket",
		Description: "Get a bucket by ID, including its data and files",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in struct {
		ID string `json:"id" jsonschema:"the bucket ID"`
	}) (*mcp.CallToolResult, struct {
		Bucket *sdk.Bucket    `json:"bucket"`
		Data   map[string]any `json:"data"`
		Files  []sdk.File     `json:"files"`
	}, error) {
		bucket, data, files, err := c.Buckets.Get(ctx, in.ID)
		return adaptResult(struct {
			Bucket *sdk.Bucket    `json:"bucket"`
			Data   map[string]any `json:"data"`
			Files  []sdk.File     `json:"files"`
		}{Bucket: bucket, Data: data, Files: files}, err)
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "create_bucket",
		Description: "Create a bucket",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in sdk.Bucket) (*mcp.CallToolResult, *sdk.Bucket, error) {
		return adaptResult(c.Buckets.Create(ctx, &in))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "update_bucket",
		Description: "Update a bucket",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in sdk.Bucket) (*mcp.CallToolResult, *sdk.Bucket, error) {
		return adaptResult(c.Buckets.Update(ctx, &in))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "delete_bucket",
		Description: "Delete a bucket by ID",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in struct {
		ID string `json:"id" jsonschema:"the bucket ID"`
	}) (*mcp.CallToolResult, StatusResult, error) {
		return adaptStatusResult("deleted", c.Buckets.Delete(ctx, in.ID))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "write_bucket_data",
		Description: "Shallow-merge data into a bucket's data object",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in struct {
		ID    string         `json:"id" jsonschema:"the bucket ID"`
		Data  map[string]any `json:"data" jsonschema:"the data to merge"`
		Fetch bool           `json:"fetch,omitempty" jsonschema:"return the full merged data object"`
	}) (*mcp.CallToolResult, map[string]any, error) {
		return adaptResult(c.Buckets.WriteData(ctx, in.ID, in.Data, in.Fetch))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "delete_bucket_file",
		Description: "Delete files from a bucket by criteria (filename, path, date, size, username). All supplied criteria must match the same file.",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in struct {
		ID       string `json:"id" jsonschema:"the bucket ID"`
		Filename string `json:"filename,omitempty" jsonschema:"file criterion: filename"`
		Path     string `json:"path,omitempty" jsonschema:"file criterion: path"`
		Date     int64  `json:"date,omitempty" jsonschema:"file criterion: date (unix seconds)"`
		Size     int64  `json:"size,omitempty" jsonschema:"file criterion: size in bytes"`
		Username string `json:"username,omitempty" jsonschema:"file criterion: username"`
	}) (*mcp.CallToolResult, StatusResult, error) {
		criteria := map[string]any{}
		if in.Filename != "" {
			criteria["filename"] = in.Filename
		}
		if in.Path != "" {
			criteria["path"] = in.Path
		}
		if in.Date != 0 {
			criteria["date"] = in.Date
		}
		if in.Size != 0 {
			criteria["size"] = in.Size
		}
		if in.Username != "" {
			criteria["username"] = in.Username
		}
		return adaptStatusResult("deleted", c.Buckets.DeleteFile(ctx, in.ID, criteria))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "empty_bucket",
		Description: "Remove all files, all data, or both from a bucket. At least one of files or data must be true.",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in struct {
		ID          string `json:"id" jsonschema:"the bucket ID"`
		DeleteFiles bool   `json:"files,omitempty" jsonschema:"delete all files"`
		ClearData   bool   `json:"data,omitempty" jsonschema:"delete all data"`
	}) (*mcp.CallToolResult, StatusResult, error) {
		return adaptStatusResult("emptied", c.Buckets.Empty(ctx, in.ID, in.DeleteFiles, in.ClearData))
	})
}
