package apiclient

import (
	"context"
	"fmt"
)

func (c *Client) GetSubscriber(ctx context.Context, pageID, subscriberID string) (*Subscriber, error) {
	var result Subscriber
	err := c.Get(ctx, fmt.Sprintf("/pages/%s/subscribers/%s", pageID, subscriberID), &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) CreateSubscriber(ctx context.Context, pageID string, body SubscriberBody) (*Subscriber, error) {
	var result Subscriber
	req := SubscriberRequest{Subscriber: body}
	err := c.Post(ctx, fmt.Sprintf("/pages/%s/subscribers", pageID), req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) UpdateSubscriber(ctx context.Context, pageID, subscriberID string, body SubscriberBody) (*Subscriber, error) {
	var result Subscriber
	req := SubscriberRequest{Subscriber: body}
	err := c.Patch(ctx, fmt.Sprintf("/pages/%s/subscribers/%s", pageID, subscriberID), req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) DeleteSubscriber(ctx context.Context, pageID, subscriberID string) error {
	return c.Delete(ctx, fmt.Sprintf("/pages/%s/subscribers/%s", pageID, subscriberID))
}
