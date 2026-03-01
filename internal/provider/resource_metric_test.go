package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccMetric_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
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
