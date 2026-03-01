package apiclient

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestGetPageAccessUser(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"GET /pages/p1/page_access_users/au1": func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, 200, PageAccessUser{ID: "au1", ExternalEmail: "user@example.com"})
		},
	})

	user, err := c.GetPageAccessUser(context.Background(), "p1", "au1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.ExternalEmail != "user@example.com" {
		t.Errorf("expected ExternalEmail %q, got %q", "user@example.com", user.ExternalEmail)
	}
}

func TestCreatePageAccessUser(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"POST /pages/p1/page_access_users": func(w http.ResponseWriter, r *http.Request) {
			var req PageAccessUserRequest
			json.NewDecoder(r.Body).Decode(&req) //nolint:errcheck
			jsonResponse(w, 201, PageAccessUser{ID: "au-new", ExternalEmail: req.PageAccessUser.ExternalEmail})
		},
	})

	user, err := c.CreatePageAccessUser(context.Background(), "p1", PageAccessUserBody{ExternalEmail: "new@example.com"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.ExternalEmail != "new@example.com" {
		t.Errorf("expected ExternalEmail %q, got %q", "new@example.com", user.ExternalEmail)
	}
}

func TestUpdatePageAccessUser(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"PATCH /pages/p1/page_access_users/au1": func(w http.ResponseWriter, r *http.Request) {
			var req PageAccessUserRequest
			json.NewDecoder(r.Body).Decode(&req) //nolint:errcheck
			jsonResponse(w, 200, PageAccessUser{ID: "au1", ExternalEmail: req.PageAccessUser.ExternalEmail})
		},
	})

	user, err := c.UpdatePageAccessUser(context.Background(), "p1", "au1", PageAccessUserBody{ExternalEmail: "updated@example.com"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.ExternalEmail != "updated@example.com" {
		t.Errorf("expected ExternalEmail %q, got %q", "updated@example.com", user.ExternalEmail)
	}
}

func TestDeletePageAccessUser(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"DELETE /pages/p1/page_access_users/au1": func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(204)
		},
	})

	err := c.DeletePageAccessUser(context.Background(), "p1", "au1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
