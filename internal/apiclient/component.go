package apiclient

import (
	"context"
	"fmt"
)

func (c *Client) GetComponent(ctx context.Context, pageID, componentID string) (*Component, error) {
	var result Component
	err := c.Get(ctx, fmt.Sprintf("/pages/%s/components/%s", pageID, componentID), &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) CreateComponent(ctx context.Context, pageID string, body ComponentBody) (*Component, error) {
	var result Component
	req := ComponentRequest{Component: body}
	err := c.Post(ctx, fmt.Sprintf("/pages/%s/components", pageID), req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) UpdateComponent(ctx context.Context, pageID, componentID string, body ComponentBody) (*Component, error) {
	var result Component
	req := ComponentRequest{Component: body}
	err := c.Patch(ctx, fmt.Sprintf("/pages/%s/components/%s", pageID, componentID), req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) DeleteComponent(ctx context.Context, pageID, componentID string) error {
	return c.Delete(ctx, fmt.Sprintf("/pages/%s/components/%s", pageID, componentID))
}
