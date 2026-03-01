package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccComponentGroup_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccComponentGroupConfig("tf-test-group", "tf-test-group-comp1", "tf-test-group-comp2"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("statuspage_component_group.test", "name", "tf-test-group"),
					resource.TestCheckResourceAttr("statuspage_component_group.test", "components.#", "2"),
					resource.TestCheckResourceAttrSet("statuspage_component_group.test", "id"),
				),
			},
			// Import
			{
				ResourceName:      "statuspage_component_group.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: importStateIDFuncComponentGroup("statuspage_component_group.test"),
			},
			// Update
			{
				Config: testAccComponentGroupConfig("tf-test-group-updated", "tf-test-group-comp1", "tf-test-group-comp2"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("statuspage_component_group.test", "name", "tf-test-group-updated"),
				),
			},
		},
	})
}

func testAccComponentGroupConfig(groupName, comp1Name, comp2Name string) string {
	return fmt.Sprintf(`
resource "statuspage_component" "group_comp1" {
  page_id = %[1]q
  name    = %[2]q
  status  = "operational"
}

resource "statuspage_component" "group_comp2" {
  page_id = %[1]q
  name    = %[3]q
  status  = "operational"
}

resource "statuspage_component_group" "test" {
  page_id    = %[1]q
  name       = %[4]q
  components = [statuspage_component.group_comp1.id, statuspage_component.group_comp2.id]
}
`, testAccPageID, comp1Name, comp2Name, groupName)
}

func importStateIDFuncComponentGroup(resourceName string) resource.ImportStateIdFunc {
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
