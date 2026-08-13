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
		_, _ = w.Write([]byte(`{"code": 0, "rows": [{"id": "j1", "event": "e1"}], "list": {"length": 1}}`))
	})
	defer cleanup()

	jobs, err := client.Jobs.Search(context.Background(), "event:e1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(jobs) != 1 {
		t.Errorf("expected 1 job, got %d", len(jobs))
	}
	if jobs[0].Event != "e1" {
		t.Errorf("expected event e1, got %q", jobs[0].Event)
	}
}

func TestJobsService_SearchWithOpts_QueryParams(t *testing.T) {
	_, client, cleanup := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("query") != "event:e1" {
			t.Errorf("expected query param 'event:e1', got %q", q.Get("query"))
		}
		if q.Get("sort_by") != "elapsed" {
			t.Errorf("expected sort_by 'elapsed', got %q", q.Get("sort_by"))
		}
		if q.Get("sort_dir") != "-1" {
			t.Errorf("expected sort_dir '-1', got %q", q.Get("sort_dir"))
		}
		if q.Get("limit") != "10" {
			t.Errorf("expected limit '10', got %q", q.Get("limit"))
		}
		if q.Get("q") != "" {
			t.Errorf("did not expect legacy 'q' param, got %q", q.Get("q"))
		}
		if q.Get("sort") != "" {
			t.Errorf("did not expect legacy 'sort' param, got %q", q.Get("sort"))
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"code": 0, "rows": [], "list": {"length": 0}}`))
	})
	defer cleanup()

	_, err := client.Jobs.SearchWithOpts(context.Background(), SearchOpts{
		Query:   "event:e1",
		Limit:   10,
		SortBy:  "elapsed",
		SortDir: -1,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestJobsService_GetActive_Filter(t *testing.T) {
	_, client, cleanup := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("event") != "e1" {
			t.Errorf("expected event filter 'e1', got %q", q.Get("event"))
		}
		if q.Get("event_id") != "" {
			t.Errorf("did not expect event_id param, got %q", q.Get("event_id"))
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"code": 0, "rows": [{"id": "j1", "state": "queued"}], "list": {"length": 1}}`))
	})
	defer cleanup()

	jobs, err := client.Jobs.GetActive(context.Background(), &ActiveJobOpts{EventID: "e1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(jobs) != 1 {
		t.Errorf("expected 1 job, got %d", len(jobs))
	}
}

func TestJobsService_GetWorkflowSummary_Filter(t *testing.T) {
	_, client, cleanup := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("workflow.job") != "jtop" {
			t.Errorf("expected workflow.job filter 'jtop', got %q", q.Get("workflow.job"))
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"code": 0, "nodes": {"nmhr8zbgjiv": 3}}`))
	})
	defer cleanup()

	sum, err := client.Jobs.GetWorkflowSummary(context.Background(), "jtop")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	nodes, ok := sum["nodes"].(map[string]any)
	if !ok {
		t.Fatalf("expected nodes map, got %v", sum)
	}
	if len(nodes) != 1 {
		t.Errorf("expected 1 node, got %d", len(nodes))
	}
}

func TestJobsService_GetBatch_JobsKey(t *testing.T) {
	_, client, cleanup := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"code": 0, "jobs": [{"id": "j1", "title": "first"}, {"err": "Job not found"}]}`))
	})
	defer cleanup()

	jobs, err := client.Jobs.GetBatch(context.Background(), []string{"j1", "jmissing"}, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(jobs) != 2 {
		t.Errorf("expected 2 jobs, got %d", len(jobs))
	}
	if jobs[0].Title != "first" {
		t.Errorf("expected title 'first', got %q", jobs[0].Title)
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
