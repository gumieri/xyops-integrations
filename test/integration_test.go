package integration

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/gumieri/xyops-integrations/sdk"
)

var xyopsClient *sdk.Client

func TestMain(m *testing.M) {
	baseURL := os.Getenv("XYOPS_BASE_URL")
	apiKey := os.Getenv("XYOPS_API_KEY")

	if baseURL != "" && apiKey != "" {
		c, err := sdk.New(baseURL, apiKey)
		if err != nil {
			panic(err)
		}
		xyopsClient = c
	}

	os.Exit(m.Run())
}

func skipIfNoClient(t *testing.T) {
	t.Helper()
	if xyopsClient == nil {
		t.Skip("XYOPS_BASE_URL and XYOPS_API_KEY not set; skipping integration test")
	}
}

// cleanupCtx returns a context for use inside t.Cleanup functions.
// t.Context() is canceled just before Cleanup-registered functions run
// (Go 1.24+), so HTTP calls in cleanups must not use it.
func cleanupCtx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 30*time.Second)
}

func TestTagCRUD(t *testing.T) {
	skipIfNoClient(t)

	created, err := xyopsClient.Tags.Create(t.Context(), &sdk.Tag{Title: "it-tag"})
	if err != nil {
		t.Fatalf("create tag: %v", err)
	}
	t.Cleanup(func() {
		ctx, cancel := cleanupCtx()
		defer cancel()
		if err := xyopsClient.Tags.Delete(ctx, created.ID); err != nil {
			t.Errorf("cleanup delete tag: %v", err)
		}
	})

	if created.ID == "" {
		t.Fatal("expected created tag to have an ID")
	}

	got, err := xyopsClient.Tags.Get(t.Context(), created.ID)
	if err != nil {
		t.Fatalf("get tag: %v", err)
	}
	if got.Title != "it-tag" {
		t.Errorf("expected title %q, got %q", "it-tag", got.Title)
	}

	updated, err := xyopsClient.Tags.Update(t.Context(), &sdk.Tag{ID: created.ID, Title: "it-tag-updated"})
	if err != nil {
		t.Fatalf("update tag: %v", err)
	}
	if updated.Title != "it-tag-updated" {
		t.Errorf("expected title %q, got %q", "it-tag-updated", updated.Title)
	}

	rows, err := xyopsClient.Tags.List(t.Context())
	if err != nil {
		t.Fatalf("list tags: %v", err)
	}
	found := false
	for _, tag := range rows {
		if tag.ID == created.ID {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("created tag %s not found in list", created.ID)
	}
}

func TestGroupCRUD(t *testing.T) {
	skipIfNoClient(t)

	created, err := xyopsClient.Groups.Create(t.Context(), &sdk.Group{
		Title:         "it-group",
		HostnameMatch: `.*\.example\.com`, // v1.0.90 validates hostname_match as a regular expression
		SortOrder:     10,
	})
	if err != nil {
		t.Fatalf("create group: %v", err)
	}
	t.Cleanup(func() {
		ctx, cancel := cleanupCtx()
		defer cancel()
		if err := xyopsClient.Groups.Delete(ctx, created.ID); err != nil {
			t.Errorf("cleanup delete group: %v", err)
		}
	})

	got, err := xyopsClient.Groups.Get(t.Context(), created.ID)
	if err != nil {
		t.Fatalf("get group: %v", err)
	}
	if got.HostnameMatch != `.*\.example\.com` {
		t.Errorf("expected hostname_match %q, got %q", `.*\.example\.com`, got.HostnameMatch)
	}
}

func TestEventCRUD(t *testing.T) {
	skipIfNoClient(t)

	// v1.0.90+ requires an existing category on event create (HTTP 400
	// "Missing parameter: category" otherwise); reuse the first one.
	cats, err := xyopsClient.Categories.List(t.Context())
	if err != nil {
		t.Fatalf("list categories: %v", err)
	}
	if len(cats) == 0 {
		t.Fatal("no categories on instance; event create requires one on v1.0.90+")
	}

	// Use the first event-type plugin on the instance rather than
	// hardcoding a plugin name.
	plugins, err := xyopsClient.Plugins.List(t.Context())
	if err != nil {
		t.Fatalf("list plugins: %v", err)
	}
	pluginID := ""
	for _, p := range plugins {
		if p.Type == "event" {
			pluginID = p.ID
			break
		}
	}
	if pluginID == "" {
		t.Skip("no event-type plugins on instance; cannot create event")
	}

	// v1.0.90 validates required plugin parameters at event create
	// (HTTP 400 "Parameter 'x' is required"). Plugin params are supplied
	// under Event.Params as an object keyed by param ID — not Fields.
	definition, err := xyopsClient.Plugins.Get(t.Context(), pluginID)
	if err != nil {
		t.Fatalf("get plugin %s: %v", pluginID, err)
	}
	params := make(map[string]any)
	for _, prm := range definition.Params {
		if prm.Required == nil || !*prm.Required {
			continue
		}
		params[prm.ID] = fmt.Sprintf("%v", prm.Value)
	}

	event := &sdk.Event{
		Title:    "it-event",
		Plugin:   pluginID,
		Targets:  []string{"it-group"},
		Category: cats[0].ID,
		Algo:     "random", // v1.0.90+ requires an explicit target algorithm on create
		Params:   params,
		Enabled:  true,
		Triggers: []sdk.Trigger{
			// v1.0.90 validates an interval trigger's starting date/time
			// (`start`, unix seconds or ms); v1.0.49 accepted duration alone.
			{Type: "interval", Duration: 60, Start: time.Now().Unix()},
		},
		Limits: []sdk.Limit{
			{Type: "job", Enabled: true, Amount: 1},
			{Type: "queue", Enabled: true, Amount: 0},
			{Type: "retry", Enabled: true, Amount: 0},
		},
		// No Actions in this test: v1.0.90 email actions require a valid
		// instance user list (Action.Users), which is instance-dependent.
		// Covered by docs/api-reference.md.
	}

	created, err := xyopsClient.Events.Create(t.Context(), event)
	if err != nil {
		t.Fatalf("create event: %v", err)
	}
	t.Cleanup(func() {
		ctx, cancel := cleanupCtx()
		defer cancel()
		if err := xyopsClient.Events.Delete(ctx, created.ID, &sdk.DeleteEventOpts{DeleteJobs: true}); err != nil {
			t.Errorf("cleanup delete event: %v", err)
		}
	})

	got, err := xyopsClient.Events.Get(t.Context(), created.ID)
	if err != nil {
		t.Fatalf("get event: %v", err)
	}
	if got.Title != "it-event" {
		t.Errorf("expected title %q, got %q", "it-event", got.Title)
	}
	if len(got.Limits) == 0 {
		t.Error("expected event to have limits injected by the API")
	}

	updated, err := xyopsClient.Events.Update(t.Context(), &sdk.Event{
		ID:      created.ID,
		Title:   "it-event-updated",
		Plugin:  pluginID,
		Enabled: false,
	})
	if err != nil {
		t.Fatalf("update event: %v", err)
	}
	if updated.Title != "it-event-updated" {
		t.Errorf("expected title %q, got %q", "it-event-updated", updated.Title)
	}
}

func TestMonitorCRUD(t *testing.T) {
	skipIfNoClient(t)

	created, err := xyopsClient.Monitors.Create(t.Context(), &sdk.Monitor{
		Title:        "it-monitor",
		Source:       "cpu",
		DataType:     "float",
		Display:      true,
		MinVertScale: 0,
	})
	if err != nil {
		t.Fatalf("create monitor: %v", err)
	}
	t.Cleanup(func() {
		ctx, cancel := cleanupCtx()
		defer cancel()
		if err := xyopsClient.Monitors.Delete(ctx, created.ID); err != nil {
			t.Errorf("cleanup delete monitor: %v", err)
		}
	})

	got, err := xyopsClient.Monitors.Get(t.Context(), created.ID)
	if err != nil {
		t.Fatalf("get monitor: %v", err)
	}
	if got.Source != "cpu" {
		t.Errorf("expected source %q, got %q", "cpu", got.Source)
	}
}

func TestAlertCRUD(t *testing.T) {
	skipIfNoClient(t)

	created, err := xyopsClient.Alerts.Create(t.Context(), &sdk.Alert{
		Title:      "it-alert",
		Enabled:    true,
		Expression: "value > 90",
		Message:    "cpu high",
	})
	if err != nil {
		t.Fatalf("create alert: %v", err)
	}
	t.Cleanup(func() {
		ctx, cancel := cleanupCtx()
		defer cancel()
		if err := xyopsClient.Alerts.Delete(ctx, created.ID); err != nil {
			t.Errorf("cleanup delete alert: %v", err)
		}
	})

	got, err := xyopsClient.Alerts.Get(t.Context(), created.ID)
	if err != nil {
		t.Fatalf("get alert: %v", err)
	}
	if got.Expression != "value > 90" {
		t.Errorf("expected expression %q, got %q", "value > 90", got.Expression)
	}
}

func TestAPIKeyWriteOnly(t *testing.T) {
	skipIfNoClient(t)

	created, err := xyopsClient.APIKeys.Create(t.Context(), &sdk.APIKey{
		Title:      "it-apikey",
		Active:     1,
		Privileges: sdk.BoolMap{"all": true},
	})
	if err != nil {
		t.Fatalf("create api key: %v", err)
	}
	t.Cleanup(func() {
		ctx, cancel := cleanupCtx()
		defer cancel()
		if err := xyopsClient.APIKeys.Delete(ctx, created.ID); err != nil {
			t.Errorf("cleanup delete api key: %v", err)
		}
	})

	// The issued key value is exposed on create only, on both versions:
	// v1.0.49 returns nothing on read (write-only), while v1.0.90+ returns
	// a server-side hash (64-char hex) — never the usable key issued above.
	got, err := xyopsClient.APIKeys.Get(t.Context(), created.ID)
	if err != nil {
		t.Fatalf("get api key: %v", err)
	}
	if got.Key == created.Key {
		t.Errorf("read leaked the usable key value issued at create: %q", got.Key)
	}
	if got.Key != "" && len(got.Key) == len(created.Key) {
		t.Errorf("read returned a same-format key value instead of a hash or empty: %q", got.Key)
	}
}

func TestJobTriggerAndSearch(t *testing.T) {
	skipIfNoClient(t)

	// Fires a real job on the target instance. Additionally gated because
	// run_event/v1 on v1.0.90 rejects manual runs with "Event does not
	// allow manual job runs" even for enabled events with a stored manual
	// trigger and an admin key — the server-side gate is under
	// investigation (docs/api-reference.md).
	if os.Getenv("XYOPS_JOB_TRIGGER") == "" {
		t.Skip("set XYOPS_JOB_TRIGGER=1 to fire a real job; run_event manual-run gate under investigation")
	}

	// v1.0.90+ requires a category on event create; reuse the first one.
	cats, err := xyopsClient.Categories.List(t.Context())
	if err != nil {
		t.Fatalf("list categories: %v", err)
	}
	if len(cats) == 0 {
		t.Fatal("no categories on instance; event create requires one on v1.0.90+")
	}

	// v1.0.90+ requires targets on an enabled event; use the first group.
	groups, err := xyopsClient.Groups.List(t.Context())
	if err != nil {
		t.Fatalf("list groups: %v", err)
	}
	if len(groups) == 0 {
		t.Skip("no groups on instance; cannot target a job trigger test")
	}

	event, err := xyopsClient.Events.Create(t.Context(), &sdk.Event{
		Title:    "it-job-event",
		Plugin:   "testplug",
		Category: cats[0].ID,
		Targets:  []string{groups[0].ID},
		Triggers: []sdk.Trigger{
			// Manual-run permission comes from an explicit manual trigger;
			// without it run_event is rejected ("does not allow manual job runs").
			{Type: "manual"},
		},
		Params:  map[string]any{"duration": "1"}, // keep the triggered job short
		Algo:    "random",                        // v1.0.90+ requires an explicit target algorithm on create
		Enabled: true,
	})
	if err != nil {
		t.Fatalf("create event: %v", err)
	}
	t.Cleanup(func() {
		ctx, cancel := cleanupCtx()
		defer cancel()
		if err := xyopsClient.Events.Delete(ctx, event.ID, &sdk.DeleteEventOpts{DeleteJobs: true}); err != nil {
			t.Errorf("cleanup delete event: %v", err)
		}
	})

	res, err := xyopsClient.Jobs.Trigger(t.Context(), event.ID, nil)
	if err != nil {
		t.Fatalf("trigger event: %v", err)
	}
	if res.ID == "" {
		t.Fatal("expected trigger to return a job ID")
	}

	jobs, err := xyopsClient.Jobs.GetActive(t.Context(), nil)
	if err != nil {
		t.Fatalf("get active jobs: %v", err)
	}
	if len(jobs) == 0 {
		t.Log("warning: no active jobs (event may have completed quickly)")
	}

	job, err := xyopsClient.Jobs.Get(t.Context(), res.ID)
	if err != nil {
		t.Fatalf("get job: %v", err)
	}
	if job.ID != res.ID {
		t.Errorf("expected job ID %s, got %s", res.ID, job.ID)
	}
}

func TestErrorHandlingDeleteNonExistent(t *testing.T) {
	skipIfNoClient(t)

	if err := xyopsClient.Tags.Delete(t.Context(), "nonexistent-id-xyz"); err == nil {
		t.Fatal("expected error deleting non-existent tag")
	}
}

func TestGetMultipleBatch(t *testing.T) {
	skipIfNoClient(t)

	resp, err := xyopsClient.GetMultiple.Get(t.Context(), sdk.GetMultipleOpts{})
	if err != nil {
		t.Fatalf("get_multiple: %v", err)
	}
	if resp.Tags == nil && resp.Events == nil && resp.Groups == nil {
		t.Log("get_multiple returned empty batch (expected on fresh instance)")
	}
}
