package mcp

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gumieri/xyops-integrations/sdk"
	"github.com/modelcontextprotocol/go-sdk/mcp"
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

func TestAdaptListResult(t *testing.T) {
	result, output, err := adaptListResult([]int{1, 2, 3}, nil)
	if result != nil || err != nil {
		t.Fatalf("adaptListResult with success: result=%v err=%v", result, err)
	}
	if !reflect.DeepEqual(output.Items, []int{1, 2, 3}) {
		t.Fatalf("adaptListResult wrapped items = %v, want [1 2 3]", output.Items)
	}

	_, _, err = adaptListResult[int](nil, http.ErrNotSupported)
	if err != http.ErrNotSupported {
		t.Fatalf("adaptListResult with error: err=%v", err)
	}
}

func TestAdaptStatusResult(t *testing.T) {
	result, output, err := adaptStatusResult("deleted", nil)
	if result != nil || err != nil {
		t.Fatalf("adaptStatusResult with success: result=%v err=%v", result, err)
	}
	if output.Status != "deleted" {
		t.Fatalf("adaptStatusResult status = %q, want %q", output.Status, "deleted")
	}

	_, _, err = adaptStatusResult("deleted", http.ErrNotSupported)
	if err != http.ErrNotSupported {
		t.Fatalf("adaptStatusResult with error: err=%v", err)
	}
}

func TestAdaptHealthResult(t *testing.T) {
	result, output, err := adaptHealthResult(true, nil)
	if result != nil || err != nil {
		t.Fatalf("adaptHealthResult with success: result=%v err=%v", result, err)
	}
	if !output.Healthy {
		t.Fatalf("adaptHealthResult healthy = %v, want true", output.Healthy)
	}

	_, _, err = adaptHealthResult(false, http.ErrNotSupported)
	if err != http.ErrNotSupported {
		t.Fatalf("adaptHealthResult with error: err=%v", err)
	}
}

func TestIsObjectSchema(t *testing.T) {
	tests := []struct {
		name string
		in   any
		want bool
	}{
		{"nil", nil, false},
		{"object type", map[string]any{"type": "object"}, true},
		{"string type", map[string]any{"type": "string"}, false},
		{"array type", map[string]any{"type": "array"}, false},
		{"no type", map[string]any{"properties": map[string]any{}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isObjectSchema(tt.in); got != tt.want {
				t.Errorf("isObjectSchema(%v) = %v, want %v", tt.in, got, tt.want)
			}
		})
	}
}

func TestNormalizeSchema(t *testing.T) {
	tests := []struct {
		name string
		in   any
		want any
	}{
		{"nil", nil, nil},
		{
			"boolean true property becomes open schema",
			map[string]any{
				"type":       "object",
				"properties": map[string]any{"params": true},
			},
			map[string]any{
				"type":       "object",
				"properties": map[string]any{"params": map[string]any{}},
			},
		},
		{
			"boolean false property left unchanged",
			map[string]any{
				"type":       "object",
				"properties": map[string]any{"never": false},
			},
			map[string]any{
				"type":       "object",
				"properties": map[string]any{"never": false},
			},
		},
		{
			"additionalProperties true left unchanged",
			map[string]any{
				"type":                 "object",
				"additionalProperties": true,
			},
			map[string]any{
				"type":                 "object",
				"additionalProperties": true,
			},
		},
		{
			"nested object property recurses",
			map[string]any{
				"type": "object",
				"properties": map[string]any{
					"nested": map[string]any{
						"type":       "object",
						"properties": map[string]any{"deep": true},
					},
				},
			},
			map[string]any{
				"type": "object",
				"properties": map[string]any{
					"nested": map[string]any{
						"type":       "object",
						"properties": map[string]any{"deep": map[string]any{}},
					},
				},
			},
		},
		{
			"boolean inside array items recurses",
			map[string]any{
				"type":  "object",
				"items": map[string]any{"type": "object", "properties": map[string]any{"data": true}},
			},
			map[string]any{
				"type":  "object",
				"items": map[string]any{"type": "object", "properties": map[string]any{"data": map[string]any{}}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := normalizeSchema(tt.in)
			if !equalAny(got, tt.want) {
				t.Errorf("normalizeSchema(%v) = %v, want %v", tt.in, got, tt.want)
			}
		})
	}
}

func TestNormalizeOutputSchema(t *testing.T) {
	tests := []struct {
		name string
		in   any
		want any
	}{
		{"nil dropped", nil, nil},
		{"string schema dropped", map[string]any{"type": "string"}, nil},
		{"array schema dropped", map[string]any{"type": "array", "items": map[string]any{}}, nil},
		{
			"object schema preserved",
			map[string]any{"type": "object", "properties": map[string]any{"id": map[string]any{"type": "string"}}},
			map[string]any{"type": "object", "properties": map[string]any{"id": map[string]any{"type": "string"}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := normalizeOutputSchema(tt.in)
			if !equalAny(got, tt.want) {
				t.Errorf("normalizeOutputSchema(%v) = %v, want %v", tt.in, got, tt.want)
			}
		})
	}
}

func TestSanitizeToolSchemasMiddleware(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		tools      []*mcp.Tool
		wantOutput map[string]bool // tool name -> outputSchema nil?
		wantParams map[string]bool // tool name -> params expanded to open schema?
	}{
		{
			name:   "non tools/list method untouched",
			method: "tools/call",
			tools: []*mcp.Tool{
				{Name: "abort_job", OutputSchema: map[string]any{"type": "string"}},
			},
			wantOutput: map[string]bool{"abort_job": false},
			wantParams: map[string]bool{},
		},
		{
			name:   "primitive outputSchema dropped",
			method: "tools/list",
			tools: []*mcp.Tool{
				{Name: "abort_job", OutputSchema: map[string]any{"type": "string"}},
			},
			wantOutput: map[string]bool{"abort_job": true},
			wantParams: map[string]bool{},
		},
		{
			name:   "object outputSchema preserved",
			method: "tools/list",
			tools: []*mcp.Tool{
				{Name: "get_event", OutputSchema: map[string]any{"type": "object", "properties": map[string]any{}}},
			},
			wantOutput: map[string]bool{"get_event": false},
			wantParams: map[string]bool{},
		},
		{
			name:   "boolean true input property expanded",
			method: "tools/list",
			tools: []*mcp.Tool{
				{Name: "create_event", InputSchema: map[string]any{"type": "object", "properties": map[string]any{"params": true}}},
			},
			wantOutput: map[string]bool{},
			wantParams: map[string]bool{"create_event": true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			next := func(context.Context, string, mcp.Request) (mcp.Result, error) {
				return &mcp.ListToolsResult{Tools: tt.tools}, nil
			}
			// Track whether the middleware mutated the stored tools.
			original := make([]*mcp.Tool, len(tt.tools))
			for i, tool := range tt.tools {
				original[i] = tool
			}

			res, err := sanitizeToolSchemasMiddleware(next)(context.Background(), tt.method, nil)
			if err != nil {
				t.Fatalf("middleware returned error: %v", err)
			}

			if tt.method != "tools/list" {
				// Non-list methods must pass the result through untouched.
				gotRes, ok := res.(*mcp.ListToolsResult)
				if !ok || !reflect.DeepEqual(gotRes, &mcp.ListToolsResult{Tools: tt.tools}) {
					t.Fatalf("non-list method returned modified result: %#v", res)
				}
				return
			}

			listRes := res.(*mcp.ListToolsResult)
			byName := make(map[string]*mcp.Tool, len(listRes.Tools))
			for _, tool := range listRes.Tools {
				byName[tool.Name] = tool
			}

			for name, wantNil := range tt.wantOutput {
				gotNil := byName[name].OutputSchema == nil
				if gotNil != wantNil {
					t.Errorf("%s: OutputSchema nil = %v, want %v", name, gotNil, wantNil)
				}
			}
			for name, wantExpanded := range tt.wantParams {
				props := byName[name].InputSchema.(map[string]any)["properties"].(map[string]any)
				_, gotExpanded := props["params"].(map[string]any)
				if gotExpanded != wantExpanded {
					t.Errorf("%s: params expanded = %v, want %v", name, gotExpanded, wantExpanded)
				}
			}

			// The middleware must not mutate the server's stored tools.
			for i, orig := range original {
				if orig.OutputSchema != nil {
					if !equalAny(orig.OutputSchema, tt.tools[i].OutputSchema) {
						t.Errorf("stored tool %s outputSchema was mutated", orig.Name)
					}
				}
				if orig.InputSchema != nil {
					if !equalAny(orig.InputSchema, tt.tools[i].InputSchema) {
						t.Errorf("stored tool %s inputSchema was mutated", orig.Name)
					}
				}
			}
		})
	}
}

// equalAny reports whether a and b are deeply equal as JSON values.
func equalAny(a, b any) bool {
	return reflect.DeepEqual(a, b)
}
