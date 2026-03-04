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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	assert.Equal(t, "my-key", c.apiKey)
	assert.Equal(t, DefaultBaseURL, c.baseURL)
	assert.Equal(t, time.Second, c.rateLimitInterval)
	assert.NotNil(t, c.httpClient)
}

func TestNewClient_WithRateLimitInterval(t *testing.T) {
	c := NewClient("my-key", WithRateLimitInterval(500*time.Millisecond))
	assert.Equal(t, 500*time.Millisecond, c.rateLimitInterval)
}

func TestAPIError_Error(t *testing.T) {
	err := &APIError{StatusCode: 404, Message: "not found"}
	assert.Equal(t, "statuspage API error (HTTP 404): not found", err.Error())
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

	assert.Equal(t, "OAuth test-api-key", gotAuth)
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

	assert.Equal(t, "application/json", gotContentType)
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

	assert.Empty(t, gotContentType)
}

func TestClient_Get(t *testing.T) {
	c := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/pages/p1", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"id": "p1", "name": "My Page"}) //nolint:errcheck
	}))

	var result map[string]string
	err := c.Get(context.Background(), "/pages/p1", &result)
	require.NoError(t, err)
	assert.Equal(t, "p1", result["id"])
	assert.Equal(t, "My Page", result["name"])
}

func TestClient_Post(t *testing.T) {
	var gotBody map[string]string
	c := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		json.NewDecoder(r.Body).Decode(&gotBody) //nolint:errcheck
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"id": "new-1", "name": gotBody["name"]}) //nolint:errcheck
	}))

	body := map[string]string{"name": "new-component"}
	var result map[string]string
	err := c.Post(context.Background(), "/components", body, &result)
	require.NoError(t, err)
	assert.Equal(t, "new-1", result["id"])
}

func TestClient_Post_NilResult(t *testing.T) {
	c := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	}))

	body := map[string]string{"name": "test"}
	err := c.Post(context.Background(), "/test", body, nil)
	require.NoError(t, err)
}

func TestClient_Patch(t *testing.T) {
	c := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPatch, r.Method)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"id": "1", "name": "updated"}) //nolint:errcheck
	}))

	body := map[string]string{"name": "updated"}
	var result map[string]string
	err := c.Patch(context.Background(), "/test/1", body, &result)
	require.NoError(t, err)
	assert.Equal(t, "updated", result["name"])
}

func TestClient_Put(t *testing.T) {
	c := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"id": "1"}) //nolint:errcheck
	}))

	body := map[string]string{"name": "test"}
	var result map[string]string
	err := c.Put(context.Background(), "/test/1", body, &result)
	require.NoError(t, err)
	assert.Equal(t, "1", result["id"])
}

func TestClient_Delete(t *testing.T) {
	c := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		w.WriteHeader(http.StatusNoContent)
	}))

	err := c.Delete(context.Background(), "/test/1")
	require.NoError(t, err)
}

func TestClient_ErrorResponse_WithErrorField(t *testing.T) {
	c := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(map[string]string{"error": "validation failed"}) //nolint:errcheck
	}))

	var result map[string]string
	err := c.Get(context.Background(), "/test", &result)
	require.Error(t, err)
	apiErr, ok := err.(*APIError)
	require.True(t, ok, "expected *APIError, got %T", err)
	assert.Equal(t, 422, apiErr.StatusCode)
	assert.Equal(t, "validation failed", apiErr.Message)
}

func TestClient_ErrorResponse_WithMessageField(t *testing.T) {
	c := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{"message": "forbidden"}) //nolint:errcheck
	}))

	var result map[string]string
	err := c.Get(context.Background(), "/test", &result)
	require.Error(t, err)
	apiErr := err.(*APIError)
	assert.Equal(t, "forbidden", apiErr.Message)
}

func TestClient_ErrorResponse_RawBody(t *testing.T) {
	c := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "internal error") //nolint:errcheck
	}))

	var result map[string]string
	err := c.Get(context.Background(), "/test", &result)
	require.Error(t, err)
	apiErr := err.(*APIError)
	assert.Equal(t, 500, apiErr.StatusCode)
	assert.Equal(t, "internal error", apiErr.Message)
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
	require.NoError(t, err)
	assert.Equal(t, int32(2), atomic.LoadInt32(&callCount))
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
	require.NoError(t, err)
	assert.Equal(t, int32(2), atomic.LoadInt32(&callCount))
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

	assert.GreaterOrEqual(t, elapsed, 50*time.Millisecond)
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
	require.NoError(t, json.Unmarshal(gotBody, &parsed))
	assert.Equal(t, "web", parsed.Component.Name)
	assert.Equal(t, "operational", parsed.Component.Status)
}

func TestMaskHeader(t *testing.T) {
	input := "GET / HTTP/1.1\nHost: example.com\nAuthorization: OAuth mySecretToken\nContent-Type: application/json"
	got := maskHeader("Authorization", input)
	assert.Equal(t, "GET / HTTP/1.1\nHost: example.com\nAuthorization: ***** ********Token\nContent-Type: application/json", got)
}

func TestMaskHeader_NoMatch(t *testing.T) {
	input := "GET / HTTP/1.1\nHost: example.com\nContent-Type: application/json"
	got := maskHeader("Authorization", input)
	assert.Equal(t, input, got)
}

func TestMaskHeader_ShortValue(t *testing.T) {
	// Value is short enough that no masking occurs (len <= 5)
	input := "Authorization: abcd"
	got := maskHeader("Authorization", input)
	assert.Equal(t, "Authorization: abcd", got)
}

func TestMaskHeader_LeadingSpacePreserved(t *testing.T) {
	// Leading space after colon is preserved during masking; last 5 chars are visible
	input := "Authorization: mysecretapikey"
	got := maskHeader("Authorization", input)
	assert.Equal(t, "Authorization: *********pikey", got)
}
