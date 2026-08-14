package mcp

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gumieri/xyops-integrations/sdk"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Version is the MCP server implementation version.
// Overridden by cmd/xyops-mcp from the ldflags-injected main.version.
var Version = "0.1.0"

// NewServer creates the xyops MCP server with all tools registered.
func NewServer(c *sdk.Client) *mcp.Server {
	s := mcp.NewServer(&mcp.Implementation{Name: "xyops", Version: Version}, nil)

	registerTagTools(s, c)
	registerGroupTools(s, c)
	registerCategoryTools(s, c)
	registerPluginTools(s, c)
	registerEventTools(s, c)
	registerMonitorTools(s, c)
	registerAlertTools(s, c)
	registerRoleTools(s, c)
	registerAPIKeyTools(s, c)
	registerWebHookTools(s, c)
	registerServerTools(s, c)
	registerConfigTools(s, c)
	registerJobTools(s, c)

	s.AddReceivingMiddleware(sanitizeToolSchemasMiddleware)

	return s
}

// ListResult wraps a slice result so list tools return an object-shaped
// StructuredContent ({"items": [...]}) instead of a bare array. Strict client
// validators (e.g. the TS SDK 1.29.0 that opencode bundles) require
// StructuredContent to be a record, not an array.
type ListResult[T any] struct {
	Items []T `json:"items"`
}

// adaptResult adapts a (T, error) pair into the (result, output, error) shape
// expected by mcp.ToolHandlerFor.
func adaptResult[T any](v T, err error) (*mcp.CallToolResult, T, error) {
	return nil, v, err
}

// adaptListResult adapts a slice-returning list tool's ([]T, error) pair into
// the (result, output, error) shape expected by mcp.ToolHandlerFor, wrapping
// the slice in a [ListResult] object.
func adaptListResult[T any](v []T, err error) (*mcp.CallToolResult, ListResult[T], error) {
	if err != nil {
		return nil, ListResult[T]{}, err
	}
	return nil, ListResult[T]{Items: v}, nil
}

// sanitizeToolSchemasMiddleware rewrites the tools/list response so every
// tool's schemas pass client-side validation.
//
// The Go SDK infers JSON Schemas via reflection. Two inferred shapes break
// strict client validators (e.g. the TS SDK 1.29.0 ToolSchema that opencode
// bundles):
//
//   - outputSchema for tools returning a primitive (string, int, ...) or a
//     slice gets a non-object root type, but clients require an object schema
//     (or no schema at all). Since outputSchema is optional per the spec, we
//     drop it entirely.
//   - inputSchema/outputSchema properties for fields typed as `any`
//     (interface{}) collapse to the JSON Schema boolean `true`, which strict
//     clients reject as "not an object schema". We expand boolean property
//     schemas to the open schema `{}` (JSON Schema for "anything").
//
// listTools returns the server's stored *Tool values, so we clone each tool
// before mutating its schemas: writing in place would both race under
// concurrent tools/list requests and silently rewrite the server's source of
// truth. The clones only affect this response.
func sanitizeToolSchemasMiddleware(next mcp.MethodHandler) mcp.MethodHandler {
	return func(ctx context.Context, method string, req mcp.Request) (mcp.Result, error) {
		res, err := next(ctx, method, req)
		if err != nil {
			return nil, err
		}
		if method != "tools/list" {
			return res, nil
		}
		listRes, ok := res.(*mcp.ListToolsResult)
		if !ok {
			return nil, fmt.Errorf("tools/list returned unexpected result type %T", res)
		}
		// Build a fresh slice and clone each tool before mutating its
		// schemas: writing in place would both race under concurrent
		// tools/list requests and silently rewrite the server's source of
		// truth (and the caller's slice). The clones only affect this
		// response.
		cloned := make([]*mcp.Tool, 0, len(listRes.Tools))
		for _, tool := range listRes.Tools {
			clone := *tool
			clone.OutputSchema = normalizeOutputSchema(tool.OutputSchema)
			clone.InputSchema = normalizeSchema(tool.InputSchema)
			cloned = append(cloned, &clone)
		}
		listRes.Tools = cloned
		return res, nil
	}
}

// normalizeOutputSchema normalizes a tool's output schema, or drops it
// entirely when it is not an object schema. Non-object output schemas (from
// tools returning primitives or slices) are optional per the spec, but strict
// clients reject them, so we omit them.
func normalizeOutputSchema(s any) any {
	if !isObjectSchema(s) {
		return nil
	}
	return normalizeSchema(s)
}

// isObjectSchema reports whether s marshals to a JSON Schema whose root type
// is "object". A nil or absent schema is not an object schema.
func isObjectSchema(s any) bool {
	if s == nil {
		return false
	}
	var m map[string]any
	if err := remarshalSchema(s, &m); err != nil {
		return false
	}
	typ, _ := m["type"].(string)
	return typ == "object"
}

// normalizeSchema rewrites a JSON Schema for client compatibility: boolean
// property schemas (`true`) inside `properties` become the open schema `{}`,
// which strict clients accept while preserving "anything" semantics.
func normalizeSchema(s any) any {
	if s == nil {
		return nil
	}
	var m map[string]any
	if err := remarshalSchema(s, &m); err != nil {
		return s
	}
	normalizeSchemaMap(m)
	return m
}

// normalizeSchemaMap rewrites a schema in place so every boolean `true`
// property schema becomes the open schema `{}` (JSON Schema for "anything").
//
// Boolean schemas are only valid as values inside a `properties` map, so only
// those are expanded. Keywords like `additionalProperties` are left untouched
// so their meaning is preserved; they are still descended into, because a
// keyword may itself hold a nested schema containing its own `properties`.
func normalizeSchemaMap(m map[string]any) {
	if props, ok := m["properties"].(map[string]any); ok {
		for k, v := range props {
			switch v := v.(type) {
			case bool:
				if v {
					props[k] = map[string]any{}
				}
			case map[string]any:
				normalizeSchemaMap(v)
			}
		}
	}
	for _, key := range []string{"items", "additionalProperties", "contains", "not", "propertyNames", "if", "then", "else", "oneOf", "anyOf", "allOf"} {
		switch v := m[key].(type) {
		case map[string]any:
			normalizeSchemaMap(v)
		case []any:
			for _, item := range v {
				if cm, ok := item.(map[string]any); ok {
					normalizeSchemaMap(cm)
				}
			}
		}
	}
}

// remarshalSchema converts any schema representation (map, jsonschema.Schema,
// json.RawMessage) into the target value via JSON round-tripping.
func remarshalSchema(src any, dst any) error {
	b, err := json.Marshal(src)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, dst)
}
