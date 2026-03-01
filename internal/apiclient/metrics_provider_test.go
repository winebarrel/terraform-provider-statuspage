package apiclient

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestGetMetricsProvider(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"GET /pages/p1/metrics_providers/mp1": func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, 200, MetricsProvider{ID: "mp1", Type: "Datadog"})
		},
	})

	provider, err := c.GetMetricsProvider(context.Background(), "p1", "mp1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if provider.ID != "mp1" {
		t.Errorf("expected ID %q, got %q", "mp1", provider.ID)
	}
	if provider.Type != "Datadog" {
		t.Errorf("expected Type %q, got %q", "Datadog", provider.Type)
	}
}

func TestCreateMetricsProvider(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"POST /pages/p1/metrics_providers": func(w http.ResponseWriter, r *http.Request) {
			var req MetricsProviderRequest
			json.NewDecoder(r.Body).Decode(&req) //nolint:errcheck
			jsonResponse(w, 201, MetricsProvider{ID: "mp-new", Type: req.MetricsProvider.Type, Email: req.MetricsProvider.Email})
		},
	})

	provider, err := c.CreateMetricsProvider(context.Background(), "p1", MetricsProviderBody{Type: "Datadog", Email: "test@example.com"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if provider.Type != "Datadog" {
		t.Errorf("expected Type %q, got %q", "Datadog", provider.Type)
	}
}

func TestUpdateMetricsProvider(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"PATCH /pages/p1/metrics_providers/mp1": func(w http.ResponseWriter, r *http.Request) {
			var req MetricsProviderRequest
			json.NewDecoder(r.Body).Decode(&req) //nolint:errcheck
			jsonResponse(w, 200, MetricsProvider{ID: "mp1", Type: req.MetricsProvider.Type})
		},
	})

	provider, err := c.UpdateMetricsProvider(context.Background(), "p1", "mp1", MetricsProviderBody{Type: "NewRelic"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if provider.Type != "NewRelic" {
		t.Errorf("expected Type %q, got %q", "NewRelic", provider.Type)
	}
}

func TestDeleteMetricsProvider(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"DELETE /pages/p1/metrics_providers/mp1": func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(204)
		},
	})

	err := c.DeleteMetricsProvider(context.Background(), "p1", "mp1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
