package provider

import (
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	"github.com/winebarrel/terraform-provider-statuspage/internal/apiclient"
)

func TestMetricDataSource_basic(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	metric := &apiclient.Metric{
		ID:                 "metric-id-1",
		PageID:             "test-page-id",
		MetricsProviderID:  "provider-id-1",
		MetricIdentifier:   "test-metric",
		Name:               "Test Metric",
		Display:            true,
		TooltipDescription: "A test metric",
		YAxisMin:           0,
		YAxisMax:           100,
		YAxisHidden:        false,
		Suffix:             "ms",
		DecimalPlaces:      2,
	}

	httpmock.RegisterResponder("GET", "https://api.statuspage.io/v1/pages/test-page-id/metrics/metric-id-1",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, metric)
		})

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
data "statuspage_metric" "test" {
  page_id = "test-page-id"
  id      = "metric-id-1"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.statuspage_metric.test", "id", "metric-id-1"),
					resource.TestCheckResourceAttr("data.statuspage_metric.test", "page_id", "test-page-id"),
					resource.TestCheckResourceAttr("data.statuspage_metric.test", "metrics_provider_id", "provider-id-1"),
					resource.TestCheckResourceAttr("data.statuspage_metric.test", "metric_identifier", "test-metric"),
					resource.TestCheckResourceAttr("data.statuspage_metric.test", "name", "Test Metric"),
					resource.TestCheckResourceAttr("data.statuspage_metric.test", "display", "true"),
					resource.TestCheckResourceAttr("data.statuspage_metric.test", "tooltip_description", "A test metric"),
					resource.TestCheckResourceAttr("data.statuspage_metric.test", "y_axis_max", "100"),
					resource.TestCheckResourceAttr("data.statuspage_metric.test", "suffix", "ms"),
					resource.TestCheckResourceAttr("data.statuspage_metric.test", "decimal_places", "2"),
				),
			},
		},
	})
}
