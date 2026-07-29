package provider

import (
	"context"

	"github.com/gumieri/xyops-integrations/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type monitorResource struct {
	client *client.Client
}

func NewMonitorResource() resource.Resource {
	return &monitorResource{}
}

func (r *monitorResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *monitorResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "xyops_monitor"
}

func (r *monitorResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "xyops monitor",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Monitor ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"title": schema.StringAttribute{
				Required:    true,
				Description: "Monitor title",
			},
			"source": schema.StringAttribute{
				Required:    true,
				Description: "Monitor data source",
			},
			"data_type": schema.StringAttribute{
				Optional:    true,
				Description: "Monitor data type",
			},
			"display": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether the monitor is displayed on the dashboard",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"groups": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Monitor groups",
			},
			"suffix": schema.StringAttribute{
				Optional:    true,
				Description: "Value suffix",
			},
			"min_vert_scale": schema.Int64Attribute{
				Optional:    true,
				Description: "Minimum vertical scale",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"data_match": schema.StringAttribute{
				Optional:    true,
				Description: "Data match expression",
			},
			"delta": schema.BoolAttribute{
				Optional:    true,
				Description: "Compute delta between points",
			},
			"divide_time": schema.BoolAttribute{
				Optional:    true,
				Description: "Divide by elapsed time",
			},
			"zero_min": schema.Int64Attribute{
				Optional:    true,
				Description: "Force zero minimum on graph",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"notes": schema.StringAttribute{
				Optional:    true,
				Description: "Monitor notes",
			},
		},
	}
}

func (r *monitorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var config monitorModel
	diags := req.Plan.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	monitor := &client.Monitor{
		Title:        config.Title.ValueString(),
		Source:       config.Source.ValueString(),
		DataType:     config.DataType.ValueString(),
		Display:      config.Display.ValueBool(),
		Groups:       stringSliceFromSet(config.Groups),
		Suffix:       config.Suffix.ValueString(),
		MinVertScale: int(config.MinVertScale.ValueInt64()),
		DataMatch:    config.DataMatch.ValueString(),
		Delta:        config.Delta.ValueBool(),
		DivideTime:   config.DivideTime.ValueBool(),
		ZeroMin:      int(config.ZeroMin.ValueInt64()),
		Notes:        config.Notes.ValueString(),
	}

	created, err := r.client.Monitors.Create(ctx, monitor)
	if err != nil {
		resp.Diagnostics.AddError("Error creating monitor", err.Error())
		return
	}

	state := monitorFromModel(config, created)
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *monitorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state monitorModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	monitor, err := r.client.Monitors.Get(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading monitor", err.Error())
		return
	}

	state = monitorFromModel(state, monitor)
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *monitorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state monitorModel
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

	monitor := &client.Monitor{
		ID:           state.ID.ValueString(),
		Title:        plan.Title.ValueString(),
		Source:       plan.Source.ValueString(),
		DataType:     plan.DataType.ValueString(),
		Display:      plan.Display.ValueBool(),
		Groups:       stringSliceFromSet(plan.Groups),
		Suffix:       plan.Suffix.ValueString(),
		MinVertScale: int(plan.MinVertScale.ValueInt64()),
		DataMatch:    plan.DataMatch.ValueString(),
		Delta:        plan.Delta.ValueBool(),
		DivideTime:   plan.DivideTime.ValueBool(),
		ZeroMin:      int(plan.ZeroMin.ValueInt64()),
		Notes:        plan.Notes.ValueString(),
	}

	updated, err := r.client.Monitors.Update(ctx, monitor)
	if err != nil {
		resp.Diagnostics.AddError("Error updating monitor", err.Error())
		return
	}

	state = monitorFromModel(state, updated)
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *monitorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state monitorModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	if err := r.client.Monitors.Delete(ctx, state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting monitor", err.Error())
		return
	}
}

func (r *monitorResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id := req.ID
	state := monitorModel{
		ID: types.StringValue(id),
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func monitorFromModel(model monitorModel, monitor *client.Monitor) monitorModel {
	model.ID = types.StringValue(monitor.ID)
	model.Title = types.StringValue(monitor.Title)
	model.Source = types.StringValue(monitor.Source)
	model.DataType = types.StringValue(monitor.DataType)
	model.Display = types.BoolValue(monitor.Display)
	model.Groups = stringSetFromSlice(monitor.Groups)
	model.Suffix = types.StringValue(monitor.Suffix)
	model.MinVertScale = types.Int64Value(int64(monitor.MinVertScale))
	model.DataMatch = types.StringValue(monitor.DataMatch)
	model.Delta = types.BoolValue(monitor.Delta)
	model.DivideTime = types.BoolValue(monitor.DivideTime)
	model.ZeroMin = types.Int64Value(int64(monitor.ZeroMin))
	model.Notes = types.StringValue(monitor.Notes)
	return model
}

type monitorModel struct {
	ID           types.String `tfsdk:"id"`
	Title        types.String `tfsdk:"title"`
	Source       types.String `tfsdk:"source"`
	DataType     types.String `tfsdk:"data_type"`
	Display      types.Bool   `tfsdk:"display"`
	Groups       types.List   `tfsdk:"groups"`
	Suffix       types.String `tfsdk:"suffix"`
	MinVertScale types.Int64  `tfsdk:"min_vert_scale"`
	DataMatch    types.String `tfsdk:"data_match"`
	Delta        types.Bool   `tfsdk:"delta"`
	DivideTime   types.Bool   `tfsdk:"divide_time"`
	ZeroMin      types.Int64  `tfsdk:"zero_min"`
	Notes        types.String `tfsdk:"notes"`
}
