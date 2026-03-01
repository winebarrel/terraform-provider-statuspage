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

func TestAccMetricsProvider_basic(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	providerID := "metrics-provider-id-1"
	metricsProvider := &apiclient.MetricsProvider{
		ID:     providerID,
		PageID: "test-page-id",
	}

	httpmock.RegisterResponder("POST", "https://api.statuspage.io/v1/pages/test-page-id/metrics_providers",
		func(req *http.Request) (*http.Response, error) {
			var body apiclient.MetricsProviderRequest
			json.NewDecoder(req.Body).Decode(&body)
			metricsProvider.Type = body.MetricsProvider.Type
			metricsProvider.Email = body.MetricsProvider.Email
			return httpmock.NewJsonResponse(201, metricsProvider)
		})

	httpmock.RegisterResponder("GET", "https://api.statuspage.io/v1/pages/test-page-id/metrics_providers/"+providerID,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, metricsProvider)
		})

	httpmock.RegisterResponder("PATCH", "https://api.statuspage.io/v1/pages/test-page-id/metrics_providers/"+providerID,
		func(req *http.Request) (*http.Response, error) {
			var body apiclient.MetricsProviderRequest
			json.NewDecoder(req.Body).Decode(&body)
			if body.MetricsProvider.Type != "" {
				metricsProvider.Type = body.MetricsProvider.Type
			}
			if body.MetricsProvider.Email != "" {
				metricsProvider.Email = body.MetricsProvider.Email
			}
			return httpmock.NewJsonResponse(200, metricsProvider)
		})

	httpmock.RegisterResponder("DELETE", "https://api.statuspage.io/v1/pages/test-page-id/metrics_providers/"+providerID,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(204, ""), nil
		})

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccMetricsProviderConfig("Self"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("statuspage_metrics_provider.test", "type", "Self"),
					resource.TestCheckResourceAttrSet("statuspage_metrics_provider.test", "id"),
				),
			},
			// Import
			{
				ResourceName:      "statuspage_metrics_provider.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: importStateIDFuncMetricsProvider("statuspage_metrics_provider.test"),
				// Sensitive fields (api_key, api_token, application_key) are not returned by the API
				ImportStateVerifyIgnore: []string{"api_key", "api_token", "application_key"},
			},
			// Update (change type to Self - effectively a no-op re-apply since type is the same)
			{
				Config: testAccMetricsProviderConfigWithEmail("Self", "tf-test@example.com"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("statuspage_metrics_provider.test", "type", "Self"),
					resource.TestCheckResourceAttr("statuspage_metrics_provider.test", "email", "tf-test@example.com"),
				),
			},
		},
	})
}

func testAccMetricsProviderConfig(providerType string) string {
	return fmt.Sprintf(`
resource "statuspage_metrics_provider" "test" {
  page_id = %q
  type    = %q
}
`, "test-page-id", providerType)
}

func testAccMetricsProviderConfigWithEmail(providerType, email string) string {
	return fmt.Sprintf(`
resource "statuspage_metrics_provider" "test" {
  page_id = %q
  type    = %q
  email   = %q
}
`, "test-page-id", providerType, email)
}

func importStateIDFuncMetricsProvider(resourceName string) resource.ImportStateIdFunc {
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
