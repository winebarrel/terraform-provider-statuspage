package apiclient

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetPostmortem(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"GET /pages/p1/incidents/i1/postmortem": func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, 200, Postmortem{Body: "Root cause analysis", BodyDraft: "Draft content"})
		},
	})

	pm, err := c.GetPostmortem(context.Background(), "p1", "i1")
	require.NoError(t, err)
	assert.Equal(t, "Root cause analysis", pm.Body)
}

func TestCreateOrUpdatePostmortem(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"PUT /pages/p1/incidents/i1/postmortem": func(w http.ResponseWriter, r *http.Request) {
			var req PostmortemRequest
			json.NewDecoder(r.Body).Decode(&req) //nolint:errcheck
			jsonResponse(w, 200, Postmortem{Body: req.Postmortem.Body})
		},
	})

	pm, err := c.CreateOrUpdatePostmortem(context.Background(), "p1", "i1", PostmortemBody{Body: "Updated analysis"})
	require.NoError(t, err)
	assert.Equal(t, "Updated analysis", pm.Body)
}

func TestDeletePostmortem(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"DELETE /pages/p1/incidents/i1/postmortem": func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(204)
		},
	})

	err := c.DeletePostmortem(context.Background(), "p1", "i1")
	require.NoError(t, err)
}
