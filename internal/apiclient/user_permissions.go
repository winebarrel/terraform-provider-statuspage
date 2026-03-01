package apiclient

import (
	"context"
	"fmt"
)

func (c *Client) GetPermissions(ctx context.Context, organizationID, userID string) (*Permissions, error) {
	var result Permissions
	err := c.Get(ctx, fmt.Sprintf("/organizations/%s/permissions/%s", organizationID, userID), &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) UpdatePermissions(ctx context.Context, organizationID, userID string, body PermissionsRequest) (*Permissions, error) {
	var result Permissions
	err := c.Put(ctx, fmt.Sprintf("/organizations/%s/permissions/%s", organizationID, userID), body, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
