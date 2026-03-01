package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccIncident_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccIncidentConfig("tf-test-incident", "investigating", "Initial investigation"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("statuspage_incident.test", "name", "tf-test-incident"),
					resource.TestCheckResourceAttr("statuspage_incident.test", "status", "investigating"),
					resource.TestCheckResourceAttr("statuspage_incident.test", "body", "Initial investigation"),
					resource.TestCheckResourceAttrSet("statuspage_incident.test", "id"),
					resource.TestCheckResourceAttrSet("statuspage_incident.test", "shortlink"),
				),
			},
			// Import
			{
				ResourceName:      "statuspage_incident.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: importStateIDFuncIncident("statuspage_incident.test"),
			},
			// Update
			{
				Config: testAccIncidentConfig("tf-test-incident-updated", "identified", "Issue has been identified"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("statuspage_incident.test", "name", "tf-test-incident-updated"),
					resource.TestCheckResourceAttr("statuspage_incident.test", "status", "identified"),
					resource.TestCheckResourceAttr("statuspage_incident.test", "body", "Issue has been identified"),
				),
			},
		},
	})
}

func testAccIncidentConfig(name, status, body string) string {
	return fmt.Sprintf(`
resource "statuspage_incident" "test" {
  page_id = %q
  name    = %q
  status  = %q
  body    = %q
}
`, testAccPageID, name, status, body)
}

func importStateIDFuncIncident(resourceName string) resource.ImportStateIdFunc {
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
