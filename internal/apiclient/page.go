package apiclient

import (
	"context"
	"fmt"
)

func (c *Client) GetPage(ctx context.Context, pageID string) (*Page, error) {
	var result Page
	err := c.Get(ctx, fmt.Sprintf("/pages/%s", pageID), &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) UpdatePage(ctx context.Context, pageID string, body PageBody) (*Page, error) {
	var result Page
	req := PageRequest{Page: body}
	err := c.Patch(ctx, fmt.Sprintf("/pages/%s", pageID), req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
