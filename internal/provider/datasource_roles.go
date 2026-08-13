package provider

import (
	"context"

	"github.com/gumieri/xyops-integrations/sdk"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type rolesDataSource struct {
	client *sdk.Client
}

func NewRolesDataSource() datasource.DataSource {
	return &rolesDataSource{}
}

func (d *rolesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *rolesDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "xyops_roles"
}

func (d *rolesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List all xyops roles",
		Attributes: map[string]schema.Attribute{
			"roles": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of roles",
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
						"privileges": schema.MapAttribute{
							Computed:    true,
							ElementType: types.BoolType,
						},
						"categories": schema.ListAttribute{
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

func (d *rolesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	roles, err := d.client.Roles.List(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error listing roles", err.Error())
		return
	}

	elems := make([]attr.Value, 0, len(roles))
	for _, r := range roles {
		objVal, ds := types.ObjectValue(map[string]attr.Type{
			"id":         types.StringType,
			"title":      types.StringType,
			"enabled":    types.BoolType,
			"privileges": types.MapType{ElemType: types.BoolType},
			"categories": types.ListType{ElemType: types.StringType},
			"notes":      types.StringType,
		}, map[string]attr.Value{
			"id":         types.StringValue(r.ID),
			"title":      types.StringValue(r.Title),
			"enabled":    types.BoolValue(r.Enabled),
			"privileges": intMapToBoolMap(r.Privileges),
			"categories": stringSetFromSlice(r.Categories),
			"notes":      types.StringValue(r.Notes),
		})
		resp.Diagnostics.Append(ds...)
		elems = append(elems, objVal)
	}

	var data rolesDataModel
	listVal, ds := types.ListValue(types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":         types.StringType,
			"title":      types.StringType,
			"enabled":    types.BoolType,
			"privileges": types.MapType{ElemType: types.BoolType},
			"categories": types.ListType{ElemType: types.StringType},
			"notes":      types.StringType,
		},
	}, elems)
	data.Roles = listVal
	resp.Diagnostics.Append(ds...)

	diags := resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

type rolesDataModel struct {
	Roles types.List `tfsdk:"roles"`
}
