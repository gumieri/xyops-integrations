package provider

import (
	"context"
	"os"

	"github.com/gumieri/xyops-integrations/sdk"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func New(version string) provider.Provider {
	return &xyopsProvider{}
}

type xyopsProvider struct{}

func (p *xyopsProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "xyops"
}

func (p *xyopsProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"base_url": schema.StringAttribute{
				Required:    true,
				Description: "xyops conductor base URL (e.g., https://xyops.example.com)",
			},
			"api_key": schema.StringAttribute{
				Required:    true,
				Sensitive:   true,
				Description: "xyops API key for authentication",
			},
		},
	}
}

func (p *xyopsProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var baseURL, apiKey types.String

	diags := req.Config.GetAttribute(ctx, path.Root("base_url"), &baseURL)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	diags = req.Config.GetAttribute(ctx, path.Root("api_key"), &apiKey)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	if baseURL.IsNull() || baseURL.IsUnknown() {
		baseURL = types.StringValue(os.Getenv("XYOPS_BASE_URL"))
	}

	if apiKey.IsNull() || apiKey.IsUnknown() {
		apiKey = types.StringValue(os.Getenv("XYOPS_API_KEY"))
	}

	xyopsClient, err := sdk.New(baseURL.ValueString(), apiKey.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error creating xyops client", err.Error())
		return
	}

	resp.DataSourceData = xyopsClient
	resp.ResourceData = xyopsClient
}

func (p *xyopsProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewTagResource,
		NewGroupResource,
		NewCategoryResource,
		NewMonitorResource,
		NewAlertResource,
		NewWebHookResource,
		NewRoleResource,
		NewAPIKeyResource,
		NewPluginResource,
		NewEventResource,
		NewChannelResource,
		NewSecretResource,
		NewBucketResource,
		NewTicketResource,
	}
}

func (p *xyopsProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewTagsDataSource,
		NewGroupsDataSource,
		NewCategoriesDataSource,
		NewPluginsDataSource,
		NewEventsDataSource,
		NewMonitorsDataSource,
		NewAlertsDataSource,
		NewRolesDataSource,
		NewAPIKeysDataSource,
		NewWebHooksDataSource,
		NewServersDataSource,
		NewConfigDataSource,
		NewChannelsDataSource,
		NewSecretsDataSource,
		NewBucketsDataSource,
		NewTicketsDataSource,
	}
}
