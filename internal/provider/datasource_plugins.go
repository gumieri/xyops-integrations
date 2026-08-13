package provider

import (
	"context"

	"github.com/gumieri/xyops-integrations/sdk"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type pluginsDataSource struct {
	client *sdk.Client
}

func NewPluginsDataSource() datasource.DataSource {
	return &pluginsDataSource{}
}

func (d *pluginsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *pluginsDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "xyops_plugins"
}

func (d *pluginsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List all xyops plugins",
		Attributes: map[string]schema.Attribute{
			"plugins": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of plugins",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
						},
						"title": schema.StringAttribute{
							Computed: true,
						},
						"type": schema.StringAttribute{
							Computed: true,
						},
						"command": schema.StringAttribute{
							Computed: true,
						},
						"script": schema.StringAttribute{
							Computed: true,
						},
						"enabled": schema.BoolAttribute{
							Computed: true,
						},
						"uid": schema.StringAttribute{
							Computed: true,
						},
						"gid": schema.StringAttribute{
							Computed: true,
						},
						"groups": schema.ListAttribute{
							Computed:    true,
							ElementType: types.StringType,
						},
						"format": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func (d *pluginsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	plugins, err := d.client.Plugins.List(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error listing plugins", err.Error())
		return
	}

	elems := make([]attr.Value, 0, len(plugins))
	for _, p := range plugins {
		objVal, ds := types.ObjectValue(map[string]attr.Type{
			"id":      types.StringType,
			"title":   types.StringType,
			"type":    types.StringType,
			"command": types.StringType,
			"script":  types.StringType,
			"enabled": types.BoolType,
			"uid":     types.StringType,
			"gid":     types.StringType,
			"groups":  types.ListType{ElemType: types.StringType},
			"format":  types.StringType,
		}, map[string]attr.Value{
			"id":      types.StringValue(p.ID),
			"title":   types.StringValue(p.Title),
			"type":    types.StringValue(p.Type),
			"command": types.StringValue(p.Command),
			"script":  types.StringValue(p.Script),
			"enabled": types.BoolValue(intToBool(int(p.Enabled))),
			"uid":     types.StringValue(p.UID),
			"gid":     types.StringValue(p.GID),
			"groups":  stringSetFromSlice(p.Groups),
			"format":  types.StringValue(p.Format),
		})
		resp.Diagnostics.Append(ds...)
		elems = append(elems, objVal)
	}

	var data pluginsDataModel
	listVal, ds := types.ListValue(types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":      types.StringType,
			"title":   types.StringType,
			"type":    types.StringType,
			"command": types.StringType,
			"script":  types.StringType,
			"enabled": types.BoolType,
			"uid":     types.StringType,
			"gid":     types.StringType,
			"groups":  types.ListType{ElemType: types.StringType},
			"format":  types.StringType,
		},
	}, elems)
	data.Plugins = listVal
	resp.Diagnostics.Append(ds...)

	diags := resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

type pluginsDataModel struct {
	Plugins types.List `tfsdk:"plugins"`
}
