package apiclient

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

func newTestClient(t *testing.T, handler http.Handler) *Client {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	c := NewClient("test-api-key", WithRateLimitInterval(0))
	c.baseURL = server.URL
	return c
}

func TestNewClient(t *testing.T) {
	c := NewClient("my-key")
	if c.apiKey != "my-key" {
		t.Errorf("expected apiKey %q, got %q", "my-key", c.apiKey)
	}
	if c.baseURL != DefaultBaseURL {
		t.Errorf("expected baseURL %q, got %q", DefaultBaseURL, c.baseURL)
	}
	if c.rateLimitInterval != time.Second {
		t.Errorf("expected rateLimitInterval %v, got %v", time.Second, c.rateLimitInterval)
	}
	if c.httpClient == nil {
		t.Error("expected httpClient to be non-nil")
	}
}

func TestNewClient_WithRateLimitInterval(t *testing.T) {
	c := NewClient("my-key", WithRateLimitInterval(500*time.Millisecond))
	if c.rateLimitInterval != 500*time.Millisecond {
		t.Errorf("expected rateLimitInterval %v, got %v", 500*time.Millisecond, c.rateLimitInterval)
	}
}

func TestAPIError_Error(t *testing.T) {
	err := &APIError{StatusCode: 404, Message: "not found"}
	expected := "statuspage API error (HTTP 404): not found"
	if err.Error() != expected {
		t.Errorf("expected %q, got %q", expected, err.Error())
	}
}

func TestClient_AuthorizationHeader(t *testing.T) {
	var gotAuth string
	c := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"id": "1"}) //nolint:errcheck
	}))

	var result map[string]string
	_ = c.Get(context.Background(), "/test", &result)

	if gotAuth != "OAuth test-api-key" {
		t.Errorf("expected Authorization %q, got %q", "OAuth test-api-key", gotAuth)
	}
}

func TestClient_ContentTypeHeader(t *testing.T) {
	var gotContentType string
	c := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotContentType = r.Header.Get("Content-Type")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"id": "1"}) //nolint:errcheck
	}))

	body := map[string]string{"name": "test"}
	var result map[string]string
	_ = c.Post(context.Background(), "/test", body, &result)

	if gotContentType != "application/json" {
		t.Errorf("expected Content-Type %q, got %q", "application/json", gotContentType)
	}
}

func TestClient_NoContentTypeOnGetRequest(t *testing.T) {
	var gotContentType string
	c := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotContentType = r.Header.Get("Content-Type")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"id": "1"}) //nolint:errcheck
	}))

	var result map[string]string
	_ = c.Get(context.Background(), "/test", &result)

	if gotContentType != "" {
		t.Errorf("expected empty Content-Type for GET, got %q", gotContentType)
	}
}

func TestClient_Get(t *testing.T) {
	c := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}
		if r.URL.Path != "/pages/p1" {
			t.Errorf("expected path /pages/p1, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"id": "p1", "name": "My Page"}) //nolint:errcheck
	}))

	var result map[string]string
	err := c.Get(context.Background(), "/pages/p1", &result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["id"] != "p1" {
		t.Errorf("expected id %q, got %q", "p1", result["id"])
	}
	if result["name"] != "My Page" {
		t.Errorf("expected name %q, got %q", "My Page", result["name"])
	}
}

func TestClient_Post(t *testing.T) {
	var gotBody map[string]string
	c := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}
		json.NewDecoder(r.Body).Decode(&gotBody) //nolint:errcheck
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"id": "new-1", "name": gotBody["name"]}) //nolint:errcheck
	}))

	body := map[string]string{"name": "new-component"}
	var result map[string]string
	err := c.Post(context.Background(), "/components", body, &result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["id"] != "new-1" {
		t.Errorf("expected id %q, got %q", "new-1", result["id"])
	}
}

func TestClient_Post_NilResult(t *testing.T) {
	c := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	}))

	body := map[string]string{"name": "test"}
	err := c.Post(context.Background(), "/test", body, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestClient_Patch(t *testing.T) {
	c := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected method PATCH, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"id": "1", "name": "updated"}) //nolint:errcheck
	}))

	body := map[string]string{"name": "updated"}
	var result map[string]string
	err := c.Patch(context.Background(), "/test/1", body, &result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["name"] != "updated" {
		t.Errorf("expected name %q, got %q", "updated", result["name"])
	}
}

func TestClient_Put(t *testing.T) {
	c := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected method PUT, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"id": "1"}) //nolint:errcheck
	}))

	body := map[string]string{"name": "test"}
	var result map[string]string
	err := c.Put(context.Background(), "/test/1", body, &result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["id"] != "1" {
		t.Errorf("expected id %q, got %q", "1", result["id"])
	}
}

func TestClient_Delete(t *testing.T) {
	c := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected method DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))

	err := c.Delete(context.Background(), "/test/1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestClient_ErrorResponse_WithErrorField(t *testing.T) {
	c := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(map[string]string{"error": "validation failed"}) //nolint:errcheck
	}))

	var result map[string]string
	err := c.Get(context.Background(), "/test", &result)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != 422 {
		t.Errorf("expected status 422, got %d", apiErr.StatusCode)
	}
	if apiErr.Message != "validation failed" {
		t.Errorf("expected message %q, got %q", "validation failed", apiErr.Message)
	}
}

func TestClient_ErrorResponse_WithMessageField(t *testing.T) {
	c := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{"message": "forbidden"}) //nolint:errcheck
	}))

	var result map[string]string
	err := c.Get(context.Background(), "/test", &result)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	apiErr := err.(*APIError)
	if apiErr.Message != "forbidden" {
		t.Errorf("expected message %q, got %q", "forbidden", apiErr.Message)
	}
}

func TestClient_ErrorResponse_RawBody(t *testing.T) {
	c := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "internal error") //nolint:errcheck
	}))

	var result map[string]string
	err := c.Get(context.Background(), "/test", &result)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	apiErr := err.(*APIError)
	if apiErr.StatusCode != 500 {
		t.Errorf("expected status 500, got %d", apiErr.StatusCode)
	}
	if apiErr.Message != "internal error" {
		t.Errorf("expected message %q, got %q", "internal error", apiErr.Message)
	}
}

func TestClient_RateLimitRetry_429(t *testing.T) {
	var callCount int32
	c := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&callCount, 1)
		if n == 1 {
			w.WriteHeader(429)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"id": "1"}) //nolint:errcheck
	}))

	var result map[string]string
	err := c.Get(context.Background(), "/test", &result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if atomic.LoadInt32(&callCount) != 2 {
		t.Errorf("expected 2 requests (1 retry), got %d", callCount)
	}
}

func TestClient_RateLimitRetry_420(t *testing.T) {
	var callCount int32
	c := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&callCount, 1)
		if n == 1 {
			w.WriteHeader(420)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"id": "1"}) //nolint:errcheck
	}))

	var result map[string]string
	err := c.Get(context.Background(), "/test", &result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if atomic.LoadInt32(&callCount) != 2 {
		t.Errorf("expected 2 requests (1 retry), got %d", callCount)
	}
}

func TestClient_RateLimit(t *testing.T) {
	c := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"id": "1"}) //nolint:errcheck
	}))
	c.rateLimitInterval = 50 * time.Millisecond

	var result map[string]string
	start := time.Now()
	_ = c.Get(context.Background(), "/test", &result)
	_ = c.Get(context.Background(), "/test", &result)
	elapsed := time.Since(start)

	if elapsed < 50*time.Millisecond {
		t.Errorf("expected at least 50ms for rate limiting, got %v", elapsed)
	}
}

func TestClient_RequestBodySerialization(t *testing.T) {
	var gotBody []byte
	c := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotBody, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"id": "1"}) //nolint:errcheck
	}))

	body := ComponentRequest{
		Component: ComponentBody{
			Name:   "web",
			Status: "operational",
		},
	}
	var result map[string]string
	_ = c.Post(context.Background(), "/test", body, &result)

	var parsed ComponentRequest
	if err := json.Unmarshal(gotBody, &parsed); err != nil {
		t.Fatalf("failed to parse request body: %v", err)
	}
	if parsed.Component.Name != "web" {
		t.Errorf("expected component name %q, got %q", "web", parsed.Component.Name)
	}
	if parsed.Component.Status != "operational" {
		t.Errorf("expected component status %q, got %q", "operational", parsed.Component.Status)
	}
}
