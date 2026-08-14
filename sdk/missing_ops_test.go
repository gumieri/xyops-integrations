package sdk

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

// decodeRequestJSON decodes a request body into v for asserting request payloads.
func decodeRequestJSON(t *testing.T, r *http.Request, v any) {
	t.Helper()
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		t.Fatalf("decode request body: %v", err)
	}
}

func TestAlertsService_Invocations(t *testing.T) {
	_, client, cleanup := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/app/get_alert_invocations/v1" {
			t.Errorf("unexpected path %q", r.URL.Path)
		}
		if r.URL.Query().Get("ids") != "a1,a2" {
			t.Errorf("expected ids=a1,a2, got %q", r.URL.Query().Get("ids"))
		}
		_, _ = w.Write([]byte(`{"code": 0, "alerts": [{"id": "a1", "active": true, "alert": "cpu-high"}, {"err": "Alert invocation not found"}]}`))
	})
	defer cleanup()

	invocations, err := client.Alerts.Invocations(context.Background(), []string{"a1", "a2"})
	if err != nil {
		t.Fatalf("invocations: %v", err)
	}
	if len(invocations) != 2 {
		t.Fatalf("expected 2 invocations, got %d", len(invocations))
	}
	if invocations[0]["alert"] != "cpu-high" {
		t.Errorf("expected alert cpu-high, got %v", invocations[0]["alert"])
	}
	if _, ok := invocations[1]["err"]; !ok {
		t.Errorf("expected error object in second invocation, got %+v", invocations[1])
	}
}

func TestAlertsService_Test(t *testing.T) {
	_, client, cleanup := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/app/test_alert/v1" {
			t.Errorf("unexpected path %q", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"code": 0, "result": true, "message": "CPU high on server1"}`))
	})
	defer cleanup()

	result, msg, err := client.Alerts.Test(context.Background(), "server1", "cpu > 90", "CPU high")
	if err != nil {
		t.Fatalf("test: %v", err)
	}
	if !result {
		t.Error("expected result=true")
	}
	if msg != "CPU high on server1" {
		t.Errorf("unexpected message %q", msg)
	}
}

func TestMonitorsService_MultiUpdate(t *testing.T) {
	var body map[string]any
	_, client, cleanup := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/app/multi_update_monitor/v1" {
			t.Errorf("unexpected path %q", r.URL.Path)
		}
		decodeRequestJSON(t, r, &body)
		_, _ = w.Write([]byte(`{"code": 0}`))
	})
	defer cleanup()

	err := client.Monitors.MultiUpdate(context.Background(), []*Monitor{{ID: "m1", Title: "cpu", Source: "node_exporter"}})
	if err != nil {
		t.Fatalf("multi update: %v", err)
	}
	items, ok := body["items"].([]any)
	if !ok || len(items) != 1 {
		t.Fatalf("expected items array with 1 element, got %+v", body)
	}
	item := items[0].(map[string]any)
	if item["id"] != "m1" || item["source"] != "node_exporter" {
		t.Errorf("unexpected item: %+v", item)
	}
}

func TestMonitorsService_QuickMonData(t *testing.T) {
	_, client, cleanup := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/app/get_quickmon_data/v1" {
			t.Errorf("unexpected path %q", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"code": 0, "servers": {"s1": [{"title": "cpu"}]}}`))
	})
	defer cleanup()

	servers, err := client.Monitors.QuickMonData(context.Background(), "", "")
	if err != nil {
		t.Fatalf("quickmon: %v", err)
	}
	if len(servers["s1"]) != 1 {
		t.Errorf("expected 1 quickmon entry for s1, got %d", len(servers["s1"]))
	}
}

func TestPluginsService_TestScheduler(t *testing.T) {
	_, client, cleanup := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/app/test_scheduler_plugin/v1" {
			t.Errorf("unexpected path %q", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"code": 0, "data": {"items": [{"server": "s1"}]}, "stdout": "ok"}`))
	})
	defer cleanup()

	data, stdout, err := client.Plugins.TestScheduler(context.Background(), "p1", nil, "", 0)
	if err != nil {
		t.Fatalf("test scheduler: %v", err)
	}
	if stdout != "ok" {
		t.Errorf("expected stdout ok, got %q", stdout)
	}
	if data == nil {
		t.Error("expected data object")
	}
}

func TestWebHooksService_Test(t *testing.T) {
	_, client, cleanup := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/app/test_web_hook/v1" {
			t.Errorf("unexpected path %q", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"code": 0, "result": {"code": 0, "description": "delivered", "details": "## Report"}}`))
	})
	defer cleanup()

	code, desc, details, err := client.WebHooks.Test(context.Background(), &WebHook{
		ID: "w1", Title: "alerts", Method: "POST", URL: "https://example.com/hook",
	})
	if err != nil {
		t.Fatalf("test: %v", err)
	}
	if code != float64(0) {
		t.Errorf("expected result code 0, got %v", code)
	}
	if desc != "delivered" {
		t.Errorf("expected description delivered, got %q", desc)
	}
	if details == "" {
		t.Error("expected details report")
	}
}

func TestGroupsService_CreateSnapshot(t *testing.T) {
	_, client, cleanup := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/app/create_group_snapshot/v1" {
			t.Errorf("unexpected path %q", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"code": 0, "id": "snap1"}`))
	})
	defer cleanup()

	id, err := client.Groups.CreateSnapshot(context.Background(), "g1")
	if err != nil {
		t.Fatalf("create snapshot: %v", err)
	}
	if id != "snap1" {
		t.Errorf("expected snap1, got %q", id)
	}
}

func TestServersService_Get(t *testing.T) {
	_, client, cleanup := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/app/get_server/v1" {
			t.Errorf("unexpected path %q", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"code": 0, "server": {"id": "s1", "hostname": "worker01"}, "data": {"cpu": 42}, "online": true}`))
	})
	defer cleanup()

	server, data, online, err := client.Servers.Get(context.Background(), "s1")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if server.Hostname != "worker01" {
		t.Errorf("expected hostname worker01, got %s", server.Hostname)
	}
	if data["cpu"] != float64(42) {
		t.Errorf("expected data cpu=42, got %v", data["cpu"])
	}
	if !online {
		t.Error("expected online=true")
	}
}

func TestServersService_Active(t *testing.T) {
	_, client, cleanup := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/app/get_active_servers/v1" {
			t.Errorf("unexpected path %q", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"code": 0, "rows": [{"id": "s1", "hostname": "worker01"}], "list": {"length": 1}}`))
	})
	defer cleanup()

	servers, err := client.Servers.Active(context.Background())
	if err != nil {
		t.Fatalf("active: %v", err)
	}
	if len(servers) != 1 || servers[0].Hostname != "worker01" {
		t.Errorf("unexpected servers: %+v", servers)
	}
}

func TestServersService_Summaries(t *testing.T) {
	_, client, cleanup := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/app/get_server_summaries/v1" {
			t.Errorf("unexpected path %q", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"code": 0, "summaries": {"os_platform": {"Linux": 3}}}`))
	})
	defer cleanup()

	summaries, err := client.Servers.Summaries(context.Background())
	if err != nil {
		t.Fatalf("summaries: %v", err)
	}
	if summaries["os_platform"]["Linux"] != 3 {
		t.Errorf("expected Linux count 3, got %d", summaries["os_platform"]["Linux"])
	}
}

func TestJobsService_GetLog(t *testing.T) {
	_, client, cleanup := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/app/get_job_log/v1" {
			t.Errorf("unexpected path %q", r.URL.Path)
		}
		w.Header().Set("Content-Type", "text/plain")
		_, _ = w.Write([]byte("line1\nline2\n"))
	})
	defer cleanup()

	log, err := client.Jobs.GetLog(context.Background(), "j1")
	if err != nil {
		t.Fatalf("get log: %v", err)
	}
	if log != "line1\nline2\n" {
		t.Errorf("unexpected log %q", log)
	}
}

func TestJobsService_ToggleNotifyMe(t *testing.T) {
	_, client, cleanup := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/app/job_toggle_notify_me/v1" {
			t.Errorf("unexpected path %q", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"code": 0, "enabled": true}`))
	})
	defer cleanup()

	enabled, err := client.Jobs.ToggleNotifyMe(context.Background(), "j1")
	if err != nil {
		t.Fatalf("toggle notify: %v", err)
	}
	if !enabled {
		t.Error("expected enabled=true")
	}
}

func TestSearchService_SearchServers(t *testing.T) {
	_, client, cleanup := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/app/search_servers/v1" {
			t.Errorf("unexpected path %q", r.URL.Path)
		}
		if r.URL.Query().Get("query") != "keywords:build" {
			t.Errorf("expected query param, got %q", r.URL.Query().Get("query"))
		}
		_, _ = w.Write([]byte(`{"code": 0, "rows": [{"id": "s1", "hostname": "build-01"}], "list": {"length": 1}}`))
	})
	defer cleanup()

	servers, err := client.Search.SearchServers(context.Background(), SearchQueryOpts{Query: "keywords:build"})
	if err != nil {
		t.Fatalf("search servers: %v", err)
	}
	if len(servers) != 1 || servers[0].Hostname != "build-01" {
		t.Errorf("unexpected servers: %+v", servers)
	}
}

func TestSearchService_SearchStatHistory(t *testing.T) {
	_, client, cleanup := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/app/search_stat_history/v1" {
			t.Errorf("unexpected path %q", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"code": 0, "items": [{"epoch": 1, "date": "2026-08-01", "data": {"jobs": 5}}], "list": {"length": 1}}`))
	})
	defer cleanup()

	items, err := client.Search.SearchStatHistory(context.Background(), 1, 7, "daily.jobs", "", false)
	if err != nil {
		t.Fatalf("search stat history: %v", err)
	}
	if len(items) != 1 || items[0]["data"].(map[string]any)["jobs"] != float64(5) {
		t.Errorf("unexpected items: %+v", items)
	}
}

func TestSearchService_BulkExport(t *testing.T) {
	_, client, cleanup := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/app/bulk_search_export/v1" {
			t.Errorf("unexpected path %q", r.URL.Path)
		}
		if r.URL.Query().Get("index") != "jobs" || r.URL.Query().Get("format") != "csv" {
			t.Errorf("unexpected params: %s", r.URL.RawQuery)
		}
		_, _ = w.Write([]byte("\xef\xbb\xbfid,title\nj1,backup\n"))
	})
	defer cleanup()

	data, err := client.Search.BulkExport(context.Background(), "jobs", []string{"id", "title"}, "csv", "", false)
	if err != nil {
		t.Fatalf("bulk export: %v", err)
	}
	if len(data) == 0 {
		t.Error("expected non-empty export bytes")
	}
}

func TestSearchService_Marketplace(t *testing.T) {
	_, client, cleanup := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/app/marketplace/v1" {
			t.Errorf("unexpected path %q", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"code": 0, "rows": [{"id": "p1", "title": "echo", "type": "plugin"}], "list": {"length": 1}}`))
	})
	defer cleanup()

	items, err := client.Search.Marketplace(context.Background(), MarketplaceOpts{Query: "echo", PluginType: "event"})
	if err != nil {
		t.Fatalf("marketplace: %v", err)
	}
	if len(items) != 1 || items[0]["title"] != "echo" {
		t.Errorf("unexpected marketplace items: %+v", items)
	}
}
