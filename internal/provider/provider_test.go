package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestProvider_Resources(t *testing.T) {
	p := New("test")
	ctx := t.Context()

	resources := p.Resources(ctx)
	if resources == nil {
		t.Fatal("expected non-nil resources")
	}

	names := make(map[string]bool)
	for _, factory := range resources {
		r := factory()
		var resp resource.MetadataResponse
		r.Metadata(ctx, resource.MetadataRequest{}, &resp)
		names[resp.TypeName] = true
	}

	expected := []string{
		"xyops_tag", "xyops_group", "xyops_category", "xyops_monitor",
		"xyops_alert", "xyops_webhook", "xyops_role", "xyops_api_key",
		"xyops_plugin", "xyops_event",
	}
	for _, name := range expected {
		if !names[name] {
			t.Errorf("missing resource: %s", name)
		}
	}
	if len(names) != 10 {
		t.Errorf("expected 10 resources, got %d", len(names))
	}
}

func TestProvider_DataSources(t *testing.T) {
	p := New("test")
	ctx := t.Context()

	datasources := p.DataSources(ctx)
	if datasources == nil {
		t.Fatal("expected non-nil datasources")
	}

	names := make(map[string]bool)
	for _, factory := range datasources {
		d := factory()
		var resp datasource.MetadataResponse
		d.Metadata(ctx, datasource.MetadataRequest{}, &resp)
		names[resp.TypeName] = true
	}

	expected := []string{
		"xyops_tags", "xyops_groups", "xyops_categories", "xyops_plugins",
		"xyops_events", "xyops_monitors", "xyops_alerts", "xyops_roles",
		"xyops_api_keys", "xyops_webhooks", "xyops_servers", "xyops_config",
	}
	for _, name := range expected {
		if !names[name] {
			t.Errorf("missing data source: %s", name)
		}
	}
	if len(names) != 12 {
		t.Errorf("expected 12 datasources, got %d", len(names))
	}
}

func TestProvider_Metadata(t *testing.T) {
	p := New("test")
	var resp provider.MetadataResponse
	p.Metadata(t.Context(), provider.MetadataRequest{}, &resp)
	if resp.TypeName != "xyops" {
		t.Errorf("expected type name 'xyops', got %q", resp.TypeName)
	}
}
