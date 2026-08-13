package sdk

import (
	"context"
	"fmt"
	"net/http"
	"testing"
)

func TestGroupsService_CRUD(t *testing.T) {
	calls := 0
	_, client, cleanup := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
		calls++
		switch r.URL.Path {
		case "/api/app/get_groups/v1":
			_, _ = w.Write([]byte(`{"code": 0, "rows": [{"id": "g1", "title": "vps"}], "list": {"length": 1}}`))
		case "/api/app/get_group/v1":
			_, _ = w.Write([]byte(`{"code": 0, "group": {"id": "g1", "title": "vps"}}`))
		case "/api/app/create_group/v1":
			_, _ = w.Write([]byte(`{"code": 0, "group": {"id": "g1", "title": "vps"}}`))
		case "/api/app/update_group/v1":
			_, _ = w.Write([]byte(`{"code": 0}`))
		case "/api/app/delete_group/v1":
			_, _ = w.Write([]byte(`{"code": 0}`))
		default:
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(fmt.Sprintf(`{"code": 1, "message": "unexpected path %q"}`, r.URL.Path)))
		}
	})
	defer cleanup()

	ctx := context.Background()

	if _, err := client.Groups.List(ctx); err != nil {
		t.Fatalf("list: %v", err)
	}
	if _, err := client.Groups.Get(ctx, "g1"); err != nil {
		t.Fatalf("get: %v", err)
	}
	if _, err := client.Groups.Create(ctx, &Group{Title: "vps"}); err != nil {
		t.Fatalf("create: %v", err)
	}
	if _, err := client.Groups.Update(ctx, &Group{ID: "g1", Title: "vps"}); err != nil {
		t.Fatalf("update: %v", err)
	}
	if err := client.Groups.Delete(ctx, "g1"); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if calls != 6 {
		t.Errorf("expected 6 calls (list+get+create+update+get+delete), got %d", calls)
	}
}

func TestCategoriesService_CRUD(t *testing.T) {
	_, client, cleanup := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"code": 0, "rows": [{"id": "c1", "title": "maintenance", "enabled": true}], "list": {"length": 1}}`))
	})
	defer cleanup()

	cats, err := client.Categories.List(context.Background())
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(cats) != 1 || cats[0].Title != "maintenance" {
		t.Errorf("unexpected categories: %+v", cats)
	}
}

func TestPluginsService_Create(t *testing.T) {
	_, client, cleanup := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"code": 0, "plugin": {"id": "p1", "title": "echo", "type": "event"}}`))
	})
	defer cleanup()

	plugin, err := client.Plugins.Create(context.Background(), &Plugin{Title: "echo", Type: "event"})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if plugin.ID != "p1" {
		t.Errorf("expected ID p1, got %s", plugin.ID)
	}
}

func TestEventsService_Create_WorkflowInjection(t *testing.T) {
	_, client, cleanup := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"code": 0, "event": {"id": "e1", "title": "wf", "plugin": "_workflow", "type": "workflow"}}`))
	})
	defer cleanup()

	event, err := client.Events.Create(context.Background(), &Event{
		Title:    "wf",
		Workflow: &WorkflowGraph{},
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if event.Plugin != "_workflow" {
		t.Errorf("expected injected plugin _workflow, got %s", event.Plugin)
	}
}

func TestMonitorsService_Get(t *testing.T) {
	_, client, cleanup := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"code": 0, "monitor": {"id": "m1", "title": "cpu", "source": "node_exporter"}}`))
	})
	defer cleanup()

	m, err := client.Monitors.Get(context.Background(), "m1")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if m.Title != "cpu" {
		t.Errorf("expected title cpu, got %s", m.Title)
	}
}

func TestRolesService_List(t *testing.T) {
	_, client, cleanup := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"code": 0, "rows": [{"id": "r1", "title": "admin", "enabled": true, "privileges": {"config": 1}}], "list": {"length": 1}}`))
	})
	defer cleanup()

	roles, err := client.Roles.List(context.Background())
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(roles) != 1 || roles[0].Privileges["config"] != 1 {
		t.Errorf("unexpected roles: %+v", roles)
	}
}

func TestAPIKeysService_Get(t *testing.T) {
	_, client, cleanup := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"code": 0, "api_key": {"id": "k1", "title": "deploy", "active": 1}}`))
	})
	defer cleanup()

	key, err := client.APIKeys.Get(context.Background(), "k1")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if key.Active != 1 {
		t.Errorf("expected active=1, got %d", key.Active)
	}
}

func TestWebHooksService_CRUD(t *testing.T) {
	calls := 0
	_, client, cleanup := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
		calls++
		_, _ = w.Write([]byte(`{"code": 0, "web_hook": {"id": "w1", "title": "alerts", "url": "https://example.com/hook", "enabled": true}}`))
	})
	defer cleanup()

	ctx := context.Background()
	if _, err := client.WebHooks.List(ctx); err != nil {
		t.Fatalf("list: %v", err)
	}
	if _, err := client.WebHooks.Create(ctx, &WebHook{Title: "alerts", URL: "https://example.com/hook"}); err != nil {
		t.Fatalf("create: %v", err)
	}
	if _, err := client.WebHooks.Get(ctx, "w1"); err != nil {
		t.Fatalf("get: %v", err)
	}
	if calls != 3 {
		t.Errorf("expected 3 calls, got %d", calls)
	}
}

func TestConfigService_Get(t *testing.T) {
	_, client, cleanup := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"code": 0, "config": {"check_updates": true}}`))
	})
	defer cleanup()

	cfg, err := client.Config.Get(context.Background())
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if cfg["check_updates"] != true {
		t.Errorf("expected check_updates=true, got %v", cfg["check_updates"])
	}
}

func TestServersService_List(t *testing.T) {
	_, client, cleanup := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"code": 0, "servers": [{"hostname": "s1", "os": "linux", "arch": "arm64", "status": "connected"}]}`))
	})
	defer cleanup()

	servers, err := client.Servers.List(context.Background())
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(servers) != 1 || servers[0].Hostname != "s1" {
		t.Errorf("unexpected servers: %+v", servers)
	}
}

func TestGetMultipleService_Get(t *testing.T) {
	_, client, cleanup := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"code": 0, "tags": [{"id": "t1", "title": "backup"}]}`))
	})
	defer cleanup()

	resp, err := client.GetMultiple.Get(context.Background(), GetMultipleOpts{Lists: []string{"tags"}})
	if err != nil {
		t.Fatalf("get multiple: %v", err)
	}
	if len(resp.Tags) != 1 || resp.Tags[0].Title != "backup" {
		t.Errorf("unexpected get_multiple response: %+v", resp)
	}
}

func TestSatelliteService_Token(t *testing.T) {
	_, client, cleanup := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"code": 0, "token": "tok123"}`))
	})
	defer cleanup()

	token, err := client.Satellite.Token(context.Background())
	if err != nil {
		t.Fatalf("token: %v", err)
	}
	if token != "tok123" {
		t.Errorf("expected tok123, got %s", token)
	}
}
