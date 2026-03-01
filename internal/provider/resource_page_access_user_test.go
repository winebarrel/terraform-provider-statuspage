package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccPageAccessUser_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccPageAccessUserConfig("tf-test-pau", "tf-test-pau@example.com"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("statuspage_page_access_user.test", "external_login", "tf-test-pau"),
					resource.TestCheckResourceAttr("statuspage_page_access_user.test", "external_email", "tf-test-pau@example.com"),
					resource.TestCheckResourceAttrSet("statuspage_page_access_user.test", "id"),
				),
			},
			// Import
			{
				ResourceName:      "statuspage_page_access_user.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: importStateIDFuncPageAccessUser("statuspage_page_access_user.test"),
			},
			// Update
			{
				Config: testAccPageAccessUserConfig("tf-test-pau-updated", "tf-test-pau-updated@example.com"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("statuspage_page_access_user.test", "external_login", "tf-test-pau-updated"),
					resource.TestCheckResourceAttr("statuspage_page_access_user.test", "external_email", "tf-test-pau-updated@example.com"),
				),
			},
		},
	})
}

func testAccPageAccessUserConfig(login, email string) string {
	return fmt.Sprintf(`
resource "statuspage_page_access_user" "test" {
  page_id        = %q
  external_login = %q
  external_email = %q
}
`, testAccPageID, login, email)
}

func importStateIDFuncPageAccessUser(resourceName string) resource.ImportStateIdFunc {
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
