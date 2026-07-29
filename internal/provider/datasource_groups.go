package provider

import (
	"context"

	"github.com/gumieri/xyops-integrations/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type groupsDataSource struct {
	client *client.Client
}

func NewGroupsDataSource() datasource.DataSource {
	return &groupsDataSource{}
}

func (d *groupsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *groupsDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "xyops_groups"
}

func (d *groupsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List all xyops groups",
		Attributes: map[string]schema.Attribute{
			"groups": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of groups",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
						},
						"title": schema.StringAttribute{
							Computed: true,
						},
						"hostname_match": schema.StringAttribute{
							Computed: true,
						},
						"sort_order": schema.Int64Attribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func (d *groupsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	groups, err := d.client.Groups.List(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error listing groups", err.Error())
		return
	}

	elems := make([]attr.Value, 0, len(groups))
	for _, g := range groups {
		objVal, ds := types.ObjectValue(map[string]attr.Type{
			"id":             types.StringType,
			"title":          types.StringType,
			"hostname_match": types.StringType,
			"sort_order":     types.Int64Type,
		}, map[string]attr.Value{
			"id":             types.StringValue(g.ID),
			"title":          types.StringValue(g.Title),
			"hostname_match": types.StringValue(g.HostnameMatch),
			"sort_order":     types.Int64Value(int64(g.SortOrder)),
		})
		resp.Diagnostics.Append(ds...)
		elems = append(elems, objVal)
	}

	var data groupsDataModel
	listVal, ds := types.ListValue(types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":             types.StringType,
			"title":          types.StringType,
			"hostname_match": types.StringType,
			"sort_order":     types.Int64Type,
		},
	}, elems)
	data.Groups = listVal
	resp.Diagnostics.Append(ds...)

	diags := resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

type groupsDataModel struct {
	Groups types.List `tfsdk:"groups"`
}
