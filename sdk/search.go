package sdk

import (
	"context"
	"strconv"
	"strings"
)

// SearchQueryOpts configures the search_*/v1 endpoints. All search endpoints share
// the base params (query, offset, limit, sort_by, sort_dir).
type SearchQueryOpts struct {
	Query   string // Unbase/PxQL syntax, default "*"
	Offset  int
	Limit   int
	SortBy  string // default "_id"
	SortDir int    // 1 ascending, -1 descending
	Verbose bool   // search_snapshots only: include heavy nested fields
}

// SearchService
func (s *SearchService) search(ctx context.Context, path string, opts SearchQueryOpts, extra map[string]string) ([]byte, error) {
	params := map[string]string{}
	if opts.Query != "" {
		params["query"] = opts.Query
	}
	if opts.Offset > 0 {
		params["offset"] = strconv.Itoa(opts.Offset)
	}
	if opts.Limit > 0 {
		params["limit"] = strconv.Itoa(opts.Limit)
	}
	if opts.SortBy != "" {
		params["sort_by"] = opts.SortBy
	}
	if opts.SortDir != 0 {
		params["sort_dir"] = strconv.Itoa(opts.SortDir)
	}
	if opts.Verbose {
		params["verbose"] = "1"
	}
	for k, v := range extra {
		params[k] = v
	}
	return s.client.DoGETWithQuery(ctx, path, params)
}

func (s *SearchService) SearchServers(ctx context.Context, opts SearchQueryOpts) ([]Server, error) {
	body, err := s.search(ctx, "search_servers/v1", opts, nil)
	if err != nil {
		return nil, err
	}
	var resp struct {
		Code int      `json:"code"`
		Rows []Server `json:"rows"`
	}
	if err := unmarshalJSON(body, &resp); err != nil {
		return nil, err
	}
	return resp.Rows, nil
}

func (s *SearchService) SearchAlerts(ctx context.Context, opts SearchQueryOpts) ([]AlertInvocation, error) {
	body, err := s.search(ctx, "search_alerts/v1", opts, nil)
	if err != nil {
		return nil, err
	}
	var resp struct {
		Code int               `json:"code"`
		Rows []AlertInvocation `json:"rows"`
	}
	if err := unmarshalJSON(body, &resp); err != nil {
		return nil, err
	}
	return resp.Rows, nil
}

func (s *SearchService) SearchSnapshots(ctx context.Context, opts SearchQueryOpts) ([]map[string]any, error) {
	body, err := s.search(ctx, "search_snapshots/v1", opts, nil)
	if err != nil {
		return nil, err
	}
	var resp struct {
		Code int              `json:"code"`
		Rows []map[string]any `json:"rows"`
	}
	if err := unmarshalJSON(body, &resp); err != nil {
		return nil, err
	}
	return resp.Rows, nil
}

// SearchActivity searches the admin activity log (admin only).
func (s *SearchService) SearchActivity(ctx context.Context, opts SearchQueryOpts) ([]map[string]any, error) {
	body, err := s.search(ctx, "search_activity/v1", opts, nil)
	if err != nil {
		return nil, err
	}
	var resp struct {
		Code int              `json:"code"`
		Rows []map[string]any `json:"rows"`
	}
	if err := unmarshalJSON(body, &resp); err != nil {
		return nil, err
	}
	return resp.Rows, nil
}

// SearchRevisionHistory searches revision history for one resource type.
// typ is one of: alerts, categories, channels, events, groups, monitors,
// plugins, tags, web_hooks, buckets, secrets, tickets, roles.
func (s *SearchService) SearchRevisionHistory(ctx context.Context, typ string, opts SearchQueryOpts) ([]map[string]any, error) {
	body, err := s.search(ctx, "search_revision_history/v1", opts, map[string]string{"type": typ})
	if err != nil {
		return nil, err
	}
	var resp struct {
		Code int              `json:"code"`
		Rows []map[string]any `json:"rows"`
	}
	if err := unmarshalJSON(body, &resp); err != nil {
		return nil, err
	}
	return resp.Rows, nil
}

// SearchStatHistory searches stats history. offset/limit are in DAYS (unlike
// other search endpoints). The response uses an "items" key, not "rows".
func (s *SearchService) SearchStatHistory(ctx context.Context, offset, limit int, path, keyPrefix string, currentDay bool) ([]map[string]any, error) {
	params := map[string]string{}
	if offset > 0 {
		params["offset"] = strconv.Itoa(offset)
	}
	if limit > 0 {
		params["limit"] = strconv.Itoa(limit)
	}
	if path != "" {
		params["path"] = path
	}
	if keyPrefix != "" {
		params["key_prefix"] = keyPrefix
	}
	if currentDay {
		params["current_day"] = "1"
	}
	body, err := s.client.DoGETWithQuery(ctx, "search_stat_history/v1", params)
	if err != nil {
		return nil, err
	}
	var resp struct {
		Code  int              `json:"code"`
		Items []map[string]any `json:"items"`
	}
	if err := unmarshalJSON(body, &resp); err != nil {
		return nil, err
	}
	return resp.Items, nil
}

// BulkExport streams a search result set as a file (csv/tsv/ndjson). Returns
// raw bytes. compress enables gzip encoding.
func (s *SearchService) BulkExport(ctx context.Context, index string, columns []string, format, query string, compress bool) ([]byte, error) {
	params := map[string]string{
		"index":   index,
		"columns": strings.Join(columns, ","),
		"format":  format,
	}
	if query != "" {
		params["query"] = query
	}
	if compress {
		params["compress"] = "1"
	}
	return s.client.DoGETWithQuery(ctx, "bulk_search_export/v1", params)
}

// MarketplaceOpts configures marketplace/v1 (admin only).
type MarketplaceOpts struct {
	Query      string
	Type       string
	PluginType string
	Author     string
	Status     string
	License    string
	Tags       []string
	Requires   []string
	SortBy     string
	SortDir    int
	Offset     int
	Limit      int
}

// Marketplace searches the plugin marketplace catalog (admin only).
func (s *SearchService) Marketplace(ctx context.Context, opts MarketplaceOpts) ([]map[string]any, error) {
	params := map[string]string{}
	if opts.Query != "" {
		params["query"] = opts.Query
	}
	if opts.Type != "" {
		params["type"] = opts.Type
	}
	if opts.PluginType != "" {
		params["plugin_type"] = opts.PluginType
	}
	if opts.Author != "" {
		params["author"] = opts.Author
	}
	if opts.Status != "" {
		params["status"] = opts.Status
	}
	if opts.License != "" {
		params["license"] = opts.License
	}
	if len(opts.Tags) > 0 {
		params["tags"] = strings.Join(opts.Tags, ",")
	}
	if len(opts.Requires) > 0 {
		params["requires"] = strings.Join(opts.Requires, ",")
	}
	if opts.SortBy != "" {
		params["sort_by"] = opts.SortBy
	}
	if opts.SortDir != 0 {
		params["sort_dir"] = strconv.Itoa(opts.SortDir)
	}
	if opts.Offset > 0 {
		params["offset"] = strconv.Itoa(opts.Offset)
	}
	if opts.Limit > 0 {
		params["limit"] = strconv.Itoa(opts.Limit)
	}
	body, err := s.client.DoGETWithQuery(ctx, "marketplace/v1", params)
	if err != nil {
		return nil, err
	}
	var resp struct {
		Code int              `json:"code"`
		Rows []map[string]any `json:"rows"`
	}
	if err := unmarshalJSON(body, &resp); err != nil {
		return nil, err
	}
	return resp.Rows, nil
}
