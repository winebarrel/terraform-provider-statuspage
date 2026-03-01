package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccComponent_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccComponentConfig("tf-test-component", "operational"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("statuspage_component.test", "name", "tf-test-component"),
					resource.TestCheckResourceAttr("statuspage_component.test", "status", "operational"),
					resource.TestCheckResourceAttrSet("statuspage_component.test", "id"),
					resource.TestCheckResourceAttrSet("statuspage_component.test", "automation_email"),
				),
			},
			// Import
			{
				ResourceName:      "statuspage_component.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: importStateIDFunc("statuspage_component.test"),
			},
			// Update
			{
				Config: testAccComponentConfig("tf-test-component-updated", "degraded_performance"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("statuspage_component.test", "name", "tf-test-component-updated"),
					resource.TestCheckResourceAttr("statuspage_component.test", "status", "degraded_performance"),
				),
			},
		},
	})
}

func testAccComponentConfig(name, status string) string {
	return fmt.Sprintf(`
resource "statuspage_component" "test" {
  page_id = %q
  name    = %q
  status  = %q
}
`, testAccPageID, name, status)
}

func importStateIDFunc(resourceName string) resource.ImportStateIdFunc {
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
