package apiclient

import (
	"context"
	"fmt"
)

func (c *Client) GetMetric(ctx context.Context, pageID, metricID string) (*Metric, error) {
	var result Metric
	err := c.Get(ctx, fmt.Sprintf("/pages/%s/metrics/%s", pageID, metricID), &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) CreateMetric(ctx context.Context, pageID, metricsProviderID string, body MetricBody) (*Metric, error) {
	var result Metric
	req := MetricRequest{Metric: body}
	err := c.Post(ctx, fmt.Sprintf("/pages/%s/metrics_providers/%s/metrics", pageID, metricsProviderID), req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) UpdateMetric(ctx context.Context, pageID, metricID string, body MetricBody) (*Metric, error) {
	var result Metric
	req := MetricRequest{Metric: body}
	err := c.Patch(ctx, fmt.Sprintf("/pages/%s/metrics/%s", pageID, metricID), req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) DeleteMetric(ctx context.Context, pageID, metricID string) error {
	return c.Delete(ctx, fmt.Sprintf("/pages/%s/metrics/%s", pageID, metricID))
}
