package provider

import (
	"context"

	"github.com/gumieri/xyops-integrations/sdk"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ticketResource struct {
	client *sdk.Client
}

func NewTicketResource() resource.Resource {
	return &ticketResource{}
}

func (r *ticketResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ticketResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "xyops_ticket"
}

func (r *ticketResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "xyops ticket",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Ticket ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"num": schema.Int64Attribute{
				Computed:    true,
				Description: "Ticket number",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"subject": schema.StringAttribute{
				Required:    true,
				Description: "Ticket subject",
			},
			"body": schema.StringAttribute{
				Optional:    true,
				Description: "Ticket body (markdown)",
			},
			"type": schema.StringAttribute{
				Optional:    true,
				Description: "Ticket type",
			},
			"status": schema.StringAttribute{
				Optional:    true,
				Description: "Ticket status",
			},
			"category": schema.StringAttribute{
				Optional:    true,
				Description: "Ticket category",
			},
			"server": schema.StringAttribute{
				Optional:    true,
				Description: "Server hostname the ticket refers to",
			},
			"assignees": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Ticket assignees",
			},
			"cc": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Carbon-copied users",
			},
			"notify": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Users to notify",
			},
			"due": schema.Int64Attribute{
				Optional:    true,
				Description: "Due timestamp (Unix seconds)",
			},
			"tags": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Ticket tags",
			},
		},
	}
}

func (r *ticketResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var config ticketModel
	diags := req.Plan.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	ticket := ticketToSDK(config)

	created, err := r.client.Tickets.Create(ctx, ticket)
	if err != nil {
		resp.Diagnostics.AddError("Error creating ticket", err.Error())
		return
	}

	state := ticketFromModel(config, created)
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *ticketResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ticketModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	ticket, err := r.client.Tickets.Get(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading ticket", err.Error())
		return
	}

	state = ticketFromModel(state, ticket)
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *ticketResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state ticketModel
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

	ticket := ticketToSDK(plan)
	ticket.ID = state.ID.ValueString()

	updated, err := r.client.Tickets.Update(ctx, ticket)
	if err != nil {
		resp.Diagnostics.AddError("Error updating ticket", err.Error())
		return
	}

	state = ticketFromModel(state, updated)
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *ticketResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ticketModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	if err := r.client.Tickets.Delete(ctx, state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting ticket", err.Error())
		return
	}
}

func (r *ticketResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id := req.ID
	state := ticketModel{
		ID: types.StringValue(id),
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func ticketToSDK(model ticketModel) *sdk.Ticket {
	return &sdk.Ticket{
		Subject:   model.Subject.ValueString(),
		Body:      model.Body.ValueString(),
		Type:      model.Type.ValueString(),
		Status:    model.Status.ValueString(),
		Category:  model.Category.ValueString(),
		Server:    model.Server.ValueString(),
		Assignees: stringSliceFromSet(model.Assignees),
		CC:        stringSliceFromSet(model.CC),
		Notify:    stringSliceFromSet(model.Notify),
		Due:       model.Due.ValueInt64(),
		Tags:      stringSliceFromSet(model.Tags),
	}
}

func ticketFromModel(model ticketModel, ticket *sdk.Ticket) ticketModel {
	model.ID = types.StringValue(ticket.ID)
	model.Num = types.Int64Value(int64(ticket.Num))
	model.Subject = types.StringValue(ticket.Subject)
	model.Body = types.StringValue(ticket.Body)
	model.Type = types.StringValue(ticket.Type)
	model.Status = types.StringValue(ticket.Status)
	model.Category = types.StringValue(ticket.Category)
	model.Server = types.StringValue(ticket.Server)
	model.Assignees = stringSetFromSlice(ticket.Assignees)
	model.CC = stringSetFromSlice(ticket.CC)
	model.Notify = stringSetFromSlice(ticket.Notify)
	model.Due = types.Int64Value(ticket.Due)
	model.Tags = stringSetFromSlice(ticket.Tags)
	return model
}

type ticketModel struct {
	ID        types.String `tfsdk:"id"`
	Num       types.Int64  `tfsdk:"num"`
	Subject   types.String `tfsdk:"subject"`
	Body      types.String `tfsdk:"body"`
	Type      types.String `tfsdk:"type"`
	Status    types.String `tfsdk:"status"`
	Category  types.String `tfsdk:"category"`
	Server    types.String `tfsdk:"server"`
	Assignees types.List   `tfsdk:"assignees"`
	CC        types.List   `tfsdk:"cc"`
	Notify    types.List   `tfsdk:"notify"`
	Due       types.Int64  `tfsdk:"due"`
	Tags      types.List   `tfsdk:"tags"`
}
