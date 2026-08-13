package integration

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"testing"

	"github.com/gumieri/xyops-integrations/sdk"
)

func dumpRaw(t *testing.T, label string, raw []byte) {
	t.Helper()
	const maxLen = 4096
	if len(raw) > maxLen {
		raw = raw[:maxLen]
	}
	var buf bytes.Buffer
	if json.Indent(&buf, raw, "", "  ") == nil {
		t.Logf("%s raw: %s", label, buf.String())
	} else {
		t.Logf("%s raw: %s", label, string(raw))
	}
}

func TestReadOnly_Health(t *testing.T) {
	skipIfNoClient(t)
	healthy, err := xyopsClient.Health(t.Context())
	if err != nil {
		t.Fatalf("health: %v", err)
	}
	if !healthy {
		t.Error("expected healthy")
	}
	t.Logf("ping/v1 -> healthy=%v", healthy)
}

func TestReadOnly_ErrorEnvelope(t *testing.T) {
	skipIfNoClient(t)
	raw, err := xyopsClient.DoGET(t.Context(), "error/v1")
	if err == nil {
		t.Fatal("expected error from error/v1")
	}
	dumpRaw(t, "error/v1", raw)
	var apiErr *sdk.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *sdk.APIError, got %T", err)
	}
	if apiErr.Code != "test" {
		t.Errorf("expected string code %q, got %v (%T)", "test", apiErr.Code, apiErr.Code)
	}
	if apiErr.Description == "" {
		t.Error("expected non-empty description")
	}
	if apiErr.IsAuth() {
		t.Error("expected IsAuth()==false for code 'test'")
	}
	t.Logf("envelope: code=%v (%T), status=%d, description=%q", apiErr.Code, apiErr.Code, apiErr.HTTPStatus, apiErr.Description)
}

func TestReadOnly_SessionDetection(t *testing.T) {
	skipIfNoClient(t)
	baseURL := os.Getenv("XYOPS_BASE_URL")
	bad, err := sdk.New(baseURL, "bogus-key-bogus-key-000")
	if err != nil {
		t.Fatalf("new bad client: %v", err)
	}
	raw, err := bad.DoGET(t.Context(), "get_tags/v1")
	if err == nil {
		t.Fatal("expected auth error with bogus key")
	}
	dumpRaw(t, "bogus key get_tags/v1", raw)
	var apiErr *sdk.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *sdk.APIError, got %T", err)
	}
	if !apiErr.IsAuth() {
		t.Errorf("expected IsAuth()==true, got code=%v (%T) status=%d description=%q",
			apiErr.Code, apiErr.Code, apiErr.HTTPStatus, apiErr.Description)
	}
	t.Logf("bogus key -> code=%v (%T), status=%d", apiErr.Code, apiErr.Code, apiErr.HTTPStatus)
}

func TestReadOnly_Servers(t *testing.T) {
	skipIfNoClient(t)
	servers, err := xyopsClient.Servers.List(t.Context())
	if err != nil {
		t.Fatalf("list servers: %v", err)
	}
	masters, err := xyopsClient.Servers.Masters(t.Context())
	if err != nil {
		t.Fatalf("masters: %v", err)
	}
	t.Logf("servers=%d masters=%d", len(servers), len(masters))
	for _, s := range servers {
		b, _ := json.Marshal(s)
		t.Logf("server %s: %s", s.ID, b)
		if s.ID == "" {
			t.Error("server missing id")
		}
	}
	if len(masters) > 0 {
		for id, m := range masters {
			b, _ := json.Marshal(m)
			t.Logf("master %s: %s", id, b)
		}
	}
}

func TestReadOnly_GetMultiple(t *testing.T) {
	skipIfNoClient(t)
	t.Skip("get_multiple/v1 unavailable on v1.0.49: API returns 'Unsupported API: app/api_get_multiple'")
	_ = xyopsClient.GetMultiple // keep reference if endpoint returns in a later version
	resp, err := xyopsClient.GetMultiple.Get(t.Context(), sdk.GetMultipleOpts{
		Lists:   []string{"tags", "groups", "events", "categories", "plugins", "monitors", "alerts", "roles", "api_keys", "web_hooks", "channels", "buckets", "secrets"},
		Servers: true,
		State:   true,
		Stats:   true,
	})
	if err != nil {
		t.Fatalf("get_multiple: %v", err)
	}
	t.Logf("code=%d epoch=%d", resp.Code, resp.Epoch)
	t.Logf("tags=%d groups=%d events=%d categories=%d plugins=%d monitors=%d alerts=%d roles=%d api_keys=%d web_hooks=%d channels=%d buckets=%d secrets=%d servers=%d",
		len(resp.Tags), len(resp.Groups), len(resp.Events), len(resp.Categories), len(resp.Plugins),
		len(resp.Monitors), len(resp.Alerts), len(resp.Roles), len(resp.APIKeys), len(resp.WebHooks),
		len(resp.Channels), len(resp.Buckets), len(resp.Secrets), len(resp.Servers))
}

func TestReadOnly_SearchJobs(t *testing.T) {
	skipIfNoClient(t)
	jobs, err := xyopsClient.Jobs.SearchWithOpts(t.Context(), sdk.SearchOpts{
		Query:   "*",
		Limit:   1,
		SortBy:  "elapsed",
		SortDir: -1,
	})
	if err != nil {
		t.Fatalf("search_jobs: %v", err)
	}
	t.Logf("search returned %d jobs", len(jobs))
	if len(jobs) == 0 {
		t.Log("no jobs found; skipping per-job asserts")
		return
	}
	j := jobs[0]
	b, _ := json.Marshal(j)
	t.Logf("job[0]: %s", b)
	t.Logf("event=%q title=%q state=%q elapsed=%v code=%v type=%q category=%q plugin=%q",
		j.Event, j.Title, j.State, j.Elapsed, j.Code, j.Type, j.Category, j.Plugin)
	if j.ID == "" {
		t.Error("job missing id")
	}
	if j.Event == "" && j.Title == "" {
		t.Log("note: event/title empty (possibly pruned by verbose=false)")
	}
}

func TestReadOnly_JobModel(t *testing.T) {
	skipIfNoClient(t)
	jobs, err := xyopsClient.Jobs.SearchWithOpts(t.Context(), sdk.SearchOpts{Query: "*", Limit: 1})
	if err != nil {
		t.Fatalf("search_jobs: %v", err)
	}
	if len(jobs) == 0 {
		t.Skip("no jobs to compare")
	}
	id := jobs[0].ID
	raw, err := xyopsClient.DoGET(t.Context(), "get_job/v1?id="+id)
	if err != nil {
		t.Fatalf("get_job raw: %v", err)
	}
	dumpRaw(t, "get_job/v1", raw)
	var parsed sdk.Job
	if err := json.Unmarshal(raw, &parsed); err != nil {
		t.Fatalf("unmarshal into Job: %v", err)
	}
	t.Logf("parsed get_job: id=%q event=%q title=%q state=%q elapsed=%v started=%v code=%v",
		parsed.ID, parsed.Event, parsed.Title, parsed.State, parsed.Elapsed, parsed.Started, parsed.Code)
}

func TestReadOnly_GetBatch(t *testing.T) {
	skipIfNoClient(t)
	jobs, err := xyopsClient.Jobs.SearchWithOpts(t.Context(), sdk.SearchOpts{Query: "*", Limit: 2})
	if err != nil {
		t.Fatalf("search_jobs: %v", err)
	}
	if len(jobs) == 0 {
		t.Skip("no jobs to batch-fetch")
	}
	var ids []string
	for _, j := range jobs {
		ids = append(ids, j.ID)
	}
	batch, err := xyopsClient.Jobs.GetBatch(t.Context(), ids, false)
	if err != nil {
		t.Fatalf("get_jobs: %v", err)
	}
	t.Logf("batch fetched %d jobs (requested %d)", len(batch), len(ids))
}

func TestReadOnly_GetActiveJobs(t *testing.T) {
	skipIfNoClient(t)
	jobs, err := xyopsClient.Jobs.GetActive(t.Context(), nil)
	if err != nil {
		t.Fatalf("get_active_jobs: %v", err)
	}
	t.Logf("active jobs (no filter): %d", len(jobs))
	for _, j := range jobs {
		b, _ := json.Marshal(j)
		t.Logf("active job: %s", b)
	}

	events, err := xyopsClient.Events.List(t.Context())
	if err != nil {
		t.Fatalf("events list: %v", err)
	}
	t.Logf("events=%d", len(events))
	if len(events) == 0 {
		t.Log("no events; skipping event-filtered active jobs")
		return
	}
	eventID := events[0].ID
	filtered, err := xyopsClient.Jobs.GetActive(t.Context(), &sdk.ActiveJobOpts{EventID: eventID})
	if err != nil {
		t.Fatalf("get_active_jobs event filter: %v", err)
	}
	t.Logf("active jobs filtered by event=%s: %d", eventID, len(filtered))
}

func TestReadOnly_ListEnvelopes(t *testing.T) {
	skipIfNoClient(t)
	checks := []struct {
		name string
		fn   func() (int, error)
	}{
		{"tags", func() (int, error) { r, err := xyopsClient.Tags.List(t.Context()); return len(r), err }},
		{"groups", func() (int, error) { r, err := xyopsClient.Groups.List(t.Context()); return len(r), err }},
		{"events", func() (int, error) { r, err := xyopsClient.Events.List(t.Context()); return len(r), err }},
		{"categories", func() (int, error) { r, err := xyopsClient.Categories.List(t.Context()); return len(r), err }},
		{"plugins", func() (int, error) { r, err := xyopsClient.Plugins.List(t.Context()); return len(r), err }},
		{"monitors", func() (int, error) { r, err := xyopsClient.Monitors.List(t.Context()); return len(r), err }},
		{"alerts", func() (int, error) { r, err := xyopsClient.Alerts.List(t.Context()); return len(r), err }},
		{"roles", func() (int, error) { r, err := xyopsClient.Roles.List(t.Context()); return len(r), err }},
		{"web_hooks", func() (int, error) { r, err := xyopsClient.WebHooks.List(t.Context()); return len(r), err }},
		{"api_keys", func() (int, error) { r, err := xyopsClient.APIKeys.List(t.Context()); return len(r), err }},
	}
	for _, c := range checks {
		n, err := c.fn()
		if err != nil {
			t.Errorf("%s list: %v", c.name, err)
			continue
		}
		t.Logf("%s rows=%d", c.name, n)
	}
}

func TestReadOnly_AdminReads(t *testing.T) {
	skipIfNoClient(t)
	cfg, err := xyopsClient.Config.Get(t.Context())
	if err != nil {
		t.Fatalf("admin_get_config: %v", err)
	}
	b, _ := json.Marshal(cfg)
	t.Logf("config keys=%d", len(cfg))
	if len(cfg) > 0 {
		t.Logf("config: %s", b)
	}
}

func TestReadOnly_RealLimitAndActionValues(t *testing.T) {
	skipIfNoClient(t)
	events, err := xyopsClient.Events.List(t.Context())
	if err != nil {
		t.Fatalf("events list: %v", err)
	}
	seenLimits := map[string]int{}
	seenConditions := map[string]int{}
	for _, e := range events {
		for _, l := range e.Limits {
			seenLimits[l.Type]++
		}
		for _, a := range e.Actions {
			seenConditions[a.Condition]++
		}
	}
	t.Logf("observed limit types: %v", seenLimits)
	t.Logf("observed action conditions: %v", seenConditions)
	if len(events) > 0 {
		for _, e := range events[:min(len(events), 2)] {
			b, _ := json.Marshal(e)
			t.Logf("event[%s]: %s", e.ID, b)
		}
	}
}
