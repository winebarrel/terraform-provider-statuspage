package apiclient

import (
	"context"
	"fmt"
)

func (c *Client) GetComponentGroup(ctx context.Context, pageID, groupID string) (*ComponentGroup, error) {
	var result ComponentGroup
	err := c.Get(ctx, fmt.Sprintf("/pages/%s/component-groups/%s", pageID, groupID), &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) CreateComponentGroup(ctx context.Context, pageID string, body ComponentGroupBody) (*ComponentGroup, error) {
	var result ComponentGroup
	req := ComponentGroupRequest{ComponentGroup: body}
	err := c.Post(ctx, fmt.Sprintf("/pages/%s/component-groups", pageID), req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) UpdateComponentGroup(ctx context.Context, pageID, groupID string, body ComponentGroupBody) (*ComponentGroup, error) {
	var result ComponentGroup
	req := ComponentGroupRequest{ComponentGroup: body}
	err := c.Patch(ctx, fmt.Sprintf("/pages/%s/component-groups/%s", pageID, groupID), req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) DeleteComponentGroup(ctx context.Context, pageID, groupID string) error {
	return c.Delete(ctx, fmt.Sprintf("/pages/%s/component-groups/%s", pageID, groupID))
}
