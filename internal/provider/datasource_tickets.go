package provider

import (
	"context"

	"github.com/gumieri/xyops-integrations/sdk"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ticketsDataSource struct {
	client *sdk.Client
}

func NewTicketsDataSource() datasource.DataSource {
	return &ticketsDataSource{}
}

func (d *ticketsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ticketsDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "xyops_tickets"
}

func (d *ticketsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List all xyops tickets",
		Attributes: map[string]schema.Attribute{
			"query": schema.StringAttribute{
				Optional:    true,
				Description: "Search query (Unbase/PxQL syntax). Defaults to all tickets.",
			},
			"tickets": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of tickets",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
						},
						"num": schema.Int64Attribute{
							Computed: true,
						},
						"subject": schema.StringAttribute{
							Computed: true,
						},
						"body": schema.StringAttribute{
							Computed: true,
						},
						"type": schema.StringAttribute{
							Computed: true,
						},
						"status": schema.StringAttribute{
							Computed: true,
						},
						"category": schema.StringAttribute{
							Computed: true,
						},
						"server": schema.StringAttribute{
							Computed: true,
						},
						"assignees": schema.ListAttribute{
							Computed:    true,
							ElementType: types.StringType,
						},
						"cc": schema.ListAttribute{
							Computed:    true,
							ElementType: types.StringType,
						},
						"notify": schema.ListAttribute{
							Computed:    true,
							ElementType: types.StringType,
						},
						"due": schema.Int64Attribute{
							Computed: true,
						},
						"tags": schema.ListAttribute{
							Computed:    true,
							ElementType: types.StringType,
						},
					},
				},
			},
		},
	}
}

func (d *ticketsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config ticketsDataModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	opts := sdk.SearchTicketsOpts{
		Query: config.Query.ValueString(),
	}
	tickets, err := d.client.Tickets.Search(ctx, opts)
	if err != nil {
		resp.Diagnostics.AddError("Error listing tickets", err.Error())
		return
	}

	elems := make([]attr.Value, 0, len(tickets))
	for _, t := range tickets {
		objVal, ds := types.ObjectValue(map[string]attr.Type{
			"id":        types.StringType,
			"num":       types.Int64Type,
			"subject":   types.StringType,
			"body":      types.StringType,
			"type":      types.StringType,
			"status":    types.StringType,
			"category":  types.StringType,
			"server":    types.StringType,
			"assignees": types.ListType{ElemType: types.StringType},
			"cc":        types.ListType{ElemType: types.StringType},
			"notify":    types.ListType{ElemType: types.StringType},
			"due":       types.Int64Type,
			"tags":      types.ListType{ElemType: types.StringType},
		}, map[string]attr.Value{
			"id":        types.StringValue(t.ID),
			"num":       types.Int64Value(int64(t.Num)),
			"subject":   types.StringValue(t.Subject),
			"body":      types.StringValue(t.Body),
			"type":      types.StringValue(t.Type),
			"status":    types.StringValue(t.Status),
			"category":  types.StringValue(t.Category),
			"server":    types.StringValue(t.Server),
			"assignees": stringSetFromSlice(t.Assignees),
			"cc":        stringSetFromSlice(t.CC),
			"notify":    stringSetFromSlice(t.Notify),
			"due":       types.Int64Value(t.Due),
			"tags":      stringSetFromSlice(t.Tags),
		})
		resp.Diagnostics.Append(ds...)
		elems = append(elems, objVal)
	}

	var data ticketsDataModel
	listVal, ds := types.ListValue(types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":        types.StringType,
			"num":       types.Int64Type,
			"subject":   types.StringType,
			"body":      types.StringType,
			"type":      types.StringType,
			"status":    types.StringType,
			"category":  types.StringType,
			"server":    types.StringType,
			"assignees": types.ListType{ElemType: types.StringType},
			"cc":        types.ListType{ElemType: types.StringType},
			"notify":    types.ListType{ElemType: types.StringType},
			"due":       types.Int64Type,
			"tags":      types.ListType{ElemType: types.StringType},
		},
	}, elems)
	data.Tickets = listVal
	resp.Diagnostics.Append(ds...)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

type ticketsDataModel struct {
	Query   types.String `tfsdk:"query"`
	Tickets types.List   `tfsdk:"tickets"`
}
