package apiclient

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestGetPostmortem(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"GET /pages/p1/incidents/i1/postmortem": func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, 200, Postmortem{Body: "Root cause analysis", BodyDraft: "Draft content"})
		},
	})

	pm, err := c.GetPostmortem(context.Background(), "p1", "i1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pm.Body != "Root cause analysis" {
		t.Errorf("expected Body %q, got %q", "Root cause analysis", pm.Body)
	}
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
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pm.Body != "Updated analysis" {
		t.Errorf("expected Body %q, got %q", "Updated analysis", pm.Body)
	}
}

func TestDeletePostmortem(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"DELETE /pages/p1/incidents/i1/postmortem": func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(204)
		},
	})

	err := c.DeletePostmortem(context.Background(), "p1", "i1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
