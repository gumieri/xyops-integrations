package provider

import (
	"context"

	"github.com/gumieri/xyops-integrations/sdk"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type alertsDataSource struct {
	client *sdk.Client
}

func NewAlertsDataSource() datasource.DataSource {
	return &alertsDataSource{}
}

func (d *alertsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *alertsDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "xyops_alerts"
}

func (d *alertsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List all xyops alerts",
		Attributes: map[string]schema.Attribute{
			"alerts": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of alerts",
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
						"expression": schema.StringAttribute{
							Computed: true,
						},
						"message": schema.StringAttribute{
							Computed: true,
						},
						"groups": schema.ListAttribute{
							Computed:    true,
							ElementType: types.StringType,
						},
						"monitor_id": schema.StringAttribute{
							Computed: true,
						},
						"max_count": schema.Int64Attribute{
							Computed: true,
						},
						"notes": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func (d *alertsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	alerts, err := d.client.Alerts.List(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error listing alerts", err.Error())
		return
	}

	elems := make([]attr.Value, 0, len(alerts))
	for _, a := range alerts {
		objVal, ds := types.ObjectValue(map[string]attr.Type{
			"id":         types.StringType,
			"title":      types.StringType,
			"enabled":    types.BoolType,
			"expression": types.StringType,
			"message":    types.StringType,
			"groups":     types.ListType{ElemType: types.StringType},
			"monitor_id": types.StringType,
			"max_count":  types.Int64Type,
			"notes":      types.StringType,
		}, map[string]attr.Value{
			"id":         types.StringValue(a.ID),
			"title":      types.StringValue(a.Title),
			"enabled":    types.BoolValue(a.Enabled),
			"expression": types.StringValue(a.Expression),
			"message":    types.StringValue(a.Message),
			"groups":     stringSetFromSlice(a.Groups),
			"monitor_id": types.StringValue(a.MonitorID),
			"max_count":  types.Int64Value(int64(a.MaxCount)),
			"notes":      types.StringValue(a.Notes),
		})
		resp.Diagnostics.Append(ds...)
		elems = append(elems, objVal)
	}

	var data alertsDataModel
	listVal, ds := types.ListValue(types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":         types.StringType,
			"title":      types.StringType,
			"enabled":    types.BoolType,
			"expression": types.StringType,
			"message":    types.StringType,
			"groups":     types.ListType{ElemType: types.StringType},
			"monitor_id": types.StringType,
			"max_count":  types.Int64Type,
			"notes":      types.StringType,
		},
	}, elems)
	data.Alerts = listVal
	resp.Diagnostics.Append(ds...)

	diags := resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

type alertsDataModel struct {
	Alerts types.List `tfsdk:"alerts"`
}
