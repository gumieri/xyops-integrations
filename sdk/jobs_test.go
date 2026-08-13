package sdk

import (
	"context"
	"net/http"
	"testing"
)

func TestJobsService_Trigger(t *testing.T) {
	_, client, cleanup := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"code": 0, "id": "j123"}`))
	})
	defer cleanup()

	result, err := client.Jobs.Trigger(context.Background(), "e1", map[string]any{"host": "db01"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ID != "j123" {
		t.Errorf("expected ID 'j123', got %q", result.ID)
	}
}

func TestJobsService_TriggerWithOpts_RequiresKey(t *testing.T) {
	_, client, cleanup := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("server should not be called when no match key is provided")
		w.WriteHeader(http.StatusOK)
	})
	defer cleanup()

	_, err := client.Jobs.TriggerWithOpts(context.Background(), TriggerOpts{})
	if err == nil {
		t.Fatal("expected error when neither EventID nor EventTitle is set")
	}
}

func TestJobsService_Search(t *testing.T) {
	_, client, cleanup := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"code": 0, "rows": [{"id": "j1", "event_id": "e1"}], "total": 1}`))
	})
	defer cleanup()

	jobs, err := client.Jobs.Search(context.Background(), "event_id:e1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(jobs) != 1 {
		t.Errorf("expected 1 job, got %d", len(jobs))
	}
}

func TestJobsService_Abort(t *testing.T) {
	_, client, cleanup := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"code": 0}`))
	})
	defer cleanup()

	if err := client.Jobs.Abort(context.Background(), "j1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestJobsService_Get(t *testing.T) {
	_, client, cleanup := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"code": 0, "job": {"id": "j1", "state": "completed"}}`))
	})
	defer cleanup()

	job, err := client.Jobs.Get(context.Background(), "j1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if job.State != "completed" {
		t.Errorf("expected state 'completed', got %q", job.State)
	}
}
