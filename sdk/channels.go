package sdk

import "context"

// ChannelsService
func (s *ChannelsService) List(ctx context.Context) ([]Channel, error) {
	body, err := s.client.DoGET(ctx, "get_channels/v1")
	if err != nil {
		return nil, err
	}

	var resp ChannelListResponse
	if err := unmarshalJSON(body, &resp); err != nil {
		return nil, err
	}
	return resp.Rows, nil
}

func (s *ChannelsService) Get(ctx context.Context, id string) (*Channel, error) {
	body, err := s.client.DoGETWithQuery(ctx, "get_channel/v1", map[string]string{"id": id})
	if err != nil {
		return nil, err
	}

	var resp ChannelResponse
	if err := unmarshalJSON(body, &resp); err != nil {
		return nil, err
	}
	return resp.Channel, nil
}

// Create inserts a channel. If the channel already carries an ID, the call is
// routed to Update (upsert-style, matching other resources).
func (s *ChannelsService) Create(ctx context.Context, channel *Channel) (*Channel, error) {
	if channel.ID != "" {
		return s.Update(ctx, channel)
	}

	body, err := s.client.DoPOST(ctx, "create_channel/v1", channel)
	if err != nil {
		return nil, err
	}

	var resp ChannelResponse
	if err := unmarshalJSON(body, &resp); err != nil {
		return nil, err
	}
	return resp.Channel, nil
}

// Update performs a shallow merge via update_channel/v1, then re-fetches the
// channel (the update response is only {code: 0}).
func (s *ChannelsService) Update(ctx context.Context, channel *Channel) (*Channel, error) {
	_, err := s.client.DoPOST(ctx, "update_channel/v1", channel)
	if err != nil {
		return nil, err
	}

	return s.Get(ctx, channel.ID)
}

func (s *ChannelsService) Delete(ctx context.Context, id string) error {
	_, err := s.client.DoPOST(ctx, "delete_channel/v1", map[string]string{"id": id})
	return err
}
