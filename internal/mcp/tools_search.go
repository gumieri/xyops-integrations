package mcp

import (
	"context"

	"github.com/gumieri/xyops-integrations/sdk"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func registerSearchTools(s *mcp.Server, c *sdk.Client) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "search_servers",
		Description: "Search servers by query",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in struct {
		Query   string `json:"query,omitempty" jsonschema:"search query (Unbase/PxQL syntax)"`
		Offset  int    `json:"offset,omitempty" jsonschema:"offset for pagination"`
		Limit   int    `json:"limit,omitempty" jsonschema:"limit the number of results"`
		SortBy  string `json:"sort_by,omitempty" jsonschema:"sort field"`
		SortDir int    `json:"sort_dir,omitempty" jsonschema:"1 ascending, -1 descending"`
	}) (*mcp.CallToolResult, ListResult[sdk.Server], error) {
		opts := sdk.SearchQueryOpts{
			Query:   in.Query,
			Offset:  in.Offset,
			Limit:   in.Limit,
			SortBy:  in.SortBy,
			SortDir: in.SortDir,
		}
		return adaptListResult(c.Search.SearchServers(ctx, opts))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "search_alerts",
		Description: "Search alert invocations by query",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in struct {
		Query   string `json:"query,omitempty" jsonschema:"search query (Unbase/PxQL syntax)"`
		Offset  int    `json:"offset,omitempty" jsonschema:"offset for pagination"`
		Limit   int    `json:"limit,omitempty" jsonschema:"limit the number of results"`
		SortBy  string `json:"sort_by,omitempty" jsonschema:"sort field"`
		SortDir int    `json:"sort_dir,omitempty" jsonschema:"1 ascending, -1 descending"`
	}) (*mcp.CallToolResult, ListResult[sdk.AlertInvocation], error) {
		opts := sdk.SearchQueryOpts{
			Query:   in.Query,
			Offset:  in.Offset,
			Limit:   in.Limit,
			SortBy:  in.SortBy,
			SortDir: in.SortDir,
		}
		return adaptListResult(c.Search.SearchAlerts(ctx, opts))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "search_snapshots",
		Description: "Search server snapshots by query",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in struct {
		Query   string `json:"query,omitempty" jsonschema:"search query (Unbase/PxQL syntax)"`
		Offset  int    `json:"offset,omitempty" jsonschema:"offset for pagination"`
		Limit   int    `json:"limit,omitempty" jsonschema:"limit the number of results"`
		SortBy  string `json:"sort_by,omitempty" jsonschema:"sort field"`
		SortDir int    `json:"sort_dir,omitempty" jsonschema:"1 ascending, -1 descending"`
		Verbose bool   `json:"verbose,omitempty" jsonschema:"include heavy nested fields"`
	}) (*mcp.CallToolResult, ListResult[map[string]any], error) {
		opts := sdk.SearchQueryOpts{
			Query:   in.Query,
			Offset:  in.Offset,
			Limit:   in.Limit,
			SortBy:  in.SortBy,
			SortDir: in.SortDir,
			Verbose: in.Verbose,
		}
		return adaptListResult(c.Search.SearchSnapshots(ctx, opts))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "search_activity",
		Description: "Search the activity log by query",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in struct {
		Query   string `json:"query,omitempty" jsonschema:"search query (Unbase/PxQL syntax)"`
		Offset  int    `json:"offset,omitempty" jsonschema:"offset for pagination"`
		Limit   int    `json:"limit,omitempty" jsonschema:"limit the number of results"`
		SortBy  string `json:"sort_by,omitempty" jsonschema:"sort field"`
		SortDir int    `json:"sort_dir,omitempty" jsonschema:"1 ascending, -1 descending"`
	}) (*mcp.CallToolResult, ListResult[map[string]any], error) {
		opts := sdk.SearchQueryOpts{
			Query:   in.Query,
			Offset:  in.Offset,
			Limit:   in.Limit,
			SortBy:  in.SortBy,
			SortDir: in.SortDir,
		}
		return adaptListResult(c.Search.SearchActivity(ctx, opts))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "search_revision_history",
		Description: "Search revision history by type and query",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in struct {
		Type    string `json:"type" jsonschema:"revision type (e.g. servers, jobs, alerts)"`
		Query   string `json:"query,omitempty" jsonschema:"search query (Unbase/PxQL syntax)"`
		Offset  int    `json:"offset,omitempty" jsonschema:"offset for pagination"`
		Limit   int    `json:"limit,omitempty" jsonschema:"limit the number of results"`
		SortBy  string `json:"sort_by,omitempty" jsonschema:"sort field"`
		SortDir int    `json:"sort_dir,omitempty" jsonschema:"1 ascending, -1 descending"`
	}) (*mcp.CallToolResult, ListResult[map[string]any], error) {
		opts := sdk.SearchQueryOpts{
			Query:   in.Query,
			Offset:  in.Offset,
			Limit:   in.Limit,
			SortBy:  in.SortBy,
			SortDir: in.SortDir,
		}
		return adaptListResult(c.Search.SearchRevisionHistory(ctx, in.Type, opts))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "search_stat_history",
		Description: "Search stat history (one day per result)",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in struct {
		Offset     int    `json:"offset" jsonschema:"day offset"`
		Limit      int    `json:"limit" jsonschema:"number of days to return"`
		Path       string `json:"path,omitempty" jsonschema:"stat path filter"`
		KeyPrefix  string `json:"key_prefix,omitempty" jsonschema:"key prefix filter"`
		CurrentDay bool   `json:"current_day,omitempty" jsonschema:"include today's partial stats"`
	}) (*mcp.CallToolResult, ListResult[map[string]any], error) {
		return adaptListResult(c.Search.SearchStatHistory(ctx, in.Offset, in.Limit, in.Path, in.KeyPrefix, in.CurrentDay))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "bulk_export",
		Description: "Bulk-export an index to a file format",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in struct {
		Index    string `json:"index" jsonschema:"the index to export"`
		Format   string `json:"format" jsonschema:"export format (csv, json, ...)"`
		Query    string `json:"query,omitempty" jsonschema:"filter query"`
		Compress bool   `json:"compress,omitempty" jsonschema:"gzip the export"`
		Columns  string `json:"columns,omitempty" jsonschema:"comma-separated columns to export"`
	}) (*mcp.CallToolResult, map[string]string, error) {
		var columns []string
		if in.Columns != "" {
			columns = append(columns, splitComma(in.Columns)...)
		}
		data, err := c.Search.BulkExport(ctx, in.Index, columns, in.Format, in.Query, in.Compress)
		if err != nil {
			return nil, nil, err
		}
		return nil, map[string]string{"data": string(data)}, nil
	})

	registerMarketplaceTool(s, c)
}

func registerMarketplaceTool(s *mcp.Server, c *sdk.Client) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "marketplace",
		Description: "Search the plugin marketplace",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in struct {
		Query      string   `json:"query,omitempty" jsonschema:"search query"`
		Type       string   `json:"type,omitempty" jsonschema:"filter by type"`
		PluginType string   `json:"plugin_type,omitempty" jsonschema:"filter by plugin type"`
		Author     string   `json:"author,omitempty" jsonschema:"filter by author"`
		Status     string   `json:"status,omitempty" jsonschema:"filter by status"`
		License    string   `json:"license,omitempty" jsonschema:"filter by license"`
		Tags       []string `json:"tags,omitempty" jsonschema:"filter by tags"`
		Requires   []string `json:"requires,omitempty" jsonschema:"filter by requirements"`
		SortBy     string   `json:"sort_by,omitempty" jsonschema:"sort field"`
		SortDir    int      `json:"sort_dir,omitempty" jsonschema:"1 ascending, -1 descending"`
		Offset     int      `json:"offset,omitempty" jsonschema:"offset for pagination"`
		Limit      int      `json:"limit,omitempty" jsonschema:"limit the number of results"`
	}) (*mcp.CallToolResult, ListResult[map[string]any], error) {
		opts := sdk.MarketplaceOpts{
			Query:      in.Query,
			Type:       in.Type,
			PluginType: in.PluginType,
			Author:     in.Author,
			Status:     in.Status,
			License:    in.License,
			Tags:       in.Tags,
			Requires:   in.Requires,
			SortBy:     in.SortBy,
			SortDir:    in.SortDir,
			Offset:     in.Offset,
			Limit:      in.Limit,
		}
		return adaptListResult(c.Search.Marketplace(ctx, opts))
	})
}

func splitComma(s string) []string {
	var out []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == ',' {
			if i > start {
				out = append(out, s[start:i])
			}
			start = i + 1
		}
	}
	if start < len(s) {
		out = append(out, s[start:])
	}
	return out
}
