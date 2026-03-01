package apiclient

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestGetMetric(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"GET /pages/p1/metrics/m1": func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, 200, Metric{ID: "m1", Name: "Latency", Suffix: "ms"})
		},
	})

	metric, err := c.GetMetric(context.Background(), "p1", "m1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if metric.ID != "m1" {
		t.Errorf("expected ID %q, got %q", "m1", metric.ID)
	}
	if metric.Name != "Latency" {
		t.Errorf("expected Name %q, got %q", "Latency", metric.Name)
	}
}

func TestCreateMetric(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"POST /pages/p1/metrics_providers/mp1/metrics": func(w http.ResponseWriter, r *http.Request) {
			var req MetricRequest
			json.NewDecoder(r.Body).Decode(&req) //nolint:errcheck
			jsonResponse(w, 201, Metric{ID: "m-new", Name: req.Metric.Name, MetricsProviderID: "mp1"})
		},
	})

	metric, err := c.CreateMetric(context.Background(), "p1", "mp1", MetricBody{Name: "CPU Usage"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if metric.Name != "CPU Usage" {
		t.Errorf("expected Name %q, got %q", "CPU Usage", metric.Name)
	}
}

func TestUpdateMetric(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"PATCH /pages/p1/metrics/m1": func(w http.ResponseWriter, r *http.Request) {
			var req MetricRequest
			json.NewDecoder(r.Body).Decode(&req) //nolint:errcheck
			jsonResponse(w, 200, Metric{ID: "m1", Name: req.Metric.Name})
		},
	})

	metric, err := c.UpdateMetric(context.Background(), "p1", "m1", MetricBody{Name: "Updated Metric"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if metric.Name != "Updated Metric" {
		t.Errorf("expected Name %q, got %q", "Updated Metric", metric.Name)
	}
}

func TestDeleteMetric(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"DELETE /pages/p1/metrics/m1": func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(204)
		},
	})

	err := c.DeleteMetric(context.Background(), "p1", "m1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
