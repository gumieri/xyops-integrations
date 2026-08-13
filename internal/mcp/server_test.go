package mcp

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gumieri/xyops-integrations/sdk"
)

func TestNewServer(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"code": 0}`))
	}))
	defer server.Close()

	c, err := sdk.New(server.URL, "test-key")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	mcpServer := NewServer(c)
	if mcpServer == nil {
		t.Fatal("expected non-nil server")
	}
}

func TestAdaptResult(t *testing.T) {
	result, output, err := adaptResult("value", nil)
	if result != nil || output != "value" || err != nil {
		t.Fatalf("adaptResult with success: result=%v output=%v err=%v", result, output, err)
	}

	result, output, err = adaptResult("", http.ErrNotSupported)
	if result != nil || output != "" || err != http.ErrNotSupported {
		t.Fatalf("adaptResult with error: result=%v output=%v err=%v", result, output, err)
	}
}
