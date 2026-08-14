package sdk

import (
	"context"
	"net/http"
	"testing"
)

func TestResolver_SeedAndResolve(t *testing.T) {
	_, client, cleanup := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
		t.Error("no API calls expected when seeded")
	})
	defer cleanup()

	r := NewResolver(client)
	r.Seed("group", map[string]string{"vps": "g123", "ecs": "g456"})

	id, err := r.Resolve(context.Background(), "group", "vps")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != "g123" {
		t.Errorf("expected g123, got %s", id)
	}
}

func TestResolver_APIRefresh(t *testing.T) {
	_, client, cleanup := setupTestServer(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"code": 0, "rows": [{"id": "g1", "title": "vps"}, {"id": "g2", "title": "ecs"}], "list": {"length": 2}}`))
	})
	defer cleanup()

	r := NewResolver(client)

	id, err := r.Resolve(context.Background(), "group", "ecs")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != "g2" {
		t.Errorf("expected g2, got %s", id)
	}
}

func TestResolver_AutoRegister(t *testing.T) {
	_, client, cleanup := setupTestServer(func(w http.ResponseWriter, req *http.Request) {
		t.Error("no API calls expected after AutoRegister")
	})
	defer cleanup()

	r := NewResolver(client)
	r.AutoRegister("tag", "backup", "t123")

	id, err := r.Resolve(context.Background(), "tag", "backup")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != "t123" {
		t.Errorf("expected t123, got %s", id)
	}
}

func TestResolver_UnknownKey(t *testing.T) {
	_, client, cleanup := setupTestServer(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"code": 0, "rows": [], "list": {"length": 0}}`))
	})
	defer cleanup()

	r := NewResolver(client)
	_, err := r.Resolve(context.Background(), "group", "nonexistent")
	if err == nil {
		t.Fatal("expected error for unknown key")
	}
}

func TestResolver_CacheInvalidation(t *testing.T) {
	calls := 0
	_, client, cleanup := setupTestServer(func(w http.ResponseWriter, req *http.Request) {
		calls++
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"code": 0, "rows": [{"id": "g1", "title": "vps"}], "list": {"length": 1}}`))
	})
	defer cleanup()

	r := NewResolver(client)
	ctx := context.Background()
	if _, err := r.Resolve(ctx, "group", "vps"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := r.Resolve(ctx, "group", "vps"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if calls != 1 {
		t.Errorf("expected 1 API call, got %d", calls)
	}

	r.Clear("group")
	if _, err := r.Resolve(ctx, "group", "vps"); err != nil {
		t.Fatalf("unexpected error after clear: %v", err)
	}
	if calls != 2 {
		t.Errorf("expected 2 API calls after clear, got %d", calls)
	}
}

func TestResolver_EmptyKey(t *testing.T) {
	_, client, cleanup := setupTestServer(func(w http.ResponseWriter, req *http.Request) {
		t.Error("no API calls expected for empty key")
	})
	defer cleanup()

	r := NewResolver(client)
	_, err := r.Resolve(context.Background(), "group", "")
	if err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestResolver_UnknownType(t *testing.T) {
	_, client, cleanup := setupTestServer(func(w http.ResponseWriter, req *http.Request) {
		t.Error("no API calls expected for unknown type")
	})
	defer cleanup()

	r := NewResolver(client)
	_, err := r.Resolve(context.Background(), "bogus", "x")
	if err == nil {
		t.Fatal("expected error for unknown resource type")
	}
}

func TestResolver_NewTypes(t *testing.T) {
	tests := []struct {
		name       string
		pathPrefix string
		rows       string
		resolveKey string
		resolveID  string
	}{
		{
			name:       "channel",
			pathPrefix: "/api/app/get_channels/v1",
			rows:       `[{"id": "ch1", "title": "oncall"}, {"id": "ch2", "title": "alerts"}]`,
			resolveKey: "alerts",
			resolveID:  "ch2",
		},
		{
			name:       "secret",
			pathPrefix: "/api/app/get_secrets/v1",
			rows:       `[{"id": "se1", "title": "db-credentials"}, {"id": "se2", "title": "api-token"}]`,
			resolveKey: "db-credentials",
			resolveID:  "se1",
		},
		{
			name:       "bucket",
			pathPrefix: "/api/app/get_buckets/v1",
			rows:       `[{"id": "b1", "title": "artifacts"}, {"id": "b2", "title": "logs"}]`,
			resolveKey: "logs",
			resolveID:  "b2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, client, cleanup := setupTestServer(func(w http.ResponseWriter, req *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(`{"code": 0, "rows": ` + tt.rows + `, "list": {"length": 2}}`))
			})
			defer cleanup()

			r := NewResolver(client)
			id, err := r.Resolve(context.Background(), tt.name, tt.resolveKey)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if id != tt.resolveID {
				t.Errorf("expected %s, got %s", tt.resolveID, id)
			}
		})
	}
}
