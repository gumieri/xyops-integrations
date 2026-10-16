package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

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
		Description: "xyops configuration (string-encoded; nested values are JSON)",
		Attributes: map[string]schema.Attribute{
			"config": schema.MapAttribute{
				Computed:    true,
				Description: "xyops configuration; scalar values as strings, nested structures JSON-encoded",
				ElementType: types.StringType,
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
		elems[k] = types.StringValue(configValueString(v))
	}

	var diags diag.Diagnostics
	data.Config, diags = types.MapValueFrom(ctx, types.StringType, elems)
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

// configValueString renders a config value as a string; nested structures are
// JSON-encoded. The framework forbids dynamic types inside collections, so
// the heterogeneous config map is stringly-typed.
func configValueString(v any) string {
	switch val := v.(type) {
	case string:
		return val
	case bool:
		return strconv.FormatBool(val)
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64)
	default:
		b, err := json.Marshal(val)
		if err != nil {
			return fmt.Sprintf("%v", val)
		}
		return string(b)
	}
}

type configDataModel struct {
	Config types.Map `tfsdk:"config"`
}
