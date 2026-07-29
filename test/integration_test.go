package integration

import (
	"os"
	"testing"

	"github.com/gumieri/xyops-integrations/internal/client"
)

var xyopsClient *client.Client

func TestMain(m *testing.M) {
	baseURL := os.Getenv("XYOPS_BASE_URL")
	apiKey := os.Getenv("XYOPS_API_KEY")

	if baseURL != "" && apiKey != "" {
		c, err := client.New(baseURL, apiKey)
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

func TestTagCRUD(t *testing.T) {
	skipIfNoClient(t)

	created, err := xyopsClient.Tags.Create(t.Context(), &client.Tag{Title: "it-tag"})
	if err != nil {
		t.Fatalf("create tag: %v", err)
	}
	t.Cleanup(func() {
		if err := xyopsClient.Tags.Delete(t.Context(), created.ID); err != nil {
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

	updated, err := xyopsClient.Tags.Update(t.Context(), &client.Tag{ID: created.ID, Title: "it-tag-updated"})
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

	created, err := xyopsClient.Groups.Create(t.Context(), &client.Group{
		Title:         "it-group",
		HostnameMatch: "*.example.com",
		SortOrder:     10,
	})
	if err != nil {
		t.Fatalf("create group: %v", err)
	}
	t.Cleanup(func() {
		if err := xyopsClient.Groups.Delete(t.Context(), created.ID); err != nil {
			t.Errorf("cleanup delete group: %v", err)
		}
	})

	got, err := xyopsClient.Groups.Get(t.Context(), created.ID)
	if err != nil {
		t.Fatalf("get group: %v", err)
	}
	if got.HostnameMatch != "*.example.com" {
		t.Errorf("expected hostname_match %q, got %q", "*.example.com", got.HostnameMatch)
	}

	if err := xyopsClient.Groups.Delete(t.Context(), created.ID); err != nil {
		t.Fatalf("delete group: %v", err)
	}
}

func TestEventCRUD(t *testing.T) {
	skipIfNoClient(t)

	event := &client.Event{
		Title:   "it-event",
		Plugin:  "demo_echo",
		Targets: []string{"it-group"},
		Enabled: true,
		Triggers: []client.Trigger{
			{Type: "interval", Duration: 60},
		},
		Limits: []client.Limit{
			{Type: "job", Enabled: true, Amount: 1},
			{Type: "queue", Enabled: true, Amount: 0},
			{Type: "retry", Enabled: true, Amount: 0},
		},
		Actions: []client.Action{
			{Condition: "failure", Type: "email", Email: "ops@example.com", Enabled: true},
		},
	}

	created, err := xyopsClient.Events.Create(t.Context(), event)
	if err != nil {
		t.Fatalf("create event: %v", err)
	}
	t.Cleanup(func() {
		if err := xyopsClient.Events.Delete(t.Context(), created.ID, &client.DeleteEventOpts{DeleteJobs: true}); err != nil {
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

	updated, err := xyopsClient.Events.Update(t.Context(), &client.Event{
		ID:      created.ID,
		Title:   "it-event-updated",
		Plugin:  "demo_echo",
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

	created, err := xyopsClient.Monitors.Create(t.Context(), &client.Monitor{
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
		if err := xyopsClient.Monitors.Delete(t.Context(), created.ID); err != nil {
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

	created, err := xyopsClient.Alerts.Create(t.Context(), &client.Alert{
		Title:      "it-alert",
		Enabled:    true,
		Expression: "value > 90",
		Message:    "cpu high",
	})
	if err != nil {
		t.Fatalf("create alert: %v", err)
	}
	t.Cleanup(func() {
		if err := xyopsClient.Alerts.Delete(t.Context(), created.ID); err != nil {
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

	created, err := xyopsClient.APIKeys.Create(t.Context(), &client.APIKey{
		Title:      "it-apikey",
		Active:     1,
		Privileges: map[string]int{"all": 1},
	})
	if err != nil {
		t.Fatalf("create api key: %v", err)
	}
	t.Cleanup(func() {
		if err := xyopsClient.APIKeys.Delete(t.Context(), created.ID); err != nil {
			t.Errorf("cleanup delete api key: %v", err)
		}
	})

	// The key value is write-only — the create response may include it,
	// but a subsequent Get must NOT leak it.
	got, err := xyopsClient.APIKeys.Get(t.Context(), created.ID)
	if err != nil {
		t.Fatalf("get api key: %v", err)
	}
	if got.Key != "" {
		t.Errorf("expected write-only key to be empty on read, got %q", got.Key)
	}
}

func TestJobTriggerAndSearch(t *testing.T) {
	skipIfNoClient(t)

	event, err := xyopsClient.Events.Create(t.Context(), &client.Event{
		Title:   "it-job-event",
		Plugin:  "demo_echo",
		Enabled: true,
	})
	if err != nil {
		t.Fatalf("create event: %v", err)
	}
	t.Cleanup(func() {
		if err := xyopsClient.Events.Delete(t.Context(), event.ID, &client.DeleteEventOpts{DeleteJobs: true}); err != nil {
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

	resp, err := xyopsClient.GetMultiple.Get(t.Context(), client.GetMultipleOpts{})
	if err != nil {
		t.Fatalf("get_multiple: %v", err)
	}
	if resp.Tags == nil && resp.Events == nil && resp.Groups == nil {
		t.Log("get_multiple returned empty batch (expected on fresh instance)")
	}
}
