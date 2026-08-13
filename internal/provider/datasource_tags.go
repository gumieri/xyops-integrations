package provider

import (
	"context"

	"github.com/gumieri/xyops-integrations/sdk"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type tagsDataSource struct {
	client *sdk.Client
}

func NewTagsDataSource() datasource.DataSource {
	return &tagsDataSource{}
}

func (d *tagsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *tagsDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "xyops_tags"
}

func (d *tagsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List all xyops tags",
		Attributes: map[string]schema.Attribute{
			"tags": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of tags",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
						},
						"title": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func (d *tagsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tags, err := d.client.Tags.List(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error listing tags", err.Error())
		return
	}

	elems := make([]attr.Value, 0, len(tags))
	for _, tag := range tags {
		objVal, ds := types.ObjectValue(map[string]attr.Type{
			"id":    types.StringType,
			"title": types.StringType,
		}, map[string]attr.Value{
			"id":    types.StringValue(tag.ID),
			"title": types.StringValue(tag.Title),
		})
		resp.Diagnostics.Append(ds...)
		elems = append(elems, objVal)
	}

	var data tagsDataModel
	listVal, ds := types.ListValue(types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":    types.StringType,
			"title": types.StringType,
		},
	}, elems)
	data.Tags = listVal
	resp.Diagnostics.Append(ds...)

	diags := resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

type tagsDataModel struct {
	Tags types.List `tfsdk:"tags"`
}
