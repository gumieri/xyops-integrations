package provider

import (
	"context"
	"fmt"

	"github.com/gumieri/xyops-integrations/sdk"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type pluginResource struct {
	client *sdk.Client
}

func NewPluginResource() resource.Resource {
	return &pluginResource{}
}

func (r *pluginResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *pluginResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "xyops_plugin"
}

func (r *pluginResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "xyops plugin",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Plugin ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"title": schema.StringAttribute{
				Required:    true,
				Description: "Plugin title",
			},
			"type": schema.StringAttribute{
				Required:    true,
				Description: "Plugin type (event or monitor)",
			},
			"command": schema.StringAttribute{
				Required:    true,
				Description: "Plugin command to execute",
			},
			"script": schema.StringAttribute{
				Optional:    true,
				Description: "Plugin script",
			},
			"enabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether the plugin is enabled",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"params": schema.ListNestedAttribute{
				Optional:    true,
				Description: "Plugin parameters",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
						},
						"type": schema.StringAttribute{
							Required: true,
						},
						"title": schema.StringAttribute{
							Required: true,
						},
						"value": schema.StringAttribute{
							Optional: true,
						},
						"required": schema.BoolAttribute{
							Optional: true,
							Computed: true,
						},
						"locked": schema.BoolAttribute{
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"uid": schema.StringAttribute{
				Optional:    true,
				Description: "User ID the plugin runs as",
			},
			"gid": schema.StringAttribute{
				Optional:    true,
				Description: "Group ID the plugin runs as",
			},
			"groups": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Plugin groups",
			},
			"format": schema.StringAttribute{
				Optional:    true,
				Description: "Output format",
			},
		},
	}
}

func (r *pluginResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var config pluginModel
	diags := req.Plan.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	plugin := &sdk.Plugin{
		Title:   config.Title.ValueString(),
		Type:    config.Type.ValueString(),
		Command: config.Command.ValueString(),
		Script:  config.Script.ValueString(),
		Enabled: sdk.EnabledInt(boolToInt(config.Enabled.ValueBool())),
		Params:  pluginParamSliceFromSet(ctx, config.Params),
		UID:     config.UID.ValueString(),
		GID:     config.GID.ValueString(),
		Groups:  stringSliceFromSet(config.Groups),
		Format:  config.Format.ValueString(),
	}

	created, err := r.client.Plugins.Create(ctx, plugin)
	if err != nil {
		resp.Diagnostics.AddError("Error creating plugin", err.Error())
		return
	}

	state := pluginFromModel(config, created)
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *pluginResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state pluginModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	plugin, err := r.client.Plugins.Get(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading plugin", err.Error())
		return
	}

	state = pluginFromModel(state, plugin)
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *pluginResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state pluginModel
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

	plugin := &sdk.Plugin{
		ID:      state.ID.ValueString(),
		Title:   plan.Title.ValueString(),
		Type:    plan.Type.ValueString(),
		Command: plan.Command.ValueString(),
		Script:  plan.Script.ValueString(),
		Enabled: sdk.EnabledInt(boolToInt(plan.Enabled.ValueBool())),
		Params:  pluginParamSliceFromSet(ctx, plan.Params),
		UID:     plan.UID.ValueString(),
		GID:     plan.GID.ValueString(),
		Groups:  stringSliceFromSet(plan.Groups),
		Format:  plan.Format.ValueString(),
	}

	updated, err := r.client.Plugins.Update(ctx, plugin)
	if err != nil {
		resp.Diagnostics.AddError("Error updating plugin", err.Error())
		return
	}

	state = pluginFromModel(state, updated)
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *pluginResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state pluginModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	if err := r.client.Plugins.Delete(ctx, state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting plugin", err.Error())
		return
	}
}

func (r *pluginResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id := req.ID
	state := pluginModel{
		ID: types.StringValue(id),
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func pluginParamSliceFromSet(ctx context.Context, list types.List) []sdk.PluginParam {
	if list.IsNull() || list.IsUnknown() {
		return nil
	}

	elems := make([]pluginParamModel, 0, len(list.Elements()))
	diags := list.ElementsAs(ctx, &elems, false)
	if diags.HasError() {
		return nil
	}

	out := make([]sdk.PluginParam, 0, len(elems))
	for _, e := range elems {
		required := e.Required.ValueBool()
		locked := e.Locked.ValueBool()
		out = append(out, sdk.PluginParam{
			ID:       e.ID.ValueString(),
			Type:     e.Type.ValueString(),
			Title:    e.Title.ValueString(),
			Value:    e.Value.ValueString(),
			Required: &required,
			Locked:   &locked,
		})
	}
	return out
}

func pluginParamSetFromSlice(slice []sdk.PluginParam) types.List {
	if slice == nil {
		return types.ListNull(types.ObjectType{
			AttrTypes: pluginParamAttrTypes(),
		})
	}

	elems := make([]attr.Value, 0, len(slice))
	for _, p := range slice {
		required := false
		if p.Required != nil {
			required = *p.Required
		}
		locked := false
		if p.Locked != nil {
			locked = *p.Locked
		}
		objVal, _ := types.ObjectValue(pluginParamAttrTypes(), map[string]attr.Value{
			"id":       types.StringValue(p.ID),
			"type":     types.StringValue(p.Type),
			"title":    types.StringValue(p.Title),
			"value":    types.StringValue(fmt.Sprintf("%v", p.Value)),
			"required": types.BoolValue(required),
			"locked":   types.BoolValue(locked),
		})
		elems = append(elems, objVal)
	}

	list, _ := types.ListValue(types.ObjectType{AttrTypes: pluginParamAttrTypes()}, elems)
	return list
}

func pluginParamAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":       types.StringType,
		"type":     types.StringType,
		"title":    types.StringType,
		"value":    types.StringType,
		"required": types.BoolType,
		"locked":   types.BoolType,
	}
}

func pluginFromModel(model pluginModel, plugin *sdk.Plugin) pluginModel {
	model.ID = types.StringValue(plugin.ID)
	model.Title = types.StringValue(plugin.Title)
	model.Type = types.StringValue(plugin.Type)
	model.Command = types.StringValue(plugin.Command)
	model.Script = types.StringValue(plugin.Script)
	model.Enabled = types.BoolValue(intToBool(int(plugin.Enabled)))
	model.Params = pluginParamSetFromSlice(plugin.Params)
	model.UID = types.StringValue(plugin.UID)
	model.GID = types.StringValue(plugin.GID)
	model.Groups = stringSetFromSlice(plugin.Groups)
	model.Format = types.StringValue(plugin.Format)
	return model
}

type pluginParamModel struct {
	ID       types.String `tfsdk:"id"`
	Type     types.String `tfsdk:"type"`
	Title    types.String `tfsdk:"title"`
	Value    types.String `tfsdk:"value"`
	Required types.Bool   `tfsdk:"required"`
	Locked   types.Bool   `tfsdk:"locked"`
}

type pluginModel struct {
	ID      types.String `tfsdk:"id"`
	Title   types.String `tfsdk:"title"`
	Type    types.String `tfsdk:"type"`
	Command types.String `tfsdk:"command"`
	Script  types.String `tfsdk:"script"`
	Enabled types.Bool   `tfsdk:"enabled"`
	Params  types.List   `tfsdk:"params"`
	UID     types.String `tfsdk:"uid"`
	GID     types.String `tfsdk:"gid"`
	Groups  types.List   `tfsdk:"groups"`
	Format  types.String `tfsdk:"format"`
}
