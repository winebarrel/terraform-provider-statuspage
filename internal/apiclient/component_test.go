package apiclient

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestGetComponent(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"GET /pages/p1/components/c1": func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, 200, Component{ID: "c1", PageID: "p1", Name: "API", Status: "operational"})
		},
	})

	comp, err := c.GetComponent(context.Background(), "p1", "c1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if comp.ID != "c1" {
		t.Errorf("expected ID %q, got %q", "c1", comp.ID)
	}
	if comp.Name != "API" {
		t.Errorf("expected Name %q, got %q", "API", comp.Name)
	}
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
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if comp.ID != "c-new" {
		t.Errorf("expected ID %q, got %q", "c-new", comp.ID)
	}
	if comp.Name != "Web" {
		t.Errorf("expected Name %q, got %q", "Web", comp.Name)
	}
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
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if comp.Name != "API Updated" {
		t.Errorf("expected Name %q, got %q", "API Updated", comp.Name)
	}
}

func TestDeleteComponent(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"DELETE /pages/p1/components/c1": func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(204)
		},
	})

	err := c.DeleteComponent(context.Background(), "p1", "c1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetComponent_Error(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"GET /pages/p1/components/bad": func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, 404, map[string]string{"error": "not found"})
		},
	})

	_, err := c.GetComponent(context.Background(), "p1", "bad")
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
