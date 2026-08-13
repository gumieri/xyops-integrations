package provider

import (
	"context"

	"github.com/gumieri/xyops-integrations/sdk"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type categoriesDataSource struct {
	client *sdk.Client
}

func NewCategoriesDataSource() datasource.DataSource {
	return &categoriesDataSource{}
}

func (d *categoriesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *categoriesDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "xyops_categories"
}

func (d *categoriesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List all xyops categories",
		Attributes: map[string]schema.Attribute{
			"categories": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of categories",
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
					},
				},
			},
		},
	}
}

func (d *categoriesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	categories, err := d.client.Categories.List(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error listing categories", err.Error())
		return
	}

	elems := make([]attr.Value, 0, len(categories))
	for _, c := range categories {
		objVal, ds := types.ObjectValue(map[string]attr.Type{
			"id":      types.StringType,
			"title":   types.StringType,
			"enabled": types.BoolType,
		}, map[string]attr.Value{
			"id":      types.StringValue(c.ID),
			"title":   types.StringValue(c.Title),
			"enabled": types.BoolValue(c.Enabled),
		})
		resp.Diagnostics.Append(ds...)
		elems = append(elems, objVal)
	}

	var data categoriesDataModel
	listVal, ds := types.ListValue(types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":      types.StringType,
			"title":   types.StringType,
			"enabled": types.BoolType,
		},
	}, elems)
	data.Categories = listVal
	resp.Diagnostics.Append(ds...)

	diags := resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

type categoriesDataModel struct {
	Categories types.List `tfsdk:"categories"`
}
