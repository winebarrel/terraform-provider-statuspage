package apiclient

import (
	"context"
	"fmt"
)

func (c *Client) GetPageAccessUser(ctx context.Context, pageID, userID string) (*PageAccessUser, error) {
	var result PageAccessUser
	err := c.Get(ctx, fmt.Sprintf("/pages/%s/page_access_users/%s", pageID, userID), &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) CreatePageAccessUser(ctx context.Context, pageID string, body PageAccessUserBody) (*PageAccessUser, error) {
	var result PageAccessUser
	req := PageAccessUserRequest{PageAccessUser: body}
	err := c.Post(ctx, fmt.Sprintf("/pages/%s/page_access_users", pageID), req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) UpdatePageAccessUser(ctx context.Context, pageID, userID string, body PageAccessUserBody) (*PageAccessUser, error) {
	var result PageAccessUser
	req := PageAccessUserRequest{PageAccessUser: body}
	err := c.Patch(ctx, fmt.Sprintf("/pages/%s/page_access_users/%s", pageID, userID), req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) DeletePageAccessUser(ctx context.Context, pageID, userID string) error {
	return c.Delete(ctx, fmt.Sprintf("/pages/%s/page_access_users/%s", pageID, userID))
}
