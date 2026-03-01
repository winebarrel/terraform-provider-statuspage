package apiclient

import (
	"context"
	"fmt"
)

func (c *Client) GetIncident(ctx context.Context, pageID, incidentID string) (*Incident, error) {
	var result Incident
	err := c.Get(ctx, fmt.Sprintf("/pages/%s/incidents/%s", pageID, incidentID), &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) CreateIncident(ctx context.Context, pageID string, body IncidentBody) (*Incident, error) {
	var result Incident
	req := IncidentRequest{Incident: body}
	err := c.Post(ctx, fmt.Sprintf("/pages/%s/incidents", pageID), req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) UpdateIncident(ctx context.Context, pageID, incidentID string, body IncidentBody) (*Incident, error) {
	var result Incident
	req := IncidentRequest{Incident: body}
	err := c.Patch(ctx, fmt.Sprintf("/pages/%s/incidents/%s", pageID, incidentID), req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) DeleteIncident(ctx context.Context, pageID, incidentID string) error {
	return c.Delete(ctx, fmt.Sprintf("/pages/%s/incidents/%s", pageID, incidentID))
}
