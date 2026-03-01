package apiclient

import (
	"context"
	"fmt"
)

func (c *Client) GetStatusEmbedConfig(ctx context.Context, pageID string) (*StatusEmbedConfig, error) {
	var result StatusEmbedConfig
	err := c.Get(ctx, fmt.Sprintf("/pages/%s/status_embed_config", pageID), &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) UpdateStatusEmbedConfig(ctx context.Context, pageID string, body StatusEmbedConfigBody) (*StatusEmbedConfig, error) {
	var result StatusEmbedConfig
	req := StatusEmbedConfigRequest{StatusEmbedConfig: body}
	err := c.Patch(ctx, fmt.Sprintf("/pages/%s/status_embed_config", pageID), req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
