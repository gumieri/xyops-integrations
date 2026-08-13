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

type roleResource struct {
	client *sdk.Client
}

func NewRoleResource() resource.Resource {
	return &roleResource{}
}

func (r *roleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *roleResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "xyops_role"
}

func (r *roleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "xyops role",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Role ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"title": schema.StringAttribute{
				Required:    true,
				Description: "Role title",
			},
			"enabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether the role is enabled",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"privileges": schema.MapAttribute{
				ElementType: types.BoolType,
				Required:    true,
				Description: "Role privileges. Map of privilege name to whether it is granted.",
			},
			"categories": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Categories the role applies to",
			},
			"notes": schema.StringAttribute{
				Optional:    true,
				Description: "Role notes",
			},
		},
	}
}

func (r *roleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var config roleModel
	diags := req.Plan.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	role := &sdk.Role{
		Title:      config.Title.ValueString(),
		Enabled:    config.Enabled.ValueBool(),
		Privileges: boolMapToIntMap(config.Privileges),
		Categories: stringSliceFromSet(config.Categories),
		Notes:      config.Notes.ValueString(),
	}

	created, err := r.client.Roles.Create(ctx, role)
	if err != nil {
		resp.Diagnostics.AddError("Error creating role", err.Error())
		return
	}

	state := roleFromModel(config, created)
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *roleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state roleModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	role, err := r.client.Roles.Get(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading role", err.Error())
		return
	}

	state = roleFromModel(state, role)
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *roleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state roleModel
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

	role := &sdk.Role{
		ID:         state.ID.ValueString(),
		Title:      plan.Title.ValueString(),
		Enabled:    plan.Enabled.ValueBool(),
		Privileges: boolMapToIntMap(plan.Privileges),
		Categories: stringSliceFromSet(plan.Categories),
		Notes:      plan.Notes.ValueString(),
	}

	updated, err := r.client.Roles.Update(ctx, role)
	if err != nil {
		resp.Diagnostics.AddError("Error updating role", err.Error())
		return
	}

	state = roleFromModel(state, updated)
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *roleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state roleModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	if err := r.client.Roles.Delete(ctx, state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting role", err.Error())
		return
	}
}

func (r *roleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id := req.ID
	state := roleModel{
		ID: types.StringValue(id),
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func roleFromModel(model roleModel, role *sdk.Role) roleModel {
	model.ID = types.StringValue(role.ID)
	model.Title = types.StringValue(role.Title)
	model.Enabled = types.BoolValue(role.Enabled)
	model.Privileges = intMapToBoolMap(role.Privileges)
	model.Categories = stringSetFromSlice(role.Categories)
	model.Notes = types.StringValue(role.Notes)
	return model
}

type roleModel struct {
	ID         types.String `tfsdk:"id"`
	Title      types.String `tfsdk:"title"`
	Enabled    types.Bool   `tfsdk:"enabled"`
	Privileges types.Map    `tfsdk:"privileges"`
	Categories types.List   `tfsdk:"categories"`
	Notes      types.String `tfsdk:"notes"`
}
