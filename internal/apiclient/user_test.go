package apiclient

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestGetUser(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"GET /organizations/org1/users": func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, 200, []User{
				{ID: "u1", Email: "alice@example.com", FirstName: "Alice"},
				{ID: "u2", Email: "bob@example.com", FirstName: "Bob"},
			})
		},
	})

	user, err := c.GetUser(context.Background(), "org1", "u2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.ID != "u2" {
		t.Errorf("expected ID %q, got %q", "u2", user.ID)
	}
	if user.FirstName != "Bob" {
		t.Errorf("expected FirstName %q, got %q", "Bob", user.FirstName)
	}
}

func TestGetUser_NotFound(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"GET /organizations/org1/users": func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, 200, []User{
				{ID: "u1", Email: "alice@example.com"},
			})
		},
	})

	_, err := c.GetUser(context.Background(), "org1", "nonexistent")
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

func TestCreateUser(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"POST /organizations/org1/users": func(w http.ResponseWriter, r *http.Request) {
			var req UserRequest
			json.NewDecoder(r.Body).Decode(&req) //nolint:errcheck
			jsonResponse(w, 201, User{ID: "u-new", Email: req.User.Email, FirstName: req.User.FirstName})
		},
	})

	user, err := c.CreateUser(context.Background(), "org1", UserBody{Email: "new@example.com", FirstName: "New"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.Email != "new@example.com" {
		t.Errorf("expected Email %q, got %q", "new@example.com", user.Email)
	}
}

func TestDeleteUser(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"DELETE /organizations/org1/users/u1": func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(204)
		},
	})

	err := c.DeleteUser(context.Background(), "org1", "u1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
