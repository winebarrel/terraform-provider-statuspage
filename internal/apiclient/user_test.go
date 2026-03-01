package apiclient

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetUser(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"GET /organizations/org1/users": func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, 200, []User{
				{ID: "u1", Email: "alice@example.com", FirstName: "Alice"},
				{ID: "u2", Email: "bob@example.com", FirstName: "Bob"},
			})
		},
	})

	user, err := c.GetUser(context.Background(), "org1", "u2")
	require.NoError(t, err)
	assert.Equal(t, "u2", user.ID)
	assert.Equal(t, "Bob", user.FirstName)
}

func TestGetUser_NotFound(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"GET /organizations/org1/users": func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, 200, []User{
				{ID: "u1", Email: "alice@example.com"},
			})
		},
	})

	_, err := c.GetUser(context.Background(), "org1", "nonexistent")
	require.Error(t, err)
	apiErr, ok := err.(*APIError)
	require.True(t, ok, "expected *APIError, got %T", err)
	assert.Equal(t, 404, apiErr.StatusCode)
}

func TestCreateUser(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"POST /organizations/org1/users": func(w http.ResponseWriter, r *http.Request) {
			var req UserRequest
			json.NewDecoder(r.Body).Decode(&req) //nolint:errcheck
			jsonResponse(w, 201, User{ID: "u-new", Email: req.User.Email, FirstName: req.User.FirstName})
		},
	})

	user, err := c.CreateUser(context.Background(), "org1", UserBody{Email: "new@example.com", FirstName: "New"})
	require.NoError(t, err)
	assert.Equal(t, "new@example.com", user.Email)
}

func TestDeleteUser(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"DELETE /organizations/org1/users/u1": func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(204)
		},
	})

	err := c.DeleteUser(context.Background(), "org1", "u1")
	require.NoError(t, err)
}
