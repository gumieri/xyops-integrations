package client

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

// listPairs returns key→ID pairs from the list endpoint for a resource type.
func (r *Resolver) listPairs(ctx context.Context, resourceType string) ([]keyIDPair, error) {
	switch resourceType {
	case "tag":
		tags, err := r.client.Tags.List(ctx)
		if err != nil {
			return nil, err
		}
		out := make([]keyIDPair, 0, len(tags))
		for _, t := range tags {
			out = append(out, keyIDPair{key: t.Title, id: t.ID})
		}
		return out, nil
	case "group":
		groups, err := r.client.Groups.List(ctx)
		if err != nil {
			return nil, err
		}
		out := make([]keyIDPair, 0, len(groups))
		for _, g := range groups {
			out = append(out, keyIDPair{key: g.Title, id: g.ID})
		}
		return out, nil
	case "category":
		cats, err := r.client.Categories.List(ctx)
		if err != nil {
			return nil, err
		}
		out := make([]keyIDPair, 0, len(cats))
		for _, c := range cats {
			out = append(out, keyIDPair{key: c.Title, id: c.ID})
		}
		return out, nil
	case "plugin":
		plugins, err := r.client.Plugins.List(ctx)
		if err != nil {
			return nil, err
		}
		out := make([]keyIDPair, 0, len(plugins))
		for _, p := range plugins {
			out = append(out, keyIDPair{key: p.Title, id: p.ID})
		}
		return out, nil
	case "monitor":
		monitors, err := r.client.Monitors.List(ctx)
		if err != nil {
			return nil, err
		}
		out := make([]keyIDPair, 0, len(monitors))
		for _, m := range monitors {
			out = append(out, keyIDPair{key: m.Title, id: m.ID})
		}
		return out, nil
	case "role":
		roles, err := r.client.Roles.List(ctx)
		if err != nil {
			return nil, err
		}
		out := make([]keyIDPair, 0, len(roles))
		for _, rl := range roles {
			out = append(out, keyIDPair{key: rl.Title, id: rl.ID})
		}
		return out, nil
	default:
		return nil, fmt.Errorf("unknown resource type %q", resourceType)
	}
}
