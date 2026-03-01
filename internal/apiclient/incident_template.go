package apiclient

import (
	"context"
	"fmt"
)

func (c *Client) GetIncidentTemplate(ctx context.Context, pageID, templateID string) (*IncidentTemplate, error) {
	// The API doesn't have a single-template GET endpoint, so list and filter.
	var results []IncidentTemplate
	err := c.Get(ctx, fmt.Sprintf("/pages/%s/incident_templates", pageID), &results)
	if err != nil {
		return nil, err
	}
	for _, t := range results {
		if t.ID == templateID {
			return &t, nil
		}
	}
	return nil, &APIError{StatusCode: 404, Message: "incident template not found"}
}

func (c *Client) CreateIncidentTemplate(ctx context.Context, pageID string, body IncidentTemplateBody) (*IncidentTemplate, error) {
	var result IncidentTemplate
	req := IncidentTemplateRequest{Template: body}
	err := c.Post(ctx, fmt.Sprintf("/pages/%s/incident_templates", pageID), req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) UpdateIncidentTemplate(ctx context.Context, pageID, templateID string, body IncidentTemplateBody) (*IncidentTemplate, error) {
	var result IncidentTemplate
	req := IncidentTemplateRequest{Template: body}
	err := c.Patch(ctx, fmt.Sprintf("/pages/%s/incident_templates/%s", pageID, templateID), req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) DeleteIncidentTemplate(ctx context.Context, pageID, templateID string) error {
	return c.Delete(ctx, fmt.Sprintf("/pages/%s/incident_templates/%s", pageID, templateID))
}
