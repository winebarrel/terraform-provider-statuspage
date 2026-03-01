package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccIncidentPostmortem_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			// Create and Read (creates an incident first, then attaches postmortem)
			{
				Config: testAccIncidentPostmortemConfig("Initial postmortem body"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("statuspage_incident_postmortem.test", "body", "Initial postmortem body"),
					resource.TestCheckResourceAttrSet("statuspage_incident_postmortem.test", "incident_id"),
				),
			},
			// Import
			{
				ResourceName:      "statuspage_incident_postmortem.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: importStateIDFuncPostmortem("statuspage_incident_postmortem.test"),
			},
			// Update
			{
				Config: testAccIncidentPostmortemConfig("Updated postmortem body"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("statuspage_incident_postmortem.test", "body", "Updated postmortem body"),
				),
			},
		},
	})
}

func testAccIncidentPostmortemConfig(body string) string {
	return fmt.Sprintf(`
resource "statuspage_incident" "postmortem_incident" {
  page_id = %[1]q
  name    = "tf-test-postmortem-incident"
  status  = "resolved"
  body    = "Incident for postmortem testing."
}

resource "statuspage_incident_postmortem" "test" {
  page_id            = %[1]q
  incident_id        = statuspage_incident.postmortem_incident.id
  body               = %[2]q
  notify_subscribers = false
  notify_twitter     = false
}
`, testAccPageID, body)
}

func importStateIDFuncPostmortem(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource not found: %s", resourceName)
		}
		pageID := rs.Primary.Attributes["page_id"]
		incidentID := rs.Primary.Attributes["incident_id"]
		return fmt.Sprintf("%s/%s", pageID, incidentID), nil
	}
}
