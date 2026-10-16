package mcp

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// isBasicTypeName reports whether name is a predeclared type a tool handler
// must not return as its output. The MCP spec requires StructuredContent to
// be a JSON object, so handlers must wrap primitives in an object type
// (ListResult, StatusResult, HealthResult, a map, or a struct).
func isBasicTypeName(name string) bool {
	switch name {
	case "string", "bool", "byte", "rune",
		"int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64", "uintptr",
		"float32", "float64", "error", "any":
		return true
	}
	return false
}

// checkToolOutputTypes scans parsed files for mcp.AddTool calls and returns a
// violation for every handler whose second result type cannot serialize as a
// JSON object: a basic type, a bare slice/array, a pointer to a basic type,
// or an interface. It is the build-time guardrail for the bug class fixed by
// StatusResult/HealthResult (primitives) and ListResult (bare arrays).
func checkToolOutputTypes(fset *token.FileSet, files []*ast.File) []string {
	var violations []string
	found := 0
	for _, f := range files {
		ast.Inspect(f, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			sel, ok := call.Fun.(*ast.SelectorExpr)
			if !ok || sel.Sel.Name != "AddTool" {
				return true
			}
			if pkg, ok := sel.X.(*ast.Ident); !ok || pkg.Name != "mcp" {
				return true
			}
			if len(call.Args) < 3 {
				return true
			}
			fn, ok := call.Args[2].(*ast.FuncLit)
			if !ok {
				violations = append(violations,
					fset.Position(call.Pos()).String()+": AddTool handler is not a function literal; inspect its output type manually")
				return true
			}
			if fn.Type.Results == nil || len(fn.Type.Results.List) < 2 {
				return true
			}
			found++
			out := fn.Type.Results.List[1].Type
			for {
				paren, ok := out.(*ast.ParenExpr)
				if !ok {
					break
				}
				out = paren.X
			}
			pos := fset.Position(out.Pos()).String()
			switch t := out.(type) {
			case *ast.Ident:
				if isBasicTypeName(t.Name) {
					violations = append(violations,
						pos+": tool output type "+t.Name+" is not an object; wrap it (ListResult/StatusResult/HealthResult/map/struct)")
				}
			case *ast.ArrayType:
				violations = append(violations,
					pos+": tool output is a bare slice/array; wrap it in ListResult[T]")
			case *ast.StarExpr:
				if id, ok := t.X.(*ast.Ident); ok && isBasicTypeName(id.Name) {
					violations = append(violations,
						pos+": tool output type *"+id.Name+" points at a non-object; wrap it in an object type")
				}
			case *ast.InterfaceType:
				violations = append(violations,
					pos+": tool output type is an interface; use a concrete object type")
			}
			return true
		})
	}
	if found == 0 {
		violations = append(violations, "no mcp.AddTool handlers found; the guardrail may be scanning the wrong directory")
	}
	return violations
}

// packageFiles parses every non-test .go file in the given directory.
func packageFiles(t *testing.T, dir string) (*token.FileSet, []*ast.File) {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read dir: %v", err)
	}
	fset := token.NewFileSet()
	var files []*ast.File
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}
		f, err := parser.ParseFile(fset, filepath.Join(dir, name), nil, 0)
		if err != nil {
			t.Fatalf("parse %s: %v", name, err)
		}
		files = append(files, f)
	}
	return fset, files
}

// TestToolOutputsAreObjects fails the build if any registered tool returns a
// primitive, bare array, or interface as its typed output — those emit a
// non-object StructuredContent that strict MCP clients (TS SDK 1.29+)
// reject.
func TestToolOutputsAreObjects(t *testing.T) {
	fset, files := packageFiles(t, ".")
	if len(files) == 0 {
		t.Fatal("no source files found in package directory")
	}
	if v := checkToolOutputTypes(fset, files); len(v) > 0 {
		t.Errorf("found %d tool output violations:\n%s", len(v), strings.Join(v, "\n"))
	}
}

// TestToolOutputGuardrailCatchesPrimitives proves the guardrail actually
// fires: a handler returning a bare string must produce a violation.
func TestToolOutputGuardrailCatchesPrimitives(t *testing.T) {
	const src = `package mcp

import (
	"context"

	"github.com/gumieri/xyops-integrations/sdk"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func registerBad(s *mcp.Server, c *sdk.Client) {
	mcp.AddTool(s, &mcp.Tool{
		Name: "bad_tool",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, string, error) {
		return adaptResult("ok", nil)
	})
}
`
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "fixture.go", src, 0)
	if err != nil {
		t.Fatalf("parse fixture: %v", err)
	}
	v := checkToolOutputTypes(fset, []*ast.File{f})
	if len(v) == 0 {
		t.Fatal("guardrail did not flag a primitive-output handler")
	}
	if !strings.Contains(strings.Join(v, "\n"), "string") {
		t.Errorf("violation does not mention the offending type: %v", v)
	}
}

// TestToolOutputGuardrailAcceptsWrappers proves the guardrail passes handlers
// whose outputs are object-shaped.
func TestToolOutputGuardrailAcceptsWrappers(t *testing.T) {
	const src = `package mcp

import (
	"context"

	"github.com/gumieri/xyops-integrations/sdk"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func registerGood(s *mcp.Server, c *sdk.Client) {
	mcp.AddTool(s, &mcp.Tool{
		Name: "good_status",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in struct {
		ID string ` + "`" + `json:"id"` + "`" + `
	}) (*mcp.CallToolResult, StatusResult, error) {
		return adaptStatusResult("deleted", nil)
	})
	mcp.AddTool(s, &mcp.Tool{
		Name: "good_list",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, ListResult[sdk.Tag], error) {
		return adaptListResult(c.Tags.List(ctx))
	})
	mcp.AddTool(s, &mcp.Tool{
		Name: "good_map",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, map[string]any, error) {
		return adaptResult(map[string]any{}, nil)
	})
	mcp.AddTool(s, &mcp.Tool{
		Name: "good_struct",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, *sdk.Tag, error) {
		return adaptResult(c.Tags.Get(ctx, "t1"))
	})
}
`
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "fixture.go", src, 0)
	if err != nil {
		t.Fatalf("parse fixture: %v", err)
	}
	if v := checkToolOutputTypes(fset, []*ast.File{f}); len(v) > 0 {
		t.Errorf("guardrail flagged object-shaped outputs: %v", v)
	}
}
