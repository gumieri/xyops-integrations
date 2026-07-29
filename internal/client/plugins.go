package client

import "context"

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
	body, err := s.client.DoPOST(ctx, "update_plugin/v1", plugin)
	if err != nil {
		return nil, err
	}

	if err := checkSuccessCode(body); err != nil {
		return nil, err
	}

	return s.Get(ctx, plugin.ID)
}

func (s *PluginsService) Delete(ctx context.Context, id string) error {
	_, err := s.client.DoPOST(ctx, "delete_plugin/v1", map[string]string{"id": id})
	return err
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

	// Inject workflow markers for workflow-type events.
	if event.Workflow != nil {
		event.Plugin = "_workflow"
		event.Type = "workflow"
	}

	// Inject default limits when empty.
	if len(event.Limits) == 0 {
		event.Limits = DefaultLimits()
	}

	body, err := s.client.DoPOST(ctx, "create_event/v1", event)
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
	body, err := s.client.DoPOST(ctx, "update_event/v1", event)
	if err != nil {
		return nil, err
	}

	if err := checkSuccessCode(body); err != nil {
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
	body, err := s.client.DoPOST(ctx, "update_monitor/v1", monitor)
	if err != nil {
		return nil, err
	}

	if err := checkSuccessCode(body); err != nil {
		return nil, err
	}

	return s.Get(ctx, monitor.ID)
}

func (s *MonitorsService) Delete(ctx context.Context, id string) error {
	_, err := s.client.DoPOST(ctx, "delete_monitor/v1", map[string]string{"id": id})
	return err
}
