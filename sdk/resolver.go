package sdk

import (
	"context"
	"fmt"
	"sync"
)

// Resolver resolves human-readable keys to internal IDs.
// It queries list endpoints on cache miss and caches results.
type Resolver struct {
	client *Client
	mu     sync.RWMutex
	cache  map[string]map[string]string // resourceType -> key -> id
}

// NewResolver creates a key→ID resolver backed by the client.
func NewResolver(client *Client) *Resolver {
	return &Resolver{
		client: client,
		cache:  make(map[string]map[string]string),
	}
}

// Seed pre-populates the cache with known key→ID mappings.
func (r *Resolver) Seed(resourceType string, idMap map[string]string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.cache[resourceType] == nil {
		r.cache[resourceType] = make(map[string]string)
	}
	for k, v := range idMap {
		r.cache[resourceType][k] = v
	}
}

// AutoRegister caches a key→ID mapping for a newly created resource.
func (r *Resolver) AutoRegister(resourceType, key, id string) {
	if key == "" || id == "" {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.cache[resourceType] == nil {
		r.cache[resourceType] = make(map[string]string)
	}
	r.cache[resourceType][key] = id
}

// Clear invalidates the cache for a resource type (or all types if empty).
func (r *Resolver) Clear(resourceType string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if resourceType == "" {
		r.cache = make(map[string]map[string]string)
		return
	}
	delete(r.cache, resourceType)
}

// Resolve returns the internal ID for a key, querying the list endpoint on miss.
func (r *Resolver) Resolve(ctx context.Context, resourceType, key string) (string, error) {
	if key == "" {
		return "", fmt.Errorf("empty key for %s", resourceType)
	}

	// Check cache first.
	r.mu.RLock()
	id, ok := r.cache[resourceType][key]
	r.mu.RUnlock()
	if ok {
		return id, nil
	}

	// Query the list endpoint and rebuild the cache for this type.
	// All map access is guarded by r.mu, so concurrent Resolve calls may each
	// trigger a redundant refresh, but they never read partially-written data.
	// This is a deliberate trade-off: single-flight refresh would add locking
	// complexity for negligible benefit on list endpoints.
	if err := r.refresh(ctx, resourceType); err != nil {
		return "", err
	}

	r.mu.RLock()
	id, ok = r.cache[resourceType][key]
	r.mu.RUnlock()
	if !ok {
		return "", fmt.Errorf("key %q not found in %s", key, resourceType)
	}
	return id, nil
}

// refresh rebuilds the cache for a single resource type from its list endpoint.
func (r *Resolver) refresh(ctx context.Context, resourceType string) error {
	pairs, err := r.listPairs(ctx, resourceType)
	if err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	if r.cache[resourceType] == nil {
		r.cache[resourceType] = make(map[string]string)
	}
	for _, p := range pairs {
		r.cache[resourceType][p.key] = p.id
	}
	return nil
}

type keyIDPair struct{ key, id string }

// pairList lists items and maps each to a key→ID pair. The extractor
// functions isolate the differing field accessors so the dispatch switch in
// listPairs stays a single-statement-per-case.
func pairList[T any](ctx context.Context, list func(context.Context) ([]T, error), id, title func(T) string) ([]keyIDPair, error) {
	items, err := list(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]keyIDPair, 0, len(items))
	for _, it := range items {
		out = append(out, keyIDPair{key: title(it), id: id(it)})
	}
	return out, nil
}

// listPairs returns key→ID pairs from the list endpoint for a resource type.
func (r *Resolver) listPairs(ctx context.Context, resourceType string) ([]keyIDPair, error) {
	switch resourceType {
	case "tag":
		return pairList(ctx, r.client.Tags.List, func(t Tag) string { return t.ID }, func(t Tag) string { return t.Title })
	case "group":
		return pairList(ctx, r.client.Groups.List, func(g Group) string { return g.ID }, func(g Group) string { return g.Title })
	case "category":
		return pairList(ctx, r.client.Categories.List, func(c Category) string { return c.ID }, func(c Category) string { return c.Title })
	case "plugin":
		return pairList(ctx, r.client.Plugins.List, func(p Plugin) string { return p.ID }, func(p Plugin) string { return p.Title })
	case "monitor":
		return pairList(ctx, r.client.Monitors.List, func(m Monitor) string { return m.ID }, func(m Monitor) string { return m.Title })
	case "role":
		return pairList(ctx, r.client.Roles.List, func(rl Role) string { return rl.ID }, func(rl Role) string { return rl.Title })
	case "alert":
		return pairList(ctx, r.client.Alerts.List, func(a Alert) string { return a.ID }, func(a Alert) string { return a.Title })
	case "api_key":
		return pairList(ctx, r.client.APIKeys.List, func(k APIKey) string { return k.ID }, func(k APIKey) string { return k.Title })
	case "webhook":
		return pairList(ctx, r.client.WebHooks.List, func(h WebHook) string { return h.ID }, func(h WebHook) string { return h.Title })
	case "event":
		return pairList(ctx, r.client.Events.List, func(e Event) string { return e.ID }, func(e Event) string { return e.Title })
	case "channel":
		return pairList(ctx, r.client.Channels.List, func(ch Channel) string { return ch.ID }, func(ch Channel) string { return ch.Title })
	case "secret":
		return pairList(ctx, r.client.Secrets.List, func(sc Secret) string { return sc.ID }, func(sc Secret) string { return sc.Title })
	case "bucket":
		return pairList(ctx, r.client.Buckets.List, func(b Bucket) string { return b.ID }, func(b Bucket) string { return b.Title })
	default:
		return nil, fmt.Errorf("unknown resource type %q", resourceType)
	}
}
