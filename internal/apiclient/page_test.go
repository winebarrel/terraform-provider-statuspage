package apiclient

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetPage(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"GET /pages/page1": func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, 200, Page{ID: "page1", Name: "My Page", Subdomain: "mypage"})
		},
	})

	page, err := c.GetPage(context.Background(), "page1")
	require.NoError(t, err)
	assert.Equal(t, "page1", page.ID)
	assert.Equal(t, "My Page", page.Name)
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
	require.NoError(t, err)
	assert.Equal(t, "Updated Page", page.Name)
}
