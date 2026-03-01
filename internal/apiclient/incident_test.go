package apiclient

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestGetIncident(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"GET /pages/p1/incidents/i1": func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, 200, Incident{ID: "i1", PageID: "p1", Name: "Major Outage", Status: "investigating"})
		},
	})

	incident, err := c.GetIncident(context.Background(), "p1", "i1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if incident.Name != "Major Outage" {
		t.Errorf("expected Name %q, got %q", "Major Outage", incident.Name)
	}
	if incident.Status != "investigating" {
		t.Errorf("expected Status %q, got %q", "investigating", incident.Status)
	}
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
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if incident.ID != "i-new" {
		t.Errorf("expected ID %q, got %q", "i-new", incident.ID)
	}
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
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if incident.Status != "resolved" {
		t.Errorf("expected Status %q, got %q", "resolved", incident.Status)
	}
}

func TestDeleteIncident(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"DELETE /pages/p1/incidents/i1": func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(204)
		},
	})

	err := c.DeleteIncident(context.Background(), "p1", "i1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
