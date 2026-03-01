package apiclient

import (
	"context"
	"net/http"
	"testing"
)

func TestGetPermissions(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"GET /organizations/org1/permissions/u1": func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, 200, Permissions{UserID: "u1", Pages: map[string]string{"p1": "admin", "p2": "viewer"}})
		},
	})

	perms, err := c.GetPermissions(context.Background(), "org1", "u1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if perms.UserID != "u1" {
		t.Errorf("expected UserID %q, got %q", "u1", perms.UserID)
	}
	if perms.Pages["p1"] != "admin" {
		t.Errorf("expected page p1 role %q, got %q", "admin", perms.Pages["p1"])
	}
}

func TestUpdatePermissions(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"PUT /organizations/org1/permissions/u1": func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, 200, Permissions{UserID: "u1", Pages: map[string]string{"p1": "manager"}})
		},
	})

	perms, err := c.UpdatePermissions(context.Background(), "org1", "u1", PermissionsRequest{
		Pages: map[string][]string{"p1": {"manager"}},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if perms.Pages["p1"] != "manager" {
		t.Errorf("expected page p1 role %q, got %q", "manager", perms.Pages["p1"])
	}
}
