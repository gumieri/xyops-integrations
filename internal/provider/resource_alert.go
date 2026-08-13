package provider

import (
	"context"

	"github.com/gumieri/xyops-integrations/sdk"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type alertResource struct {
	client *sdk.Client
}

func NewAlertResource() resource.Resource {
	return &alertResource{}
}

func (r *alertResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *alertResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "xyops_alert"
}

func (r *alertResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "xyops alert",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Alert ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"title": schema.StringAttribute{
				Required:    true,
				Description: "Alert title",
			},
			"enabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether the alert is enabled",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"expression": schema.StringAttribute{
				Required:    true,
				Description: "Alert expression",
			},
			"message": schema.StringAttribute{
				Optional:    true,
				Description: "Alert message",
			},
			"groups": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Alert groups",
			},
			"monitor_id": schema.StringAttribute{
				Optional:    true,
				Description: "Monitor ID the alert is attached to",
			},
			"max_count": schema.Int64Attribute{
				Optional:    true,
				Description: "Maximum number of alerts",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"notes": schema.StringAttribute{
				Optional:    true,
				Description: "Alert notes",
			},
		},
	}
}

func (r *alertResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var config alertModel
	diags := req.Plan.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	alert := &sdk.Alert{
		Title:      config.Title.ValueString(),
		Enabled:    config.Enabled.ValueBool(),
		Expression: config.Expression.ValueString(),
		Message:    config.Message.ValueString(),
		Groups:     stringSliceFromSet(config.Groups),
		MonitorID:  config.MonitorID.ValueString(),
		MaxCount:   int(config.MaxCount.ValueInt64()),
		Notes:      config.Notes.ValueString(),
	}

	created, err := r.client.Alerts.Create(ctx, alert)
	if err != nil {
		resp.Diagnostics.AddError("Error creating alert", err.Error())
		return
	}

	state := alertFromModel(config, created)
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *alertResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state alertModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	alert, err := r.client.Alerts.Get(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading alert", err.Error())
		return
	}

	state = alertFromModel(state, alert)
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *alertResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state alertModel
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

	alert := &sdk.Alert{
		ID:         state.ID.ValueString(),
		Title:      plan.Title.ValueString(),
		Enabled:    plan.Enabled.ValueBool(),
		Expression: plan.Expression.ValueString(),
		Message:    plan.Message.ValueString(),
		Groups:     stringSliceFromSet(plan.Groups),
		MonitorID:  plan.MonitorID.ValueString(),
		MaxCount:   int(plan.MaxCount.ValueInt64()),
		Notes:      plan.Notes.ValueString(),
	}

	updated, err := r.client.Alerts.Update(ctx, alert)
	if err != nil {
		resp.Diagnostics.AddError("Error updating alert", err.Error())
		return
	}

	state = alertFromModel(state, updated)
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *alertResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state alertModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	if err := r.client.Alerts.Delete(ctx, state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting alert", err.Error())
		return
	}
}

func (r *alertResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id := req.ID
	state := alertModel{
		ID: types.StringValue(id),
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func alertFromModel(model alertModel, alert *sdk.Alert) alertModel {
	model.ID = types.StringValue(alert.ID)
	model.Title = types.StringValue(alert.Title)
	model.Enabled = types.BoolValue(alert.Enabled)
	model.Expression = types.StringValue(alert.Expression)
	model.Message = types.StringValue(alert.Message)
	model.Groups = stringSetFromSlice(alert.Groups)
	model.MonitorID = types.StringValue(alert.MonitorID)
	model.MaxCount = types.Int64Value(int64(alert.MaxCount))
	model.Notes = types.StringValue(alert.Notes)
	return model
}

type alertModel struct {
	ID         types.String `tfsdk:"id"`
	Title      types.String `tfsdk:"title"`
	Enabled    types.Bool   `tfsdk:"enabled"`
	Expression types.String `tfsdk:"expression"`
	Message    types.String `tfsdk:"message"`
	Groups     types.List   `tfsdk:"groups"`
	MonitorID  types.String `tfsdk:"monitor_id"`
	MaxCount   types.Int64  `tfsdk:"max_count"`
	Notes      types.String `tfsdk:"notes"`
}
