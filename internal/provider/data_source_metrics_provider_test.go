package provider

import (
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	"github.com/winebarrel/terraform-provider-statuspage/internal/apiclient"
)

func TestMetricsProviderDataSource_basic(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	provider := &apiclient.MetricsProvider{
		ID:            "provider-id-1",
		PageID:        "test-page-id",
		Type:          "Datadog",
		Email:         "metrics@example.com",
		MetricBaseURI: "https://app.datadoghq.com",
	}

	httpmock.RegisterResponder("GET", "https://api.statuspage.io/v1/pages/test-page-id/metrics_providers/provider-id-1",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, provider)
		})

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
data "statuspage_metrics_provider" "test" {
  page_id = "test-page-id"
  id      = "provider-id-1"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.statuspage_metrics_provider.test", "id", "provider-id-1"),
					resource.TestCheckResourceAttr("data.statuspage_metrics_provider.test", "page_id", "test-page-id"),
					resource.TestCheckResourceAttr("data.statuspage_metrics_provider.test", "type", "Datadog"),
					resource.TestCheckResourceAttr("data.statuspage_metrics_provider.test", "email", "metrics@example.com"),
					resource.TestCheckResourceAttr("data.statuspage_metrics_provider.test", "metric_base_uri", "https://app.datadoghq.com"),
				),
			},
		},
	})
}
