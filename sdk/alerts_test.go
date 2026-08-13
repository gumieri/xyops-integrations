package sdk

import (
	"context"
	"net/http"
	"testing"
)

func TestAlertsService_List(t *testing.T) {
	_, client, cleanup := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"code": 0, "rows": [{"id": "a1", "title": "cpu-alert", "enabled": true, "expression": "cpu > 80", "monitor_id": "m1"}], "list": {"length": 1}}`))
	})
	defer cleanup()

	alerts, err := client.Alerts.List(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(alerts) != 1 {
		t.Errorf("expected 1 alert, got %d", len(alerts))
	}
	if alerts[0].Title != "cpu-alert" {
		t.Errorf("expected title 'cpu-alert', got %q", alerts[0].Title)
	}
}

func TestAlertsService_Create(t *testing.T) {
	_, client, cleanup := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"code": 0, "alert": {"id": "a123", "title": "mem-alert", "enabled": true}}`))
	})
	defer cleanup()

	alert, err := client.Alerts.Create(context.Background(), &Alert{Title: "mem-alert", Enabled: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if alert.ID != "a123" {
		t.Errorf("expected ID 'a123', got %q", alert.ID)
	}
}

func TestAlertsService_Update(t *testing.T) {
	calls := 0
	_, client, cleanup := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
		calls++
		if r.URL.Path == "/api/app/update_alert/v1" {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"code": 0}`))
		} else {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"code": 0, "alert": {"id": "a1", "title": "updated"}}`))
		}
	})
	defer cleanup()

	alert, err := client.Alerts.Update(context.Background(), &Alert{ID: "a1", Title: "updated"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if alert.Title != "updated" {
		t.Errorf("expected title 'updated', got %q", alert.Title)
	}
	if calls != 2 {
		t.Errorf("expected 2 calls, got %d", calls)
	}
}

func TestAlertsService_Delete(t *testing.T) {
	_, client, cleanup := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"code": 0}`))
	})
	defer cleanup()

	if err := client.Alerts.Delete(context.Background(), "a1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
