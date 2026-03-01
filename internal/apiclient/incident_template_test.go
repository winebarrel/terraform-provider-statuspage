package apiclient

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
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
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tmpl.ID != "t2" {
		t.Errorf("expected ID %q, got %q", "t2", tmpl.ID)
	}
	if tmpl.Name != "Template B" {
		t.Errorf("expected Name %q, got %q", "Template B", tmpl.Name)
	}
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
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != 404 {
		t.Errorf("expected status 404, got %d", apiErr.StatusCode)
	}
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
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tmpl.ID != "t-new" {
		t.Errorf("expected ID %q, got %q", "t-new", tmpl.ID)
	}
	if tmpl.Name != "Outage Template" {
		t.Errorf("expected Name %q, got %q", "Outage Template", tmpl.Name)
	}
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
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tmpl.Name != "Updated Template" {
		t.Errorf("expected Name %q, got %q", "Updated Template", tmpl.Name)
	}
}

func TestDeleteIncidentTemplate(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"DELETE /pages/p1/incident_templates/t1": func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(204)
		},
	})

	err := c.DeleteIncidentTemplate(context.Background(), "p1", "t1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
