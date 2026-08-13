package sdk

import "context"

// TagsService
func (s *TagsService) List(ctx context.Context) ([]Tag, error) {
	body, err := s.client.DoGET(ctx, "get_tags/v1")
	if err != nil {
		return nil, err
	}

	var resp TagListResponse
	if err := unmarshalJSON(body, &resp); err != nil {
		return nil, err
	}
	return resp.Rows, nil
}

func (s *TagsService) Get(ctx context.Context, id string) (*Tag, error) {
	body, err := s.client.DoGETWithQuery(ctx, "get_tag/v1", map[string]string{"id": id})
	if err != nil {
		return nil, err
	}

	var resp TagResponse
	if err := unmarshalJSON(body, &resp); err != nil {
		return nil, err
	}
	return resp.Tag, nil
}

func (s *TagsService) Create(ctx context.Context, tag *Tag) (*Tag, error) {
	if tag.ID != "" {
		return s.Update(ctx, tag)
	}

	body, err := s.client.DoPOST(ctx, "create_tag/v1", tag)
	if err != nil {
		return nil, err
	}

	var resp TagResponse
	if err := unmarshalJSON(body, &resp); err != nil {
		return nil, err
	}
	return resp.Tag, nil
}

func (s *TagsService) Update(ctx context.Context, tag *Tag) (*Tag, error) {
	_, err := s.client.DoPOST(ctx, "update_tag/v1", tag)
	if err != nil {
		return nil, err
	}

	return s.Get(ctx, tag.ID)
}

func (s *TagsService) Delete(ctx context.Context, id string) error {
	_, err := s.client.DoPOST(ctx, "delete_tag/v1", map[string]string{"id": id})
	return err
}

// GroupsService
func (s *GroupsService) List(ctx context.Context) ([]Group, error) {
	body, err := s.client.DoGET(ctx, "get_groups/v1")
	if err != nil {
		return nil, err
	}

	var resp GroupListResponse
	if err := unmarshalJSON(body, &resp); err != nil {
		return nil, err
	}
	return resp.Rows, nil
}

func (s *GroupsService) Get(ctx context.Context, id string) (*Group, error) {
	body, err := s.client.DoGETWithQuery(ctx, "get_group/v1", map[string]string{"id": id})
	if err != nil {
		return nil, err
	}

	var resp GroupResponse
	if err := unmarshalJSON(body, &resp); err != nil {
		return nil, err
	}
	return resp.Group, nil
}

func (s *GroupsService) Create(ctx context.Context, group *Group) (*Group, error) {
	if group.ID != "" {
		return s.Update(ctx, group)
	}

	body, err := s.client.DoPOST(ctx, "create_group/v1", group)
	if err != nil {
		return nil, err
	}

	var resp GroupResponse
	if err := unmarshalJSON(body, &resp); err != nil {
		return nil, err
	}
	return resp.Group, nil
}

func (s *GroupsService) Update(ctx context.Context, group *Group) (*Group, error) {
	_, err := s.client.DoPOST(ctx, "update_group/v1", group)
	if err != nil {
		return nil, err
	}

	return s.Get(ctx, group.ID)
}

func (s *GroupsService) Delete(ctx context.Context, id string) error {
	_, err := s.client.DoPOST(ctx, "delete_group/v1", map[string]string{"id": id})
	return err
}

// CategoriesService
func (s *CategoriesService) List(ctx context.Context) ([]Category, error) {
	body, err := s.client.DoGET(ctx, "get_categories/v1")
	if err != nil {
		return nil, err
	}

	var resp CategoryListResponse
	if err := unmarshalJSON(body, &resp); err != nil {
		return nil, err
	}
	return resp.Rows, nil
}

func (s *CategoriesService) Get(ctx context.Context, id string) (*Category, error) {
	body, err := s.client.DoGETWithQuery(ctx, "get_category/v1", map[string]string{"id": id})
	if err != nil {
		return nil, err
	}

	var resp CategoryResponse
	if err := unmarshalJSON(body, &resp); err != nil {
		return nil, err
	}
	return resp.Category, nil
}

func (s *CategoriesService) Create(ctx context.Context, category *Category) (*Category, error) {
	if category.ID != "" {
		return s.Update(ctx, category)
	}

	body, err := s.client.DoPOST(ctx, "create_category/v1", category)
	if err != nil {
		return nil, err
	}

	var resp CategoryResponse
	if err := unmarshalJSON(body, &resp); err != nil {
		return nil, err
	}
	return resp.Category, nil
}

func (s *CategoriesService) Update(ctx context.Context, category *Category) (*Category, error) {
	_, err := s.client.DoPOST(ctx, "update_category/v1", category)
	if err != nil {
		return nil, err
	}

	return s.Get(ctx, category.ID)
}

func (s *CategoriesService) Delete(ctx context.Context, id string) error {
	_, err := s.client.DoPOST(ctx, "delete_category/v1", map[string]string{"id": id})
	return err
}
