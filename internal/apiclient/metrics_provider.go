package apiclient

import (
	"context"
	"fmt"
)

func (c *Client) GetMetricsProvider(ctx context.Context, pageID, providerID string) (*MetricsProvider, error) {
	var result MetricsProvider
	err := c.Get(ctx, fmt.Sprintf("/pages/%s/metrics_providers/%s", pageID, providerID), &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) CreateMetricsProvider(ctx context.Context, pageID string, body MetricsProviderBody) (*MetricsProvider, error) {
	var result MetricsProvider
	req := MetricsProviderRequest{MetricsProvider: body}
	err := c.Post(ctx, fmt.Sprintf("/pages/%s/metrics_providers", pageID), req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) UpdateMetricsProvider(ctx context.Context, pageID, providerID string, body MetricsProviderBody) (*MetricsProvider, error) {
	var result MetricsProvider
	req := MetricsProviderRequest{MetricsProvider: body}
	err := c.Patch(ctx, fmt.Sprintf("/pages/%s/metrics_providers/%s", pageID, providerID), req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) DeleteMetricsProvider(ctx context.Context, pageID, providerID string) error {
	return c.Delete(ctx, fmt.Sprintf("/pages/%s/metrics_providers/%s", pageID, providerID))
}
