package provider

import (
	"context"

	"github.com/gumieri/xyops-integrations/sdk"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type webHooksDataSource struct {
	client *sdk.Client
}

func NewWebHooksDataSource() datasource.DataSource {
	return &webHooksDataSource{}
}

func (d *webHooksDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	c, ok := req.ProviderData.(*sdk.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			"Expected *sdk.Client, got: %T. Please report this bug to the provider developers.",
		)
		return
	}

	d.client = c
}

func (d *webHooksDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "xyops_webhooks"
}

func (d *webHooksDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List all xyops webhooks",
		Attributes: map[string]schema.Attribute{
			"webhooks": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of webhooks",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
						},
						"title": schema.StringAttribute{
							Computed: true,
						},
						"enabled": schema.BoolAttribute{
							Computed: true,
						},
						"url": schema.StringAttribute{
							Computed: true,
						},
						"method": schema.StringAttribute{
							Computed: true,
						},
						"body": schema.StringAttribute{
							Computed: true,
						},
						"timeout": schema.Int64Attribute{
							Computed: true,
						},
						"retries": schema.Int64Attribute{
							Computed: true,
						},
						"follow": schema.BoolAttribute{
							Computed: true,
						},
						"ssl_cert_bypass": schema.BoolAttribute{
							Computed: true,
						},
						"max_per_day": schema.Int64Attribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func (d *webHooksDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	hooks, err := d.client.WebHooks.List(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error listing webhooks", err.Error())
		return
	}

	elems := make([]attr.Value, 0, len(hooks))
	for _, h := range hooks {
		objVal, ds := types.ObjectValue(map[string]attr.Type{
			"id":              types.StringType,
			"title":           types.StringType,
			"enabled":         types.BoolType,
			"url":             types.StringType,
			"method":          types.StringType,
			"body":            types.StringType,
			"timeout":         types.Int64Type,
			"retries":         types.Int64Type,
			"follow":          types.BoolType,
			"ssl_cert_bypass": types.BoolType,
			"max_per_day":     types.Int64Type,
		}, map[string]attr.Value{
			"id":              types.StringValue(h.ID),
			"title":           types.StringValue(h.Title),
			"enabled":         types.BoolValue(h.Enabled),
			"url":             types.StringValue(h.URL),
			"method":          types.StringValue(h.Method),
			"body":            types.StringValue(h.Body),
			"timeout":         types.Int64Value(int64(h.Timeout)),
			"retries":         types.Int64Value(int64(h.Retries)),
			"follow":          types.BoolValue(h.Follow),
			"ssl_cert_bypass": types.BoolValue(h.SSLCertBypass),
			"max_per_day":     types.Int64Value(int64(h.MaxPerDay)),
		})
		resp.Diagnostics.Append(ds...)
		elems = append(elems, objVal)
	}

	var data webHooksDataModel
	listVal, ds := types.ListValue(types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":              types.StringType,
			"title":           types.StringType,
			"enabled":         types.BoolType,
			"url":             types.StringType,
			"method":          types.StringType,
			"body":            types.StringType,
			"timeout":         types.Int64Type,
			"retries":         types.Int64Type,
			"follow":          types.BoolType,
			"ssl_cert_bypass": types.BoolType,
			"max_per_day":     types.Int64Type,
		},
	}, elems)
	data.WebHooks = listVal
	resp.Diagnostics.Append(ds...)

	diags := resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

type webHooksDataModel struct {
	WebHooks types.List `tfsdk:"webhooks"`
}
