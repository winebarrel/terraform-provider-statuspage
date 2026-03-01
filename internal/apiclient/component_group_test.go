package apiclient

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestGetComponentGroup(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"GET /pages/p1/component-groups/g1": func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, 200, ComponentGroup{ID: "g1", PageID: "p1", Name: "Infrastructure"})
		},
	})

	group, err := c.GetComponentGroup(context.Background(), "p1", "g1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if group.ID != "g1" {
		t.Errorf("expected ID %q, got %q", "g1", group.ID)
	}
	if group.Name != "Infrastructure" {
		t.Errorf("expected Name %q, got %q", "Infrastructure", group.Name)
	}
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
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if group.Name != "Backend" {
		t.Errorf("expected Name %q, got %q", "Backend", group.Name)
	}
	if len(group.Components) != 2 {
		t.Errorf("expected 2 components, got %d", len(group.Components))
	}
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
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if group.Name != "Updated Group" {
		t.Errorf("expected Name %q, got %q", "Updated Group", group.Name)
	}
}

func TestDeleteComponentGroup(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"DELETE /pages/p1/component-groups/g1": func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(204)
		},
	})

	err := c.DeleteComponentGroup(context.Background(), "p1", "g1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
