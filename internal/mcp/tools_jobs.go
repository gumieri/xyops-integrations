package mcp

import (
	"context"

	"github.com/gumieri/xyops-integrations/internal/client"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func registerJobTools(s *mcp.Server, c *client.Client) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_job",
		Description: "Get a job by ID",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in struct {
		ID string `json:"id" jsonschema:"the job ID"`
	}) (*mcp.CallToolResult, *client.Job, error) {
		return wrapErr(c.Jobs.Get(ctx, in.ID))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_jobs_batch",
		Description: "Get multiple jobs by ID",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in struct {
		IDs     []string `json:"ids" jsonschema:"the job IDs"`
		Verbose bool     `json:"verbose,omitempty" jsonschema:"include verbose output"`
	}) (*mcp.CallToolResult, []client.Job, error) {
		return wrapErr(c.Jobs.GetBatch(ctx, in.IDs, in.Verbose))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "search_jobs",
		Description: "Search jobs by query",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in struct {
		Query string `json:"query" jsonschema:"search query"`
	}) (*mcp.CallToolResult, []client.Job, error) {
		return wrapErr(c.Jobs.Search(ctx, in.Query))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_active_jobs",
		Description: "Get currently active jobs",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in struct {
		EventID string `json:"event_id,omitempty" jsonschema:"filter by event ID"`
		Limit   int    `json:"limit,omitempty" jsonschema:"limit the number of results"`
		Offset  int    `json:"offset,omitempty" jsonschema:"offset for pagination"`
	}) (*mcp.CallToolResult, []client.Job, error) {
		opts := &client.ActiveJobOpts{
			EventID: in.EventID,
			Limit:   in.Limit,
			Offset:  in.Offset,
		}
		return wrapErr(c.Jobs.GetActive(ctx, opts))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "trigger_event",
		Description: "Trigger an event",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in struct {
		EventID string         `json:"event_id" jsonschema:"the event ID to trigger"`
		Params  map[string]any `json:"params,omitempty" jsonschema:"event parameters"`
	}) (*mcp.CallToolResult, *client.TriggerResult, error) {
		return wrapErr(c.Jobs.Trigger(ctx, in.EventID, in.Params))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "abort_job",
		Description: "Abort a job",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in struct {
		JobID string `json:"job_id" jsonschema:"the job ID to abort"`
	}) (*mcp.CallToolResult, string, error) {
		return wrapErr("aborted", c.Jobs.Abort(ctx, in.JobID))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "delete_job",
		Description: "Delete a job",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in struct {
		JobID string `json:"job_id" jsonschema:"the job ID to delete"`
	}) (*mcp.CallToolResult, string, error) {
		return wrapErr("deleted", c.Jobs.Delete(ctx, in.JobID))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "flush_event_queue",
		Description: "Flush all queued jobs for an event",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in struct {
		EventID string `json:"event_id" jsonschema:"the event ID"`
	}) (*mcp.CallToolResult, string, error) {
		return wrapErr("flushed", c.Jobs.FlushQueue(ctx, in.EventID))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "manage_job_tags",
		Description: "Manage job tags",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in struct {
		JobID string   `json:"job_id" jsonschema:"the job ID"`
		Tags  []string `json:"tags" jsonschema:"the tags to set"`
	}) (*mcp.CallToolResult, string, error) {
		return wrapErr("updated", c.Jobs.ManageTags(ctx, in.JobID, in.Tags))
	})
}
