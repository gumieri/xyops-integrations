package provider

import (
	"context"

	"github.com/gumieri/xyops-integrations/sdk"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type secretResource struct {
	client *sdk.Client
}

func NewSecretResource() resource.Resource {
	return &secretResource{}
}

func (r *secretResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *secretResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "xyops_secret"
}

func (r *secretResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "xyops secret. Field values are encrypted and write-only.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Secret ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"title": schema.StringAttribute{
				Required:    true,
				Description: "Secret title",
			},
			"enabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether the secret is enabled",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"notes": schema.StringAttribute{
				Optional:    true,
				Description: "Secret notes",
			},
			"fields": schema.ListNestedAttribute{
				Optional:    true,
				Description: "Secret fields. Values are encrypted; on read the stored values are preserved.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required: true,
						},
						"value": schema.StringAttribute{
							Required:  true,
							Sensitive: true,
						},
					},
				},
			},
		},
	}
}

func (r *secretResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var config secretModel
	diags := req.Plan.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	secret := &sdk.Secret{
		Title:   config.Title.ValueString(),
		Enabled: config.Enabled.ValueBool(),
		Notes:   config.Notes.ValueString(),
		Fields:  secretFieldSliceFromSet(ctx, config.Fields),
	}

	created, err := r.client.Secrets.Create(ctx, secret)
	if err != nil {
		resp.Diagnostics.AddError("Error creating secret", err.Error())
		return
	}

	state := secretFromModel(config, created)
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *secretResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state secretModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	secret, err := r.client.Secrets.Get(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading secret", err.Error())
		return
	}

	// Fields values are encrypted write-only data; preserve the state's
	// fields rather than overwriting them with (unreadable) API values.
	fields := state.Fields
	state = secretFromModel(state, secret)
	state.Fields = fields

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *secretResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state secretModel
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

	secret := &sdk.Secret{
		ID:      state.ID.ValueString(),
		Title:   plan.Title.ValueString(),
		Enabled: plan.Enabled.ValueBool(),
		Notes:   plan.Notes.ValueString(),
		Fields:  secretFieldSliceFromSet(ctx, plan.Fields),
	}

	updated, err := r.client.Secrets.Update(ctx, secret)
	if err != nil {
		resp.Diagnostics.AddError("Error updating secret", err.Error())
		return
	}

	state = secretFromModel(state, updated)
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *secretResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state secretModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	if err := r.client.Secrets.Delete(ctx, state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting secret", err.Error())
		return
	}
}

func (r *secretResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id := req.ID
	state := secretModel{
		ID: types.StringValue(id),
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func secretFieldSliceFromSet(ctx context.Context, list types.List) []sdk.SecretField {
	if list.IsNull() || list.IsUnknown() {
		return nil
	}

	elems := make([]secretFieldModel, 0, len(list.Elements()))
	diags := list.ElementsAs(ctx, &elems, false)
	if diags.HasError() {
		return nil
	}

	out := make([]sdk.SecretField, 0, len(elems))
	for _, e := range elems {
		out = append(out, sdk.SecretField{
			Name:  e.Name.ValueString(),
			Value: e.Value.ValueString(),
		})
	}
	return out
}

func secretFieldSetFromSlice(slice []sdk.SecretField) types.List {
	if slice == nil {
		return types.ListNull(types.ObjectType{
			AttrTypes: secretFieldAttrTypes(),
		})
	}

	elems := make([]attr.Value, 0, len(slice))
	for _, f := range slice {
		objVal, _ := types.ObjectValue(secretFieldAttrTypes(), map[string]attr.Value{
			"name":  types.StringValue(f.Name),
			"value": types.StringValue(f.Value),
		})
		elems = append(elems, objVal)
	}

	list, _ := types.ListValue(types.ObjectType{AttrTypes: secretFieldAttrTypes()}, elems)
	return list
}

func secretFieldAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":  types.StringType,
		"value": types.StringType,
	}
}

func secretFromModel(model secretModel, secret *sdk.Secret) secretModel {
	model.ID = types.StringValue(secret.ID)
	model.Title = types.StringValue(secret.Title)
	model.Enabled = types.BoolValue(secret.Enabled)
	model.Notes = types.StringValue(secret.Notes)
	model.Fields = secretFieldSetFromSlice(secret.Fields)
	return model
}

type secretFieldModel struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}

type secretModel struct {
	ID      types.String `tfsdk:"id"`
	Title   types.String `tfsdk:"title"`
	Enabled types.Bool   `tfsdk:"enabled"`
	Notes   types.String `tfsdk:"notes"`
	Fields  types.List   `tfsdk:"fields"`
}
