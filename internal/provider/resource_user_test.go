package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccUser_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck: func() {
			testAccPreCheck(t)
			if testAccOrganizationID == "" {
				t.Skip("STATUSPAGE_ORGANIZATION_ID must be set for user acceptance tests")
			}
		},
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccUserConfig("tf-test-user@example.com", "TestFirst", "TestLast"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("statuspage_user.test", "email", "tf-test-user@example.com"),
					resource.TestCheckResourceAttr("statuspage_user.test", "first_name", "TestFirst"),
					resource.TestCheckResourceAttr("statuspage_user.test", "last_name", "TestLast"),
					resource.TestCheckResourceAttrSet("statuspage_user.test", "id"),
				),
			},
			// Import
			{
				ResourceName:      "statuspage_user.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: importStateIDFuncUser("statuspage_user.test"),
				// password is sensitive and not returned by the API
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func testAccUserConfig(email, firstName, lastName string) string {
	return fmt.Sprintf(`
resource "statuspage_user" "test" {
  organization_id = %q
  email           = %q
  password        = "TestP@ssw0rd123!"
  first_name      = %q
  last_name       = %q
}
`, testAccOrganizationID, email, firstName, lastName)
}

func importStateIDFuncUser(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource not found: %s", resourceName)
		}
		orgID := rs.Primary.Attributes["organization_id"]
		id := rs.Primary.ID
		return fmt.Sprintf("%s/%s", orgID, id), nil
	}
}
