package mcp

import (
	"context"

	"github.com/gumieri/xyops-integrations/sdk"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func registerJobTools(s *mcp.Server, c *sdk.Client) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_job",
		Description: "Get a job by ID",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in struct {
		ID string `json:"id" jsonschema:"the job ID"`
	}) (*mcp.CallToolResult, *sdk.Job, error) {
		return adaptResult(c.Jobs.Get(ctx, in.ID))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_jobs_batch",
		Description: "Get multiple jobs by ID",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in struct {
		IDs     []string `json:"ids" jsonschema:"the job IDs"`
		Verbose bool     `json:"verbose,omitempty" jsonschema:"include verbose output"`
	}) (*mcp.CallToolResult, ListResult[sdk.Job], error) {
		return adaptListResult(c.Jobs.GetBatch(ctx, in.IDs, in.Verbose))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "search_jobs",
		Description: "Search jobs by query",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in struct {
		Query string `json:"query" jsonschema:"search query"`
	}) (*mcp.CallToolResult, ListResult[sdk.Job], error) {
		return adaptListResult(c.Jobs.Search(ctx, in.Query))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_active_jobs",
		Description: "Get currently active jobs",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in struct {
		EventID string `json:"event_id,omitempty" jsonschema:"filter by event ID"`
		Limit   int    `json:"limit,omitempty" jsonschema:"limit the number of results"`
		Offset  int    `json:"offset,omitempty" jsonschema:"offset for pagination"`
	}) (*mcp.CallToolResult, ListResult[sdk.Job], error) {
		opts := &sdk.ActiveJobOpts{
			EventID: in.EventID,
			Limit:   in.Limit,
			Offset:  in.Offset,
		}
		return adaptListResult(c.Jobs.GetActive(ctx, opts))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "trigger_event",
		Description: "Trigger an event",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in struct {
		EventID string         `json:"event_id" jsonschema:"the event ID to trigger"`
		Params  map[string]any `json:"params,omitempty" jsonschema:"event parameters"`
	}) (*mcp.CallToolResult, *sdk.TriggerResult, error) {
		return adaptResult(c.Jobs.Trigger(ctx, in.EventID, in.Params))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "abort_job",
		Description: "Abort a job",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in struct {
		JobID string `json:"job_id" jsonschema:"the job ID to abort"`
	}) (*mcp.CallToolResult, StatusResult, error) {
		return adaptStatusResult("aborted", c.Jobs.Abort(ctx, in.JobID))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "delete_job",
		Description: "Delete a job",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in struct {
		JobID string `json:"job_id" jsonschema:"the job ID to delete"`
	}) (*mcp.CallToolResult, StatusResult, error) {
		return adaptStatusResult("deleted", c.Jobs.Delete(ctx, in.JobID))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "flush_event_queue",
		Description: "Flush all queued jobs for an event",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in struct {
		EventID string `json:"event_id" jsonschema:"the event ID"`
	}) (*mcp.CallToolResult, StatusResult, error) {
		return adaptStatusResult("flushed", c.Jobs.FlushQueue(ctx, in.EventID))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "manage_job_tags",
		Description: "Manage job tags",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in struct {
		JobID string   `json:"job_id" jsonschema:"the job ID"`
		Tags  []string `json:"tags" jsonschema:"the tags to set"`
	}) (*mcp.CallToolResult, StatusResult, error) {
		return adaptStatusResult("updated", c.Jobs.ManageTags(ctx, in.JobID, in.Tags))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "resume_job",
		Description: "Resume a failed or paused job",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in struct {
		JobID  string         `json:"job_id" jsonschema:"the job ID to resume"`
		Params map[string]any `json:"params,omitempty" jsonschema:"optional params for the resumed job"`
	}) (*mcp.CallToolResult, StatusResult, error) {
		return adaptStatusResult("resumed", c.Jobs.Resume(ctx, in.JobID, in.Params))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_active_job_summary",
		Description: "Get a summary of active jobs (counts per event)",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, map[string]any, error) {
		return adaptResult(c.Jobs.GetActiveSummary(ctx))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_workflow_job_summary",
		Description: "Get summary stats for a workflow job",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in struct {
		JobID string `json:"job_id" jsonschema:"the workflow job ID"`
	}) (*mcp.CallToolResult, map[string]any, error) {
		return adaptResult(c.Jobs.GetWorkflowSummary(ctx, in.JobID))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_event_history",
		Description: "Get the event execution history for an event",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in struct {
		EventID string `json:"event_id" jsonschema:"the event ID"`
	}) (*mcp.CallToolResult, map[string]any, error) {
		return adaptResult(c.Jobs.GetEventHistory(ctx, in.EventID))
	})
}
