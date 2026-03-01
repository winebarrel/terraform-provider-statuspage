package apiclient

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetStatusEmbedConfig(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"GET /pages/p1/status_embed_config": func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, 200, StatusEmbedConfig{PageID: "p1", Position: "bottom_left"})
		},
	})

	config, err := c.GetStatusEmbedConfig(context.Background(), "p1")
	require.NoError(t, err)
	assert.Equal(t, "bottom_left", config.Position)
}

func TestUpdateStatusEmbedConfig(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"PATCH /pages/p1/status_embed_config": func(w http.ResponseWriter, r *http.Request) {
			var req StatusEmbedConfigRequest
			json.NewDecoder(r.Body).Decode(&req) //nolint:errcheck
			jsonResponse(w, 200, StatusEmbedConfig{PageID: "p1", Position: req.StatusEmbedConfig.Position})
		},
	})

	config, err := c.UpdateStatusEmbedConfig(context.Background(), "p1", StatusEmbedConfigBody{Position: "top_right"})
	require.NoError(t, err)
	assert.Equal(t, "top_right", config.Position)
}
