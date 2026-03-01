package apiclient

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// testServer creates a test HTTP server and returns a client configured to use it.
// The handler map keys are "METHOD /path" strings.
func testServer(t *testing.T, handlers map[string]http.HandlerFunc) *Client {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Method + " " + r.URL.Path
		if h, ok := handlers[key]; ok {
			h(w, r)
			return
		}
		t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		w.WriteHeader(http.StatusNotFound)
	}))
	t.Cleanup(server.Close)
	c := NewClient("test-key", WithRateLimitInterval(0))
	c.baseURL = server.URL
	return c
}

func jsonResponse(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v) //nolint:errcheck
}
