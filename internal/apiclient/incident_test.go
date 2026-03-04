package apiclient

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetIncident(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"GET /pages/p1/incidents/i1": func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, 200, Incident{ID: "i1", PageID: "p1", Name: "Major Outage", Status: "investigating"})
		},
	})

	incident, err := c.GetIncident(context.Background(), "p1", "i1")
	require.NoError(t, err)
	assert.Equal(t, "Major Outage", incident.Name)
	assert.Equal(t, "investigating", incident.Status)
}

func TestGetIncident_withComponents(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"GET /pages/p1/incidents/i1": func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, 200, Incident{
				ID:     "i1",
				PageID: "p1",
				Name:   "Major Outage",
				Status: "investigating",
				Components: []Component{
					{ID: "c1", PageID: "p1", Name: "API", Status: "major_outage"},
					{ID: "c2", PageID: "p1", Name: "Web App", Status: "degraded_performance"},
				},
				ComponentIDs: []string{"c1", "c2"},
			})
		},
	})

	incident, err := c.GetIncident(context.Background(), "p1", "i1")
	require.NoError(t, err)
	assert.Equal(t, "Major Outage", incident.Name)
	require.Len(t, incident.Components, 2)
	assert.Equal(t, "c1", incident.Components[0].ID)
	assert.Equal(t, "API", incident.Components[0].Name)
	assert.Equal(t, "major_outage", incident.Components[0].Status)
	assert.Equal(t, "c2", incident.Components[1].ID)
	assert.Equal(t, "Web App", incident.Components[1].Name)
	assert.Equal(t, "degraded_performance", incident.Components[1].Status)
	assert.Equal(t, []string{"c1", "c2"}, incident.ComponentIDs)
}

func TestCreateIncident(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"POST /pages/p1/incidents": func(w http.ResponseWriter, r *http.Request) {
			var req IncidentRequest
			json.NewDecoder(r.Body).Decode(&req) //nolint:errcheck
			jsonResponse(w, 201, Incident{ID: "i-new", PageID: "p1", Name: req.Incident.Name, Status: req.Incident.Status})
		},
	})

	incident, err := c.CreateIncident(context.Background(), "p1", IncidentBody{Name: "API Down", Status: "investigating"})
	require.NoError(t, err)
	assert.Equal(t, "i-new", incident.ID)
}

func TestUpdateIncident(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"PATCH /pages/p1/incidents/i1": func(w http.ResponseWriter, r *http.Request) {
			var req IncidentRequest
			json.NewDecoder(r.Body).Decode(&req) //nolint:errcheck
			jsonResponse(w, 200, Incident{ID: "i1", PageID: "p1", Name: "API Down", Status: req.Incident.Status})
		},
	})

	incident, err := c.UpdateIncident(context.Background(), "p1", "i1", IncidentBody{Status: "resolved"})
	require.NoError(t, err)
	assert.Equal(t, "resolved", incident.Status)
}

func TestDeleteIncident(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"DELETE /pages/p1/incidents/i1": func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(204)
		},
	})

	err := c.DeleteIncident(context.Background(), "p1", "i1")
	require.NoError(t, err)
}
