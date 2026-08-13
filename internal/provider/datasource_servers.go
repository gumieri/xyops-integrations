package provider

import (
	"context"

	"github.com/gumieri/xyops-integrations/sdk"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type serversDataSource struct {
	client *sdk.Client
}

func NewServersDataSource() datasource.DataSource {
	return &serversDataSource{}
}

func (d *serversDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
						"id": schema.StringAttribute{
							Computed: true,
						},
						"hostname": schema.StringAttribute{
							Computed: true,
						},
						"os_platform": schema.StringAttribute{
							Computed: true,
						},
						"os_distro": schema.StringAttribute{
							Computed: true,
						},
						"os_release": schema.StringAttribute{
							Computed: true,
						},
						"os_arch": schema.StringAttribute{
							Computed: true,
						},
						"cpu_cores": schema.Int64Attribute{
							Computed: true,
						},
						"cpu_combo": schema.StringAttribute{
							Computed: true,
						},
						"enabled": schema.BoolAttribute{
							Computed: true,
						},
						"groups": schema.ListAttribute{
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
			"id":          types.StringType,
			"hostname":    types.StringType,
			"os_platform": types.StringType,
			"os_distro":   types.StringType,
			"os_release":  types.StringType,
			"os_arch":     types.StringType,
			"cpu_cores":   types.Int64Type,
			"cpu_combo":   types.StringType,
			"enabled":     types.BoolType,
			"groups":      types.ListType{ElemType: types.StringType},
		}, map[string]attr.Value{
			"id":          types.StringValue(s.ID),
			"hostname":    types.StringValue(s.Hostname),
			"os_platform": types.StringValue(s.Info.OS.Platform),
			"os_distro":   types.StringValue(s.Info.OS.Distro),
			"os_release":  types.StringValue(s.Info.OS.Release),
			"os_arch":     types.StringValue(s.Info.OS.Arch),
			"cpu_cores":   types.Int64Value(int64(s.Info.CPU.Cores)),
			"cpu_combo":   types.StringValue(s.Info.CPU.Combo),
			"enabled":     types.BoolValue(s.Enabled),
			"groups":      stringSetFromSlice(s.Groups),
		})
		resp.Diagnostics.Append(ds...)
		elems = append(elems, objVal)
	}

	var data serversDataModel
	listVal, ds := types.ListValue(types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":          types.StringType,
			"hostname":    types.StringType,
			"os_platform": types.StringType,
			"os_distro":   types.StringType,
			"os_release":  types.StringType,
			"os_arch":     types.StringType,
			"cpu_cores":   types.Int64Type,
			"cpu_combo":   types.StringType,
			"enabled":     types.BoolType,
			"groups":      types.ListType{ElemType: types.StringType},
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
