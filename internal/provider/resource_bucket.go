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

type bucketResource struct {
	client *sdk.Client
}

func NewBucketResource() resource.Resource {
	return &bucketResource{}
}

func (r *bucketResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *bucketResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "xyops_bucket"
}

func (r *bucketResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "xyops bucket",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Bucket ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"title": schema.StringAttribute{
				Required:    true,
				Description: "Bucket title",
			},
			"enabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether the bucket is enabled",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"notes": schema.StringAttribute{
				Optional:    true,
				Description: "Bucket notes",
			},
		},
	}
}

func (r *bucketResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var config bucketModel
	diags := req.Plan.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	bucket := &sdk.Bucket{
		Title:   config.Title.ValueString(),
		Enabled: config.Enabled.ValueBool(),
		Notes:   config.Notes.ValueString(),
	}

	created, err := r.client.Buckets.Create(ctx, bucket)
	if err != nil {
		resp.Diagnostics.AddError("Error creating bucket", err.Error())
		return
	}

	state := bucketFromModel(config, created)
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *bucketResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state bucketModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	bucket, _, _, err := r.client.Buckets.Get(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading bucket", err.Error())
		return
	}

	state = bucketFromModel(state, bucket)
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *bucketResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state bucketModel
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

	bucket := &sdk.Bucket{
		ID:      state.ID.ValueString(),
		Title:   plan.Title.ValueString(),
		Enabled: plan.Enabled.ValueBool(),
		Notes:   plan.Notes.ValueString(),
	}

	updated, err := r.client.Buckets.Update(ctx, bucket)
	if err != nil {
		resp.Diagnostics.AddError("Error updating bucket", err.Error())
		return
	}

	state = bucketFromModel(state, updated)
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *bucketResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state bucketModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	if err := r.client.Buckets.Delete(ctx, state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting bucket", err.Error())
		return
	}
}

func (r *bucketResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id := req.ID
	state := bucketModel{
		ID: types.StringValue(id),
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func bucketFromModel(model bucketModel, bucket *sdk.Bucket) bucketModel {
	model.ID = types.StringValue(bucket.ID)
	model.Title = types.StringValue(bucket.Title)
	model.Enabled = types.BoolValue(bucket.Enabled)
	model.Notes = types.StringValue(bucket.Notes)
	return model
}

type bucketModel struct {
	ID      types.String `tfsdk:"id"`
	Title   types.String `tfsdk:"title"`
	Enabled types.Bool   `tfsdk:"enabled"`
	Notes   types.String `tfsdk:"notes"`
}
