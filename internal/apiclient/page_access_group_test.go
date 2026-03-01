package apiclient

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestGetPageAccessGroup(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"GET /pages/p1/page_access_groups/ag1": func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, 200, PageAccessGroup{ID: "ag1", Name: "VIP Group"})
		},
	})

	group, err := c.GetPageAccessGroup(context.Background(), "p1", "ag1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if group.Name != "VIP Group" {
		t.Errorf("expected Name %q, got %q", "VIP Group", group.Name)
	}
}

func TestCreatePageAccessGroup(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"POST /pages/p1/page_access_groups": func(w http.ResponseWriter, r *http.Request) {
			var req PageAccessGroupRequest
			json.NewDecoder(r.Body).Decode(&req) //nolint:errcheck
			jsonResponse(w, 201, PageAccessGroup{ID: "ag-new", Name: req.PageAccessGroup.Name})
		},
	})

	group, err := c.CreatePageAccessGroup(context.Background(), "p1", PageAccessGroupBody{Name: "New Group"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if group.Name != "New Group" {
		t.Errorf("expected Name %q, got %q", "New Group", group.Name)
	}
}

func TestUpdatePageAccessGroup(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"PATCH /pages/p1/page_access_groups/ag1": func(w http.ResponseWriter, r *http.Request) {
			var req PageAccessGroupRequest
			json.NewDecoder(r.Body).Decode(&req) //nolint:errcheck
			jsonResponse(w, 200, PageAccessGroup{ID: "ag1", Name: req.PageAccessGroup.Name})
		},
	})

	group, err := c.UpdatePageAccessGroup(context.Background(), "p1", "ag1", PageAccessGroupBody{Name: "Updated Group"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if group.Name != "Updated Group" {
		t.Errorf("expected Name %q, got %q", "Updated Group", group.Name)
	}
}

func TestDeletePageAccessGroup(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"DELETE /pages/p1/page_access_groups/ag1": func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(204)
		},
	})

	err := c.DeletePageAccessGroup(context.Background(), "p1", "ag1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
