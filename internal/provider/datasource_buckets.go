package provider

import (
	"context"

	"github.com/gumieri/xyops-integrations/sdk"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type bucketsDataSource struct {
	client *sdk.Client
}

func NewBucketsDataSource() datasource.DataSource {
	return &bucketsDataSource{}
}

func (d *bucketsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *bucketsDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "xyops_buckets"
}

func (d *bucketsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List all xyops buckets",
		Attributes: map[string]schema.Attribute{
			"buckets": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of buckets",
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

func (d *bucketsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	buckets, err := d.client.Buckets.List(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error listing buckets", err.Error())
		return
	}

	elems := make([]attr.Value, 0, len(buckets))
	for _, b := range buckets {
		objVal, ds := types.ObjectValue(map[string]attr.Type{
			"id":       types.StringType,
			"title":    types.StringType,
			"enabled":  types.BoolType,
			"icon":     types.StringType,
			"notes":    types.StringType,
			"username": types.StringType,
			"created":  types.Int64Type,
			"modified": types.Int64Type,
			"revision": types.Int64Type,
		}, map[string]attr.Value{
			"id":       types.StringValue(b.ID),
			"title":    types.StringValue(b.Title),
			"enabled":  types.BoolValue(b.Enabled),
			"icon":     types.StringValue(b.Icon),
			"notes":    types.StringValue(b.Notes),
			"username": types.StringValue(b.Username),
			"created":  types.Int64Value(b.Created),
			"modified": types.Int64Value(b.Modified),
			"revision": types.Int64Value(int64(b.Revision)),
		})
		resp.Diagnostics.Append(ds...)
		elems = append(elems, objVal)
	}

	var data bucketsDataModel
	listVal, ds := types.ListValue(types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":       types.StringType,
			"title":    types.StringType,
			"enabled":  types.BoolType,
			"icon":     types.StringType,
			"notes":    types.StringType,
			"username": types.StringType,
			"created":  types.Int64Type,
			"modified": types.Int64Type,
			"revision": types.Int64Type,
		},
	}, elems)
	data.Buckets = listVal
	resp.Diagnostics.Append(ds...)

	diags := resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

type bucketsDataModel struct {
	Buckets types.List `tfsdk:"buckets"`
}
