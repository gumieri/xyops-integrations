package mcp

import (
	"context"

	"github.com/gumieri/xyops-integrations/sdk"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func registerTicketTools(s *mcp.Server, c *sdk.Client) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "search_tickets",
		Description: "Search tickets by query",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in struct {
		Query   string `json:"query,omitempty" jsonschema:"search query (Unbase/PxQL syntax)"`
		Offset  int    `json:"offset,omitempty" jsonschema:"offset for pagination"`
		Limit   int    `json:"limit,omitempty" jsonschema:"limit the number of results"`
		SortBy  string `json:"sort_by,omitempty" jsonschema:"sort field"`
		SortDir int    `json:"sort_dir,omitempty" jsonschema:"1 ascending, -1 descending"`
		Compact bool   `json:"compact,omitempty" jsonschema:"omit body and replace changes with a count"`
	}) (*mcp.CallToolResult, ListResult[sdk.Ticket], error) {
		opts := sdk.SearchTicketsOpts{
			Query:   in.Query,
			Offset:  in.Offset,
			Limit:   in.Limit,
			SortBy:  in.SortBy,
			SortDir: in.SortDir,
			Compact: in.Compact,
		}
		return adaptListResult(c.Tickets.Search(ctx, opts))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_ticket",
		Description: "Get a ticket by ID",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in struct {
		ID string `json:"id" jsonschema:"the ticket ID"`
	}) (*mcp.CallToolResult, *sdk.Ticket, error) {
		return adaptResult(c.Tickets.Get(ctx, in.ID))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "create_ticket",
		Description: "Create a ticket",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in sdk.Ticket) (*mcp.CallToolResult, *sdk.Ticket, error) {
		return adaptResult(c.Tickets.Create(ctx, &in))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "update_ticket",
		Description: "Update a ticket",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in sdk.Ticket) (*mcp.CallToolResult, *sdk.Ticket, error) {
		return adaptResult(c.Tickets.Update(ctx, &in))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "delete_ticket",
		Description: "Delete a ticket by ID",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in struct {
		ID string `json:"id" jsonschema:"the ticket ID"`
	}) (*mcp.CallToolResult, string, error) {
		return adaptResult("deleted", c.Tickets.Delete(ctx, in.ID))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "add_ticket_change",
		Description: "Add a change (e.g. a comment) to a ticket",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in struct {
		ID     string           `json:"id" jsonschema:"the ticket ID"`
		Change sdk.TicketChange `json:"change" jsonschema:"the change object"`
	}) (*mcp.CallToolResult, *sdk.Ticket, error) {
		return adaptResult(c.Tickets.AddChange(ctx, in.ID, in.Change))
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "update_ticket_change",
		Description: "Edit or delete a ticket change",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in struct {
		ID       string           `json:"id" jsonschema:"the ticket ID"`
		ChangeID string           `json:"change_id" jsonschema:"the change ID to edit or delete"`
		Change   sdk.TicketChange `json:"change,omitempty" jsonschema:"new change fields to merge"`
		Delete   bool             `json:"delete,omitempty" jsonschema:"delete the change"`
	}) (*mcp.CallToolResult, *sdk.Ticket, error) {
		return adaptResult(c.Tickets.UpdateChange(ctx, in.ID, in.ChangeID, in.Change, in.Delete))
	})
}
