package apiclient

import (
	"context"
	"fmt"
)

func (c *Client) GetPageAccessGroup(ctx context.Context, pageID, groupID string) (*PageAccessGroup, error) {
	var result PageAccessGroup
	err := c.Get(ctx, fmt.Sprintf("/pages/%s/page_access_groups/%s", pageID, groupID), &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) CreatePageAccessGroup(ctx context.Context, pageID string, body PageAccessGroupBody) (*PageAccessGroup, error) {
	var result PageAccessGroup
	req := PageAccessGroupRequest{PageAccessGroup: body}
	err := c.Post(ctx, fmt.Sprintf("/pages/%s/page_access_groups", pageID), req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) UpdatePageAccessGroup(ctx context.Context, pageID, groupID string, body PageAccessGroupBody) (*PageAccessGroup, error) {
	var result PageAccessGroup
	req := PageAccessGroupRequest{PageAccessGroup: body}
	err := c.Patch(ctx, fmt.Sprintf("/pages/%s/page_access_groups/%s", pageID, groupID), req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) DeletePageAccessGroup(ctx context.Context, pageID, groupID string) error {
	return c.Delete(ctx, fmt.Sprintf("/pages/%s/page_access_groups/%s", pageID, groupID))
}
