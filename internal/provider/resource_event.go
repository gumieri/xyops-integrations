package provider

import (
	"context"
	"encoding/json"

	"github.com/gumieri/xyops-integrations/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type eventResource struct {
	client *client.Client
}

func NewEventResource() resource.Resource {
	return &eventResource{}
}

func (r *eventResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *eventResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "xyops_event"
}

func (r *eventResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "xyops event",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Event ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"title": schema.StringAttribute{
				Required:    true,
				Description: "Event title",
			},
			"plugin": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Plugin ID the event is attached to",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"targets": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Event targets",
			},
			"category": schema.StringAttribute{
				Optional:    true,
				Description: "Event category ID",
			},
			"algo": schema.StringAttribute{
				Optional:    true,
				Description: "Event algorithm",
			},
			"enabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether the event is enabled",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"fields": schema.ListNestedAttribute{
				Optional:    true,
				Description: "Event fields",
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
			"triggers": schema.ListNestedAttribute{
				Optional:    true,
				Description: "Event triggers",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Required: true,
						},
						"years": schema.ListAttribute{
							ElementType: types.Int64Type,
							Optional:    true,
						},
						"months": schema.ListAttribute{
							ElementType: types.Int64Type,
							Optional:    true,
						},
						"days": schema.ListAttribute{
							ElementType: types.Int64Type,
							Optional:    true,
						},
						"weekdays": schema.ListAttribute{
							ElementType: types.Int64Type,
							Optional:    true,
						},
						"hours": schema.ListAttribute{
							ElementType: types.Int64Type,
							Optional:    true,
						},
						"minutes": schema.ListAttribute{
							ElementType: types.Int64Type,
							Optional:    true,
						},
						"timezone": schema.StringAttribute{
							Optional: true,
						},
						"duration": schema.Int64Attribute{
							Optional: true,
						},
						"epoch": schema.Int64Attribute{
							Optional: true,
						},
					},
				},
			},
			"limits": schema.ListNestedAttribute{
				Optional:    true,
				Description: "Event limits (job, queue, retry)",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Required: true,
						},
						"enabled": schema.BoolAttribute{
							Optional: true,
							Computed: true,
						},
						"amount": schema.Int64Attribute{
							Optional: true,
						},
						"duration": schema.Int64Attribute{
							Optional: true,
						},
						"period": schema.StringAttribute{
							Optional: true,
						},
					},
				},
			},
			"actions": schema.ListNestedAttribute{
				Optional:    true,
				Description: "Event actions",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"condition": schema.StringAttribute{
							Optional: true,
						},
						"type": schema.StringAttribute{
							Required: true,
						},
						"email": schema.StringAttribute{
							Optional: true,
						},
						"web_hook": schema.StringAttribute{
							Optional: true,
						},
						"event": schema.StringAttribute{
							Optional: true,
						},
						"enabled": schema.BoolAttribute{
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"tags": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Event tags",
			},
			"notes": schema.StringAttribute{
				Optional:    true,
				Description: "Event notes",
			},
			"params_json": schema.StringAttribute{
				Optional:    true,
				Description: "Event parameters as a JSON string",
			},
			"expression": schema.StringAttribute{
				Optional:    true,
				Description: "Event expression (always sent as empty string)",
			},
			"icon": schema.StringAttribute{
				Optional:    true,
				Description: "Event icon (always sent as empty string)",
			},
			"workflow_json": schema.StringAttribute{
				Optional:    true,
				Description: "Workflow graph as a JSON string. Setting this makes the event a workflow.",
			},
		},
	}
}

func (r *eventResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var config eventModel
	diags := req.Plan.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	event, ok := eventFromModel(ctx, config)
	if !ok {
		resp.Diagnostics.AddError("Error building event", "invalid params_json or workflow_json")
		return
	}

	created, err := r.client.Events.Create(ctx, event)
	if err != nil {
		resp.Diagnostics.AddError("Error creating event", err.Error())
		return
	}

	state := eventToModel(config, created)
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *eventResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state eventModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	event, err := r.client.Events.Get(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading event", err.Error())
		return
	}

	state = eventToModel(state, event)
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *eventResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state eventModel
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

	event, ok := eventFromModel(ctx, plan)
	if !ok {
		resp.Diagnostics.AddError("Error building event", "invalid params_json or workflow_json")
		return
	}
	event.ID = state.ID.ValueString()

	updated, err := r.client.Events.Update(ctx, event)
	if err != nil {
		resp.Diagnostics.AddError("Error updating event", err.Error())
		return
	}

	state = eventToModel(state, updated)
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *eventResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state eventModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	if err := r.client.Events.Delete(ctx, state.ID.ValueString(), nil); err != nil {
		resp.Diagnostics.AddError("Error deleting event", err.Error())
		return
	}
}

func (r *eventResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id := req.ID
	state := eventModel{
		ID: types.StringValue(id),
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func eventFromModel(ctx context.Context, m eventModel) (*client.Event, bool) {
	event := &client.Event{
		Title:      m.Title.ValueString(),
		Plugin:     m.Plugin.ValueString(),
		Targets:    stringSliceFromSet(m.Targets),
		Category:   m.Category.ValueString(),
		Algo:       m.Algo.ValueString(),
		Enabled:    m.Enabled.ValueBool(),
		Fields:     fieldSliceFromSet(ctx, m.Fields),
		Triggers:   triggerSliceFromSet(ctx, m.Triggers),
		Limits:     limitSliceFromSet(ctx, m.Limits),
		Actions:    actionSliceFromSet(ctx, m.Actions),
		Tags:       stringSliceFromSet(m.Tags),
		Notes:      m.Notes.ValueString(),
		Expression: m.Expression.ValueString(),
		Icon:       m.Icon.ValueString(),
	}

	if !m.ParamsJSON.IsNull() && m.ParamsJSON.ValueString() != "" {
		var params any
		if err := json.Unmarshal([]byte(m.ParamsJSON.ValueString()), &params); err != nil {
			return nil, false
		}
		event.Params = params
	}

	if !m.WorkflowJSON.IsNull() && m.WorkflowJSON.ValueString() != "" {
		var wf client.WorkflowGraph
		if err := json.Unmarshal([]byte(m.WorkflowJSON.ValueString()), &wf); err != nil {
			return nil, false
		}
		event.Workflow = &wf
	}

	return event, true
}

func eventToModel(m eventModel, event *client.Event) eventModel {
	m.ID = types.StringValue(event.ID)
	m.Title = types.StringValue(event.Title)
	m.Plugin = types.StringValue(event.Plugin)
	m.Targets = stringSetFromSlice(event.Targets)
	m.Category = types.StringValue(event.Category)
	m.Algo = types.StringValue(event.Algo)
	m.Enabled = types.BoolValue(event.Enabled)
	m.Fields = fieldSetFromSlice(event.Fields)
	m.Triggers = triggerSetFromSlice(event.Triggers)
	m.Limits = limitSetFromSlice(event.Limits)
	m.Actions = actionSetFromSlice(event.Actions)
	m.Tags = stringSetFromSlice(event.Tags)
	m.Notes = types.StringValue(event.Notes)
	m.Expression = types.StringValue(event.Expression)
	m.Icon = types.StringValue(event.Icon)

	if event.Params != nil {
		if b, err := json.Marshal(event.Params); err == nil {
			m.ParamsJSON = types.StringValue(string(b))
		}
	}

	if event.Workflow != nil {
		if b, err := json.Marshal(event.Workflow); err == nil {
			m.WorkflowJSON = types.StringValue(string(b))
		}
	}

	return m
}

func fieldSliceFromSet(ctx context.Context, list types.List) []client.Field {
	if list.IsNull() || list.IsUnknown() {
		return nil
	}

	elems := make([]fieldModel, 0, len(list.Elements()))
	diags := list.ElementsAs(ctx, &elems, false)
	if diags.HasError() {
		return nil
	}

	out := make([]client.Field, 0, len(elems))
	for _, e := range elems {
		out = append(out, client.Field{
			ID:       e.ID.ValueString(),
			Type:     e.Type.ValueString(),
			Title:    e.Title.ValueString(),
			Value:    e.Value.ValueString(),
			Required: e.Required.ValueBool(),
			Locked:   e.Locked.ValueBool(),
		})
	}
	return out
}

func fieldSetFromSlice(slice []client.Field) types.List {
	if slice == nil {
		return types.ListNull(types.ObjectType{AttrTypes: fieldAttrTypes()})
	}

	elems := make([]attr.Value, 0, len(slice))
	for _, f := range slice {
		objVal, _ := types.ObjectValue(fieldAttrTypes(), map[string]attr.Value{
			"id":       types.StringValue(f.ID),
			"type":     types.StringValue(f.Type),
			"title":    types.StringValue(f.Title),
			"value":    types.StringValue(f.Value),
			"required": types.BoolValue(f.Required),
			"locked":   types.BoolValue(f.Locked),
		})
		elems = append(elems, objVal)
	}

	list, _ := types.ListValue(types.ObjectType{AttrTypes: fieldAttrTypes()}, elems)
	return list
}

func fieldAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":       types.StringType,
		"type":     types.StringType,
		"title":    types.StringType,
		"value":    types.StringType,
		"required": types.BoolType,
		"locked":   types.BoolType,
	}
}

func triggerSliceFromSet(ctx context.Context, list types.List) []client.Trigger {
	if list.IsNull() || list.IsUnknown() {
		return nil
	}

	elems := make([]triggerModel, 0, len(list.Elements()))
	diags := list.ElementsAs(ctx, &elems, false)
	if diags.HasError() {
		return nil
	}

	out := make([]client.Trigger, 0, len(elems))
	for _, e := range elems {
		out = append(out, client.Trigger{
			Type:     e.Type.ValueString(),
			Years:    intSliceFromSet(e.Years),
			Months:   intSliceFromSet(e.Months),
			Days:     intSliceFromSet(e.Days),
			Weekdays: intSliceFromSet(e.Weekdays),
			Hours:    intSliceFromSet(e.Hours),
			Minutes:  intSliceFromSet(e.Minutes),
			Timezone: e.Timezone.ValueString(),
			Duration: int(e.Duration.ValueInt64()),
			Epoch:    e.Epoch.ValueInt64(),
		})
	}
	return out
}

func triggerSetFromSlice(slice []client.Trigger) types.List {
	if slice == nil {
		return types.ListNull(types.ObjectType{AttrTypes: triggerAttrTypes()})
	}

	elems := make([]attr.Value, 0, len(slice))
	for _, t := range slice {
		objVal, _ := types.ObjectValue(triggerAttrTypes(), map[string]attr.Value{
			"type":     types.StringValue(t.Type),
			"years":    intSetFromSlice(t.Years),
			"months":   intSetFromSlice(t.Months),
			"days":     intSetFromSlice(t.Days),
			"weekdays": intSetFromSlice(t.Weekdays),
			"hours":    intSetFromSlice(t.Hours),
			"minutes":  intSetFromSlice(t.Minutes),
			"timezone": types.StringValue(t.Timezone),
			"duration": types.Int64Value(int64(t.Duration)),
			"epoch":    types.Int64Value(t.Epoch),
		})
		elems = append(elems, objVal)
	}

	list, _ := types.ListValue(types.ObjectType{AttrTypes: triggerAttrTypes()}, elems)
	return list
}

func triggerAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"type":     types.StringType,
		"years":    types.ListType{ElemType: types.Int64Type},
		"months":   types.ListType{ElemType: types.Int64Type},
		"days":     types.ListType{ElemType: types.Int64Type},
		"weekdays": types.ListType{ElemType: types.Int64Type},
		"hours":    types.ListType{ElemType: types.Int64Type},
		"minutes":  types.ListType{ElemType: types.Int64Type},
		"timezone": types.StringType,
		"duration": types.Int64Type,
		"epoch":    types.Int64Type,
	}
}

func limitSliceFromSet(ctx context.Context, list types.List) []client.Limit {
	if list.IsNull() || list.IsUnknown() {
		return nil
	}

	elems := make([]limitModel, 0, len(list.Elements()))
	diags := list.ElementsAs(ctx, &elems, false)
	if diags.HasError() {
		return nil
	}

	out := make([]client.Limit, 0, len(elems))
	for _, e := range elems {
		out = append(out, client.Limit{
			Type:     e.Type.ValueString(),
			Enabled:  e.Enabled.ValueBool(),
			Amount:   int(e.Amount.ValueInt64()),
			Duration: int(e.Duration.ValueInt64()),
			Period:   e.Period.ValueString(),
		})
	}
	return out
}

func limitSetFromSlice(slice []client.Limit) types.List {
	if slice == nil {
		return types.ListNull(types.ObjectType{AttrTypes: limitAttrTypes()})
	}

	elems := make([]attr.Value, 0, len(slice))
	for _, l := range slice {
		objVal, _ := types.ObjectValue(limitAttrTypes(), map[string]attr.Value{
			"type":     types.StringValue(l.Type),
			"enabled":  types.BoolValue(l.Enabled),
			"amount":   types.Int64Value(int64(l.Amount)),
			"duration": types.Int64Value(int64(l.Duration)),
			"period":   types.StringValue(l.Period),
		})
		elems = append(elems, objVal)
	}

	list, _ := types.ListValue(types.ObjectType{AttrTypes: limitAttrTypes()}, elems)
	return list
}

func limitAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"type":     types.StringType,
		"enabled":  types.BoolType,
		"amount":   types.Int64Type,
		"duration": types.Int64Type,
		"period":   types.StringType,
	}
}

func actionSliceFromSet(ctx context.Context, list types.List) []client.Action {
	if list.IsNull() || list.IsUnknown() {
		return nil
	}

	elems := make([]actionModel, 0, len(list.Elements()))
	diags := list.ElementsAs(ctx, &elems, false)
	if diags.HasError() {
		return nil
	}

	out := make([]client.Action, 0, len(elems))
	for _, e := range elems {
		out = append(out, client.Action{
			Condition: e.Condition.ValueString(),
			Type:      e.Type.ValueString(),
			Email:     e.Email.ValueString(),
			WebHook:   e.WebHook.ValueString(),
			Event:     e.Event.ValueString(),
			Enabled:   e.Enabled.ValueBool(),
		})
	}
	return out
}

func actionSetFromSlice(slice []client.Action) types.List {
	if slice == nil {
		return types.ListNull(types.ObjectType{AttrTypes: actionAttrTypes()})
	}

	elems := make([]attr.Value, 0, len(slice))
	for _, a := range slice {
		objVal, _ := types.ObjectValue(actionAttrTypes(), map[string]attr.Value{
			"condition": types.StringValue(a.Condition),
			"type":      types.StringValue(a.Type),
			"email":     types.StringValue(a.Email),
			"web_hook":  types.StringValue(a.WebHook),
			"event":     types.StringValue(a.Event),
			"enabled":   types.BoolValue(a.Enabled),
		})
		elems = append(elems, objVal)
	}

	list, _ := types.ListValue(types.ObjectType{AttrTypes: actionAttrTypes()}, elems)
	return list
}

func actionAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"condition": types.StringType,
		"type":      types.StringType,
		"email":     types.StringType,
		"web_hook":  types.StringType,
		"event":     types.StringType,
		"enabled":   types.BoolType,
	}
}

type fieldModel struct {
	ID       types.String `tfsdk:"id"`
	Type     types.String `tfsdk:"type"`
	Title    types.String `tfsdk:"title"`
	Value    types.String `tfsdk:"value"`
	Required types.Bool   `tfsdk:"required"`
	Locked   types.Bool   `tfsdk:"locked"`
}

type triggerModel struct {
	Type     types.String `tfsdk:"type"`
	Years    types.List   `tfsdk:"years"`
	Months   types.List   `tfsdk:"months"`
	Days     types.List   `tfsdk:"days"`
	Weekdays types.List   `tfsdk:"weekdays"`
	Hours    types.List   `tfsdk:"hours"`
	Minutes  types.List   `tfsdk:"minutes"`
	Timezone types.String `tfsdk:"timezone"`
	Duration types.Int64  `tfsdk:"duration"`
	Epoch    types.Int64  `tfsdk:"epoch"`
}

type limitModel struct {
	Type     types.String `tfsdk:"type"`
	Enabled  types.Bool   `tfsdk:"enabled"`
	Amount   types.Int64  `tfsdk:"amount"`
	Duration types.Int64  `tfsdk:"duration"`
	Period   types.String `tfsdk:"period"`
}

type actionModel struct {
	Condition types.String `tfsdk:"condition"`
	Type      types.String `tfsdk:"type"`
	Email     types.String `tfsdk:"email"`
	WebHook   types.String `tfsdk:"web_hook"`
	Event     types.String `tfsdk:"event"`
	Enabled   types.Bool   `tfsdk:"enabled"`
}

type eventModel struct {
	ID           types.String `tfsdk:"id"`
	Title        types.String `tfsdk:"title"`
	Plugin       types.String `tfsdk:"plugin"`
	Targets      types.List   `tfsdk:"targets"`
	Category     types.String `tfsdk:"category"`
	Algo         types.String `tfsdk:"algo"`
	Enabled      types.Bool   `tfsdk:"enabled"`
	Fields       types.List   `tfsdk:"fields"`
	Triggers     types.List   `tfsdk:"triggers"`
	Limits       types.List   `tfsdk:"limits"`
	Actions      types.List   `tfsdk:"actions"`
	Tags         types.List   `tfsdk:"tags"`
	Notes        types.String `tfsdk:"notes"`
	ParamsJSON   types.String `tfsdk:"params_json"`
	Expression   types.String `tfsdk:"expression"`
	Icon         types.String `tfsdk:"icon"`
	WorkflowJSON types.String `tfsdk:"workflow_json"`
}
