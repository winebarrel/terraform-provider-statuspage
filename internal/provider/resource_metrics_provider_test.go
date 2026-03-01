package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccMetricsProvider_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
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
`, testAccPageID, providerType)
}

func testAccMetricsProviderConfigWithEmail(providerType, email string) string {
	return fmt.Sprintf(`
resource "statuspage_metrics_provider" "test" {
  page_id = %q
  type    = %q
  email   = %q
}
`, testAccPageID, providerType, email)
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
