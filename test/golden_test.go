package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"sort"
	"strings"
	"testing"

	"github.com/gumieri/xyops-integrations/sdk"
)

// Golden wire contracts (XYOINT-15). xyops ships no machine-readable API
// spec, so the external truth is captured instead of hand-duplicated: raw
// response bodies from a pinned version live in testdata/goldens/<version>/.
//
// Two modes:
//
//   - Capture (live): XYOPS_CAPTURE_GOLDENS=1 plus credentials re-fetches
//     every endpoint below and rewrites testdata/goldens/<version>/. Run it
//     on every xyops version bump — the capture diff is the API changelog.
//     Reads only: capture never mutates the instance.
//   - Offline (CI, the default): every golden on disk must be a success
//     envelope AND unmarshal into the matching SDK response type. Runs with
//     no xyops instance, so drift against captured versions surfaces in CI.

type goldenEndpoint struct {
	name   string // file stem
	path   string // request path; relative paths are resolved under /api/app/
	method string // "" → GET; "POST" → POST (body may be nil)
	body   any    // optional JSON body for POST
	// parse, when non-nil, additionally unmarshals the body into the SDK's
	// response type. Endpoints with variable shapes (search, active jobs)
	// leave it nil and only get the envelope check.
	parse func(t *testing.T, body []byte)
	// strip lists top-level or dotted-path response keys (e.g.
	// "config.intl") removed before writing the golden. Used to drop large
	// non-contract payload (locale catalogs, etc.) that would drown
	// re-capture diffs in noise.
	strip []string
}

// typedParse returns a parse func that unmarshals into T.
func typedParse[T any]() func(t *testing.T, body []byte) {
	return func(t *testing.T, body []byte) {
		t.Helper()
		var v T
		if err := json.Unmarshal(body, &v); err != nil {
			t.Errorf("unmarshal %T: %v", v, err)
		}
	}
}

var goldenEndpoints = []goldenEndpoint{
	{name: "ping", path: "ping/v1", parse: typedParse[struct {
		Code int `json:"code"`
	}]()},
	{name: "tags_list", path: "get_tags/v1", parse: typedParse[sdk.TagListResponse]()},
	{name: "groups_list", path: "get_groups/v1", parse: typedParse[sdk.GroupListResponse]()},
	{name: "categories_list", path: "get_categories/v1", parse: typedParse[sdk.CategoryListResponse]()},
	{name: "plugins_list", path: "get_plugins/v1", parse: typedParse[sdk.PluginListResponse]()},
	{name: "events_list", path: "get_events/v1", parse: typedParse[sdk.EventListResponse]()},
	{name: "monitors_list", path: "get_monitors/v1", parse: typedParse[sdk.MonitorListResponse]()},
	{name: "alerts_list", path: "get_alerts/v1", parse: typedParse[sdk.AlertListResponse]()},
	{name: "roles_list", path: "get_roles/v1", parse: typedParse[sdk.RoleListResponse]()},
	{name: "api_keys_list", path: "get_api_keys/v1", parse: typedParse[sdk.APIKeyListResponse]()},
	{name: "web_hooks_list", path: "get_web_hooks/v1", parse: typedParse[sdk.WebHookListResponse]()},
	{name: "channels_list", path: "get_channels/v1", parse: typedParse[sdk.ChannelListResponse]()},
	{name: "secrets_list", path: "get_secrets/v1", parse: typedParse[sdk.SecretListResponse]()},
	{name: "buckets_list", path: "get_buckets/v1", parse: typedParse[sdk.BucketListResponse]()},
	{name: "servers_list", path: "get_servers/v1", parse: typedParse[sdk.ServerListResponse]()},
	{name: "config_get", path: "admin_get_config/v1", method: "POST", parse: typedParse[sdk.ConfigResponse](), strip: []string{"config.intl"}},
	{name: "get_multiple", path: "get_multiple/v1?lists=tags,groups,events", parse: typedParse[sdk.GetMultipleResponse]()},
	{name: "search_jobs", path: "search_jobs/v1?limit=5"},       // shape varies with content; envelope check only
	{name: "search_tickets", path: "search_tickets/v1?limit=5"}, // shape varies with content; envelope check only
	{name: "active_jobs", path: "get_active_jobs/v1"},           // shape varies with content; envelope check only
}

// goldenVersion returns the conductor version used to name capture
// directories, or "" when unknown.
func goldenVersion(t *testing.T) string {
	t.Helper()
	masters, err := xyopsClient.Servers.Masters(t.Context())
	if err != nil {
		t.Fatalf("list masters: %v", err)
	}
	for _, m := range masters {
		if m.Version != "" {
			return m.Version
		}
	}
	return ""
}

// TestGoldenWireContracts runs in whichever mode the environment selects.
func TestGoldenWireContracts(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping golden contracts in short mode")
	}

	if os.Getenv("XYOPS_CAPTURE_GOLDENS") == "1" {
		skipIfNoClient(t)
		captureGoldens(t)
		return
	}

	assertGoldensOffline(t)
}

// captureGoldens re-fetches every endpoint and rewrites
// testdata/goldens/<version>/.
func captureGoldens(t *testing.T) {
	t.Helper()
	version := goldenVersion(t)
	if version == "" {
		t.Fatal("cannot capture: conductor version not reported")
	}

	dir := filepath.Join("testdata", "goldens", version)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("create %s: %v", dir, err)
	}

	ctx := context.Background()
	for _, ep := range goldenEndpoints {
		var (
			body []byte
			err  error
		)
		if ep.method == "POST" {
			body, err = xyopsClient.DoPOST(ctx, ep.path, ep.body)
		} else {
			body, err = xyopsClient.DoGET(ctx, ep.path)
		}
		if err != nil {
			t.Errorf("capture %s: %v", ep.name, err)
			continue
		}

		var canonical any
		if err := json.Unmarshal(body, &canonical); err != nil {
			t.Errorf("capture %s: response is not JSON: %v", ep.name, err)
			continue
		}
		if obj, ok := canonical.(map[string]any); ok && len(ep.strip) > 0 {
			stripKeys(obj, ep.strip)
		}
		out, err := json.MarshalIndent(canonical, "", "  ")
		if err != nil {
			t.Fatalf("canonicalize %s: %v", ep.name, err)
		}

		// The repo is public: instance-identifying values (hosts, internal
		// IPs) must never be tracked (AGENTS.md privacy convention).
		out = sanitizeGoldenIdentities(t, out)

		file := filepath.Join(dir, ep.name+".json")
		if err := os.WriteFile(file, append(out, '\n'), 0o644); err != nil {
			t.Fatalf("write %s: %v", file, err)
		}
		t.Logf("captured %s (%d bytes)", file, len(out))
	}
}

// goldenScrubPair is one instance-identifying string and its placeholder.
type goldenScrubPair struct{ from, to string }

// goldenPrivateIP matches RFC1918 literals (optionally ::ffff:-prefixed);
// they are rewritten to the TEST-NET-3 documentation address.
var goldenPrivateIP = regexp.MustCompile(`(::ffff:)?((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)\.){3}(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)`)

func isPrivateIPv4(ip string) bool {
	ip = strings.TrimPrefix(ip, "::ffff:")
	parts := strings.Split(ip, ".")
	if len(parts) != 4 {
		return false
	}
	n := make([]int, 4)
	for i, p := range parts {
		v := 0
		for _, c := range p {
			if c < '0' || c > '9' {
				return false
			}
			v = v*10 + int(c-'0')
		}
		if v > 255 {
			return false
		}
		n[i] = v
	}
	return n[0] == 10 || (n[0] == 172 && n[1] >= 16 && n[1] <= 31) || (n[0] == 192 && n[1] == 168)
}

// sanitizeGoldenIdentities replaces instance-identifying values harvested
// from the live environment (base URL host, master IDs, server
// hostnames/IPs/keywords) with nxample.com / TEST-NET-3 placeholders, then
// fails closed if any harvested identity survived.
func sanitizeGoldenIdentities(t *testing.T, out []byte) []byte {
	t.Helper()

	var pairs []goldenScrubPair
	add := func(from, to string) {
		if from != "" && from != to {
			pairs = append(pairs, goldenScrubPair{from, to})
		}
	}
	if base := os.Getenv("XYOPS_BASE_URL"); base != "" {
		if u, err := url.Parse(base); err == nil && u.Hostname() != "" {
			host := u.Hostname()
			add(host, "xyops.nxample.com")
			// The registrable domain also appears inside email addresses
			// and derived values; replace the suffix too (longest rules
			// run first).
			if labels := strings.Split(host, "."); len(labels) >= 2 {
				add(strings.Join(labels[len(labels)-2:], "."), "nxample.com")
			}
		}
	}
	if xyopsClient != nil {
		ctx := context.Background()
		if masters, err := xyopsClient.Servers.Masters(ctx); err == nil {
			for _, m := range masters {
				add(m.ID, "xyops.nxample.com")
			}
		}
		if servers, err := xyopsClient.Servers.List(ctx); err == nil {
			for _, s := range servers {
				add(s.Hostname, "satellite.nxample.com")
				add(s.Keywords, "satellite.nxample.com,203.0.113.10")
				if isPrivateIPv4(s.IP) {
					add(s.IP, "203.0.113.10")
					add("::ffff:"+s.IP, "::ffff:203.0.113.10")
				}
			}
		}
	}

	// Longest first so overlapping identities replace completely.
	sort.Slice(pairs, func(i, j int) bool { return len(pairs[i].from) > len(pairs[j].from) })
	for _, p := range pairs {
		out = bytes.ReplaceAll(out, []byte(p.from), []byte(p.to))
	}
	out = goldenPrivateIP.ReplaceAll(out, []byte("203.0.113.10"))

	// Fail closed: no harvested identity may survive sanitization.
	for _, p := range pairs {
		if bytes.Contains(out, []byte(p.from)) {
			t.Fatalf("sanitization failed: %q still present in golden capture", p.from)
		}
	}
	return out
}

// assertGoldensOffline validates every captured golden against the current
// SDK types. Requires no live instance.
func assertGoldensOffline(t *testing.T) {
	t.Helper()

	root := filepath.Join("testdata", "goldens")
	versions, err := os.ReadDir(root)
	if err != nil || len(versions) == 0 {
		t.Skip("no golden captures yet; run with XYOPS_CAPTURE_GOLDENS=1 against a live instance")
	}

	for _, ver := range versions {
		if !ver.IsDir() {
			continue
		}
		for _, ep := range goldenEndpoints {
			file := filepath.Join(root, ver.Name(), ep.name+".json")
			body, err := os.ReadFile(file)
			if err != nil {
				continue // not captured for this version
			}
			t.Run(ver.Name()+"/"+ep.name, func(t *testing.T) {
				assertGoldenEnvelope(t, body)
				if ep.parse != nil {
					ep.parse(t, body)
				}
			})
		}
	}
}

// stripKeys deletes dotted-path keys ("a.b.c") from a decoded JSON object.
// Missing path segments are ignored silently.
func stripKeys(obj map[string]any, paths []string) {
	for _, p := range paths {
		parts := splitPath(p)
		cur := obj
		for i, key := range parts {
			if i == len(parts)-1 {
				delete(cur, key)
				break
			}
			next, ok := cur[key].(map[string]any)
			if !ok {
				break
			}
			cur = next
		}
	}
}

func splitPath(p string) []string {
	var parts []string
	for _, s := range strings.Split(p, ".") {
		if s != "" {
			parts = append(parts, s)
		}
	}
	return parts
}

// assertGoldenEnvelope asserts the body is a success envelope: valid JSON
// whose envelope carries no error description (success bodies never do —
// see sdk.checkErrorResponse). Captures only ever store success responses,
// so a drifted re-parse into an error body is a contract break.
func assertGoldenEnvelope(t *testing.T, body []byte) {
	t.Helper()
	var env struct {
		Code        any    `json:"code"`
		Description string `json:"description"`
	}
	if err := json.Unmarshal(body, &env); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if env.Description != "" {
		t.Fatalf("envelope carries an error description: %s", env.Description)
	}
	switch v := env.Code.(type) {
	case nil:
		// Envelope-less documents (e.g. satellite/config) are acceptable.
	case float64:
		if v != 0 {
			t.Fatalf("non-zero envelope code: %v", v)
		}
	case string:
		if !slices.Contains([]string{"0"}, v) {
			t.Fatalf("non-success envelope code: %q", v)
		}
	}
}
