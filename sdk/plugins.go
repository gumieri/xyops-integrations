package sdk

import (
	"context"
	"strconv"
)

// PluginsService
func (s *PluginsService) List(ctx context.Context) ([]Plugin, error) {
	body, err := s.client.DoGET(ctx, "get_plugins/v1")
	if err != nil {
		return nil, err
	}

	var resp PluginListResponse
	if err := unmarshalJSON(body, &resp); err != nil {
		return nil, err
	}
	return resp.Rows, nil
}

func (s *PluginsService) Get(ctx context.Context, id string) (*Plugin, error) {
	body, err := s.client.DoGETWithQuery(ctx, "get_plugin/v1", map[string]string{"id": id})
	if err != nil {
		return nil, err
	}

	var resp PluginResponse
	if err := unmarshalJSON(body, &resp); err != nil {
		return nil, err
	}
	return resp.Plugin, nil
}

// Create always calls create_plugin/v1 — plugins have NO upsert.
func (s *PluginsService) Create(ctx context.Context, plugin *Plugin) (*Plugin, error) {
	body, err := s.client.DoPOST(ctx, "create_plugin/v1", plugin)
	if err != nil {
		return nil, err
	}

	var resp PluginResponse
	if err := unmarshalJSON(body, &resp); err != nil {
		return nil, err
	}
	return resp.Plugin, nil
}

func (s *PluginsService) Update(ctx context.Context, plugin *Plugin) (*Plugin, error) {
	_, err := s.client.DoPOST(ctx, "update_plugin/v1", plugin)
	if err != nil {
		return nil, err
	}

	return s.Get(ctx, plugin.ID)
}

func (s *PluginsService) Delete(ctx context.Context, id string) error {
	_, err := s.client.DoPOST(ctx, "delete_plugin/v1", map[string]string{"id": id})
	return err
}

// TestMonitor runs a monitor plugin against a live server via
// test_monitor_plugin/v1. Returns the raw plugin output and stderr.
func (s *PluginsService) TestMonitor(ctx context.Context, id, server string) (result map[string]any, stderr string, err error) {
	body, err := s.client.DoPOST(ctx, "test_monitor_plugin/v1", map[string]string{"id": id, "server": server})
	if err != nil {
		return nil, "", err
	}

	var resp struct {
		Code   int            `json:"code"`
		Result map[string]any `json:"result"`
		Stderr string         `json:"stderr"`
	}
	if err := unmarshalJSON(body, &resp); err != nil {
		return nil, "", err
	}
	return resp.Result, resp.Stderr, nil
}

// TestScheduler runs a scheduler plugin in simulation via
// test_scheduler_plugin/v1. Returns the output data and stdout. Note the API
// returns code 0 even when the plugin itself errors (err:true in the body).
func (s *PluginsService) TestScheduler(ctx context.Context, id string, params map[string]any, timezone string, epoch int64) (data map[string]any, stdout string, err error) {
	payload := map[string]any{"id": id}
	if params != nil {
		payload["params"] = params
	}
	if timezone != "" {
		payload["timezone"] = timezone
	}
	if epoch != 0 {
		payload["epoch"] = epoch
	}
	body, err := s.client.DoPOST(ctx, "test_scheduler_plugin/v1", payload)
	if err != nil {
		return nil, "", err
	}

	var resp struct {
		Code   int            `json:"code"`
		Data   map[string]any `json:"data"`
		Stdout string         `json:"stdout"`
	}
	if err := unmarshalJSON(body, &resp); err != nil {
		return nil, "", err
	}
	return resp.Data, resp.Stdout, nil
}

// EventsService
func (s *EventsService) List(ctx context.Context) ([]Event, error) {
	body, err := s.client.DoGET(ctx, "get_events/v1")
	if err != nil {
		return nil, err
	}

	var resp EventListResponse
	if err := unmarshalJSON(body, &resp); err != nil {
		return nil, err
	}
	return resp.Rows, nil
}

func (s *EventsService) Get(ctx context.Context, id string) (*Event, error) {
	body, err := s.client.DoGETWithQuery(ctx, "get_event/v1", map[string]string{"id": id})
	if err != nil {
		return nil, err
	}

	var resp EventResponse
	if err := unmarshalJSON(body, &resp); err != nil {
		return nil, err
	}
	return resp.Event, nil
}

func (s *EventsService) Create(ctx context.Context, event *Event) (*Event, error) {
	if event.ID != "" {
		return s.Update(ctx, event)
	}

	// Inject workflow markers and default limits on a copy so the caller's
	// struct is never mutated, even if the request fails. The copy is shallow,
	// which is safe here: Event fields are value types (the Workflow pointer is
	// only read, never written by this method).
	create := *event
	if create.Workflow != nil {
		create.Plugin = "_workflow"
		create.Type = "workflow"
	}
	if len(create.Limits) == 0 {
		create.Limits = DefaultLimits()
	}

	body, err := s.client.DoPOST(ctx, "create_event/v1", create)
	if err != nil {
		return nil, err
	}

	var resp EventResponse
	if err := unmarshalJSON(body, &resp); err != nil {
		return nil, err
	}
	return resp.Event, nil
}

func (s *EventsService) Update(ctx context.Context, event *Event) (*Event, error) {
	_, err := s.client.DoPOST(ctx, "update_event/v1", event)
	if err != nil {
		return nil, err
	}

	return s.Get(ctx, event.ID)
}

// DeleteEventOpts controls event deletion behavior.
type DeleteEventOpts struct {
	DeleteJobs bool
}

func (s *EventsService) Delete(ctx context.Context, id string, opts *DeleteEventOpts) error {
	payload := map[string]any{"id": id}
	if opts != nil && opts.DeleteJobs {
		payload["delete_jobs"] = true
	}
	_, err := s.client.DoPOST(ctx, "delete_event/v1", payload)
	return err
}

// MonitorsService
func (s *MonitorsService) List(ctx context.Context) ([]Monitor, error) {
	body, err := s.client.DoGET(ctx, "get_monitors/v1")
	if err != nil {
		return nil, err
	}

	var resp MonitorListResponse
	if err := unmarshalJSON(body, &resp); err != nil {
		return nil, err
	}
	return resp.Rows, nil
}

func (s *MonitorsService) Get(ctx context.Context, id string) (*Monitor, error) {
	body, err := s.client.DoGETWithQuery(ctx, "get_monitor/v1", map[string]string{"id": id})
	if err != nil {
		return nil, err
	}

	var resp MonitorResponse
	if err := unmarshalJSON(body, &resp); err != nil {
		return nil, err
	}
	return resp.Monitor, nil
}

func (s *MonitorsService) Create(ctx context.Context, monitor *Monitor) (*Monitor, error) {
	if monitor.ID != "" {
		return s.Update(ctx, monitor)
	}

	body, err := s.client.DoPOST(ctx, "create_monitor/v1", monitor)
	if err != nil {
		return nil, err
	}

	var resp MonitorResponse
	if err := unmarshalJSON(body, &resp); err != nil {
		return nil, err
	}
	return resp.Monitor, nil
}

func (s *MonitorsService) Update(ctx context.Context, monitor *Monitor) (*Monitor, error) {
	_, err := s.client.DoPOST(ctx, "update_monitor/v1", monitor)
	if err != nil {
		return nil, err
	}

	return s.Get(ctx, monitor.ID)
}

func (s *MonitorsService) Delete(ctx context.Context, id string) error {
	_, err := s.client.DoPOST(ctx, "delete_monitor/v1", map[string]string{"id": id})
	return err
}

// Test evaluates a monitor source expression against live server data via
// test_monitor/v1. Returns the numeric result, or (0, true) when the
// expression cannot be evaluated.
func (s *MonitorsService) Test(ctx context.Context, server, source, dataType, dataMatch string) (value float64, fail bool, err error) {
	payload := map[string]string{"server": server, "source": source, "data_type": dataType}
	if dataMatch != "" {
		payload["data_match"] = dataMatch
	}
	body, err := s.client.DoPOST(ctx, "test_monitor/v1", payload)
	if err != nil {
		return 0, false, err
	}

	var resp struct {
		Code  int     `json:"code"`
		Value float64 `json:"value"`
		Fail  bool    `json:"fail"`
	}
	if err := unmarshalJSON(body, &resp); err != nil {
		return 0, false, err
	}
	return resp.Value, resp.Fail, nil
}

// MultiUpdate applies shallow merges to several monitors in one call.
func (s *MonitorsService) MultiUpdate(ctx context.Context, monitors []*Monitor) error {
	if len(monitors) == 0 {
		return nil
	}
	items := make([]map[string]any, 0, len(monitors))
	for _, m := range monitors {
		if m == nil {
			continue
		}
		items = append(items, monitorUpdateItem(m))
	}
	if len(items) == 0 {
		return nil
	}
	_, err := s.client.DoPOST(ctx, "multi_update_monitor/v1", map[string]any{"items": items})
	return err
}

// monitorUpdateItem renders a Monitor as a shallow-merge update item (id + all
// non-zero fields).
func monitorUpdateItem(m *Monitor) map[string]any {
	item := map[string]any{"id": m.ID}
	if m.Title != "" {
		item["title"] = m.Title
	}
	if m.Source != "" {
		item["source"] = m.Source
	}
	if m.DataType != "" {
		item["data_type"] = m.DataType
	}
	item["display"] = m.Display
	if len(m.Groups) > 0 {
		item["groups"] = m.Groups
	}
	if m.Suffix != "" {
		item["suffix"] = m.Suffix
	}
	if m.MinVertScale != 0 {
		item["min_vert_scale"] = m.MinVertScale
	}
	if m.DataMatch != "" {
		item["data_match"] = m.DataMatch
	}
	item["delta"] = m.Delta
	item["divide_time"] = m.DivideTime
	item["zero_min"] = m.ZeroMin
	if m.Notes != "" {
		item["notes"] = m.Notes
	}
	return item
}

// QuickMonData returns the last-60-seconds monitor snapshots keyed by server ID.
func (s *MonitorsService) QuickMonData(ctx context.Context, server, group string) (map[string][]map[string]any, error) {
	params := map[string]string{}
	if server != "" {
		params["server"] = server
	}
	if group != "" {
		params["group"] = group
	}
	body, err := s.client.DoGETWithQuery(ctx, "get_quickmon_data/v1", params)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Code    int                         `json:"code"`
		Servers map[string][]map[string]any `json:"servers"`
	}
	if err := unmarshalJSON(body, &resp); err != nil {
		return nil, err
	}
	return resp.Servers, nil
}

// LatestData returns a server's latest timeline data and full host record via
// get_latest_monitor_data/v1. sys is one of hourly/daily/monthly/yearly.
func (s *MonitorsService) LatestData(ctx context.Context, server, sys string, limit int) (rows []map[string]any, data map[string]any, err error) {
	params := map[string]string{"server": server, "sys": sys, "limit": strconv.Itoa(limit)}
	body, err := s.client.DoGETWithQuery(ctx, "get_latest_monitor_data/v1", params)
	if err != nil {
		return nil, nil, err
	}

	var resp struct {
		Code int              `json:"code"`
		Rows []map[string]any `json:"rows"`
		Data map[string]any   `json:"data"`
	}
	if err := unmarshalJSON(body, &resp); err != nil {
		return nil, nil, err
	}
	return resp.Rows, resp.Data, nil
}

// HistoricalData returns timeline entries for a server starting at a given Unix
// date via get_historical_monitor_data/v1.
func (s *MonitorsService) HistoricalData(ctx context.Context, server, sys string, date int64, limit int) ([]map[string]any, error) {
	params := map[string]string{
		"server": server,
		"sys":    sys,
		"date":   strconv.FormatInt(date, 10),
		"limit":  strconv.Itoa(limit),
	}
	body, err := s.client.DoGETWithQuery(ctx, "get_historical_monitor_data/v1", params)
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
