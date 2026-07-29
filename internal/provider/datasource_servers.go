package provider

import (
	"context"

	"github.com/gumieri/xyops-integrations/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type serversDataSource struct {
	client *client.Client
}

func NewServersDataSource() datasource.DataSource {
	return &serversDataSource{}
}

func (d *serversDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *serversDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "xyops_servers"
}

func (d *serversDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List all xyops satellite servers",
		Attributes: map[string]schema.Attribute{
			"servers": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of servers",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"hostname": schema.StringAttribute{
							Computed: true,
						},
						"os": schema.StringAttribute{
							Computed: true,
						},
						"arch": schema.StringAttribute{
							Computed: true,
						},
						"status": schema.StringAttribute{
							Computed: true,
						},
						"groups": schema.ListAttribute{
							Computed:    true,
							ElementType: types.StringType,
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

func (d *serversDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	servers, err := d.client.Servers.List(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error listing servers", err.Error())
		return
	}

	elems := make([]attr.Value, 0, len(servers))
	for _, s := range servers {
		objVal, ds := types.ObjectValue(map[string]attr.Type{
			"hostname": types.StringType,
			"os":       types.StringType,
			"arch":     types.StringType,
			"status":   types.StringType,
			"groups":   types.ListType{ElemType: types.StringType},
			"tags":     types.ListType{ElemType: types.StringType},
		}, map[string]attr.Value{
			"hostname": types.StringValue(s.Hostname),
			"os":       types.StringValue(s.OS),
			"arch":     types.StringValue(s.Arch),
			"status":   types.StringValue(s.Status),
			"groups":   stringSetFromSlice(s.Groups),
			"tags":     stringSetFromSlice(s.Tags),
		})
		resp.Diagnostics.Append(ds...)
		elems = append(elems, objVal)
	}

	var data serversDataModel
	listVal, ds := types.ListValue(types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"hostname": types.StringType,
			"os":       types.StringType,
			"arch":     types.StringType,
			"status":   types.StringType,
			"groups":   types.ListType{ElemType: types.StringType},
			"tags":     types.ListType{ElemType: types.StringType},
		},
	}, elems)
	data.Servers = listVal
	resp.Diagnostics.Append(ds...)

	diags := resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

type serversDataModel struct {
	Servers types.List `tfsdk:"servers"`
}
