package provider

import (
	"context"

	"github.com/gumieri/xyops-integrations/sdk"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type apiKeyResource struct {
	client *sdk.Client
}

func NewAPIKeyResource() resource.Resource {
	return &apiKeyResource{}
}

func (r *apiKeyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	c, ok := req.ProviderData.(*sdk.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			"Expected *sdk.Client, got: %T. Please report this bug to the provider developers.",
		)
		return
	}

	r.client = c
}

func (r *apiKeyResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "xyops_api_key"
}

func (r *apiKeyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "xyops API key",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "API key ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"title": schema.StringAttribute{
				Required:    true,
				Description: "API key title",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "API key description",
			},
			"key": schema.StringAttribute{
				Computed:    true,
				Sensitive:   true,
				Description: "API key value (sensitive; the usable value is returned on create only — v1.0.90+ reads return a server-side hash, v1.0.49 returns nothing)",
			},
			"active": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether the API key is active",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"privileges": schema.MapAttribute{
				ElementType: types.BoolType,
				Optional:    true,
				Description: "API key privileges. Map of privilege name to whether it is granted.",
			},
			"roles": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Roles granted to the API key",
			},
			"expires": schema.Int64Attribute{
				Optional:    true,
				Description: "Expiration timestamp (Unix seconds). 0 for no expiration.",
			},
		},
	}
}

func (r *apiKeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var config apiKeyModel
	diags := req.Plan.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	key := &sdk.APIKey{
		Title:       config.Title.ValueString(),
		Description: config.Description.ValueString(),
		Active:      boolToInt(config.Active.ValueBool()),
		Privileges:  frameworkMapToBoolMap(config.Privileges),
		Roles:       stringSliceFromSet(config.Roles),
		Expires:     config.Expires.ValueInt64(),
	}

	created, err := r.client.APIKeys.Create(ctx, key)
	if err != nil {
		resp.Diagnostics.AddError("Error creating API key", err.Error())
		return
	}

	state := apiKeyFromModel(config, created)
	state.Key = types.StringValue(created.Key)
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *apiKeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state apiKeyModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	key, err := r.client.APIKeys.Get(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading API key", err.Error())
		return
	}

	state = apiKeyFromModel(state, key)
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *apiKeyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state apiKeyModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	key := &sdk.APIKey{
		ID:          state.ID.ValueString(),
		Title:       plan.Title.ValueString(),
		Description: plan.Description.ValueString(),
		Active:      boolToInt(plan.Active.ValueBool()),
		Privileges:  frameworkMapToBoolMap(plan.Privileges),
		Roles:       stringSliceFromSet(plan.Roles),
		Expires:     plan.Expires.ValueInt64(),
	}

	updated, err := r.client.APIKeys.Update(ctx, key)
	if err != nil {
		resp.Diagnostics.AddError("Error updating API key", err.Error())
		return
	}

	// Key is write-only — apiKeyFromModel does not set it, so the prior state value is preserved.
	state = apiKeyFromModel(state, updated)
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *apiKeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state apiKeyModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	if err := r.client.APIKeys.Delete(ctx, state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting API key", err.Error())
		return
	}
}

func (r *apiKeyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id := req.ID
	state := apiKeyModel{
		ID: types.StringValue(id),
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func apiKeyFromModel(model apiKeyModel, key *sdk.APIKey) apiKeyModel {
	model.ID = types.StringValue(key.ID)
	model.Title = types.StringValue(key.Title)
	model.Description = types.StringValue(key.Description)
	model.Active = types.BoolValue(key.Active == 1)
	model.Privileges = boolMapToFrameworkMap(key.Privileges)
	model.Roles = stringSetFromSlice(key.Roles)
	model.Expires = types.Int64Value(key.Expires)
	return model
}

type apiKeyModel struct {
	ID          types.String `tfsdk:"id"`
	Title       types.String `tfsdk:"title"`
	Description types.String `tfsdk:"description"`
	Key         types.String `tfsdk:"key"`
	Active      types.Bool   `tfsdk:"active"`
	Privileges  types.Map    `tfsdk:"privileges"`
	Roles       types.List   `tfsdk:"roles"`
	Expires     types.Int64  `tfsdk:"expires"`
}
