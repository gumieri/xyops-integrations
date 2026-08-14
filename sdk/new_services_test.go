package sdk

import (
	"context"
	"net/http"
	"testing"
)

func TestChannelsService_CRUD(t *testing.T) {
	calls := 0
	_, client, cleanup := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
		calls++
		switch r.URL.Path {
		case "/api/app/get_channels/v1":
			_, _ = w.Write([]byte(`{"code": 0, "rows": [{"id": "c1", "title": "sev1", "enabled": true, "users": ["admin"]}], "list": {"length": 1}}`))
		case "/api/app/get_channel/v1":
			_, _ = w.Write([]byte(`{"code": 0, "channel": {"id": "c1", "title": "sev1", "enabled": true}}`))
		case "/api/app/create_channel/v1":
			_, _ = w.Write([]byte(`{"code": 0, "channel": {"id": "c1", "title": "sev1", "enabled": true}}`))
		case "/api/app/update_channel/v1":
			_, _ = w.Write([]byte(`{"code": 0}`))
		case "/api/app/delete_channel/v1":
			_, _ = w.Write([]byte(`{"code": 0}`))
		default:
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"code": 1, "message": "unexpected path"}`))
		}
	})
	defer cleanup()

	ctx := context.Background()
	if _, err := client.Channels.List(ctx); err != nil {
		t.Fatalf("list: %v", err)
	}
	if _, err := client.Channels.Get(ctx, "c1"); err != nil {
		t.Fatalf("get: %v", err)
	}
	if _, err := client.Channels.Create(ctx, &Channel{Title: "sev1"}); err != nil {
		t.Fatalf("create: %v", err)
	}
	if _, err := client.Channels.Update(ctx, &Channel{ID: "c1", Title: "sev1"}); err != nil {
		t.Fatalf("update: %v", err)
	}
	if err := client.Channels.Delete(ctx, "c1"); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if calls != 6 {
		t.Errorf("expected 6 calls (list+get+create+update+get+delete), got %d", calls)
	}
}

func TestSecretsService_CRUD(t *testing.T) {
	calls := 0
	_, client, cleanup := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
		calls++
		switch r.URL.Path {
		case "/api/app/get_secrets/v1":
			_, _ = w.Write([]byte(`{"code": 0, "rows": [{"id": "s1", "title": "db", "enabled": true, "names": ["DB_HOST"]}], "list": {"length": 1}}`))
		case "/api/app/get_secret/v1":
			_, _ = w.Write([]byte(`{"code": 0, "secret": {"id": "s1", "title": "db", "enabled": true}}`))
		case "/api/app/decrypt_secret/v1":
			_, _ = w.Write([]byte(`{"code": 0, "fields": [{"name": "DB_HOST", "value": "db.dev.internal"}]}`))
		case "/api/app/create_secret/v1":
			_, _ = w.Write([]byte(`{"code": 0, "secret": {"id": "s1", "title": "db", "enabled": true}}`))
		case "/api/app/update_secret/v1":
			_, _ = w.Write([]byte(`{"code": 0}`))
		case "/api/app/delete_secret/v1":
			_, _ = w.Write([]byte(`{"code": 0}`))
		default:
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"code": 1, "message": "unexpected path"}`))
		}
	})
	defer cleanup()

	ctx := context.Background()
	if _, err := client.Secrets.List(ctx); err != nil {
		t.Fatalf("list: %v", err)
	}
	if _, err := client.Secrets.Get(ctx, "s1"); err != nil {
		t.Fatalf("get: %v", err)
	}
	fields, err := client.Secrets.Decrypt(ctx, "s1")
	if err != nil {
		t.Fatalf("decrypt: %v", err)
	}
	if len(fields) != 1 || fields[0].Name != "DB_HOST" || fields[0].Value != "db.dev.internal" {
		t.Errorf("unexpected fields: %+v", fields)
	}
	if _, err := client.Secrets.Create(ctx, &Secret{Title: "db"}); err != nil {
		t.Fatalf("create: %v", err)
	}
	if _, err := client.Secrets.Update(ctx, &Secret{ID: "s1", Title: "db"}); err != nil {
		t.Fatalf("update: %v", err)
	}
	if err := client.Secrets.Delete(ctx, "s1"); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if calls != 7 {
		t.Errorf("expected 7 calls (list+get+decrypt+create+update+get+delete), got %d", calls)
	}
}

func TestBucketsService_Get_Composite(t *testing.T) {
	_, client, cleanup := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/app/get_bucket/v1" {
			t.Errorf("unexpected path %q", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"code": 0, "bucket": {"id": "b1", "title": "the-void"}, "data": {"foo": "bar"}, "files": [{"id": "f1", "filename": "test.png", "size": 92615, "path": "files/bucket/b1/f1/test.png"}]}`))
	})
	defer cleanup()

	bucket, data, files, err := client.Buckets.Get(context.Background(), "b1")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if bucket.Title != "the-void" {
		t.Errorf("expected title the-void, got %s", bucket.Title)
	}
	if data["foo"] != "bar" {
		t.Errorf("expected data foo=bar, got %v", data)
	}
	if len(files) != 1 || files[0].Filename != "test.png" {
		t.Errorf("unexpected files: %+v", files)
	}
}

func TestBucketsService_WriteData_Fetch(t *testing.T) {
	_, client, cleanup := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/app/write_bucket_data/v1" {
			t.Errorf("unexpected path %q", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"code": 0, "data": {"foo": "bar", "other": 12345}}`))
	})
	defer cleanup()

	data, err := client.Buckets.WriteData(context.Background(), "b1", map[string]any{"foo": "bar"}, true)
	if err != nil {
		t.Fatalf("write data: %v", err)
	}
	if data["foo"] != "bar" || data["other"] != float64(12345) {
		t.Errorf("unexpected merged data: %+v", data)
	}
}

func TestTicketsService_CRUD(t *testing.T) {
	calls := 0
	_, client, cleanup := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
		calls++
		switch r.URL.Path {
		case "/api/app/get_ticket/v1":
			_, _ = w.Write([]byte(`{"code": 0, "ticket": {"id": "t1", "num": 24, "subject": "backup failed", "status": "open"}}`))
		case "/api/app/search_tickets/v1":
			_, _ = w.Write([]byte(`{"code": 0, "rows": [{"id": "t1", "num": 24, "subject": "backup failed", "status": "open"}], "list": {"length": 1}}`))
		case "/api/app/create_ticket/v1":
			_, _ = w.Write([]byte(`{"code": 0, "ticket": {"id": "t1", "num": 24, "subject": "backup failed", "status": "open"}}`))
		case "/api/app/update_ticket/v1":
			_, _ = w.Write([]byte(`{"code": 0}`))
		case "/api/app/delete_ticket/v1":
			_, _ = w.Write([]byte(`{"code": 0}`))
		default:
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"code": 1, "message": "unexpected path"}`))
		}
	})
	defer cleanup()

	ctx := context.Background()
	if _, err := client.Tickets.Search(ctx, SearchTicketsOpts{Query: "status:open"}); err != nil {
		t.Fatalf("search: %v", err)
	}
	if _, err := client.Tickets.Get(ctx, "t1"); err != nil {
		t.Fatalf("get: %v", err)
	}
	if _, err := client.Tickets.Create(ctx, &Ticket{Subject: "backup failed"}); err != nil {
		t.Fatalf("create: %v", err)
	}
	if _, err := client.Tickets.Update(ctx, &Ticket{ID: "t1", Status: "closed"}); err != nil {
		t.Fatalf("update: %v", err)
	}
	if err := client.Tickets.Delete(ctx, "t1"); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if calls != 6 {
		t.Errorf("expected 6 calls (search+get+create+update+get+delete), got %d", calls)
	}
}

func TestTicketsService_AddChange(t *testing.T) {
	_, client, cleanup := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/app/add_ticket_change/v1" {
			t.Errorf("unexpected path %q", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"code": 0, "ticket": {"id": "t1", "changes": [{"type": "comment", "body": "Investigating"}]}}`))
	})
	defer cleanup()

	ticket, err := client.Tickets.AddChange(context.Background(), "t1", TicketChange{Type: "comment", Body: "Investigating"})
	if err != nil {
		t.Fatalf("add change: %v", err)
	}
	if len(ticket.Changes) != 1 || ticket.Changes[0].Type != "comment" {
		t.Errorf("unexpected changes: %+v", ticket.Changes)
	}
}
