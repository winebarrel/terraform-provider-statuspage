package apiclient

import (
	"context"
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

// --- Page ---

func TestGetPage(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"GET /pages/page1": func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, 200, Page{ID: "page1", Name: "My Page", Subdomain: "mypage"})
		},
	})

	page, err := c.GetPage(context.Background(), "page1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if page.ID != "page1" {
		t.Errorf("expected ID %q, got %q", "page1", page.ID)
	}
	if page.Name != "My Page" {
		t.Errorf("expected Name %q, got %q", "My Page", page.Name)
	}
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
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if page.Name != "Updated Page" {
		t.Errorf("expected Name %q, got %q", "Updated Page", page.Name)
	}
}

// --- Component ---

func TestGetComponent(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"GET /pages/p1/components/c1": func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, 200, Component{ID: "c1", PageID: "p1", Name: "API", Status: "operational"})
		},
	})

	comp, err := c.GetComponent(context.Background(), "p1", "c1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if comp.ID != "c1" {
		t.Errorf("expected ID %q, got %q", "c1", comp.ID)
	}
	if comp.Name != "API" {
		t.Errorf("expected Name %q, got %q", "API", comp.Name)
	}
}

func TestCreateComponent(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"POST /pages/p1/components": func(w http.ResponseWriter, r *http.Request) {
			var req ComponentRequest
			json.NewDecoder(r.Body).Decode(&req) //nolint:errcheck
			jsonResponse(w, 201, Component{ID: "c-new", PageID: "p1", Name: req.Component.Name, Status: req.Component.Status})
		},
	})

	comp, err := c.CreateComponent(context.Background(), "p1", ComponentBody{Name: "Web", Status: "operational"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if comp.ID != "c-new" {
		t.Errorf("expected ID %q, got %q", "c-new", comp.ID)
	}
	if comp.Name != "Web" {
		t.Errorf("expected Name %q, got %q", "Web", comp.Name)
	}
}

func TestUpdateComponent(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"PATCH /pages/p1/components/c1": func(w http.ResponseWriter, r *http.Request) {
			var req ComponentRequest
			json.NewDecoder(r.Body).Decode(&req) //nolint:errcheck
			jsonResponse(w, 200, Component{ID: "c1", PageID: "p1", Name: req.Component.Name, Status: "degraded_performance"})
		},
	})

	comp, err := c.UpdateComponent(context.Background(), "p1", "c1", ComponentBody{Name: "API Updated"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if comp.Name != "API Updated" {
		t.Errorf("expected Name %q, got %q", "API Updated", comp.Name)
	}
}

func TestDeleteComponent(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"DELETE /pages/p1/components/c1": func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(204)
		},
	})

	err := c.DeleteComponent(context.Background(), "p1", "c1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetComponent_Error(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"GET /pages/p1/components/bad": func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, 404, map[string]string{"error": "not found"})
		},
	})

	_, err := c.GetComponent(context.Background(), "p1", "bad")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != 404 {
		t.Errorf("expected status 404, got %d", apiErr.StatusCode)
	}
}

// --- ComponentGroup ---

func TestGetComponentGroup(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"GET /pages/p1/component-groups/g1": func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, 200, ComponentGroup{ID: "g1", PageID: "p1", Name: "Infrastructure"})
		},
	})

	group, err := c.GetComponentGroup(context.Background(), "p1", "g1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if group.ID != "g1" {
		t.Errorf("expected ID %q, got %q", "g1", group.ID)
	}
	if group.Name != "Infrastructure" {
		t.Errorf("expected Name %q, got %q", "Infrastructure", group.Name)
	}
}

func TestCreateComponentGroup(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"POST /pages/p1/component-groups": func(w http.ResponseWriter, r *http.Request) {
			var req ComponentGroupRequest
			json.NewDecoder(r.Body).Decode(&req) //nolint:errcheck
			jsonResponse(w, 201, ComponentGroup{ID: "g-new", PageID: "p1", Name: req.ComponentGroup.Name, Components: req.ComponentGroup.Components})
		},
	})

	group, err := c.CreateComponentGroup(context.Background(), "p1", ComponentGroupBody{Name: "Backend", Components: []string{"c1", "c2"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if group.Name != "Backend" {
		t.Errorf("expected Name %q, got %q", "Backend", group.Name)
	}
	if len(group.Components) != 2 {
		t.Errorf("expected 2 components, got %d", len(group.Components))
	}
}

func TestUpdateComponentGroup(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"PATCH /pages/p1/component-groups/g1": func(w http.ResponseWriter, r *http.Request) {
			var req ComponentGroupRequest
			json.NewDecoder(r.Body).Decode(&req) //nolint:errcheck
			jsonResponse(w, 200, ComponentGroup{ID: "g1", PageID: "p1", Name: req.ComponentGroup.Name})
		},
	})

	group, err := c.UpdateComponentGroup(context.Background(), "p1", "g1", ComponentGroupBody{Name: "Updated Group"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if group.Name != "Updated Group" {
		t.Errorf("expected Name %q, got %q", "Updated Group", group.Name)
	}
}

func TestDeleteComponentGroup(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"DELETE /pages/p1/component-groups/g1": func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(204)
		},
	})

	err := c.DeleteComponentGroup(context.Background(), "p1", "g1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// --- Incident ---

func TestGetIncident(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"GET /pages/p1/incidents/i1": func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, 200, Incident{ID: "i1", PageID: "p1", Name: "Major Outage", Status: "investigating"})
		},
	})

	incident, err := c.GetIncident(context.Background(), "p1", "i1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if incident.Name != "Major Outage" {
		t.Errorf("expected Name %q, got %q", "Major Outage", incident.Name)
	}
	if incident.Status != "investigating" {
		t.Errorf("expected Status %q, got %q", "investigating", incident.Status)
	}
}

func TestCreateIncident(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"POST /pages/p1/incidents": func(w http.ResponseWriter, r *http.Request) {
			var req IncidentRequest
			json.NewDecoder(r.Body).Decode(&req) //nolint:errcheck
			jsonResponse(w, 201, Incident{ID: "i-new", PageID: "p1", Name: req.Incident.Name, Status: req.Incident.Status})
		},
	})

	incident, err := c.CreateIncident(context.Background(), "p1", IncidentBody{Name: "API Down", Status: "investigating"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if incident.ID != "i-new" {
		t.Errorf("expected ID %q, got %q", "i-new", incident.ID)
	}
}

func TestUpdateIncident(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"PATCH /pages/p1/incidents/i1": func(w http.ResponseWriter, r *http.Request) {
			var req IncidentRequest
			json.NewDecoder(r.Body).Decode(&req) //nolint:errcheck
			jsonResponse(w, 200, Incident{ID: "i1", PageID: "p1", Name: "API Down", Status: req.Incident.Status})
		},
	})

	incident, err := c.UpdateIncident(context.Background(), "p1", "i1", IncidentBody{Status: "resolved"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if incident.Status != "resolved" {
		t.Errorf("expected Status %q, got %q", "resolved", incident.Status)
	}
}

func TestDeleteIncident(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"DELETE /pages/p1/incidents/i1": func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(204)
		},
	})

	err := c.DeleteIncident(context.Background(), "p1", "i1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// --- IncidentTemplate ---

func TestGetIncidentTemplate(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"GET /pages/p1/incident_templates": func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, 200, []IncidentTemplate{
				{ID: "t1", Name: "Template A"},
				{ID: "t2", Name: "Template B"},
			})
		},
	})

	tmpl, err := c.GetIncidentTemplate(context.Background(), "p1", "t2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tmpl.ID != "t2" {
		t.Errorf("expected ID %q, got %q", "t2", tmpl.ID)
	}
	if tmpl.Name != "Template B" {
		t.Errorf("expected Name %q, got %q", "Template B", tmpl.Name)
	}
}

func TestGetIncidentTemplate_NotFound(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"GET /pages/p1/incident_templates": func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, 200, []IncidentTemplate{
				{ID: "t1", Name: "Template A"},
			})
		},
	})

	_, err := c.GetIncidentTemplate(context.Background(), "p1", "nonexistent")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != 404 {
		t.Errorf("expected status 404, got %d", apiErr.StatusCode)
	}
}

func TestCreateIncidentTemplate(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"POST /pages/p1/incident_templates": func(w http.ResponseWriter, r *http.Request) {
			var req IncidentTemplateRequest
			json.NewDecoder(r.Body).Decode(&req) //nolint:errcheck
			jsonResponse(w, 201, IncidentTemplate{ID: "t-new", Name: req.Template.Name, Title: req.Template.Title})
		},
	})

	tmpl, err := c.CreateIncidentTemplate(context.Background(), "p1", IncidentTemplateBody{Name: "Outage Template", Title: "Service Outage"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tmpl.ID != "t-new" {
		t.Errorf("expected ID %q, got %q", "t-new", tmpl.ID)
	}
	if tmpl.Name != "Outage Template" {
		t.Errorf("expected Name %q, got %q", "Outage Template", tmpl.Name)
	}
}

func TestUpdateIncidentTemplate(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"PATCH /pages/p1/incident_templates/t1": func(w http.ResponseWriter, r *http.Request) {
			var req IncidentTemplateRequest
			json.NewDecoder(r.Body).Decode(&req) //nolint:errcheck
			jsonResponse(w, 200, IncidentTemplate{ID: "t1", Name: req.Template.Name})
		},
	})

	tmpl, err := c.UpdateIncidentTemplate(context.Background(), "p1", "t1", IncidentTemplateBody{Name: "Updated Template"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tmpl.Name != "Updated Template" {
		t.Errorf("expected Name %q, got %q", "Updated Template", tmpl.Name)
	}
}

func TestDeleteIncidentTemplate(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"DELETE /pages/p1/incident_templates/t1": func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(204)
		},
	})

	err := c.DeleteIncidentTemplate(context.Background(), "p1", "t1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// --- Postmortem ---

func TestGetPostmortem(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"GET /pages/p1/incidents/i1/postmortem": func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, 200, Postmortem{Body: "Root cause analysis", BodyDraft: "Draft content"})
		},
	})

	pm, err := c.GetPostmortem(context.Background(), "p1", "i1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pm.Body != "Root cause analysis" {
		t.Errorf("expected Body %q, got %q", "Root cause analysis", pm.Body)
	}
}

func TestCreateOrUpdatePostmortem(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"PUT /pages/p1/incidents/i1/postmortem": func(w http.ResponseWriter, r *http.Request) {
			var req PostmortemRequest
			json.NewDecoder(r.Body).Decode(&req) //nolint:errcheck
			jsonResponse(w, 200, Postmortem{Body: req.Postmortem.Body})
		},
	})

	pm, err := c.CreateOrUpdatePostmortem(context.Background(), "p1", "i1", PostmortemBody{Body: "Updated analysis"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pm.Body != "Updated analysis" {
		t.Errorf("expected Body %q, got %q", "Updated analysis", pm.Body)
	}
}

func TestDeletePostmortem(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"DELETE /pages/p1/incidents/i1/postmortem": func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(204)
		},
	})

	err := c.DeletePostmortem(context.Background(), "p1", "i1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// --- Metric ---

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

// --- MetricsProvider ---

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

// --- Subscriber ---

func TestGetSubscriber(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"GET /pages/p1/subscribers/s1": func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, 200, Subscriber{ID: "s1", Email: "user@example.com", Mode: "email"})
		},
	})

	sub, err := c.GetSubscriber(context.Background(), "p1", "s1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sub.Email != "user@example.com" {
		t.Errorf("expected Email %q, got %q", "user@example.com", sub.Email)
	}
}

func TestCreateSubscriber(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"POST /pages/p1/subscribers": func(w http.ResponseWriter, r *http.Request) {
			var req SubscriberRequest
			json.NewDecoder(r.Body).Decode(&req) //nolint:errcheck
			jsonResponse(w, 201, Subscriber{ID: "s-new", Email: req.Subscriber.Email})
		},
	})

	sub, err := c.CreateSubscriber(context.Background(), "p1", SubscriberBody{Email: "new@example.com"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sub.Email != "new@example.com" {
		t.Errorf("expected Email %q, got %q", "new@example.com", sub.Email)
	}
}

func TestUpdateSubscriber(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"PATCH /pages/p1/subscribers/s1": func(w http.ResponseWriter, r *http.Request) {
			var req SubscriberRequest
			json.NewDecoder(r.Body).Decode(&req) //nolint:errcheck
			jsonResponse(w, 200, Subscriber{ID: "s1", Email: req.Subscriber.Email})
		},
	})

	sub, err := c.UpdateSubscriber(context.Background(), "p1", "s1", SubscriberBody{Email: "updated@example.com"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sub.Email != "updated@example.com" {
		t.Errorf("expected Email %q, got %q", "updated@example.com", sub.Email)
	}
}

func TestDeleteSubscriber(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"DELETE /pages/p1/subscribers/s1": func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(204)
		},
	})

	err := c.DeleteSubscriber(context.Background(), "p1", "s1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// --- PageAccessGroup ---

func TestGetPageAccessGroup(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"GET /pages/p1/page_access_groups/ag1": func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, 200, PageAccessGroup{ID: "ag1", Name: "VIP Group"})
		},
	})

	group, err := c.GetPageAccessGroup(context.Background(), "p1", "ag1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if group.Name != "VIP Group" {
		t.Errorf("expected Name %q, got %q", "VIP Group", group.Name)
	}
}

func TestCreatePageAccessGroup(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"POST /pages/p1/page_access_groups": func(w http.ResponseWriter, r *http.Request) {
			var req PageAccessGroupRequest
			json.NewDecoder(r.Body).Decode(&req) //nolint:errcheck
			jsonResponse(w, 201, PageAccessGroup{ID: "ag-new", Name: req.PageAccessGroup.Name})
		},
	})

	group, err := c.CreatePageAccessGroup(context.Background(), "p1", PageAccessGroupBody{Name: "New Group"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if group.Name != "New Group" {
		t.Errorf("expected Name %q, got %q", "New Group", group.Name)
	}
}

func TestUpdatePageAccessGroup(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"PATCH /pages/p1/page_access_groups/ag1": func(w http.ResponseWriter, r *http.Request) {
			var req PageAccessGroupRequest
			json.NewDecoder(r.Body).Decode(&req) //nolint:errcheck
			jsonResponse(w, 200, PageAccessGroup{ID: "ag1", Name: req.PageAccessGroup.Name})
		},
	})

	group, err := c.UpdatePageAccessGroup(context.Background(), "p1", "ag1", PageAccessGroupBody{Name: "Updated Group"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if group.Name != "Updated Group" {
		t.Errorf("expected Name %q, got %q", "Updated Group", group.Name)
	}
}

func TestDeletePageAccessGroup(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"DELETE /pages/p1/page_access_groups/ag1": func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(204)
		},
	})

	err := c.DeletePageAccessGroup(context.Background(), "p1", "ag1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// --- PageAccessUser ---

func TestGetPageAccessUser(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"GET /pages/p1/page_access_users/au1": func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, 200, PageAccessUser{ID: "au1", ExternalEmail: "user@example.com"})
		},
	})

	user, err := c.GetPageAccessUser(context.Background(), "p1", "au1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.ExternalEmail != "user@example.com" {
		t.Errorf("expected ExternalEmail %q, got %q", "user@example.com", user.ExternalEmail)
	}
}

func TestCreatePageAccessUser(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"POST /pages/p1/page_access_users": func(w http.ResponseWriter, r *http.Request) {
			var req PageAccessUserRequest
			json.NewDecoder(r.Body).Decode(&req) //nolint:errcheck
			jsonResponse(w, 201, PageAccessUser{ID: "au-new", ExternalEmail: req.PageAccessUser.ExternalEmail})
		},
	})

	user, err := c.CreatePageAccessUser(context.Background(), "p1", PageAccessUserBody{ExternalEmail: "new@example.com"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.ExternalEmail != "new@example.com" {
		t.Errorf("expected ExternalEmail %q, got %q", "new@example.com", user.ExternalEmail)
	}
}

func TestUpdatePageAccessUser(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"PATCH /pages/p1/page_access_users/au1": func(w http.ResponseWriter, r *http.Request) {
			var req PageAccessUserRequest
			json.NewDecoder(r.Body).Decode(&req) //nolint:errcheck
			jsonResponse(w, 200, PageAccessUser{ID: "au1", ExternalEmail: req.PageAccessUser.ExternalEmail})
		},
	})

	user, err := c.UpdatePageAccessUser(context.Background(), "p1", "au1", PageAccessUserBody{ExternalEmail: "updated@example.com"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.ExternalEmail != "updated@example.com" {
		t.Errorf("expected ExternalEmail %q, got %q", "updated@example.com", user.ExternalEmail)
	}
}

func TestDeletePageAccessUser(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"DELETE /pages/p1/page_access_users/au1": func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(204)
		},
	})

	err := c.DeletePageAccessUser(context.Background(), "p1", "au1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// --- StatusEmbedConfig ---

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

// --- User ---

func TestGetUser(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"GET /organizations/org1/users": func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, 200, []User{
				{ID: "u1", Email: "alice@example.com", FirstName: "Alice"},
				{ID: "u2", Email: "bob@example.com", FirstName: "Bob"},
			})
		},
	})

	user, err := c.GetUser(context.Background(), "org1", "u2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.ID != "u2" {
		t.Errorf("expected ID %q, got %q", "u2", user.ID)
	}
	if user.FirstName != "Bob" {
		t.Errorf("expected FirstName %q, got %q", "Bob", user.FirstName)
	}
}

func TestGetUser_NotFound(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"GET /organizations/org1/users": func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, 200, []User{
				{ID: "u1", Email: "alice@example.com"},
			})
		},
	})

	_, err := c.GetUser(context.Background(), "org1", "nonexistent")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != 404 {
		t.Errorf("expected status 404, got %d", apiErr.StatusCode)
	}
}

func TestCreateUser(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"POST /organizations/org1/users": func(w http.ResponseWriter, r *http.Request) {
			var req UserRequest
			json.NewDecoder(r.Body).Decode(&req) //nolint:errcheck
			jsonResponse(w, 201, User{ID: "u-new", Email: req.User.Email, FirstName: req.User.FirstName})
		},
	})

	user, err := c.CreateUser(context.Background(), "org1", UserBody{Email: "new@example.com", FirstName: "New"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.Email != "new@example.com" {
		t.Errorf("expected Email %q, got %q", "new@example.com", user.Email)
	}
}

func TestDeleteUser(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"DELETE /organizations/org1/users/u1": func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(204)
		},
	})

	err := c.DeleteUser(context.Background(), "org1", "u1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// --- Permissions ---

func TestGetPermissions(t *testing.T) {
	c := testServer(t, map[string]http.HandlerFunc{
		"GET /organizations/org1/permissions/u1": func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, 200, Permissions{UserID: "u1", Pages: map[string]string{"p1": "admin", "p2": "viewer"}})
		},
	})

	perms, err := c.GetPermissions(context.Background(), "org1", "u1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if perms.UserID != "u1" {
		t.Errorf("expected UserID %q, got %q", "u1", perms.UserID)
	}
	if perms.Pages["p1"] != "admin" {
		t.Errorf("expected page p1 role %q, got %q", "admin", perms.Pages["p1"])
	}
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
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if perms.Pages["p1"] != "manager" {
		t.Errorf("expected page p1 role %q, got %q", "manager", perms.Pages["p1"])
	}
}
