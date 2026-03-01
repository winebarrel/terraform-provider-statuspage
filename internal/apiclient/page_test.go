package apiclient

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestGetPage(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"GET /pages/page1": func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, 200, Page{ID: "page1", Name: "My Page", Subdomain: "mypage"})
		},
	})

	page, err := c.GetPage(context.Background(), "page1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if page.ID != "page1" {
		t.Errorf("expected ID %q, got %q", "page1", page.ID)
	}
	if page.Name != "My Page" {
		t.Errorf("expected Name %q, got %q", "My Page", page.Name)
	}
}

func TestUpdatePage(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"PATCH /pages/page1": func(w http.ResponseWriter, r *http.Request) {
			var req PageRequest
			json.NewDecoder(r.Body).Decode(&req) //nolint:errcheck
			jsonResponse(w, 200, Page{ID: "page1", Name: req.Page.Name})
		},
	})

	page, err := c.UpdatePage(context.Background(), "page1", PageBody{Name: "Updated Page"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if page.Name != "Updated Page" {
		t.Errorf("expected Name %q, got %q", "Updated Page", page.Name)
	}
}
