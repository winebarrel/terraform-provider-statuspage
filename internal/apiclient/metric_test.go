package apiclient

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetMetric(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"GET /pages/p1/metrics/m1": func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, 200, Metric{ID: "m1", Name: "Latency", Suffix: "ms"})
		},
	})

	metric, err := c.GetMetric(context.Background(), "p1", "m1")
	require.NoError(t, err)
	assert.Equal(t, "m1", metric.ID)
	assert.Equal(t, "Latency", metric.Name)
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
	require.NoError(t, err)
	assert.Equal(t, "CPU Usage", metric.Name)
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
	require.NoError(t, err)
	assert.Equal(t, "Updated Metric", metric.Name)
}

func TestDeleteMetric(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"DELETE /pages/p1/metrics/m1": func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(204)
		},
	})

	err := c.DeleteMetric(context.Background(), "p1", "m1")
	require.NoError(t, err)
}
