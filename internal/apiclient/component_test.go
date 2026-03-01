package apiclient

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetComponent(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"GET /pages/p1/components/c1": func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, 200, Component{ID: "c1", PageID: "p1", Name: "API", Status: "operational"})
		},
	})

	comp, err := c.GetComponent(context.Background(), "p1", "c1")
	require.NoError(t, err)
	assert.Equal(t, "c1", comp.ID)
	assert.Equal(t, "API", comp.Name)
}

func TestCreateComponent(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"POST /pages/p1/components": func(w http.ResponseWriter, r *http.Request) {
			var req ComponentRequest
			json.NewDecoder(r.Body).Decode(&req) //nolint:errcheck
			jsonResponse(w, 201, Component{ID: "c-new", PageID: "p1", Name: req.Component.Name, Status: req.Component.Status})
		},
	})

	comp, err := c.CreateComponent(context.Background(), "p1", ComponentBody{Name: "Web", Status: "operational"})
	require.NoError(t, err)
	assert.Equal(t, "c-new", comp.ID)
	assert.Equal(t, "Web", comp.Name)
}

func TestUpdateComponent(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"PATCH /pages/p1/components/c1": func(w http.ResponseWriter, r *http.Request) {
			var req ComponentRequest
			json.NewDecoder(r.Body).Decode(&req) //nolint:errcheck
			jsonResponse(w, 200, Component{ID: "c1", PageID: "p1", Name: req.Component.Name, Status: "degraded_performance"})
		},
	})

	comp, err := c.UpdateComponent(context.Background(), "p1", "c1", ComponentBody{Name: "API Updated"})
	require.NoError(t, err)
	assert.Equal(t, "API Updated", comp.Name)
}

func TestDeleteComponent(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"DELETE /pages/p1/components/c1": func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(204)
		},
	})

	err := c.DeleteComponent(context.Background(), "p1", "c1")
	require.NoError(t, err)
}

func TestGetComponent_Error(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"GET /pages/p1/components/bad": func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, 404, map[string]string{"error": "not found"})
		},
	})

	_, err := c.GetComponent(context.Background(), "p1", "bad")
	require.Error(t, err)
	apiErr, ok := err.(*APIError)
	require.True(t, ok, "expected *APIError, got %T", err)
	assert.Equal(t, 404, apiErr.StatusCode)
}
