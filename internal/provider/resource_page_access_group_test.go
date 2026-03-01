package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccPageAccessGroup_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccPageAccessGroupConfig("tf-test-pag"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("statuspage_page_access_group.test", "name", "tf-test-pag"),
					resource.TestCheckResourceAttrSet("statuspage_page_access_group.test", "id"),
				),
			},
			// Import
			{
				ResourceName:      "statuspage_page_access_group.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: importStateIDFuncPageAccessGroup("statuspage_page_access_group.test"),
			},
			// Update
			{
				Config: testAccPageAccessGroupConfig("tf-test-pag-updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("statuspage_page_access_group.test", "name", "tf-test-pag-updated"),
				),
			},
		},
	})
}

func testAccPageAccessGroupConfig(name string) string {
	return fmt.Sprintf(`
resource "statuspage_page_access_group" "test" {
  page_id = %q
  name    = %q
}
`, testAccPageID, name)
}

func importStateIDFuncPageAccessGroup(resourceName string) resource.ImportStateIdFunc {
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
