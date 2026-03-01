package provider

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/jarcoal/httpmock"
	"github.com/winebarrel/terraform-provider-statuspage/internal/apiclient"
)

func TestAccMetric_basic(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	providerID := "metrics-provider-for-metric-1"
	metricsProvider := &apiclient.MetricsProvider{
		ID:     providerID,
		PageID: testAccPageID,
		Type:   "Self",
	}

	metricID := "metric-id-1"
	metric := &apiclient.Metric{
		ID:                metricID,
		PageID:            testAccPageID,
		MetricsProviderID: providerID,
	}

	// Metrics Provider responders
	httpmock.RegisterResponder("POST", testBaseURL+"/pages/"+testAccPageID+"/metrics_providers",
		func(req *http.Request) (*http.Response, error) {
			var body apiclient.MetricsProviderRequest
			json.NewDecoder(req.Body).Decode(&body)
			metricsProvider.Type = body.MetricsProvider.Type
			return httpmock.NewJsonResponse(201, metricsProvider)
		})

	httpmock.RegisterResponder("GET", testBaseURL+"/pages/"+testAccPageID+"/metrics_providers/"+providerID,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, metricsProvider)
		})

	httpmock.RegisterResponder("DELETE", testBaseURL+"/pages/"+testAccPageID+"/metrics_providers/"+providerID,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(204, ""), nil
		})

	// Metric responders
	httpmock.RegisterResponder("POST", testBaseURL+"/pages/"+testAccPageID+"/metrics_providers/"+providerID+"/metrics",
		func(req *http.Request) (*http.Response, error) {
			var body apiclient.MetricRequest
			json.NewDecoder(req.Body).Decode(&body)
			metric.Name = body.Metric.Name
			metric.Suffix = body.Metric.Suffix
			return httpmock.NewJsonResponse(201, metric)
		})

	httpmock.RegisterResponder("GET", testBaseURL+"/pages/"+testAccPageID+"/metrics/"+metricID,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, metric)
		})

	httpmock.RegisterResponder("PATCH", testBaseURL+"/pages/"+testAccPageID+"/metrics/"+metricID,
		func(req *http.Request) (*http.Response, error) {
			var body apiclient.MetricRequest
			json.NewDecoder(req.Body).Decode(&body)
			if body.Metric.Name != "" {
				metric.Name = body.Metric.Name
			}
			if body.Metric.Suffix != "" {
				metric.Suffix = body.Metric.Suffix
			}
			return httpmock.NewJsonResponse(200, metric)
		})

	httpmock.RegisterResponder("DELETE", testBaseURL+"/pages/"+testAccPageID+"/metrics/"+metricID,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(204, ""), nil
		})

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read (creates a Self metrics provider first, then a metric)
			{
				Config: testAccMetricConfig("tf-test-metric", "ms"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("statuspage_metric.test", "name", "tf-test-metric"),
					resource.TestCheckResourceAttr("statuspage_metric.test", "suffix", "ms"),
					resource.TestCheckResourceAttrSet("statuspage_metric.test", "id"),
					resource.TestCheckResourceAttrSet("statuspage_metric.test", "metrics_provider_id"),
				),
			},
			// Import
			{
				ResourceName:      "statuspage_metric.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: importStateIDFuncMetric("statuspage_metric.test"),
			},
			// Update
			{
				Config: testAccMetricConfig("tf-test-metric-updated", "req/s"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("statuspage_metric.test", "name", "tf-test-metric-updated"),
					resource.TestCheckResourceAttr("statuspage_metric.test", "suffix", "req/s"),
				),
			},
		},
	})
}

func testAccMetricConfig(name, suffix string) string {
	return fmt.Sprintf(`
resource "statuspage_metrics_provider" "metric_provider" {
  page_id = %[1]q
  type    = "Self"
}

resource "statuspage_metric" "test" {
  page_id             = %[1]q
  metrics_provider_id = statuspage_metrics_provider.metric_provider.id
  name                = %[2]q
  suffix              = %[3]q
}
`, testAccPageID, name, suffix)
}

func importStateIDFuncMetric(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource not found: %s", resourceName)
		}
		pageID := rs.Primary.Attributes["page_id"]
		id := rs.Primary.ID
		return fmt.Sprintf("%s/%s", pageID, id), nil
	}
}
