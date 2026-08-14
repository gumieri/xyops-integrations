package provider

import (
	"context"

	"github.com/gumieri/xyops-integrations/sdk"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type channelsDataSource struct {
	client *sdk.Client
}

func NewChannelsDataSource() datasource.DataSource {
	return &channelsDataSource{}
}

func (d *channelsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *channelsDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "xyops_channels"
}

func (d *channelsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List all xyops channels",
		Attributes: map[string]schema.Attribute{
			"channels": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of channels",
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
						"icon": schema.StringAttribute{
							Computed: true,
						},
						"users": schema.ListAttribute{
							Computed:    true,
							ElementType: types.StringType,
						},
						"email": schema.StringAttribute{
							Computed: true,
						},
						"web_hook": schema.StringAttribute{
							Computed: true,
						},
						"run_event": schema.StringAttribute{
							Computed: true,
						},
						"sound": schema.StringAttribute{
							Computed: true,
						},
						"max_per_day": schema.Int64Attribute{
							Computed: true,
						},
						"notes": schema.StringAttribute{
							Computed: true,
						},
						"username": schema.StringAttribute{
							Computed: true,
						},
						"created": schema.Int64Attribute{
							Computed: true,
						},
						"modified": schema.Int64Attribute{
							Computed: true,
						},
						"revision": schema.Int64Attribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func (d *channelsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	channels, err := d.client.Channels.List(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error listing channels", err.Error())
		return
	}

	elems := make([]attr.Value, 0, len(channels))
	for _, ch := range channels {
		objVal, ds := types.ObjectValue(map[string]attr.Type{
			"id":          types.StringType,
			"title":       types.StringType,
			"enabled":     types.BoolType,
			"icon":        types.StringType,
			"users":       types.ListType{ElemType: types.StringType},
			"email":       types.StringType,
			"web_hook":    types.StringType,
			"run_event":   types.StringType,
			"sound":       types.StringType,
			"max_per_day": types.Int64Type,
			"notes":       types.StringType,
			"username":    types.StringType,
			"created":     types.Int64Type,
			"modified":    types.Int64Type,
			"revision":    types.Int64Type,
		}, map[string]attr.Value{
			"id":          types.StringValue(ch.ID),
			"title":       types.StringValue(ch.Title),
			"enabled":     types.BoolValue(ch.Enabled),
			"icon":        types.StringValue(ch.Icon),
			"users":       stringSetFromSlice(ch.Users),
			"email":       types.StringValue(ch.Email),
			"web_hook":    types.StringValue(ch.WebHook),
			"run_event":   types.StringValue(ch.RunEvent),
			"sound":       types.StringValue(ch.Sound),
			"max_per_day": types.Int64Value(int64(ch.MaxPerDay)),
			"notes":       types.StringValue(ch.Notes),
			"username":    types.StringValue(ch.Username),
			"created":     types.Int64Value(ch.Created),
			"modified":    types.Int64Value(ch.Modified),
			"revision":    types.Int64Value(int64(ch.Revision)),
		})
		resp.Diagnostics.Append(ds...)
		elems = append(elems, objVal)
	}

	var data channelsDataModel
	listVal, ds := types.ListValue(types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":          types.StringType,
			"title":       types.StringType,
			"enabled":     types.BoolType,
			"icon":        types.StringType,
			"users":       types.ListType{ElemType: types.StringType},
			"email":       types.StringType,
			"web_hook":    types.StringType,
			"run_event":   types.StringType,
			"sound":       types.StringType,
			"max_per_day": types.Int64Type,
			"notes":       types.StringType,
			"username":    types.StringType,
			"created":     types.Int64Type,
			"modified":    types.Int64Type,
			"revision":    types.Int64Type,
		},
	}, elems)
	data.Channels = listVal
	resp.Diagnostics.Append(ds...)

	diags := resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

type channelsDataModel struct {
	Channels types.List `tfsdk:"channels"`
}
