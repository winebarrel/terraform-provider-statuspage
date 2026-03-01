package apiclient

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetComponentGroup(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"GET /pages/p1/component-groups/g1": func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, 200, ComponentGroup{ID: "g1", PageID: "p1", Name: "Infrastructure"})
		},
	})

	group, err := c.GetComponentGroup(context.Background(), "p1", "g1")
	require.NoError(t, err)
	assert.Equal(t, "g1", group.ID)
	assert.Equal(t, "Infrastructure", group.Name)
}

func TestCreateComponentGroup(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"POST /pages/p1/component-groups": func(w http.ResponseWriter, r *http.Request) {
			var req ComponentGroupRequest
			json.NewDecoder(r.Body).Decode(&req) //nolint:errcheck
			jsonResponse(w, 201, ComponentGroup{ID: "g-new", PageID: "p1", Name: req.ComponentGroup.Name, Components: req.ComponentGroup.Components})
		},
	})

	group, err := c.CreateComponentGroup(context.Background(), "p1", ComponentGroupBody{Name: "Backend", Components: []string{"c1", "c2"}})
	require.NoError(t, err)
	assert.Equal(t, "Backend", group.Name)
	assert.Len(t, group.Components, 2)
}

func TestUpdateComponentGroup(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"PATCH /pages/p1/component-groups/g1": func(w http.ResponseWriter, r *http.Request) {
			var req ComponentGroupRequest
			json.NewDecoder(r.Body).Decode(&req) //nolint:errcheck
			jsonResponse(w, 200, ComponentGroup{ID: "g1", PageID: "p1", Name: req.ComponentGroup.Name})
		},
	})

	group, err := c.UpdateComponentGroup(context.Background(), "p1", "g1", ComponentGroupBody{Name: "Updated Group"})
	require.NoError(t, err)
	assert.Equal(t, "Updated Group", group.Name)
}

func TestDeleteComponentGroup(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"DELETE /pages/p1/component-groups/g1": func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(204)
		},
	})

	err := c.DeleteComponentGroup(context.Background(), "p1", "g1")
	require.NoError(t, err)
}
