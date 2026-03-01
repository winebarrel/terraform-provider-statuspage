package apiclient

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetPermissions(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"GET /organizations/org1/permissions/u1": func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, 200, Permissions{UserID: "u1", Pages: map[string]string{"p1": "admin", "p2": "viewer"}})
		},
	})

	perms, err := c.GetPermissions(context.Background(), "org1", "u1")
	require.NoError(t, err)
	assert.Equal(t, "u1", perms.UserID)
	assert.Equal(t, "admin", perms.Pages["p1"])
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
	require.NoError(t, err)
	assert.Equal(t, "manager", perms.Pages["p1"])
}
