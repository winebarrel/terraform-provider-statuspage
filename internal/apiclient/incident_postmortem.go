package apiclient

import (
	"context"
	"fmt"
)

func (c *Client) GetPostmortem(ctx context.Context, pageID, incidentID string) (*Postmortem, error) {
	var result Postmortem
	err := c.Get(ctx, fmt.Sprintf("/pages/%s/incidents/%s/postmortem", pageID, incidentID), &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) CreateOrUpdatePostmortem(ctx context.Context, pageID, incidentID string, body PostmortemBody) (*Postmortem, error) {
	var result Postmortem
	req := PostmortemRequest{Postmortem: body}
	err := c.Put(ctx, fmt.Sprintf("/pages/%s/incidents/%s/postmortem", pageID, incidentID), req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) DeletePostmortem(ctx context.Context, pageID, incidentID string) error {
	return c.Delete(ctx, fmt.Sprintf("/pages/%s/incidents/%s/postmortem", pageID, incidentID))
}
