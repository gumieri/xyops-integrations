package provider

import (
	"context"

	"github.com/gumieri/xyops-integrations/sdk"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type webHookResource struct {
	client *sdk.Client
}

func NewWebHookResource() resource.Resource {
	return &webHookResource{}
}

func (r *webHookResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *webHookResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "xyops_webhook"
}

func (r *webHookResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "xyops webhook",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Webhook ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"title": schema.StringAttribute{
				Required:    true,
				Description: "Webhook title",
			},
			"enabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether the webhook is enabled",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"url": schema.StringAttribute{
				Required:    true,
				Description: "Webhook destination URL",
			},
			"method": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "HTTP method (POST, GET, PUT, PATCH)",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"headers": schema.ListNestedAttribute{
				Optional:    true,
				Description: "Webhook headers",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required: true,
						},
						"value": schema.StringAttribute{
							Required:    true,
							Sensitive:   true,
							Description: "Header value. Webhook headers commonly carry auth tokens.",
						},
					},
				},
			},
			"body": schema.StringAttribute{
				Optional:    true,
				Description: "Request body",
			},
			"timeout": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Request timeout in seconds",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"retries": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Number of retries",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"follow": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Follow redirects",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"ssl_cert_bypass": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Bypass SSL certificate verification",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"max_per_day": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Maximum calls per day",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *webHookResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var config webHookModel
	diags := req.Plan.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	hook := &sdk.WebHook{
		Title:         config.Title.ValueString(),
		Enabled:       config.Enabled.ValueBool(),
		URL:           config.URL.ValueString(),
		Method:        config.Method.ValueString(),
		Headers:       headerSliceFromSet(config.Headers),
		Body:          config.Body.ValueString(),
		Timeout:       int(config.Timeout.ValueInt64()),
		Retries:       int(config.Retries.ValueInt64()),
		Follow:        config.Follow.ValueBool(),
		SSLCertBypass: config.SSLCertBypass.ValueBool(),
		MaxPerDay:     int(config.MaxPerDay.ValueInt64()),
	}

	created, err := r.client.WebHooks.Create(ctx, hook)
	if err != nil {
		resp.Diagnostics.AddError("Error creating webhook", err.Error())
		return
	}

	state := webHookFromModel(config, created)
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *webHookResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state webHookModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	hook, err := r.client.WebHooks.Get(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading webhook", err.Error())
		return
	}

	state = webHookFromModel(state, hook)
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *webHookResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state webHookModel
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

	hook := &sdk.WebHook{
		ID:            state.ID.ValueString(),
		Title:         plan.Title.ValueString(),
		Enabled:       plan.Enabled.ValueBool(),
		URL:           plan.URL.ValueString(),
		Method:        plan.Method.ValueString(),
		Headers:       headerSliceFromSet(plan.Headers),
		Body:          plan.Body.ValueString(),
		Timeout:       int(plan.Timeout.ValueInt64()),
		Retries:       int(plan.Retries.ValueInt64()),
		Follow:        plan.Follow.ValueBool(),
		SSLCertBypass: plan.SSLCertBypass.ValueBool(),
		MaxPerDay:     int(plan.MaxPerDay.ValueInt64()),
	}

	updated, err := r.client.WebHooks.Update(ctx, hook)
	if err != nil {
		resp.Diagnostics.AddError("Error updating webhook", err.Error())
		return
	}

	state = webHookFromModel(state, updated)
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *webHookResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state webHookModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	if err := r.client.WebHooks.Delete(ctx, state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting webhook", err.Error())
		return
	}
}

func (r *webHookResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id := req.ID
	state := webHookModel{
		ID: types.StringValue(id),
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func headerSliceFromSet(list types.List) []sdk.Header {
	if list.IsNull() || list.IsUnknown() {
		return nil
	}

	elems := make([]headerModel, 0, len(list.Elements()))
	diags := list.ElementsAs(context.Background(), &elems, false)
	if diags.HasError() {
		return nil
	}

	out := make([]sdk.Header, 0, len(elems))
	for _, e := range elems {
		out = append(out, sdk.Header{
			Name:  e.Name.ValueString(),
			Value: e.Value.ValueString(),
		})
	}
	return out
}

func headerSetFromSlice(slice []sdk.Header) types.List {
	if slice == nil {
		return types.ListNull(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"name":  types.StringType,
				"value": types.StringType,
			},
		})
	}

	elems := make([]attr.Value, 0, len(slice))
	for _, h := range slice {
		objVal, _ := types.ObjectValue(map[string]attr.Type{
			"name":  types.StringType,
			"value": types.StringType,
		}, map[string]attr.Value{
			"name":  types.StringValue(h.Name),
			"value": types.StringValue(h.Value),
		})
		elems = append(elems, objVal)
	}

	list, _ := types.ListValue(types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"name":  types.StringType,
			"value": types.StringType,
		},
	}, elems)
	return list
}

func webHookFromModel(model webHookModel, hook *sdk.WebHook) webHookModel {
	model.ID = types.StringValue(hook.ID)
	model.Title = types.StringValue(hook.Title)
	model.Enabled = types.BoolValue(hook.Enabled)
	model.URL = types.StringValue(hook.URL)
	model.Method = types.StringValue(hook.Method)
	model.Headers = headerSetFromSlice(hook.Headers)
	model.Body = types.StringValue(hook.Body)
	model.Timeout = types.Int64Value(int64(hook.Timeout))
	model.Retries = types.Int64Value(int64(hook.Retries))
	model.Follow = types.BoolValue(hook.Follow)
	model.SSLCertBypass = types.BoolValue(hook.SSLCertBypass)
	model.MaxPerDay = types.Int64Value(int64(hook.MaxPerDay))
	return model
}

type headerModel struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}

type webHookModel struct {
	ID            types.String `tfsdk:"id"`
	Title         types.String `tfsdk:"title"`
	Enabled       types.Bool   `tfsdk:"enabled"`
	URL           types.String `tfsdk:"url"`
	Method        types.String `tfsdk:"method"`
	Headers       types.List   `tfsdk:"headers"`
	Body          types.String `tfsdk:"body"`
	Timeout       types.Int64  `tfsdk:"timeout"`
	Retries       types.Int64  `tfsdk:"retries"`
	Follow        types.Bool   `tfsdk:"follow"`
	SSLCertBypass types.Bool   `tfsdk:"ssl_cert_bypass"`
	MaxPerDay     types.Int64  `tfsdk:"max_per_day"`
}
