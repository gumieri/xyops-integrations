package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupTestServer(handler http.HandlerFunc) (*httptest.Server, *Client, func()) {
	server := httptest.NewServer(handler)
	client, err := NewWithHTTPClient(server.URL, "test-api-key", server.Client())
	if err != nil {
		server.Close()
		panic(err)
	}
	return server, client, server.Close
}

func TestNew_InvalidURL(t *testing.T) {
	_, err := New("://invalid", "key")
	if err == nil {
		t.Fatal("expected error for invalid URL")
	}
}

func TestNewWithHTTPClient(t *testing.T) {
	_, err := NewWithHTTPClient("http://example.com", "key", http.DefaultClient)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAuthHeaderInjection(t *testing.T) {
	_, client, cleanup := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-API-Key") != "test-api-key" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"code": 0}`))
	})
	defer cleanup()

	_, err := client.DoGET(context.Background(), "test/v1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestContentTypeOnPOST(t *testing.T) {
	_, client, cleanup := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		if r.Header.Get("Content-Type") != "application/json" {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"code": 0}`))
	})
	defer cleanup()

	_, err := client.DoPOST(context.Background(), "test/v1", map[string]string{"foo": "bar"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNoContentTypeOnGET(t *testing.T) {
	_, client, cleanup := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "" {
			t.Error("GET request should not have Content-Type")
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"code": 0}`))
	})
	defer cleanup()

	_, err := client.DoGET(context.Background(), "test/v1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDoGETWithQuery(t *testing.T) {
	_, client, cleanup := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		if r.URL.Query().Get("id") != "abc123" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"code": 0}`))
	})
	defer cleanup()

	_, err := client.DoGETWithQuery(context.Background(), "test/v1", map[string]string{"id": "abc123"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestErrorResponseParsing(t *testing.T) {
	_, client, cleanup := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"code": 1, "message": "not found"}`))
	})
	defer cleanup()

	_, err := client.DoGET(context.Background(), "test/v1")
	if err == nil {
		t.Fatal("expected error")
	}

	var apiErr *APIError
	if !asAPIError(err, &apiErr) {
		t.Fatalf("expected APIError, got %T", err)
	}
	if apiErr.Code != 1 {
		t.Errorf("expected code 1, got %d", apiErr.Code)
	}
	if apiErr.Message != "not found" {
		t.Errorf("expected message 'not found', got %q", apiErr.Message)
	}
}

func TestHTTPErrorResponse(t *testing.T) {
	_, client, cleanup := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"code": 5, "message": "internal error"}`))
	})
	defer cleanup()

	_, err := client.DoGET(context.Background(), "test/v1")
	if err == nil {
		t.Fatal("expected error")
	}

	var apiErr *APIError
	if !asAPIError(err, &apiErr) {
		t.Fatalf("expected APIError, got %T", err)
	}
	if apiErr.HTTPStatus != 500 {
		t.Errorf("expected HTTP 500, got %d", apiErr.HTTPStatus)
	}
	if apiErr.Code != 5 {
		t.Errorf("expected code 5, got %d", apiErr.Code)
	}
}

func TestHealth_Healthy(t *testing.T) {
	_, client, cleanup := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok", "primary": true}`))
	})
	defer cleanup()

	healthy, err := client.Health(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !healthy {
		t.Error("expected healthy")
	}
}

func TestHealth_NotPrimary(t *testing.T) {
	_, client, cleanup := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte(`{"status": "not primary"}`))
	})
	defer cleanup()

	healthy, err := client.Health(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if healthy {
		t.Error("expected not healthy")
	}
}

func TestBodySerialization(t *testing.T) {
	input := Tag{Title: "test-tag"}
	_, client, cleanup := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
		var received Tag
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if received.Title != "test-tag" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"code": 0, "tag": {"id": "t123", "title": "test-tag"}}`))
	})
	defer cleanup()

	body, err := client.DoPOST(context.Background(), "test/v1", input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var resp TagResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if resp.Tag == nil {
		t.Fatal("expected tag in response")
	}
	if resp.Tag.ID != "t123" {
		t.Errorf("expected id t123, got %s", resp.Tag.ID)
	}
}

func TestDoGETWithEmptyQuery(t *testing.T) {
	_, client, cleanup := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.RawQuery != "" {
			t.Errorf("expected no query params, got %s", r.URL.RawQuery)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"code": 0}`))
	})
	defer cleanup()

	_, err := client.DoGETWithQuery(context.Background(), "test/v1", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestURLResolution(t *testing.T) {
	// Test that relative paths resolve to /api/app/NAME
	_, client, cleanup := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/app/get_tag/v1" {
			t.Errorf("expected /api/app/get_tag/v1, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"code": 0}`))
	})
	defer cleanup()

	_, err := client.DoGET(context.Background(), "get_tag/v1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAbsoluteURLResolution(t *testing.T) {
	// Test that absolute paths don't get /api/app/ prepended
	_, client, cleanup := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/health" {
			t.Errorf("expected /health, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	})
	defer cleanup()

	_, err := client.DoGET(context.Background(), "/health")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAPIError_IsAuth(t *testing.T) {
	tests := []struct {
		err    *APIError
		isAuth bool
	}{
		{NewAPIError(0, "session", 200), true},
		{NewAPIError(1, "session", 200), true},
		{NewAPIError(0, "not found", 404), false},
		{NewAPIError(1, "invalid", 401), true},
		{NewAPIError(1, "invalid", 403), true},
		{NewAPIError(1, "invalid", 200), false},
	}
	for _, tt := range tests {
		if got := tt.err.IsAuth(); got != tt.isAuth {
			t.Errorf("IsAuth() = %v, want %v for %v", got, tt.isAuth, tt.err)
		}
	}
}

// asAPIError extracts an *APIError from an error chain.
func asAPIError(err error, target **APIError) bool {
	if e, ok := err.(*APIError); ok {
		*target = e
		return true
	}
	return false
}

func TestDefaultLimits(t *testing.T) {
	limits := DefaultLimits()
	if len(limits) != 3 {
		t.Fatalf("expected 3 default limits, got %d", len(limits))
	}
	expected := map[string]int{"job": 1, "queue": 0, "retry": 0}
	for _, l := range limits {
		if expected[l.Type] != l.Amount {
			t.Errorf("limit %s: expected amount %d, got %d", l.Type, expected[l.Type], l.Amount)
		}
	}
}

func TestServicesInitialized(t *testing.T) {
	c, err := New("http://localhost", "test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Tags == nil {
		t.Error("Tags service not initialized")
	}
	if c.Groups == nil {
		t.Error("Groups service not initialized")
	}
	if c.Categories == nil {
		t.Error("Categories service not initialized")
	}
	if c.Plugins == nil {
		t.Error("Plugins service not initialized")
	}
	if c.Events == nil {
		t.Error("Events service not initialized")
	}
	if c.Monitors == nil {
		t.Error("Monitors service not initialized")
	}
	if c.Alerts == nil {
		t.Error("Alerts service not initialized")
	}
	if c.Roles == nil {
		t.Error("Roles service not initialized")
	}
	if c.APIKeys == nil {
		t.Error("APIKeys service not initialized")
	}
	if c.WebHooks == nil {
		t.Error("WebHooks service not initialized")
	}
	if c.Config == nil {
		t.Error("Config service not initialized")
	}
	if c.Servers == nil {
		t.Error("Servers service not initialized")
	}
	if c.Satellite == nil {
		t.Error("Satellite service not initialized")
	}
	if c.Jobs == nil {
		t.Error("Jobs service not initialized")
	}
	if c.GetMultiple == nil {
		t.Error("GetMultiple service not initialized")
	}
}
