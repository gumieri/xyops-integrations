package provider

import (
	"context"

	"github.com/gumieri/xyops-integrations/sdk"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type secretsDataSource struct {
	client *sdk.Client
}

func NewSecretsDataSource() datasource.DataSource {
	return &secretsDataSource{}
}

func (d *secretsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *secretsDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "xyops_secrets"
}

func (d *secretsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List all xyops secrets",
		Attributes: map[string]schema.Attribute{
			"secrets": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of secrets",
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
						"names": schema.ListAttribute{
							Computed:    true,
							ElementType: types.StringType,
						},
						"events": schema.ListAttribute{
							Computed:    true,
							ElementType: types.StringType,
						},
						"categories": schema.ListAttribute{
							Computed:    true,
							ElementType: types.StringType,
						},
						"plugins": schema.ListAttribute{
							Computed:    true,
							ElementType: types.StringType,
						},
						"web_hooks": schema.ListAttribute{
							Computed:    true,
							ElementType: types.StringType,
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

func (d *secretsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	secrets, err := d.client.Secrets.List(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error listing secrets", err.Error())
		return
	}

	elems := make([]attr.Value, 0, len(secrets))
	for _, s := range secrets {
		objVal, ds := types.ObjectValue(map[string]attr.Type{
			"id":         types.StringType,
			"title":      types.StringType,
			"enabled":    types.BoolType,
			"icon":       types.StringType,
			"notes":      types.StringType,
			"names":      types.ListType{ElemType: types.StringType},
			"events":     types.ListType{ElemType: types.StringType},
			"categories": types.ListType{ElemType: types.StringType},
			"plugins":    types.ListType{ElemType: types.StringType},
			"web_hooks":  types.ListType{ElemType: types.StringType},
			"username":   types.StringType,
			"created":    types.Int64Type,
			"modified":   types.Int64Type,
			"revision":   types.Int64Type,
		}, map[string]attr.Value{
			"id":         types.StringValue(s.ID),
			"title":      types.StringValue(s.Title),
			"enabled":    types.BoolValue(s.Enabled),
			"icon":       types.StringValue(s.Icon),
			"notes":      types.StringValue(s.Notes),
			"names":      stringSetFromSlice(s.Names),
			"events":     stringSetFromSlice(s.Events),
			"categories": stringSetFromSlice(s.Categories),
			"plugins":    stringSetFromSlice(s.Plugins),
			"web_hooks":  stringSetFromSlice(s.WebHooks),
			"username":   types.StringValue(s.Username),
			"created":    types.Int64Value(s.Created),
			"modified":   types.Int64Value(s.Modified),
			"revision":   types.Int64Value(int64(s.Revision)),
		})
		resp.Diagnostics.Append(ds...)
		elems = append(elems, objVal)
	}

	var data secretsDataModel
	listVal, ds := types.ListValue(types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":         types.StringType,
			"title":      types.StringType,
			"enabled":    types.BoolType,
			"icon":       types.StringType,
			"notes":      types.StringType,
			"names":      types.ListType{ElemType: types.StringType},
			"events":     types.ListType{ElemType: types.StringType},
			"categories": types.ListType{ElemType: types.StringType},
			"plugins":    types.ListType{ElemType: types.StringType},
			"web_hooks":  types.ListType{ElemType: types.StringType},
			"username":   types.StringType,
			"created":    types.Int64Type,
			"modified":   types.Int64Type,
			"revision":   types.Int64Type,
		},
	}, elems)
	data.Secrets = listVal
	resp.Diagnostics.Append(ds...)

	diags := resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

type secretsDataModel struct {
	Secrets types.List `tfsdk:"secrets"`
}
