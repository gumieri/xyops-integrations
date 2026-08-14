package sdk

import (
	"context"
	"fmt"
	"strconv"
)

// SearchTicketsOpts configures search_tickets/v1.
type SearchTicketsOpts struct {
	Query   string // Unbase/PxQL syntax, default "*"
	Offset  int
	Limit   int
	SortBy  string // default "_id"
	SortDir int    // 1 ascending, -1 descending
	Compact bool   // omit body, replace changes with count
}

// TicketsService
func (s *TicketsService) Get(ctx context.Context, id string) (*Ticket, error) {
	body, err := s.client.DoGETWithQuery(ctx, "get_ticket/v1", map[string]string{"id": id})
	if err != nil {
		return nil, err
	}

	var resp TicketResponse
	if err := unmarshalJSON(body, &resp); err != nil {
		return nil, err
	}
	return resp.Ticket, nil
}

// GetByNum retrieves a ticket by its auto-assigned number.
func (s *TicketsService) GetByNum(ctx context.Context, num int) (*Ticket, error) {
	body, err := s.client.DoGETWithQuery(ctx, "get_ticket/v1", map[string]string{"num": strconv.Itoa(num)})
	if err != nil {
		return nil, err
	}

	var resp TicketResponse
	if err := unmarshalJSON(body, &resp); err != nil {
		return nil, err
	}
	return resp.Ticket, nil
}

// Search queries tickets with PxQL syntax via search_tickets/v1 (GET only).
func (s *TicketsService) Search(ctx context.Context, opts SearchTicketsOpts) ([]Ticket, error) {
	params := map[string]string{}
	if opts.Query != "" {
		params["query"] = opts.Query
	}
	if opts.Offset > 0 {
		params["offset"] = strconv.Itoa(opts.Offset)
	}
	if opts.Limit > 0 {
		params["limit"] = strconv.Itoa(opts.Limit)
	}
	if opts.SortBy != "" {
		params["sort_by"] = opts.SortBy
	}
	if opts.SortDir != 0 {
		params["sort_dir"] = strconv.Itoa(opts.SortDir)
	}
	if opts.Compact {
		params["compact"] = "1"
	}

	body, err := s.client.DoGETWithQuery(ctx, "search_tickets/v1", params)
	if err != nil {
		return nil, err
	}

	var resp TicketListResponse
	if err := unmarshalJSON(body, &resp); err != nil {
		return nil, err
	}
	return resp.Rows, nil
}

// Create inserts a ticket. Multipart file uploads on create are not supported.
func (s *TicketsService) Create(ctx context.Context, ticket *Ticket) (*Ticket, error) {
	body, err := s.client.DoPOST(ctx, "create_ticket/v1", ticket)
	if err != nil {
		return nil, err
	}

	var resp TicketResponse
	if err := unmarshalJSON(body, &resp); err != nil {
		return nil, err
	}
	return resp.Ticket, nil
}

// Update performs a shallow merge via update_ticket/v1, then re-fetches the
// ticket. The update response includes a ticket object but its completeness is
// not guaranteed, so Get is used to return the authoritative record.
func (s *TicketsService) Update(ctx context.Context, ticket *Ticket) (*Ticket, error) {
	_, err := s.client.DoPOST(ctx, "update_ticket/v1", ticket)
	if err != nil {
		return nil, err
	}

	return s.Get(ctx, ticket.ID)
}

func (s *TicketsService) Delete(ctx context.Context, id string) error {
	_, err := s.client.DoPOST(ctx, "delete_ticket/v1", map[string]string{"id": id})
	return err
}

// AddChange appends a change (e.g. a comment) to a ticket via
// add_ticket_change/v1 and returns the updated ticket.
func (s *TicketsService) AddChange(ctx context.Context, ticketID string, change TicketChange) (*Ticket, error) {
	payload := map[string]any{"id": ticketID, "change": change}
	body, err := s.client.DoPOST(ctx, "add_ticket_change/v1", payload)
	if err != nil {
		return nil, err
	}

	var resp TicketResponse
	if err := unmarshalJSON(body, &resp); err != nil {
		return nil, err
	}
	if resp.Ticket == nil {
		return nil, fmt.Errorf("add ticket change %s: response did not include a ticket", ticketID)
	}
	return resp.Ticket, nil
}

// UpdateChange edits or deletes a ticket change via update_ticket_change/v1
// and returns the updated ticket.
func (s *TicketsService) UpdateChange(ctx context.Context, ticketID, changeID string, change TicketChange, del bool) (*Ticket, error) {
	payload := map[string]any{"id": ticketID, "change_id": changeID}
	if del {
		payload["delete"] = true
	} else {
		payload["change"] = change
	}
	body, err := s.client.DoPOST(ctx, "update_ticket_change/v1", payload)
	if err != nil {
		return nil, err
	}

	var resp TicketResponse
	if err := unmarshalJSON(body, &resp); err != nil {
		return nil, err
	}
	if resp.Ticket == nil {
		return nil, fmt.Errorf("update ticket change %s: response did not include a ticket", ticketID)
	}
	return resp.Ticket, nil
}
