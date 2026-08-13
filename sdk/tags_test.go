package sdk

import (
	"context"
	"net/http"
	"testing"
)

func TestTagsService_List(t *testing.T) {
	_, client, cleanup := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"code": 0, "rows": [{"id": "t1", "title": "backup"}, {"id": "t2", "title": "prod"}], "list": {"length": 2}}`))
	})
	defer cleanup()

	tags, err := client.Tags.List(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tags) != 2 {
		t.Errorf("expected 2 tags, got %d", len(tags))
	}
	if tags[0].Title != "backup" {
		t.Errorf("expected first tag title 'backup', got %q", tags[0].Title)
	}
}

func TestTagsService_Get(t *testing.T) {
	_, client, cleanup := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("id") != "t1" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"code": 0, "tag": {"id": "t1", "title": "backup"}}`))
	})
	defer cleanup()

	tag, err := client.Tags.Get(context.Background(), "t1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tag.Title != "backup" {
		t.Errorf("expected title 'backup', got %q", tag.Title)
	}
}

func TestTagsService_Create(t *testing.T) {
	_, client, cleanup := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"code": 0, "tag": {"id": "t123", "title": "new-tag"}}`))
	})
	defer cleanup()

	tag, err := client.Tags.Create(context.Background(), &Tag{Title: "new-tag"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tag.ID != "t123" {
		t.Errorf("expected ID 't123', got %q", tag.ID)
	}
}

func TestTagsService_Update(t *testing.T) {
	calls := 0
	_, client, cleanup := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
		calls++
		if r.URL.Path == "/api/app/update_tag/v1" {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"code": 0}`))
		} else {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"code": 0, "tag": {"id": "t123", "title": "updated"}}`))
		}
	})
	defer cleanup()

	tag, err := client.Tags.Update(context.Background(), &Tag{ID: "t123", Title: "updated"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tag.Title != "updated" {
		t.Errorf("expected title 'updated', got %q", tag.Title)
	}
	if calls != 2 { // update then get
		t.Errorf("expected 2 calls (update + get), got %d", calls)
	}
}

func TestTagsService_Delete(t *testing.T) {
	_, client, cleanup := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"code": 0}`))
	})
	defer cleanup()

	if err := client.Tags.Delete(context.Background(), "t123"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
