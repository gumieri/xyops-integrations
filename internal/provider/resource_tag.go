package provider

import (
	"context"

	"github.com/gumieri/xyops-integrations/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type tagResource struct {
	client *client.Client
}

func NewTagResource() resource.Resource {
	return &tagResource{}
}

func (r *tagResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			"Expected *client.Client, got: %T. Please report this bug to the provider developers.",
		)
		return
	}

	r.client = c
}

func (r *tagResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "xyops_tag"
}

func (r *tagResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "xyops tag",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Tag ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"title": schema.StringAttribute{
				Required:    true,
				Description: "Tag title",
			},
		},
	}
}

func (r *tagResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var config tagModel
	diags := req.Plan.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	tag := &client.Tag{
		Title: config.Title.ValueString(),
	}

	created, err := r.client.Tags.Create(ctx, tag)
	if err != nil {
		resp.Diagnostics.AddError("Error creating tag", err.Error())
		return
	}

	config.ID = types.StringValue(created.ID)
	config.Title = types.StringValue(created.Title)

	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
}

func (r *tagResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state tagModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	tag, err := r.client.Tags.Get(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading tag", err.Error())
		return
	}

	state.ID = types.StringValue(tag.ID)
	state.Title = types.StringValue(tag.Title)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *tagResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state tagModel
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

	tag := &client.Tag{
		ID:    state.ID.ValueString(),
		Title: plan.Title.ValueString(),
	}

	updated, err := r.client.Tags.Update(ctx, tag)
	if err != nil {
		resp.Diagnostics.AddError("Error updating tag", err.Error())
		return
	}

	state.ID = types.StringValue(updated.ID)
	state.Title = types.StringValue(updated.Title)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *tagResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state tagModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	if err := r.client.Tags.Delete(ctx, state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting tag", err.Error())
		return
	}
}

func (r *tagResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id := req.ID
	state := tagModel{
		ID: types.StringValue(id),
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

type tagModel struct {
	ID    types.String `tfsdk:"id"`
	Title types.String `tfsdk:"title"`
}
