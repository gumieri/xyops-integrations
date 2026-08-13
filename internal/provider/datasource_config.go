package provider

import (
	"context"
	"fmt"
	"math/big"

	"github.com/gumieri/xyops-integrations/sdk"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type configDataSource struct {
	client *sdk.Client
}

func NewConfigDataSource() datasource.DataSource {
	return &configDataSource{}
}

func (d *configDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *configDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "xyops_config"
}

func (d *configDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "xyops configuration",
		Attributes: map[string]schema.Attribute{
			"config": schema.MapAttribute{
				Computed:    true,
				Description: "xyops configuration as a dynamic map",
				ElementType: types.DynamicType,
			},
		},
	}
}

func (d *configDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	cfg, err := d.client.Config.Get(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error reading config", err.Error())
		return
	}

	var data configDataModel

	elems := make(map[string]attr.Value)
	for k, v := range cfg {
		switch val := v.(type) {
		case string:
			elems[k] = types.StringValue(val)
		case bool:
			elems[k] = types.BoolValue(val)
		case float64:
			elems[k] = types.NumberValue(new(big.Float).SetFloat64(val))
		case int:
			elems[k] = types.Int64Value(int64(val))
		default:
			elems[k] = types.DynamicValue(types.StringValue(""))
			resp.Diagnostics.AddWarning(
				"Unsupported config value type",
				"Config key "+k+" has type "+fmt.Sprintf("%T", v)+" which is not supported; set as empty string.",
			)
		}
	}

	var diags diag.Diagnostics
	data.Config, diags = types.MapValueFrom(ctx, types.DynamicType, elems)
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

type configDataModel struct {
	Config types.Map `tfsdk:"config"`
}
