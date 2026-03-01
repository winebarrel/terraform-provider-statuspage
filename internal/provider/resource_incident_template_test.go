package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccIncidentTemplate_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccIncidentTemplateConfig("tf-test-template", "Test Incident", "investigating"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("statuspage_incident_template.test", "name", "tf-test-template"),
					resource.TestCheckResourceAttr("statuspage_incident_template.test", "title", "Test Incident"),
					resource.TestCheckResourceAttr("statuspage_incident_template.test", "update_status", "investigating"),
					resource.TestCheckResourceAttrSet("statuspage_incident_template.test", "id"),
				),
			},
			// Import
			{
				ResourceName:      "statuspage_incident_template.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: importStateIDFuncIncidentTemplate("statuspage_incident_template.test"),
			},
			// Update
			{
				Config: testAccIncidentTemplateConfig("tf-test-template-updated", "Updated Incident", "identified"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("statuspage_incident_template.test", "name", "tf-test-template-updated"),
					resource.TestCheckResourceAttr("statuspage_incident_template.test", "title", "Updated Incident"),
					resource.TestCheckResourceAttr("statuspage_incident_template.test", "update_status", "identified"),
				),
			},
		},
	})
}

func testAccIncidentTemplateConfig(name, title, updateStatus string) string {
	return fmt.Sprintf(`
resource "statuspage_incident_template" "test" {
  page_id       = %q
  name          = %q
  title         = %q
  body          = "This is a template body."
  update_status = %q
}
`, testAccPageID, name, title, updateStatus)
}

func importStateIDFuncIncidentTemplate(resourceName string) resource.ImportStateIdFunc {
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
