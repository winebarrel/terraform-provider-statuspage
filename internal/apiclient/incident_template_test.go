package apiclient

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetIncidentTemplate(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"GET /pages/p1/incident_templates": func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, 200, []IncidentTemplate{
				{ID: "t1", Name: "Template A"},
				{ID: "t2", Name: "Template B"},
			})
		},
	})

	tmpl, err := c.GetIncidentTemplate(context.Background(), "p1", "t2")
	require.NoError(t, err)
	assert.Equal(t, "t2", tmpl.ID)
	assert.Equal(t, "Template B", tmpl.Name)
}

func TestGetIncidentTemplate_NotFound(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"GET /pages/p1/incident_templates": func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, 200, []IncidentTemplate{
				{ID: "t1", Name: "Template A"},
			})
		},
	})

	_, err := c.GetIncidentTemplate(context.Background(), "p1", "nonexistent")
	require.Error(t, err)
	apiErr, ok := err.(*APIError)
	require.True(t, ok, "expected *APIError, got %T", err)
	assert.Equal(t, 404, apiErr.StatusCode)
}

func TestCreateIncidentTemplate(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"POST /pages/p1/incident_templates": func(w http.ResponseWriter, r *http.Request) {
			var req IncidentTemplateRequest
			json.NewDecoder(r.Body).Decode(&req) //nolint:errcheck
			jsonResponse(w, 201, IncidentTemplate{ID: "t-new", Name: req.Template.Name, Title: req.Template.Title})
		},
	})

	tmpl, err := c.CreateIncidentTemplate(context.Background(), "p1", IncidentTemplateBody{Name: "Outage Template", Title: "Service Outage"})
	require.NoError(t, err)
	assert.Equal(t, "t-new", tmpl.ID)
	assert.Equal(t, "Outage Template", tmpl.Name)
}

func TestUpdateIncidentTemplate(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"PATCH /pages/p1/incident_templates/t1": func(w http.ResponseWriter, r *http.Request) {
			var req IncidentTemplateRequest
			json.NewDecoder(r.Body).Decode(&req) //nolint:errcheck
			jsonResponse(w, 200, IncidentTemplate{ID: "t1", Name: req.Template.Name})
		},
	})

	tmpl, err := c.UpdateIncidentTemplate(context.Background(), "p1", "t1", IncidentTemplateBody{Name: "Updated Template"})
	require.NoError(t, err)
	assert.Equal(t, "Updated Template", tmpl.Name)
}

func TestDeleteIncidentTemplate(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"DELETE /pages/p1/incident_templates/t1": func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(204)
		},
	})

	err := c.DeleteIncidentTemplate(context.Background(), "p1", "t1")
	require.NoError(t, err)
}
