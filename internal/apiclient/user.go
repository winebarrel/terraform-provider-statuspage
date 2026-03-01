package apiclient

import (
	"context"
	"fmt"
)

func (c *Client) GetUser(ctx context.Context, organizationID, userID string) (*User, error) {
	// The API doesn't have a single-user GET endpoint, so list and filter.
	var results []User
	err := c.Get(ctx, fmt.Sprintf("/organizations/%s/users", organizationID), &results)
	if err != nil {
		return nil, err
	}
	for _, u := range results {
		if u.ID == userID {
			return &u, nil
		}
	}
	return nil, &APIError{StatusCode: 404, Message: "user not found"}
}

func (c *Client) CreateUser(ctx context.Context, organizationID string, body UserBody) (*User, error) {
	var result User
	req := UserRequest{User: body}
	err := c.Post(ctx, fmt.Sprintf("/organizations/%s/users", organizationID), req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) DeleteUser(ctx context.Context, organizationID, userID string) error {
	return c.Delete(ctx, fmt.Sprintf("/organizations/%s/users/%s", organizationID, userID))
}
