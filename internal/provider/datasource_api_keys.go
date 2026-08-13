package provider

import (
	"context"

	"github.com/gumieri/xyops-integrations/sdk"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type apiKeysDataSource struct {
	client *sdk.Client
}

func NewAPIKeysDataSource() datasource.DataSource {
	return &apiKeysDataSource{}
}

func (d *apiKeysDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *apiKeysDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "xyops_api_keys"
}

func (d *apiKeysDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List all xyops API keys",
		Attributes: map[string]schema.Attribute{
			"api_keys": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of API keys",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
						},
						"title": schema.StringAttribute{
							Computed: true,
						},
						"description": schema.StringAttribute{
							Computed: true,
						},
						"active": schema.BoolAttribute{
							Computed: true,
						},
						"privileges": schema.MapAttribute{
							Computed:    true,
							ElementType: types.BoolType,
						},
						"roles": schema.ListAttribute{
							Computed:    true,
							ElementType: types.StringType,
						},
						"expires": schema.Int64Attribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func (d *apiKeysDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	keys, err := d.client.APIKeys.List(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error listing API keys", err.Error())
		return
	}

	elems := make([]attr.Value, 0, len(keys))
	for _, k := range keys {
		objVal, ds := types.ObjectValue(map[string]attr.Type{
			"id":          types.StringType,
			"title":       types.StringType,
			"description": types.StringType,
			"active":      types.BoolType,
			"privileges":  types.MapType{ElemType: types.BoolType},
			"roles":       types.ListType{ElemType: types.StringType},
			"expires":     types.Int64Type,
		}, map[string]attr.Value{
			"id":          types.StringValue(k.ID),
			"title":       types.StringValue(k.Title),
			"description": types.StringValue(k.Description),
			"active":      types.BoolValue(k.Active == 1),
			"privileges":  boolMapToFrameworkMap(k.Privileges),
			"roles":       stringSetFromSlice(k.Roles),
			"expires":     types.Int64Value(k.Expires),
		})
		resp.Diagnostics.Append(ds...)
		elems = append(elems, objVal)
	}

	var data apiKeysDataModel
	listVal, ds := types.ListValue(types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":          types.StringType,
			"title":       types.StringType,
			"description": types.StringType,
			"active":      types.BoolType,
			"privileges":  types.MapType{ElemType: types.BoolType},
			"roles":       types.ListType{ElemType: types.StringType},
			"expires":     types.Int64Type,
		},
	}, elems)
	data.APIKeys = listVal
	resp.Diagnostics.Append(ds...)

	diags := resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

type apiKeysDataModel struct {
	APIKeys types.List `tfsdk:"api_keys"`
}
