package apiclient

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetSubscriber(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"GET /pages/p1/subscribers/s1": func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, 200, Subscriber{ID: "s1", Email: "user@example.com", Mode: "email"})
		},
	})

	sub, err := c.GetSubscriber(context.Background(), "p1", "s1")
	require.NoError(t, err)
	assert.Equal(t, "user@example.com", sub.Email)
}

func TestCreateSubscriber(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"POST /pages/p1/subscribers": func(w http.ResponseWriter, r *http.Request) {
			var req SubscriberRequest
			json.NewDecoder(r.Body).Decode(&req) //nolint:errcheck
			jsonResponse(w, 201, Subscriber{ID: "s-new", Email: req.Subscriber.Email})
		},
	})

	sub, err := c.CreateSubscriber(context.Background(), "p1", SubscriberBody{Email: "new@example.com"})
	require.NoError(t, err)
	assert.Equal(t, "new@example.com", sub.Email)
}

func TestUpdateSubscriber(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"PATCH /pages/p1/subscribers/s1": func(w http.ResponseWriter, r *http.Request) {
			var req SubscriberRequest
			json.NewDecoder(r.Body).Decode(&req) //nolint:errcheck
			jsonResponse(w, 200, Subscriber{ID: "s1", Email: req.Subscriber.Email})
		},
	})

	sub, err := c.UpdateSubscriber(context.Background(), "p1", "s1", SubscriberBody{Email: "updated@example.com"})
	require.NoError(t, err)
	assert.Equal(t, "updated@example.com", sub.Email)
}

func TestDeleteSubscriber(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"DELETE /pages/p1/subscribers/s1": func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(204)
		},
	})

	err := c.DeleteSubscriber(context.Background(), "p1", "s1")
	require.NoError(t, err)
}
