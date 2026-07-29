package provider

import (
	"context"

	"github.com/gumieri/xyops-integrations/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type eventsDataSource struct {
	client *client.Client
}

func NewEventsDataSource() datasource.DataSource {
	return &eventsDataSource{}
}

func (d *eventsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *eventsDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "xyops_events"
}

func (d *eventsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List all xyops events",
		Attributes: map[string]schema.Attribute{
			"events": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of events",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
						},
						"title": schema.StringAttribute{
							Computed: true,
						},
						"plugin": schema.StringAttribute{
							Computed: true,
						},
						"targets": schema.ListAttribute{
							Computed:    true,
							ElementType: types.StringType,
						},
						"category": schema.StringAttribute{
							Computed: true,
						},
						"algo": schema.StringAttribute{
							Computed: true,
						},
						"enabled": schema.BoolAttribute{
							Computed: true,
						},
						"tags": schema.ListAttribute{
							Computed:    true,
							ElementType: types.StringType,
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

func (d *eventsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	events, err := d.client.Events.List(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error listing events", err.Error())
		return
	}

	elems := make([]attr.Value, 0, len(events))
	for _, e := range events {
		objVal, ds := types.ObjectValue(map[string]attr.Type{
			"id":       types.StringType,
			"title":    types.StringType,
			"plugin":   types.StringType,
			"targets":  types.ListType{ElemType: types.StringType},
			"category": types.StringType,
			"algo":     types.StringType,
			"enabled":  types.BoolType,
			"tags":     types.ListType{ElemType: types.StringType},
			"notes":    types.StringType,
		}, map[string]attr.Value{
			"id":       types.StringValue(e.ID),
			"title":    types.StringValue(e.Title),
			"plugin":   types.StringValue(e.Plugin),
			"targets":  stringSetFromSlice(e.Targets),
			"category": types.StringValue(e.Category),
			"algo":     types.StringValue(e.Algo),
			"enabled":  types.BoolValue(e.Enabled),
			"tags":     stringSetFromSlice(e.Tags),
			"notes":    types.StringValue(e.Notes),
		})
		resp.Diagnostics.Append(ds...)
		elems = append(elems, objVal)
	}

	var data eventsDataModel
	listVal, ds := types.ListValue(types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":       types.StringType,
			"title":    types.StringType,
			"plugin":   types.StringType,
			"targets":  types.ListType{ElemType: types.StringType},
			"category": types.StringType,
			"algo":     types.StringType,
			"enabled":  types.BoolType,
			"tags":     types.ListType{ElemType: types.StringType},
			"notes":    types.StringType,
		},
	}, elems)
	data.Events = listVal
	resp.Diagnostics.Append(ds...)

	diags := resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

type eventsDataModel struct {
	Events types.List `tfsdk:"events"`
}
