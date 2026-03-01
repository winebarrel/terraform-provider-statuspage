package apiclient

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestGetStatusEmbedConfig(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"GET /pages/p1/status_embed_config": func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, 200, StatusEmbedConfig{PageID: "p1", Position: "bottom_left"})
		},
	})

	config, err := c.GetStatusEmbedConfig(context.Background(), "p1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if config.Position != "bottom_left" {
		t.Errorf("expected Position %q, got %q", "bottom_left", config.Position)
	}
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
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if config.Position != "top_right" {
		t.Errorf("expected Position %q, got %q", "top_right", config.Position)
	}
}
