package provider

import (
	"context"

	"github.com/gumieri/xyops-integrations/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type monitorsDataSource struct {
	client *client.Client
}

func NewMonitorsDataSource() datasource.DataSource {
	return &monitorsDataSource{}
}

func (d *monitorsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			"Expected *client.Client, got: %T. Please report this bug to the provider developers.",
		)
		return
	}

	d.client = c
}

func (d *monitorsDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "xyops_monitors"
}

func (d *monitorsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List all xyops monitors",
		Attributes: map[string]schema.Attribute{
			"monitors": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of monitors",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
						},
						"title": schema.StringAttribute{
							Computed: true,
						},
						"source": schema.StringAttribute{
							Computed: true,
						},
						"data_type": schema.StringAttribute{
							Computed: true,
						},
						"display": schema.BoolAttribute{
							Computed: true,
						},
						"groups": schema.ListAttribute{
							Computed:    true,
							ElementType: types.StringType,
						},
						"suffix": schema.StringAttribute{
							Computed: true,
						},
						"min_vert_scale": schema.Int64Attribute{
							Computed: true,
						},
						"data_match": schema.StringAttribute{
							Computed: true,
						},
						"delta": schema.BoolAttribute{
							Computed: true,
						},
						"divide_time": schema.BoolAttribute{
							Computed: true,
						},
						"zero_min": schema.Int64Attribute{
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

func (d *monitorsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	monitors, err := d.client.Monitors.List(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error listing monitors", err.Error())
		return
	}

	elems := make([]attr.Value, 0, len(monitors))
	for _, m := range monitors {
		objVal, ds := types.ObjectValue(map[string]attr.Type{
			"id":             types.StringType,
			"title":          types.StringType,
			"source":         types.StringType,
			"data_type":      types.StringType,
			"display":        types.BoolType,
			"groups":         types.ListType{ElemType: types.StringType},
			"suffix":         types.StringType,
			"min_vert_scale": types.Int64Type,
			"data_match":     types.StringType,
			"delta":          types.BoolType,
			"divide_time":    types.BoolType,
			"zero_min":       types.Int64Type,
			"notes":          types.StringType,
		}, map[string]attr.Value{
			"id":             types.StringValue(m.ID),
			"title":          types.StringValue(m.Title),
			"source":         types.StringValue(m.Source),
			"data_type":      types.StringValue(m.DataType),
			"display":        types.BoolValue(m.Display),
			"groups":         stringSetFromSlice(m.Groups),
			"suffix":         types.StringValue(m.Suffix),
			"min_vert_scale": types.Int64Value(int64(m.MinVertScale)),
			"data_match":     types.StringValue(m.DataMatch),
			"delta":          types.BoolValue(m.Delta),
			"divide_time":    types.BoolValue(m.DivideTime),
			"zero_min":       types.Int64Value(int64(m.ZeroMin)),
			"notes":          types.StringValue(m.Notes),
		})
		resp.Diagnostics.Append(ds...)
		elems = append(elems, objVal)
	}

	var data monitorsDataModel
	listVal, ds := types.ListValue(types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":             types.StringType,
			"title":          types.StringType,
			"source":         types.StringType,
			"data_type":      types.StringType,
			"display":        types.BoolType,
			"groups":         types.ListType{ElemType: types.StringType},
			"suffix":         types.StringType,
			"min_vert_scale": types.Int64Type,
			"data_match":     types.StringType,
			"delta":          types.BoolType,
			"divide_time":    types.BoolType,
			"zero_min":       types.Int64Type,
			"notes":          types.StringType,
		},
	}, elems)
	data.Monitors = listVal
	resp.Diagnostics.Append(ds...)

	diags := resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

type monitorsDataModel struct {
	Monitors types.List `tfsdk:"monitors"`
}
