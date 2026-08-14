package sdk

import (
	"context"
	"fmt"
	"strconv"
	"strings"
)

// TriggerOpts configures a run_event/v1 trigger call.
type TriggerOpts struct {
	EventID    string
	EventTitle string
	Params     map[string]any
	Tags       []string
	Input      map[string]any
	Test       bool
}

// SearchOpts configures search_jobs/v1.
type SearchOpts struct {
	Query   string // API uses "query", not "q"
	Limit   int
	Offset  int
	SortBy  string // API uses "sort_by", not "sort"
	SortDir int    // API uses "sort_dir" (1 for ascending, -1 for descending)
	Select  []string
	Verbose bool
}

// ActiveJobOpts configures get_active_jobs/v1.
type ActiveJobOpts struct {
	EventID string // Note: API filter name is "event", not "event_id"
	Limit   int
	Offset  int
}

// JobsService
func (s *JobsService) Trigger(ctx context.Context, eventID string, params map[string]any) (*TriggerResult, error) {
	opts := TriggerOpts{EventID: eventID, Params: params}
	return s.TriggerWithOpts(ctx, opts)
}

// TriggerWithOpts dispatches an event execution on demand via run_event/v1.
func (s *JobsService) TriggerWithOpts(ctx context.Context, opts TriggerOpts) (*TriggerResult, error) {
	if opts.EventID == "" && opts.EventTitle == "" {
		return nil, fmt.Errorf("trigger event: at least one of EventID or EventTitle is required")
	}
	payload := map[string]any{}
	if opts.EventID != "" {
		payload["id"] = opts.EventID
	}
	if opts.EventTitle != "" {
		payload["title"] = opts.EventTitle
	}
	if opts.Params != nil {
		payload["params"] = opts.Params
	}
	if opts.Tags != nil {
		payload["tags"] = opts.Tags
	}
	if opts.Input != nil {
		payload["input"] = opts.Input
	}
	if opts.Test {
		payload["test"] = true
	}

	body, err := s.client.DoPOST(ctx, "run_event/v1", payload)
	if err != nil {
		return nil, err
	}

	var resp TriggerResult
	if err := unmarshalJSON(body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetActive lists active jobs with optional filters.
func (s *JobsService) GetActive(ctx context.Context, opts *ActiveJobOpts) ([]Job, error) {
	params := map[string]string{}
	if opts != nil {
		if opts.EventID != "" {
			params["event"] = opts.EventID // API uses "event", not "event_id"
		}
		if opts.Limit > 0 {
			params["limit"] = strconv.Itoa(opts.Limit)
		}
		if opts.Offset > 0 {
			params["offset"] = strconv.Itoa(opts.Offset)
		}
	}

	body, err := s.client.DoGETWithQuery(ctx, "get_active_jobs/v1", params)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Code int   `json:"code"`
		Rows []Job `json:"rows"`
	}
	if err := unmarshalJSON(body, &resp); err != nil {
		return nil, err
	}
	return resp.Rows, nil
}

func (s *JobsService) GetActiveSummary(ctx context.Context) (map[string]any, error) {
	body, err := s.client.DoGET(ctx, "get_active_job_summary/v1")
	if err != nil {
		return nil, err
	}

	var result map[string]any
	if err := unmarshalJSON(body, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *JobsService) GetWorkflowSummary(ctx context.Context, jobID string) (map[string]any, error) {
	body, err := s.client.DoGETWithQuery(ctx, "get_workflow_job_summary/v1", map[string]string{"workflow.job": jobID})
	if err != nil {
		return nil, err
	}

	var result map[string]any
	if err := unmarshalJSON(body, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *JobsService) Get(ctx context.Context, jobID string) (*Job, error) {
	body, err := s.client.DoGETWithQuery(ctx, "get_job/v1", map[string]string{"id": jobID})
	if err != nil {
		return nil, err
	}

	var resp struct {
		Code int  `json:"code"`
		Job  *Job `json:"job"`
	}
	if err := unmarshalJSON(body, &resp); err != nil {
		return nil, err
	}
	if resp.Job == nil {
		return nil, fmt.Errorf("get job %s: job not found", jobID)
	}
	return resp.Job, nil
}

func (s *JobsService) GetBatch(ctx context.Context, ids []string, verbose bool) ([]Job, error) {
	if len(ids) == 0 {
		return nil, fmt.Errorf("get jobs: at least one job ID is required")
	}
	payload := map[string]any{"ids": ids}
	if verbose {
		payload["verbose"] = true
	}

	body, err := s.client.DoPOST(ctx, "get_jobs/v1", payload)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Code int   `json:"code"`
		Jobs []Job `json:"jobs"`
	}
	if err := unmarshalJSON(body, &resp); err != nil {
		return nil, err
	}
	return resp.Jobs, nil
}

// Search queries completed jobs with PxQL syntax.
func (s *JobsService) Search(ctx context.Context, query string) ([]Job, error) {
	return s.SearchWithOpts(ctx, SearchOpts{Query: query})
}

// SearchWithOpts queries completed jobs with full control.
func (s *JobsService) SearchWithOpts(ctx context.Context, opts SearchOpts) ([]Job, error) {
	params := map[string]string{}
	if opts.Query != "" {
		params["query"] = opts.Query
	}
	if opts.Limit > 0 {
		params["limit"] = strconv.Itoa(opts.Limit)
	}
	if opts.Offset > 0 {
		params["offset"] = strconv.Itoa(opts.Offset)
	}
	if opts.SortBy != "" {
		params["sort_by"] = opts.SortBy
	}
	if opts.SortDir != 0 {
		params["sort_dir"] = strconv.Itoa(opts.SortDir)
	}
	if opts.Verbose {
		params["verbose"] = "1"
	}
	if len(opts.Select) > 0 {
		params["select"] = strings.Join(opts.Select, ",")
	}

	body, err := s.client.DoGETWithQuery(ctx, "search_jobs/v1", params)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Code int   `json:"code"`
		Rows []Job `json:"rows"`
	}
	if err := unmarshalJSON(body, &resp); err != nil {
		return nil, err
	}
	return resp.Rows, nil
}

func (s *JobsService) GetEventHistory(ctx context.Context, eventID string) (map[string]any, error) {
	body, err := s.client.DoGETWithQuery(ctx, "get_event_history/v1", map[string]string{"id": eventID})
	if err != nil {
		return nil, err
	}

	var result map[string]any
	if err := unmarshalJSON(body, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *JobsService) Abort(ctx context.Context, jobID string) error {
	_, err := s.client.DoPOST(ctx, "abort_job/v1", map[string]string{"id": jobID})
	if err != nil {
		return fmt.Errorf("abort job %s: %w", jobID, err)
	}
	return nil
}

func (s *JobsService) Resume(ctx context.Context, jobID string, params map[string]any) error {
	payload := map[string]any{"id": jobID}
	if params != nil {
		payload["params"] = params
	}
	_, err := s.client.DoPOST(ctx, "resume_job/v1", payload)
	if err != nil {
		return fmt.Errorf("resume job %s: %w", jobID, err)
	}
	return nil
}

func (s *JobsService) Delete(ctx context.Context, jobID string) error {
	_, err := s.client.DoPOST(ctx, "delete_job/v1", map[string]string{"id": jobID})
	if err != nil {
		return fmt.Errorf("delete job %s: %w", jobID, err)
	}
	return nil
}

func (s *JobsService) FlushQueue(ctx context.Context, eventID string) error {
	_, err := s.client.DoPOST(ctx, "flush_event_queue/v1", map[string]string{"id": eventID})
	if err != nil {
		return fmt.Errorf("flush event queue %s: %w", eventID, err)
	}
	return nil
}

func (s *JobsService) ManageTags(ctx context.Context, jobID string, tags []string) error {
	if tags == nil {
		// Replace the job's tags with an empty array; omitting the field
		// is undefined by the API, so we send an explicit empty slice.
		tags = []string{}
	}
	_, err := s.client.DoPOST(ctx, "manage_job_tags/v1", map[string]any{"id": jobID, "tags": tags})
	if err != nil {
		return fmt.Errorf("manage job tags %s: %w", jobID, err)
	}
	return nil
}

// GetLog returns a job's log as raw text (text/plain, not JSON).
// Returns ("", nil) when there is no log.
func (s *JobsService) GetLog(ctx context.Context, jobID string) (string, error) {
	body, err := s.client.DoGETWithQuery(ctx, "get_job_log/v1", map[string]string{"id": jobID})
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// ViewLog returns a job's log as raw text using a shareable token.
func (s *JobsService) ViewLog(ctx context.Context, jobID, token string) (string, error) {
	body, err := s.client.DoGETWithQuery(ctx, "view_job_log/v1", map[string]string{"id": jobID, "t": token})
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// DownloadLog returns a job's log as raw bytes using a shareable token.
func (s *JobsService) DownloadLog(ctx context.Context, jobID, token string) ([]byte, error) {
	return s.client.DoGETWithQuery(ctx, "download_job_log/v1", map[string]string{"id": jobID, "t": token})
}

// Update modifies a job record (admin only). Running jobs are updated
// in-memory; completed jobs are updated in persisted storage.
func (s *JobsService) Update(ctx context.Context, job *Job) error {
	_, err := s.client.DoPOST(ctx, "update_job/v1", job)
	if err != nil {
		return fmt.Errorf("update job %s: %w", job.ID, err)
	}
	return nil
}

// UpdateActive patches a live in-memory job record via update_active_job/v1.
func (s *JobsService) UpdateActive(ctx context.Context, jobID string, updates map[string]any) error {
	payload := map[string]any{"id": jobID}
	for k, v := range updates {
		payload[k] = v
	}
	_, err := s.client.DoPOST(ctx, "update_active_job/v1", payload)
	if err != nil {
		return fmt.Errorf("update active job %s: %w", jobID, err)
	}
	return nil
}

// SkipDelay moves a delayed job's deadline to now, skipping its delay.
func (s *JobsService) SkipDelay(ctx context.Context, jobID string) error {
	_, err := s.client.DoPOST(ctx, "job_skip_delay/v1", map[string]string{"id": jobID})
	if err != nil {
		return fmt.Errorf("skip delay job %s: %w", jobID, err)
	}
	return nil
}

// ToggleNotifyMe toggles the current user's completion-notification email for
// a job and returns the new state.
func (s *JobsService) ToggleNotifyMe(ctx context.Context, jobID string) (bool, error) {
	body, err := s.client.DoPOST(ctx, "job_toggle_notify_me/v1", map[string]string{"id": jobID})
	if err != nil {
		return false, err
	}

	var resp struct {
		Code    int  `json:"code"`
		Enabled bool `json:"enabled"`
	}
	if err := unmarshalJSON(body, &resp); err != nil {
		return false, err
	}
	return resp.Enabled, nil
}

// ManageTickets replaces the tickets attached to a completed job.
func (s *JobsService) ManageTickets(ctx context.Context, jobID string, tickets []string) error {
	if tickets == nil {
		tickets = []string{}
	}
	_, err := s.client.DoPOST(ctx, "manage_job_tickets/v1", map[string]any{"id": jobID, "tickets": tickets})
	if err != nil {
		return fmt.Errorf("manage job tickets %s: %w", jobID, err)
	}
	return nil
}
